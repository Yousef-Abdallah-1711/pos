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

// OperationalTransformationService handles conflict-free collaborative editing
type OperationalTransformationService struct {
	editRepository     repository.CollaborativeEditRepository
	otRepository       repository.OperationalTransformRepository
	conflictRepository repository.ConflictResolutionRepository
	versionRepository  repository.DocumentVersionRepository
	sessionManager     *CollaborationSessionManager
	mu                 sync.RWMutex
	operationQueues    map[string][]*models.OperationalTransform // sessionID -> operations
	versionVectors     map[string]map[int64]int64                 // sessionID -> userID -> version
}

// NewOperationalTransformationService creates new OT service
func NewOperationalTransformationService(
	editRepository repository.CollaborativeEditRepository,
	otRepository repository.OperationalTransformRepository,
	conflictRepository repository.ConflictResolutionRepository,
	versionRepository repository.DocumentVersionRepository,
	sessionManager *CollaborationSessionManager,
) *OperationalTransformationService {
	return &OperationalTransformationService{
		editRepository:     editRepository,
		otRepository:       otRepository,
		conflictRepository: conflictRepository,
		versionRepository:  versionRepository,
		sessionManager:     sessionManager,
		operationQueues:    make(map[string][]*models.OperationalTransform),
		versionVectors:     make(map[string]map[int64]int64),
	}
}

// Operation represents an edit operation
type Operation struct {
	Type     string // insert, delete, retain
	Position int32  // position in document
	Content  string // for insert
	Length   int32  // for delete/retain
}

// ApplyOperation applies an operation to a document
func (ots *OperationalTransformationService) ApplyOperation(
	ctx context.Context,
	sessionID string,
	userID int64,
	operation *Operation,
	currentDoc string,
) (string, error) {
	ots.mu.Lock()
	defer ots.mu.Unlock()

	// Apply operation to document
	result := ots.applyOpToString(operation, currentDoc)

	// Create OT record
	otOp := &models.OperationalTransform{
		SessionID: sessionID,
		UserID:    userID,
		Operation: operation.Type,
		Position:  operation.Position,
		Content:   operation.Content,
		Length:    operation.Length,
		Timestamp: time.Now(),
		VersionVector: ots.versionVectors[sessionID],
	}

	// Save operation
	if err := ots.otRepository.Create(ctx, otOp); err != nil {
		return "", fmt.Errorf("failed to save operation: %w", err)
	}

	// Queue operation for distribution
	ots.operationQueues[sessionID] = append(
		ots.operationQueues[sessionID],
		otOp,
	)

	// Update version vector
	if ots.versionVectors[sessionID] == nil {
		ots.versionVectors[sessionID] = make(map[int64]int64)
	}
	ots.versionVectors[sessionID][userID]++

	return result, nil
}

// Transform transforms an operation against another operation (OT algorithm)
// Implements Google Docs-style OT
func (ots *OperationalTransformationService) Transform(
	op1 *Operation,
	op2 *Operation,
	priority string, // "left" or "right" for conflict resolution
) (*Operation, *Operation) {
	op1Prime := *op1
	op2Prime := *op2

	if op1.Type == "insert" && op2.Type == "insert" {
		if op1.Position < op2.Position {
			op2Prime.Position += int32(len(op1.Content))
		} else if op1.Position > op2.Position {
			op1Prime.Position += int32(len(op2.Content))
		} else {
			// Same position - use priority
			if priority == "left" {
				op2Prime.Position += int32(len(op1.Content))
			} else {
				op1Prime.Position += int32(len(op2.Content))
			}
		}
	} else if op1.Type == "insert" && op2.Type == "delete" {
		if op1.Position <= op2.Position {
			op2Prime.Position += int32(len(op1.Content))
		} else if op1.Position >= op2.Position+op2.Length {
			op1Prime.Position -= op2.Length
		} else {
			// Overlapping - delete part of insert
			deleteCount := op2.Position + op2.Length - op1.Position
			op1Prime.Content = op1.Content[:len(op1.Content)-int(deleteCount)]
		}
	} else if op1.Type == "delete" && op2.Type == "insert" {
		if op2.Position <= op1.Position {
			op1Prime.Position += int32(len(op2.Content))
		} else if op2.Position >= op1.Position+op1.Length {
			op2Prime.Position -= op1.Length
		} else {
			// Overlapping
			deleteCount := op1.Position + op1.Length - op2.Position
			// Position adjustment
			op1Prime.Position += int32(len(op2.Content))
		}
	} else if op1.Type == "delete" && op2.Type == "delete" {
		if op1.Position < op2.Position {
			op2Prime.Position -= op1.Length
		} else if op1.Position > op2.Position {
			op1Prime.Position -= op2.Length
		} else {
			// Same position - delete minimum
			minLen := op1.Length
			if op2.Length < minLen {
				minLen = op2.Length
			}
			op1Prime.Length -= minLen
			op2Prime.Length -= minLen
		}
	}

	return &op1Prime, &op2Prime
}

// ResolveConflict resolves conflicting operations
func (ots *OperationalTransformationService) ResolveConflict(
	ctx context.Context,
	sessionID string,
	conflictingOps []*models.OperationalTransform,
	strategy string,
) (*models.OperationalTransform, error) {
	if len(conflictingOps) == 0 {
		return nil, fmt.Errorf("no conflicting operations")
	}

	var resolved *models.OperationalTransform

	switch strategy {
	case "ot":
		// Use operational transformation to resolve
		resolved = ots.resolveViaOT(conflictingOps)
	case "vote":
		// Use voting/consensus
		resolved = ots.resolveViaVoting(conflictingOps)
	default:
		// First-come-first-served
		resolved = conflictingOps[0]
	}

	// Record resolution
	conflictRes := &models.ConflictResolution{
		SessionID:        sessionID,
		ResolutionMethod: strategy,
		ResolvedAt:       time.Now(),
		ResolvedBy:       0, // automated
	}

	if err := ots.conflictRepository.Create(ctx, conflictRes); err != nil {
		return nil, fmt.Errorf("failed to record conflict resolution: %w", err)
	}

	return resolved, nil
}

// resolveViaOT resolves using operational transformation
func (ots *OperationalTransformationService) resolveViaOT(
	ops []*models.OperationalTransform,
) *models.OperationalTransform {
	if len(ops) == 0 {
		return nil
	}

	// Start with first operation
	result := ops[0]

	// Transform all subsequent operations against the first
	for i := 1; i < len(ops); i++ {
		otOp := &Operation{
			Type:     ops[i].Operation,
			Position: ops[i].Position,
			Content:  ops[i].Content,
			Length:   ops[i].Length,
		}

		resultOp := &Operation{
			Type:     result.Operation,
			Position: result.Position,
			Content:  result.Content,
			Length:   result.Length,
		}

		// Determine priority based on timestamp
		priority := "left"
		if ops[i].Timestamp.Before(result.Timestamp) {
			priority = "right"
		}

		transformed1, _ := ots.Transform(resultOp, otOp, priority)

		result.Operation = transformed1.Type
		result.Position = transformed1.Position
		result.Content = transformed1.Content
		result.Length = transformed1.Length
	}

	return result
}

// resolveViaVoting resolves using voting
func (ots *OperationalTransformationService) resolveViaVoting(
	ops []*models.OperationalTransform,
) *models.OperationalTransform {
	// Use the operation with most recent timestamp
	var latest *models.OperationalTransform
	for _, op := range ops {
		if latest == nil || op.Timestamp.After(latest.Timestamp) {
			latest = op
		}
	}
	return latest
}

// applyOpToString applies an operation to a string document
func (ots *OperationalTransformationService) applyOpToString(
	op *Operation,
	doc string,
) string {
	pos := op.Position
	switch op.Type {
	case "insert":
		if pos > int32(len(doc)) {
			pos = int32(len(doc))
		}
		return doc[:pos] + op.Content + doc[pos:]

	case "delete":
		endPos := pos + op.Length
		if endPos > int32(len(doc)) {
			endPos = int32(len(doc))
		}
		if pos > endPos {
			pos = endPos
		}
		return doc[:pos] + doc[endPos:]

	case "retain":
		// No change
		return doc

	default:
		return doc
	}
}

// GetOperationHistory retrieves operation history for a session
func (ots *OperationalTransformationService) GetOperationHistory(
	ctx context.Context,
	sessionID string,
	limit int,
) ([]*models.OperationalTransform, error) {
	return ots.otRepository.GetBySessionID(ctx, sessionID, limit)
}

// CreateSnapshot creates a document snapshot at current version
func (ots *OperationalTransformationService) CreateSnapshot(
	ctx context.Context,
	sessionID string,
	resourceID int64,
	content string,
	createdBy int64,
	opCount int64,
) (*models.DocumentVersion, error) {
	version := &models.DocumentVersion{
		SessionID:     sessionID,
		ResourceID:    resourceID,
		Content:       content,
		SnapshotAt:    opCount,
		CreatedBy:     createdBy,
		CreatedAt:     time.Now(),
	}

	// Calculate version number
	versions, err := ots.versionRepository.GetByResourceID(ctx, resourceID, 1)
	if err == nil && len(versions) > 0 {
		version.VersionNumber = versions[0].VersionNumber + 1
	} else {
		version.VersionNumber = 1
	}

	if err := ots.versionRepository.Create(ctx, version); err != nil {
		return nil, fmt.Errorf("failed to create snapshot: %w", err)
	}

	return version, nil
}

// GetLatestVersion retrieves the latest document version
func (ots *OperationalTransformationService) GetLatestVersion(
	ctx context.Context,
	resourceID int64,
) (*models.DocumentVersion, error) {
	versions, err := ots.versionRepository.GetByResourceID(ctx, resourceID, 1)
	if err != nil {
		return nil, err
	}

	if len(versions) == 0 {
		return nil, fmt.Errorf("no versions found")
	}

	return versions[0], nil
}

// ReconstructDocument reconstructs document from snapshot + operations
func (ots *OperationalTransformationService) ReconstructDocument(
	ctx context.Context,
	resourceID int64,
) (string, error) {
	// Get latest snapshot
	version, err := ots.GetLatestVersion(ctx, resourceID)
	if err != nil {
		return "", err
	}

	doc := version.Content

	// Apply operations since snapshot
	ops, err := ots.otRepository.GetSinceSnapshot(ctx, version.SessionID, version.SnapshotAt, 10000)
	if err != nil {
		return "", err
	}

	// Apply each operation
	for _, op := range ops {
		operation := &Operation{
			Type:     op.Operation,
			Position: op.Position,
			Content:  op.Content,
			Length:   op.Length,
		}
		doc = ots.applyOpToString(operation, doc)
	}

	return doc, nil
}

// ValidateOperation validates an operation
func (ots *OperationalTransformationService) ValidateOperation(
	op *Operation,
	docLength int32,
) error {
	switch op.Type {
	case "insert":
		if op.Position > docLength {
			return fmt.Errorf("position %d out of bounds (doc length %d)", op.Position, docLength)
		}
		if len(op.Content) == 0 {
			return fmt.Errorf("insert operation must have content")
		}

	case "delete":
		if op.Position > docLength {
			return fmt.Errorf("position %d out of bounds", op.Position)
		}
		if op.Position+op.Length > docLength {
			return fmt.Errorf("delete range exceeds document length")
		}
		if op.Length <= 0 {
			return fmt.Errorf("delete length must be positive")
		}

	case "retain":
		if op.Length <= 0 {
			return fmt.Errorf("retain length must be positive")
		}

	default:
		return fmt.Errorf("unknown operation type: %s", op.Type)
	}

	return nil
}
