package workspaceappearance

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"pos-saas/internal/pkg/database"
)

var ErrConflict = errors.New("workspace appearance revision conflict")

type Published struct {
	SchemaVersion int             `json:"schemaVersion"`
	Revision      int64           `json:"revision"`
	Overrides     json.RawMessage `json:"overrides"`
}

type Repository struct{ DB *sql.DB }

// Get returns an empty override document for a tenant without a published row.
func (r Repository) Get(ctx context.Context, tenantID, actorID int64) (Published, error) {
	result := Published{SchemaVersion: 1, Revision: 0, Overrides: json.RawMessage(`{}`)}
	err := database.WithTenantUserTx(ctx, r.DB, tenantID, 0, actorID, func(tx *sql.Tx) error {
		err := tx.QueryRowContext(ctx, `SELECT schema_version, revision, published_config
			FROM tenant_workspace_appearance WHERE tenant_id = $1`, tenantID).
			Scan(&result.SchemaVersion, &result.Revision, &result.Overrides)
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	})
	if err != nil {
		return Published{}, fmt.Errorf("read workspace appearance: %w", err)
	}
	return result, nil
}

// Save validates and records the publication and audit record atomically.
func (r Repository) Save(ctx context.Context, tenantID, actorID, expected int64, overrides json.RawMessage) (Published, error) {
	if expected < 0 {
		return Published{}, errors.New("negative revision")
	}
	if err := ValidateOverrides(overrides); err != nil {
		return Published{}, fmt.Errorf("validate workspace appearance: %w", err)
	}
	result := Published{SchemaVersion: 1, Revision: expected + 1, Overrides: append(json.RawMessage(nil), overrides...)}
	err := database.WithTenantUserTx(ctx, r.DB, tenantID, 0, actorID, func(tx *sql.Tx) error {
		previous, err := updateAppearance(ctx, tx, tenantID, actorID, expected, overrides)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `INSERT INTO tenant_workspace_appearance_audit
			(tenant_id, actor_id, revision, previous_config, published_config) VALUES ($1, $2, $3, $4, $5)`,
			tenantID, actorID, result.Revision, previous, string(overrides))
		return err
	})
	if err != nil {
		return Published{}, fmt.Errorf("save workspace appearance: %w", err)
	}
	return result, nil
}

func updateAppearance(ctx context.Context, tx *sql.Tx, tenantID, actorID, expected int64, config []byte) (any, error) {
	var previous []byte
	err := tx.QueryRowContext(ctx, `SELECT published_config FROM tenant_workspace_appearance
		WHERE tenant_id = $1 AND revision = $2 FOR UPDATE`, tenantID, expected).Scan(&previous)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) && expected != 0 {
		return nil, ErrConflict
	}
	var result sql.Result
	if expected == 0 {
		result, err = tx.ExecContext(ctx, `INSERT INTO tenant_workspace_appearance
			(tenant_id, revision, published_config, updated_by) VALUES ($1, 1, $2, $3)
			ON CONFLICT (tenant_id) DO NOTHING`, tenantID, string(config), actorID)
	} else {
		result, err = tx.ExecContext(ctx, `UPDATE tenant_workspace_appearance SET revision = revision + 1,
			published_config = $1, updated_by = $2, updated_at = now() WHERE tenant_id = $3 AND revision = $4`,
			string(config), actorID, tenantID, expected)
	}
	if err != nil {
		return nil, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if count != 1 {
		return nil, ErrConflict
	}
	if expected == 0 {
		return nil, nil
	}
	return string(previous), nil
}
