package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

// AuthLookupIdentity contains the credential material required to verify a
// password. Callers must never serialize or log PasswordHash.
type AuthLookupIdentity struct {
	ID              string
	NormalizedEmail string
	EmailVerifiedAt time.Time
	PasswordHash    string
	SecurityVersion int64
}

// AuthLookupMembership is a server-resolved tenant and restaurant choice. It
// contains no client-supplied authorization state.
type AuthLookupMembership struct {
	MembershipID   string
	TenantID       int
	LegacyUserID   sql.NullInt64
	TenantName     string
	RestaurantID   int
	RestaurantName string
	RoleID         int
	RoleCode       string
}

// AuthLookupRepository uses only the isolated lookup database pool. It must
// never be constructed with the ordinary pos_app database connection.
type AuthLookupRepository struct {
	db *sql.DB
}

func NewAuthLookupRepository(db *sql.DB) *AuthLookupRepository {
	return &AuthLookupRepository{db: db}
}

func (r *AuthLookupRepository) FindVerifiedIdentityByEmail(ctx context.Context, email string) (*AuthLookupIdentity, error) {
	if r == nil || r.db == nil {
		return nil, sql.ErrConnDone
	}
	var identity AuthLookupIdentity
	err := r.db.QueryRowContext(ctx, `
		SELECT identity_id, normalized_email, email_verified_at, password_hash, security_version
		FROM auth_private.lookup_verified_identity_by_email($1)
	`, strings.ToLower(strings.TrimSpace(email))).Scan(
		&identity.ID, &identity.NormalizedEmail, &identity.EmailVerifiedAt, &identity.PasswordHash, &identity.SecurityVersion,
	)
	if err != nil {
		return nil, err
	}
	return &identity, nil
}

func (r *AuthLookupRepository) FindEligibleMemberships(ctx context.Context, identityID string) ([]AuthLookupMembership, error) {
	if r == nil || r.db == nil {
		return nil, sql.ErrConnDone
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT membership_id, tenant_id, legacy_user_id, tenant_name,
		       restaurant_id, restaurant_name, role_id, role_code
		FROM auth_private.lookup_eligible_memberships($1::uuid)
		ORDER BY tenant_id, restaurant_id, role_id
	`, identityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	memberships := make([]AuthLookupMembership, 0)
	for rows.Next() {
		var membership AuthLookupMembership
		if err := rows.Scan(
			&membership.MembershipID, &membership.TenantID, &membership.LegacyUserID, &membership.TenantName,
			&membership.RestaurantID, &membership.RestaurantName, &membership.RoleID, &membership.RoleCode,
		); err != nil {
			return nil, err
		}
		memberships = append(memberships, membership)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return memberships, nil
}
