package repository

import (
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"pos-saas/internal/domain"
)

func roundAccountingAmount(value float64) float64 {
	return math.Round(value*100) / 100
}

func effectiveInvoiceStatus(current string, dueDate time.Time, balanceDue float64, sentState string) string {
	current = strings.TrimSpace(strings.ToLower(current))
	if current == domain.InvoiceStatusVoid || current == domain.InvoiceStatusPaid {
		return current
	}
	if balanceDue <= 0 {
		return domain.InvoiceStatusPaid
	}
	if balanceDue > 0 && current == domain.InvoiceStatusPartiallyPaid {
		if !dueDate.IsZero() && dueDate.Before(time.Now()) {
			return domain.InvoiceStatusOverdue
		}
		return domain.InvoiceStatusPartiallyPaid
	}
	if !dueDate.IsZero() && dueDate.Before(time.Now()) && current != domain.InvoiceStatusDraft {
		return domain.InvoiceStatusOverdue
	}
	if current == "" {
		return sentState
	}
	return current
}

func fallbackCollectionAccountCode(paymentMethod string) string {
	switch strings.ToLower(strings.TrimSpace(paymentMethod)) {
	case "cash":
		return "1110"
	case "card":
		return "1113"
	default:
		return "1111"
	}
}

func lookupAccountByCode(accountRepo *AccountRepository, tenantID, restaurantID int, code string) (*domain.Account, error) {
	if accountRepo == nil {
		return nil, fmt.Errorf("account repository not configured")
	}
	account, err := accountRepo.GetAccountByCode(tenantID, restaurantID, code)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, fmt.Errorf("account %s not found", code)
	}
	return account, nil
}

func scanNullableInt(scanner interface {
	Scan(dest ...interface{}) error
}) (*int, error) {
	var value sql.NullInt64
	if err := scanner.Scan(&value); err != nil {
		return nil, err
	}
	if !value.Valid {
		return nil, nil
	}
	converted := int(value.Int64)
	return &converted, nil
}
