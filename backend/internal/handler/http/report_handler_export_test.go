package http

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"pos-saas/internal/middleware"
	appjwt "pos-saas/internal/pkg/jwt"
)

func TestReportHandler_ExportReport_ResponseHeadersAndPayloadShape(t *testing.T) {
	type testCase struct {
		name             string
		reportType       string
		format           string
		contentType      string
		filenameContains string
		stub             func(h *ReportHandler)
		assertBody       func(t *testing.T, body []byte)
	}

	now := "2026-02-23T00:00:00Z"

	tests := []testCase{
		{
			name:             "sales csv",
			reportType:       "sales",
			format:           "csv",
			contentType:      "text/csv",
			filenameContains: "sales-report-",
			stub: func(h *ReportHandler) {
				h.salesReportDataBuilder = func(_, _ int64, _, _, _ time.Time, _ string, _ int64) (salesReportResponse, error) {
					return salesReportResponse{
						Filters:       map[string]interface{}{"date_from": "2026-02-01", "date_to": "2026-02-23"},
						Summary:       salesSummaryResponse{TotalRevenue: 100, TotalOrders: 4, AvgTicket: 25},
						Trend:         []salesTrendPoint{{Bucket: "2026-02-01", Revenue: 100, Orders: 4}},
						TopItems:      []salesTopItem{{ProductID: 1, Name: "Burger", Quantity: 5, Revenue: 80}},
						TopCategories: []salesTopCategory{{CategoryID: 1, Name: "Main", Quantity: 5, Revenue: 80}},
						GeneratedAt:   now,
					}, nil
				}
			},
			assertBody: func(t *testing.T, body []byte) {
				t.Helper()
				rows, err := csv.NewReader(bytes.NewReader(body)).ReadAll()
				if err != nil {
					t.Fatalf("parse csv: %v", err)
				}
				if len(rows) < 2 {
					t.Fatalf("expected at least 2 csv rows, got %d", len(rows))
				}
				if got := strings.TrimPrefix(rows[0][0], "\uFEFF"); got != "section" {
					t.Fatalf("expected csv first header 'section', got %q", got)
				}
				if !strings.Contains(string(body), "top_items") {
					t.Fatalf("expected sales csv to include top_items section")
				}
			},
		},
		{
			name:             "profitability pdf",
			reportType:       "profitability",
			format:           "pdf",
			contentType:      "application/pdf",
			filenameContains: "profitability-report-summary-",
			stub: func(h *ReportHandler) {
				h.profitabilityReportDataBuilder = func(_, _ int64, _, _, _ time.Time, _ int64) (profitabilityReportResponse, error) {
					return profitabilityReportResponse{
						Filters:     map[string]interface{}{"date_from": "2026-02-01", "date_to": "2026-02-23"},
						Summary:     profitabilitySummaryResponse{TotalRevenue: 200, TotalCOGS: 70, GrossProfit: 130, GrossMarginPct: 65, FoodCostPct: 35},
						TopItems:    []profitabilityItem{{ProductID: 1, Name: "Burger", Quantity: 10, Revenue: 200, COGS: 70, GrossProfit: 130, GrossMarginPct: 65}},
						LowestItems: []profitabilityItem{{ProductID: 2, Name: "Salad", Quantity: 3, Revenue: 30, COGS: 20, GrossProfit: 10, GrossMarginPct: 33.33}},
						GeneratedAt: now,
					}, nil
				}
			},
			assertBody: func(t *testing.T, body []byte) {
				t.Helper()
				if len(body) < 64 {
					t.Fatalf("expected non-trivial pdf body, got %d bytes", len(body))
				}
				if !bytes.HasPrefix(body, []byte("%PDF")) {
					t.Fatalf("expected PDF payload prefix, got %q", string(body[:min(len(body), 8)]))
				}
			},
		},
		{
			name:             "inventory csv",
			reportType:       "inventory",
			format:           "csv",
			contentType:      "text/csv",
			filenameContains: "inventory-report-",
			stub: func(h *ReportHandler) {
				h.inventoryReportDataBuilder = func(_, _ int64, _, _, _ time.Time, _ int64, _ int64) (inventoryReportResponse, error) {
					item := inventoryTurnoverItem{
						ProductID: 1, Name: "Potatoes", CurrentStockQty: 100, SoldQty: 10,
						AverageInventoryValue: 50, EstimatedCOGS: 20, MonthlyTurnoverRatio: 1.2, IsSlowMoving: true,
					}
					return inventoryReportResponse{
						Filters:     map[string]interface{}{"date_from": "2026-02-01", "date_to": "2026-02-23"},
						Summary:     inventorySummaryResponse{TrackedProductsCount: 1, TotalCurrentInventoryValue: 50, TotalEstimatedCOGS: 20, AverageMonthlyTurnover: 1.2, SlowMovingCount: 1},
						Items:       []inventoryTurnoverItem{item},
						SlowMovers:  []inventoryTurnoverItem{item},
						TopTurnover: []inventoryTurnoverItem{item},
						GeneratedAt: now,
					}, nil
				}
			},
			assertBody: func(t *testing.T, body []byte) {
				t.Helper()
				rows, err := csv.NewReader(bytes.NewReader(body)).ReadAll()
				if err != nil {
					t.Fatalf("parse csv: %v", err)
				}
				if len(rows) < 2 {
					t.Fatalf("expected at least 2 csv rows, got %d", len(rows))
				}
				if rows[0][1] != "product" {
					t.Fatalf("expected inventory csv product column, got %q", rows[0][1])
				}
				if !strings.Contains(string(body), "Potatoes") {
					t.Fatalf("expected inventory csv to include item name")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := &ReportHandler{db: &sql.DB{}}
			tc.stub(h)

			var auditCalls []reportExportAuditRecord
			h.reportExportAuditLogger = func(_ context.Context, rec reportExportAuditRecord) error {
				auditCalls = append(auditCalls, rec)
				return nil
			}

			reqBody := `{"report_type":"` + tc.reportType + `","format":"` + tc.format + `","date_from":"2026-02-01","date_to":"2026-02-23","granularity":"day","top_limit":10}`
			req := newReportExportRequestForTest(reqBody)
			rr := httptest.NewRecorder()

			h.ExportReport(rr, req)

			if rr.Code != http.StatusOK {
				t.Fatalf("expected 200, got %d body=%s", rr.Code, rr.Body.String())
			}
			if ct := rr.Header().Get("Content-Type"); !strings.Contains(ct, tc.contentType) {
				t.Fatalf("expected content-type containing %q, got %q", tc.contentType, ct)
			}
			cd := rr.Header().Get("Content-Disposition")
			if !strings.Contains(cd, "attachment;") || !strings.Contains(cd, tc.filenameContains) {
				t.Fatalf("unexpected content-disposition: %q", cd)
			}
			tc.assertBody(t, rr.Body.Bytes())

			if len(auditCalls) != 1 {
				t.Fatalf("expected 1 audit call, got %d", len(auditCalls))
			}
			if auditCalls[0].ReportType != tc.reportType {
				t.Fatalf("expected audit report_type=%s, got %s", tc.reportType, auditCalls[0].ReportType)
			}
			if auditCalls[0].Format != tc.format {
				t.Fatalf("expected audit format=%s, got %s", tc.format, auditCalls[0].Format)
			}
			if auditCalls[0].TenantID != 8 || auditCalls[0].RestaurantID != 8 || auditCalls[0].UserID != 7 {
				t.Fatalf("unexpected audit IDs: %+v", auditCalls[0])
			}
			if auditCalls[0].Filters["date_from"] != "2026-02-01" || auditCalls[0].Filters["date_to"] != "2026-02-23" {
				t.Fatalf("unexpected audit filters: %+v", auditCalls[0].Filters)
			}
		})
	}
}

func TestReportHandler_ExportReport_AuditFailureDoesNotFailExport(t *testing.T) {
	h := &ReportHandler{
		db: &sql.DB{},
		salesReportDataBuilder: func(_, _ int64, _, _, _ time.Time, _ string, _ int64) (salesReportResponse, error) {
			return salesReportResponse{
				Filters:     map[string]interface{}{"date_from": "2026-02-01", "date_to": "2026-02-23"},
				Summary:     salesSummaryResponse{TotalRevenue: 10, TotalOrders: 1, AvgTicket: 10},
				GeneratedAt: "2026-02-23T00:00:00Z",
			}, nil
		},
		reportExportAuditLogger: func(context.Context, reportExportAuditRecord) error {
			return context.DeadlineExceeded
		},
	}

	req := newReportExportRequestForTest(`{"report_type":"sales","format":"csv"}`)
	rr := httptest.NewRecorder()

	h.ExportReport(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected export to succeed despite audit failure, got %d body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Header().Get("Content-Type"), "text/csv") {
		t.Fatalf("expected csv content-type, got %q", rr.Header().Get("Content-Type"))
	}
}

func TestReportHandler_ExportReport_ArabicCSVLocaleLabels(t *testing.T) {
	h := &ReportHandler{
		db: &sql.DB{},
		salesReportDataBuilder: func(_, _ int64, _, _, _ time.Time, _ string, _ int64) (salesReportResponse, error) {
			return salesReportResponse{
				Filters:       map[string]interface{}{"date_from": "2026-02-01", "date_to": "2026-02-23"},
				Summary:       salesSummaryResponse{TotalRevenue: 100, TotalOrders: 4, AvgTicket: 25},
				Trend:         []salesTrendPoint{{Bucket: "2026-02-01", Revenue: 100, Orders: 4}},
				TopItems:      []salesTopItem{{ProductID: 1, Name: "Burger", Quantity: 5, Revenue: 80}},
				TopCategories: []salesTopCategory{{CategoryID: 1, Name: "Appetizers", Quantity: 5, Revenue: 80}},
				GeneratedAt:   "2026-02-23T00:00:00Z",
			}, nil
		},
		reportExportAuditLogger: func(context.Context, reportExportAuditRecord) error { return nil },
	}

	req := newReportExportRequestForTest(`{"report_type":"sales","format":"csv","locale":"ar","date_from":"2026-02-01","date_to":"2026-02-23"}`)
	rr := httptest.NewRecorder()

	h.ExportReport(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rr.Code, rr.Body.String())
	}
	if got := rr.Header().Get("Content-Language"); got != "ar" {
		t.Fatalf("expected Content-Language=ar, got %q", got)
	}
	rows, err := csv.NewReader(bytes.NewReader(rr.Body.Bytes())).ReadAll()
	if err != nil {
		t.Fatalf("parse csv: %v", err)
	}
	if len(rows) < 2 {
		t.Fatalf("expected at least 2 rows, got %d", len(rows))
	}
	if got := strings.TrimPrefix(rows[0][0], "\uFEFF"); got != "القسم" {
		t.Fatalf("expected Arabic csv header, got %q", got)
	}
	if got := rows[1][0]; got != "ملخص" {
		t.Fatalf("expected Arabic section label, got %q", got)
	}
}

func newReportExportRequestForTest(body string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reports/export", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", "203.0.113.10")
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantContextKey, int64(8)))
	req = req.WithContext(context.WithValue(req.Context(), middleware.RestaurantContextKey, int64(8)))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserContextKey, &appjwt.Claims{UserID: 7, TenantID: 8, RestaurantID: 8}))
	return req
}
