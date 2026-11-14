//go:build integration

package usecase

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
	"pos-saas/tests/testdb"
)

func TestCFR010PostgresBalanceSheetAndEarningsReconcileByScope(t *testing.T) {
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(12)

	fixture := createCFR010Fixture(t, db)
	postCFR010Entry(t, db, fixture, "2025-01-01", "owner contribution", cfr010Line{"cash", 1000, 0}, cfr010Line{"capital", 0, 1000})
	postCFR010Entry(t, db, fixture, "2025-06-01", "prior period revenue", cfr010Line{"cash", 500, 0}, cfr010Line{"revenue", 0, 500})
	postCFR010Entry(t, db, fixture, "2025-06-05", "prior period expenses", cfr010Line{"expense", 300, 0}, cfr010Line{"cash", 0, 300})
	postCFR010Entry(t, db, fixture, "2025-06-10", "owner distribution", cfr010Line{"distribution", 50, 0}, cfr010Line{"cash", 0, 50})

	priorAsOf := cfr010Date("2025-06-30")
	priorBalanceSheet := generateCFR010BalanceSheet(t, db, fixture, priorAsOf)
	if priorBalanceSheet.TotalAssets != 1150 || priorBalanceSheet.TotalLiabilities != 0 || priorBalanceSheet.TotalEquity != 950 {
		t.Fatalf("pre-close balances = assets/liabilities/equity %.2f/%.2f/%.2f, want 1150/0/950", priorBalanceSheet.TotalAssets, priorBalanceSheet.TotalLiabilities, priorBalanceSheet.TotalEquity)
	}
	if priorBalanceSheet.EquationBalanced || priorBalanceSheet.EquationDifference != 200 {
		t.Fatalf("pre-close integrity = %t, difference %.2f; want unreconciled difference 200 for unclosed earnings", priorBalanceSheet.EquationBalanced, priorBalanceSheet.EquationDifference)
	}
	if amountForAccount(priorBalanceSheet.Equity, "3190") != -50 {
		t.Fatalf("owner distributions contribution = %.2f, want -50", amountForAccount(priorBalanceSheet.Equity, "3190"))
	}

	priorProfitLoss := generateCFR010ProfitLoss(t, db, fixture, cfr010Date("2025-01-01"), priorAsOf)
	if priorProfitLoss.TotalRevenue != 500 || priorProfitLoss.TotalExpenses != 300 || priorProfitLoss.NetProfit != 200 {
		t.Fatalf("prior period P&L = revenue/expenses/net %.2f/%.2f/%.2f, want 500/300/200", priorProfitLoss.TotalRevenue, priorProfitLoss.TotalExpenses, priorProfitLoss.NetProfit)
	}
	priorTrialBalance := generateCFR010TrialBalance(t, db, fixture, priorAsOf)
	if !priorTrialBalance.IsBalanced || priorTrialBalance.TotalDebits != 1500 || priorTrialBalance.TotalCredits != 1500 {
		t.Fatalf("prior trial balance = %.2f/%.2f balanced=%t, want 1500/1500/true", priorTrialBalance.TotalDebits, priorTrialBalance.TotalCredits, priorTrialBalance.IsBalanced)
	}

	// The explicit synthetic close transfers the prior-period net income into retained earnings.
	postCFR010Entry(t, db, fixture, "2026-01-01", "synthetic period close", cfr010Line{"revenue", 500, 0}, cfr010Line{"expense", 0, 300}, cfr010Line{"retained", 0, 200})
	closedBalanceSheet := generateCFR010BalanceSheet(t, db, fixture, cfr010Date("2026-01-01"))
	if !closedBalanceSheet.EquationBalanced || closedBalanceSheet.TotalAssets != 1150 || closedBalanceSheet.TotalEquity != 1150 {
		t.Fatalf("post-close balances = assets %.2f, equity %.2f, balanced=%t; want 1150/1150/true", closedBalanceSheet.TotalAssets, closedBalanceSheet.TotalEquity, closedBalanceSheet.EquationBalanced)
	}
	if amountForAccount(closedBalanceSheet.Equity, "3200") != 200 || amountForAccount(closedBalanceSheet.Equity, "3300") != 0 {
		t.Fatalf("post-close retained/current earnings = %.2f/%.2f, want 200/0", amountForAccount(closedBalanceSheet.Equity, "3200"), amountForAccount(closedBalanceSheet.Equity, "3300"))
	}

	postCFR010Entry(t, db, fixture, "2026-01-02", "current period revenue", cfr010Line{"cash", 100, 0}, cfr010Line{"revenue", 0, 100})
	postCFR010Entry(t, db, fixture, "2026-01-05", "current period net loss", cfr010Line{"expense", 120, 0}, cfr010Line{"cash", 0, 120})
	postCFR010Entry(t, db, fixture, "2026-01-08", "current period distribution", cfr010Line{"distribution", 30, 0}, cfr010Line{"cash", 0, 30})
	currentFrom, currentTo := cfr010Date("2026-01-02"), cfr010Date("2026-01-31")
	currentProfitLoss := generateCFR010ProfitLoss(t, db, fixture, currentFrom, currentTo)
	if currentProfitLoss.TotalRevenue != 100 || currentProfitLoss.TotalExpenses != 120 || currentProfitLoss.NetProfit != -20 {
		t.Fatalf("current period P&L = revenue/expenses/net %.2f/%.2f/%.2f, want 100/120/-20", currentProfitLoss.TotalRevenue, currentProfitLoss.TotalExpenses, currentProfitLoss.NetProfit)
	}
	currentAsOf := cfr010Date("2026-01-31")
	currentBalanceSheet := generateCFR010BalanceSheet(t, db, fixture, currentAsOf)
	if currentBalanceSheet.TotalAssets != 1100 || currentBalanceSheet.TotalEquity != 1120 || currentBalanceSheet.EquationDifference != -20 {
		t.Fatalf("open-period balances = assets %.2f, equity %.2f, difference %.2f; want 1100/1120/-20", currentBalanceSheet.TotalAssets, currentBalanceSheet.TotalEquity, currentBalanceSheet.EquationDifference)
	}
	if currentBalanceSheet.EquationBalanced || amountForAccount(currentBalanceSheet.Equity, "3190") != -80 {
		t.Fatalf("current balance sheet balanced=%t, distributions=%.2f; want unresolved earnings and distributions -80", currentBalanceSheet.EquationBalanced, amountForAccount(currentBalanceSheet.Equity, "3190"))
	}

	zeroPeriod := generateCFR010ProfitLoss(t, db, fixture, cfr010Date("2027-01-01"), cfr010Date("2027-12-31"))
	if zeroPeriod.TotalRevenue != 0 || zeroPeriod.TotalExpenses != 0 || zeroPeriod.NetProfit != 0 {
		t.Fatalf("zero-activity P&L = %.2f/%.2f/%.2f, want 0/0/0", zeroPeriod.TotalRevenue, zeroPeriod.TotalExpenses, zeroPeriod.NetProfit)
	}
	firstRepeat := generateCFR010BalanceSheet(t, db, fixture, currentAsOf)
	secondRepeat := generateCFR010BalanceSheet(t, db, fixture, currentAsOf)
	if !reflect.DeepEqual(firstRepeat, secondRepeat) {
		t.Fatalf("repeated report results differ: first=%+v second=%+v", firstRepeat, secondRepeat)
	}

	otherTenant := createCFR010Fixture(t, db)
	postCFR010Entry(t, db, otherTenant, "2025-01-01", "isolated other tenant", cfr010Line{"cash", 700, 0}, cfr010Line{"capital", 0, 700})
	siblingRestaurant := createCFR010RestaurantFixture(t, db, fixture.tenantID)
	postCFR010Entry(t, db, siblingRestaurant, "2025-01-01", "isolated sibling restaurant", cfr010Line{"cash", 900, 0}, cfr010Line{"capital", 0, 900})
	if isolated := generateCFR010BalanceSheet(t, db, fixture, currentAsOf); !reflect.DeepEqual(isolated, currentBalanceSheet) {
		t.Fatalf("other tenant/restaurant postings leaked into report: before=%+v after=%+v", currentBalanceSheet, isolated)
	}

	var postedCount, nonPostedCount int
	if err := db.QueryRow(`SELECT COUNT(*) FILTER (WHERE status = 'posted'), COUNT(*) FILTER (WHERE status <> 'posted') FROM journal_entries WHERE tenant_id = $1 AND restaurant_id = $2`, fixture.tenantID, fixture.restaurantID).Scan(&postedCount, &nonPostedCount); err != nil {
		t.Fatalf("verify report calls did not mutate journal states: %v", err)
	}
	if postedCount != 8 || nonPostedCount != 0 {
		t.Fatalf("journal states after reporting = posted %d / other %d, want 8/0", postedCount, nonPostedCount)
	}
}

type cfr010Fixture struct {
	tenantID, restaurantID int
	accounts               map[string]int
	accountCodes           map[string]string
}

type cfr010Line struct {
	account       string
	debit, credit float64
}

func createCFR010Fixture(t *testing.T, db *sql.DB) *cfr010Fixture {
	t.Helper()
	name := fmt.Sprintf("cfr010-%d", time.Now().UnixNano())
	var tenantID, restaurantID int
	if err := db.QueryRow(`INSERT INTO tenants (name, slug, domain) VALUES ($1, $2, $3) RETURNING id`, name, name, name+".example.test").Scan(&tenantID); err != nil {
		t.Fatalf("create isolated tenant: %v", err)
	}
	if err := db.QueryRow(`INSERT INTO restaurants (tenant_id, name, slug, address) VALUES ($1, $2, $3, 'test') RETURNING id`, tenantID, name, name).Scan(&restaurantID); err != nil {
		t.Fatalf("create isolated restaurant: %v", err)
	}
	fixture := &cfr010Fixture{tenantID: tenantID, restaurantID: restaurantID, accounts: make(map[string]int), accountCodes: make(map[string]string)}
	accounts := []struct{ key, code, label, kind, subtype, normal string }{
		{"cash", "1101", "Cash", "asset", "cash", "debit"},
		{"payable", "2101", "Accounts Payable", "liability", "payable", "credit"},
		{"capital", "3100", "Owner Capital", "equity", "capital", "credit"},
		{"retained", "3200", "Retained Earnings", "equity", "retained_earnings", "credit"},
		{"current", "3300", "Current Year Earnings", "equity", "current_year_earnings", "credit"},
		{"distribution", "3190", "Owner Distributions", "equity", "owner_distribution", "debit"},
		{"revenue", "4100", "Sales", "revenue", "food_sales", "credit"},
		{"expense", "5100", "Operating Expenses", "expense", "operating", "debit"},
	}
	for _, account := range accounts {
		var accountID int
		if err := db.QueryRow(`
			INSERT INTO chart_of_accounts (tenant_id, restaurant_id, account_code, account_name, account_type, account_subtype, normal_balance)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
		`, tenantID, restaurantID, account.code, account.label, account.kind, account.subtype, account.normal).Scan(&accountID); err != nil {
			t.Fatalf("create %s account: %v", account.key, err)
		}
		fixture.accounts[account.key] = accountID
		fixture.accountCodes[account.key] = account.code
	}
	return fixture
}

func createCFR010RestaurantFixture(t *testing.T, db *sql.DB, tenantID int) *cfr010Fixture {
	t.Helper()
	name := fmt.Sprintf("cfr010-sibling-%d", time.Now().UnixNano())
	var restaurantID int
	if err := db.QueryRow(`INSERT INTO restaurants (tenant_id, name, slug, address) VALUES ($1, $2, $3, 'test') RETURNING id`, tenantID, name, name).Scan(&restaurantID); err != nil {
		t.Fatalf("create sibling restaurant: %v", err)
	}
	fixture := &cfr010Fixture{tenantID: tenantID, restaurantID: restaurantID, accounts: make(map[string]int), accountCodes: make(map[string]string)}
	for _, account := range []struct{ key, code, label, kind, normal string }{
		{"cash", "1101", "Cash", "asset", "debit"},
		{"capital", "3100", "Owner Capital", "equity", "credit"},
	} {
		var accountID int
		if err := db.QueryRow(`
			INSERT INTO chart_of_accounts (tenant_id, restaurant_id, account_code, account_name, account_type, normal_balance)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
		`, tenantID, restaurantID, account.code, account.label, account.kind, account.normal).Scan(&accountID); err != nil {
			t.Fatalf("create sibling %s account: %v", account.key, err)
		}
		fixture.accounts[account.key] = accountID
		fixture.accountCodes[account.key] = account.code
	}
	return fixture
}

func postCFR010Entry(t *testing.T, db *sql.DB, fixture *cfr010Fixture, date, description string, lines ...cfr010Line) {
	t.Helper()
	entryDate := cfr010Date(date)
	entry := &domain.JournalEntry{
		TenantID: fixture.tenantID, RestaurantID: fixture.restaurantID, EntryDate: entryDate,
		Description: description, EntryType: "manual", Status: domain.JournalStatusDraft,
		Lines: make([]domain.JournalEntryLine, 0, len(lines)),
	}
	for index, line := range lines {
		accountID, ok := fixture.accounts[line.account]
		if !ok {
			t.Fatalf("fixture account %q does not exist", line.account)
		}
		entry.Lines = append(entry.Lines, domain.JournalEntryLine{
			AccountID: accountID, AccountCode: fixture.accountCodes[line.account],
			DebitAmount: line.debit, CreditAmount: line.credit, LineOrder: index + 1,
		})
	}
	journalRepository := repository.NewJournalEntryRepository(db)
	journalID, err := journalRepository.CreateJournalEntry(entry)
	if err != nil {
		t.Fatalf("create %s journal: %v", description, err)
	}
	if err := journalRepository.PostEntry(fixture.tenantID, fixture.restaurantID, journalID, 1); err != nil {
		t.Fatalf("post %s journal: %v", description, err)
	}
}

func generateCFR010BalanceSheet(t *testing.T, db *sql.DB, fixture *cfr010Fixture, asOf time.Time) BalanceSheetReport {
	t.Helper()
	report, err := NewFinancialReportUseCase(db, nil).GenerateBalanceSheet(fixture.tenantID, fixture.restaurantID, asOf)
	if err != nil {
		t.Fatalf("generate Balance Sheet: %v", err)
	}
	return report
}

func generateCFR010ProfitLoss(t *testing.T, db *sql.DB, fixture *cfr010Fixture, from, to time.Time) ProfitLossReport {
	t.Helper()
	report, err := NewFinancialReportUseCase(db, nil).GenerateProfitLoss(fixture.tenantID, fixture.restaurantID, from, to)
	if err != nil {
		t.Fatalf("generate Profit & Loss: %v", err)
	}
	return report
}

func generateCFR010TrialBalance(t *testing.T, db *sql.DB, fixture *cfr010Fixture, asOf time.Time) TrialBalanceReport {
	t.Helper()
	journalRepository := repository.NewJournalEntryRepository(db)
	ledgerRepository := repository.NewGeneralLedgerRepository(db, journalRepository)
	report, err := NewFinancialReportUseCase(db, ledgerRepository).GenerateTrialBalance(fixture.tenantID, fixture.restaurantID, asOf)
	if err != nil {
		t.Fatalf("generate Trial Balance: %v", err)
	}
	return report
}

func amountForAccount(lines []FinancialReportLine, accountCode string) float64 {
	for _, line := range lines {
		if line.AccountCode == accountCode {
			return line.Amount
		}
	}
	return 0
}

func cfr010Date(value string) time.Time {
	date, _ := time.Parse("2006-01-02", value)
	return date
}
