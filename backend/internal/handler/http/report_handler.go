package http

import (
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/netip"
	"sort"
	"strings"
	"sync"
	"time"

	"pos-saas/internal/middleware"

	goarabic "github.com/01walid/goarabic"
	"github.com/go-pdf/fpdf"
)

//go:embed report_fonts/NotoNaskhArabic.ttf
var reportArabicFontTTF []byte

// ReportHandler serves analytics/report endpoints for the dashboard.
type ReportHandler struct {
	db *sql.DB

	cache *reportResponseCache

	// test hooks / overrides for endpoint tests
	salesReportDataBuilder         salesReportDataBuilderFn
	profitabilityReportDataBuilder profitabilityReportDataBuilderFn
	inventoryReportDataBuilder     inventoryReportDataBuilderFn
	reportExportAuditLogger        reportExportAuditLoggerFn
}

func NewReportHandler(db *sql.DB) *ReportHandler {
	return &ReportHandler{
		db:    db,
		cache: newReportResponseCache(),
	}
}

type reportResponseCache struct {
	mu    sync.RWMutex
	items map[string]cachedReportResponse
}

type cachedReportResponse struct {
	payload    []byte
	expiresAt  time.Time
	cachedAt   time.Time
	contentTyp string
}

func newReportResponseCache() *reportResponseCache {
	return &reportResponseCache{
		items: make(map[string]cachedReportResponse),
	}
}

func (c *reportResponseCache) Get(key string) ([]byte, string, bool) {
	if c == nil || key == "" {
		return nil, "", false
	}

	now := time.Now()
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()
	if !ok {
		return nil, "", false
	}
	if now.After(item.expiresAt) {
		c.mu.Lock()
		// Re-check under write lock before deleting.
		if latest, exists := c.items[key]; exists && now.After(latest.expiresAt) {
			delete(c.items, key)
		}
		c.mu.Unlock()
		return nil, "", false
	}
	return append([]byte(nil), item.payload...), item.contentTyp, true
}

func (c *reportResponseCache) Set(key string, payload []byte, contentType string, ttl time.Duration) {
	if c == nil || key == "" || ttl <= 0 || len(payload) == 0 {
		return
	}
	now := time.Now()
	c.mu.Lock()
	c.items[key] = cachedReportResponse{
		payload:    append([]byte(nil), payload...),
		expiresAt:  now.Add(ttl),
		cachedAt:   now,
		contentTyp: contentType,
	}
	c.mu.Unlock()
}

type salesReportDataBuilderFn func(
	tenantID, restaurantID int64,
	dateFrom, dateToExclusive, dateToDisplay time.Time,
	granularity string,
	topLimit int64,
) (salesReportResponse, error)

type profitabilityReportDataBuilderFn func(
	tenantID, restaurantID int64,
	dateFrom, dateToExclusive, dateToDisplay time.Time,
	topLimit int64,
) (profitabilityReportResponse, error)

type inventoryReportDataBuilderFn func(
	tenantID, restaurantID int64,
	dateFrom, dateToExclusive, dateToDisplay time.Time,
	topLimit int64,
	daysInRange int64,
) (inventoryReportResponse, error)

type reportExportAuditRecord struct {
	TenantID      int64
	RestaurantID  int64
	UserID        int64
	ReportType    string
	Format        string
	FileName      string
	ExportedBytes int
	Filters       map[string]interface{}
	RequestPath   string
	UserAgent     string
	IPAddress     string
}

type reportExportAuditLoggerFn func(ctx context.Context, record reportExportAuditRecord) error

type salesSummaryResponse struct {
	TotalRevenue float64 `json:"total_revenue"`
	TotalOrders  int64   `json:"total_orders"`
	AvgTicket    float64 `json:"avg_ticket"`
}

type salesTrendPoint struct {
	Bucket  string  `json:"bucket"`
	Revenue float64 `json:"revenue"`
	Orders  int64   `json:"orders"`
}

type salesHourlyPoint struct {
	Hour    int     `json:"hour"`
	Revenue float64 `json:"revenue"`
	Orders  int64   `json:"orders"`
}

type salesTopItem struct {
	ProductID int64   `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int64   `json:"quantity"`
	Revenue   float64 `json:"revenue"`
}

type salesTopCategory struct {
	CategoryID int64   `json:"category_id,omitempty"`
	Name       string  `json:"name"`
	Quantity   int64   `json:"quantity"`
	Revenue    float64 `json:"revenue"`
}

type salesReportResponse struct {
	Filters       map[string]interface{} `json:"filters"`
	Summary       salesSummaryResponse   `json:"summary"`
	Trend         []salesTrendPoint      `json:"trend"`
	Hourly        []salesHourlyPoint     `json:"hourly"`
	TopItems      []salesTopItem         `json:"top_items"`
	TopCategories []salesTopCategory     `json:"top_categories"`
	GeneratedAt   string                 `json:"generated_at"`
}

type profitabilitySummaryResponse struct {
	TotalRevenue            float64 `json:"total_revenue"`
	TotalCOGS               float64 `json:"total_cogs"`
	GrossProfit             float64 `json:"gross_profit"`
	GrossMarginPct          float64 `json:"gross_margin_pct"`
	FoodCostPct             float64 `json:"food_cost_pct"`
	UnitsMissingCost        int64   `json:"units_missing_cost"`
	ProductsMissingCostRows int64   `json:"products_missing_cost_rows"`
}

type profitabilityItem struct {
	ProductID       int64   `json:"product_id"`
	Name            string  `json:"name"`
	Quantity        int64   `json:"quantity"`
	Revenue         float64 `json:"revenue"`
	COGS            float64 `json:"cogs"`
	GrossProfit     float64 `json:"gross_profit"`
	GrossMarginPct  float64 `json:"gross_margin_pct"`
	FoodCostPct     float64 `json:"food_cost_pct"`
	CostPerUnit     float64 `json:"cost_per_unit"`
	MissingCostData bool    `json:"missing_cost_data"`
}

type profitabilityReportResponse struct {
	Filters     map[string]interface{}       `json:"filters"`
	Summary     profitabilitySummaryResponse `json:"summary"`
	TopItems    []profitabilityItem          `json:"top_items"`
	LowestItems []profitabilityItem          `json:"lowest_margin_items"`
	GeneratedAt string                       `json:"generated_at"`
}

type inventorySummaryResponse struct {
	TrackedProductsCount       int64   `json:"tracked_products_count"`
	TotalCurrentInventoryValue float64 `json:"total_current_inventory_value"`
	TotalEstimatedCOGS         float64 `json:"total_estimated_cogs"`
	AverageMonthlyTurnover     float64 `json:"average_monthly_turnover"`
	SlowMovingCount            int64   `json:"slow_moving_count"`
}

type inventoryTurnoverItem struct {
	ProductID             int64   `json:"product_id"`
	Name                  string  `json:"name"`
	CostPerUnit           float64 `json:"cost_per_unit"`
	CurrentStockQty       float64 `json:"current_stock_qty"`
	OpeningStockQty       float64 `json:"opening_stock_qty"`
	AverageInventoryQty   float64 `json:"average_inventory_qty"`
	AverageInventoryValue float64 `json:"average_inventory_value"`
	SoldQty               int64   `json:"sold_qty"`
	EstimatedCOGS         float64 `json:"estimated_cogs"`
	WasteQty              int64   `json:"waste_qty"`
	TurnoverRatio         float64 `json:"turnover_ratio"`
	MonthlyTurnoverRatio  float64 `json:"monthly_turnover_ratio"`
	IsSlowMoving          bool    `json:"is_slow_moving"`
}

type inventoryReportResponse struct {
	Filters     map[string]interface{}   `json:"filters"`
	Summary     inventorySummaryResponse `json:"summary"`
	Items       []inventoryTurnoverItem  `json:"items"`
	SlowMovers  []inventoryTurnoverItem  `json:"slow_movers"`
	TopTurnover []inventoryTurnoverItem  `json:"top_turnover"`
	GeneratedAt string                   `json:"generated_at"`
}

type reportsOverviewResponse struct {
	Filters map[string]interface{} `json:"filters"`
	Sales   struct {
		Summary salesSummaryResponse `json:"summary"`
		TopItem *salesTopItem        `json:"top_item,omitempty"`
	} `json:"sales"`
	Profitability struct {
		Summary        profitabilitySummaryResponse `json:"summary"`
		MostProfitable *profitabilityItem           `json:"most_profitable,omitempty"`
	} `json:"profitability"`
	Inventory struct {
		Summary       inventorySummaryResponse `json:"summary"`
		SlowestMoving *inventoryTurnoverItem   `json:"slowest_moving,omitempty"`
	} `json:"inventory"`
	GeneratedAt string `json:"generated_at"`
}

type reportExportRequest struct {
	ReportType   string `json:"report_type"`
	Format       string `json:"format"`
	DateFrom     string `json:"date_from"`
	DateTo       string `json:"date_to"`
	Granularity  string `json:"granularity"`
	TopLimit     int64  `json:"top_limit"`
	Locale       string `json:"locale"`
	IncludeItems bool   `json:"include_items"`
}

const maxReportDateRangeDays = 366
const reportDBQueryTimeout = 5 * time.Second
const (
	reportCacheTTLOverview      = 30 * time.Second
	reportCacheTTLSales         = 30 * time.Second
	reportCacheTTLProfitability = 60 * time.Second
	reportCacheTTLInventory     = 60 * time.Second
	reportCacheContentTypeJSON  = "application/json"
)

func withReportDBTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), reportDBQueryTimeout)
}

func withReportDBTimeoutFrom(parent context.Context) (context.Context, context.CancelFunc) {
	if parent == nil {
		return withReportDBTimeout()
	}
	return context.WithTimeout(parent, reportDBQueryTimeout)
}

func writeJSONPayload(w http.ResponseWriter, status int, payload []byte) {
	w.Header().Set("Content-Type", reportCacheContentTypeJSON)
	w.WriteHeader(status)
	_, _ = w.Write(payload)
}

func (h *ReportHandler) writeCachedOrStoreJSON(
	w http.ResponseWriter,
	cacheKey string,
	ttl time.Duration,
	responseBody map[string]interface{},
) error {
	payload, err := json.Marshal(responseBody)
	if err != nil {
		return fmt.Errorf("marshal response: %w", err)
	}
	if ttl > 0 {
		h.cache.Set(cacheKey, payload, reportCacheContentTypeJSON, ttl)
	}
	writeJSONPayload(w, http.StatusOK, payload)
	return nil
}

func (h *ReportHandler) logReportEndpointTiming(
	endpoint string,
	tenantID, restaurantID int64,
	startedAt time.Time,
	cacheHit bool,
	err error,
) {
	duration := time.Since(startedAt)
	status := "ok"
	if err != nil {
		status = "error"
	}
	log.Printf(
		"[ReportHandler] endpoint=%s tenant_id=%d restaurant_id=%d status=%s cache_hit=%t duration_ms=%d err=%v",
		endpoint,
		tenantID,
		restaurantID,
		status,
		cacheHit,
		duration.Milliseconds(),
		err,
	)
}

func reportCacheKey(parts ...string) string {
	return strings.Join(parts, "|")
}

// GetSalesReport returns aggregated sales analytics for the reports dashboard.
// GET /api/v1/reports/sales?date_from=2026-02-01&date_to=2026-02-29&granularity=day
func (h *ReportHandler) GetSalesReport(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		respondError(w, http.StatusServiceUnavailable, "Reports unavailable: database connection not initialized")
		return
	}

	tenantID, err := middleware.VerifyTenantID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}
	restaurantID, err := middleware.VerifyRestaurantID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Missing restaurant context")
		return
	}

	q := r.URL.Query()
	granularity := normalizeGranularity(q.Get("granularity"))
	dateFrom, dateToExclusive, dateToDisplay, err := parseReportDateRange(q.Get("date_from"), q.Get("date_to"))
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	topLimit := parseInt64(q.Get("top_limit"), 10)
	if topLimit <= 0 {
		topLimit = 10
	}
	if topLimit > 50 {
		topLimit = 50
	}

	startedAt := time.Now()
	cacheHit := false
	var endpointErr error
	defer func() {
		h.logReportEndpointTiming("sales", tenantID, restaurantID, startedAt, cacheHit, endpointErr)
	}()

	cacheKey := reportCacheKey(
		"reports:sales",
		fmt.Sprintf("tenant=%d", tenantID),
		fmt.Sprintf("restaurant=%d", restaurantID),
		fmt.Sprintf("date_from=%s", dateFrom.Format("2006-01-02")),
		fmt.Sprintf("date_to=%s", dateToDisplay.Format("2006-01-02")),
		fmt.Sprintf("granularity=%s", granularity),
		fmt.Sprintf("top_limit=%d", topLimit),
	)
	if payload, contentType, ok := h.cache.Get(cacheKey); ok {
		cacheHit = true
		if contentType == "" {
			contentType = reportCacheContentTypeJSON
		}
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(payload)
		return
	}

	summary, err := h.querySalesSummary(tenantID, restaurantID, dateFrom, dateToExclusive)
	if err != nil {
		endpointErr = err
		respondError(w, http.StatusInternalServerError, "Failed to load sales summary")
		return
	}

	trend, err := h.querySalesTrend(tenantID, restaurantID, dateFrom, dateToExclusive, granularity)
	if err != nil {
		endpointErr = err
		respondError(w, http.StatusInternalServerError, "Failed to load sales trend")
		return
	}

	hourly, err := h.querySalesHourly(tenantID, restaurantID, dateFrom, dateToExclusive)
	if err != nil {
		endpointErr = err
		respondError(w, http.StatusInternalServerError, "Failed to load hourly sales")
		return
	}

	topItems, err := h.queryTopItems(tenantID, restaurantID, dateFrom, dateToExclusive, topLimit)
	if err != nil {
		endpointErr = err
		respondError(w, http.StatusInternalServerError, "Failed to load top items")
		return
	}

	topCategories, err := h.queryTopCategories(tenantID, restaurantID, dateFrom, dateToExclusive, topLimit)
	if err != nil {
		endpointErr = err
		respondError(w, http.StatusInternalServerError, "Failed to load top categories")
		return
	}

	if err := h.writeCachedOrStoreJSON(w, cacheKey, reportCacheTTLSales, map[string]interface{}{
		"success": true,
		"data": salesReportResponse{
			Filters: map[string]interface{}{
				"date_from":         dateFrom.Format("2006-01-02"),
				"date_to":           dateToDisplay.Format("2006-01-02"),
				"granularity":       granularity,
				"top_limit":         topLimit,
				"included_statuses": []string{"all_except_cancelled"},
			},
			Summary:       summary,
			Trend:         trend,
			Hourly:        hourly,
			TopItems:      topItems,
			TopCategories: topCategories,
			GeneratedAt:   time.Now().UTC().Format(time.RFC3339),
		},
	}); err != nil {
		endpointErr = err
		respondError(w, http.StatusInternalServerError, "Failed to serialize sales report")
		return
	}
}

// GetProfitabilityReport returns product profitability analytics for a date range.
// GET /api/v1/reports/profitability?date_from=2026-02-01&date_to=2026-02-29&top_limit=10
func (h *ReportHandler) GetProfitabilityReport(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		respondError(w, http.StatusServiceUnavailable, "Reports unavailable: database connection not initialized")
		return
	}
	tenantID, err := middleware.VerifyTenantID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}
	restaurantID, err := middleware.VerifyRestaurantID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Missing restaurant context")
		return
	}

	q := r.URL.Query()
	dateFrom, dateToExclusive, dateToDisplay, err := parseReportDateRange(q.Get("date_from"), q.Get("date_to"))
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	topLimit := parseReportTopLimit(q.Get("top_limit"), 10)

	startedAt := time.Now()
	cacheHit := false
	var endpointErr error
	defer func() {
		h.logReportEndpointTiming("profitability", tenantID, restaurantID, startedAt, cacheHit, endpointErr)
	}()

	cacheKey := reportCacheKey(
		"reports:profitability",
		fmt.Sprintf("tenant=%d", tenantID),
		fmt.Sprintf("restaurant=%d", restaurantID),
		fmt.Sprintf("date_from=%s", dateFrom.Format("2006-01-02")),
		fmt.Sprintf("date_to=%s", dateToDisplay.Format("2006-01-02")),
		fmt.Sprintf("top_limit=%d", topLimit),
	)
	if payload, contentType, ok := h.cache.Get(cacheKey); ok {
		cacheHit = true
		if contentType == "" {
			contentType = reportCacheContentTypeJSON
		}
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(payload)
		return
	}

	summary, err := h.queryProfitabilitySummary(tenantID, restaurantID, dateFrom, dateToExclusive)
	if err != nil {
		endpointErr = err
		respondError(w, http.StatusInternalServerError, "Failed to load profitability summary")
		return
	}
	topItems, lowestItems, err := h.queryProfitabilityItems(tenantID, restaurantID, dateFrom, dateToExclusive, topLimit)
	if err != nil {
		endpointErr = err
		respondError(w, http.StatusInternalServerError, "Failed to load profitability items")
		return
	}

	if err := h.writeCachedOrStoreJSON(w, cacheKey, reportCacheTTLProfitability, map[string]interface{}{
		"success": true,
		"data": profitabilityReportResponse{
			Filters: map[string]interface{}{
				"date_from":  dateFrom.Format("2006-01-02"),
				"date_to":    dateToDisplay.Format("2006-01-02"),
				"top_limit":  topLimit,
				"cost_field": "products.cost",
			},
			Summary:     summary,
			TopItems:    topItems,
			LowestItems: lowestItems,
			GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		},
	}); err != nil {
		endpointErr = err
		respondError(w, http.StatusInternalServerError, "Failed to serialize profitability report")
		return
	}
}

// GetInventoryReport returns inventory turnover analytics for tracked products.
// GET /api/v1/reports/inventory?date_from=YYYY-MM-DD&date_to=YYYY-MM-DD&top_limit=10
func (h *ReportHandler) GetInventoryReport(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		respondError(w, http.StatusServiceUnavailable, "Reports unavailable: database connection not initialized")
		return
	}
	tenantID, err := middleware.VerifyTenantID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}
	restaurantID, err := middleware.VerifyRestaurantID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Missing restaurant context")
		return
	}

	q := r.URL.Query()
	dateFrom, dateToExclusive, dateToDisplay, err := parseReportDateRange(q.Get("date_from"), q.Get("date_to"))
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	topLimit := parseReportTopLimit(q.Get("top_limit"), 10)
	daysInRange := int64(dateToExclusive.Sub(dateFrom).Hours() / 24)
	if daysInRange <= 0 {
		daysInRange = 1
	}

	startedAt := time.Now()
	cacheHit := false
	var endpointErr error
	defer func() {
		h.logReportEndpointTiming("inventory", tenantID, restaurantID, startedAt, cacheHit, endpointErr)
	}()

	cacheKey := reportCacheKey(
		"reports:inventory",
		fmt.Sprintf("tenant=%d", tenantID),
		fmt.Sprintf("restaurant=%d", restaurantID),
		fmt.Sprintf("date_from=%s", dateFrom.Format("2006-01-02")),
		fmt.Sprintf("date_to=%s", dateToDisplay.Format("2006-01-02")),
		fmt.Sprintf("top_limit=%d", topLimit),
		fmt.Sprintf("days_in_range=%d", daysInRange),
	)
	if payload, contentType, ok := h.cache.Get(cacheKey); ok {
		cacheHit = true
		if contentType == "" {
			contentType = reportCacheContentTypeJSON
		}
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(payload)
		return
	}

	allItems, err := h.queryInventoryTurnoverItems(tenantID, restaurantID, dateFrom, dateToExclusive, daysInRange)
	if err != nil {
		endpointErr = err
		respondError(w, http.StatusInternalServerError, "Failed to load inventory analytics")
		return
	}
	summary := summarizeInventoryItems(allItems)
	slowMovers, topTurnover := splitInventoryItems(allItems, topLimit)

	if err := h.writeCachedOrStoreJSON(w, cacheKey, reportCacheTTLInventory, map[string]interface{}{
		"success": true,
		"data": inventoryReportResponse{
			Filters: map[string]interface{}{
				"date_from":     dateFrom.Format("2006-01-02"),
				"date_to":       dateToDisplay.Format("2006-01-02"),
				"top_limit":     topLimit,
				"days_in_range": daysInRange,
				"formula_note":  "Estimated COGS uses products.cost; turnover normalized to monthly ratio",
			},
			Summary:     summary,
			Items:       allItems,
			SlowMovers:  slowMovers,
			TopTurnover: topTurnover,
			GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		},
	}); err != nil {
		endpointErr = err
		respondError(w, http.StatusInternalServerError, "Failed to serialize inventory report")
		return
	}
}

// GetOverviewReport returns lightweight summaries for the reports overview dashboard.
// GET /api/v1/reports/overview?date_from=YYYY-MM-DD&date_to=YYYY-MM-DD
func (h *ReportHandler) GetOverviewReport(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		respondError(w, http.StatusServiceUnavailable, "Reports unavailable: database connection not initialized")
		return
	}
	tenantID, err := middleware.VerifyTenantID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}
	restaurantID, err := middleware.VerifyRestaurantID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Missing restaurant context")
		return
	}

	q := r.URL.Query()
	dateFrom, dateToExclusive, dateToDisplay, err := parseReportDateRange(q.Get("date_from"), q.Get("date_to"))
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	daysInRange := int64(dateToExclusive.Sub(dateFrom).Hours() / 24)
	if daysInRange <= 0 {
		daysInRange = 1
	}

	startedAt := time.Now()
	cacheHit := false
	var endpointErr error
	defer func() {
		h.logReportEndpointTiming("overview", tenantID, restaurantID, startedAt, cacheHit, endpointErr)
	}()

	cacheKey := reportCacheKey(
		"reports:overview",
		fmt.Sprintf("tenant=%d", tenantID),
		fmt.Sprintf("restaurant=%d", restaurantID),
		fmt.Sprintf("date_from=%s", dateFrom.Format("2006-01-02")),
		fmt.Sprintf("date_to=%s", dateToDisplay.Format("2006-01-02")),
		fmt.Sprintf("days_in_range=%d", daysInRange),
	)
	if payload, contentType, ok := h.cache.Get(cacheKey); ok {
		cacheHit = true
		if contentType == "" {
			contentType = reportCacheContentTypeJSON
		}
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(payload)
		return
	}

	var (
		salesSummary   salesSummaryResponse
		salesTop       []salesTopItem
		profitSummary  profitabilitySummaryResponse
		profitTop      []profitabilityItem
		inventoryItems []inventoryTurnoverItem
		firstErr       error
		errMu          sync.Mutex
	)

	setErr := func(err error) {
		if err == nil {
			return
		}
		errMu.Lock()
		if firstErr == nil {
			firstErr = err
		}
		errMu.Unlock()
	}

	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		defer wg.Done()
		out, qErr := h.querySalesSummary(tenantID, restaurantID, dateFrom, dateToExclusive)
		if qErr != nil {
			setErr(qErr)
			return
		}
		salesSummary = out
	}()

	go func() {
		defer wg.Done()
		out, qErr := h.queryTopItems(tenantID, restaurantID, dateFrom, dateToExclusive, 1)
		if qErr != nil {
			setErr(qErr)
			return
		}
		salesTop = out
	}()

	go func() {
		defer wg.Done()
		out, qErr := h.queryProfitabilitySummary(tenantID, restaurantID, dateFrom, dateToExclusive)
		if qErr != nil {
			setErr(qErr)
			return
		}
		profitSummary = out
	}()

	go func() {
		defer wg.Done()
		top, _, qErr := h.queryProfitabilityItems(tenantID, restaurantID, dateFrom, dateToExclusive, 1)
		if qErr != nil {
			setErr(qErr)
			return
		}
		profitTop = top
	}()

	go func() {
		defer wg.Done()
		out, qErr := h.queryInventoryTurnoverItems(tenantID, restaurantID, dateFrom, dateToExclusive, daysInRange)
		if qErr != nil {
			setErr(qErr)
			return
		}
		inventoryItems = out
	}()

	wg.Wait()
	if firstErr != nil {
		endpointErr = firstErr
		respondError(w, http.StatusInternalServerError, "Failed to load reports overview")
		return
	}
	inventorySummary := summarizeInventoryItems(inventoryItems)
	slowMovers, _ := splitInventoryItems(inventoryItems, 1)

	resp := reportsOverviewResponse{
		Filters: map[string]interface{}{
			"date_from": dateFrom.Format("2006-01-02"),
			"date_to":   dateToDisplay.Format("2006-01-02"),
		},
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
	}
	resp.Sales.Summary = salesSummary
	if len(salesTop) > 0 {
		resp.Sales.TopItem = &salesTop[0]
	}
	resp.Profitability.Summary = profitSummary
	if len(profitTop) > 0 {
		resp.Profitability.MostProfitable = &profitTop[0]
	}
	resp.Inventory.Summary = inventorySummary
	if len(slowMovers) > 0 {
		resp.Inventory.SlowestMoving = &slowMovers[0]
	}

	if err := h.writeCachedOrStoreJSON(w, cacheKey, reportCacheTTLOverview, map[string]interface{}{
		"success": true,
		"data":    resp,
	}); err != nil {
		endpointErr = err
		respondError(w, http.StatusInternalServerError, "Failed to serialize reports overview")
		return
	}
}

// ExportReport exports reports data as CSV or PDF summary.
// POST /api/v1/reports/export
func (h *ReportHandler) ExportReport(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		respondError(w, http.StatusServiceUnavailable, "Reports unavailable: database connection not initialized")
		return
	}

	tenantID, err := middleware.VerifyTenantID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}
	restaurantID, err := middleware.VerifyRestaurantID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Missing restaurant context")
		return
	}

	req, err := parseReportExportRequest(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	reportType := strings.ToLower(strings.TrimSpace(req.ReportType))
	format := strings.ToLower(strings.TrimSpace(req.Format))
	if reportType == "" {
		reportType = "sales"
	}
	if format == "" {
		format = "csv"
	}
	if reportType != "sales" && reportType != "profitability" && reportType != "inventory" {
		respondError(w, http.StatusBadRequest, "Unsupported report_type (expected sales, profitability, or inventory)")
		return
	}
	if format != "csv" && format != "pdf" {
		respondError(w, http.StatusBadRequest, "Unsupported format (expected csv or pdf)")
		return
	}

	dateFrom, dateToExclusive, dateToDisplay, err := parseReportDateRange(req.DateFrom, req.DateTo)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	granularity := normalizeGranularity(req.Granularity)
	exportLocale := normalizeReportLocale(req.Locale)
	topLimit := parseReportTopLimit(fmt.Sprintf("%d", req.TopLimit), 10)
	daysInRange := int64(dateToExclusive.Sub(dateFrom).Hours() / 24)
	if daysInRange <= 0 {
		daysInRange = 1
	}

	var (
		payload     []byte
		contentType string
		filename    string
	)

	switch reportType {
	case "sales":
		report, err := h.buildSalesReportDataForExport(tenantID, restaurantID, dateFrom, dateToExclusive, dateToDisplay, granularity, topLimit)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to generate sales report export")
			return
		}
		if format == "csv" {
			payload, err = buildSalesCSV(report, exportLocale)
			contentType = "text/csv; charset=utf-8"
			filename = fmt.Sprintf("sales-report-%s-to-%s.csv", req.DateFromOr(dateFrom), req.DateToOr(dateToDisplay))
		} else {
			payload, err = buildSalesPDFSummary(report, exportLocale)
			contentType = "application/pdf"
			filename = fmt.Sprintf("sales-report-summary-%s-to-%s.pdf", req.DateFromOr(dateFrom), req.DateToOr(dateToDisplay))
		}
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to export sales report")
			return
		}

	case "profitability":
		report, err := h.buildProfitabilityReportDataForExport(tenantID, restaurantID, dateFrom, dateToExclusive, dateToDisplay, topLimit)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to generate profitability report export")
			return
		}
		if format == "csv" {
			payload, err = buildProfitabilityCSV(report, exportLocale)
			contentType = "text/csv; charset=utf-8"
			filename = fmt.Sprintf("profitability-report-%s-to-%s.csv", req.DateFromOr(dateFrom), req.DateToOr(dateToDisplay))
		} else {
			payload, err = buildProfitabilityPDFSummary(report, exportLocale)
			contentType = "application/pdf"
			filename = fmt.Sprintf("profitability-report-summary-%s-to-%s.pdf", req.DateFromOr(dateFrom), req.DateToOr(dateToDisplay))
		}
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to export profitability report")
			return
		}

	case "inventory":
		report, err := h.buildInventoryReportDataForExport(tenantID, restaurantID, dateFrom, dateToExclusive, dateToDisplay, topLimit, daysInRange)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to generate inventory report export")
			return
		}
		if format == "csv" {
			payload, err = buildInventoryCSV(report, exportLocale)
			contentType = "text/csv; charset=utf-8"
			filename = fmt.Sprintf("inventory-report-%s-to-%s.csv", req.DateFromOr(dateFrom), req.DateToOr(dateToDisplay))
		} else {
			payload, err = buildInventoryPDFSummary(report, exportLocale)
			contentType = "application/pdf"
			filename = fmt.Sprintf("inventory-report-summary-%s-to-%s.pdf", req.DateFromOr(dateFrom), req.DateToOr(dateToDisplay))
		}
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to export inventory report")
			return
		}
	}

	auditFilters := map[string]interface{}{
		"date_from":     req.DateFromOr(dateFrom),
		"date_to":       req.DateToOr(dateToDisplay),
		"granularity":   granularity,
		"top_limit":     topLimit,
		"days_in_range": daysInRange,
		"locale":        exportLocale,
	}
	if err := h.logReportExportAudit(r.Context(), reportExportAuditRecord{
		TenantID:      tenantID,
		RestaurantID:  restaurantID,
		UserID:        middleware.GetUserID(r),
		ReportType:    reportType,
		Format:        format,
		FileName:      filename,
		ExportedBytes: len(payload),
		Filters:       auditFilters,
		RequestPath:   r.URL.Path,
		UserAgent:     r.UserAgent(),
		IPAddress:     clientIPFromRequest(r),
	}); err != nil {
		log.Printf("[ReportHandler] WARN: failed to audit report export: %v", err)
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Language", exportLocale)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(payload); err != nil {
		// headers already sent; nothing else to do
		return
	}
}

func parseReportExportRequest(r *http.Request) (reportExportRequest, error) {
	req := reportExportRequest{}
	if r.Body != nil {
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != io.EOF {
			return req, fmt.Errorf("invalid request body")
		}
	}

	q := r.URL.Query()
	if strings.TrimSpace(req.ReportType) == "" {
		req.ReportType = q.Get("report_type")
	}
	if strings.TrimSpace(req.Format) == "" {
		req.Format = q.Get("format")
	}
	if strings.TrimSpace(req.DateFrom) == "" {
		req.DateFrom = q.Get("date_from")
	}
	if strings.TrimSpace(req.DateTo) == "" {
		req.DateTo = q.Get("date_to")
	}
	if strings.TrimSpace(req.Granularity) == "" {
		req.Granularity = q.Get("granularity")
	}
	if req.TopLimit == 0 {
		req.TopLimit = parseInt64(q.Get("top_limit"), 10)
	}
	if strings.TrimSpace(req.Locale) == "" {
		req.Locale = q.Get("locale")
	}
	if strings.TrimSpace(req.Locale) == "" {
		req.Locale = r.Header.Get("X-Restaurant-Locale")
	}
	if strings.TrimSpace(req.Locale) == "" {
		req.Locale = r.Header.Get("Accept-Language")
	}
	return req, nil
}

func (r reportExportRequest) DateFromOr(fallback time.Time) string {
	if strings.TrimSpace(r.DateFrom) != "" {
		return r.DateFrom
	}
	return fallback.Format("2006-01-02")
}

func (r reportExportRequest) DateToOr(fallback time.Time) string {
	if strings.TrimSpace(r.DateTo) != "" {
		return r.DateTo
	}
	return fallback.Format("2006-01-02")
}

func (h *ReportHandler) buildSalesReportDataForExport(
	tenantID, restaurantID int64,
	dateFrom, dateToExclusive, dateToDisplay time.Time,
	granularity string,
	topLimit int64,
) (salesReportResponse, error) {
	if h.salesReportDataBuilder != nil {
		return h.salesReportDataBuilder(tenantID, restaurantID, dateFrom, dateToExclusive, dateToDisplay, granularity, topLimit)
	}
	return h.buildSalesReportData(tenantID, restaurantID, dateFrom, dateToExclusive, dateToDisplay, granularity, topLimit)
}

func (h *ReportHandler) buildProfitabilityReportDataForExport(
	tenantID, restaurantID int64,
	dateFrom, dateToExclusive, dateToDisplay time.Time,
	topLimit int64,
) (profitabilityReportResponse, error) {
	if h.profitabilityReportDataBuilder != nil {
		return h.profitabilityReportDataBuilder(tenantID, restaurantID, dateFrom, dateToExclusive, dateToDisplay, topLimit)
	}
	return h.buildProfitabilityReportData(tenantID, restaurantID, dateFrom, dateToExclusive, dateToDisplay, topLimit)
}

func (h *ReportHandler) buildInventoryReportDataForExport(
	tenantID, restaurantID int64,
	dateFrom, dateToExclusive, dateToDisplay time.Time,
	topLimit int64,
	daysInRange int64,
) (inventoryReportResponse, error) {
	if h.inventoryReportDataBuilder != nil {
		return h.inventoryReportDataBuilder(tenantID, restaurantID, dateFrom, dateToExclusive, dateToDisplay, topLimit, daysInRange)
	}
	return h.buildInventoryReportData(tenantID, restaurantID, dateFrom, dateToExclusive, dateToDisplay, topLimit, daysInRange)
}

func (h *ReportHandler) logReportExportAudit(ctx context.Context, record reportExportAuditRecord) error {
	if h.reportExportAuditLogger != nil {
		return h.reportExportAuditLogger(ctx, record)
	}
	if h.db == nil {
		return nil
	}

	filtersJSON, err := json.Marshal(record.Filters)
	if err != nil {
		return fmt.Errorf("marshal export filters: %w", err)
	}

	execCtx, cancel := withReportDBTimeoutFrom(ctx)
	defer cancel()
	_, err = h.db.ExecContext(execCtx, `
		INSERT INTO report_export_audit_logs (
			tenant_id, restaurant_id, user_id, report_type, format, file_name,
			exported_bytes, filters, request_path, user_agent, ip_address, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW())
	`,
		record.TenantID,
		record.RestaurantID,
		record.UserID,
		record.ReportType,
		record.Format,
		record.FileName,
		record.ExportedBytes,
		string(filtersJSON),
		record.RequestPath,
		record.UserAgent,
		record.IPAddress,
	)
	if err != nil {
		return fmt.Errorf("insert report export audit: %w", err)
	}
	return nil
}

func clientIPFromRequest(r *http.Request) string {
	for _, header := range []string{"X-Forwarded-For", "X-Real-IP"} {
		if raw := strings.TrimSpace(r.Header.Get(header)); raw != "" {
			if header == "X-Forwarded-For" && strings.Contains(raw, ",") {
				raw = strings.TrimSpace(strings.Split(raw, ",")[0])
			}
			if addr, err := netip.ParseAddr(raw); err == nil {
				return addr.String()
			}
			return raw
		}
	}

	hostPort := strings.TrimSpace(r.RemoteAddr)
	if hostPort == "" {
		return ""
	}
	if addr, err := netip.ParseAddrPort(hostPort); err == nil {
		return addr.Addr().String()
	}
	return hostPort
}

func (h *ReportHandler) buildSalesReportData(
	tenantID, restaurantID int64,
	dateFrom, dateToExclusive, dateToDisplay time.Time,
	granularity string,
	topLimit int64,
) (salesReportResponse, error) {
	summary, err := h.querySalesSummary(tenantID, restaurantID, dateFrom, dateToExclusive)
	if err != nil {
		return salesReportResponse{}, err
	}
	trend, err := h.querySalesTrend(tenantID, restaurantID, dateFrom, dateToExclusive, granularity)
	if err != nil {
		return salesReportResponse{}, err
	}
	hourly, err := h.querySalesHourly(tenantID, restaurantID, dateFrom, dateToExclusive)
	if err != nil {
		return salesReportResponse{}, err
	}
	topItems, err := h.queryTopItems(tenantID, restaurantID, dateFrom, dateToExclusive, topLimit)
	if err != nil {
		return salesReportResponse{}, err
	}
	topCategories, err := h.queryTopCategories(tenantID, restaurantID, dateFrom, dateToExclusive, topLimit)
	if err != nil {
		return salesReportResponse{}, err
	}

	return salesReportResponse{
		Filters: map[string]interface{}{
			"date_from":         dateFrom.Format("2006-01-02"),
			"date_to":           dateToDisplay.Format("2006-01-02"),
			"granularity":       granularity,
			"top_limit":         topLimit,
			"included_statuses": []string{"all_except_cancelled"},
		},
		Summary:       summary,
		Trend:         trend,
		Hourly:        hourly,
		TopItems:      topItems,
		TopCategories: topCategories,
		GeneratedAt:   time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func (h *ReportHandler) buildProfitabilityReportData(
	tenantID, restaurantID int64,
	dateFrom, dateToExclusive, dateToDisplay time.Time,
	topLimit int64,
) (profitabilityReportResponse, error) {
	summary, err := h.queryProfitabilitySummary(tenantID, restaurantID, dateFrom, dateToExclusive)
	if err != nil {
		return profitabilityReportResponse{}, err
	}
	topItems, lowestItems, err := h.queryProfitabilityItems(tenantID, restaurantID, dateFrom, dateToExclusive, topLimit)
	if err != nil {
		return profitabilityReportResponse{}, err
	}

	return profitabilityReportResponse{
		Filters: map[string]interface{}{
			"date_from":  dateFrom.Format("2006-01-02"),
			"date_to":    dateToDisplay.Format("2006-01-02"),
			"top_limit":  topLimit,
			"cost_field": "products.cost",
		},
		Summary:     summary,
		TopItems:    topItems,
		LowestItems: lowestItems,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func (h *ReportHandler) buildInventoryReportData(
	tenantID, restaurantID int64,
	dateFrom, dateToExclusive, dateToDisplay time.Time,
	topLimit int64,
	daysInRange int64,
) (inventoryReportResponse, error) {
	allItems, err := h.queryInventoryTurnoverItems(tenantID, restaurantID, dateFrom, dateToExclusive, daysInRange)
	if err != nil {
		return inventoryReportResponse{}, err
	}
	summary := summarizeInventoryItems(allItems)
	slowMovers, topTurnover := splitInventoryItems(allItems, topLimit)

	return inventoryReportResponse{
		Filters: map[string]interface{}{
			"date_from":     dateFrom.Format("2006-01-02"),
			"date_to":       dateToDisplay.Format("2006-01-02"),
			"top_limit":     topLimit,
			"days_in_range": daysInRange,
			"formula_note":  "Estimated COGS uses products.cost; turnover normalized to monthly ratio",
		},
		Summary:     summary,
		Items:       allItems,
		SlowMovers:  slowMovers,
		TopTurnover: topTurnover,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

type reportExportTextSet struct {
	CSVSection              string
	CSVLabel1               string
	CSVLabel2               string
	CSVRevenue              string
	CSVOrders               string
	CSVQuantity             string
	CSVProduct              string
	CSVCategory             string
	CSVCOGS                 string
	CSVGrossProfit          string
	CSVGrossMarginPct       string
	CSVFoodCostPct          string
	CSVMissingCostData      string
	CSVCurrentStockQty      string
	CSVSoldQty              string
	CSVWasteQty             string
	CSVAvgInventoryValue    string
	CSVEstimatedCOGS        string
	CSVMonthlyTurnoverRatio string
	CSVIsSlowMoving         string
	CSVHour                 string
	CSVTotal                string

	SectionSummary           string
	SectionTrend             string
	SectionHourly            string
	SectionTopItems          string
	SectionTopCategories     string
	SectionLowestMarginItems string
	SectionItems             string

	LabelTotalRevenue string
	LabelTotalOrders  string
	LabelAvgTicket    string

	PDFSalesTitle         string
	PDFTopItems           string
	PDFTopCategories      string
	PDFPeakHours          string
	PDFProfitabilityTitle string
	PDFMissingCostNotice  string
	PDFTopProfitableItems string
	PDFLowestMarginItems  string
	PDFInventoryTitle     string
	PDFTrackedProducts    string
	PDFInventoryValue     string
	PDFAvgMonthlyTurnover string
	PDFSlowMovers         string
	PDFTopTurnover        string
	PDFSlowMovingItems    string
}

func normalizeReportLocale(value string) string {
	v := strings.ToLower(strings.TrimSpace(value))
	if v == "" {
		return "en"
	}
	if strings.HasPrefix(v, "ar") {
		return "ar"
	}
	return "en"
}

func reportExportTexts(locale string) reportExportTextSet {
	// CSV and PDF labels are locale-aware (Arabic and English) for report exports.
	if normalizeReportLocale(locale) == "ar" {
		return reportExportTextSet{
			CSVSection:              "القسم",
			CSVLabel1:               "التسمية_1",
			CSVLabel2:               "التسمية_2",
			CSVRevenue:              "الإيراد",
			CSVOrders:               "الطلبات",
			CSVQuantity:             "الكمية",
			CSVProduct:              "المنتج",
			CSVCategory:             "الفئة",
			CSVCOGS:                 "تكلفة_البضاعة_المباعة",
			CSVGrossProfit:          "الربح_الإجمالي",
			CSVGrossMarginPct:       "نسبة_الهامش_الإجمالي",
			CSVFoodCostPct:          "نسبة_تكلفة_الطعام",
			CSVMissingCostData:      "بيانات_تكلفة_مفقودة",
			CSVCurrentStockQty:      "كمية_المخزون_الحالي",
			CSVSoldQty:              "الكمية_المباعة",
			CSVWasteQty:             "كمية_الهدر",
			CSVAvgInventoryValue:    "متوسط_قيمة_المخزون",
			CSVEstimatedCOGS:        "تكلفة_تقديرية_للبضاعة_المباعة",
			CSVMonthlyTurnoverRatio: "معدل_الدوران_الشهري",
			CSVIsSlowMoving:         "بطيء_الحركة",
			CSVHour:                 "الساعة",
			CSVTotal:                "الإجمالي",

			SectionSummary:           "ملخص",
			SectionTrend:             "الاتجاه",
			SectionHourly:            "بالساعة",
			SectionTopItems:          "أفضل_المنتجات",
			SectionTopCategories:     "أفضل_الفئات",
			SectionLowestMarginItems: "أقل_هامش_ربح",
			SectionItems:             "العناصر",

			LabelTotalRevenue: "إجمالي الإيراد",
			LabelTotalOrders:  "إجمالي الطلبات",
			LabelAvgTicket:    "متوسط قيمة الطلب",

			PDFSalesTitle:         "ملخص تحليلات المبيعات",
			PDFTopItems:           "أفضل المنتجات",
			PDFTopCategories:      "أفضل الفئات",
			PDFPeakHours:          "ساعات الذروة (أعلى إيراد)",
			PDFProfitabilityTitle: "ملخص الربحية",
			PDFMissingCostNotice:  "تحذير: وحدات بدون بيانات تكلفة (products.cost)",
			PDFTopProfitableItems: "أعلى المنتجات ربحية",
			PDFLowestMarginItems:  "أقل المنتجات هامشًا",
			PDFInventoryTitle:     "ملخص دوران المخزون",
			PDFTrackedProducts:    "المنتجات المتتبعة",
			PDFInventoryValue:     "قيمة المخزون",
			PDFAvgMonthlyTurnover: "متوسط الدوران الشهري",
			PDFSlowMovers:         "بطيئة الحركة",
			PDFTopTurnover:        "أعلى دوران",
			PDFSlowMovingItems:    "المنتجات بطيئة الحركة",
		}
	}

	return reportExportTextSet{
		CSVSection:              "section",
		CSVLabel1:               "label_1",
		CSVLabel2:               "label_2",
		CSVRevenue:              "revenue",
		CSVOrders:               "orders",
		CSVQuantity:             "quantity",
		CSVProduct:              "product",
		CSVCategory:             "category",
		CSVCOGS:                 "cogs",
		CSVGrossProfit:          "gross_profit",
		CSVGrossMarginPct:       "gross_margin_pct",
		CSVFoodCostPct:          "food_cost_pct",
		CSVMissingCostData:      "missing_cost_data",
		CSVCurrentStockQty:      "current_stock_qty",
		CSVSoldQty:              "sold_qty",
		CSVWasteQty:             "waste_qty",
		CSVAvgInventoryValue:    "avg_inventory_value",
		CSVEstimatedCOGS:        "estimated_cogs",
		CSVMonthlyTurnoverRatio: "monthly_turnover_ratio",
		CSVIsSlowMoving:         "is_slow_moving",
		CSVHour:                 "hour",
		CSVTotal:                "TOTAL",

		SectionSummary:           "summary",
		SectionTrend:             "trend",
		SectionHourly:            "hourly",
		SectionTopItems:          "top_items",
		SectionTopCategories:     "top_categories",
		SectionLowestMarginItems: "lowest_margin_items",
		SectionItems:             "items",

		LabelTotalRevenue: "Total Revenue",
		LabelTotalOrders:  "Total Orders",
		LabelAvgTicket:    "Average Ticket",

		PDFSalesTitle:         "Sales Analytics Summary",
		PDFTopItems:           "Top Items",
		PDFTopCategories:      "Top Categories",
		PDFPeakHours:          "Peak Hours (Top Revenue)",
		PDFProfitabilityTitle: "Profitability Summary",
		PDFMissingCostNotice:  "Warning: units missing cost data (products.cost)",
		PDFTopProfitableItems: "Top Profitable Items",
		PDFLowestMarginItems:  "Lowest Margin Items",
		PDFInventoryTitle:     "Inventory Turnover Summary",
		PDFTrackedProducts:    "Tracked Products",
		PDFInventoryValue:     "Inventory Value",
		PDFAvgMonthlyTurnover: "Avg Monthly Turnover",
		PDFSlowMovers:         "Slow Movers",
		PDFTopTurnover:        "Top Turnover",
		PDFSlowMovingItems:    "Slow Moving Items",
	}
}

func reportExportBool(value bool, locale string) string {
	if normalizeReportLocale(locale) == "ar" {
		if value {
			return "نعم"
		}
		return "لا"
	}
	if value {
		return "true"
	}
	return "false"
}

func formatReportCurrency(amount float64, locale string) string {
	if normalizeReportLocale(locale) == "ar" {
		return fmt.Sprintf("%.2f ج.م", amount)
	}
	return fmt.Sprintf("EGP %.2f", amount)
}

func writeArabicPDFFallbackNoticeIfNeeded(pdf *fpdf.Fpdf, locale string) {
	_ = pdf
	_ = locale
}

func buildSalesCSV(report salesReportResponse, locale string) ([]byte, error) {
	t := reportExportTexts(locale)
	buf := &bytes.Buffer{}
	buf.Write([]byte{0xEF, 0xBB, 0xBF}) // UTF-8 BOM for Excel compatibility
	w := csv.NewWriter(buf)
	w.UseCRLF = true

	rows := [][]string{
		{t.CSVSection, t.CSVLabel1, t.CSVLabel2, t.CSVRevenue, t.CSVOrders, t.CSVQuantity},
		{t.SectionSummary, t.LabelTotalRevenue, "", fmt.Sprintf("%.2f", report.Summary.TotalRevenue), fmt.Sprintf("%d", report.Summary.TotalOrders), ""},
		{t.SectionSummary, t.LabelAvgTicket, "", fmt.Sprintf("%.2f", report.Summary.AvgTicket), "", ""},
	}
	for _, p := range report.Trend {
		rows = append(rows, []string{t.SectionTrend, p.Bucket, "", fmt.Sprintf("%.2f", p.Revenue), fmt.Sprintf("%d", p.Orders), ""})
	}
	for _, p := range report.Hourly {
		rows = append(rows, []string{t.SectionHourly, fmt.Sprintf("%d", p.Hour), "", fmt.Sprintf("%.2f", p.Revenue), fmt.Sprintf("%d", p.Orders), ""})
	}
	for _, item := range report.TopItems {
		rows = append(rows, []string{t.SectionTopItems, item.Name, "", fmt.Sprintf("%.2f", item.Revenue), "", fmt.Sprintf("%d", item.Quantity)})
	}
	for _, item := range report.TopCategories {
		rows = append(rows, []string{t.SectionTopCategories, item.Name, "", fmt.Sprintf("%.2f", item.Revenue), "", fmt.Sprintf("%d", item.Quantity)})
	}

	if err := w.WriteAll(rows); err != nil {
		return nil, err
	}
	return buf.Bytes(), w.Error()
}

func buildProfitabilityCSV(report profitabilityReportResponse, locale string) ([]byte, error) {
	t := reportExportTexts(locale)
	buf := &bytes.Buffer{}
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(buf)
	w.UseCRLF = true
	rows := [][]string{
		{t.CSVSection, t.CSVProduct, t.CSVQuantity, t.CSVRevenue, t.CSVCOGS, t.CSVGrossProfit, t.CSVGrossMarginPct, t.CSVFoodCostPct, t.CSVMissingCostData},
		{t.SectionSummary, t.CSVTotal, "", fmt.Sprintf("%.2f", report.Summary.TotalRevenue), fmt.Sprintf("%.2f", report.Summary.TotalCOGS), fmt.Sprintf("%.2f", report.Summary.GrossProfit), fmt.Sprintf("%.2f", report.Summary.GrossMarginPct), fmt.Sprintf("%.2f", report.Summary.FoodCostPct), ""},
	}
	for _, item := range report.TopItems {
		rows = append(rows, []string{
			t.SectionTopItems,
			item.Name,
			fmt.Sprintf("%d", item.Quantity),
			fmt.Sprintf("%.2f", item.Revenue),
			fmt.Sprintf("%.2f", item.COGS),
			fmt.Sprintf("%.2f", item.GrossProfit),
			fmt.Sprintf("%.2f", item.GrossMarginPct),
			fmt.Sprintf("%.2f", item.FoodCostPct),
			reportExportBool(item.MissingCostData, locale),
		})
	}
	for _, item := range report.LowestItems {
		rows = append(rows, []string{
			t.SectionLowestMarginItems,
			item.Name,
			fmt.Sprintf("%d", item.Quantity),
			fmt.Sprintf("%.2f", item.Revenue),
			fmt.Sprintf("%.2f", item.COGS),
			fmt.Sprintf("%.2f", item.GrossProfit),
			fmt.Sprintf("%.2f", item.GrossMarginPct),
			fmt.Sprintf("%.2f", item.FoodCostPct),
			reportExportBool(item.MissingCostData, locale),
		})
	}
	if err := w.WriteAll(rows); err != nil {
		return nil, err
	}
	return buf.Bytes(), w.Error()
}

func buildInventoryCSV(report inventoryReportResponse, locale string) ([]byte, error) {
	t := reportExportTexts(locale)
	buf := &bytes.Buffer{}
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(buf)
	w.UseCRLF = true
	rows := [][]string{
		{t.CSVSection, t.CSVProduct, t.CSVCurrentStockQty, t.CSVSoldQty, t.CSVWasteQty, t.CSVAvgInventoryValue, t.CSVEstimatedCOGS, t.CSVMonthlyTurnoverRatio, t.CSVIsSlowMoving},
		{t.SectionSummary, t.CSVTotal, "", "", "", fmt.Sprintf("%.2f", report.Summary.TotalCurrentInventoryValue), fmt.Sprintf("%.2f", report.Summary.TotalEstimatedCOGS), fmt.Sprintf("%.2f", report.Summary.AverageMonthlyTurnover), ""},
	}
	for _, item := range report.Items {
		rows = append(rows, []string{
			t.SectionItems,
			item.Name,
			fmt.Sprintf("%.2f", item.CurrentStockQty),
			fmt.Sprintf("%d", item.SoldQty),
			fmt.Sprintf("%d", item.WasteQty),
			fmt.Sprintf("%.2f", item.AverageInventoryValue),
			fmt.Sprintf("%.2f", item.EstimatedCOGS),
			fmt.Sprintf("%.2f", item.MonthlyTurnoverRatio),
			reportExportBool(item.IsSlowMoving, locale),
		})
	}
	if err := w.WriteAll(rows); err != nil {
		return nil, err
	}
	return buf.Bytes(), w.Error()
}

func buildSalesPDFSummary(report salesReportResponse, locale string) ([]byte, error) {
	t := reportExportTexts(locale)
	pdf := newReportPDF(t.PDFSalesTitle, report.Filters, report.GeneratedAt, locale)
	pdf.AddPage()
	writePDFMetricCards(pdf, []pdfMetricCard{
		{Label: t.LabelTotalRevenue, Value: formatReportCurrency(report.Summary.TotalRevenue, locale)},
		{Label: t.LabelTotalOrders, Value: fmt.Sprintf("%d", report.Summary.TotalOrders)},
		{Label: t.LabelAvgTicket, Value: formatReportCurrency(report.Summary.AvgTicket, locale)},
	})

	topItemRows := make([][]string, 0, minInt(len(report.TopItems), 10))
	for i, item := range report.TopItems {
		if i >= 10 {
			break
		}
		topItemRows = append(topItemRows, []string{
			fmt.Sprintf("%d", i+1),
			item.Name,
			fmt.Sprintf("%d", item.Quantity),
			formatReportCurrency(item.Revenue, locale),
		})
	}
	writePDFTable(pdf, t.PDFTopItems, []pdfTableColumn{
		{Header: "#", Width: 12, Align: "C"},
		{Header: pdfHumanLabel(t.CSVProduct, locale), Width: 140, Align: "L"},
		{Header: pdfHumanLabel(t.CSVQuantity, locale), Width: 35, Align: "R"},
		{Header: pdfHumanLabel(t.CSVRevenue, locale), Width: 55, Align: "R"},
	}, topItemRows, locale)

	topCategoryRows := make([][]string, 0, minInt(len(report.TopCategories), 10))
	for i, item := range report.TopCategories {
		if i >= 10 {
			break
		}
		topCategoryRows = append(topCategoryRows, []string{
			fmt.Sprintf("%d", i+1),
			item.Name,
			fmt.Sprintf("%d", item.Quantity),
			formatReportCurrency(item.Revenue, locale),
		})
	}
	writePDFTable(pdf, t.PDFTopCategories, []pdfTableColumn{
		{Header: "#", Width: 12, Align: "C"},
		{Header: pdfHumanLabel(t.CSVCategory, locale), Width: 140, Align: "L"},
		{Header: pdfHumanLabel(t.CSVQuantity, locale), Width: 35, Align: "R"},
		{Header: pdfHumanLabel(t.CSVRevenue, locale), Width: 55, Align: "R"},
	}, topCategoryRows, locale)

	if len(report.Hourly) > 0 {
		hourly := append([]salesHourlyPoint(nil), report.Hourly...)
		sort.Slice(hourly, func(i, j int) bool {
			if hourly[i].Revenue == hourly[j].Revenue {
				return hourly[i].Hour < hourly[j].Hour
			}
			return hourly[i].Revenue > hourly[j].Revenue
		})
		rows := make([][]string, 0, minInt(len(hourly), 8))
		for i, pt := range hourly {
			if i >= 8 {
				break
			}
			rows = append(rows, []string{
				fmt.Sprintf("%d", i+1),
				fmt.Sprintf("%02d:00", pt.Hour),
				fmt.Sprintf("%d", pt.Orders),
				formatReportCurrency(pt.Revenue, locale),
			})
		}
		writePDFTable(pdf, t.PDFPeakHours, []pdfTableColumn{
			{Header: "#", Width: 12, Align: "C"},
			{Header: pdfHumanLabel(t.CSVHour, locale), Width: 70, Align: "L"},
			{Header: pdfHumanLabel(t.CSVOrders, locale), Width: 55, Align: "R"},
			{Header: pdfHumanLabel(t.CSVRevenue, locale), Width: 65, Align: "R"},
		}, rows, locale)
	}
	return outputPDF(pdf)
}

func buildProfitabilityPDFSummary(report profitabilityReportResponse, locale string) ([]byte, error) {
	t := reportExportTexts(locale)
	pdf := newReportPDF(t.PDFProfitabilityTitle, report.Filters, report.GeneratedAt, locale)
	pdf.AddPage()
	writePDFMetricCards(pdf, []pdfMetricCard{
		{Label: t.LabelTotalRevenue, Value: formatReportCurrency(report.Summary.TotalRevenue, locale)},
		{Label: pdfHumanLabel(t.CSVCOGS, locale), Value: formatReportCurrency(report.Summary.TotalCOGS, locale)},
		{Label: pdfHumanLabel(t.CSVGrossProfit, locale), Value: formatReportCurrency(report.Summary.GrossProfit, locale)},
		{Label: pdfHumanLabel(t.CSVGrossMarginPct, locale), Value: fmt.Sprintf("%.2f%%", report.Summary.GrossMarginPct)},
		{Label: pdfHumanLabel(t.CSVFoodCostPct, locale), Value: fmt.Sprintf("%.2f%%", report.Summary.FoodCostPct)},
	})
	if report.Summary.UnitsMissingCost > 0 {
		writePDFNotice(pdf, fmt.Sprintf("%s: %d", t.PDFMissingCostNotice, report.Summary.UnitsMissingCost))
	}

	topRows := make([][]string, 0, minInt(len(report.TopItems), 10))
	for i, item := range report.TopItems {
		if i >= 10 {
			break
		}
		topRows = append(topRows, []string{
			fmt.Sprintf("%d", i+1),
			item.Name,
			fmt.Sprintf("%d", item.Quantity),
			formatReportCurrency(item.Revenue, locale),
			formatReportCurrency(item.COGS, locale),
			formatReportCurrency(item.GrossProfit, locale),
			fmt.Sprintf("%.2f%%", item.GrossMarginPct),
		})
	}
	writePDFTable(pdf, t.PDFTopProfitableItems, []pdfTableColumn{
		{Header: "#", Width: 10, Align: "C"},
		{Header: pdfHumanLabel(t.CSVProduct, locale), Width: 95, Align: "L"},
		{Header: pdfHumanLabel(t.CSVQuantity, locale), Width: 18, Align: "R"},
		{Header: pdfHumanLabel(t.CSVRevenue, locale), Width: 34, Align: "R"},
		{Header: pdfHumanLabel(t.CSVCOGS, locale), Width: 34, Align: "R"},
		{Header: pdfHumanLabel(t.CSVGrossProfit, locale), Width: 38, Align: "R"},
		{Header: pdfHumanLabel(t.CSVGrossMarginPct, locale), Width: 28, Align: "R"},
	}, topRows, locale)

	lowRows := make([][]string, 0, minInt(len(report.LowestItems), 10))
	for i, item := range report.LowestItems {
		if i >= 10 {
			break
		}
		lowRows = append(lowRows, []string{
			fmt.Sprintf("%d", i+1),
			item.Name,
			fmt.Sprintf("%d", item.Quantity),
			formatReportCurrency(item.Revenue, locale),
			formatReportCurrency(item.COGS, locale),
			formatReportCurrency(item.GrossProfit, locale),
			fmt.Sprintf("%.2f%%", item.GrossMarginPct),
		})
	}
	writePDFTable(pdf, t.PDFLowestMarginItems, []pdfTableColumn{
		{Header: "#", Width: 10, Align: "C"},
		{Header: pdfHumanLabel(t.CSVProduct, locale), Width: 95, Align: "L"},
		{Header: pdfHumanLabel(t.CSVQuantity, locale), Width: 18, Align: "R"},
		{Header: pdfHumanLabel(t.CSVRevenue, locale), Width: 34, Align: "R"},
		{Header: pdfHumanLabel(t.CSVCOGS, locale), Width: 34, Align: "R"},
		{Header: pdfHumanLabel(t.CSVGrossProfit, locale), Width: 38, Align: "R"},
		{Header: pdfHumanLabel(t.CSVGrossMarginPct, locale), Width: 28, Align: "R"},
	}, lowRows, locale)
	return outputPDF(pdf)
}

func buildInventoryPDFSummary(report inventoryReportResponse, locale string) ([]byte, error) {
	t := reportExportTexts(locale)
	pdf := newReportPDF(t.PDFInventoryTitle, report.Filters, report.GeneratedAt, locale)
	pdf.AddPage()
	writePDFMetricCards(pdf, []pdfMetricCard{
		{Label: t.PDFTrackedProducts, Value: fmt.Sprintf("%d", report.Summary.TrackedProductsCount)},
		{Label: t.PDFInventoryValue, Value: formatReportCurrency(report.Summary.TotalCurrentInventoryValue, locale)},
		{Label: pdfHumanLabel(t.CSVEstimatedCOGS, locale), Value: formatReportCurrency(report.Summary.TotalEstimatedCOGS, locale)},
		{Label: t.PDFAvgMonthlyTurnover, Value: fmt.Sprintf("%.2fx", report.Summary.AverageMonthlyTurnover)},
		{Label: t.PDFSlowMovers, Value: fmt.Sprintf("%d", report.Summary.SlowMovingCount)},
	})

	topRows := make([][]string, 0, minInt(len(report.TopTurnover), 10))
	for i, item := range report.TopTurnover {
		if i >= 10 {
			break
		}
		topRows = append(topRows, []string{
			fmt.Sprintf("%d", i+1),
			item.Name,
			fmt.Sprintf("%d", item.SoldQty),
			fmt.Sprintf("%.2f", item.CurrentStockQty),
			fmt.Sprintf("%.2fx", item.MonthlyTurnoverRatio),
			formatReportCurrency(item.EstimatedCOGS, locale),
		})
	}
	writePDFTable(pdf, t.PDFTopTurnover, []pdfTableColumn{
		{Header: "#", Width: 10, Align: "C"},
		{Header: pdfHumanLabel(t.CSVProduct, locale), Width: 105, Align: "L"},
		{Header: pdfHumanLabel(t.CSVSoldQty, locale), Width: 28, Align: "R"},
		{Header: pdfHumanLabel(t.CSVCurrentStockQty, locale), Width: 32, Align: "R"},
		{Header: pdfHumanLabel(t.CSVMonthlyTurnoverRatio, locale), Width: 34, Align: "R"},
		{Header: pdfHumanLabel(t.CSVEstimatedCOGS, locale), Width: 42, Align: "R"},
	}, topRows, locale)

	slowRows := make([][]string, 0, minInt(len(report.SlowMovers), 10))
	for i, item := range report.SlowMovers {
		if i >= 10 {
			break
		}
		slowRows = append(slowRows, []string{
			fmt.Sprintf("%d", i+1),
			item.Name,
			fmt.Sprintf("%.2f", item.CurrentStockQty),
			fmt.Sprintf("%d", item.SoldQty),
			fmt.Sprintf("%.2fx", item.MonthlyTurnoverRatio),
			formatReportCurrency(item.AverageInventoryValue, locale),
		})
	}
	writePDFTable(pdf, t.PDFSlowMovingItems, []pdfTableColumn{
		{Header: "#", Width: 10, Align: "C"},
		{Header: pdfHumanLabel(t.CSVProduct, locale), Width: 105, Align: "L"},
		{Header: pdfHumanLabel(t.CSVCurrentStockQty, locale), Width: 32, Align: "R"},
		{Header: pdfHumanLabel(t.CSVSoldQty, locale), Width: 28, Align: "R"},
		{Header: pdfHumanLabel(t.CSVMonthlyTurnoverRatio, locale), Width: 34, Align: "R"},
		{Header: pdfHumanLabel(t.CSVAvgInventoryValue, locale), Width: 42, Align: "R"},
	}, slowRows, locale)
	return outputPDF(pdf)
}

func newReportPDF(title string, filters map[string]interface{}, generatedAt string, locale string) *fpdf.Fpdf {
	pdf := fpdf.New("L", "mm", "A4", "")
	registerReportPDFFonts(pdf)
	pdf.SetMargins(10, 10, 10)
	pdf.SetAutoPageBreak(true, 12)
	pdf.SetTitle(title, false)
	pdf.SetAuthor("POS SaaS", false)
	pdf.AliasNbPages("")
	isArabic := normalizeReportLocale(locale) == "ar"
	headerAlign := "L"
	footerAlign := "R"
	dateRangeLabel := "Date Range"
	granularityLabel := "Granularity"
	topLimitLabel := "Top Limit"
	generatedAtLabel := "Generated at"
	pageLabel := "Page"
	if isArabic {
		headerAlign = "R"
		footerAlign = "L"
		dateRangeLabel = asciiSafe("الفترة")
		granularityLabel = asciiSafe("التجميع")
		topLimitLabel = asciiSafe("الحد الأعلى")
		generatedAtLabel = asciiSafe("وقت التوليد")
		pageLabel = "صفحة"
	}
	pdf.SetHeaderFuncMode(func() {
		pdf.SetFillColor(31, 41, 55)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetFont("Arial", "B", 14)
		pdf.CellFormat(0, 10, asciiSafe(title), "", 1, headerAlign, true, 0, "")

		rangeLine := dateRangeLabel
		if from, ok := filters["date_from"]; ok {
			rangeLine += fmt.Sprintf(": %s", formatPDFHeaderValue(fmt.Sprint(from), isArabic))
		}
		if to, ok := filters["date_to"]; ok {
			toValue := formatPDFHeaderValue(fmt.Sprint(to), isArabic)
			if isArabic {
				rangeLine += fmt.Sprintf(" %s %s", asciiSafe("إلى"), toValue)
			} else {
				rangeLine += fmt.Sprintf(" to %s", toValue)
			}
		}
		metaParts := []string{rangeLine}
		if granularity, ok := filters["granularity"]; ok && fmt.Sprint(granularity) != "" {
			metaParts = append(metaParts, fmt.Sprintf("%s: %s", granularityLabel, localizePDFGranularity(fmt.Sprint(granularity), isArabic)))
		}
		if topLimit, ok := filters["top_limit"]; ok && fmt.Sprint(topLimit) != "" {
			metaParts = append(metaParts, fmt.Sprintf("%s: %v", topLimitLabel, topLimit))
		}

		pdf.SetTextColor(75, 85, 99)
		pdf.SetFont("Arial", "", 9)
		if isArabic {
			for _, part := range metaParts {
				pdf.CellFormat(0, 5, fitPDFText(pdf, asciiSafe(part), pdfContentWidth(pdf)), "", 1, headerAlign, false, 0, "")
			}
			pdf.CellFormat(
				0,
				5,
				fitPDFText(pdf, fmt.Sprintf("%s: %s", asciiSafe(generatedAtLabel), formatPDFHeaderDateTime(generatedAt, isArabic)), pdfContentWidth(pdf)),
				"",
				1,
				headerAlign,
				false,
				0,
				"",
			)
		} else {
			pdf.CellFormat(0, 5, fitPDFText(pdf, strings.Join(metaParts, " | "), pdfContentWidth(pdf)), "", 1, headerAlign, false, 0, "")
			pdf.CellFormat(0, 5, fitPDFText(pdf, fmt.Sprintf("%s: %s", generatedAtLabel, formatPDFHeaderDateTime(generatedAt, isArabic)), pdfContentWidth(pdf)), "", 1, headerAlign, false, 0, "")
		}
		pdf.SetDrawColor(229, 231, 235)
		x1, y := pdf.GetX(), pdf.GetY()
		pageW, _, _, _, _, _ := pdfBox(pdf)
		pdf.Line(10, y+1, pageW-10, y+1)
		pdf.SetXY(x1, y+3)
	}, true)
	pdf.SetFooterFunc(func() {
		_, pageH, _, _, _, bottom := pdfBox(pdf)
		pdf.SetY(pageH - bottom + 2)
		pdf.SetFont("Arial", "", 8)
		pdf.SetTextColor(107, 114, 128)
		if isArabic {
			// Render Arabic page label and page counter separately to avoid bidi placeholder reversal.
			pdf.CellFormat(0, 5, fmt.Sprintf("%s %d/{nb}", asciiSafe(pageLabel), pdf.PageNo()), "", 0, footerAlign, false, 0, "")
			return
		}
		pdf.CellFormat(0, 5, fmt.Sprintf("%s %d/{nb}", pageLabel, pdf.PageNo()), "", 0, footerAlign, false, 0, "")
	})
	return pdf
}

func registerReportPDFFonts(pdf *fpdf.Fpdf) {
	if pdf == nil || len(reportArabicFontTTF) == 0 {
		return
	}
	pdf.AddUTF8FontFromBytes("Arial", "", reportArabicFontTTF)
	pdf.AddUTF8FontFromBytes("Arial", "B", reportArabicFontTTF)
}

func writePDFSectionTitle(pdf *fpdf.Fpdf, title string) {
	ensurePDFSpace(pdf, 10)
	pdf.SetDrawColor(229, 231, 235)
	pdf.SetFillColor(249, 250, 251)
	pdf.SetTextColor(17, 24, 39)
	pdf.SetFont("Arial", "B", 11)
	align := "L"
	cellTitle := title
	if containsArabicText(title) {
		align = "R"
		cellTitle = asciiSafe(title)
	}
	pdf.CellFormat(0, 8, cellTitle, "1", 1, align, true, 0, "")
	pdf.Ln(1)
}

func writePDFLine(pdf *fpdf.Fpdf, line string) {
	ensurePDFSpace(pdf, 6)
	pdf.SetTextColor(31, 41, 55)
	pdf.SetFont("Arial", "", 9)
	align := "L"
	if containsArabicText(line) {
		align = "R"
	}
	pdf.MultiCell(0, 5, asciiSafe(line), "", align, false)
}

func outputPDF(pdf *fpdf.Fpdf) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := pdf.Output(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type pdfMetricCard struct {
	Label string
	Value string
}

type pdfTableColumn struct {
	Header string
	Width  float64
	Align  string
}

func writePDFMetricCards(pdf *fpdf.Fpdf, cards []pdfMetricCard) {
	if len(cards) == 0 {
		return
	}
	_, _, left, _, right, _ := pdfBox(pdf)
	contentW := pdfContentWidth(pdf)
	const (
		columnsPerRow = 3
		gap           = 4.0
		cardH         = 16.0
	)
	cardW := (contentW - gap*float64(columnsPerRow-1)) / float64(columnsPerRow)

	for i, card := range cards {
		if i%columnsPerRow == 0 {
			ensurePDFSpace(pdf, cardH+2)
		}
		row := i / columnsPerRow
		col := i % columnsPerRow
		baseY := pdf.GetY() + float64(row*0) // keeps linter quiet; actual Y reset below
		_ = baseY
		if col == 0 {
			// row starts at current Y
		}
		rowStartY := pdf.GetY()
		x := left + float64(col)*(cardW+gap)
		y := rowStartY
		pdf.SetDrawColor(229, 231, 235)
		pdf.SetFillColor(255, 255, 255)
		pdf.Rect(x, y, cardW, cardH, "FD")

		pdf.SetXY(x+2, y+2)
		pdf.SetFont("Arial", "", 8)
		pdf.SetTextColor(107, 114, 128)
		pdf.CellFormat(cardW-4, 4, fitPDFText(pdf, asciiSafe(card.Label), cardW-4), "", 0, "L", false, 0, "")

		pdf.SetXY(x+2, y+7)
		pdf.SetFont("Arial", "B", 11)
		pdf.SetTextColor(17, 24, 39)
		pdf.CellFormat(cardW-4, 6, fitPDFText(pdf, asciiSafe(card.Value), cardW-4), "", 0, "L", false, 0, "")

		if col == columnsPerRow-1 || i == len(cards)-1 {
			pdf.SetXY(left, y+cardH+2)
		} else {
			// keep current row y for next card
			pdf.SetXY(left, y)
		}
	}

	// Reset x after row rendering.
	pdf.SetXY(left, pdf.GetY())
	_ = right
}

func writePDFNotice(pdf *fpdf.Fpdf, text string) {
	ensurePDFSpace(pdf, 8)
	pdf.SetDrawColor(245, 158, 11)
	pdf.SetFillColor(255, 251, 235)
	pdf.SetTextColor(146, 64, 14)
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(0, 7, fitPDFText(pdf, asciiSafe(text), pdfContentWidth(pdf)-2), "1", 1, "L", true, 0, "")
	pdf.Ln(1)
}

func writePDFTable(pdf *fpdf.Fpdf, title string, cols []pdfTableColumn, rows [][]string, locale string) {
	writePDFSectionTitle(pdf, title)
	if len(rows) == 0 {
		writePDFLine(pdf, pdfNoDataText(locale))
		pdf.Ln(1)
		return
	}
	writePDFTableHeader(pdf, cols)
	const rowH = 6.0
	for rowIdx, row := range rows {
		ensurePDFSpaceForTableRow(pdf, rowH, cols)
		fill := rowIdx%2 == 1
		if fill {
			pdf.SetFillColor(249, 250, 251)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		pdf.SetDrawColor(229, 231, 235)
		pdf.SetFont("Arial", "", 8)
		pdf.SetTextColor(31, 41, 55)
		for i, col := range cols {
			val := ""
			if i < len(row) {
				val = row[i]
			}
			align := col.Align
			if align == "" {
				align = "L"
			}
			cellText := fitPDFText(pdf, asciiSafe(val), col.Width-2)
			pdf.CellFormat(col.Width, rowH, cellText, "1", 0, align, fill, 0, "")
		}
		pdf.Ln(-1)
	}
	pdf.Ln(2)
}

func writePDFTableHeader(pdf *fpdf.Fpdf, cols []pdfTableColumn) {
	ensurePDFSpace(pdf, 8)
	pdf.SetFillColor(243, 244, 246)
	pdf.SetDrawColor(209, 213, 219)
	pdf.SetFont("Arial", "B", 8)
	pdf.SetTextColor(17, 24, 39)
	for _, col := range cols {
		align := col.Align
		if align == "" {
			align = "L"
		}
		pdf.CellFormat(col.Width, 7, fitPDFText(pdf, asciiSafe(col.Header), col.Width-2), "1", 0, align, true, 0, "")
	}
	pdf.Ln(-1)
}

func ensurePDFSpaceForTableRow(pdf *fpdf.Fpdf, rowH float64, cols []pdfTableColumn) {
	if !needsNewPDFPage(pdf, rowH) {
		return
	}
	pdf.AddPage()
	writePDFTableHeader(pdf, cols)
}

func ensurePDFSpace(pdf *fpdf.Fpdf, needed float64) {
	if needsNewPDFPage(pdf, needed) {
		pdf.AddPage()
	}
}

func needsNewPDFPage(pdf *fpdf.Fpdf, needed float64) bool {
	_, pageH, _, _, _, bottom := pdfBox(pdf)
	return pdf.GetY()+needed > pageH-bottom
}

func pdfContentWidth(pdf *fpdf.Fpdf) float64 {
	pageW, _, left, _, right, _ := pdfBox(pdf)
	return pageW - left - right
}

func pdfBox(pdf *fpdf.Fpdf) (pageW, pageH, left, top, right, bottom float64) {
	pageW, pageH = pdf.GetPageSize()
	left, top, right, bottom = pdf.GetMargins()
	return
}

func fitPDFText(pdf *fpdf.Fpdf, text string, width float64) string {
	if width <= 0 {
		return ""
	}
	clean := strings.TrimSpace(text)
	if clean == "" {
		return ""
	}
	if pdf.GetStringWidth(clean) <= width {
		return clean
	}
	runes := []rune(clean)
	for len(runes) > 0 {
		runes = runes[:len(runes)-1]
		candidate := strings.TrimSpace(string(runes))
		if candidate == "" {
			break
		}
		if pdf.GetStringWidth(candidate+"...") <= width {
			return candidate + "..."
		}
	}
	return "..."
}

func pdfNoDataText(locale string) string {
	if normalizeReportLocale(locale) == "ar" {
		return "لا توجد بيانات متاحة"
	}
	return "No data available"
}

func pdfHumanLabel(value string, locale string) string {
	label := strings.ReplaceAll(strings.TrimSpace(value), "_", " ")
	if label == "" {
		return value
	}
	if normalizeReportLocale(locale) == "en" {
		switch strings.ToLower(label) {
		case "cogs":
			return "COGS"
		case "gross profit":
			return "Gross Profit"
		case "gross margin pct":
			return "Gross Margin %"
		case "food cost pct":
			return "Food Cost %"
		case "monthly turnover ratio":
			return "Monthly Turnover"
		case "current stock qty":
			return "Current Stock"
		case "sold qty":
			return "Sold Qty"
		case "avg inventory value":
			return "Avg Inventory Value"
		case "estimated cogs":
			return "Estimated COGS"
		}
	}
	return label
}

func localizePDFGranularity(value string, isArabic bool) string {
	if !isArabic {
		return value
	}
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "day":
		return asciiSafe("يومي")
	case "week":
		return asciiSafe("أسبوعي")
	case "month":
		return asciiSafe("شهري")
	default:
		return value
	}
}

func formatPDFHeaderValue(value string, isArabic bool) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if t, err := time.Parse("2006-01-02", value); err == nil {
		if isArabic {
			return t.Format("02/01/2006")
		}
		return t.Format("2006-01-02")
	}
	return value
}

func formatPDFHeaderDateTime(value string, isArabic bool) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		if isArabic {
			return t.UTC().Format("2006-01-02 15:04:05 UTC")
		}
		return t.UTC().Format(time.RFC3339)
	}
	return value
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func asciiSafe(value string) string {
	safe := strings.TrimSpace(value)
	if safe == "" {
		return "N/A"
	}
	var b strings.Builder
	for _, r := range safe {
		if r < 32 && r != '\n' && r != '\r' && r != '\t' {
			continue
		}
		b.WriteRune(r)
	}
	safe = strings.TrimSpace(b.String())
	if safe == "" {
		return "N/A"
	}
	if containsArabicText(safe) {
		return goarabic.Reverse(goarabic.ToGlyph(safe))
	}
	return safe
}

func containsArabicText(s string) bool {
	for _, r := range s {
		if (r >= 0x0600 && r <= 0x06FF) || (r >= 0x0750 && r <= 0x077F) || (r >= 0x08A0 && r <= 0x08FF) {
			return true
		}
	}
	return false
}

func (h *ReportHandler) querySalesSummary(tenantID, restaurantID int64, dateFrom, dateToExclusive time.Time) (salesSummaryResponse, error) {
	const query = `
SELECT
  COALESCE(SUM(o.total_amount), 0) AS total_revenue,
  COUNT(*) AS total_orders,
  COALESCE(AVG(o.total_amount), 0) AS avg_ticket
FROM orders o
WHERE o.tenant_id = $1
  AND o.restaurant_id = $2
  AND o.status <> 'cancelled'
  AND o.created_at >= $3
  AND o.created_at < $4`

	var out salesSummaryResponse
	ctx, cancel := withReportDBTimeout()
	defer cancel()
	err := h.db.QueryRowContext(ctx, query, tenantID, restaurantID, dateFrom, dateToExclusive).
		Scan(&out.TotalRevenue, &out.TotalOrders, &out.AvgTicket)
	return out, err
}

func (h *ReportHandler) querySalesTrend(tenantID, restaurantID int64, dateFrom, dateToExclusive time.Time, granularity string) ([]salesTrendPoint, error) {
	bucketExpr := "to_char(date_trunc('day', o.created_at), 'YYYY-MM-DD')"
	switch granularity {
	case "week":
		bucketExpr = "to_char(date_trunc('week', o.created_at), 'IYYY-\"W\"IW')"
	case "month":
		bucketExpr = "to_char(date_trunc('month', o.created_at), 'YYYY-MM')"
	}

	query := fmt.Sprintf(`
SELECT
  %s AS bucket,
  COALESCE(SUM(o.total_amount), 0) AS revenue,
  COUNT(*) AS orders_count
FROM orders o
WHERE o.tenant_id = $1
  AND o.restaurant_id = $2
  AND o.status <> 'cancelled'
  AND o.created_at >= $3
  AND o.created_at < $4
GROUP BY 1
ORDER BY 1 ASC`, bucketExpr)

	ctx, cancel := withReportDBTimeout()
	defer cancel()
	rows, err := h.db.QueryContext(ctx, query, tenantID, restaurantID, dateFrom, dateToExclusive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	points := make([]salesTrendPoint, 0)
	for rows.Next() {
		var p salesTrendPoint
		if err := rows.Scan(&p.Bucket, &p.Revenue, &p.Orders); err != nil {
			return nil, err
		}
		points = append(points, p)
	}
	return points, rows.Err()
}

func (h *ReportHandler) querySalesHourly(tenantID, restaurantID int64, dateFrom, dateToExclusive time.Time) ([]salesHourlyPoint, error) {
	const query = `
SELECT
  EXTRACT(HOUR FROM o.created_at)::INT AS hour_of_day,
  COALESCE(SUM(o.total_amount), 0) AS revenue,
  COUNT(*) AS orders_count
FROM orders o
WHERE o.tenant_id = $1
  AND o.restaurant_id = $2
  AND o.status <> 'cancelled'
  AND o.created_at >= $3
  AND o.created_at < $4
GROUP BY 1
ORDER BY 1 ASC`

	ctx, cancel := withReportDBTimeout()
	defer cancel()
	rows, err := h.db.QueryContext(ctx, query, tenantID, restaurantID, dateFrom, dateToExclusive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	raw := make(map[int]salesHourlyPoint)
	for rows.Next() {
		var p salesHourlyPoint
		if err := rows.Scan(&p.Hour, &p.Revenue, &p.Orders); err != nil {
			return nil, err
		}
		raw[p.Hour] = p
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	points := make([]salesHourlyPoint, 0, 24)
	for hour := 0; hour < 24; hour++ {
		if p, ok := raw[hour]; ok {
			points = append(points, p)
			continue
		}
		points = append(points, salesHourlyPoint{Hour: hour, Revenue: 0, Orders: 0})
	}
	return points, nil
}

func (h *ReportHandler) queryTopItems(tenantID, restaurantID int64, dateFrom, dateToExclusive time.Time, limit int64) ([]salesTopItem, error) {
	const query = `
SELECT
  oi.product_id,
  COALESCE(NULLIF(TRIM(p.name_en), ''), NULLIF(TRIM(p.name), ''), 'Product #' || oi.product_id::text) AS product_name,
  COALESCE(SUM(oi.quantity), 0)::BIGINT AS qty_sold,
  COALESCE(SUM(oi.total_price), 0) AS revenue
FROM order_items oi
JOIN orders o ON o.id = oi.order_id
LEFT JOIN products p ON p.id = oi.product_id
WHERE oi.tenant_id = $1
  AND oi.restaurant_id = $2
  AND o.status <> 'cancelled'
  AND oi.created_at >= $3
  AND oi.created_at < $4
GROUP BY
  oi.product_id,
  COALESCE(NULLIF(TRIM(p.name_en), ''), NULLIF(TRIM(p.name), ''), 'Product #' || oi.product_id::text)
ORDER BY revenue DESC, qty_sold DESC
LIMIT $5`

	ctx, cancel := withReportDBTimeout()
	defer cancel()
	rows, err := h.db.QueryContext(ctx, query, tenantID, restaurantID, dateFrom, dateToExclusive, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]salesTopItem, 0)
	for rows.Next() {
		var row salesTopItem
		if err := rows.Scan(&row.ProductID, &row.Name, &row.Quantity, &row.Revenue); err != nil {
			return nil, err
		}
		items = append(items, row)
	}
	return items, rows.Err()
}

func (h *ReportHandler) queryTopCategories(tenantID, restaurantID int64, dateFrom, dateToExclusive time.Time, limit int64) ([]salesTopCategory, error) {
	const query = `
SELECT
  COALESCE(c.id, 0) AS category_id,
  COALESCE(NULLIF(TRIM(c.name_en), ''), NULLIF(TRIM(c.name), ''), 'Uncategorized') AS category_name,
  COALESCE(SUM(oi.quantity), 0)::BIGINT AS qty_sold,
  COALESCE(SUM(oi.total_price), 0) AS revenue
FROM order_items oi
JOIN orders o ON o.id = oi.order_id
LEFT JOIN products p ON p.id = oi.product_id
LEFT JOIN categories c ON c.id = p.category_id
WHERE oi.tenant_id = $1
  AND oi.restaurant_id = $2
  AND o.status <> 'cancelled'
  AND oi.created_at >= $3
  AND oi.created_at < $4
GROUP BY
  COALESCE(c.id, 0),
  COALESCE(NULLIF(TRIM(c.name_en), ''), NULLIF(TRIM(c.name), ''), 'Uncategorized')
ORDER BY revenue DESC, qty_sold DESC
LIMIT $5`

	ctx, cancel := withReportDBTimeout()
	defer cancel()
	rows, err := h.db.QueryContext(ctx, query, tenantID, restaurantID, dateFrom, dateToExclusive, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]salesTopCategory, 0)
	for rows.Next() {
		var row salesTopCategory
		if err := rows.Scan(&row.CategoryID, &row.Name, &row.Quantity, &row.Revenue); err != nil {
			return nil, err
		}
		items = append(items, row)
	}
	return items, rows.Err()
}

func (h *ReportHandler) queryProfitabilitySummary(tenantID, restaurantID int64, dateFrom, dateToExclusive time.Time) (profitabilitySummaryResponse, error) {
	const query = `
SELECT
  COALESCE(SUM(oi.total_price), 0) AS total_revenue,
  COALESCE(SUM(oi.quantity * COALESCE(p.cost, 0)), 0) AS total_cogs,
  COALESCE(SUM(CASE WHEN p.cost IS NULL THEN oi.quantity ELSE 0 END), 0)::BIGINT AS units_missing_cost,
  COALESCE(SUM(CASE WHEN p.cost IS NULL THEN 1 ELSE 0 END), 0)::BIGINT AS rows_missing_cost
FROM order_items oi
JOIN orders o ON o.id = oi.order_id
LEFT JOIN products p ON p.id = oi.product_id
WHERE oi.tenant_id = $1
  AND oi.restaurant_id = $2
  AND o.status <> 'cancelled'
  AND oi.created_at >= $3
  AND oi.created_at < $4`

	var out profitabilitySummaryResponse
	ctx, cancel := withReportDBTimeout()
	defer cancel()
	if err := h.db.QueryRowContext(ctx, query, tenantID, restaurantID, dateFrom, dateToExclusive).
		Scan(&out.TotalRevenue, &out.TotalCOGS, &out.UnitsMissingCost, &out.ProductsMissingCostRows); err != nil {
		return out, err
	}
	out.GrossProfit = out.TotalRevenue - out.TotalCOGS
	if out.TotalRevenue > 0 {
		out.GrossMarginPct = (out.GrossProfit / out.TotalRevenue) * 100
		out.FoodCostPct = (out.TotalCOGS / out.TotalRevenue) * 100
	}
	return out, nil
}

func (h *ReportHandler) queryProfitabilityItems(tenantID, restaurantID int64, dateFrom, dateToExclusive time.Time, limit int64) ([]profitabilityItem, []profitabilityItem, error) {
	const query = `
WITH profitability AS (
  SELECT
    oi.product_id,
    COALESCE(NULLIF(TRIM(p.name_en), ''), NULLIF(TRIM(p.name), ''), 'Product #' || oi.product_id::text) AS product_name,
    COALESCE(p.cost, 0) AS cost_per_unit,
    (p.cost IS NULL) AS missing_cost_data,
    COALESCE(SUM(oi.quantity), 0)::BIGINT AS qty_sold,
    COALESCE(SUM(oi.total_price), 0) AS revenue,
    COALESCE(SUM(oi.quantity * COALESCE(p.cost, 0)), 0) AS cogs
  FROM order_items oi
  JOIN orders o ON o.id = oi.order_id
  LEFT JOIN products p ON p.id = oi.product_id
  WHERE oi.tenant_id = $1
    AND oi.restaurant_id = $2
    AND o.status <> 'cancelled'
    AND oi.created_at >= $3
    AND oi.created_at < $4
  GROUP BY
    oi.product_id,
    COALESCE(NULLIF(TRIM(p.name_en), ''), NULLIF(TRIM(p.name), ''), 'Product #' || oi.product_id::text),
    COALESCE(p.cost, 0),
    (p.cost IS NULL)
)
SELECT
  product_id,
  product_name,
  qty_sold,
  revenue,
  cogs,
  (revenue - cogs) AS gross_profit,
  CASE WHEN revenue > 0 THEN ((revenue - cogs) / revenue) * 100 ELSE 0 END AS gross_margin_pct,
  CASE WHEN revenue > 0 THEN (cogs / revenue) * 100 ELSE 0 END AS food_cost_pct,
  cost_per_unit,
  missing_cost_data
FROM profitability
ORDER BY gross_profit DESC, revenue DESC`

	ctx, cancel := withReportDBTimeout()
	defer cancel()
	rows, err := h.db.QueryContext(ctx, query, tenantID, restaurantID, dateFrom, dateToExclusive)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	all := make([]profitabilityItem, 0)
	for rows.Next() {
		var row profitabilityItem
		if err := rows.Scan(
			&row.ProductID,
			&row.Name,
			&row.Quantity,
			&row.Revenue,
			&row.COGS,
			&row.GrossProfit,
			&row.GrossMarginPct,
			&row.FoodCostPct,
			&row.CostPerUnit,
			&row.MissingCostData,
		); err != nil {
			return nil, nil, err
		}
		all = append(all, row)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	top := all
	if int64(len(top)) > limit {
		top = top[:limit]
	}

	lowest := make([]profitabilityItem, 0, len(all))
	for i := len(all) - 1; i >= 0; i-- {
		lowest = append(lowest, all[i])
		if int64(len(lowest)) >= limit {
			break
		}
	}
	return top, lowest, nil
}

func (h *ReportHandler) queryInventoryTurnoverItems(tenantID, restaurantID int64, dateFrom, dateToExclusive time.Time, daysInRange int64) ([]inventoryTurnoverItem, error) {
	const query = `
WITH product_scope AS (
  SELECT
    p.id,
    COALESCE(NULLIF(TRIM(p.name_en), ''), NULLIF(TRIM(p.name), ''), 'Product #' || p.id::text) AS product_name,
    COALESCE(p.cost, 0) AS cost_per_unit,
    COALESCE(p.quantity_in_stock, 0)::DOUBLE PRECISION AS current_qty
  FROM products p
  WHERE p.tenant_id = $1
    AND p.restaurant_id = $2
    AND COALESCE(p.track_inventory, false) = true
    AND COALESCE(p.status, 'active') <> 'deleted'
),
sales AS (
  SELECT
    oi.product_id,
    COALESCE(SUM(oi.quantity), 0)::BIGINT AS sold_qty,
    COALESCE(SUM(oi.quantity * COALESCE(p.cost, 0)), 0) AS estimated_cogs
  FROM order_items oi
  JOIN orders o ON o.id = oi.order_id
  JOIN products p ON p.id = oi.product_id
  WHERE oi.tenant_id = $1
    AND oi.restaurant_id = $2
    AND o.status <> 'cancelled'
    AND oi.created_at >= $3
    AND oi.created_at < $4
  GROUP BY oi.product_id
),
waste AS (
  SELECT
    i.product_id,
    COALESCE(SUM(CASE WHEN i.reason = 'waste' AND i.quantity_change < 0 THEN ABS(i.quantity_change) ELSE 0 END), 0)::BIGINT AS waste_qty
  FROM inventory i
  JOIN product_scope ps ON ps.id = i.product_id
  WHERE i.created_at >= $3
    AND i.created_at < $4
  GROUP BY i.product_id
)
SELECT
  ps.id,
  ps.product_name,
  ps.cost_per_unit,
  COALESCE(opening.opening_qty, ps.current_qty) AS opening_qty,
  ps.current_qty AS current_qty,
  ((COALESCE(opening.opening_qty, ps.current_qty) + ps.current_qty) / 2.0) AS avg_inventory_qty,
  (((COALESCE(opening.opening_qty, ps.current_qty) + ps.current_qty) / 2.0) * ps.cost_per_unit) AS avg_inventory_value,
  COALESCE(s.sold_qty, 0) AS sold_qty,
  COALESCE(s.estimated_cogs, 0) AS estimated_cogs,
  COALESCE(w.waste_qty, 0) AS waste_qty
FROM product_scope ps
LEFT JOIN sales s ON s.product_id = ps.id
LEFT JOIN waste w ON w.product_id = ps.id
LEFT JOIN LATERAL (
  SELECT i.quantity_after::DOUBLE PRECISION AS opening_qty
  FROM inventory i
  WHERE i.product_id = ps.id
    AND i.created_at < $3
  ORDER BY i.created_at DESC, i.id DESC
  LIMIT 1
) opening ON TRUE
ORDER BY estimated_cogs DESC, sold_qty DESC, ps.id ASC`

	ctx, cancel := withReportDBTimeout()
	defer cancel()
	rows, err := h.db.QueryContext(ctx, query, tenantID, restaurantID, dateFrom, dateToExclusive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]inventoryTurnoverItem, 0)
	monthlyMultiplier := 30.0 / float64(daysInRange)
	slowThreshold := 3.0

	for rows.Next() {
		var row inventoryTurnoverItem
		if err := rows.Scan(
			&row.ProductID,
			&row.Name,
			&row.CostPerUnit,
			&row.OpeningStockQty,
			&row.CurrentStockQty,
			&row.AverageInventoryQty,
			&row.AverageInventoryValue,
			&row.SoldQty,
			&row.EstimatedCOGS,
			&row.WasteQty,
		); err != nil {
			return nil, err
		}

		if row.AverageInventoryValue > 0 {
			row.TurnoverRatio = row.EstimatedCOGS / row.AverageInventoryValue
			row.MonthlyTurnoverRatio = row.TurnoverRatio * monthlyMultiplier
		}
		if row.CurrentStockQty > 0 && row.MonthlyTurnoverRatio > 0 && row.MonthlyTurnoverRatio < slowThreshold {
			row.IsSlowMoving = true
		}
		if row.CurrentStockQty > 0 && row.SoldQty == 0 {
			row.IsSlowMoving = true
		}

		items = append(items, row)
	}
	return items, rows.Err()
}

func summarizeInventoryItems(items []inventoryTurnoverItem) inventorySummaryResponse {
	var out inventorySummaryResponse
	var turnoverSum float64
	var turnoverCount int64
	for _, item := range items {
		out.TrackedProductsCount++
		out.TotalCurrentInventoryValue += item.CurrentStockQty * item.CostPerUnit
		out.TotalEstimatedCOGS += item.EstimatedCOGS
		if item.MonthlyTurnoverRatio > 0 {
			turnoverSum += item.MonthlyTurnoverRatio
			turnoverCount++
		}
		if item.IsSlowMoving {
			out.SlowMovingCount++
		}
	}
	if turnoverCount > 0 {
		out.AverageMonthlyTurnover = turnoverSum / float64(turnoverCount)
	}
	return out
}

func splitInventoryItems(items []inventoryTurnoverItem, topLimit int64) ([]inventoryTurnoverItem, []inventoryTurnoverItem) {
	slow := make([]inventoryTurnoverItem, 0)
	for _, item := range items {
		if item.IsSlowMoving {
			slow = append(slow, item)
			if int64(len(slow)) >= topLimit {
				break
			}
		}
	}

	top := make([]inventoryTurnoverItem, 0)
	for _, item := range items {
		if item.MonthlyTurnoverRatio <= 0 {
			continue
		}
		top = append(top, item)
		if int64(len(top)) >= topLimit {
			break
		}
	}
	return slow, top
}

func parseReportTopLimit(raw string, fallback int64) int64 {
	limit := parseInt64(raw, fallback)
	if limit <= 0 {
		return fallback
	}
	if limit > 50 {
		return 50
	}
	return limit
}

func normalizeGranularity(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "week", "weekly":
		return "week"
	case "month", "monthly":
		return "month"
	default:
		return "day"
	}
}

func parseReportDateRange(fromRaw, toRaw string) (time.Time, time.Time, time.Time, error) {
	now := time.Now().UTC()
	toDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	fromDate := toDate.AddDate(0, 0, -29)

	if strings.TrimSpace(fromRaw) != "" {
		parsed, err := time.Parse("2006-01-02", strings.TrimSpace(fromRaw))
		if err != nil {
			return time.Time{}, time.Time{}, time.Time{}, fmt.Errorf("invalid date_from format (expected YYYY-MM-DD)")
		}
		fromDate = parsed.UTC()
	}
	if strings.TrimSpace(toRaw) != "" {
		parsed, err := time.Parse("2006-01-02", strings.TrimSpace(toRaw))
		if err != nil {
			return time.Time{}, time.Time{}, time.Time{}, fmt.Errorf("invalid date_to format (expected YYYY-MM-DD)")
		}
		toDate = parsed.UTC()
	}

	if toDate.Before(fromDate) {
		return time.Time{}, time.Time{}, time.Time{}, fmt.Errorf("date_to must be on or after date_from")
	}
	if toDate.Sub(fromDate).Hours()/24 > maxReportDateRangeDays {
		return time.Time{}, time.Time{}, time.Time{}, fmt.Errorf("date range exceeds maximum of %d days", maxReportDateRangeDays)
	}

	// Exclusive upper bound for SQL filtering.
	toExclusive := toDate.AddDate(0, 0, 1)
	return fromDate, toExclusive, toDate, nil
}
