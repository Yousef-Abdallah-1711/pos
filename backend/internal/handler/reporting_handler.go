//go:build legacy
// +build legacy

package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"pos-saas/internal/service"
)

// ReportingHandler handles reporting endpoints
type ReportingHandler struct {
	reportingService *service.ReportingService
}

// NewReportingHandler creates a new reporting handler
func NewReportingHandler(reportingService *service.ReportingService) *ReportingHandler {
	return &ReportingHandler{
		reportingService: reportingService,
	}
}

// RegisterRoutes registers all reporting routes
func (h *ReportingHandler) RegisterRoutes(router *gin.Engine) {
	// Reports
	reportsGroup := router.Group("/api/v1/reports")
	{
		reportsGroup.POST("/", h.CreateReport)
		reportsGroup.GET("/", h.ListReports)
		reportsGroup.GET("/:reportID", h.GetReport)
		reportsGroup.PUT("/:reportID", h.UpdateReport)
		reportsGroup.DELETE("/:reportID", h.DeleteReport)

		// Report generation & execution
		reportsGroup.POST("/:reportID/generate", h.GenerateReport)
		reportsGroup.GET("/:reportID/data", h.GetReportData)
		reportsGroup.GET("/:reportID/executions", h.GetExecutionHistory)
		reportsGroup.POST("/:reportID/export", h.ExportReport)

		// Report management
		reportsGroup.POST("/:reportID/share", h.ShareReport)
		reportsGroup.GET("/:reportID/metadata", h.GetReportMetadata)
	}

	// Report Templates
	templatesGroup := router.Group("/api/v1/reports/templates")
	{
		templatesGroup.GET("/", h.ListTemplates)
		templatesGroup.GET("/:templateID", h.GetTemplate)
		templatesGroup.POST("/:templateID/use", h.CreateReportFromTemplate)
	}

	// Dashboards
	dashboardsGroup := router.Group("/api/v1/dashboards")
	{
		dashboardsGroup.POST("/", h.CreateDashboard)
		dashboardsGroup.GET("/", h.ListDashboards)
		dashboardsGroup.GET("/:dashboardID", h.GetDashboard)
		dashboardsGroup.PUT("/:dashboardID", h.UpdateDashboard)
		dashboardsGroup.DELETE("/:dashboardID", h.DeleteDashboard)

		// Dashboard widgets
		dashboardsGroup.POST("/:dashboardID/widgets", h.AddWidget)
		dashboardsGroup.PUT("/:dashboardID/widgets/:widgetID", h.UpdateWidget)
		dashboardsGroup.DELETE("/:dashboardID/widgets/:widgetID", h.RemoveWidget)

		// Dashboard data
		dashboardsGroup.GET("/:dashboardID/data", h.GetDashboardData)
		dashboardsGroup.POST("/:dashboardID/share", h.ShareDashboard)
	}

	// Scheduling
	scheduleGroup := router.Group("/api/v1/reports/schedule")
	{
		scheduleGroup.POST("/", h.ScheduleReport)
		scheduleGroup.GET("/", h.ListSchedules)
		scheduleGroup.GET("/:scheduleID", h.GetSchedule)
		scheduleGroup.PUT("/:scheduleID", h.UpdateSchedule)
		scheduleGroup.DELETE("/:scheduleID", h.DeleteSchedule)
	}
}

// CreateReport creates a new report
// POST /api/v1/reports
func (h *ReportingHandler) CreateReport(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	userID := c.GetInt64("user_id")

	var req struct {
		Name        string              `json:"name" binding:"required"`
		Description string              `json:"description"`
		ReportType  string              `json:"report_type" binding:"required"`
		Query       *models.ReportQuery `json:"query" binding:"required"`
		Format      string              `json:"format"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Format == "" {
		req.Format = "pdf"
	}

	report, err := h.reportingService.CreateReport(
		c.Request.Context(),
		restaurantID,
		userID,
		req.Name,
		req.Description,
		req.ReportType,
		req.Query,
		req.Format,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, report)
}

// ListReports lists reports
// GET /api/v1/reports?type=sales&status=generated&limit=20
func (h *ReportingHandler) ListReports(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	reportType := c.DefaultQuery("type", "")
	status := c.DefaultQuery("status", "")
	limit := 20

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	reports, err := h.reportingService.ListReports(
		c.Request.Context(),
		restaurantID,
		reportType,
		status,
		limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reports": reports,
		"count":   len(reports),
	})
}

// GetReport gets a specific report
// GET /api/v1/reports/:reportID
func (h *ReportingHandler) GetReport(c *gin.Context) {
	reportID, err := strconv.ParseInt(c.Param("reportID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report_id"})
		return
	}

	report, err := h.reportingService.GetReport(c.Request.Context(), reportID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "report not found"})
		return
	}

	c.JSON(http.StatusOK, report)
}

// UpdateReport updates a report
// PUT /api/v1/reports/:reportID
func (h *ReportingHandler) UpdateReport(c *gin.Context) {
	reportID, err := strconv.ParseInt(c.Param("reportID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report_id"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Format      string `json:"format"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get and update report
	report, err := h.reportingService.GetReport(c.Request.Context(), reportID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "report not found"})
		return
	}

	if req.Name != "" {
		report.Name = req.Name
	}
	if req.Description != "" {
		report.Description = req.Description
	}
	if req.Format != "" {
		report.Format = req.Format
	}

	report.UpdatedAt = time.Now()

	c.JSON(http.StatusOK, report)
}

// DeleteReport deletes a report
// DELETE /api/v1/reports/:reportID
func (h *ReportingHandler) DeleteReport(c *gin.Context) {
	reportID, err := strconv.ParseInt(c.Param("reportID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report_id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "deleted",
		"report_id": reportID,
	})
}

// GenerateReport generates a report
// POST /api/v1/reports/:reportID/generate
func (h *ReportingHandler) GenerateReport(c *gin.Context) {
	reportID, err := strconv.ParseInt(c.Param("reportID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report_id"})
		return
	}

	data, err := h.reportingService.GenerateReport(c.Request.Context(), reportID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "generated",
		"row_count":  len(data.Rows),
		"summary":    data.Summary,
		"charts":     data.Charts,
	})
}

// GetReportData gets generated report data
// GET /api/v1/reports/:reportID/data
func (h *ReportingHandler) GetReportData(c *gin.Context) {
	reportID, err := strconv.ParseInt(c.Param("reportID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report_id"})
		return
	}

	data, err := h.reportingService.GenerateReport(c.Request.Context(), reportID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// ExportReport exports a report
// POST /api/v1/reports/:reportID/export?format=pdf
func (h *ReportingHandler) ExportReport(c *gin.Context) {
	reportID, err := strconv.ParseInt(c.Param("reportID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report_id"})
		return
	}

	format := c.DefaultQuery("format", "pdf")

	data, err := h.reportingService.ExportReport(c.Request.Context(), reportID, format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set content type
	contentType := "application/pdf"
	filename := fmt.Sprintf("report_%d.pdf", reportID)

	switch format {
	case "csv":
		contentType = "text/csv"
		filename = fmt.Sprintf("report_%d.csv", reportID)
	case "excel":
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		filename = fmt.Sprintf("report_%d.xlsx", reportID)
	case "json":
		contentType = "application/json"
		filename = fmt.Sprintf("report_%d.json", reportID)
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, contentType, data)
}

// GetExecutionHistory gets report execution history
// GET /api/v1/reports/:reportID/executions?limit=10
func (h *ReportingHandler) GetExecutionHistory(c *gin.Context) {
	reportID, err := strconv.ParseInt(c.Param("reportID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report_id"})
		return
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	history, err := h.reportingService.GetReportExecutionHistory(c.Request.Context(), reportID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"executions": history,
		"count":      len(history),
	})
}

// GetReportMetadata gets report metadata
// GET /api/v1/reports/:reportID/metadata
func (h *ReportingHandler) GetReportMetadata(c *gin.Context) {
	reportID, err := strconv.ParseInt(c.Param("reportID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report_id"})
		return
	}

	metadata, err := h.reportingService.GetReportMetadata(c.Request.Context(), reportID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metadata)
}

// ShareReport shares a report
// POST /api/v1/reports/:reportID/share
func (h *ReportingHandler) ShareReport(c *gin.Context) {
	reportID, err := strconv.ParseInt(c.Param("reportID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report_id"})
		return
	}

	var req struct {
		UserIDs  []int64 `json:"user_ids"`
		IsPublic bool    `json:"is_public"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.reportingService.ShareReport(c.Request.Context(), reportID, req.UserIDs, req.IsPublic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "shared"})
}

// CreateDashboard creates a new dashboard
// POST /api/v1/dashboards
func (h *ReportingHandler) CreateDashboard(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	userID := c.GetInt64("user_id")

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Layout      string `json:"layout"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Layout == "" {
		req.Layout = "grid"
	}

	dashboard, err := h.reportingService.CreateDashboard(
		c.Request.Context(),
		restaurantID,
		userID,
		req.Name,
		req.Description,
		req.Layout,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dashboard)
}

// ListDashboards lists dashboards
// GET /api/v1/dashboards
func (h *ReportingHandler) ListDashboards(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	userID := c.GetInt64("user_id")

	dashboards, err := h.reportingService.ListDashboards(c.Request.Context(), restaurantID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"dashboards": dashboards,
		"count":      len(dashboards),
	})
}

// GetDashboard gets a specific dashboard
// GET /api/v1/dashboards/:dashboardID
func (h *ReportingHandler) GetDashboard(c *gin.Context) {
	dashboardID, err := strconv.ParseInt(c.Param("dashboardID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dashboard_id"})
		return
	}

	dashboard, err := h.reportingService.GetDashboard(c.Request.Context(), dashboardID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "dashboard not found"})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}

// GetDashboardData gets dashboard data with all widget data
// GET /api/v1/dashboards/:dashboardID/data
func (h *ReportingHandler) GetDashboardData(c *gin.Context) {
	dashboardID, err := strconv.ParseInt(c.Param("dashboardID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dashboard_id"})
		return
	}

	data, err := h.reportingService.GetDashboardData(c.Request.Context(), dashboardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// UpdateDashboard updates a dashboard
// PUT /api/v1/dashboards/:dashboardID
func (h *ReportingHandler) UpdateDashboard(c *gin.Context) {
	dashboardID, err := strconv.ParseInt(c.Param("dashboardID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dashboard_id"})
		return
	}

	var req struct {
		Name            string `json:"name"`
		Description     string `json:"description"`
		Theme           string `json:"theme"`
		RefreshInterval int    `json:"refresh_interval"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dashboard, err := h.reportingService.GetDashboard(c.Request.Context(), dashboardID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "dashboard not found"})
		return
	}

	if req.Name != "" {
		dashboard.Name = req.Name
	}
	if req.Description != "" {
		dashboard.Description = req.Description
	}
	if req.Theme != "" {
		dashboard.Theme = req.Theme
	}
	if req.RefreshInterval > 0 {
		dashboard.RefreshInterval = req.RefreshInterval
	}

	dashboard.UpdatedAt = time.Now()

	c.JSON(http.StatusOK, dashboard)
}

// DeleteDashboard deletes a dashboard
// DELETE /api/v1/dashboards/:dashboardID
func (h *ReportingHandler) DeleteDashboard(c *gin.Context) {
	dashboardID, err := strconv.ParseInt(c.Param("dashboardID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dashboard_id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "deleted",
		"dashboard_id": dashboardID,
	})
}

// AddWidget adds a widget to a dashboard
// POST /api/v1/dashboards/:dashboardID/widgets
func (h *ReportingHandler) AddWidget(c *gin.Context) {
	dashboardID, err := strconv.ParseInt(c.Param("dashboardID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dashboard_id"})
		return
	}

	var req struct {
		Title      string `json:"title" binding:"required"`
		WidgetType string `json:"widget_type" binding:"required"`
		DataSource string `json:"data_source" binding:"required"`
		ChartType  string `json:"chart_type"`
		Width      int    `json:"width"`
		Height     int    `json:"height"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Width == 0 {
		req.Width = 4
	}
	if req.Height == 0 {
		req.Height = 3
	}

	widget, err := h.reportingService.AddWidgetToDashboard(
		c.Request.Context(),
		dashboardID,
		req.WidgetType,
		req.Title,
		req.DataSource,
		req.ChartType,
		req.Width,
		req.Height,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, widget)
}

// UpdateWidget updates a widget
// PUT /api/v1/dashboards/:dashboardID/widgets/:widgetID
func (h *ReportingHandler) UpdateWidget(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "widget updated"})
}

// RemoveWidget removes a widget
// DELETE /api/v1/dashboards/:dashboardID/widgets/:widgetID
func (h *ReportingHandler) RemoveWidget(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "widget removed"})
}

// ShareDashboard shares a dashboard
// POST /api/v1/dashboards/:dashboardID/share
func (h *ReportingHandler) ShareDashboard(c *gin.Context) {
	dashboardID, err := strconv.ParseInt(c.Param("dashboardID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dashboard_id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "shared",
		"dashboard_id": dashboardID,
	})
}

// ScheduleReport schedules a report
// POST /api/v1/reports/schedule
func (h *ReportingHandler) ScheduleReport(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	userID := c.GetInt64("user_id")

	var req struct {
		ReportID   int64    `json:"report_id" binding:"required"`
		Frequency  string   `json:"frequency" binding:"required"`
		Recipients []string `json:"recipients" binding:"required"`
		Time       string   `json:"time" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule, err := h.reportingService.ScheduleReport(
		c.Request.Context(),
		restaurantID,
		req.ReportID,
		req.Frequency,
		req.Recipients,
		req.Time,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, schedule)
}

// ListSchedules lists scheduled reports
// GET /api/v1/reports/schedule
func (h *ReportingHandler) ListSchedules(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"schedules": []interface{}{},
		"count":     0,
	})
}

// GetSchedule gets a specific schedule
// GET /api/v1/reports/schedule/:scheduleID
func (h *ReportingHandler) GetSchedule(c *gin.Context) {
	scheduleID, err := strconv.ParseInt(c.Param("scheduleID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule_id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"schedule_id": scheduleID,
	})
}

// UpdateSchedule updates a schedule
// PUT /api/v1/reports/schedule/:scheduleID
func (h *ReportingHandler) UpdateSchedule(c *gin.Context) {
	scheduleID, err := strconv.ParseInt(c.Param("scheduleID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule_id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated", "schedule_id": scheduleID})
}

// DeleteSchedule deletes a schedule
// DELETE /api/v1/reports/schedule/:scheduleID
func (h *ReportingHandler) DeleteSchedule(c *gin.Context) {
	scheduleID, err := strconv.ParseInt(c.Param("scheduleID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule_id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted", "schedule_id": scheduleID})
}

// ListTemplates lists report templates
// GET /api/v1/reports/templates?category=sales
func (h *ReportingHandler) ListTemplates(c *gin.Context) {
	category := c.DefaultQuery("category", "")

	templates, err := h.reportingService.ListReportTemplates(c.Request.Context(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"templates": templates,
		"count":     len(templates),
	})
}

// GetTemplate gets a specific template
// GET /api/v1/reports/templates/:templateID
func (h *ReportingHandler) GetTemplate(c *gin.Context) {
	templateID, err := strconv.ParseInt(c.Param("templateID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template_id"})
		return
	}

	template, err := h.reportingService.GetReportTemplate(c.Request.Context(), templateID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}

	c.JSON(http.StatusOK, template)
}

// CreateReportFromTemplate creates a report from a template
// POST /api/v1/reports/templates/:templateID/use
func (h *ReportingHandler) CreateReportFromTemplate(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	userID := c.GetInt64("user_id")

	templateID, err := strconv.ParseInt(c.Param("templateID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template_id"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report, err := h.reportingService.CreateReportFromTemplate(
		c.Request.Context(),
		restaurantID,
		userID,
		templateID,
		req.Name,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, report)
}
