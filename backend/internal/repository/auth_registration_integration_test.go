package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"pos-saas/internal/domain"
	"pos-saas/tests/testdb"
)

func TestCreateTenantWithUserRegistrationContract(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	t.Setenv("ENV", "development")
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}

	stamp := time.Now().UnixNano()
	name := fmt.Sprintf("Cafe \"Nile\" \\ O'Brien 🍜\nKitchen %d", stamp)
	slug := fmt.Sprintf("regtest-%d", stamp)
	request := &domain.RegisterRequest{
		Name: "Test Owner", Email: fmt.Sprintf("registration-%d@example.test", stamp), Password: "synthetic-password-123",
		RestaurantName: name, RestaurantSlug: slug, Phone: "+16505550123",
	}
	user, err := NewAuthRepository(db).CreateTenantWithUserContext(ctx, request)
	if err != nil {
		t.Fatalf("create synthetic registration: %v", err)
	}
	if user.Phone != "+16505550123" || user.RestaurantID == nil {
		t.Fatalf("unexpected created owner: %#v", user)
	}
	allowed, reason, err := NewSubscriptionRepository(db).CheckTenantAccess(ctx, int64(user.TenantID))
	if err != nil || !allowed {
		t.Fatalf("new registration has no valid entitlement: reason=%s err=%v", reason, err)
	}
	var trialStart, trialEnd time.Time
	var trialDays int
	if err := db.QueryRow(`SELECT trial_start_date, trial_end_date, trial_days FROM tenants WHERE id = $1`, user.TenantID).Scan(&trialStart, &trialEnd, &trialDays); err != nil {
		t.Fatal(err)
	}
	if trialDays <= 0 || !trialEnd.Equal(trialStart.AddDate(0, 0, trialDays)) {
		t.Fatal("registration trial dates do not match configured duration")
	}

	var storedTenantPhone, storedUserPhone sql.NullString
	if err := db.QueryRow(`SELECT phone FROM tenants WHERE id = $1`, user.TenantID).Scan(&storedTenantPhone); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(`SELECT phone FROM users WHERE id = $1`, user.ID).Scan(&storedUserPhone); err != nil {
		t.Fatal(err)
	}
	if !storedTenantPhone.Valid || storedTenantPhone.String != "+16505550123" || !storedUserPhone.Valid || storedUserPhone.String != "+16505550123" {
		t.Fatalf("canonical phone not persisted: tenant=%#v user=%#v", storedTenantPhone, storedUserPhone)
	}

	var themeConfig []byte
	if err := db.QueryRow(`SELECT config FROM themes_v2 WHERE tenant_id = $1 AND restaurant_id = $2`, user.TenantID, *user.RestaurantID).Scan(&themeConfig); err != nil {
		t.Fatal(err)
	}
	var theme struct {
		Identity struct {
			SiteTitle string `json:"siteTitle"`
		} `json:"identity"`
		Footer struct {
			Text string `json:"text"`
		} `json:"footer"`
	}
	if err := json.Unmarshal(themeConfig, &theme); err != nil {
		t.Fatalf("default theme is not valid JSON: %v", err)
	}
	if theme.Identity.SiteTitle != name || theme.Footer.Text == "" {
		t.Fatalf("user name was not safely retained in theme: %#v", theme)
	}

	var roleCount int
	if err := db.QueryRow(`SELECT count(*) FROM user_roles WHERE tenant_id = $1 AND user_id = $2`, user.TenantID, user.ID).Scan(&roleCount); err != nil {
		t.Fatal(err)
	}
	if roleCount == 0 {
		t.Fatal("owner registration did not assign a role")
	}

	optional := &domain.RegisterRequest{
		Name: "Optional Phone Owner", Email: fmt.Sprintf("registration-optional-%d@example.test", stamp), Password: "synthetic-password-123",
		RestaurantName: fmt.Sprintf("Optional Phone Cafe %d", stamp), RestaurantSlug: fmt.Sprintf("regtest-optional-%d", stamp),
	}
	optionalUser, err := NewAuthRepository(db).CreateTenantWithUserContext(ctx, optional)
	if err != nil {
		t.Fatalf("create registration without optional phone: %v", err)
	}
	var optionalPhone sql.NullString
	if err := db.QueryRow(`SELECT phone FROM tenants WHERE id = $1`, optionalUser.TenantID).Scan(&optionalPhone); err != nil {
		t.Fatal(err)
	}
	if optionalPhone.Valid {
		t.Fatalf("blank optional phone should be stored as NULL, got %q", optionalPhone.String)
	}

	// Concurrent duplicate-email attempts must serialize and allow exactly one
	// account, without leaving tenant rows behind from the rejected request.
	raceEmail := fmt.Sprintf("registration-race-%d@example.test", stamp)
	requests := []*domain.RegisterRequest{
		{Name: "Race Owner A", Email: raceEmail, Password: "synthetic-password-123", RestaurantName: fmt.Sprintf("Race Cafe A %d", stamp), RestaurantSlug: fmt.Sprintf("race-a-%d", stamp)},
		{Name: "Race Owner B", Email: raceEmail, Password: "synthetic-password-123", RestaurantName: fmt.Sprintf("Race Cafe B %d", stamp), RestaurantSlug: fmt.Sprintf("race-b-%d", stamp)},
	}
	start := make(chan struct{})
	results := make(chan error, len(requests))
	var wait sync.WaitGroup
	for _, registration := range requests {
		wait.Add(1)
		go func(req *domain.RegisterRequest) {
			defer wait.Done()
			<-start
			_, err := NewAuthRepository(db).CreateTenantWithUserContext(ctx, req)
			results <- err
		}(registration)
	}
	close(start)
	wait.Wait()
	close(results)
	created, conflicts := 0, 0
	for err := range results {
		switch {
		case err == nil:
			created++
		case errors.Is(err, domain.ErrEmailAlreadyExists):
			conflicts++
		default:
			t.Fatalf("unexpected concurrent registration error: %v", err)
		}
	}
	if created != 1 || conflicts != 1 {
		t.Fatalf("concurrent duplicate email results: created=%d conflicts=%d", created, conflicts)
	}
	var matchingTenants int
	if err := db.QueryRow(`SELECT count(*) FROM tenants WHERE email = $1`, raceEmail).Scan(&matchingTenants); err != nil {
		t.Fatal(err)
	}
	if matchingTenants != 1 {
		t.Fatalf("concurrent duplicate email left %d tenants, want exactly one", matchingTenants)
	}

	// Concurrent registrations with different emails must also serialize on
	// the unique restaurant slug and roll back the losing tenant transaction.
	raceSlug := fmt.Sprintf("registration-slug-race-%d", stamp)
	requests = []*domain.RegisterRequest{
		{Name: "Slug Race Owner A", Email: fmt.Sprintf("registration-slug-a-%d@example.test", stamp), Password: "synthetic-password-123", RestaurantName: fmt.Sprintf("Slug Race Cafe A %d", stamp), RestaurantSlug: raceSlug},
		{Name: "Slug Race Owner B", Email: fmt.Sprintf("registration-slug-b-%d@example.test", stamp), Password: "synthetic-password-123", RestaurantName: fmt.Sprintf("Slug Race Cafe B %d", stamp), RestaurantSlug: raceSlug},
	}
	start = make(chan struct{})
	results = make(chan error, len(requests))
	wait = sync.WaitGroup{}
	for _, registration := range requests {
		wait.Add(1)
		go func(req *domain.RegisterRequest) {
			defer wait.Done()
			<-start
			_, err := NewAuthRepository(db).CreateTenantWithUserContext(ctx, req)
			results <- err
		}(registration)
	}
	close(start)
	wait.Wait()
	close(results)
	created, slugConflicts := 0, 0
	for err := range results {
		switch {
		case err == nil:
			created++
		case errors.Is(err, ErrTenantSlugExists):
			slugConflicts++
		default:
			t.Fatalf("unexpected concurrent slug registration error: %v", err)
		}
	}
	if created != 1 || slugConflicts != 1 {
		t.Fatalf("concurrent duplicate slug results: created=%d conflicts=%d", created, slugConflicts)
	}
	if err := db.QueryRow(`SELECT count(*) FROM tenants WHERE slug = $1`, raceSlug).Scan(&matchingTenants); err != nil {
		t.Fatal(err)
	}
	if matchingTenants != 1 {
		t.Fatalf("concurrent duplicate slug left %d tenants, want exactly one", matchingTenants)
	}
}
