package service

import (
	"strings"
	"sync"
	"time"
)

// PaymentWebhookEvent is a lightweight runtime event for webhook observability.
type PaymentWebhookEvent struct {
	Timestamp          time.Time `json:"timestamp"`
	PaymobOrderID      int64     `json:"paymob_order_id"`
	Status             string    `json:"status"`
	HMACValid          bool      `json:"hmac_valid"`
	ProcessingTimeMS   int64     `json:"processing_time_ms"`
	Source             string    `json:"source"`
	Error              string    `json:"error,omitempty"`
	RetryQueueRecordID *int64    `json:"retry_queue_record_id,omitempty"`
}

// PaymentWebhookMonitor stores recent webhook events in memory for dashboards/metrics.
type PaymentWebhookMonitor struct {
	mu        sync.RWMutex
	maxEvents int
	events    []PaymentWebhookEvent
}

// NewPaymentWebhookMonitor creates an in-memory monitor with bounded history.
func NewPaymentWebhookMonitor(maxEvents int) *PaymentWebhookMonitor {
	if maxEvents <= 0 {
		maxEvents = 200
	}

	return &PaymentWebhookMonitor{
		maxEvents: maxEvents,
		events:    make([]PaymentWebhookEvent, 0, maxEvents),
	}
}

// Record appends an event and keeps only the latest maxEvents entries.
func (m *PaymentWebhookMonitor) Record(event PaymentWebhookEvent) {
	if m == nil {
		return
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}
	event.Status = strings.ToLower(strings.TrimSpace(event.Status))
	if event.Status == "" {
		event.Status = "unknown"
	}
	event.Source = strings.TrimSpace(event.Source)
	if event.Source == "" {
		event.Source = "webhook"
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.events = append(m.events, event)
	if len(m.events) > m.maxEvents {
		excess := len(m.events) - m.maxEvents
		m.events = m.events[excess:]
	}
}

// List returns the latest events in reverse chronological order (newest first).
func (m *PaymentWebhookMonitor) List(limit int) []PaymentWebhookEvent {
	if m == nil {
		return []PaymentWebhookEvent{}
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	if limit <= 0 || limit > len(m.events) {
		limit = len(m.events)
	}

	out := make([]PaymentWebhookEvent, 0, limit)
	for i := len(m.events) - 1; i >= 0 && len(out) < limit; i-- {
		out = append(out, m.events[i])
	}
	return out
}

// Aggregate returns total, success, failed, and average processing duration for a period.
func (m *PaymentWebhookMonitor) Aggregate(since *time.Time) (int, int, int, float64) {
	if m == nil {
		return 0, 0, 0, 0
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var total int
	var success int
	var failed int
	var durationSum int64

	for _, event := range m.events {
		if since != nil && event.Timestamp.Before(*since) {
			continue
		}
		total++
		durationSum += event.ProcessingTimeMS

		if event.Status == "success" {
			success++
		} else if strings.Contains(event.Status, "fail") ||
			strings.Contains(event.Status, "error") ||
			strings.Contains(event.Status, "invalid") ||
			strings.Contains(event.Status, "reject") {
			failed++
		}
	}

	avg := 0.0
	if total > 0 {
		avg = float64(durationSum) / float64(total)
	}

	return total, success, failed, avg
}
