//go:build integration

package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/lib/pq"
	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
)

type financialBenchmarkFixture struct {
	tenantID, restaurantID, cashAccountID, revenueAccountID, inventoryAccountID, cogsAccountID, postedBy int
	productID                                                                                            int
	documentCount, salesOrderCount                                                                       int
	backdatedJournalCount                                                                                int
	postingDuration                                                                                      time.Duration
}

type financialBenchmarkTenant struct {
	tenantID, restaurantID, cashAccountID int
	journalCount                          int
}

type financialBenchmarkTenantCounts struct {
	journalEntries, journalLines, generalLedgerRows                            int
	journalDebitCents, journalCreditCents, ledgerDebitCents, ledgerCreditCents int64
}

func readFinancialBenchmarkTenantCounts(t *testing.T, db *sql.DB, tenantID, restaurantID int) financialBenchmarkTenantCounts {
	t.Helper()
	var counts financialBenchmarkTenantCounts
	if err := db.QueryRow(`
		SELECT COUNT(*), COALESCE(SUM(ROUND(total_debit * 100)::bigint), 0), COALESCE(SUM(ROUND(total_credit * 100)::bigint), 0)
		FROM journal_entries WHERE tenant_id=$1 AND restaurant_id=$2
	`, tenantID, restaurantID).Scan(&counts.journalEntries, &counts.journalDebitCents, &counts.journalCreditCents); err != nil {
		t.Fatalf("count benchmark journals for tenant %d: %v", tenantID, err)
	}
	if err := db.QueryRow(`
		SELECT COUNT(*) FROM journal_entry_lines jel
		JOIN journal_entries je ON je.id=jel.journal_entry_id
		WHERE je.tenant_id=$1 AND je.restaurant_id=$2
	`, tenantID, restaurantID).Scan(&counts.journalLines); err != nil {
		t.Fatalf("count benchmark journal lines for tenant %d: %v", tenantID, err)
	}
	if err := db.QueryRow(`
		SELECT COUNT(*), COALESCE(SUM(ROUND(gl.debit_amount * 100)::bigint), 0), COALESCE(SUM(ROUND(gl.credit_amount * 100)::bigint), 0)
		FROM general_ledger gl WHERE gl.tenant_id=$1 AND gl.restaurant_id=$2
	`, tenantID, restaurantID).Scan(&counts.generalLedgerRows, &counts.ledgerDebitCents, &counts.ledgerCreditCents); err != nil {
		t.Fatalf("count benchmark ledger rows for tenant %d: %v", tenantID, err)
	}
	return counts
}

func financialBenchmarkJournalCount(t *testing.T) int {
	t.Helper()
	count := 1000
	if raw := os.Getenv("POS_FINANCIAL_BENCHMARK_JOURNALS"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 100 || parsed > 50000 {
			t.Fatalf("POS_FINANCIAL_BENCHMARK_JOURNALS must be between 100 and 50000")
		}
		count = parsed
	}
	return count
}

func seedFinancialBenchmarkDataset(t *testing.T, db *sql.DB, journalCount int) financialBenchmarkFixture {
	t.Helper()
	base := seedLiveFinancialOrderFixture(t, db)
	fixture := financialBenchmarkFixture{tenantID: int(base.tenantID), restaurantID: int(base.restaurantID)}
	if err := db.QueryRow(`SELECT id FROM users WHERE tenant_id=$1 AND restaurant_id=$2 AND email=$3`, base.tenantID, base.restaurantID, base.email).Scan(&fixture.postedBy); err != nil {
		t.Fatalf("read synthetic posting actor: %v", err)
	}
	accountIDs := make(map[string]int64, 7)
	rows, err := db.Query(`SELECT account_code, id FROM chart_of_accounts WHERE tenant_id=$1 AND restaurant_id=$2 AND account_code IN ('1110','1120','1130','2101','4100','2120','5110')`, base.tenantID, base.restaurantID)
	if err != nil {
		t.Fatalf("read benchmark accounts: %v", err)
	}
	for rows.Next() {
		var code string
		var id int64
		if err := rows.Scan(&code, &id); err != nil {
			rows.Close()
			t.Fatalf("scan benchmark account: %v", err)
		}
		accountIDs[code] = id
	}
	if err := rows.Close(); err != nil {
		t.Fatalf("close benchmark account rows: %v", err)
	}
	if len(accountIDs) != 7 {
		t.Fatalf("benchmark accounts found=%d, want 7", len(accountIDs))
	}
	fixture.cashAccountID = int(accountIDs["1110"])
	fixture.revenueAccountID = int(accountIDs["4100"])
	fixture.inventoryAccountID = int(accountIDs["1130"])
	fixture.cogsAccountID = int(accountIDs["5110"])
	fixture.productID = int(base.productID)
	fixture.documentCount = minInt(journalCount, 2000)
	fixture.salesOrderCount = fixture.documentCount
	seedFinancialBenchmarkDocuments(t, db, fixture, accountIDs)

	tx, err := db.BeginTx(t.Context(), nil)
	if err != nil {
		t.Fatalf("begin synthetic financial seed: %v", err)
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`CREATE TEMP TABLE financial_benchmark_entries (ordinal INTEGER PRIMARY KEY, entry_number VARCHAR(20) NOT NULL, entry_date DATE NOT NULL, fiscal_period VARCHAR(7) NOT NULL) ON COMMIT DROP`); err != nil {
		t.Fatalf("create benchmark entry staging table: %v", err)
	}
	copyEntries, err := tx.Prepare(pq.CopyIn("financial_benchmark_entries", "ordinal", "entry_number", "entry_date", "fiscal_period"))
	if err != nil {
		t.Fatalf("prepare benchmark entry copy: %v", err)
	}
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < journalCount; i++ {
		entryDate := startDate.AddDate(0, 0, i*364/journalCount)
		if _, err := copyEntries.Exec(i, fmt.Sprintf("BM%017d", i), entryDate, entryDate.Format("2006-01")); err != nil {
			copyEntries.Close()
			t.Fatalf("stage benchmark journal %d: %v", i, err)
		}
	}
	if _, err := copyEntries.Exec(); err != nil {
		copyEntries.Close()
		t.Fatalf("flush benchmark journal copy: %v", err)
	}
	if err := copyEntries.Close(); err != nil {
		t.Fatalf("close benchmark journal copy: %v", err)
	}
	if _, err := tx.Exec(`
		INSERT INTO journal_entries (tenant_id, restaurant_id, entry_number, entry_date, fiscal_period, description, reference_type, entry_type, status)
		SELECT $1, $2, entry_number, entry_date, fiscal_period, 'Synthetic performance fixture', 'financial_benchmark', 'manual', 'draft'
		FROM financial_benchmark_entries ORDER BY ordinal`, base.tenantID, base.restaurantID); err != nil {
		t.Fatalf("insert balanced benchmark journals: %v", err)
	}
	if _, err := tx.Exec(`
		INSERT INTO journal_entry_lines (journal_entry_id, tenant_id, account_id, account_code, debit_amount, credit_amount, description, line_order)
		SELECT je.id, je.tenant_id, line.account_id, line.account_code, line.debit_amount, line.credit_amount, 'Synthetic benchmark line', line.line_order
		FROM financial_benchmark_entries staged
		JOIN journal_entries je ON je.tenant_id=$1 AND je.restaurant_id=$2 AND je.entry_number=staged.entry_number
		CROSS JOIN (VALUES ($3::bigint, '1110', 100::numeric, 0::numeric, 1), ($4::bigint, '4100', 0::numeric, 90::numeric, 2), ($5::bigint, '2120', 0::numeric, 10::numeric, 3), ($6::bigint, '5110', 10::numeric, 0::numeric, 4), ($7::bigint, '1130', 0::numeric, 10::numeric, 5)) AS line(account_id, account_code, debit_amount, credit_amount, line_order)`, base.tenantID, base.restaurantID, accountIDs["1110"], accountIDs["4100"], accountIDs["2120"], accountIDs["5110"], accountIDs["1130"]); err != nil {
		t.Fatalf("insert balanced benchmark journal lines: %v", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit synthetic financial draft data: %v", err)
	}
	entryRows, err := db.Query(`SELECT id FROM journal_entries WHERE tenant_id=$1 AND restaurant_id=$2 AND reference_type='financial_benchmark' ORDER BY entry_date, id`, base.tenantID, base.restaurantID)
	if err != nil {
		t.Fatalf("list benchmark drafts for posting: %v", err)
	}
	entryIDs := make([]int, 0, journalCount)
	for entryRows.Next() {
		var id int
		if err := entryRows.Scan(&id); err != nil {
			entryRows.Close()
			t.Fatalf("scan benchmark draft ID: %v", err)
		}
		entryIDs = append(entryIDs, id)
	}
	if err := entryRows.Close(); err != nil {
		t.Fatalf("close benchmark draft rows: %v", err)
	}
	if len(entryIDs) != journalCount {
		t.Fatalf("benchmark draft count=%d, want %d", len(entryIDs), journalCount)
	}
	postingTx, err := db.BeginTx(t.Context(), nil)
	if err != nil {
		t.Fatalf("begin supported journal posting transaction: %v", err)
	}
	defer postingTx.Rollback()
	journalRepository := repository.NewJournalEntryRepository(db)
	postingStart := time.Now()
	for _, id := range entryIDs {
		if err := journalRepository.PostEntryTx(postingTx, int(base.tenantID), int(base.restaurantID), id, fixture.postedBy); err != nil {
			t.Fatalf("post synthetic balanced journal %d through repository workflow: %v", id, err)
		}
	}
	if err := postingTx.Commit(); err != nil {
		t.Fatalf("commit synthetic financial postings: %v", err)
	}
	fixture.postingDuration = time.Since(postingStart)
	var entries, lines, ledgerRows int
	if err := db.QueryRow(`SELECT COUNT(*) FROM journal_entries WHERE tenant_id=$1 AND restaurant_id=$2 AND reference_type='financial_benchmark'`, base.tenantID, base.restaurantID).Scan(&entries); err != nil {
		t.Fatalf("verify benchmark journal count: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM journal_entry_lines jel JOIN journal_entries je ON je.id=jel.journal_entry_id WHERE je.tenant_id=$1 AND je.restaurant_id=$2 AND je.reference_type='financial_benchmark'`, base.tenantID, base.restaurantID).Scan(&lines); err != nil {
		t.Fatalf("verify benchmark line count: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM general_ledger WHERE tenant_id=$1 AND restaurant_id=$2 AND journal_entry_id IN (SELECT id FROM journal_entries WHERE reference_type='financial_benchmark' AND tenant_id=$1 AND restaurant_id=$2)`, base.tenantID, base.restaurantID).Scan(&ledgerRows); err != nil {
		t.Fatalf("verify benchmark GL count: %v", err)
	}
	if entries != journalCount || lines != journalCount*5 || ledgerRows != journalCount*5 {
		t.Fatalf("benchmark fixture counts entries/lines/GL=%d/%d/%d, expected %d/%d/%d", entries, lines, ledgerRows, journalCount, journalCount*5, journalCount*5)
	}
	var debitCents, creditCents int64
	if err := db.QueryRow(`SELECT ROUND(SUM(total_debit)*100)::bigint, ROUND(SUM(total_credit)*100)::bigint FROM journal_entries WHERE tenant_id=$1 AND restaurant_id=$2 AND reference_type='financial_benchmark'`, base.tenantID, base.restaurantID).Scan(&debitCents, &creditCents); err != nil {
		t.Fatalf("independently reconcile seeded journal totals: %v", err)
	}
	expectedTotalCents := int64(journalCount) * financialBenchmarkJournalTotalCents
	debitDifferenceCents := debitCents - expectedTotalCents
	creditDifferenceCents := creditCents - expectedTotalCents
	t.Logf("benchmark seed reconciliation cents: expected=%d observed_debit=%d debit_difference=%d observed_credit=%d credit_difference=%d", expectedTotalCents, debitCents, debitDifferenceCents, creditCents, creditDifferenceCents)
	if debitCents != creditCents || debitCents != expectedTotalCents {
		t.Fatalf("seeded journal debit/credit cents=%d/%d, expected=%d debit_difference=%d credit_difference=%d", debitCents, creditCents, expectedTotalCents, debitDifferenceCents, creditDifferenceCents)
	}
	return fixture
}

func seedFinancialBenchmarkEmptyRestaurant(t *testing.T, db *sql.DB, tenantID int) int {
	t.Helper()
	runID := strconv.FormatInt(time.Now().UnixNano(), 10)
	name := "financial-benchmark-empty-" + runID
	var restaurantID int
	if err := db.QueryRow(`INSERT INTO restaurants (tenant_id, name, slug, address) VALUES ($1, $2, $3, 'Synthetic benchmark scope') RETURNING id`, tenantID, name, name).Scan(&restaurantID); err != nil {
		t.Fatalf("seed empty reconciliation restaurant: %v", err)
	}
	return restaurantID
}

func seedFinancialBenchmarkDocuments(t *testing.T, db *sql.DB, fixture financialBenchmarkFixture, accountIDs map[string]int64) {
	t.Helper()
	var customerID, supplierID int64
	if err := db.QueryRow(`INSERT INTO customers_accounting (tenant_id, restaurant_id, customer_code, customer_name, receivable_account_id) VALUES ($1,$2,'BM-CUSTOMER','Benchmark customer',$3) RETURNING id`, fixture.tenantID, fixture.restaurantID, accountIDs["1120"]).Scan(&customerID); err != nil {
		t.Fatalf("create benchmark customer: %v", err)
	}
	if err := db.QueryRow(`INSERT INTO suppliers (tenant_id, restaurant_id, supplier_code, company_name, payable_account_id) VALUES ($1,$2,'BM-SUPPLIER','Benchmark supplier',$3) RETURNING id`, fixture.tenantID, fixture.restaurantID, accountIDs["2101"]).Scan(&supplierID); err != nil {
		t.Fatalf("create benchmark supplier: %v", err)
	}
	if _, err := db.Exec(`
		INSERT INTO invoices_customer (tenant_id, restaurant_id, invoice_number, invoice_date, due_date, customer_id, subtotal, total_amount, amount_paid, balance_due, status)
		SELECT $1, $2, 'BMC-' || lpad(seq::text, 16, '0'), DATE '2024-01-01' + ((seq * 17 % 730)::int),
		       DATE '2024-01-01' + ((seq * 17 % 730)::int) + ((seq % 120) - 60)::int, $3, 100, 100, 25, 75, 'partially_paid'
		FROM generate_series(1, $4) AS seq`, fixture.tenantID, fixture.restaurantID, customerID, fixture.documentCount); err != nil {
		t.Fatalf("seed benchmark customer invoices: %v", err)
	}
	if _, err := db.Exec(`
		INSERT INTO payments_received (tenant_id, restaurant_id, payment_number, payment_date, customer_id, invoice_id, amount, payment_method)
		SELECT $1, $2, 'BMR-' || lpad(inv.id::text, 16, '0'), inv.invoice_date + 10, $3, inv.id, 25, 'cash'
		FROM invoices_customer inv
		WHERE inv.tenant_id=$1 AND inv.restaurant_id=$2 AND inv.customer_id=$3
		ORDER BY inv.id`, fixture.tenantID, fixture.restaurantID, customerID); err != nil {
		t.Fatalf("seed benchmark customer receipts: %v", err)
	}
	if _, err := db.Exec(`
		INSERT INTO invoices_supplier (tenant_id, restaurant_id, invoice_number, supplier_ref, invoice_date, due_date, supplier_id, subtotal, total_amount, amount_paid, balance_due, status)
		SELECT $1, $2, 'BMS-' || lpad(seq::text, 16, '0'), 'SUP-' || seq::text,
		       DATE '2024-01-01' + ((seq * 29 % 730)::int), DATE '2024-01-01' + ((seq * 29 % 730)::int) + ((seq % 120) - 60)::int, $3, 100, 100, 40, 60, 'partially_paid'
		FROM generate_series(1, $4) AS seq`, fixture.tenantID, fixture.restaurantID, supplierID, fixture.documentCount); err != nil {
		t.Fatalf("seed benchmark supplier invoices: %v", err)
	}
	if _, err := db.Exec(`
		INSERT INTO payments_made (tenant_id, restaurant_id, payment_number, payment_date, supplier_id, invoice_id, amount, payment_method)
		SELECT $1, $2, 'BMP-' || lpad(inv.id::text, 16, '0'), inv.invoice_date + 15, $3, inv.id, 40, 'cash'
		FROM invoices_supplier inv
		WHERE inv.tenant_id=$1 AND inv.restaurant_id=$2 AND inv.supplier_id=$3
		ORDER BY inv.id`, fixture.tenantID, fixture.restaurantID, supplierID); err != nil {
		t.Fatalf("seed benchmark supplier payments: %v", err)
	}
	if _, err := db.Exec(`UPDATE products SET track_inventory=TRUE, quantity_in_stock=5000, cost=12.50 WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3`, fixture.productID, fixture.tenantID, fixture.restaurantID); err != nil {
		t.Fatalf("seed benchmark inventory valuation product: %v", err)
	}
	if _, err := db.Exec(`
		INSERT INTO inventory_movements (tenant_id, restaurant_id, product_id, reference_type, reference_id, movement_type, quantity, unit_cost, total_cost, notes, created_by, created_at)
		SELECT $1, $2, $3, 'financial_benchmark', seq, CASE WHEN seq % 5 = 0 THEN 'sale' ELSE 'purchase' END,
		       CASE WHEN seq % 5 = 0 THEN -1 ELSE 1 END, 12.50, 12.50, 'Synthetic benchmark inventory movement', $4,
		       TIMESTAMP '2020-01-01' + (seq * INTERVAL '1 day')
		FROM generate_series(1, $5) AS seq`, fixture.tenantID, fixture.restaurantID, fixture.productID, fixture.postedBy, fixture.documentCount); err != nil {
		t.Fatalf("seed benchmark inventory movements: %v", err)
	}
	seedFinancialBenchmarkSalesOrders(t, db, fixture)
}

func seedFinancialBenchmarkSalesOrders(t *testing.T, db *sql.DB, fixture financialBenchmarkFixture) {
	t.Helper()
	_, err := db.Exec(`
		INSERT INTO orders (
			tenant_id, restaurant_id, order_number, customer_name, customer_phone,
			subtotal, tax_amount, discount_amount, delivery_fee, total_amount,
			payment_method, payment_status, status, order_source, created_at, updated_at
		)
		SELECT $1, $2, 'FINBM-SALES-' || lpad(sequence::text, 8, '0'),
		       'Synthetic benchmark customer', '5550000000', 20, 0, 0, 0, 20,
		       'cash', 'paid', 'pending', 'website',
		       TIMESTAMP '2025-01-01' + (sequence % 365) * INTERVAL '1 day' + (sequence % 24) * INTERVAL '1 hour',
		       TIMESTAMP '2025-01-01' + (sequence % 365) * INTERVAL '1 day' + (sequence % 24) * INTERVAL '1 hour'
		FROM generate_series(1, $3) AS sequence`, fixture.tenantID, fixture.restaurantID, fixture.salesOrderCount)
	if err != nil {
		t.Fatalf("seed benchmark sales orders: %v", err)
	}
	if _, err := db.Exec(`
		INSERT INTO order_items (
			order_id, product_id, tenant_id, restaurant_id, product_name,
			quantity, unit_price, discount_amount, total_price, created_at, updated_at
		)
		SELECT orders.id, $3, orders.tenant_id, orders.restaurant_id,
		       'Financial benchmark product', 1, 10, 0, 10, orders.created_at, orders.created_at
		FROM orders
		CROSS JOIN generate_series(1, 2) AS item_number
		WHERE orders.tenant_id=$1 AND orders.restaurant_id=$2
		  AND orders.order_number LIKE 'FINBM-SALES-%'`, fixture.tenantID, fixture.restaurantID, fixture.productID); err != nil {
		t.Fatalf("seed benchmark sales items: %v", err)
	}
	// Payment status is synthetic to exercise the sales report/export query; these
	// rows do not represent payment/outbox or GL posting work. Fulfillment follows
	// the valid transition path so tracked inventory reservations remain realistic.
	for _, transition := range [][2]string{
		{"pending", "confirmed"},
		{"confirmed", "preparing"},
		{"preparing", "ready"},
		{"ready", "out_for_delivery"},
		{"out_for_delivery", "delivered"},
	} {
		if _, err := db.Exec(`
			UPDATE orders SET status=$1, updated_at=created_at
			WHERE tenant_id=$2 AND restaurant_id=$3 AND order_number LIKE 'FINBM-SALES-%' AND status=$4`,
			transition[1], fixture.tenantID, fixture.restaurantID, transition[0]); err != nil {
			t.Fatalf("advance benchmark sales orders from %s to %s: %v", transition[0], transition[1], err)
		}
	}
}

func seedFinancialBenchmarkSecondaryTenant(t *testing.T, db *sql.DB, journalCount int) financialBenchmarkTenant {
	t.Helper()
	base := seedLiveFinancialOrderFixture(t, db)
	var actorID, cashAccountID, revenueAccountID, taxAccountID int
	if err := db.QueryRow(`SELECT id FROM users WHERE tenant_id=$1 AND restaurant_id=$2 AND email=$3`, base.tenantID, base.restaurantID, base.email).Scan(&actorID); err != nil {
		t.Fatalf("read secondary benchmark actor: %v", err)
	}
	if err := db.QueryRow(`SELECT id FROM chart_of_accounts WHERE tenant_id=$1 AND restaurant_id=$2 AND account_code='1110'`, base.tenantID, base.restaurantID).Scan(&cashAccountID); err != nil {
		t.Fatalf("read secondary benchmark cash account: %v", err)
	}
	if err := db.QueryRow(`SELECT id FROM chart_of_accounts WHERE tenant_id=$1 AND restaurant_id=$2 AND account_code='4100'`, base.tenantID, base.restaurantID).Scan(&revenueAccountID); err != nil {
		t.Fatalf("read secondary benchmark revenue account: %v", err)
	}
	if err := db.QueryRow(`SELECT id FROM chart_of_accounts WHERE tenant_id=$1 AND restaurant_id=$2 AND account_code='2120'`, base.tenantID, base.restaurantID).Scan(&taxAccountID); err != nil {
		t.Fatalf("read secondary benchmark tax account: %v", err)
	}
	ledger := repository.NewGeneralLedgerRepository(db, repository.NewJournalEntryRepository(db))
	for sequence := 1; sequence <= journalCount; sequence++ {
		referenceID := sequence
		entry := &domain.JournalEntry{
			TenantID: int(base.tenantID), RestaurantID: int(base.restaurantID),
			EntryDate:   time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, sequence%365),
			Description: "Synthetic uneven-tenant benchmark journal", ReferenceType: "financial_benchmark_secondary",
			ReferenceID: &referenceID, EntryType: "manual", Status: domain.JournalStatusDraft,
			Lines: []domain.JournalEntryLine{
				{AccountID: cashAccountID, DebitAmount: 10, LineOrder: 1},
				{AccountID: revenueAccountID, CreditAmount: 9, LineOrder: 2},
				{AccountID: taxAccountID, CreditAmount: 1, LineOrder: 3},
			},
		}
		if _, err := ledger.CreateAndPostAutoEntry(entry, actorID); err != nil {
			t.Fatalf("post secondary tenant benchmark journal %d: %v", sequence, err)
		}
	}
	return financialBenchmarkTenant{tenantID: int(base.tenantID), restaurantID: int(base.restaurantID), cashAccountID: cashAccountID, journalCount: journalCount}
}
