package http

import "testing"

func TestValidateStringLengthsCountsUnicodeRunes(t *testing.T) {
	if err := validateStringLengths(stringLengthRule{Name: "name", Value: "مطعم", Max: 4}); err != nil {
		t.Fatalf("expected four Arabic characters to fit, got %v", err)
	}
}

func TestValidateStringLengthsRejectsOverLimit(t *testing.T) {
	if err := validateStringLengths(stringLengthRule{Name: "name", Value: "abc", Max: 2}); err == nil {
		t.Fatal("expected an over-limit string to be rejected")
	}
}

func TestValidateStringLengthsRejectsInvalidRule(t *testing.T) {
	if err := validateStringLengths(stringLengthRule{Name: "name", Value: "abc", Max: 0}); err == nil {
		t.Fatal("expected an invalid maximum length to be rejected")
	}
}
