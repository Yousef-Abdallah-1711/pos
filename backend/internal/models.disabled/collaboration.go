//go:build legacy
// +build legacy

package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// CollaborationSession represents an active collaboration session
type CollaborationSession struct {
	ID                 int64                  `json:"id" db:"id"`
	RestaurantID       int64                  `json:"restaurant_id" db:"restaurant_id"`
	ResourceID         int64                  `json:"resource_id" db:"resource_id"`                   // component/theme ID
	ResourceType       string                 `json:"resource_type" db:"resource_type"`               // component, theme
	SessionID          string                 `json:"session_id" db:"session_id"`                     // unique session key
	IsActive           bool                   `json:"is_active" db:"is_active"`
	ParticipantCount   int32                  `json:"participant_count" db:"participant_count"`
	Participants       map[string]*Participant `json:"participants" db:"participants"`
	CreatedAt          time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at" db:"updated_at"`
	ExpiresAt          time.Time              `json:"expires_at" db:"expires_at"`
}

// Scan implements sql.Scanner
func (cs *CollaborationSession) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), cs)
}

// Value implements driver.Valuer
func (cs CollaborationSession) Value() (driver.Value, error) {
	return json.Marshal(cs)
}

// Participant represents a user in a collaboration session
type Participant struct {
	UserID       int64     `json:"user_id"`
	Username     string    `json:"username"`
	CursorPos    int32     `json:"cursor_pos"`    // position in document
	CursorLine   int32     `json:"cursor_line"`   // line number
	Color        string    `json:"color"`         // for cursor visualization
	LastActivity time.Time `json:"last_activity"`
}

// OperationalTransform represents an edit operation in OT
type OperationalTransform struct {
	ID              int64                  `json:"id" db:"id"`
	SessionID       string                 `json:"session_id" db:"session_id"`
	RestaurantID    int64                  `json:"restaurant_id" db:"restaurant_id"`
	UserID          int64                  `json:"user_id" db:"user_id"`
	Operation       string                 `json:"operation" db:"operation"` // insert, delete, retain
	Position        int32                  `json:"position" db:"position"`   // position in document
	Content         string                 `json:"content" db:"content"`     // inserted/deleted text
	Length          int32                  `json:"length" db:"length"`       // for delete/retain
	Timestamp       time.Time              `json:"timestamp" db:"timestamp"`
	VersionVector   map[string]int64       `json:"version_vector" db:"version_vector"` // for causality tracking
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
}

// Scan implements sql.Scanner
func (ot *OperationalTransform) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), ot)
}

// Value implements driver.Valuer
func (ot OperationalTransform) Value() (driver.Value, error) {
	return json.Marshal(ot)
}

// CollaborativeEdit represents a merged edit
type CollaborativeEdit struct {
	ID            int64     `json:"id" db:"id"`
	SessionID     string    `json:"session_id" db:"session_id"`
	ResourceID    int64     `json:"resource_id" db:"resource_id"`
	ResourceType  string    `json:"resource_type" db:"resource_type"`
	UserID        int64     `json:"user_id" db:"user_id"`
	VersionNumber int64     `json:"version_number" db:"version_number"`
	Operation     string    `json:"operation" db:"operation"` // JSON-encoded OT
	Conflicts     int32     `json:"conflicts" db:"conflicts"` // number of conflicts resolved
	AppliedAt     time.Time `json:"applied_at" db:"applied_at"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// Comment represents a collaboration comment/annotation
type Comment struct {
	ID             int64                  `json:"id" db:"id"`
	SessionID      string                 `json:"session_id" db:"session_id"`
	ResourceID     int64                  `json:"resource_id" db:"resource_id"`
	ResourceType   string                 `json:"resource_type" db:"resource_type"`
	UserID         int64                  `json:"user_id" db:"user_id"`
	Username       string                 `json:"username" db:"username"`
	Position       int32                  `json:"position" db:"position"`       // position in document
	LineNumber     int32                  `json:"line_number" db:"line_number"` // for code editors
	Text           string                 `json:"text" db:"text"`               // comment content
	Resolved       bool                   `json:"resolved" db:"resolved"`
	ResolvedBy     *int64                 `json:"resolved_by" db:"resolved_by"`
	ResolvedAt     *time.Time             `json:"resolved_at" db:"resolved_at"`
	Reactions      map[string][]int64     `json:"reactions" db:"reactions"` // emoji -> userIDs
	Replies        []*Comment             `json:"replies" db:"replies"`
	CreatedAt      time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at" db:"updated_at"`
}

// Scan implements sql.Scanner
func (c *Comment) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), c)
}

// Value implements driver.Valuer
func (c Comment) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// ActivityEvent represents an activity in the collaboration stream
type ActivityEvent struct {
	ID            int64                  `json:"id" db:"id"`
	SessionID     string                 `json:"session_id" db:"session_id"`
	RestaurantID  int64                  `json:"restaurant_id" db:"restaurant_id"`
	UserID        int64                  `json:"user_id" db:"user_id"`
	Username      string                 `json:"username" db:"username"`
	EventType     string                 `json:"event_type" db:"event_type"` // edit, comment, join, leave, resolve
	ResourceID    int64                  `json:"resource_id" db:"resource_id"`
	ResourceType  string                 `json:"resource_type" db:"resource_type"`
	Data          map[string]interface{} `json:"data" db:"data"` // event-specific data
	Timestamp     time.Time              `json:"timestamp" db:"timestamp"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
}

// Scan implements sql.Scanner
func (ae *ActivityEvent) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), ae)
}

// Value implements driver.Valuer
func (ae ActivityEvent) Value() (driver.Value, error) {
	return json.Marshal(ae)
}

// ConflictResolution represents a resolved conflict
type ConflictResolution struct {
	ID                int64     `json:"id" db:"id"`
	SessionID         string    `json:"session_id" db:"session_id"`
	ResourceID        int64     `json:"resource_id" db:"resource_id"`
	ConflictingOps    []string  `json:"conflicting_ops" db:"conflicting_ops"` // OT operation IDs
	ResolutionMethod  string    `json:"resolution_method" db:"resolution_method"` // ot, vote, manual
	ResolvedAt        time.Time `json:"resolved_at" db:"resolved_at"`
	ResolvedBy        int64     `json:"resolved_by" db:"resolved_by"` // user who resolved (0 = automated)
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

// DocumentVersion represents a version of a collaborative document
type DocumentVersion struct {
	ID             int64     `json:"id" db:"id"`
	SessionID      string    `json:"session_id" db:"session_id"`
	ResourceID     int64     `json:"resource_id" db:"resource_id"`
	VersionNumber  int64     `json:"version_number" db:"version_number"`
	Content        string    `json:"content" db:"content"`
	SnapshotAt     int64     `json:"snapshot_at" db:"snapshot_at"` // operation count at snapshot
	CreatedBy      int64     `json:"created_by" db:"created_by"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// CollaborationConfig stores collaboration settings
type CollaborationConfig struct {
	ID                          int64                  `json:"id" db:"id"`
	RestaurantID                int64                  `json:"restaurant_id" db:"restaurant_id"`
	IsEnabled                   bool                   `json:"is_enabled" db:"is_enabled"`
	MaxParticipants             int32                  `json:"max_participants" db:"max_participants"`
	SessionTimeout              int32                  `json:"session_timeout" db:"session_timeout"`           // seconds
	ConflictResolutionStrategy  string                 `json:"conflict_resolution_strategy" db:"conflict_resolution_strategy"` // ot, vote
	EnableComments              bool                   `json:"enable_comments" db:"enable_comments"`
	EnableVersionHistory        bool                   `json:"enable_version_history" db:"enable_version_history"`
	SnapshotInterval            int32                  `json:"snapshot_interval" db:"snapshot_interval"`     // operations between snapshots
	MaxConcurrentEdits          int32                  `json:"max_concurrent_edits" db:"max_concurrent_edits"`
	CreatedAt                   time.Time              `json:"created_at" db:"created_at"`
}

// CursorPosition represents real-time cursor/selection position
type CursorPosition struct {
	UserID     int64 `json:"user_id"`
	Username   string `json:"username"`
	Position   int32 `json:"position"`   // character position
	Line       int32 `json:"line"`       // line number for multiline
	Column     int32 `json:"column"`     // column number
	Color      string `json:"color"`     // for visualization
	Timestamp  time.Time `json:"timestamp"`
}

// SelectionRange represents a selected range of text
type SelectionRange struct {
	UserID    int64 `json:"user_id"`
	StartPos  int32 `json:"start_pos"`
	EndPos    int32 `json:"end_pos"`
	Color     string `json:"color"`
	Timestamp time.Time `json:"timestamp"`
}
