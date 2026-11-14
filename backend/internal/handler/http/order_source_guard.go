package http

import (
	"net/http"
	"strings"

	"pos-saas/internal/domain"
)

// POS source triggers immediate confirmation and stock reservation. Only the
// permission-protected route selected by ServeMux may choose that workflow.
func applyTrustedOrderSource(r *http.Request, request *domain.CreateOrderRequest) bool {
	if r.Pattern == "POST /api/v1/pos/orders" {
		request.OrderSource = "pos"
		return true
	}
	if strings.EqualFold(strings.TrimSpace(request.OrderSource), "pos") {
		return false
	}
	if request.OrderSource == "" {
		request.OrderSource = "website"
	}
	return true
}
