package usecase

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
)

type AutoPostingService interface {
	PostSalesEntry(order *domain.Order, postedBy *int) error
	PostCOGSEntry(order *domain.Order, postedBy *int) error
	PostPayrollEntry(salary *domain.Salary, postedBy *int) error
}

type AutoPostingEngine struct {
	autoPostingRepo       *repository.AutoPostingRepository
	accountRepo           *repository.AccountRepository
	journalRepo           *repository.JournalEntryRepository
	inventoryMovementRepo *repository.InventoryMovementRepository
	generalLedgerRepo     *repository.GeneralLedgerRepository
	productRepo           *repository.ProductRepository
}

func NewAutoPostingEngine(
	autoPostingRepo *repository.AutoPostingRepository,
	accountRepo *repository.AccountRepository,
	journalRepo *repository.JournalEntryRepository,
	inventoryMovementRepo *repository.InventoryMovementRepository,
	generalLedgerRepo *repository.GeneralLedgerRepository,
	productRepo *repository.ProductRepository,
) *AutoPostingEngine {
	return &AutoPostingEngine{
		autoPostingRepo:       autoPostingRepo,
		accountRepo:           accountRepo,
		journalRepo:           journalRepo,
		inventoryMovementRepo: inventoryMovementRepo,
		generalLedgerRepo:     generalLedgerRepo,
		productRepo:           productRepo,
	}
}

func (e *AutoPostingEngine) PostSalesEntry(order *domain.Order, postedBy *int) error {
	if order == nil {
		return errors.New("order is required")
	}
	if !strings.EqualFold(strings.TrimSpace(order.PaymentStatus), "paid") {
		return nil
	}

	accounts, enabled, err := e.ruleAccounts(int(order.TenantID), int(order.RestaurantID), domain.AutoPostingEventOrderPaid)
	if err != nil || !enabled {
		return err
	}

	existingID, err := e.autoPostingRepo.FindJournalEntryByReference(
		int(order.TenantID),
		int(order.RestaurantID),
		domain.AutoPostingReferenceTypeOrder,
		domain.AutoPostingEntryTypeSales,
		int(order.ID),
	)
	if err != nil {
		return err
	}
	if existingID != nil {
		return nil
	}

	collectionAccount, err := e.resolvePaymentAccount(int(order.TenantID), int(order.RestaurantID), strings.ToLower(strings.TrimSpace(order.PaymentMethod)), accounts)
	if err != nil {
		return err
	}
	salesAccount, err := e.resolveAccount(int(order.TenantID), int(order.RestaurantID), accounts["sales"], "4100")
	if err != nil {
		return err
	}
	taxAccount, err := e.resolveAccount(int(order.TenantID), int(order.RestaurantID), accounts["tax"], "2120")
	if err != nil {
		return err
	}

	lines := make([]domain.JournalEntryLine, 0, 4)
	lines = append(lines, domain.JournalEntryLine{
		AccountID:   collectionAccount.ID,
		DebitAmount: roundedAmount(order.TotalAmount),
		Description: fmt.Sprintf("Order payment %s", order.OrderNumber),
		LineOrder:   len(lines) + 1,
	})
	if order.DiscountAmount > 0 {
		discountAccount, err := e.resolveAccount(int(order.TenantID), int(order.RestaurantID), accounts["discount"], "4150")
		if err != nil {
			return err
		}
		lines = append(lines, domain.JournalEntryLine{
			AccountID:   discountAccount.ID,
			DebitAmount: roundedAmount(order.DiscountAmount),
			Description: fmt.Sprintf("Sales discount %s", order.OrderNumber),
			LineOrder:   len(lines) + 1,
		})
	}
	if order.Subtotal > 0 {
		lines = append(lines, domain.JournalEntryLine{
			AccountID:    salesAccount.ID,
			CreditAmount: roundedAmount(order.Subtotal),
			Description:  fmt.Sprintf("Sales revenue %s", order.OrderNumber),
			LineOrder:    len(lines) + 1,
		})
	}
	if order.DeliveryFee > 0 {
		deliveryAccount, err := e.resolveAccount(int(order.TenantID), int(order.RestaurantID), accounts["delivery_sales"], "4120", "4100")
		if err != nil {
			return err
		}
		lines = append(lines, domain.JournalEntryLine{
			AccountID:    deliveryAccount.ID,
			CreditAmount: roundedAmount(order.DeliveryFee),
			Description:  fmt.Sprintf("Delivery revenue %s", order.OrderNumber),
			LineOrder:    len(lines) + 1,
		})
	}
	if order.TaxAmount > 0 {
		lines = append(lines, domain.JournalEntryLine{
			AccountID:    taxAccount.ID,
			CreditAmount: roundedAmount(order.TaxAmount),
			Description:  fmt.Sprintf("Sales tax %s", order.OrderNumber),
			LineOrder:    len(lines) + 1,
		})
	}

	description := fmt.Sprintf("Auto-posted sales for order %s", order.OrderNumber)
	_, err = e.createAndPostEntry(
		int(order.TenantID),
		int(order.RestaurantID),
		domain.PostingEvent{
			Code:            domain.AutoPostingEventOrderPaid,
			EntryType:       domain.AutoPostingEntryTypeSales,
			ReferenceType:   domain.AutoPostingReferenceTypeOrder,
			ReferenceID:     int(order.ID),
			ReferenceNumber: order.OrderNumber,
			Description:     description,
			PostedBy:        postedBy,
		},
		order.CreatedAt,
		lines,
	)
	return err
}

func (e *AutoPostingEngine) PostCOGSEntry(order *domain.Order, postedBy *int) error {
	if order == nil {
		return errors.New("order is required")
	}
	if !strings.EqualFold(strings.TrimSpace(order.PaymentStatus), "paid") {
		return nil
	}

	accounts, enabled, err := e.ruleAccounts(int(order.TenantID), int(order.RestaurantID), domain.AutoPostingEventOrderCOGS)
	if err != nil || !enabled {
		return err
	}

	existingID, err := e.autoPostingRepo.FindJournalEntryByReference(
		int(order.TenantID),
		int(order.RestaurantID),
		domain.AutoPostingReferenceTypeOrder,
		domain.AutoPostingEntryTypeCOGS,
		int(order.ID),
	)
	if err != nil {
		return err
	}
	if existingID != nil {
		return nil
	}

	cogsAccount, err := e.resolveAccount(int(order.TenantID), int(order.RestaurantID), accounts["food_cogs"], "5110")
	if err != nil {
		return err
	}
	inventoryAccount, err := e.resolveAccount(int(order.TenantID), int(order.RestaurantID), accounts["inventory"], "1130")
	if err != nil {
		return err
	}

	totalCost := 0.0
	movements := make([]domain.InventoryMovement, 0)
	for _, item := range order.Items {
		product, err := e.productRepo.GetProductByID(int(order.TenantID), int(order.RestaurantID), int(item.ProductID))
		if err != nil {
			return fmt.Errorf("load product %d: %w", item.ProductID, err)
		}
		if !product.TrackInventory {
			continue
		}

		unitCost, err := e.inventoryMovementRepo.GetProductCost(int(order.TenantID), int(order.RestaurantID), item.ProductID)
		if err != nil {
			return err
		}
		if unitCost <= 0 {
			continue
		}

		itemCost := roundedAmount(unitCost * float64(item.Quantity))
		if itemCost <= 0 {
			continue
		}

		totalCost += itemCost
		refID := int(order.ID)
		movements = append(movements, domain.InventoryMovement{
			TenantID:      int(order.TenantID),
			RestaurantID:  int(order.RestaurantID),
			ProductID:     item.ProductID,
			OrderID:       &order.ID,
			ReferenceType: domain.AutoPostingReferenceTypeOrder,
			ReferenceID:   &refID,
			MovementType:  domain.InventoryMovementTypeSale,
			Quantity:      -float64(item.Quantity),
			UnitCost:      roundedAmount(unitCost),
			TotalCost:     itemCost,
			Notes:         fmt.Sprintf("Auto-posted sale movement for order %s", order.OrderNumber),
			CreatedBy:     postedBy,
		})
	}

	totalCost = roundedAmount(totalCost)
	if totalCost <= 0 {
		return nil
	}

	description := fmt.Sprintf("Auto-posted COGS for order %s", order.OrderNumber)
	if _, err := e.createAndPostEntry(
		int(order.TenantID),
		int(order.RestaurantID),
		domain.PostingEvent{
			Code:            domain.AutoPostingEventOrderCOGS,
			EntryType:       domain.AutoPostingEntryTypeCOGS,
			ReferenceType:   domain.AutoPostingReferenceTypeOrder,
			ReferenceID:     int(order.ID),
			ReferenceNumber: order.OrderNumber,
			Description:     description,
			PostedBy:        postedBy,
		},
		order.CreatedAt,
		[]domain.JournalEntryLine{
			{
				AccountID:   cogsAccount.ID,
				DebitAmount: totalCost,
				Description: fmt.Sprintf("COGS for order %s", order.OrderNumber),
				LineOrder:   1,
			},
			{
				AccountID:    inventoryAccount.ID,
				CreditAmount: totalCost,
				Description:  fmt.Sprintf("Inventory reduction for order %s", order.OrderNumber),
				LineOrder:    2,
			},
		},
	); err != nil {
		return err
	}

	for _, movement := range movements {
		if _, err := e.inventoryMovementRepo.CreateMovement(&movement); err != nil {
			return err
		}
	}

	return nil
}

func (e *AutoPostingEngine) PostPayrollEntry(salary *domain.Salary, postedBy *int) error {
	if salary == nil {
		return errors.New("salary is required")
	}
	if !salary.IsPaid && !strings.EqualFold(strings.TrimSpace(salary.Status), "paid") {
		return nil
	}

	accounts, enabled, err := e.ruleAccounts(salary.TenantID, salary.RestaurantID, domain.AutoPostingEventSalaryPaid)
	if err != nil || !enabled {
		return err
	}

	existingID, err := e.autoPostingRepo.FindJournalEntryByReference(
		salary.TenantID,
		salary.RestaurantID,
		domain.AutoPostingReferenceTypeSalary,
		domain.AutoPostingEntryTypePayroll,
		salary.ID,
	)
	if err != nil {
		return err
	}
	if existingID != nil {
		return nil
	}

	paymentKey := "bank"
	if strings.EqualFold(strings.TrimSpace(salary.PaymentMethod), "cash") {
		paymentKey = "cash"
	}
	paymentAccount, err := e.resolvePaymentAccount(salary.TenantID, salary.RestaurantID, paymentKey, accounts)
	if err != nil {
		return err
	}

	salaryExpenseAccount, err := e.resolveAccount(salary.TenantID, salary.RestaurantID, accounts["salary_expense"], "5210")
	if err != nil {
		return err
	}

	baseExpense := roundedAmount(salary.GrossSalary - salary.OvertimeAmount - salary.Tips)
	if baseExpense < 0 {
		baseExpense = 0
	}

	lines := make([]domain.JournalEntryLine, 0, 10)
	if baseExpense > 0 {
		lines = append(lines, domain.JournalEntryLine{
			AccountID:   salaryExpenseAccount.ID,
			DebitAmount: baseExpense,
			Description: fmt.Sprintf("Payroll expense salary #%d", salary.ID),
			LineOrder:   len(lines) + 1,
		})
	}
	if salary.OvertimeAmount > 0 {
		overtimeExpenseAccount, err := e.resolveAccount(salary.TenantID, salary.RestaurantID, accounts["overtime_expense"], "5220", "5210")
		if err != nil {
			return err
		}
		lines = append(lines, domain.JournalEntryLine{
			AccountID:   overtimeExpenseAccount.ID,
			DebitAmount: roundedAmount(salary.OvertimeAmount),
			Description: fmt.Sprintf("Overtime salary #%d", salary.ID),
			LineOrder:   len(lines) + 1,
		})
	}
	if salary.Tips > 0 {
		tipsPayableAccount, err := e.resolveAccount(salary.TenantID, salary.RestaurantID, accounts["tips_payable"], "2140", "5210")
		if err != nil {
			return err
		}
		lines = append(lines, domain.JournalEntryLine{
			AccountID:   tipsPayableAccount.ID,
			DebitAmount: roundedAmount(salary.Tips),
			Description: fmt.Sprintf("Tips payout salary #%d", salary.ID),
			LineOrder:   len(lines) + 1,
		})
	}
	if salary.Tax > 0 {
		taxPayableAccount, err := e.resolveAccount(salary.TenantID, salary.RestaurantID, accounts["tax_payable"], accounts["accrued_expenses"], "2110")
		if err != nil {
			return err
		}
		lines = append(lines, domain.JournalEntryLine{
			AccountID:    taxPayableAccount.ID,
			CreditAmount: roundedAmount(salary.Tax),
			Description:  fmt.Sprintf("Payroll tax salary #%d", salary.ID),
			LineOrder:    len(lines) + 1,
		})
	}
	if salary.SocialInsurance > 0 {
		socialPayableAccount, err := e.resolveAccount(salary.TenantID, salary.RestaurantID, accounts["social_insurance_payable"], accounts["accrued_expenses"], "2110")
		if err != nil {
			return err
		}
		lines = append(lines, domain.JournalEntryLine{
			AccountID:    socialPayableAccount.ID,
			CreditAmount: roundedAmount(salary.SocialInsurance),
			Description:  fmt.Sprintf("Social insurance salary #%d", salary.ID),
			LineOrder:    len(lines) + 1,
		})
	}
	if salary.HealthInsurance > 0 {
		healthPayableAccount, err := e.resolveAccount(salary.TenantID, salary.RestaurantID, accounts["health_insurance_payable"], accounts["accrued_expenses"], "2110")
		if err != nil {
			return err
		}
		lines = append(lines, domain.JournalEntryLine{
			AccountID:    healthPayableAccount.ID,
			CreditAmount: roundedAmount(salary.HealthInsurance),
			Description:  fmt.Sprintf("Health insurance salary #%d", salary.ID),
			LineOrder:    len(lines) + 1,
		})
	}
	if salary.Pension > 0 {
		pensionPayableAccount, err := e.resolveAccount(salary.TenantID, salary.RestaurantID, accounts["pension_payable"], accounts["accrued_expenses"], "2110")
		if err != nil {
			return err
		}
		lines = append(lines, domain.JournalEntryLine{
			AccountID:    pensionPayableAccount.ID,
			CreditAmount: roundedAmount(salary.Pension),
			Description:  fmt.Sprintf("Pension salary #%d", salary.ID),
			LineOrder:    len(lines) + 1,
		})
	}
	if salary.LoanDeduction > 0 {
		loanClearingAccount, err := e.resolveAccount(salary.TenantID, salary.RestaurantID, accounts["loan_clearing"], accounts["accrued_expenses"], "2110")
		if err != nil {
			return err
		}
		lines = append(lines, domain.JournalEntryLine{
			AccountID:    loanClearingAccount.ID,
			CreditAmount: roundedAmount(salary.LoanDeduction),
			Description:  fmt.Sprintf("Loan deduction salary #%d", salary.ID),
			LineOrder:    len(lines) + 1,
		})
	}
	if salary.AdvanceDeduction > 0 {
		advanceAccount, err := e.resolveAccount(salary.TenantID, salary.RestaurantID, accounts["employee_advances"], "1121", "2110")
		if err != nil {
			return err
		}
		lines = append(lines, domain.JournalEntryLine{
			AccountID:    advanceAccount.ID,
			CreditAmount: roundedAmount(salary.AdvanceDeduction),
			Description:  fmt.Sprintf("Advance deduction salary #%d", salary.ID),
			LineOrder:    len(lines) + 1,
		})
	}

	otherDeductions := roundedAmount(salary.OtherDeductions)
	if otherDeductions > 0 {
		otherDeductionsAccount, err := e.resolveAccount(salary.TenantID, salary.RestaurantID, accounts["other_deductions_payable"], accounts["accrued_expenses"], "2110")
		if err != nil {
			return err
		}
		lines = append(lines, domain.JournalEntryLine{
			AccountID:    otherDeductionsAccount.ID,
			CreditAmount: otherDeductions,
			Description:  fmt.Sprintf("Other deductions salary #%d", salary.ID),
			LineOrder:    len(lines) + 1,
		})
	}
	if salary.NetSalary > 0 {
		lines = append(lines, domain.JournalEntryLine{
			AccountID:    paymentAccount.ID,
			CreditAmount: roundedAmount(salary.NetSalary),
			Description:  fmt.Sprintf("Salary payment #%d", salary.ID),
			LineOrder:    len(lines) + 1,
		})
	}

	description := fmt.Sprintf("Auto-posted payroll for salary #%d", salary.ID)
	_, err = e.createAndPostEntry(
		salary.TenantID,
		salary.RestaurantID,
		domain.PostingEvent{
			Code:            domain.AutoPostingEventSalaryPaid,
			EntryType:       domain.AutoPostingEntryTypePayroll,
			ReferenceType:   domain.AutoPostingReferenceTypeSalary,
			ReferenceID:     salary.ID,
			ReferenceNumber: fmt.Sprintf("SAL-%d", salary.ID),
			Description:     description,
			PostedBy:        postedBy,
		},
		payrollEntryDate(salary),
		lines,
	)
	return err
}

func (e *AutoPostingEngine) PostSupplierInvoiceEntry(event *domain.PostingEvent, lines []domain.JournalEntryLine) error {
	if event == nil {
		return errors.New("posting event is required")
	}
	if event.TenantID <= 0 || event.RestaurantID <= 0 {
		return errors.New("tenant and restaurant are required")
	}
	event.Code = domain.AutoPostingEventSupplierInvoice
	if event.EntryType == "" {
		event.EntryType = domain.AutoPostingEntryTypeSupplierInvoice
	}
	_, err := e.createAndPostEntry(event.TenantID, event.RestaurantID, *event, time.Now(), lines)
	return err
}

func (e *AutoPostingEngine) PostSupplierPaymentEntry(event *domain.PostingEvent, lines []domain.JournalEntryLine) error {
	if event == nil {
		return errors.New("posting event is required")
	}
	if event.TenantID <= 0 || event.RestaurantID <= 0 {
		return errors.New("tenant and restaurant are required")
	}
	event.Code = domain.AutoPostingEventSupplierPayment
	if event.EntryType == "" {
		event.EntryType = domain.AutoPostingEntryTypeSupplierPayment
	}
	_, err := e.createAndPostEntry(event.TenantID, event.RestaurantID, *event, time.Now(), lines)
	return err
}

func (e *AutoPostingEngine) PostDepreciationEntries(entryDate time.Time, tenantID, restaurantID int, description string, postedBy *int, lines []domain.JournalEntryLine) error {
	if strings.TrimSpace(description) == "" {
		description = "Auto-posted depreciation"
	}
	_, err := e.createAndPostEntry(
		tenantID,
		restaurantID,
		domain.PostingEvent{
			Code:          domain.AutoPostingEventDepreciation,
			EntryType:     domain.AutoPostingEntryTypeDepreciation,
			ReferenceType: domain.AutoPostingReferenceTypeAsset,
			Description:   description,
			PostedBy:      postedBy,
		},
		entryDate,
		lines,
	)
	return err
}

func (e *AutoPostingEngine) createAndPostEntry(
	tenantID, restaurantID int,
	event domain.PostingEvent,
	entryDate time.Time,
	lines []domain.JournalEntryLine,
) (int, error) {
	if tenantID <= 0 || restaurantID <= 0 {
		return 0, errors.New("tenant and restaurant are required")
	}
	if len(lines) < 2 {
		return 0, errors.New("at least two journal lines are required")
	}
	validatedLines, err := validateAutoPostingLines(lines)
	if err != nil {
		return 0, err
	}

	if entryDate.IsZero() {
		entryDate = time.Now()
	}

	createdBy := event.PostedBy
	entry := domain.JournalEntry{
		TenantID:        tenantID,
		RestaurantID:    restaurantID,
		EntryDate:       entryDate,
		Description:     event.Description,
		ReferenceType:   event.ReferenceType,
		ReferenceID:     intPtr(event.ReferenceID),
		ReferenceNumber: event.ReferenceNumber,
		EntryType:       event.EntryType,
		Status:          domain.JournalStatusDraft,
		CreatedBy:       createdBy,
		UpdatedBy:       createdBy,
		Lines:           validatedLines,
	}

	id, err := e.journalRepo.CreateJournalEntry(&entry)
	if err != nil {
		return 0, err
	}

	postedBy := 0
	if event.PostedBy != nil {
		postedBy = *event.PostedBy
	}
	if err := e.generalLedgerRepo.AppendEntry(tenantID, restaurantID, id, postedBy); err != nil {
		return id, err
	}

	log.Printf("[AccountingAutoPost] posted code=%s entry_type=%s ref=%s:%d tenant=%d restaurant=%d journal_id=%d",
		event.Code,
		event.EntryType,
		event.ReferenceType,
		event.ReferenceID,
		tenantID,
		restaurantID,
		id,
	)

	return id, nil
}

func (e *AutoPostingEngine) ruleAccounts(tenantID, restaurantID int, ruleCode string) (map[string]string, bool, error) {
	rule, err := e.autoPostingRepo.GetRuleByCode(tenantID, restaurantID, ruleCode)
	if err != nil {
		return nil, false, err
	}
	accounts := make(map[string]string, len(domain.DefaultPostingRuleAccounts[ruleCode]))
	for key, value := range domain.DefaultPostingRuleAccounts[ruleCode] {
		accounts[key] = value
	}
	if rule == nil {
		return accounts, true, nil
	}
	for key, value := range rule.DefaultAccounts {
		accounts[key] = value
	}
	return accounts, rule.IsActive, nil
}

func (e *AutoPostingEngine) resolvePaymentAccount(tenantID, restaurantID int, paymentMethod string, accounts map[string]string) (*domain.Account, error) {
	switch paymentMethod {
	case "cash":
		return e.resolveAccount(tenantID, restaurantID, accounts["cash"], "1110")
	case "card":
		return e.resolveAccount(tenantID, restaurantID, accounts["card"], "1113", "1111")
	case "online":
		return e.resolveAccount(tenantID, restaurantID, accounts["online"], accounts["bank"], "1111")
	case "wallet":
		return e.resolveAccount(tenantID, restaurantID, accounts["wallet"], accounts["bank"], "1111")
	case "bank_transfer", "check", "mobile_money", "bank":
		return e.resolveAccount(tenantID, restaurantID, accounts["bank"], "1111")
	default:
		return e.resolveAccount(tenantID, restaurantID, accounts["receivable"], "1120", "1111")
	}
}

func (e *AutoPostingEngine) resolveAccount(tenantID, restaurantID int, candidates ...string) (*domain.Account, error) {
	for _, code := range candidates {
		code = strings.TrimSpace(code)
		if code == "" {
			continue
		}
		account, err := e.accountRepo.GetAccountByCode(tenantID, restaurantID, code)
		if err != nil {
			return nil, err
		}
		if account != nil {
			return account, nil
		}
	}
	return nil, fmt.Errorf("account not found for candidates %v", candidates)
}

func roundedAmount(value float64) float64 {
	return math.Round(value*100) / 100
}

func intPtr(value int) *int {
	if value <= 0 {
		return nil
	}
	v := value
	return &v
}

func payrollEntryDate(salary *domain.Salary) time.Time {
	if salary == nil {
		return time.Now()
	}
	if salary.PaymentDate != nil && !salary.PaymentDate.IsZero() {
		return *salary.PaymentDate
	}
	if salary.PaidAt != nil && !salary.PaidAt.IsZero() {
		return *salary.PaidAt
	}
	return time.Now()
}
