//go:build legacy
// +build legacy

package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"pos-saas/internal/repository"
)

// CollaborationCommentService handles comments in collaborative editing
type CollaborationCommentService struct {
	commentRepository repository.CommentRepository
	subscribers       map[string]chan *models.Comment // sessionID -> subscribers
	mu                sync.RWMutex
}

// NewCollaborationCommentService creates new comment service
func NewCollaborationCommentService(
	commentRepository repository.CommentRepository,
) *CollaborationCommentService {
	return &CollaborationCommentService{
		commentRepository: commentRepository,
		subscribers:       make(map[string]chan *models.Comment),
	}
}

// AddComment adds a comment to a document
func (ccs *CollaborationCommentService) AddComment(
	ctx context.Context,
	sessionID string,
	resourceID int64,
	resourceType string,
	userID int64,
	username string,
	position int32,
	lineNumber int32,
	text string,
) (*models.Comment, error) {
	comment := &models.Comment{
		SessionID:    sessionID,
		ResourceID:   resourceID,
		ResourceType: resourceType,
		UserID:       userID,
		Username:     username,
		Position:     position,
		LineNumber:   lineNumber,
		Text:         text,
		Resolved:     false,
		Reactions:    make(map[string][]int64),
		Replies:      []*models.Comment{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := ccs.commentRepository.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Notify subscribers
	ccs.notifySubscribers(sessionID, comment)

	return comment, nil
}

// ReplyToComment adds a reply to a comment
func (ccs *CollaborationCommentService) ReplyToComment(
	ctx context.Context,
	sessionID string,
	parentCommentID int64,
	userID int64,
	username string,
	text string,
) (*models.Comment, error) {
	// Get parent comment
	parentComments, err := ccs.commentRepository.GetByID(ctx, parentCommentID)
	if err != nil || len(parentComments) == 0 {
		return nil, fmt.Errorf("parent comment not found")
	}

	parentComment := parentComments[0]

	// Create reply
	reply := &models.Comment{
		SessionID:    parentComment.SessionID,
		ResourceID:   parentComment.ResourceID,
		ResourceType: parentComment.ResourceType,
		UserID:       userID,
		Username:     username,
		Position:     parentComment.Position,
		LineNumber:   parentComment.LineNumber,
		Text:         text,
		Resolved:     false,
		Reactions:    make(map[string][]int64),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := ccs.commentRepository.Create(ctx, reply); err != nil {
		return nil, fmt.Errorf("failed to create reply: %w", err)
	}

	// Add reply to parent
	parentComment.Replies = append(parentComment.Replies, reply)

	// Notify subscribers
	ccs.notifySubscribers(sessionID, reply)

	return reply, nil
}

// AddReaction adds an emoji reaction to a comment
func (ccs *CollaborationCommentService) AddReaction(
	ctx context.Context,
	commentID int64,
	userID int64,
	emoji string,
) error {
	// Get comment
	comments, err := ccs.commentRepository.GetByID(ctx, commentID)
	if err != nil || len(comments) == 0 {
		return fmt.Errorf("comment not found")
	}

	comment := comments[0]

	// Add reaction
	if comment.Reactions == nil {
		comment.Reactions = make(map[string][]int64)
	}

	// Check if user already reacted
	users := comment.Reactions[emoji]
	for _, uid := range users {
		if uid == userID {
			// Already reacted, remove it
			var newUsers []int64
			for _, u := range users {
				if u != userID {
					newUsers = append(newUsers, u)
				}
			}
			comment.Reactions[emoji] = newUsers
			if len(newUsers) == 0 {
				delete(comment.Reactions, emoji)
			}
			break
		}
	}

	// Add new reaction
	comment.Reactions[emoji] = append(comment.Reactions[emoji], userID)

	// Update in database
	if err := ccs.commentRepository.Update(ctx, comment); err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

// ResolveComment marks a comment as resolved
func (ccs *CollaborationCommentService) ResolveComment(
	ctx context.Context,
	commentID int64,
	resolvedBy int64,
) error {
	// Get comment
	comments, err := ccs.commentRepository.GetByID(ctx, commentID)
	if err != nil || len(comments) == 0 {
		return fmt.Errorf("comment not found")
	}

	comment := comments[0]
	now := time.Now()

	comment.Resolved = true
	comment.ResolvedBy = &resolvedBy
	comment.ResolvedAt = &now
	comment.UpdatedAt = now

	if err := ccs.commentRepository.Update(ctx, comment); err != nil {
		return fmt.Errorf("failed to resolve comment: %w", err)
	}

	return nil
}

// GetComments retrieves comments for a resource
func (ccs *CollaborationCommentService) GetComments(
	ctx context.Context,
	resourceID int64,
	resourceType string,
	unresolvedOnly bool,
) ([]*models.Comment, error) {
	return ccs.commentRepository.GetByResource(
		ctx,
		resourceID,
		resourceType,
		unresolvedOnly,
	)
}

// GetCommentsByPosition retrieves comments at a specific position
func (ccs *CollaborationCommentService) GetCommentsByPosition(
	ctx context.Context,
	sessionID string,
	position int32,
) ([]*models.Comment, error) {
	return ccs.commentRepository.GetByPosition(ctx, sessionID, position)
}

// SubscribeToComments subscribes to comment updates for a session
func (ccs *CollaborationCommentService) SubscribeToComments(
	sessionID string,
) <-chan *models.Comment {
	ccs.mu.Lock()
	defer ccs.mu.Unlock()

	channel := make(chan *models.Comment, 50)
	ccs.subscribers[sessionID] = channel

	return channel
}

// UnsubscribeFromComments unsubscribes from comment updates
func (ccs *CollaborationCommentService) UnsubscribeFromComments(
	sessionID string,
) {
	ccs.mu.Lock()
	if ch, ok := ccs.subscribers[sessionID]; ok {
		close(ch)
		delete(ccs.subscribers, sessionID)
	}
	ccs.mu.Unlock()
}

// notifySubscribers broadcasts comment to subscribers
func (ccs *CollaborationCommentService) notifySubscribers(
	sessionID string,
	comment *models.Comment,
) {
	ccs.mu.RLock()
	ch, ok := ccs.subscribers[sessionID]
	ccs.mu.RUnlock()

	if !ok {
		return
	}

	select {
	case ch <- comment:
	default:
		// Channel full, skip
	}
}

// EditComment edits comment text
func (ccs *CollaborationCommentService) EditComment(
	ctx context.Context,
	commentID int64,
	newText string,
) error {
	comments, err := ccs.commentRepository.GetByID(ctx, commentID)
	if err != nil || len(comments) == 0 {
		return fmt.Errorf("comment not found")
	}

	comment := comments[0]
	comment.Text = newText
	comment.UpdatedAt = time.Now()

	return ccs.commentRepository.Update(ctx, comment)
}

// DeleteComment deletes a comment
func (ccs *CollaborationCommentService) DeleteComment(
	ctx context.Context,
	commentID int64,
) error {
	return ccs.commentRepository.Delete(ctx, commentID)
}

// MentionedUsers extracts mentioned users from comment text
func (ccs *CollaborationCommentService) MentionedUsers(text string) []string {
	// Parse @mentions from comment
	// Returns slice of usernames
	var mentioned []string

	// Simple regex-based parsing (would be more robust in production)
	// Matches @username patterns
	// for _, match := range regexp.FindAllString(text, -1) {
	//     mentioned = append(mentioned, match[1:]) // Remove @
	// }

	return mentioned
}

// CommentStats returns statistics for a comment thread
type CommentStats struct {
	TotalComments    int
	ResolvedComments int
	OpenComments     int
	ParticipatingUsers []int64
	LastCommentAt    time.Time
}

// GetCommentStats returns statistics for a resource's comments
func (ccs *CollaborationCommentService) GetCommentStats(
	ctx context.Context,
	resourceID int64,
	resourceType string,
) (*CommentStats, error) {
	comments, err := ccs.GetComments(ctx, resourceID, resourceType, false)
	if err != nil {
		return nil, err
	}

	stats := &CommentStats{
		TotalComments:     len(comments),
		ParticipatingUsers: make([]int64, 0),
	}

	userSet := make(map[int64]bool)
	var lastTime time.Time

	for _, comment := range comments {
		if comment.Resolved {
			stats.ResolvedComments++
		} else {
			stats.OpenComments++
		}

		if !userSet[comment.UserID] {
			stats.ParticipatingUsers = append(stats.ParticipatingUsers, comment.UserID)
			userSet[comment.UserID] = true
		}

		if comment.UpdatedAt.After(lastTime) {
			lastTime = comment.UpdatedAt
		}
	}

	stats.LastCommentAt = lastTime

	return stats, nil
}
