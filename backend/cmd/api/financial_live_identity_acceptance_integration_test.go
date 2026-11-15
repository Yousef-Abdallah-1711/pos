//go:build integration

package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type liveFinancialIdentity struct {
	userID   int64
	email    string
	password string
	role     string
}

type liveFinancialPermission struct {
	moduleID int
	level    string
}

func seedLiveFinancialIdentity(t *testing.T, db *sql.DB, fixture liveFinancialFixture, role, roleCode string, permissions []liveFinancialPermission) liveFinancialIdentity {
	t.Helper()
	runID := strconv.FormatInt(time.Now().UnixNano(), 10)
	identity := liveFinancialIdentity{
		email:    "financial-" + role + "-" + runID + "@example.test",
		password: "synthetic-financial-" + runID,
		role:     role,
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(identity.password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash synthetic %s password: %v", role, err)
	}
	if err := db.QueryRow(`
		INSERT INTO users (tenant_id, restaurant_id, email, password_hash, name, role, status)
		VALUES ($1, $2, $3, $4, $5, $6, 'active') RETURNING id`,
		fixture.tenantID, fixture.restaurantID, identity.email, string(hash), "Synthetic "+role, role).Scan(&identity.userID); err != nil {
		t.Fatalf("seed synthetic %s identity: %v", role, err)
	}
	var roleID int64
	if err := db.QueryRow(`
		INSERT INTO roles (tenant_id, restaurant_id, role_name, role_code, access_level, is_active)
		VALUES ($1, $2, $3, $4, 'basic', TRUE) RETURNING id`,
		fixture.tenantID, fixture.restaurantID, "Acceptance "+role+" "+runID, roleCode+"-"+runID).Scan(&roleID); err != nil {
		t.Fatalf("seed synthetic %s RBAC role: %v", role, err)
	}
	for _, permission := range permissions {
		if _, err := db.Exec(`
			INSERT INTO role_permissions (tenant_id, role_id, module_id, permission_level)
			VALUES ($1, $2, $3, $4)`, fixture.tenantID, roleID, permission.moduleID, permission.level); err != nil {
			t.Fatalf("grant synthetic %s module %d permission: %v", role, permission.moduleID, err)
		}
	}
	ownerID := seedLiveFinancialOwner(t, db, fixture)
	if _, err := db.Exec(`
		INSERT INTO user_roles (tenant_id, user_id, role_id, assigned_by)
		VALUES ($1, $2, $3, $4)`, fixture.tenantID, identity.userID, roleID, ownerID); err != nil {
		t.Fatalf("assign synthetic %s RBAC role: %v", role, err)
	}
	return identity
}

func seedLiveFinancialOwner(t *testing.T, db *sql.DB, fixture liveFinancialFixture) int64 {
	t.Helper()
	if fixture.email != "" {
		var ownerID int64
		if err := db.QueryRow(`SELECT id FROM users WHERE tenant_id=$1 AND restaurant_id=$2 AND email=$3`,
			fixture.tenantID, fixture.restaurantID, fixture.email).Scan(&ownerID); err != nil {
			t.Fatalf("read synthetic owner for role assignment: %v", err)
		}
		return ownerID
	}
	runID := strconv.FormatInt(time.Now().UnixNano(), 10)
	email := "financial-owner-" + runID + "@example.test"
	hash, err := bcrypt.GenerateFromPassword([]byte("synthetic-owner-"+runID), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash synthetic owner password: %v", err)
	}
	var ownerID int64
	if err := db.QueryRow(`INSERT INTO users (tenant_id, restaurant_id, email, password_hash, name, role, status) VALUES ($1,$2,$3,$4,'Synthetic Owner','owner','active') RETURNING id`, fixture.tenantID, fixture.restaurantID, email, string(hash)).Scan(&ownerID); err != nil {
		t.Fatalf("seed synthetic owner for RBAC assignment: %v", err)
	}
	return ownerID
}

func loginLiveFinancialIdentity(t *testing.T, baseURL string, identity liveFinancialIdentity) *http.Client {
	t.Helper()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("create isolated %s cookie jar: %v", identity.role, err)
	}
	client := &http.Client{Jar: jar, Timeout: 10 * time.Second}
	status, body := acceptanceJSONRequest(t, client, http.MethodPost, baseURL+"/api/v1/auth/login", map[string]string{
		"email": identity.email, "password": identity.password,
	})
	if status != http.StatusOK {
		t.Fatalf("synthetic %s login returned %d: %s", identity.role, status, body)
	}
	return client
}

func financialAcceptancePath(baseURL, pattern string, id int64) string {
	return fmt.Sprintf(baseURL+pattern, id)
}
