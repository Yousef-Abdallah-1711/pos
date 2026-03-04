package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"pos-saas/internal/domain"
)

func parseAccountingInputDate(value, field string) (time.Time, error) {
	date, err := parseOptionalDate(value)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid %s format", field)
	}
	return date, nil
}

func requireAccountingString(value, field string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", field)
	}
	return nil
}

func mapAccountingInvoiceLines(inputs []domain.InvoiceLineInput) []domain.InvoiceLine {
	lines := make([]domain.InvoiceLine, 0, len(inputs))
	for _, line := range inputs {
		lines = append(lines, domain.InvoiceLine{
			Description:    strings.TrimSpace(line.Description),
			DescriptionAr:  strings.TrimSpace(line.DescriptionAr),
			Quantity:       line.Quantity,
			UnitPrice:      line.UnitPrice,
			TaxRate:        line.TaxRate,
			TaxAmount:      line.TaxAmount,
			DiscountRate:   line.DiscountRate,
			DiscountAmount: line.DiscountAmount,
			LineTotal:      line.LineTotal,
			AccountID:      line.AccountID,
			LineOrder:      line.LineOrder,
		})
	}
	return lines
}

func decodeAccountingActionBody(r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return nil
	}
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil && err != io.EOF {
		return err
	}
	return nil
}

func writeAccountingMutationError(w http.ResponseWriter, err error) {
	statusCode := http.StatusBadRequest
	if strings.Contains(err.Error(), "not found") {
		statusCode = http.StatusNotFound
	}
	writeJSON(w, statusCode, map[string]interface{}{
		"success": false,
		"error":   err.Error(),
	})
}
