//go:build legacy
// +build legacy

package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"pos-saas/internal/repository"
)

// CollaborationSessionManager manages active collaboration sessions
type CollaborationSessionManager struct {
	sessionRepository  repository.CollaborationSessionRepository
	activityRepository repository.ActivityEventRepository
	commentRepository  repository.CommentRepository
	configRepository   repository.CollaborationConfigRepository

	sessions      map[string]*CollaborationSession // sessionID -> session
	subscribers   map[string]map[string]chan interface{} // sessionID -> userID -> channel
	cursorUpdates map[string]chan *models.CursorPosition // sessionID -> cursor updates
	mu            sync.RWMutex
}

// CollaborationSession represents an active collaboration session
type CollaborationSession struct {
	ID              string
	ResourceID      int64
	ResourceType    string
	RestaurantID    int64
	Participants    map[int64]*Participant // userID -> participant
	Document        string
	VersionNumber   int64
	CreatedAt       time.Time
	LastActivityAt  time.Time
	mu              sync.RWMutex
}

// Participant represents a collaborator
type Participant struct {
	UserID       int64
	Username     string
	Color        string
	CursorPos    int32
	CursorLine   int32
	LastActivity time.Time
	IsActive     bool
}

// NewCollaborationSessionManager creates a new session manager
func NewCollaborationSessionManager(
	sessionRepository repository.CollaborationSessionRepository,
	activityRepository repository.ActivityEventRepository,
	commentRepository repository.CommentRepository,
	configRepository repository.CollaborationConfigRepository,
) *CollaborationSessionManager {
	csm := &CollaborationSessionManager{
		sessionRepository:  sessionRepository,
		activityRepository: activityRepository,
		commentRepository:  commentRepository,
		configRepository:   configRepository,
		sessions:           make(map[string]*CollaborationSession),
		subscribers:        make(map[string]map[string]chan interface{}),
		cursorUpdates:      make(map[string]chan *models.CursorPosition),
	}

	// Start cleanup goroutine
	go csm.cleanupInactiveSessions()

	return csm
}

// CreateSession creates a new collaboration session
func (csm *CollaborationSessionManager) CreateSession(
	ctx context.Context,
	restaurantID int64,
	resourceID int64,
	resourceType string,
	initialDocument string,
) (string, error) {
	sessionID := uuid.New().String()

	session := &CollaborationSession{
		ID:             sessionID,
		ResourceID:     resourceID,
		ResourceType:   resourceType,
		RestaurantID:   restaurantID,
		Participants:   make(map[int64]*Participant),
		Document:       initialDocument,
		VersionNumber:  1,
		CreatedAt:      time.Now(),
		LastActivityAt: time.Now(),
	}

	csm.mu.Lock()
	csm.sessions[sessionID] = session
	csm.subscribers[sessionID] = make(map[string]chan interface{})
	csm.cursorUpdates[sessionID] = make(chan *models.CursorPosition, 100)
	csm.mu.Unlock()

	// Save to database
	dbSession := &models.CollaborationSession{
		RestaurantID:   restaurantID,
		ResourceID:     resourceID,
		ResourceType:   resourceType,
		SessionID:      sessionID,
		IsActive:       true,
		ParticipantCount: 0,
		Participants:   make(map[string]*models.Participant),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	if err := csm.sessionRepository.Create(ctx, dbSession); err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return sessionID, nil
}

// JoinSession adds a user to a collaboration session
func (csm *CollaborationSessionManager) JoinSession(
	ctx context.Context,
	sessionID string,
	userID int64,
	username string,
	color string,
) (chan interface{}, error) {
	csm.mu.Lock()
	defer csm.mu.Unlock()

	session, ok := csm.sessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session not found")
	}

	session.mu.Lock()
	defer session.mu.Unlock()

	// Add participant
	participant := &Participant{
		UserID:       userID,
		Username:     username,
		Color:        color,
		CursorPos:    0,
		LastActivity: time.Now(),
		IsActive:     true,
	}

	session.Participants[userID] = participant

	// Create subscription channel
	if csm.subscribers[sessionID] == nil {
		csm.subscribers[sessionID] = make(map[string]chan interface{})
	}

	userKey := fmt.Sprintf("user_%d", userID)
	channel := make(chan interface{}, 50)
	csm.subscribers[sessionID][userKey] = channel

	// Record activity
	go func() {
		csm.recordActivity(
			context.Background(),
			sessionID,
			session.ResourceID,
			session.ResourceType,
			userID,
			username,
			"join",
			map[string]interface{}{"participant_count": len(session.Participants)},
		)
	}()

	// Broadcast join event to other participants
	csm.broadcastToSession(sessionID, map[string]interface{}{
		"type":      "user_joined",
		"user_id":   userID,
		"username":  username,
		"color":     color,
		"timestamp": time.Now(),
	}, userID)

	return channel, nil
}

// LeaveSession removes a user from a collaboration session
func (csm *CollaborationSessionManager) LeaveSession(
	ctx context.Context,
	sessionID string,
	userID int64,
) error {
	csm.mu.Lock()
	defer csm.mu.Unlock()

	session, ok := csm.sessions[sessionID]
	if !ok {
		return fmt.Errorf("session not found")
	}

	session.mu.Lock()
	participant, exists := session.Participants[userID]
	if !exists {
		session.mu.Unlock()
		return fmt.Errorf("participant not found")
	}

	username := participant.Username
	delete(session.Participants, userID)
	session.mu.Unlock()

	// Close subscription channel
	userKey := fmt.Sprintf("user_%d", userID)
	if ch, ok := csm.subscribers[sessionID][userKey]; ok {
		close(ch)
		delete(csm.subscribers[sessionID], userKey)
	}

	// Record activity
	go csm.recordActivity(
		context.Background(),
		sessionID,
		session.ResourceID,
		session.ResourceType,
		userID,
		username,
		"leave",
		nil,
	)

	// Broadcast leave event
	csm.broadcastToSession(sessionID, map[string]interface{}{
		"type":      "user_left",
		"user_id":   userID,
		"username":  username,
		"timestamp": time.Now(),
	}, 0)

	return nil
}

// UpdateCursor updates user's cursor position
func (csm *CollaborationSessionManager) UpdateCursor(
	ctx context.Context,
	sessionID string,
	userID int64,
	position int32,
	line int32,
	column int32,
) error {
	csm.mu.RLock()
	session, ok := csm.sessions[sessionID]
	cursorChan, hasChannel := csm.cursorUpdates[sessionID]
	csm.mu.RUnlock()

	if !ok {
		return fmt.Errorf("session not found")
	}

	session.mu.Lock()
	participant, exists := session.Participants[userID]
	if !exists {
		session.mu.Unlock()
		return fmt.Errorf("participant not found")
	}

	participant.CursorPos = position
	participant.CursorLine = line
	session.mu.Unlock()

	// Send cursor update
	if hasChannel {
		update := &models.CursorPosition{
			UserID:    userID,
			Username:  participant.Username,
			Position:  position,
			Line:      line,
			Column:    column,
			Color:     participant.Color,
			Timestamp: time.Now(),
		}

		select {
		case cursorChan <- update:
		default:
			// Channel full, skip
		}
	}

	return nil
}

// broadcastToSession broadcasts a message to all participants except sender
func (csm *CollaborationSessionManager) broadcastToSession(
	sessionID string,
	message interface{},
	senderID int64,
) {
	csm.mu.RLock()
	subscribers, ok := csm.subscribers[sessionID]
	csm.mu.RUnlock()

	if !ok {
		return
	}

	for userKey, channel := range subscribers {
		// Skip sender
		if senderID > 0 && userKey == fmt.Sprintf("user_%d", senderID) {
			continue
		}

		select {
		case channel <- message:
		default:
			// Channel full, skip
		}
	}
}

// GetSession retrieves session information
func (csm *CollaborationSessionManager) GetSession(
	ctx context.Context,
	sessionID string,
) (*CollaborationSession, error) {
	csm.mu.RLock()
	session, ok := csm.sessions[sessionID]
	csm.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("session not found")
	}

	return session, nil
}

// GetParticipants returns list of active participants
func (csm *CollaborationSessionManager) GetParticipants(
	sessionID string,
) ([]*Participant, error) {
	csm.mu.RLock()
	session, ok := csm.sessions[sessionID]
	csm.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("session not found")
	}

	session.mu.RLock()
	defer session.mu.RUnlock()

	var participants []*Participant
	for _, p := range session.Participants {
		participants = append(participants, p)
	}

	return participants, nil
}

// recordActivity records an activity event
func (csm *CollaborationSessionManager) recordActivity(
	ctx context.Context,
	sessionID string,
	resourceID int64,
	resourceType string,
	userID int64,
	username string,
	eventType string,
	data map[string]interface{},
) {
	event := &models.ActivityEvent{
		SessionID:    sessionID,
		UserID:       userID,
		Username:     username,
		EventType:    eventType,
		ResourceID:   resourceID,
		ResourceType: resourceType,
		Data:         data,
		Timestamp:    time.Now(),
		CreatedAt:    time.Now(),
	}

	_ = csm.activityRepository.Create(ctx, event)
}

// cleanupInactiveSessions removes inactive sessions
func (csm *CollaborationSessionManager) cleanupInactiveSessions() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		csm.mu.Lock()

		now := time.Now()
		for sessionID, session := range csm.sessions {
			// Remove if no activity for 30 minutes
			if now.Sub(session.LastActivityAt) > 30*time.Minute {
				session.mu.Lock()
				participants := len(session.Participants)
				session.mu.Unlock()

				// Also close all subscriptions
				if subs, ok := csm.subscribers[sessionID]; ok {
					for _, ch := range subs {
						close(ch)
					}
					delete(csm.subscribers, sessionID)
				}

				if cursorChan, ok := csm.cursorUpdates[sessionID]; ok {
					close(cursorChan)
					delete(csm.cursorUpdates, sessionID)
				}

				delete(csm.sessions, sessionID)

				// Mark as inactive in DB
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				csm.sessionRepository.MarkInactive(ctx, sessionID)
				cancel()
			}
		}

		csm.mu.Unlock()
	}
}

// GetActivityStream returns activity events for a session
func (csm *CollaborationSessionManager) GetActivityStream(
	ctx context.Context,
	sessionID string,
	limit int,
) ([]*models.ActivityEvent, error) {
	return csm.activityRepository.GetBySessionID(ctx, sessionID, limit)
}

// GetParticipantCursors returns all participant cursors
func (csm *CollaborationSessionManager) GetParticipantCursors(
	sessionID string,
) ([]models.CursorPosition, error) {
	csm.mu.RLock()
	session, ok := csm.sessions[sessionID]
	csm.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("session not found")
	}

	session.mu.RLock()
	defer session.mu.RUnlock()

	var cursors []models.CursorPosition
	for _, p := range session.Participants {
		if p.IsActive {
			cursors = append(cursors, models.CursorPosition{
				UserID:    p.UserID,
				Username:  p.Username,
				Position:  p.CursorPos,
				Line:      p.CursorLine,
				Color:     p.Color,
				Timestamp: time.Now(),
			})
		}
	}

	return cursors, nil
}
