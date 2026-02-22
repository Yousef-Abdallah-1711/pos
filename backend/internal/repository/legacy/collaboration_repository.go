//go:build legacy
// +build legacy

package repository

import (
	"context"
)

// CollaborationRepository handles collaboration session persistence
type CollaborationRepository interface {
	// Session operations
	CreateSession(ctx context.Context, session *models.CollaborationSession) error
	GetSession(ctx context.Context, sessionID string) (*models.CollaborationSession, error)
	UpdateSession(ctx context.Context, session *models.CollaborationSession) error
	DeleteSession(ctx context.Context, sessionID string) error
	GetActiveSessions(ctx context.Context, restaurantID int64) ([]*models.CollaborationSession, error)

	// Operation management
	AddOperation(ctx context.Context, sessionID string, operation *models.OperationalTransform) error
	GetOperations(ctx context.Context, sessionID string, limit int) ([]*models.OperationalTransform, error)
	GetOperationsSince(ctx context.Context, resourceID int64, version int64) ([]*models.OperationalTransform, error)
	ClearOperations(ctx context.Context, sessionID string) error

	// Conflict management
	LogConflict(ctx context.Context, conflict *models.ConflictResolution) error
	GetConflicts(ctx context.Context, sessionID string) ([]*models.ConflictResolution, error)
}

// CommentRepository handles comment persistence
type CommentRepository interface {
	// Comment CRUD
	Create(ctx context.Context, comment *models.Comment) error
	GetByID(ctx context.Context, commentID int64) (*models.Comment, error)
	Update(ctx context.Context, comment *models.Comment) error
	Delete(ctx context.Context, commentID int64) error

	// Query operations
	GetByResource(ctx context.Context, resourceID int64, resourceType string) ([]*models.Comment, error)
	GetByPosition(ctx context.Context, sessionID string, position int64) ([]*models.Comment, error)
	GetUnresolved(ctx context.Context, resourceID int64, resourceType string) ([]*models.Comment, error)
	GetByAuthor(ctx context.Context, userID int64) ([]*models.Comment, error)

	// Thread management
	GetReplies(ctx context.Context, parentCommentID int64) ([]*models.Comment, error)
	AddReply(ctx context.Context, parentCommentID int64, reply *models.Comment) error

	// Reactions
	UpdateReactions(ctx context.Context, commentID int64, reactions map[string][]int64) error
	GetReactions(ctx context.Context, commentID int64) (map[string][]int64, error)
}

// VersionRepository handles document version snapshots
type VersionRepository interface {
	// Version management
	CreateVersion(ctx context.Context, version *models.DocumentVersion) error
	GetVersion(ctx context.Context, versionID int64) (*models.DocumentVersion, error)
	GetLatestVersion(ctx context.Context, resourceID int64) (*models.DocumentVersion, error)
	GetVersionsForResource(ctx context.Context, resourceID int64, limit int) ([]*models.DocumentVersion, error)

	// Reconstruction
	GetVersionsSince(ctx context.Context, resourceID int64, opCount int64) ([]*models.DocumentVersion, error)
}

// ActivityRepository handles activity stream events
type ActivityRepository interface {
	// Event management
	LogEvent(ctx context.Context, event *models.ActivityEvent) error
	GetEvents(ctx context.Context, sessionID string, limit int) ([]*models.ActivityEvent, error)
	GetEventsByType(ctx context.Context, sessionID string, eventType string) ([]*models.ActivityEvent, error)
	GetEventsByUser(ctx context.Context, sessionID string, userID int64) ([]*models.ActivityEvent, error)
	ClearEvents(ctx context.Context, sessionID string) error
}

// CursorRepository handles cursor position tracking
type CursorRepository interface {
	// Cursor tracking
	UpdateCursor(ctx context.Context, sessionID string, userID int64, position *models.CursorPosition) error
	GetCursor(ctx context.Context, sessionID string, userID int64) (*models.CursorPosition, error)
	GetAllCursors(ctx context.Context, sessionID string) ([]*models.CursorPosition, error)
	RemoveCursor(ctx context.Context, sessionID string, userID int64) error
}
