//go:build legacy
// +build legacy

package repository

import (
	"context"
	"time"
)

// ReportingRepository handles reporting data persistence
type ReportingRepository interface {
	// Reports
	CreateReport(ctx context.Context, report *models.Report) error
	GetReport(ctx context.Context, id int64) (*models.Report, error)
	GetReports(ctx context.Context, restaurantID int64, reportType string, status string, limit int) ([]*models.Report, error)
	UpdateReport(ctx context.Context, report *models.Report) error
	DeleteReport(ctx context.Context, id int64) error

	// Report Execution
	CreateExecution(ctx context.Context, execution *models.ReportExecution) (int64, error)
	GetExecution(ctx context.Context, id int64) (*models.ReportExecution, error)
	UpdateExecution(ctx context.Context, execution *models.ReportExecution) error
	GetExecutions(ctx context.Context, reportID int64, limit int) ([]*models.ReportExecution, error)

	// Report Data
	SaveReportData(ctx context.Context, data *models.ReportData) error
	GetReportData(ctx context.Context, reportID int64) (*models.ReportData, error)

	// Dashboards
	CreateDashboard(ctx context.Context, dashboard *models.Dashboard) error
	GetDashboard(ctx context.Context, id int64) (*models.Dashboard, error)
	GetDashboards(ctx context.Context, restaurantID int64, userID int64) ([]*models.Dashboard, error)
	UpdateDashboard(ctx context.Context, dashboard *models.Dashboard) error
	DeleteDashboard(ctx context.Context, id int64) error

	// Dashboard Widgets
	CreateWidget(ctx context.Context, widget *models.DashboardWidget) error
	GetWidget(ctx context.Context, id int64) (*models.DashboardWidget, error)
	UpdateWidget(ctx context.Context, widget *models.DashboardWidget) error
	DeleteWidget(ctx context.Context, id int64) error
	GetWidgets(ctx context.Context, dashboardID int64) ([]*models.DashboardWidget, error)

	// Report Templates
	GetTemplate(ctx context.Context, id int64) (*models.ReportTemplate, error)
	GetTemplates(ctx context.Context, category string) ([]*models.ReportTemplate, error)
	CreateTemplate(ctx context.Context, template *models.ReportTemplate) error

	// Scheduled Reports
	CreateSchedule(ctx context.Context, schedule *models.ReportSchedule) error
	GetSchedule(ctx context.Context, id int64) (*models.ReportSchedule, error)
	GetSchedules(ctx context.Context, restaurantID int64) ([]*models.ReportSchedule, error)
	UpdateSchedule(ctx context.Context, schedule *models.ReportSchedule) error
	DeleteSchedule(ctx context.Context, id int64) error

	// Report Distribution
	CreateDistribution(ctx context.Context, distribution *models.ReportDistribution) error
	GetDistributions(ctx context.Context, reportID int64) ([]*models.ReportDistribution, error)
	UpdateDistribution(ctx context.Context, distribution *models.ReportDistribution) error

	// Report Metadata
	GetMetadata(ctx context.Context, reportID int64) (*models.ReportMetadata, error)
	UpdateMetadata(ctx context.Context, metadata *models.ReportMetadata) error

	// Configuration
	GetConfig(ctx context.Context, restaurantID int64) (*models.ReportConfig, error)
	UpdateConfig(ctx context.Context, config *models.ReportConfig) error

	// Cleanup/Archive
	ArchiveOldReports(ctx context.Context, restaurantID int64, beforeDate time.Time) error
	GetArchivedReports(ctx context.Context, restaurantID int64, limit int) ([]*models.Report, error)
}
