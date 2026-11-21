//go:build integration

package repository

import (
	"context"
	"testing"
)

func TestReconcileJournalGL_PostedAndVoidedEntriesHaveExpectedEffects(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)

	entryID, err := fixture.ledger.CreateAndPostAutoEntry(fixture.entry(), 0)
	if err != nil {
		t.Fatalf("create synthetic posted journal: %v", err)
	}
	repo := NewJournalEntryRepository(db)
	got, err := repo.ReconcileJournalGL(context.Background(), fixture.tenantID, fixture.restaurantID)
	if err != nil {
		t.Fatalf("reconcile matched journal: %v", err)
	}
	if got.JournalLinesChecked != 2 || got.PostedLinesMissingGL != 0 || got.DuplicateOrInvalidGL != 0 || got.GLScopeOrAccountMismatch != 0 || got.GLDateMismatch != 0 || got.UnbalancedJournals != 0 {
		t.Fatalf("matched posted journal reconciliation = %+v", got)
	}

	if err := repo.VoidEntry(fixture.tenantID, fixture.restaurantID, entryID, 1, "synthetic reconciliation reversal"); err != nil {
		t.Fatalf("void synthetic journal: %v", err)
	}
	got, err = repo.ReconcileJournalGL(context.Background(), fixture.tenantID, fixture.restaurantID)
	if err != nil {
		t.Fatalf("reconcile reversed journal: %v", err)
	}
	if got.DuplicateOrInvalidGL != 0 || got.PostedLinesMissingGL != 0 || got.UnbalancedJournals != 0 {
		t.Fatalf("valid reversed journal reported as invalid: %+v", got)
	}
}

func TestReconcileJournalGL_IsTenantAndRestaurantScoped(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)
	other := openAtomicPostingRestaurantFixture(t, db, fixture.tenantID)
	defer cleanupAtomicPostingFixtureRows(t, db, other)

	if _, err := fixture.ledger.CreateAndPostAutoEntry(fixture.entry(), 0); err != nil {
		t.Fatalf("post first synthetic journal: %v", err)
	}
	if _, err := other.ledger.CreateAndPostAutoEntry(other.entry(), 0); err != nil {
		t.Fatalf("post sibling synthetic journal: %v", err)
	}

	got, err := NewJournalEntryRepository(db).ReconcileJournalGL(context.Background(), fixture.tenantID, fixture.restaurantID)
	if err != nil {
		t.Fatalf("reconcile scoped journals: %v", err)
	}
	if got.JournalLinesChecked != 2 || got.PostedLinesMissingGL != 0 || got.DuplicateOrInvalidGL != 0 {
		t.Fatalf("scoped reconciliation included unexpected data or exceptions: %+v", got)
	}
}
