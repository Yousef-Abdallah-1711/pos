package usecase

import (
	"testing"

	"pos-saas/internal/repository"
)

func TestAccountingEventFailureStateIsBounded(t *testing.T) {
	if got := accountingEventFailureStatus(&repository.AccountingEvent{Attempts: 1, MaxAttempts: 8}); got != repository.AccountingEventRetryable {
		t.Fatalf("status after transient failure = %q, want retryable", got)
	}
	if got := accountingEventFailureStatus(&repository.AccountingEvent{Attempts: 8, MaxAttempts: 8}); got != repository.AccountingEventFailed {
		t.Fatalf("status after final failure = %q, want failed", got)
	}
}
