package observability

import (
	"sort"
	"sync"
	"time"
)

const maxSamplesPerMetric = 240

type metricWindow struct {
	samples []float64
}

func newMetricWindow() *metricWindow {
	return &metricWindow{
		samples: make([]float64, 0, maxSamplesPerMetric),
	}
}

func (w *metricWindow) add(valueMs float64) {
	if valueMs < 0 {
		return
	}

	w.samples = append(w.samples, valueMs)
	if len(w.samples) > maxSamplesPerMetric {
		overflow := len(w.samples) - maxSamplesPerMetric
		w.samples = w.samples[overflow:]
	}
}

type MetricSummary struct {
	Count  int64   `json:"count"`
	AvgMs  float64 `json:"avg_ms"`
	P95Ms  float64 `json:"p95_ms"`
	MaxMs  float64 `json:"max_ms"`
	LastMs float64 `json:"last_ms"`
}

type MetricsSnapshot struct {
	GeneratedAt      string                   `json:"generated_at"`
	UptimeSec        int64                    `json:"uptime_sec"`
	API              map[string]MetricSummary `json:"api"`
	DB               map[string]MetricSummary `json:"db"`
	AccountingEvents *AccountingEventMetrics  `json:"accounting_events,omitempty"`
}

type AccountingEventMetrics struct {
	ObservedAt                   string `json:"observed_at"`
	MissingIntents               int64  `json:"missing_intents"`
	Pending                      int64  `json:"pending"`
	Processing                   int64  `json:"processing"`
	Retryable                    int64  `json:"retryable"`
	Failed                       int64  `json:"failed"`
	CompletedWithoutSalesJournal int64  `json:"completed_without_sales_journal"`
}

type collector struct {
	mu               sync.RWMutex
	startedAt        time.Time
	api              map[string]*metricWindow
	db               map[string]*metricWindow
	accountingEvents *AccountingEventMetrics
}

func newCollector() *collector {
	return &collector{
		startedAt: time.Now(),
		api:       make(map[string]*metricWindow),
		db:        make(map[string]*metricWindow),
	}
}

var globalCollector = newCollector()

func percentile(values []float64, p float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	index := int(((p / 100) * float64(len(sorted))))
	if index <= 0 {
		return sorted[0]
	}
	if index >= len(sorted) {
		return sorted[len(sorted)-1]
	}
	return sorted[index-1]
}

func summarize(window *metricWindow) MetricSummary {
	if window == nil || len(window.samples) == 0 {
		return MetricSummary{}
	}

	var total float64
	max := window.samples[0]
	for _, value := range window.samples {
		total += value
		if value > max {
			max = value
		}
	}

	return MetricSummary{
		Count:  int64(len(window.samples)),
		AvgMs:  total / float64(len(window.samples)),
		P95Ms:  percentile(window.samples, 95),
		MaxMs:  max,
		LastMs: window.samples[len(window.samples)-1],
	}
}

// RecordAPI tracks request latency for API routes.
func RecordAPI(metric string, duration time.Duration) {
	if metric == "" {
		metric = "unknown"
	}

	globalCollector.mu.Lock()
	defer globalCollector.mu.Unlock()

	window := globalCollector.api[metric]
	if window == nil {
		window = newMetricWindow()
		globalCollector.api[metric] = window
	}
	window.add(float64(duration.Milliseconds()))
}

// RecordDB tracks database query/operation latency.
func RecordDB(metric string, duration time.Duration) {
	if metric == "" {
		metric = "unknown"
	}

	globalCollector.mu.Lock()
	defer globalCollector.mu.Unlock()

	window := globalCollector.db[metric]
	if window == nil {
		window = newMetricWindow()
		globalCollector.db[metric] = window
	}
	window.add(float64(duration.Milliseconds()))
}

// RecordAccountingEvents publishes the latest durable-recovery reconciliation state.
func RecordAccountingEvents(snapshot AccountingEventMetrics) {
	snapshot.ObservedAt = time.Now().Format(time.RFC3339)
	globalCollector.mu.Lock()
	defer globalCollector.mu.Unlock()
	globalCollector.accountingEvents = &snapshot
}

// Snapshot returns an in-memory performance summary including p95 values.
func Snapshot() MetricsSnapshot {
	globalCollector.mu.RLock()
	defer globalCollector.mu.RUnlock()

	api := make(map[string]MetricSummary, len(globalCollector.api))
	for key, window := range globalCollector.api {
		api[key] = summarize(window)
	}

	db := make(map[string]MetricSummary, len(globalCollector.db))
	for key, window := range globalCollector.db {
		db[key] = summarize(window)
	}
	var accountingEvents *AccountingEventMetrics
	if globalCollector.accountingEvents != nil {
		copy := *globalCollector.accountingEvents
		accountingEvents = &copy
	}

	return MetricsSnapshot{
		GeneratedAt:      time.Now().Format(time.RFC3339),
		UptimeSec:        int64(time.Since(globalCollector.startedAt).Seconds()),
		API:              api,
		DB:               db,
		AccountingEvents: accountingEvents,
	}
}
