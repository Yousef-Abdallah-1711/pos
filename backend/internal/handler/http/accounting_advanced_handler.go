package http

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/database"
	"pos-saas/internal/repository"

	"github.com/go-pdf/fpdf"
)

type AccountingAdvancedHandler struct {
	db              *sql.DB
	repo            *repository.AccountingAdvancedRepository
	customerInvRepo *repository.CustomerInvoiceRepository
}

func NewAccountingAdvancedHandler(
	db *sql.DB,
	repo *repository.AccountingAdvancedRepository,
	customerInvRepo *repository.CustomerInvoiceRepository,
) *AccountingAdvancedHandler {
	return &AccountingAdvancedHandler{
		db:              db,
		repo:            repo,
		customerInvRepo: customerInvRepo,
	}
}

func (h *AccountingAdvancedHandler) ListBankReconciliationSessions(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	items, err := h.repo.ListBankReconciliationSessions(tenantID, restaurantID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to load reconciliation sessions",
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": items})
}

func (h *AccountingAdvancedHandler) CreateBankReconciliationSession(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	var req struct {
		BankAccountID    int     `json:"bank_account_id"`
		StatementDate    string  `json:"statement_date"`
		StatementBalance float64 `json:"statement_balance"`
		BookBalance      float64 `json:"book_balance"`
		Notes            string  `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Invalid request body"})
		return
	}
	if req.BankAccountID <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "bank_account_id is required"})
		return
	}
	statementDate, err := parseAccountingInputDate(req.StatementDate, "statement_date")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}

	item := &domain.BankReconciliationSession{
		TenantID:         tenantID,
		RestaurantID:     restaurantID,
		BankAccountID:    req.BankAccountID,
		StatementDate:    statementDate,
		StatementBalance: req.StatementBalance,
		BookBalance:      req.BookBalance,
		Difference:       req.StatementBalance - req.BookBalance,
		Status:           "draft",
		Notes:            strings.TrimSpace(req.Notes),
		CreatedBy:        &userID,
		UpdatedBy:        &userID,
	}
	id, err := h.repo.CreateBankReconciliationSession(item)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	created, err := h.repo.GetBankReconciliationSessionByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": map[string]int{"id": id}})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": created})
}

func (h *AccountingAdvancedHandler) AddBankReconciliationLine(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	reconciliationID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Invalid reconciliation ID"})
		return
	}

	var req struct {
		JournalEntryID  *int    `json:"journal_entry_id"`
		LineDescription string  `json:"line_description"`
		StatementAmount float64 `json:"statement_amount"`
		BookAmount      float64 `json:"book_amount"`
		MatchStatus     string  `json:"match_status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Invalid request body"})
		return
	}

	line := &domain.BankReconciliationLine{
		JournalEntryID:  req.JournalEntryID,
		LineDescription: strings.TrimSpace(req.LineDescription),
		StatementAmount: req.StatementAmount,
		BookAmount:      req.BookAmount,
		MatchStatus:     defaultAccountingString(req.MatchStatus, "matched"),
	}
	if err := h.repo.AddBankReconciliationLine(tenantID, restaurantID, reconciliationID, line); err != nil {
		writeAccountingMutationError(w, err)
		return
	}
	session, err := h.repo.GetBankReconciliationSessionByID(tenantID, restaurantID, reconciliationID)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{"success": true})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": session})
}

func (h *AccountingAdvancedHandler) CompleteBankReconciliationSession(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	reconciliationID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Invalid reconciliation ID"})
		return
	}
	if err := h.repo.CompleteBankReconciliationSession(tenantID, restaurantID, reconciliationID, userID); err != nil {
		writeAccountingMutationError(w, err)
		return
	}
	session, err := h.repo.GetBankReconciliationSessionByID(tenantID, restaurantID, reconciliationID)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{"success": true})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": session})
}

func (h *AccountingAdvancedHandler) ListRecurringJournalTemplates(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	items, err := h.repo.ListRecurringJournalTemplates(tenantID, restaurantID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to load recurring templates"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": items})
}

func (h *AccountingAdvancedHandler) CreateRecurringJournalTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	var req struct {
		TemplateCode string                         `json:"template_code"`
		TemplateName string                         `json:"template_name"`
		Description  string                         `json:"description"`
		Frequency    string                         `json:"frequency"`
		StartDate    string                         `json:"start_date"`
		NextRunDate  string                         `json:"next_run_date"`
		AutoPost     bool                           `json:"auto_post"`
		IsActive     *bool                          `json:"is_active"`
		Lines        []domain.JournalEntryLineInput `json:"lines"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Invalid request body"})
		return
	}
	if strings.TrimSpace(req.TemplateCode) == "" || strings.TrimSpace(req.TemplateName) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "template_code and template_name are required"})
		return
	}
	if len(req.Lines) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "at least one journal line is required"})
		return
	}
	startDate, err := parseAccountingInputDate(req.StartDate, "start_date")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	nextRunDate := startDate
	if strings.TrimSpace(req.NextRunDate) != "" {
		nextRunDate, err = parseAccountingInputDate(req.NextRunDate, "next_run_date")
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
			return
		}
	}

	linesJSON, _ := json.Marshal(req.Lines)
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	item := &domain.RecurringJournalTemplate{
		TenantID:     tenantID,
		RestaurantID: restaurantID,
		TemplateCode: strings.TrimSpace(req.TemplateCode),
		TemplateName: strings.TrimSpace(req.TemplateName),
		Description:  strings.TrimSpace(req.Description),
		Frequency:    defaultAccountingString(strings.ToLower(req.Frequency), "monthly"),
		StartDate:    startDate,
		NextRunDate:  nextRunDate,
		IsActive:     isActive,
		AutoPost:     req.AutoPost,
		LinesJSON:    linesJSON,
		CreatedBy:    &userID,
		UpdatedBy:    &userID,
	}
	id, err := h.repo.CreateRecurringJournalTemplate(item)
	if err != nil {
		writeAccountingMutationError(w, err)
		return
	}
	created, err := h.repo.GetRecurringJournalTemplateByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": map[string]int{"id": id}})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": created})
}

func (h *AccountingAdvancedHandler) UpdateRecurringJournalTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Invalid template ID"})
		return
	}

	existing, err := h.repo.GetRecurringJournalTemplateByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to load template"})
		return
	}
	if existing == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{"success": false, "error": "Recurring template not found"})
		return
	}

	var req struct {
		TemplateCode string                         `json:"template_code"`
		TemplateName string                         `json:"template_name"`
		Description  string                         `json:"description"`
		Frequency    string                         `json:"frequency"`
		StartDate    string                         `json:"start_date"`
		NextRunDate  string                         `json:"next_run_date"`
		AutoPost     bool                           `json:"auto_post"`
		IsActive     *bool                          `json:"is_active"`
		Lines        []domain.JournalEntryLineInput `json:"lines"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Invalid request body"})
		return
	}

	startDate, err := parseAccountingInputDate(req.StartDate, "start_date")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	nextRunDate := existing.NextRunDate
	if strings.TrimSpace(req.NextRunDate) != "" {
		nextRunDate, err = parseAccountingInputDate(req.NextRunDate, "next_run_date")
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
			return
		}
	}
	linesJSON, _ := json.Marshal(req.Lines)
	isActive := existing.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	item := &domain.RecurringJournalTemplate{
		ID:           id,
		TenantID:     tenantID,
		RestaurantID: restaurantID,
		TemplateCode: strings.TrimSpace(req.TemplateCode),
		TemplateName: strings.TrimSpace(req.TemplateName),
		Description:  strings.TrimSpace(req.Description),
		Frequency:    defaultAccountingString(strings.ToLower(req.Frequency), "monthly"),
		StartDate:    startDate,
		NextRunDate:  nextRunDate,
		IsActive:     isActive,
		AutoPost:     req.AutoPost,
		LinesJSON:    linesJSON,
		UpdatedBy:    &userID,
	}
	if item.TemplateCode == "" {
		item.TemplateCode = existing.TemplateCode
	}
	if item.TemplateName == "" {
		item.TemplateName = existing.TemplateName
	}
	if len(item.LinesJSON) == 0 {
		item.LinesJSON = existing.LinesJSON
	}

	if err := h.repo.UpdateRecurringJournalTemplate(item); err != nil {
		writeAccountingMutationError(w, err)
		return
	}
	updated, err := h.repo.GetRecurringJournalTemplateByID(tenantID, restaurantID, id)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{"success": true})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": updated})
}

func (h *AccountingAdvancedHandler) DeleteRecurringJournalTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Invalid template ID"})
		return
	}
	if err := h.repo.DeleteRecurringJournalTemplate(tenantID, restaurantID, id); err != nil {
		writeAccountingMutationError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true})
}

func (h *AccountingAdvancedHandler) RunRecurringJournalTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Invalid template ID"})
		return
	}

	var req struct {
		RunDate string `json:"run_date"`
	}
	_ = decodeAccountingActionBody(r, &req)

	runDate := time.Now()
	if strings.TrimSpace(req.RunDate) != "" {
		runDate, err = parseAccountingInputDate(req.RunDate, "run_date")
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
			return
		}
	}
	entryID, err := h.repo.RunRecurringJournalTemplate(tenantID, restaurantID, id, userID, runDate)
	if err != nil {
		writeAccountingMutationError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"journal_entry_id": entryID,
			"run_date":         runDate.Format("2006-01-02"),
		},
	})
}

func (h *AccountingAdvancedHandler) ListBudgets(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	fiscalYear, _ := parseIntQueryParam(r, "fiscal_year", 0)
	items, err := h.repo.ListBudgetVariance(tenantID, restaurantID, fiscalYear)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to load budgets"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": items})
}

func (h *AccountingAdvancedHandler) UpsertBudget(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))

	var req struct {
		AccountID    int     `json:"account_id"`
		FiscalYear   int     `json:"fiscal_year"`
		PeriodMonth  int     `json:"period_month"`
		BudgetAmount float64 `json:"budget_amount"`
		Currency     string  `json:"currency"`
		Notes        string  `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Invalid request body"})
		return
	}
	if req.AccountID <= 0 || req.FiscalYear <= 0 || req.PeriodMonth <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "account_id, fiscal_year, and period_month are required"})
		return
	}
	id, err := h.repo.UpsertBudget(
		tenantID,
		restaurantID,
		req.AccountID,
		req.FiscalYear,
		req.PeriodMonth,
		req.BudgetAmount,
		req.Currency,
		req.Notes,
	)
	if err != nil {
		writeAccountingMutationError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": map[string]int{"id": id}})
}

func (h *AccountingAdvancedHandler) ListExchangeRates(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	baseCurrency := strings.TrimSpace(r.URL.Query().Get("base_currency"))
	quoteCurrency := strings.TrimSpace(r.URL.Query().Get("quote_currency"))
	items, err := h.repo.ListExchangeRates(tenantID, restaurantID, baseCurrency, quoteCurrency)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to load exchange rates"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": items})
}

func (h *AccountingAdvancedHandler) UpsertExchangeRate(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	var req struct {
		BaseCurrency  string  `json:"base_currency"`
		QuoteCurrency string  `json:"quote_currency"`
		Rate          float64 `json:"rate"`
		EffectiveDate string  `json:"effective_date"`
		Source        string  `json:"source"`
		IsActive      *bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Invalid request body"})
		return
	}
	if strings.TrimSpace(req.BaseCurrency) == "" || strings.TrimSpace(req.QuoteCurrency) == "" || req.Rate <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "base_currency, quote_currency and positive rate are required"})
		return
	}
	effectiveDate, err := parseAccountingInputDate(req.EffectiveDate, "effective_date")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	id, err := h.repo.UpsertExchangeRate(&domain.ExchangeRate{
		TenantID:      tenantID,
		RestaurantID:  restaurantID,
		BaseCurrency:  req.BaseCurrency,
		QuoteCurrency: req.QuoteCurrency,
		Rate:          req.Rate,
		EffectiveDate: effectiveDate,
		Source:        req.Source,
		IsActive:      isActive,
	})
	if err != nil {
		writeAccountingMutationError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": map[string]int{"id": id}})
}

func (h *AccountingAdvancedHandler) ConvertCurrency(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	base := strings.TrimSpace(r.URL.Query().Get("base_currency"))
	quote := strings.TrimSpace(r.URL.Query().Get("quote_currency"))
	amount, _ := strconv.ParseFloat(strings.TrimSpace(r.URL.Query().Get("amount")), 64)
	asOf := time.Now()
	if rawDate := strings.TrimSpace(r.URL.Query().Get("as_of")); rawDate != "" {
		var err error
		asOf, err = parseAccountingInputDate(rawDate, "as_of")
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
			return
		}
	}
	converted, rate, err := h.repo.ConvertAmount(tenantID, restaurantID, base, quote, amount, asOf)
	if err != nil {
		writeAccountingMutationError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"base_currency":  strings.ToUpper(base),
			"quote_currency": strings.ToUpper(quote),
			"amount":         amount,
			"rate":           rate,
			"converted":      converted,
			"as_of":          asOf.Format("2006-01-02"),
		},
	})
}

func (h *AccountingAdvancedHandler) ListSupplierPaymentSchedules(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	items, err := h.repo.ListSupplierPaymentSchedules(tenantID, restaurantID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to load supplier payment schedules"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": items})
}

func (h *AccountingAdvancedHandler) CreateSupplierPaymentSchedule(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))
	var req struct {
		SupplierInvoiceID *int    `json:"supplier_invoice_id"`
		SupplierID        int     `json:"supplier_id"`
		ScheduledDate     string  `json:"scheduled_date"`
		Amount            float64 `json:"amount"`
		Currency          string  `json:"currency"`
		ReminderDays      int     `json:"reminder_days"`
		Notes             string  `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Invalid request body"})
		return
	}
	if req.SupplierID <= 0 || req.Amount <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "supplier_id and positive amount are required"})
		return
	}
	scheduledDate, err := parseAccountingInputDate(req.ScheduledDate, "scheduled_date")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	item := &domain.SupplierPaymentSchedule{
		TenantID:          tenantID,
		RestaurantID:      restaurantID,
		SupplierInvoiceID: req.SupplierInvoiceID,
		SupplierID:        req.SupplierID,
		ScheduledDate:     scheduledDate,
		Amount:            req.Amount,
		Currency:          req.Currency,
		ReminderDays:      req.ReminderDays,
		Status:            "scheduled",
		Notes:             req.Notes,
		CreatedBy:         &userID,
		UpdatedBy:         &userID,
	}
	id, err := h.repo.CreateSupplierPaymentSchedule(item)
	if err != nil {
		writeAccountingMutationError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": map[string]int{"id": id}})
}

func (h *AccountingAdvancedHandler) ListSupplierPaymentReminders(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	dueWithinDays, _ := parseIntQueryParam(r, "due_days", 7)
	items, err := h.repo.ListSupplierPaymentReminders(tenantID, restaurantID, dueWithinDays)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to load reminders"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": items})
}

func (h *AccountingAdvancedHandler) TriggerSupplierPaymentReminder(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	scheduleID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Invalid schedule ID"})
		return
	}
	if err := h.repo.MarkSupplierPaymentReminderSent(tenantID, restaurantID, scheduleID); err != nil {
		writeAccountingMutationError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true})
}

func (h *AccountingAdvancedHandler) DownloadCustomerInvoicePDF(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	invoiceID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Invalid invoice ID"})
		return
	}
	invoice, err := h.customerInvRepo.GetInvoiceByID(tenantID, restaurantID, invoiceID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to load invoice"})
		return
	}
	if invoice == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{"success": false, "error": "Invoice not found"})
		return
	}

	pdfBytes, err := buildCustomerInvoicePDF(invoice)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to generate invoice PDF"})
		return
	}

	filename := fmt.Sprintf("invoice-%s.pdf", strings.TrimSpace(invoice.InvoiceNumber))
	if filename == "invoice-.pdf" {
		filename = fmt.Sprintf("invoice-%d.pdf", invoice.ID)
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pdfBytes)
}

func (h *AccountingAdvancedHandler) SendCustomerInvoiceEmail(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))
	invoiceID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Invalid invoice ID"})
		return
	}
	var req struct {
		ToEmail    string `json:"to_email"`
		Subject    string `json:"subject"`
		Message    string `json:"message"`
		IncludePDF *bool  `json:"include_pdf"`
	}
	_ = decodeAccountingActionBody(r, &req)
	includePDF := true
	if req.IncludePDF != nil {
		includePDF = *req.IncludePDF
	}
	invoice, err := h.customerInvRepo.GetInvoiceByID(tenantID, restaurantID, invoiceID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to load invoice"})
		return
	}
	if invoice == nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{"success": false, "error": "Invoice not found"})
		return
	}
	pdfSize := 0
	if includePDF {
		pdfBytes, err := buildCustomerInvoicePDF(invoice)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to generate invoice PDF"})
			return
		}
		pdfSize = len(pdfBytes)
	}

	recipient := strings.TrimSpace(req.ToEmail)
	if recipient == "" {
		recipient = h.lookupCustomerEmail(tenantID, restaurantID, invoice.CustomerID)
	}
	if recipient == "" {
		recipient = fmt.Sprintf("customer-%d@pending.local", invoice.CustomerID)
	}

	subject := strings.TrimSpace(req.Subject)
	if subject == "" {
		subject = fmt.Sprintf("Invoice %s", invoice.InvoiceNumber)
	}
	body := strings.TrimSpace(req.Message)
	if body == "" {
		body = "Please find your invoice attached."
	}
	messageID := fmt.Sprintf("inv-%d-%d", invoice.ID, time.Now().Unix())

	payload, _ := json.Marshal(map[string]interface{}{
		"subject":     subject,
		"message":     body,
		"include_pdf": includePDF,
		"pdf_size":    pdfSize,
		"message_id":  messageID,
	})
	err = database.WithTenantTx(r.Context(), h.db, int64(tenantID), int64(restaurantID), func(tx *sql.Tx) error {
		_, err := tx.ExecContext(r.Context(), `
			INSERT INTO invoice_delivery_logs (
				tenant_id, restaurant_id, invoice_id, channel, recipient_email, status, payload, created_by
			) VALUES (
				$1, $2, $3, 'email', $4, 'sent', $5::jsonb, $6
			)
		`, tenantID, restaurantID, invoiceID, recipient, string(payload), userID)
		return err
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to log delivery"})
		return
	}

	log.Printf("[INVOICE_EMAIL] tenant=%d restaurant=%d invoice=%d recipient=%s message_id=%s", tenantID, restaurantID, invoiceID, recipient, messageID)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"message_id":   messageID,
			"to_email":     recipient,
			"subject":      subject,
			"pdf_attached": includePDF,
			"pdf_size":     pdfSize,
		},
	})
}

func (h *AccountingAdvancedHandler) GetFinancialKPIs(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	dateFrom, dateTo, err := parseDateRangeForReports(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	data, err := h.repo.GetFinancialKPIs(tenantID, restaurantID, dateFrom, dateTo)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to generate KPI dashboard"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": data})
}

func (h *AccountingAdvancedHandler) GetMenuProfitability(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	dateFrom, dateTo, err := parseDateRangeForReports(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	limit, _ := parseIntQueryParam(r, "limit", 20)
	items, err := h.repo.GetMenuItemProfitability(tenantID, restaurantID, dateFrom, dateTo, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to generate menu profitability"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": map[string]interface{}{"items": items, "generated_at": time.Now().UTC()}})
}

func (h *AccountingAdvancedHandler) GetBreakEvenAnalysis(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	dateFrom, dateTo, err := parseDateRangeForReports(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	data, err := h.repo.GetBreakEvenAnalysis(tenantID, restaurantID, dateFrom, dateTo)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to generate break-even analysis"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": data})
}

func (h *AccountingAdvancedHandler) GetPredictiveCashFlow(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	historyMonths, _ := parseIntQueryParam(r, "history_months", 6)
	forecastMonths, _ := parseIntQueryParam(r, "forecast_months", 3)
	data, err := h.repo.GetPredictiveCashFlow(tenantID, restaurantID, historyMonths, forecastMonths)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to generate predictive cash flow"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": data})
}

func (h *AccountingAdvancedHandler) GetHRCostAnalysis(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	dateFrom, dateTo, err := parseDateRangeForReports(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	data, err := h.repo.GetHRCostAnalysis(tenantID, restaurantID, dateFrom, dateTo)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to generate HR cost analysis"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": data})
}

func (h *AccountingAdvancedHandler) SubmitEInvoice(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))

	invoiceID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "Invalid invoice ID"})
		return
	}

	var req struct {
		InvoiceType string `json:"invoice_type"`
	}
	_ = decodeAccountingActionBody(r, &req)
	invoiceType := strings.ToLower(strings.TrimSpace(req.InvoiceType))
	if invoiceType == "" {
		invoiceType = "customer"
	}
	if invoiceType != "customer" && invoiceType != "supplier" {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "invoice_type must be customer or supplier"})
		return
	}

	mode := strings.ToLower(strings.TrimSpace(readEnvDefault("ETA_EINVOICE_MODE", "mock")))
	status := "accepted"
	if mode != "mock" {
		status = "pending"
	}
	etaUUID := fmt.Sprintf("ETA-%d-%d", invoiceID, time.Now().Unix())
	responsePayload, _ := json.Marshal(map[string]interface{}{
		"mode":       mode,
		"invoice_id": invoiceID,
		"status":     status,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	})

	submission := &domain.EInvoiceSubmission{
		TenantID:        tenantID,
		RestaurantID:    restaurantID,
		InvoiceType:     invoiceType,
		InvoiceID:       invoiceID,
		ETAUUID:         etaUUID,
		Status:          status,
		ResponsePayload: responsePayload,
		CreatedBy:       &userID,
	}
	id, err := h.repo.CreateEInvoiceSubmission(r.Context(), submission)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to submit e-invoice"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":         id,
			"eta_uuid":   etaUUID,
			"mode":       mode,
			"status":     status,
			"invoice_id": invoiceID,
		},
	})
}

func (h *AccountingAdvancedHandler) ListEInvoiceSubmissions(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	items, err := h.repo.ListEInvoiceSubmissions(r.Context(), tenantID, restaurantID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to load e-invoice submissions"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": items})
}

func (h *AccountingAdvancedHandler) ExportSocialInsuranceCSV(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	now := time.Now()
	year, _ := parseIntQueryParam(r, "year", now.Year())
	month, _ := parseIntQueryParam(r, "month", int(now.Month()))

	rows, err := h.repo.GetSocialInsuranceReportRows(tenantID, restaurantID, year, month)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to generate social insurance report"})
		return
	}

	csvRows := [][]string{{
		"employee_code",
		"employee_name",
		"national_id",
		"base_salary",
		"social_insurance_deduction",
		"health_insurance_deduction",
		"policy_number",
		"provider_name",
	}}
	for _, row := range rows {
		csvRows = append(csvRows, []string{
			fmt.Sprintf("%v", row["employee_code"]),
			fmt.Sprintf("%v", row["employee_name"]),
			fmt.Sprintf("%v", row["national_id"]),
			fmt.Sprintf("%.2f", toFloat64(row["base_salary"])),
			fmt.Sprintf("%.2f", toFloat64(row["social_insurance_deduction"])),
			fmt.Sprintf("%.2f", toFloat64(row["health_insurance_deduction"])),
			fmt.Sprintf("%v", row["policy_number"]),
			fmt.Sprintf("%v", row["provider_name"]),
		})
	}
	writeCSVAttachment(w, fmt.Sprintf("social-insurance-%04d-%02d.csv", year, month), csvRows)
}

func (h *AccountingAdvancedHandler) ExportAuditTrailCSV(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))
	exportType := strings.TrimSpace(r.URL.Query().Get("type"))
	dateFrom, dateTo, err := parseDateRangeForReports(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	rows, err := h.repo.GetAuditExportRows(r.Context(), tenantID, restaurantID, dateFrom, dateTo, exportType)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to export audit trail"})
		return
	}

	csvRows := [][]string{{"entity_type", "timestamp", "reference", "description", "amount", "actor"}}
	for _, row := range rows {
		csvRows = append(csvRows, []string{
			fmt.Sprintf("%v", row["entity_type"]),
			fmt.Sprintf("%v", row["timestamp"]),
			fmt.Sprintf("%v", row["reference"]),
			fmt.Sprintf("%v", row["description"]),
			fmt.Sprintf("%.2f", toFloat64(row["amount"])),
			fmt.Sprintf("%v", row["actor"]),
		})
	}
	_ = h.repo.LogAuditExport(r.Context(), tenantID, restaurantID, userID, defaultAccountingString(exportType, "all"), "csv", map[string]interface{}{
		"date_from": dateFrom.Format("2006-01-02"),
		"date_to":   dateTo.Format("2006-01-02"),
		"type":      defaultAccountingString(exportType, "all"),
	}, len(rows))
	writeCSVAttachment(w, fmt.Sprintf("audit-trail-%s.csv", time.Now().UTC().Format("20060102-150405")), csvRows)
}

func (h *AccountingAdvancedHandler) ExportChartOfAccountsExcel(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	rows, err := h.repo.GetChartOfAccountsExcelRows(tenantID, restaurantID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to export chart of accounts"})
		return
	}
	writeCSVAttachment(w, "chart-of-accounts-export.csv", rows)
}

func (h *AccountingAdvancedHandler) ExportJournalEntriesExcel(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	dateFrom, dateTo, err := parseDateRangeForReports(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	rows, err := h.repo.GetJournalEntriesExcelRows(tenantID, restaurantID, dateFrom, dateTo)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to export journal entries"})
		return
	}
	writeCSVAttachment(w, "journal-entries-export.csv", rows)
}

func (h *AccountingAdvancedHandler) ImportChartOfAccountsExcel(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))
	_ = userID

	rows, err := parseUploadedCSVRows(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	if len(rows) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "CSV file is empty"})
		return
	}
	headers := normalizeCSVHeaders(rows[0])
	if len(headers) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "CSV header is required"})
		return
	}

	createdOrUpdated := 0
	errorsCount := 0
	for idx := 1; idx < len(rows); idx++ {
		record := mapCSVRecord(headers, rows[idx])
		isActive := strings.ToLower(record["is_active"]) != "false"
		err := h.repo.UpsertAccountFromImport(
			tenantID,
			restaurantID,
			record["account_code"],
			record["account_name"],
			record["account_name_ar"],
			record["account_type"],
			record["normal_balance"],
			record["currency"],
			record["account_subtype"],
			record["description"],
			isActive,
		)
		if err != nil {
			errorsCount++
			continue
		}
		createdOrUpdated++
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"processed_rows": len(rows) - 1,
			"upserted_rows":  createdOrUpdated,
			"failed_rows":    errorsCount,
		},
	})
}

func (h *AccountingAdvancedHandler) ImportJournalEntriesExcel(w http.ResponseWriter, r *http.Request) {
	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	userID := int(middleware.GetUserID(r))
	autoPost := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("auto_post"))) == "true"

	rows, err := parseUploadedCSVRows(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	if len(rows) <= 1 {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": "No journal entry rows found"})
		return
	}
	headers := normalizeCSVHeaders(rows[0])
	grouped := make(map[string][]map[string]string)
	for idx := 1; idx < len(rows); idx++ {
		record := mapCSVRecord(headers, rows[idx])
		entryKey := strings.TrimSpace(record["entry_key"])
		if entryKey == "" {
			entryKey = strings.TrimSpace(record["entry_number"])
		}
		if entryKey == "" {
			entryKey = fmt.Sprintf("row-%d", idx)
		}
		grouped[entryKey] = append(grouped[entryKey], record)
	}

	createdEntries := 0
	failedEntries := 0
	for _, group := range grouped {
		if len(group) == 0 {
			continue
		}
		entryDate, err := parseAccountingInputDate(group[0]["entry_date"], "entry_date")
		if err != nil {
			failedEntries++
			continue
		}
		description := group[0]["description"]
		entryType := defaultAccountingString(group[0]["entry_type"], "manual")

		lines := make([]domain.JournalEntryLineInput, 0, len(group))
		for _, row := range group {
			accountCode := strings.TrimSpace(row["account_code"])
			accountID, err := h.repo.GetAccountIDByCode(tenantID, restaurantID, accountCode)
			if err != nil {
				lines = nil
				break
			}
			debit, _ := strconv.ParseFloat(strings.TrimSpace(row["debit_amount"]), 64)
			credit, _ := strconv.ParseFloat(strings.TrimSpace(row["credit_amount"]), 64)
			lines = append(lines, domain.JournalEntryLineInput{
				AccountID:    accountID,
				DebitAmount:  debit,
				CreditAmount: credit,
				Description:  row["line_description"],
			})
		}
		if len(lines) == 0 {
			failedEntries++
			continue
		}
		_, err = h.repo.CreateImportedJournalEntry(
			tenantID,
			restaurantID,
			userID,
			entryDate,
			description,
			entryType,
			lines,
			autoPost,
		)
		if err != nil {
			failedEntries++
			continue
		}
		createdEntries++
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"processed_entries": len(grouped),
			"created_entries":   createdEntries,
			"failed_entries":    failedEntries,
		},
	})
}

func parseIntQueryParam(r *http.Request, key string, fallback int) (int, error) {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback, err
	}
	return value, nil
}

func parseDateRangeForReports(r *http.Request) (time.Time, time.Time, error) {
	today := time.Now().UTC()
	defaultFrom := today.AddDate(0, 0, -30)
	defaultTo := today.AddDate(0, 0, 1)

	rawFrom := strings.TrimSpace(r.URL.Query().Get("date_from"))
	rawTo := strings.TrimSpace(r.URL.Query().Get("date_to"))

	dateFrom := defaultFrom
	dateTo := defaultTo
	var err error
	if rawFrom != "" {
		dateFrom, err = parseAccountingInputDate(rawFrom, "date_from")
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}
	if rawTo != "" {
		dateTo, err = parseAccountingInputDate(rawTo, "date_to")
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		dateTo = dateTo.AddDate(0, 0, 1)
	}
	if !dateTo.After(dateFrom) {
		return time.Time{}, time.Time{}, errors.New("date_to must be after date_from")
	}
	return dateFrom, dateTo, nil
}

func writeCSVAttachment(w http.ResponseWriter, filename string, rows [][]string) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	for _, row := range rows {
		_ = writer.Write(row)
	}
	writer.Flush()

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

func parseUploadedCSVRows(r *http.Request) ([][]string, error) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return nil, errors.New("multipart form with file is required")
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, errors.New("file field is required")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("invalid CSV file: %w", err)
	}
	return rows, nil
}

func normalizeCSVHeaders(headers []string) []string {
	out := make([]string, len(headers))
	for i, value := range headers {
		out[i] = strings.ToLower(strings.TrimSpace(value))
	}
	return out
}

func mapCSVRecord(headers, row []string) map[string]string {
	record := make(map[string]string, len(headers))
	for i, header := range headers {
		if i < len(row) {
			record[header] = strings.TrimSpace(row[i])
		} else {
			record[header] = ""
		}
	}
	return record
}

func buildCustomerInvoicePDF(invoice *domain.CustomerInvoice) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(14, 14, 14)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "Customer Invoice", "", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 7, fmt.Sprintf("Invoice #: %s", invoice.InvoiceNumber), "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 7, fmt.Sprintf("Invoice Date: %s", invoice.InvoiceDate.Format("2006-01-02")), "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 7, fmt.Sprintf("Due Date: %s", invoice.DueDate.Format("2006-01-02")), "", 1, "L", false, 0, "")
	if strings.TrimSpace(invoice.CustomerName) != "" {
		pdf.CellFormat(0, 7, fmt.Sprintf("Customer: %s", invoice.CustomerName), "", 1, "L", false, 0, "")
	}
	pdf.Ln(3)

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(95, 8, "Description", "1", 0, "L", false, 0, "")
	pdf.CellFormat(20, 8, "Qty", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 8, "Unit Price", "1", 0, "R", false, 0, "")
	pdf.CellFormat(30, 8, "Line Total", "1", 1, "R", false, 0, "")
	pdf.SetFont("Arial", "", 10)

	for _, line := range invoice.Lines {
		pdf.CellFormat(95, 8, truncateForPDF(line.Description, 52), "1", 0, "L", false, 0, "")
		pdf.CellFormat(20, 8, fmt.Sprintf("%.2f", line.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("%.2f", line.UnitPrice), "1", 0, "R", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("%.2f", line.LineTotal), "1", 1, "R", false, 0, "")
	}

	pdf.Ln(3)
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(145, 7, "Subtotal", "", 0, "R", false, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", invoice.Subtotal), "", 1, "R", false, 0, "")
	pdf.CellFormat(145, 7, "Tax", "", 0, "R", false, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", invoice.TaxAmount), "", 1, "R", false, 0, "")
	pdf.CellFormat(145, 7, "Discount", "", 0, "R", false, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", invoice.DiscountAmount), "", 1, "R", false, 0, "")
	pdf.CellFormat(145, 8, "Total", "T", 0, "R", false, 0, "")
	pdf.CellFormat(30, 8, fmt.Sprintf("%.2f %s", invoice.TotalAmount, defaultAccountingString(invoice.Currency, "EGP")), "T", 1, "R", false, 0, "")

	var out bytes.Buffer
	if err := pdf.Output(&out); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func truncateForPDF(value string, max int) string {
	value = strings.TrimSpace(value)
	if len(value) <= max {
		return value
	}
	if max <= 3 {
		return value[:max]
	}
	return value[:max-3] + "..."
}

func (h *AccountingAdvancedHandler) lookupCustomerEmail(tenantID, restaurantID, customerID int) string {
	if customerID <= 0 {
		return ""
	}
	var email sql.NullString
	err := h.db.QueryRow(`
		SELECT email
		FROM customers_accounting
		WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3
	`, tenantID, restaurantID, customerID).Scan(&email)
	if err != nil || !email.Valid {
		return ""
	}
	return strings.TrimSpace(email.String)
}

func toFloat64(value interface{}) float64 {
	switch typed := value.(type) {
	case float64:
		return typed
	case float32:
		return float64(typed)
	case int:
		return float64(typed)
	case int64:
		return float64(typed)
	case json.Number:
		out, _ := typed.Float64()
		return out
	case string:
		out, _ := strconv.ParseFloat(strings.TrimSpace(typed), 64)
		return out
	default:
		return 0
	}
}

func readEnvDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
