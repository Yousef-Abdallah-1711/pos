package repository

import (
	"database/sql"
	"testing"
	"time"

	"pos-saas/internal/domain"
)

func TestRecurringTemplateRepositoryRejectsForeignScopeWithoutSideEffects(t *testing.T) {
	db, fixture := openDriverScopeFixture(t)
	repo := NewAccountingAdvancedRepository(db, NewJournalEntryRepository(db))
	targets := []struct {
		name         string
		fixtureCode  string
		tenantID     int
		restaurantID int
	}{
		{name: "sibling restaurant", fixtureCode: "sibling", tenantID: fixture.tenantA, restaurantID: fixture.restaurantA2},
		{name: "foreign tenant", fixtureCode: "foreign", tenantID: fixture.tenantB, restaurantID: fixture.restaurantB1},
	}
	templateIDs := make(map[string]int)
	for _, target := range targets {
		templateIDs[target.name] = insertRecurringTemplateFixture(t, db, target.tenantID, target.restaurantID, target.fixtureCode)
	}

	for _, target := range targets {
		t.Run(target.name, func(t *testing.T) {
			templateID := templateIDs[target.name]
			if template, err := repo.GetRecurringJournalTemplateByID(fixture.tenantA, fixture.restaurantA1, templateID); err != nil || template != nil {
				t.Fatalf("foreign template read returned template=%#v error=%v", template, err)
			}

			mutation := recurringTemplateUpdate(fixture.tenantA, fixture.restaurantA1, templateID)
			if err := repo.UpdateRecurringJournalTemplate(mutation); err == nil {
				t.Fatal("foreign template update unexpectedly succeeded")
			}
			if err := repo.DeleteRecurringJournalTemplate(fixture.tenantA, fixture.restaurantA1, templateID); err == nil {
				t.Fatal("foreign template delete unexpectedly succeeded")
			}
			if _, err := repo.RunRecurringJournalTemplate(fixture.tenantA, fixture.restaurantA1, templateID, 1, time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)); err == nil {
				t.Fatal("foreign template run unexpectedly succeeded")
			}

			var name, code string
			var active bool
			var lastRun sql.NullTime
			if err := db.QueryRow(`
				SELECT template_name, template_code, is_active, last_run_date
				FROM recurring_journal_templates
				WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3
			`, target.tenantID, target.restaurantID, templateID).Scan(&name, &code, &active, &lastRun); err != nil {
				t.Fatalf("read foreign template after rejected operations: %v", err)
			}
			if name != "Original "+target.fixtureCode || code != "ORIG-"+target.fixtureCode || !active || lastRun.Valid {
				t.Fatalf("rejected operations changed template: name=%q code=%q active=%t last_run=%v", name, code, active, lastRun)
			}
			var entries int
			if err := db.QueryRow(`SELECT COUNT(*) FROM journal_entries WHERE reference_type = 'recurring_template' AND reference_id = $1`, templateID).Scan(&entries); err != nil {
				t.Fatalf("count recurring entries after denied run: %v", err)
			}
			if entries != 0 {
				t.Fatalf("denied template run created %d journal entries", entries)
			}
		})
	}
}

func TestRecurringTemplateRepositoryRejectsMismatchedTenantRestaurantCreate(t *testing.T) {
	db, fixture := openDriverScopeFixture(t)
	repo := NewAccountingAdvancedRepository(db, NewJournalEntryRepository(db))
	input := &domain.RecurringJournalTemplate{
		TenantID: fixture.tenantA, RestaurantID: fixture.restaurantB1,
		TemplateCode: "WS02-MISMATCH", TemplateName: "Mismatched scope",
		Frequency: "monthly", StartDate: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		NextRunDate: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC), IsActive: true,
		LinesJSON: []byte("[]"),
	}
	if _, err := repo.CreateRecurringJournalTemplate(input); err == nil {
		t.Fatal("recurring template create accepted a restaurant owned by another tenant")
	}
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM recurring_journal_templates WHERE tenant_id = $1 AND restaurant_id = $2 AND template_code = $3`, fixture.tenantA, fixture.restaurantB1, input.TemplateCode).Scan(&count); err != nil {
		t.Fatalf("count mismatched recurring template after rejected create: %v", err)
	}
	if count != 0 {
		t.Fatalf("mismatched create left %d recurring template rows", count)
	}
}

func recurringTemplateUpdate(tenantID, restaurantID, id int) *domain.RecurringJournalTemplate {
	return &domain.RecurringJournalTemplate{
		ID: id, TenantID: tenantID, RestaurantID: restaurantID,
		TemplateCode: "FORGED", TemplateName: "Tampered",
		Frequency: "monthly", StartDate: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
		NextRunDate: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), IsActive: false,
		LinesJSON: []byte("[]"),
	}
}

func insertRecurringTemplateFixture(t *testing.T, db *sql.DB, tenantID, restaurantID int, suffix string) int {
	t.Helper()
	var id int
	err := db.QueryRow(`
		INSERT INTO recurring_journal_templates (
			tenant_id, restaurant_id, template_code, template_name, frequency,
			start_date, next_run_date, is_active, lines_json
		) VALUES ($1, $2, $3, $4, 'monthly', DATE '2026-01-01', DATE '2026-02-01', TRUE, '[{}]'::jsonb)
		RETURNING id
	`, tenantID, restaurantID, "ORIG-"+suffix, "Original "+suffix).Scan(&id)
	if err != nil {
		t.Fatalf("insert recurring template fixture: %v", err)
	}
	return id
}
