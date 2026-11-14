//go:build integration

package http

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	stdhttp "net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"pos-saas/internal/repository"
	"pos-saas/internal/usecase"
	"pos-saas/tests/testdb"
)

type payrollFixture struct {
	db           *sql.DB
	tenantID     int
	restaurantID int
	salaryID     int
	accounting   *usecase.AutoPostingEngine
}

func openPayrollFixture(t *testing.T) *payrollFixture {
	t.Helper()
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	db.SetMaxOpenConns(20)
	if err := db.Ping(); err != nil {
		db.Close()
		t.Fatalf("connect to disposable PostgreSQL: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	fixture := &payrollFixture{db: db}
	runID := fmt.Sprintf("%d", time.Now().UnixNano())
	if err := db.QueryRow(`INSERT INTO tenants (name, slug, domain) VALUES ($1,$2,$3) RETURNING id`, "cfr009-"+runID, "cfr009-"+runID, runID+".example.test").Scan(&fixture.tenantID); err != nil {
		t.Fatalf("create tenant fixture: %v", err)
	}
	if err := db.QueryRow(`INSERT INTO restaurants (tenant_id,name,slug,address) VALUES ($1,'CFR-009 restaurant',$2,'synthetic') RETURNING id`, fixture.tenantID, "cfr009-"+runID).Scan(&fixture.restaurantID); err != nil {
		t.Fatalf("create restaurant fixture: %v", err)
	}
	var employeeID int
	if err := db.QueryRow(`
		INSERT INTO employees (tenant_id,restaurant_id,employee_code,first_name,last_name,email,hire_date,employment_type,position)
		VALUES ($1,$2,$3,'Test','Employee',$4,CURRENT_DATE,'full_time','CFR-009 fixture') RETURNING id
	`, fixture.tenantID, fixture.restaurantID, "CFR009-"+runID, runID+"@example.test").Scan(&employeeID); err != nil {
		t.Fatalf("create employee fixture: %v", err)
	}
	if err := db.QueryRow(`
		INSERT INTO salaries (tenant_id,restaurant_id,employee_id,pay_period_start,pay_period_end,month,year,base_salary,gross_salary,net_salary,status,is_paid)
		VALUES ($1,$2,$3,DATE '2026-10-01',DATE '2026-10-31',10,2026,1000,1000,1000,'pending',FALSE) RETURNING id
	`, fixture.tenantID, fixture.restaurantID, employeeID).Scan(&fixture.salaryID); err != nil {
		t.Fatalf("create salary fixture: %v", err)
	}
	for _, account := range []struct {
		code, name, accountType, normalBalance string
	}{
		{"1111", "CFR-009 bank", "asset", "debit"},
		{"5210", "CFR-009 payroll expense", "expense", "debit"},
	} {
		if _, err := db.Exec(`INSERT INTO chart_of_accounts (tenant_id,restaurant_id,account_code,account_name,account_type,normal_balance) VALUES ($1,$2,$3,$4,$5,$6)`, fixture.tenantID, fixture.restaurantID, account.code, account.name, account.accountType, account.normalBalance); err != nil {
			t.Fatalf("create account %s: %v", account.code, err)
		}
	}
	journalRepo := repository.NewJournalEntryRepository(db)
	fixture.accounting = usecase.NewAutoPostingEngine(
		repository.NewAutoPostingRepository(db), repository.NewAccountRepository(db),
		repository.NewInventoryMovementRepository(db), repository.NewGeneralLedgerRepository(db, journalRepo),
		repository.NewProductRepository(db),
	)
	return fixture
}

func (f *payrollFixture) handler() *HRSalaryHandler {
	handler := NewHRSalaryHandler(repository.NewSalaryRepository(f.db), f.db)
	handler.SetAutoPostingEngine(f.accounting)
	return handler
}

func (f *payrollFixture) markPaid(t *testing.T, reference string) *httptest.ResponseRecorder {
	t.Helper()
	body, err := json.Marshal(map[string]string{"payment_reference": reference})
	if err != nil {
		t.Fatalf("encode payment request: %v", err)
	}
	req := httptest.NewRequest(stdhttp.MethodPost, "/api/v1/hr/salaries/"+fmt.Sprint(f.salaryID)+"/mark-paid", bytes.NewReader(body))
	req.SetPathValue("id", fmt.Sprint(f.salaryID))
	req = withAccountingAuthContext(req, int64(f.tenantID), int64(f.restaurantID), 0)
	recorder := httptest.NewRecorder()
	f.handler().MarkAsPaid(recorder, req)
	return recorder
}

func TestCFR009PostgresAllowsUnpaidEditAndDelete(t *testing.T) {
	fixture := openPayrollFixture(t)
	repo := repository.NewSalaryRepository(fixture.db)
	salary, err := repo.GetSalaryByID(fixture.tenantID, fixture.restaurantID, fixture.salaryID)
	if err != nil || salary == nil {
		t.Fatalf("load unpaid salary: salary=%v err=%v", salary, err)
	}
	salary.BaseSalary = 1250
	if err := repo.UpdateSalary(salary); err != nil {
		t.Fatalf("update unpaid salary: %v", err)
	}
	if err := repo.DeleteUnpaidSalary(fixture.tenantID, fixture.restaurantID, fixture.salaryID); err != nil {
		t.Fatalf("delete unpaid salary: %v", err)
	}
	var remaining int
	if err := fixture.db.QueryRow(`SELECT COUNT(*) FROM salaries WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3`, fixture.salaryID, fixture.tenantID, fixture.restaurantID).Scan(&remaining); err != nil || remaining != 0 {
		t.Fatalf("unpaid salary remains: count=%d err=%v", remaining, err)
	}
}

func TestCFR009PostgresPaidSalaryIsImmutableAndPostingIsIdempotent(t *testing.T) {
	fixture := openPayrollFixture(t)
	first := fixture.markPaid(t, "CFR009-PAY-1")
	if first.Code != stdhttp.StatusOK {
		t.Fatalf("first mark-paid status=%d body=%s", first.Code, first.Body.String())
	}
	second := fixture.markPaid(t, "CFR009-PAY-2")
	if second.Code != stdhttp.StatusOK {
		t.Fatalf("repeat mark-paid status=%d body=%s", second.Code, second.Body.String())
	}
	repo := repository.NewSalaryRepository(fixture.db)
	salary, err := repo.GetSalaryByID(fixture.tenantID, fixture.restaurantID, fixture.salaryID)
	if err != nil || salary == nil || salary.PaymentReference != "CFR009-PAY-1" || !salary.IsPaid {
		t.Fatalf("repeat request changed payment metadata: salary=%+v err=%v", salary, err)
	}
	salary.BaseSalary++
	if err := repo.UpdateSalary(salary); err != repository.ErrSalaryPaid {
		t.Fatalf("paid salary update error=%v, want immutable conflict", err)
	}
	if err := repo.DeleteUnpaidSalary(fixture.tenantID, fixture.restaurantID, fixture.salaryID); err != repository.ErrSalaryPaid {
		t.Fatalf("paid salary delete error=%v, want immutable conflict", err)
	}
	updateRequest := httptest.NewRequest(stdhttp.MethodPut, "/api/v1/hr/salaries/"+fmt.Sprint(fixture.salaryID), bytes.NewBufferString(`{"base_salary":2000}`))
	updateRequest.SetPathValue("id", fmt.Sprint(fixture.salaryID))
	updateRequest = withAccountingAuthContext(updateRequest, int64(fixture.tenantID), int64(fixture.restaurantID), 0)
	updateResponse := httptest.NewRecorder()
	fixture.handler().UpdateSalary(updateResponse, updateRequest)
	if updateResponse.Code != stdhttp.StatusConflict {
		t.Fatalf("paid salary update HTTP status=%d body=%s, want 409", updateResponse.Code, updateResponse.Body.String())
	}
	deleteRequest := httptest.NewRequest(stdhttp.MethodDelete, "/api/v1/hr/salaries/"+fmt.Sprint(fixture.salaryID), nil)
	deleteRequest.SetPathValue("id", fmt.Sprint(fixture.salaryID))
	deleteRequest = withAccountingAuthContext(deleteRequest, int64(fixture.tenantID), int64(fixture.restaurantID), 0)
	deleteResponse := httptest.NewRecorder()
	fixture.handler().DeleteSalary(deleteResponse, deleteRequest)
	if deleteResponse.Code != stdhttp.StatusConflict {
		t.Fatalf("paid salary delete HTTP status=%d body=%s, want 409", deleteResponse.Code, deleteResponse.Body.String())
	}
	assertPayrollPostingCounts(t, fixture, 1, 2)
	var transitions int
	if err := fixture.db.QueryRow(`SELECT COUNT(*) FROM hr_audit_logs WHERE tenant_id=$1 AND restaurant_id=$2 AND entity_type='salary' AND entity_id=$3 AND action='update'`, fixture.tenantID, fixture.restaurantID, fixture.salaryID).Scan(&transitions); err != nil || transitions != 1 {
		t.Fatalf("payment transition audit count=%d err=%v, want 1", transitions, err)
	}
}

func TestCFR009PostgresConcurrentMarkPaidCreatesOnePaymentAndJournal(t *testing.T) {
	fixture := openPayrollFixture(t)
	const attempts = 8
	start := make(chan struct{})
	results := make(chan *httptest.ResponseRecorder, attempts)
	var workers sync.WaitGroup
	for i := range attempts {
		workers.Add(1)
		go func(i int) {
			defer workers.Done()
			<-start
			results <- fixture.markPaid(t, fmt.Sprintf("CFR009-CONCURRENT-%d", i))
		}(i)
	}
	close(start)
	workers.Wait()
	close(results)
	for result := range results {
		if result.Code != stdhttp.StatusOK {
			t.Fatalf("concurrent mark-paid status=%d body=%s", result.Code, result.Body.String())
		}
	}
	assertPayrollPostingCounts(t, fixture, 1, 2)
	var transitions int
	if err := fixture.db.QueryRow(`SELECT COUNT(*) FROM hr_audit_logs WHERE tenant_id=$1 AND restaurant_id=$2 AND entity_type='salary' AND entity_id=$3 AND action='update'`, fixture.tenantID, fixture.restaurantID, fixture.salaryID).Scan(&transitions); err != nil || transitions != 1 {
		t.Fatalf("concurrent transition audit count=%d err=%v, want 1", transitions, err)
	}
}

func TestCFR009PostgresPostingFailureRollsBackPaymentState(t *testing.T) {
	fixture := openPayrollFixture(t)
	triggerName := fmt.Sprintf("cfr009_fail_journal_%d", time.Now().UnixNano())
	functionName := fmt.Sprintf("cfr009_fail_journal_fn_%d", time.Now().UnixNano())
	if _, err := fixture.db.Exec(fmt.Sprintf(`CREATE FUNCTION %s() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN IF NEW.reference_type='salary' THEN RAISE EXCEPTION 'injected payroll journal failure'; END IF; RETURN NEW; END $$`, functionName)); err != nil {
		t.Fatalf("create failure-injection function: %v", err)
	}
	if _, err := fixture.db.Exec(fmt.Sprintf(`CREATE TRIGGER %s BEFORE INSERT ON journal_entries FOR EACH ROW EXECUTE FUNCTION %s()`, triggerName, functionName)); err != nil {
		t.Fatalf("create failure-injection trigger: %v", err)
	}
	defer fixture.db.Exec(fmt.Sprintf(`DROP TRIGGER IF EXISTS %s ON journal_entries`, triggerName))
	defer fixture.db.Exec(fmt.Sprintf(`DROP FUNCTION IF EXISTS %s()`, functionName))

	result := fixture.markPaid(t, "CFR009-FAILURE")
	if result.Code != stdhttp.StatusInternalServerError {
		t.Fatalf("posting failure status=%d body=%s, want 500", result.Code, result.Body.String())
	}
	var paid bool
	var status string
	if err := fixture.db.QueryRow(`SELECT is_paid,status FROM salaries WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3`, fixture.salaryID, fixture.tenantID, fixture.restaurantID).Scan(&paid, &status); err != nil {
		t.Fatalf("read salary after rollback: %v", err)
	}
	if paid || status == "paid" {
		t.Fatalf("failed accounting post left salary paid: is_paid=%v status=%q", paid, status)
	}
	assertPayrollPostingCounts(t, fixture, 0, 0)
}

func TestCFR009PostgresFailureAtEachPaymentStageRollsBackAllEffects(t *testing.T) {
	for _, table := range []string{"salaries", "general_ledger", "hr_audit_logs"} {
		t.Run(table, func(t *testing.T) {
			fixture := openPayrollFixture(t)
			triggerName := fmt.Sprintf("cfr009_abort_%d", time.Now().UnixNano())
			functionName := fmt.Sprintf("cfr009_abort_fn_%d", time.Now().UnixNano())
			if _, err := fixture.db.Exec(fmt.Sprintf(`CREATE FUNCTION %s() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'injected %s failure'; END $$`, functionName, table)); err != nil {
				t.Fatalf("create failure-injection function: %v", err)
			}
			triggerOperation := "INSERT"
			if table == "salaries" {
				triggerOperation = "UPDATE"
			}
			if _, err := fixture.db.Exec(fmt.Sprintf(`CREATE TRIGGER %s BEFORE %s ON %s FOR EACH ROW EXECUTE FUNCTION %s()`, triggerName, triggerOperation, table, functionName)); err != nil {
				t.Fatalf("create failure-injection trigger: %v", err)
			}
			defer fixture.db.Exec(fmt.Sprintf(`DROP TRIGGER IF EXISTS %s ON %s`, triggerName, table))
			defer fixture.db.Exec(fmt.Sprintf(`DROP FUNCTION IF EXISTS %s()`, functionName))

			result := fixture.markPaid(t, "CFR009-INJECTED-"+table)
			if result.Code != stdhttp.StatusInternalServerError {
				t.Fatalf("injected %s failure status=%d body=%s, want 500", table, result.Code, result.Body.String())
			}
			var paid bool
			var status string
			if err := fixture.db.QueryRow(`SELECT is_paid,status FROM salaries WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3`, fixture.salaryID, fixture.tenantID, fixture.restaurantID).Scan(&paid, &status); err != nil {
				t.Fatalf("read salary after injected %s failure: %v", table, err)
			}
			if paid || status == "paid" {
				t.Fatalf("injected %s failure left salary paid: is_paid=%v status=%q", table, paid, status)
			}
			assertPayrollPostingCounts(t, fixture, 0, 0)
			var auditCount int
			if err := fixture.db.QueryRow(`SELECT COUNT(*) FROM hr_audit_logs WHERE tenant_id=$1 AND restaurant_id=$2 AND entity_type='salary' AND entity_id=$3`, fixture.tenantID, fixture.restaurantID, fixture.salaryID).Scan(&auditCount); err != nil || auditCount != 0 {
				t.Fatalf("injected %s failure audit count=%d err=%v, want 0", table, auditCount, err)
			}
		})
	}
}

func TestCFR009PostgresDisabledPostingRuleCannotCommitPaidState(t *testing.T) {
	fixture := openPayrollFixture(t)
	if _, err := fixture.db.Exec(`
		INSERT INTO auto_posting_rules (tenant_id,restaurant_id,rule_code,rule_name,is_active)
		VALUES ($1,$2,'SALARY_PAID','Disabled payroll posting',FALSE)
	`, fixture.tenantID, fixture.restaurantID); err != nil {
		t.Fatalf("disable payroll auto-posting rule: %v", err)
	}
	result := fixture.markPaid(t, "CFR009-DISABLED")
	if result.Code != stdhttp.StatusInternalServerError {
		t.Fatalf("disabled posting status=%d body=%s, want 500", result.Code, result.Body.String())
	}
	var paid bool
	var status string
	if err := fixture.db.QueryRow(`SELECT is_paid,status FROM salaries WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3`, fixture.salaryID, fixture.tenantID, fixture.restaurantID).Scan(&paid, &status); err != nil {
		t.Fatalf("read salary after disabled posting: %v", err)
	}
	if paid || status == "paid" {
		t.Fatalf("disabled posting left salary paid: is_paid=%v status=%q", paid, status)
	}
	assertPayrollPostingCounts(t, fixture, 0, 0)
}

func TestCFR009PostgresPaymentAndDeletionRespectTenantRestaurantScope(t *testing.T) {
	fixture := openPayrollFixture(t)
	otherTenant := insertCFR009Tenant(t, fixture.db, "tenant")
	otherRestaurant := insertCFR009Restaurant(t, fixture.db, otherTenant, "foreign")
	siblingRestaurant := insertCFR009Restaurant(t, fixture.db, fixture.tenantID, "sibling")
	repo := repository.NewSalaryRepository(fixture.db)
	for _, scope := range []struct {
		name         string
		tenantID     int
		restaurantID int
	}{
		{name: "foreign tenant", tenantID: otherTenant, restaurantID: otherRestaurant},
		{name: "sibling restaurant", tenantID: fixture.tenantID, restaurantID: siblingRestaurant},
	} {
		t.Run(scope.name, func(t *testing.T) {
			if err := repo.DeleteUnpaidSalary(scope.tenantID, scope.restaurantID, fixture.salaryID); err != repository.ErrSalaryNotFound {
				t.Fatalf("scoped delete error=%v, want not found", err)
			}
			forged, err := repo.GetSalaryByID(fixture.tenantID, fixture.restaurantID, fixture.salaryID)
			if err != nil || forged == nil {
				t.Fatalf("load salary for forged update: salary=%v err=%v", forged, err)
			}
			forged.TenantID = scope.tenantID
			forged.RestaurantID = scope.restaurantID
			if err := repo.UpdateSalary(forged); err != repository.ErrSalaryNotFound {
				t.Fatalf("scoped update error=%v, want not found", err)
			}
			tx, err := fixture.db.BeginTx(context.Background(), nil)
			if err != nil {
				t.Fatalf("begin scoped mark-paid: %v", err)
			}
			_, _, err = repo.MarkAsPaidTx(context.Background(), tx, scope.tenantID, scope.restaurantID, fixture.salaryID, 0, "forged")
			_ = tx.Rollback()
			if err != repository.ErrSalaryNotFound {
				t.Fatalf("scoped mark-paid error=%v, want not found", err)
			}
		})
	}
	var paid bool
	if err := fixture.db.QueryRow(`SELECT is_paid FROM salaries WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3`, fixture.salaryID, fixture.tenantID, fixture.restaurantID).Scan(&paid); err != nil || paid {
		t.Fatalf("foreign attempt changed salary paid=%v err=%v", paid, err)
	}
}

func insertCFR009Tenant(t *testing.T, db *sql.DB, suffix string) int {
	t.Helper()
	key := fmt.Sprintf("cfr009-%s-%d", suffix, time.Now().UnixNano())
	var id int
	if err := db.QueryRow(`INSERT INTO tenants (name,slug,domain) VALUES ($1,$2,$3) RETURNING id`, key, key, key+".example.test").Scan(&id); err != nil {
		t.Fatalf("insert isolated tenant: %v", err)
	}
	return id
}

func insertCFR009Restaurant(t *testing.T, db *sql.DB, tenantID int, suffix string) int {
	t.Helper()
	key := fmt.Sprintf("cfr009-%s-%d", suffix, time.Now().UnixNano())
	var id int
	if err := db.QueryRow(`INSERT INTO restaurants (tenant_id,name,slug,address) VALUES ($1,$2,$3,'synthetic') RETURNING id`, tenantID, key, key).Scan(&id); err != nil {
		t.Fatalf("insert isolated restaurant: %v", err)
	}
	return id
}

func assertPayrollPostingCounts(t *testing.T, fixture *payrollFixture, journals, ledger int) {
	t.Helper()
	var journalCount, ledgerCount int
	if err := fixture.db.QueryRow(`SELECT COUNT(*) FROM journal_entries WHERE tenant_id=$1 AND restaurant_id=$2 AND reference_type='salary' AND reference_id=$3 AND entry_type='auto_payroll' AND status='posted'`, fixture.tenantID, fixture.restaurantID, fixture.salaryID).Scan(&journalCount); err != nil {
		t.Fatalf("count payroll journals: %v", err)
	}
	if err := fixture.db.QueryRow(`SELECT COUNT(*) FROM general_ledger gl JOIN journal_entries je ON je.id=gl.journal_entry_id WHERE je.tenant_id=$1 AND je.restaurant_id=$2 AND je.reference_type='salary' AND je.reference_id=$3 AND je.entry_type='auto_payroll'`, fixture.tenantID, fixture.restaurantID, fixture.salaryID).Scan(&ledgerCount); err != nil {
		t.Fatalf("count payroll ledger lines: %v", err)
	}
	if journalCount != journals || ledgerCount != ledger {
		t.Fatalf("payroll effects journal=%d ledger=%d, want journal=%d ledger=%d", journalCount, ledgerCount, journals, ledger)
	}
}
