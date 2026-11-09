package http

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pos-saas/internal/service"
)

// PaymentMonitoringHandler serves super-admin payment/webhook monitoring endpoints.
type PaymentMonitoringHandler struct {
	db             *sql.DB
	webhookMonitor *service.PaymentWebhookMonitor
	retryQueue     *service.WebhookRetryQueue
}

// NewPaymentMonitoringHandler creates a monitoring handler.
func NewPaymentMonitoringHandler(
	db *sql.DB,
	webhookMonitor *service.PaymentWebhookMonitor,
	retryQueue *service.WebhookRetryQueue,
) *PaymentMonitoringHandler {
	return &PaymentMonitoringHandler{
		db:             db,
		webhookMonitor: webhookMonitor,
		retryQueue:     retryQueue,
	}
}

// GetPaymentMetrics handles GET /api/v1/super-admin/payment-metrics.
func (h *PaymentMonitoringHandler) GetPaymentMetrics(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		InternalErrorResponse(w, "payment database is not configured")
		return
	}

	period := strings.ToLower(strings.TrimSpace(GetQueryParam(r, "period", "month")))
	label, since := normalizeMetricsPeriod(period)

	query := `SELECT total_payments, successful_payments, failed_payments, total_revenue_cents
		FROM get_payment_metrics($1)`

	var total int64
	var successful int64
	var failed int64
	var revenue int64
	var sinceArg interface{}
	if since != nil {
		sinceArg = *since
	}
	if err := h.db.QueryRowContext(r.Context(), query, sinceArg).Scan(&total, &successful, &failed, &revenue); err != nil {
		log.Printf("[PaymentMonitoring] GetPaymentMetrics query error: %v", err)
		InternalErrorResponse(w, "failed to query payment metrics")
		return
	}

	successRate := 0.0
	if total > 0 {
		successRate = (float64(successful) / float64(total)) * 100
	}

	_, _, _, avgProcessing := h.webhookMonitor.Aggregate(since)

	OKResponse(w, map[string]interface{}{
		"total_payments":          total,
		"successful_payments":     successful,
		"failed_payments":         failed,
		"success_rate":            roundFloat(successRate, 2),
		"total_revenue_cents":     revenue,
		"total_revenue_formatted": formatCents(revenue, "EGP"),
		"avg_processing_time_ms":  roundFloat(avgProcessing, 2),
		"period":                  label,
		"generated_at":            time.Now().UTC().Format(time.RFC3339),
	})
}

// GetWebhookHealth handles GET /api/v1/super-admin/webhook-health.
func (h *PaymentMonitoringHandler) GetWebhookHealth(w http.ResponseWriter, r *http.Request) {
	period := strings.ToLower(strings.TrimSpace(GetQueryParam(r, "period", "month")))
	_, since := normalizeMetricsPeriod(period)

	limit := GetQueryParamInt(r, "limit", 50)
	if limit < 1 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	events := h.webhookMonitor.List(limit)
	filtered := make([]service.PaymentWebhookEvent, 0, len(events))
	for _, event := range events {
		if since != nil && event.Timestamp.Before(*since) {
			continue
		}
		filtered = append(filtered, event)
	}

	total, success, failed, avg := h.webhookMonitor.Aggregate(since)
	successRate := 0.0
	if total > 0 {
		successRate = (float64(success) / float64(total)) * 100
	}

	retryStats, err := h.retryQueue.GetStats(r.Context())
	if err != nil {
		log.Printf("[PaymentMonitoring] GetWebhookHealth stats error: %v", err)
		InternalErrorResponse(w, "failed to load webhook retry stats")
		return
	}

	OKResponse(w, map[string]interface{}{
		"events": filtered,
		"summary": map[string]interface{}{
			"total_events":           total,
			"successful_events":      success,
			"failed_events":          failed,
			"success_rate":           roundFloat(successRate, 2),
			"avg_processing_time_ms": roundFloat(avg, 2),
		},
		"retry_queue":  retryStats,
		"generated_at": time.Now().UTC().Format(time.RFC3339),
	})
}

func normalizeMetricsPeriod(period string) (string, *time.Time) {
	now := time.Now().UTC()

	switch period {
	case "today":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		return "today", &start
	case "week":
		start := now.AddDate(0, 0, -7)
		return "week", &start
	case "month":
		start := now.AddDate(0, -1, 0)
		return "month", &start
	default:
		return "all", nil
	}
}

func roundFloat(value float64, precision int) float64 {
	if precision < 0 {
		precision = 0
	}
	mul := 1.0
	for i := 0; i < precision; i++ {
		mul *= 10
	}
	return float64(int64(value*mul+0.5)) / mul
}

func formatCents(cents int64, currency string) string {
	major := cents / 100
	rem := cents % 100
	sign := ""
	if cents < 0 {
		sign = "-"
		major = -major
		rem = -rem
	}
	return fmt.Sprintf("%s%s.%02d %s", sign, commaInt(major), rem, strings.ToUpper(strings.TrimSpace(currency)))
}

func commaInt(value int64) string {
	s := strconv.FormatInt(value, 10)
	if len(s) <= 3 {
		return s
	}

	var out strings.Builder
	pre := len(s) % 3
	if pre > 0 {
		out.WriteString(s[:pre])
		if len(s) > pre {
			out.WriteString(",")
		}
	}

	for i := pre; i < len(s); i += 3 {
		out.WriteString(s[i : i+3])
		if i+3 < len(s) {
			out.WriteString(",")
		}
	}
	return out.String()
}
