package main

import (
	"context"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// The hash below identifies accounts that were not rotated after the original migration 109.
const legacySeededSuperAdminHash = "$2a$10$6tgEoGE0dusUMtJ0RUd4VO3IACTDD2vb0676WfMyPvbVQU9oo2p06"

func TestIsSeededDefaultPassword(t *testing.T) {
	if !isSeededDefaultPassword(legacySeededSuperAdminHash) {
		t.Fatal("the legacy migration 109 hash must be rejected")
	}
	rotated, err := bcrypt.GenerateFromPassword([]byte("a-new-strong-password-1"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	if isSeededDefaultPassword(string(rotated)) {
		t.Fatal("a rotated password must not be flagged")
	}
	if isSeededDefaultPassword("") || isSeededDefaultPassword("not-a-hash") {
		t.Fatal("empty or invalid hashes must not be flagged")
	}
}

func TestProductionStartupGuards_NoDatabaseIsRejected(t *testing.T) {
	if err := productionStartupGuards(context.Background(), nil); err == nil {
		t.Fatal("production without a database (stub mode) must be refused")
	}
}
