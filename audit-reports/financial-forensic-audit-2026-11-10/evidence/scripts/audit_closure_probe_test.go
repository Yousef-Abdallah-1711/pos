package usecase

import (
	"database/sql"
	"fmt"
	"encoding/json"
	"os"
	"testing"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
	_ "github.com/lib/pq"
)

func TestAuditClosureMaterialImbalancePosts(t *testing.T) {
	dsn := os.Getenv("AUDIT_DB_DSN")
	if dsn == "" { t.Fatal("AUDIT_DB_DSN is required") }
	db, err := sql.Open("postgres", dsn)
	if err != nil { t.Fatal(err) }
	defer db.Close()
	if err := db.Ping(); err != nil { t.Fatal(err) }

	var cashID, revenueID int
	err = db.QueryRow("SELECT id FROM chart_of_accounts WHERE tenant_id=2 AND restaurant_id=2 AND account_code='9911'").Scan(&cashID)
	if err != nil { t.Fatal(err) }
	err = db.QueryRow("SELECT id FROM chart_of_accounts WHERE tenant_id=2 AND restaurant_id=2 AND account_code='9912'").Scan(&revenueID)
	if err != nil { t.Fatal(err) }

	jr := repository.NewJournalEntryRepository(db)
	gl := repository.NewGeneralLedgerRepository(db, jr)
	engine := NewAutoPostingEngine(nil, nil, jr, nil, gl, nil)
	lines := []domain.JournalEntryLine{
		{AccountID: cashID, DebitAmount: 100, LineOrder: 1, Description: "Source debit"},
		{AccountID: revenueID, CreditAmount: 90, LineOrder: 2, Description: "Source credit"},
	}
	event := &domain.PostingEvent{
		TenantID: 2, RestaurantID: 2, EntryType: domain.AutoPostingEntryTypeSupplierInvoice,
		ReferenceType: "supplier_invoice", ReferenceID: 98111001,
		ReferenceNumber: "AUDIT-CLOSURE-FIN003-POSITIVE-100-90",
		Description: "Audit material imbalance probe: source debit 100 credit 90",
	}
	err = engine.PostSupplierInvoiceEntry(event, lines)
	if err == nil { t.Fatal("expected current automatic-posting call to fail with nil attachments") }
	var noAutoRows int
	if err := db.QueryRow("SELECT COUNT(*) FROM journal_entries WHERE tenant_id=2 AND restaurant_id=2 AND reference_number=$1", event.ReferenceNumber).Scan(&noAutoRows); err != nil { t.Fatal(err) }
	if noAutoRows != 0 { t.Fatalf("failed auto-post persisted %d rows", noAutoRows) }
	t.Logf("actual auto-post call failed before persistence: %v", err)

	normalized := engine.normalizeBalancedLines(lines)
	if normalized[1].CreditAmount != 100 || normalized[1].DebitAmount != 0 {
		t.Fatalf("unexpected normalizer result: line 2 debit=%v credit=%v", normalized[1].DebitAmount, normalized[1].CreditAmount)
	}
	t.Logf("normalizer changed source credit 90.00 to %.2f to absorb 10.00 difference", normalized[1].CreditAmount)

	refID := 98111003
	entry := &domain.JournalEntry{
		TenantID: 2, RestaurantID: 2, EntryDate: time.Date(2026,11,11,12,0,0,0,time.UTC),
		Description: "Audit normalized material imbalance probe", EntryType: "manual",
		ReferenceType: "supplier_invoice", ReferenceID: &refID,
		ReferenceNumber: fmt.Sprintf("AUDIT-CLOSURE-FIN003-NORM-%d", time.Now().UnixNano()),
		Attachments: json.RawMessage("[]"), Lines: normalized,
	}
	id, err := jr.CreateJournalEntry(entry)
	if err != nil { t.Fatalf("create normalized journal with valid JSON attachments: %v", err) }
	if err := gl.AppendEntry(2, 2, id, 0); err != nil { t.Fatalf("post normalized journal: %v", err) }
	var status string
	var debit, credit float64
	if err := db.QueryRow("SELECT status,total_debit,total_credit FROM journal_entries WHERE id=$1", id).Scan(&status,&debit,&credit); err != nil { t.Fatal(err) }
	var adjusted float64
	if err := db.QueryRow("SELECT credit_amount FROM journal_entry_lines WHERE journal_entry_id=$1 AND account_id=$2", id,revenueID).Scan(&adjusted); err != nil { t.Fatal(err) }
	var ledgerRows int
	if err := db.QueryRow("SELECT COUNT(*) FROM general_ledger WHERE journal_entry_id=$1", id).Scan(&ledgerRows); err != nil { t.Fatal(err) }
	if status!="posted" || debit!=100 || credit!=100 || adjusted!=100 || ledgerRows!=2 {
		t.Fatalf("unexpected posted result id=%d status=%s totals=%v/%v adjusted=%v gl=%d",id,status,debit,credit,adjusted,ledgerRows)
	}
	t.Logf("conditional persistence with valid attachments: id=%d status=%s totals=%.2f/%.2f revenue_credit=%.2f ledger_rows=%d",id,status,debit,credit,adjusted,ledgerRows)

	reverse := engine.normalizeBalancedLines([]domain.JournalEntryLine{
		{AccountID: cashID, DebitAmount: 100, LineOrder: 1},
		{AccountID: revenueID, CreditAmount: 110, LineOrder: 2},
	})
	reverseRef := 98111004
	reverseEntry := &domain.JournalEntry{
		TenantID:2, RestaurantID:2, EntryDate:time.Date(2026,11,11,12,0,0,0,time.UTC),
		Description:"Audit reverse material imbalance probe",EntryType:"manual",
		ReferenceType:"supplier_invoice",ReferenceID:&reverseRef,
		ReferenceNumber:fmt.Sprintf("AUDIT-CLOSURE-FIN003-REVERSE-%d",time.Now().UnixNano()),
		Attachments:json.RawMessage("[]"),Lines:reverse,
	}
	if _,err:=jr.CreateJournalEntry(reverseEntry);err==nil{t.Fatal("expected reverse imbalance to fail line validation")}
	t.Logf("reverse-direction normalized lines rejected before persistence: debit=%.2f credit=%.2f",reverse[1].DebitAmount,reverse[1].CreditAmount)


	tooLongRef := 98111005
	tooLong := &domain.JournalEntry{
		TenantID: 2, RestaurantID: 2, EntryDate: time.Date(2026,11,11,12,0,0,0,time.UTC),
		Description: "Audit auto entry type width probe", EntryType: domain.AutoPostingEntryTypeSupplierInvoice,
		ReferenceType: "supplier_invoice", ReferenceID: &tooLongRef,
		ReferenceNumber: "AUDIT-CLOSURE-FIN003-TYPE-WIDTH",
		Attachments: json.RawMessage("[]"), Lines: normalized,
	}
	if _, err := jr.CreateJournalEntry(tooLong); err == nil { t.Fatal("expected supplier auto entry type to exceed schema width") } else { t.Logf("supplier auto entry type rejected at varchar width: %v", err) }
}



