//go:build legacy
// +build legacy

package service

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"pos-saas/internal/repository"
)

// ReportingService handles report generation and management
type ReportingService struct {
	reportRepo      repository.ReportingRepository
	analyticsRepo   repository.AnalyticsEventRepository
	cacheService    *CacheService
	config          *models.ReportConfig
}

// NewReportingService creates a new reporting service
func NewReportingService(
	reportRepo repository.ReportingRepository,
	analyticsRepo repository.AnalyticsEventRepository,
	cacheService *CacheService,
) *ReportingService {
	return &ReportingService{
		reportRepo:    reportRepo,
		analyticsRepo: analyticsRepo,
		cacheService:  cacheService,
	}
}

// CreateReport creates a new report definition
func (rs *ReportingService) CreateReport(
	ctx context.Context,
	restaurantID int64,
	userID int64,
	name string,
	description string,
	reportType string,
	query *models.ReportQuery,
	format string,
) (*models.Report, error) {
	report := &models.Report{
		RestaurantID:    restaurantID,
		UserID:          userID,
		Name:            name,
		Description:     description,
		ReportType:      reportType,
		Status:          "draft",
		QueryDefinition: query,
		Format:          format,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err := rs.reportRepo.CreateReport(ctx, report)
	if err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	return report, nil
}

// GenerateReport generates a report with the given definition
func (rs *ReportingService) GenerateReport(
	ctx context.Context,
	reportID int64,
) (*models.ReportData, error) {
	// Get report definition
	report, err := rs.reportRepo.GetReport(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("report not found: %w", err)
	}

	// Check cache
	cacheKey := fmt.Sprintf("report:%d:data", reportID)
	if cached, ok := rs.cacheService.Get(cacheKey); ok {
		if data, ok := cached.(*models.ReportData); ok {
			return data, nil
		}
	}

	// Record execution start
	execution := &models.ReportExecution{
		ReportID:    reportID,
		TriggeredBy: "user",
		Status:      "running",
		StartedAt:   time.Now(),
	}

	execID, _ := rs.reportRepo.CreateExecution(ctx, execution)

	// Build and execute query
	data, err := rs.executeReportQuery(ctx, report, report.QueryDefinition)
	if err != nil {
		// Record failure
		execution.Status = "failed"
		execution.ErrorMessage = &[]string{err.Error()}[0]
		completedAt := time.Now()
		execution.CompletedAt = &completedAt
		_ = rs.reportRepo.UpdateExecution(ctx, execution)
		return nil, err
	}

	data.ID = execID
	data.ReportID = reportID
	data.ExecutionID = execID

	// Record execution completion
	completedAt := time.Now()
	execution.Status = "completed"
	execution.CompletedAt = &completedAt
	duration := completedAt.Sub(execution.StartedAt).Milliseconds()
	execution.Duration = &duration
	execution.RowsProcessed = int64(len(data.Rows))
	_ = rs.reportRepo.UpdateExecution(ctx, execution)

	// Update report status
	now := time.Now()
	report.Status = "generated"
	report.GeneratedAt = &now
	report.Rows = int64(len(data.Rows))
	report.Columns = len(data.Headers)
	_ = rs.reportRepo.UpdateReport(ctx, report)

	// Cache for 1 hour
	rs.cacheService.Set(cacheKey, data, 3600, []string{"report", fmt.Sprintf("report:%d", reportID)})

	return data, nil
}

// executeReportQuery executes the report query and returns data
func (rs *ReportingService) executeReportQuery(
	ctx context.Context,
	report *models.Report,
	query *models.ReportQuery,
) (*models.ReportData, error) {
	// In production, would execute actual SQL/query against data sources
	// For now, return mock data

	data := &models.ReportData{
		Headers: []string{"Date", "Revenue", "Orders", "Users", "Conversion"},
		Rows:    make([][]interface{}, 0),
	}

	// Generate mock rows based on time range
	if query.TimeRange != nil {
		current := query.TimeRange.StartDate
		for current.Before(query.TimeRange.EndDate) {
			row := []interface{}{
				current.Format("2006-01-02"),
				5000.0 + (math.Sin(float64(current.Day()))*1000 + 100),
				int64(150 + current.Day()*5),
				int64(200 + current.Day()*3),
				0.035 + (float64(current.Day()%5)*0.002),
			}
			data.Rows = append(data.Rows, row)

			// Increment based on granularity
			switch query.TimeRange.Granularity {
			case "week":
				current = current.AddDate(0, 0, 7)
			case "month":
				current = current.AddDate(0, 1, 0)
			default: // day
				current = current.AddDate(0, 0, 1)
			}
		}
	}

	// Apply limit and offset
	if query.Limit > 0 {
		start := query.Offset
		if start >= len(data.Rows) {
			start = len(data.Rows)
		}
		end := start + query.Limit
		if end > len(data.Rows) {
			end = len(data.Rows)
		}
		data.Rows = data.Rows[start:end]
	}

	// Generate summary
	data.Summary = rs.generateSummary(data, query)

	// Generate chart data if requested
	data.Charts = rs.generateCharts(data, report)

	return data, nil
}

// generateSummary generates summary statistics for the report
func (rs *ReportingService) generateSummary(
	data *models.ReportData,
	query *models.ReportQuery,
) map[string]interface{} {
	summary := make(map[string]interface{})

	if len(data.Rows) == 0 {
		return summary
	}

	// Find numeric columns
	for i, header := range data.Headers {
		switch header {
		case "Revenue", "Orders", "Users":
			total := 0.0
			count := 0

			for _, row := range data.Rows {
				if i < len(row) {
					switch v := row[i].(type) {
					case float64:
						total += v
						count++
					case int64:
						total += float64(v)
						count++
					}
				}
			}

			if count > 0 {
				summary[fmt.Sprintf("%s_total", strings.ToLower(header))] = total
				summary[fmt.Sprintf("%s_avg", strings.ToLower(header))] = total / float64(count)
			}
		}
	}

	summary["row_count"] = len(data.Rows)
	summary["generated_at"] = time.Now()

	return summary
}

// generateCharts generates chart data for visualization
func (rs *ReportingService) generateCharts(
	data *models.ReportData,
	report *models.Report,
) []*models.ChartData {
	charts := make([]*models.ChartData, 0)

	// Extract labels (first column, typically dates)
	labels := make([]string, 0)
	if len(data.Headers) > 0 {
		for _, row := range data.Rows {
			if len(row) > 0 {
				labels = append(labels, fmt.Sprintf("%v", row[0]))
			}
		}
	}

	// Create line chart for numeric columns
	if len(data.Headers) > 1 {
		datasets := make([]map[string]interface{}, 0)

		// Add datasets for each numeric column
		for colIdx := 1; colIdx < len(data.Headers); colIdx++ {
			values := make([]interface{}, 0)

			for _, row := range data.Rows {
				if colIdx < len(row) {
					values = append(values, row[colIdx])
				}
			}

			dataset := map[string]interface{}{
				"label": data.Headers[colIdx],
				"data":  values,
				"borderColor": getColorForIndex(colIdx),
				"tension":     0.4,
			}
			datasets = append(datasets, dataset)
		}

		chart := &models.ChartData{
			Title:     fmt.Sprintf("%s Overview", report.Name),
			ChartType: "line",
			Labels:    labels,
			Datasets:  datasets,
			Options: map[string]interface{}{
				"responsive": true,
				"plugins": map[string]interface{}{
					"legend": map[string]interface{}{
						"position": "top",
					},
				},
			},
		}

		charts = append(charts, chart)
	}

	return charts
}

// GetReport retrieves a report definition
func (rs *ReportingService) GetReport(ctx context.Context, reportID int64) (*models.Report, error) {
	return rs.reportRepo.GetReport(ctx, reportID)
}

// ListReports lists reports for a restaurant
func (rs *ReportingService) ListReports(
	ctx context.Context,
	restaurantID int64,
	reportType string,
	status string,
	limit int,
) ([]*models.Report, error) {
	return rs.reportRepo.GetReports(ctx, restaurantID, reportType, status, limit)
}

// CreateDashboard creates a new dashboard
func (rs *ReportingService) CreateDashboard(
	ctx context.Context,
	restaurantID int64,
	userID int64,
	name string,
	description string,
	layout string,
) (*models.Dashboard, error) {
	dashboard := &models.Dashboard{
		RestaurantID: restaurantID,
		UserID:       userID,
		Name:         name,
		Description:  description,
		Layout:       layout,
		Widgets:      make([]*models.DashboardWidget, 0),
		Theme:        "light",
		RefreshInterval: 300, // 5 minutes
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := rs.reportRepo.CreateDashboard(ctx, dashboard)
	if err != nil {
		return nil, fmt.Errorf("failed to create dashboard: %w", err)
	}

	return dashboard, nil
}

// AddWidgetToDashboard adds a widget to a dashboard
func (rs *ReportingService) AddWidgetToDashboard(
	ctx context.Context,
	dashboardID int64,
	widgetType string,
	title string,
	dataSource string,
	chartType string,
	width int,
	height int,
) (*models.DashboardWidget, error) {
	// Get dashboard to determine position
	dashboard, err := rs.reportRepo.GetDashboard(ctx, dashboardID)
	if err != nil {
		return nil, err
	}

	widget := &models.DashboardWidget{
		DashboardID: dashboardID,
		WidgetType:  widgetType,
		Title:       title,
		Position:    len(dashboard.Widgets),
		Width:       width,
		Height:      height,
		DataSource:  dataSource,
		ChartType:   chartType,
		Config:      make(map[string]interface{}),
		CacheTimeout: 300,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = rs.reportRepo.CreateWidget(ctx, widget)
	if err != nil {
		return nil, err
	}

	return widget, nil
}

// GetDashboard retrieves a dashboard with all widgets
func (rs *ReportingService) GetDashboard(ctx context.Context, dashboardID int64) (*models.Dashboard, error) {
	return rs.reportRepo.GetDashboard(ctx, dashboardID)
}

// ListDashboards lists dashboards for a user
func (rs *ReportingService) ListDashboards(
	ctx context.Context,
	restaurantID int64,
	userID int64,
) ([]*models.Dashboard, error) {
	return rs.reportRepo.GetDashboards(ctx, restaurantID, userID)
}

// ScheduleReport creates a scheduled report
func (rs *ReportingService) ScheduleReport(
	ctx context.Context,
	restaurantID int64,
	reportID int64,
	frequency string,
	recipients []string,
	time string,
) (*models.ReportSchedule, error) {
	schedule := &models.ReportSchedule{
		RestaurantID: restaurantID,
		ReportID:     reportID,
		Frequency:    frequency,
		Recipients:   recipients,
		Time:         time,
		Timezone:     "UTC",
		Format:       "pdf",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Calculate next run
	schedule.NextRun = rs.calculateNextRun(schedule)

	err := rs.reportRepo.CreateSchedule(ctx, schedule)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}

// calculateNextRun calculates when a scheduled report should next run
func (rs *ReportingService) calculateNextRun(schedule *models.ReportSchedule) *time.Time {
	now := time.Now()
	var nextRun time.Time

	switch schedule.Frequency {
	case "daily":
		nextRun = now.AddDate(0, 0, 1)
	case "weekly":
		daysUntilNextWeek := (7 - int(now.Weekday()) + 1) % 7
		if daysUntilNextWeek == 0 {
			daysUntilNextWeek = 7
		}
		nextRun = now.AddDate(0, 0, daysUntilNextWeek)
	case "monthly":
		nextRun = now.AddDate(0, 1, 0)
	case "quarterly":
		nextRun = now.AddDate(0, 3, 0)
	case "yearly":
		nextRun = now.AddDate(1, 0, 0)
	default:
		nextRun = now.AddDate(0, 0, 1)
	}

	return &nextRun
}

// GetReportTemplate gets a predefined report template
func (rs *ReportingService) GetReportTemplate(
	ctx context.Context,
	templateID int64,
) (*models.ReportTemplate, error) {
	return rs.reportRepo.GetTemplate(ctx, templateID)
}

// ListReportTemplates lists available report templates
func (rs *ReportingService) ListReportTemplates(
	ctx context.Context,
	category string,
) ([]*models.ReportTemplate, error) {
	return rs.reportRepo.GetTemplates(ctx, category)
}

// CreateReportFromTemplate creates a report from a template
func (rs *ReportingService) CreateReportFromTemplate(
	ctx context.Context,
	restaurantID int64,
	userID int64,
	templateID int64,
	name string,
) (*models.Report, error) {
	template, err := rs.GetReportTemplate(ctx, templateID)
	if err != nil {
		return nil, err
	}

	report := &models.Report{
		RestaurantID:    restaurantID,
		UserID:          userID,
		Name:            name,
		Description:     template.Description,
		ReportType:      template.Category,
		Status:          "draft",
		QueryDefinition: template.Query,
		Format:          template.Format,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err = rs.reportRepo.CreateReport(ctx, report)
	if err != nil {
		return nil, err
	}

	return report, nil
}

// ExportReport exports a report to a specific format
func (rs *ReportingService) ExportReport(
	ctx context.Context,
	reportID int64,
	format string, // pdf, excel, csv, json
) ([]byte, error) {
	// Get report data
	data, err := rs.GenerateReport(ctx, reportID)
	if err != nil {
		return nil, err
	}

	// Export based on format
	switch format {
	case "csv":
		return rs.exportCSV(data)
	case "json":
		return rs.exportJSON(data)
	case "excel":
		return rs.exportExcel(data)
	case "pdf":
		return rs.exportPDF(data)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// exportCSV exports report data as CSV
func (rs *ReportingService) exportCSV(data *models.ReportData) ([]byte, error) {
	// Build CSV
	csv := ""

	// Headers
	csv += strings.Join(data.Headers, ",") + "\n"

	// Rows
	for _, row := range data.Rows {
		values := make([]string, len(row))
		for i, v := range row {
			values[i] = fmt.Sprintf("%v", v)
		}
		csv += strings.Join(values, ",") + "\n"
	}

	return []byte(csv), nil
}

// exportJSON exports report data as JSON
func (rs *ReportingService) exportJSON(data *models.ReportData) ([]byte, error) {
	// Convert to JSON-friendly format
	rows := make([]map[string]interface{}, len(data.Rows))
	for i, row := range data.Rows {
		rowMap := make(map[string]interface{})
		for j, value := range row {
			if j < len(data.Headers) {
				rowMap[data.Headers[j]] = value
			}
		}
		rows[i] = rowMap
	}

	export := map[string]interface{}{
		"headers":  data.Headers,
		"rows":     rows,
		"summary":  data.Summary,
		"row_count": len(rows),
	}

	return []byte(fmt.Sprintf("%v", export)), nil
}

// exportExcel exports report data as Excel
func (rs *ReportingService) exportExcel(data *models.ReportData) ([]byte, error) {
	// In production, would use a library like github.com/xuri/excelize
	// For now, return a placeholder
	return []byte("EXCEL_DATA_PLACEHOLDER"), nil
}

// exportPDF exports report data as PDF
func (rs *ReportingService) exportPDF(data *models.ReportData) ([]byte, error) {
	// In production, would use a library like github.com/jung-kurt/gofpdf
	// For now, return a placeholder
	return []byte("PDF_DATA_PLACEHOLDER"), nil
}

// GetReportMetadata gets metadata about a report
func (rs *ReportingService) GetReportMetadata(
	ctx context.Context,
	reportID int64,
) (*models.ReportMetadata, error) {
	return rs.reportRepo.GetMetadata(ctx, reportID)
}

// GetDashboardData gets data for all widgets in a dashboard
func (rs *ReportingService) GetDashboardData(
	ctx context.Context,
	dashboardID int64,
) (map[string]interface{}, error) {
	dashboard, err := rs.GetDashboard(ctx, dashboardID)
	if err != nil {
		return nil, err
	}

	dashboardData := map[string]interface{}{
		"id":      dashboard.ID,
		"name":    dashboard.Name,
		"widgets": make([]map[string]interface{}, 0),
	}

	// Get data for each widget
	for _, widget := range dashboard.Widgets {
		widgetData := map[string]interface{}{
			"id":        widget.ID,
			"title":     widget.Title,
			"type":      widget.WidgetType,
			"data":      nil,
		}

		// In production, would fetch actual widget data from reports
		// For now, return empty data
		dashboardData["widgets"].([]map[string]interface{}) = append(
			dashboardData["widgets"].([]map[string]interface{}),
			widgetData,
		)
	}

	return dashboardData, nil
}

// ShareReport shares a report with users
func (rs *ReportingService) ShareReport(
	ctx context.Context,
	reportID int64,
	userIDs []int64,
	isPublic bool,
) error {
	report, err := rs.GetReport(ctx, reportID)
	if err != nil {
		return err
	}

	report.Viewers = userIDs
	report.IsPublic = isPublic

	return rs.reportRepo.UpdateReport(ctx, report)
}

// GetReportExecutionHistory gets execution history for a report
func (rs *ReportingService) GetReportExecutionHistory(
	ctx context.Context,
	reportID int64,
	limit int,
) ([]*models.ReportExecution, error) {
	return rs.reportRepo.GetExecutions(ctx, reportID, limit)
}

// getColorForIndex returns a color for chart dataset
func getColorForIndex(index int) string {
	colors := []string{
		"#FF6384",
		"#36A2EB",
		"#FFCE56",
		"#4BC0C0",
		"#9966FF",
		"#FF9F40",
	}
	return colors[index%len(colors)]
}

// BulkGenerateReports generates reports in batch
func (rs *ReportingService) BulkGenerateReports(
	ctx context.Context,
	reportIDs []int64,
) (map[int64]error, error) {
	results := make(map[int64]error)

	for _, reportID := range reportIDs {
		_, err := rs.GenerateReport(ctx, reportID)
		results[reportID] = err
	}

	return results, nil
}

// GetReportComparisonData compares two time periods
func (rs *ReportingService) GetReportComparisonData(
	ctx context.Context,
	reportID int64,
	period1Start time.Time,
	period1End time.Time,
	period2Start time.Time,
	period2End time.Time,
) (map[string]interface{}, error) {
	// Would generate two reports and compare metrics
	comparison := map[string]interface{}{
		"period1": map[string]interface{}{
			"start":      period1Start,
			"end":        period1End,
			"metrics":    make(map[string]float64),
		},
		"period2": map[string]interface{}{
			"start":      period2Start,
			"end":        period2End,
			"metrics":    make(map[string]float64),
		},
		"variance": make(map[string]float64),
	}

	return comparison, nil
}
