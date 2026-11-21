package repository

import (
	"testing"
	"time"
)

func TestCompareCrossReportMovements_ReportsExactPerAccountDifferencesInStableOrder(t *testing.T) {
	result := emptyCrossReportDiagnostics(1, 2, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC))
	journals := map[int]crossReportMovement{
		9: {AccountID: 9, AccountCode: "4200", AccountType: "revenue", AmountCents: 5000},
		4: {AccountID: 4, AccountCode: "4100", AccountType: "revenue", AmountCents: 10000},
	}
	ledger := map[int]crossReportMovement{
		9: {AccountID: 9, AccountCode: "4200", AccountType: "revenue", AmountCents: 4500},
		4: {AccountID: 4, AccountCode: "4100", AccountType: "revenue", AmountCents: 10000},
	}
	compareCrossReportMovements(&result, journals, ledger)
	if len(result.ProfitLossAccountMismatches) != 1 {
		t.Fatalf("mismatches = %+v, want exactly one", result.ProfitLossAccountMismatches)
	}
	got := result.ProfitLossAccountMismatches[0]
	if got.AccountID != 9 || got.DifferenceCents != 500 || result.JournalRevenueCents != 15000 || result.LedgerRevenueCents != 14500 {
		t.Fatalf("exact cross-report comparison = %+v totals=%d/%d", got, result.JournalRevenueCents, result.LedgerRevenueCents)
	}
}
