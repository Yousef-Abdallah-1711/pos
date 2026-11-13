package usecase

import (
	"database/sql"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
)

type FinancialReportLine struct {
	AccountCode string  `json:"account_code"`
	AccountName string  `json:"account_name"`
	Amount      float64 `json:"amount"`
}

type ProfitLossReport struct {
	DateFrom                string                `json:"date_from"`
	DateTo                  string                `json:"date_to"`
	Revenue                 []FinancialReportLine `json:"revenue"`
	Expenses                []FinancialReportLine `json:"expenses"`
	TotalRevenue            float64               `json:"total_revenue"`
	TotalExpenses           float64               `json:"total_expenses"`
	NetProfit               float64               `json:"net_profit"`
	ComparisonDateFrom      string                `json:"comparison_date_from,omitempty"`
	ComparisonDateTo        string                `json:"comparison_date_to,omitempty"`
	ComparisonTotalRevenue  float64               `json:"comparison_total_revenue,omitempty"`
	ComparisonTotalExpenses float64               `json:"comparison_total_expenses,omitempty"`
	ComparisonNetProfit     float64               `json:"comparison_net_profit,omitempty"`
	RevenueChangePercent    float64               `json:"revenue_change_percent,omitempty"`
	ExpenseChangePercent    float64               `json:"expense_change_percent,omitempty"`
	NetProfitChangePercent  float64               `json:"net_profit_change_percent,omitempty"`
}

type BalanceSheetReport struct {
	AsOfDate            string                `json:"as_of_date"`
	Assets              []FinancialReportLine `json:"assets"`
	Liabilities         []FinancialReportLine `json:"liabilities"`
	Equity              []FinancialReportLine `json:"equity"`
	TotalAssets         float64               `json:"total_assets"`
	TotalLiabilities    float64               `json:"total_liabilities"`
	TotalEquity         float64               `json:"total_equity"`
	EquationBalanced    bool                  `json:"equation_balanced"`
	IntegrityStatus     string                `json:"integrity_status"`
	EquationDifference  float64               `json:"equation_difference"`
	IntegrityCheckBasis string                `json:"integrity_check_basis"`
}

type CashFlowReport struct {
	DateFrom          string  `json:"date_from"`
	DateTo            string  `json:"date_to"`
	OperatingCashFlow float64 `json:"operating_cash_flow"`
	InvestingCashFlow float64 `json:"investing_cash_flow"`
	FinancingCashFlow float64 `json:"financing_cash_flow"`
	NetCashFlow       float64 `json:"net_cash_flow"`
}

type AgingBucket struct {
	EntityID   int     `json:"entity_id"`
	EntityName string  `json:"entity_name"`
	Current    float64 `json:"current"`
	Days30     float64 `json:"days_30"`
	Days60     float64 `json:"days_60"`
	Days90     float64 `json:"days_90"`
	Days90Plus float64 `json:"days_90_plus"`
	Total      float64 `json:"total"`
}

type AgingReport struct {
	AsOfDate string        `json:"as_of_date"`
	Rows     []AgingBucket `json:"rows"`
	Total    float64       `json:"total"`
}

type InventoryValuationRow struct {
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    float64 `json:"quantity"`
	CostPerUnit float64 `json:"cost_per_unit"`
	TotalValue  float64 `json:"total_value"`
}

type InventoryValuationReport struct {
	AsOfDate   string                  `json:"as_of_date"`
	Rows       []InventoryValuationRow `json:"rows"`
	TotalValue float64                 `json:"total_value"`
}

type TrialBalanceRow struct {
	AccountID     int     `json:"account_id"`
	AccountCode   string  `json:"account_code"`
	AccountName   string  `json:"account_name"`
	AccountType   string  `json:"account_type"`
	NormalBalance string  `json:"normal_balance"`
	DebitBalance  float64 `json:"debit_balance"`
	CreditBalance float64 `json:"credit_balance"`
	NetBalance    float64 `json:"net_balance"`
}

type TrialBalanceReport struct {
	AsOfDate     string            `json:"as_of_date"`
	Rows         []TrialBalanceRow `json:"rows"`
	TotalDebits  float64           `json:"total_debits"`
	TotalCredits float64           `json:"total_credits"`
	IsBalanced   bool              `json:"is_balanced"`
}

type FinancialReportUseCase struct {
	db                *sql.DB
	generalLedgerRepo *repository.GeneralLedgerRepository
}

func NewFinancialReportUseCase(db *sql.DB, generalLedgerRepo *repository.GeneralLedgerRepository) *FinancialReportUseCase {
	return &FinancialReportUseCase{
		db:                db,
		generalLedgerRepo: generalLedgerRepo,
	}
}

func (uc *FinancialReportUseCase) GenerateProfitLoss(tenantID, restaurantID int, from, to time.Time) (ProfitLossReport, error) {
	report, err := uc.generateProfitLossCore(tenantID, restaurantID, from, to)
	if err != nil {
		return report, err
	}

	periodDays := int(to.Sub(from).Hours()/24) + 1
	if periodDays < 1 {
		periodDays = 1
	}
	previousTo := from.AddDate(0, 0, -1)
	previousFrom := previousTo.AddDate(0, 0, -(periodDays - 1))

	comparison, err := uc.generateProfitLossCore(tenantID, restaurantID, previousFrom, previousTo)
	if err != nil {
		return report, err
	}

	report.ComparisonDateFrom = comparison.DateFrom
	report.ComparisonDateTo = comparison.DateTo
	report.ComparisonTotalRevenue = comparison.TotalRevenue
	report.ComparisonTotalExpenses = comparison.TotalExpenses
	report.ComparisonNetProfit = comparison.NetProfit
	report.RevenueChangePercent = percentChange(comparison.TotalRevenue, report.TotalRevenue)
	report.ExpenseChangePercent = percentChange(comparison.TotalExpenses, report.TotalExpenses)
	report.NetProfitChangePercent = percentChange(comparison.NetProfit, report.NetProfit)
	return report, nil
}

func (uc *FinancialReportUseCase) GenerateProfitLossWithComparison(
	tenantID, restaurantID int,
	from, to, compareFrom, compareTo time.Time,
) (ProfitLossReport, error) {
	report, err := uc.generateProfitLossCore(tenantID, restaurantID, from, to)
	if err != nil {
		return report, err
	}
	comparison, err := uc.generateProfitLossCore(tenantID, restaurantID, compareFrom, compareTo)
	if err != nil {
		return report, err
	}

	report.ComparisonDateFrom = comparison.DateFrom
	report.ComparisonDateTo = comparison.DateTo
	report.ComparisonTotalRevenue = comparison.TotalRevenue
	report.ComparisonTotalExpenses = comparison.TotalExpenses
	report.ComparisonNetProfit = comparison.NetProfit
	report.RevenueChangePercent = percentChange(comparison.TotalRevenue, report.TotalRevenue)
	report.ExpenseChangePercent = percentChange(comparison.TotalExpenses, report.TotalExpenses)
	report.NetProfitChangePercent = percentChange(comparison.NetProfit, report.NetProfit)
	return report, nil
}

func (uc *FinancialReportUseCase) generateProfitLossCore(tenantID, restaurantID int, from, to time.Time) (ProfitLossReport, error) {
	report := ProfitLossReport{
		DateFrom: from.Format("2006-01-02"),
		DateTo:   to.Format("2006-01-02"),
		Revenue:  make([]FinancialReportLine, 0),
		Expenses: make([]FinancialReportLine, 0),
	}

	rows, err := uc.db.Query(`
		SELECT
			coa.account_code,
			coa.account_name,
			coa.account_type,
			COALESCE(SUM(jel.debit_amount), 0) AS debit_total,
			COALESCE(SUM(jel.credit_amount), 0) AS credit_total
		FROM journal_entry_lines jel
		JOIN journal_entries je ON je.id = jel.journal_entry_id
		JOIN chart_of_accounts coa ON coa.id = jel.account_id
		WHERE je.tenant_id = $1
		  AND je.restaurant_id = $2
		  AND je.status = 'posted'
		  AND je.entry_date >= $3
		  AND je.entry_date <= $4
		  AND coa.account_type IN ('revenue', 'expense')
		GROUP BY coa.account_code, coa.account_name, coa.account_type
		ORDER BY coa.account_code ASC
	`, tenantID, restaurantID, from, to)
	if err != nil {
		return report, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			line        FinancialReportLine
			accountType string
			debitTotal  float64
			creditTotal float64
		)
		if err := rows.Scan(
			&line.AccountCode,
			&line.AccountName,
			&accountType,
			&debitTotal,
			&creditTotal,
		); err != nil {
			return report, err
		}

		if accountType == "revenue" {
			line.Amount = creditTotal - debitTotal
			report.TotalRevenue += line.Amount
			report.Revenue = append(report.Revenue, line)
		} else {
			line.Amount = debitTotal - creditTotal
			report.TotalExpenses += line.Amount
			report.Expenses = append(report.Expenses, line)
		}
	}
	if err := rows.Err(); err != nil {
		return report, err
	}

	report.NetProfit = report.TotalRevenue - report.TotalExpenses
	return report, nil
}

func (uc *FinancialReportUseCase) GenerateTrialBalance(tenantID, restaurantID int, asOf time.Time) (TrialBalanceReport, error) {
	report := TrialBalanceReport{
		AsOfDate: asOf.Format("2006-01-02"),
		Rows:     make([]TrialBalanceRow, 0),
	}
	if uc.generalLedgerRepo == nil {
		return report, sql.ErrConnDone
	}

	rows, err := uc.generalLedgerRepo.GetTrialBalance(tenantID, restaurantID, asOf)
	if err != nil {
		return report, err
	}

	for _, item := range rows {
		row := mapTrialBalanceRow(item)
		report.TotalDebits += row.DebitBalance
		report.TotalCredits += row.CreditBalance
		report.Rows = append(report.Rows, row)
	}
	report.IsBalanced = round2(report.TotalDebits) == round2(report.TotalCredits)
	return report, nil
}

func (uc *FinancialReportUseCase) GenerateBalanceSheet(tenantID, restaurantID int, asOf time.Time) (BalanceSheetReport, error) {
	report := BalanceSheetReport{
		AsOfDate:    asOf.Format("2006-01-02"),
		Assets:      make([]FinancialReportLine, 0),
		Liabilities: make([]FinancialReportLine, 0),
		Equity:      make([]FinancialReportLine, 0),
	}

	rows, err := uc.db.Query(`
		WITH latest AS (
			SELECT DISTINCT ON (gl.account_id)
				gl.account_id,
				gl.running_balance
			FROM general_ledger gl
			WHERE gl.tenant_id = $1
			  AND gl.restaurant_id = $2
			  AND gl.entry_date <= $3
			ORDER BY gl.account_id, gl.entry_date DESC, gl.id DESC
		)
		SELECT
			coa.account_code,
			coa.account_name,
			coa.account_type,
			coa.normal_balance,
			COALESCE(latest.running_balance, 0) AS running_balance
		FROM chart_of_accounts coa
		LEFT JOIN latest ON latest.account_id = coa.id
		WHERE coa.tenant_id = $1
		  AND coa.restaurant_id = $2
		  AND coa.is_active = TRUE
		  AND coa.is_header = FALSE
		  AND coa.account_type IN ('asset', 'liability', 'equity')
		ORDER BY coa.account_code ASC
	`, tenantID, restaurantID, asOf)
	if err != nil {
		return report, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			line          FinancialReportLine
			accountType   string
			normalBalance string
			running       float64
		)
		if err := rows.Scan(&line.AccountCode, &line.AccountName, &accountType, &normalBalance, &running); err != nil {
			return report, err
		}

		amount := running
		if normalBalance == "credit" {
			amount = -running
		}
		if amount < 0 {
			amount = -amount
		}
		line.Amount = amount

		switch accountType {
		case "asset":
			report.Assets = append(report.Assets, line)
			report.TotalAssets += amount
		case "liability":
			report.Liabilities = append(report.Liabilities, line)
			report.TotalLiabilities += amount
		case "equity":
			report.Equity = append(report.Equity, line)
			report.TotalEquity += amount
		}
	}
	if err := rows.Err(); err != nil {
		return report, err
	}

	// Keep the existing equation rule during Phase 00. This surfaces its result without
	// changing report totals or introducing a new rounding/accounting policy.
	report.EquationBalanced = round2(report.TotalAssets) == round2(report.TotalLiabilities+report.TotalEquity)
	report.EquationDifference = report.TotalAssets - report.TotalLiabilities - report.TotalEquity
	report.IntegrityCheckBasis = "existing_round2_equation_check"
	if report.EquationBalanced {
		report.IntegrityStatus = "RECONCILED"
	} else {
		report.IntegrityStatus = "UNRECONCILED"
	}
	return report, nil
}

func (uc *FinancialReportUseCase) GenerateCashFlow(tenantID, restaurantID int, from, to time.Time) (CashFlowReport, error) {
	report := CashFlowReport{
		DateFrom: from.Format("2006-01-02"),
		DateTo:   to.Format("2006-01-02"),
	}

	if err := uc.db.QueryRow(`
		SELECT
			COALESCE(SUM(CASE WHEN coa.account_type = 'revenue' THEN jel.credit_amount - jel.debit_amount ELSE 0 END), 0)
			-
			COALESCE(SUM(CASE WHEN coa.account_type = 'expense' THEN jel.debit_amount - jel.credit_amount ELSE 0 END), 0) AS operating
		FROM journal_entry_lines jel
		JOIN journal_entries je ON je.id = jel.journal_entry_id
		JOIN chart_of_accounts coa ON coa.id = jel.account_id
		WHERE je.tenant_id = $1
		  AND je.restaurant_id = $2
		  AND je.status = 'posted'
		  AND je.entry_date >= $3
		  AND je.entry_date <= $4
	`, tenantID, restaurantID, from, to).Scan(&report.OperatingCashFlow); err != nil {
		return report, err
	}

	if err := uc.db.QueryRow(`
		SELECT
			COALESCE(SUM(jel.debit_amount - jel.credit_amount), 0)
		FROM journal_entry_lines jel
		JOIN journal_entries je ON je.id = jel.journal_entry_id
		JOIN chart_of_accounts coa ON coa.id = jel.account_id
		WHERE je.tenant_id = $1
		  AND je.restaurant_id = $2
		  AND je.status = 'posted'
		  AND je.entry_date >= $3
		  AND je.entry_date <= $4
		  AND coa.account_code LIKE '12%%'
	`, tenantID, restaurantID, from, to).Scan(&report.InvestingCashFlow); err != nil {
		return report, err
	}

	if err := uc.db.QueryRow(`
		SELECT
			COALESCE(SUM(jel.credit_amount - jel.debit_amount), 0)
		FROM journal_entry_lines jel
		JOIN journal_entries je ON je.id = jel.journal_entry_id
		JOIN chart_of_accounts coa ON coa.id = jel.account_id
		WHERE je.tenant_id = $1
		  AND je.restaurant_id = $2
		  AND je.status = 'posted'
		  AND je.entry_date >= $3
		  AND je.entry_date <= $4
		  AND (coa.account_type = 'equity' OR coa.account_code LIKE '22%%')
	`, tenantID, restaurantID, from, to).Scan(&report.FinancingCashFlow); err != nil {
		return report, err
	}

	report.NetCashFlow = report.OperatingCashFlow + report.InvestingCashFlow + report.FinancingCashFlow
	return report, nil
}

func (uc *FinancialReportUseCase) GenerateReceivablesAging(tenantID, restaurantID int, asOf time.Time) (AgingReport, error) {
	return uc.generateAging(
		tenantID,
		restaurantID,
		asOf,
		"invoices_customer",
		"customer_id",
		"customers_accounting",
		"customer_name",
		"customer_name",
	)
}

func (uc *FinancialReportUseCase) GeneratePayablesAging(tenantID, restaurantID int, asOf time.Time) (AgingReport, error) {
	return uc.generateAging(
		tenantID,
		restaurantID,
		asOf,
		"invoices_supplier",
		"supplier_id",
		"suppliers",
		"company_name",
		"supplier_name",
	)
}

func (uc *FinancialReportUseCase) generateAging(
	tenantID, restaurantID int,
	asOf time.Time,
	invoiceTable, entityIDColumn, entityTable, entityNameColumn, outEntityName string,
) (AgingReport, error) {
	report := AgingReport{
		AsOfDate: asOf.Format("2006-01-02"),
		Rows:     make([]AgingBucket, 0),
	}

	query := `
		SELECT
			inv.` + entityIDColumn + `,
			COALESCE(ent.` + entityNameColumn + `, 'Unknown') AS entity_name,
			COALESCE(SUM(CASE WHEN inv.due_date >= $3 THEN inv.balance_due ELSE 0 END), 0) AS current_bucket,
			COALESCE(SUM(CASE WHEN inv.due_date < $3 AND inv.due_date >= $3 - INTERVAL '30 day' THEN inv.balance_due ELSE 0 END), 0) AS bucket_30,
			COALESCE(SUM(CASE WHEN inv.due_date < $3 - INTERVAL '30 day' AND inv.due_date >= $3 - INTERVAL '60 day' THEN inv.balance_due ELSE 0 END), 0) AS bucket_60,
			COALESCE(SUM(CASE WHEN inv.due_date < $3 - INTERVAL '60 day' AND inv.due_date >= $3 - INTERVAL '90 day' THEN inv.balance_due ELSE 0 END), 0) AS bucket_90,
			COALESCE(SUM(CASE WHEN inv.due_date < $3 - INTERVAL '90 day' THEN inv.balance_due ELSE 0 END), 0) AS bucket_90_plus
		FROM ` + invoiceTable + ` inv
		LEFT JOIN ` + entityTable + ` ent ON ent.id = inv.` + entityIDColumn + `
		WHERE inv.tenant_id = $1
		  AND inv.restaurant_id = $2
		  AND inv.status <> 'void'
		  AND inv.balance_due > 0
		GROUP BY inv.` + entityIDColumn + `, COALESCE(ent.` + entityNameColumn + `, 'Unknown')
		ORDER BY entity_name ASC
	`

	rows, err := uc.db.Query(query, tenantID, restaurantID, asOf)
	if err != nil {
		return report, err
	}
	defer rows.Close()

	for rows.Next() {
		var row AgingBucket
		if err := rows.Scan(
			&row.EntityID,
			&row.EntityName,
			&row.Current,
			&row.Days30,
			&row.Days60,
			&row.Days90,
			&row.Days90Plus,
		); err != nil {
			return report, err
		}
		row.Total = row.Current + row.Days30 + row.Days60 + row.Days90 + row.Days90Plus
		report.Total += row.Total
		report.Rows = append(report.Rows, row)
	}
	if err := rows.Err(); err != nil {
		return report, err
	}

	_ = outEntityName
	return report, nil
}

func (uc *FinancialReportUseCase) GenerateInventoryValuation(tenantID, restaurantID int, asOf time.Time) (InventoryValuationReport, error) {
	report := InventoryValuationReport{
		AsOfDate: asOf.Format("2006-01-02"),
		Rows:     make([]InventoryValuationRow, 0),
	}

	rows, err := uc.db.Query(`
		SELECT
			id,
			COALESCE(NULLIF(TRIM(name_en), ''), NULLIF(TRIM(name), ''), 'Product #' || id::text) AS product_name,
			COALESCE(quantity_in_stock, 0)::DOUBLE PRECISION AS quantity_in_stock,
			COALESCE(cost, 0)::DOUBLE PRECISION AS cost_per_unit
		FROM products
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND COALESCE(track_inventory, false) = true
		  AND COALESCE(status, 'active') <> 'deleted'
		ORDER BY product_name ASC
	`, tenantID, restaurantID)
	if err != nil {
		return report, err
	}
	defer rows.Close()

	for rows.Next() {
		var row InventoryValuationRow
		if err := rows.Scan(&row.ProductID, &row.ProductName, &row.Quantity, &row.CostPerUnit); err != nil {
			return report, err
		}
		row.TotalValue = row.Quantity * row.CostPerUnit
		report.TotalValue += row.TotalValue
		report.Rows = append(report.Rows, row)
	}
	if err := rows.Err(); err != nil {
		return report, err
	}

	return report, nil
}

func round2(value float64) float64 {
	return float64(int(value*100+0.5)) / 100
}

func percentChange(previous, current float64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}
	return ((current - previous) / previous) * 100
}

func mapTrialBalanceRow(in domain.TrialBalanceRow) TrialBalanceRow {
	return TrialBalanceRow{
		AccountID:     in.AccountID,
		AccountCode:   in.AccountCode,
		AccountName:   in.AccountName,
		AccountType:   in.AccountType,
		NormalBalance: in.NormalBalance,
		DebitBalance:  in.DebitBalance,
		CreditBalance: in.CreditBalance,
		NetBalance:    in.NetBalance,
	}
}
