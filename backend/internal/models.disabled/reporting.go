//go:build legacy
// +build legacy

package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Report represents a generated report
type Report struct {
	ID               int64                  `gorm:"primaryKey" json:"id"`
	RestaurantID     int64                  `json:"restaurant_id"`
	UserID           int64                  `json:"user_id"`
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	ReportType       string                 `json:"report_type"` // "sales", "inventory", "users", "custom"
	Status           string                 `json:"status"` // "draft", "scheduled", "generated", "archived"
	QueryDefinition  *ReportQuery           `json:"query_definition"`
	DataSource       string                 `json:"data_source"` // "database", "analytics", "combined"
	Format           string                 `json:"format"` // "pdf", "excel", "csv", "html"
	Rows             int64                  `json:"rows"` // Number of data rows
	Columns          int                    `json:"columns"` // Number of columns
	FileURL          *string                `json:"file_url"` // Location of generated file
	FileSize         *int64                 `json:"file_size"` // Bytes
	GeneratedAt      *time.Time             `json:"generated_at"`
	ExpiresAt        *time.Time             `json:"expires_at"`
	IsPublic         bool                   `json:"is_public"`
	Viewers          []int64                `json:"viewers"` // User IDs who can view
	Tags             []string               `json:"tags"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	ScheduledAt      *time.Time             `json:"scheduled_at"`
}

// ReportQuery defines the query parameters for a report
type ReportQuery struct {
	Filters      []*ReportFilter       `json:"filters"`
	Dimensions   []string              `json:"dimensions"` // Group by fields
	Metrics      []string              `json:"metrics"` // Aggregation fields
	SortBy       string                `json:"sort_by"`
	SortOrder    string                `json:"sort_order"` // asc, desc
	TimeRange    *TimeRange            `json:"time_range"`
	Limit        int                   `json:"limit"`
	Offset       int                   `json:"offset"`
}

// ReportFilter represents a single filter
type ReportFilter struct {
	Field    string        `json:"field"`
	Operator string        `json:"operator"` // equals, contains, gt, lt, between, in
	Value    interface{}   `json:"value"`
}

// TimeRange represents a time period
type TimeRange struct {
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Granularity string   `json:"granularity"` // day, week, month, quarter, year
}

// Scan implements sql.Scanner
func (rq *ReportQuery) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return ErrInvalidData
	}
	return json.Unmarshal(bytes, &rq)
}

// Value implements driver.Valuer
func (rq ReportQuery) Value() (driver.Value, error) {
	return json.Marshal(rq)
}

// Dashboard represents a custom dashboard
type Dashboard struct {
	ID               int64           `gorm:"primaryKey" json:"id"`
	RestaurantID     int64           `json:"restaurant_id"`
	UserID           int64           `json:"user_id"`
	Name             string          `json:"name"`
	Description      string          `json:"description"`
	Layout           string          `json:"layout"` // "grid", "rows", "custom"
	Widgets          []*DashboardWidget `json:"widgets"`
	IsPublic         bool            `json:"is_public"`
	SharedWith       []int64         `json:"shared_with"`
	Theme            string          `json:"theme"` // "light", "dark"
	RefreshInterval  int             `json:"refresh_interval"` // seconds
	IsDefault        bool            `json:"is_default"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

// DashboardWidget represents a widget on a dashboard
type DashboardWidget struct {
	ID              int64           `gorm:"primaryKey" json:"id"`
	DashboardID     int64           `json:"dashboard_id"`
	WidgetType      string          `json:"widget_type"` // "chart", "metric", "table", "gauge", "map"
	Title           string          `json:"title"`
	Position        int             `json:"position"` // Order in dashboard
	Width           int             `json:"width"` // Grid width (1-12)
	Height          int             `json:"height"` // Grid height (1-4)
	DataSource      string          `json:"data_source"` // Report ID or query
	ChartType       string          `json:"chart_type"` // "line", "bar", "pie", "area", etc.
	Config          map[string]interface{} `json:"config"`
	CacheTimeout    int             `json:"cache_timeout"` // seconds
	Drill           *DrillThrough   `json:"drill"` // Drill-down capability
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// DrillThrough represents drill-down capability
type DrillThrough struct {
	TargetReport  string                 `json:"target_report"`
	MappedFields  map[string]string      `json:"mapped_fields"` // source -> target
	PreserveFilters bool                 `json:"preserve_filters"`
}

// ReportTemplate represents a predefined report template
type ReportTemplate struct {
	ID            int64           `gorm:"primaryKey" json:"id"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	Category      string          `json:"category"` // "sales", "inventory", "users", "finance"
	Query         *ReportQuery    `json:"query"`
	Format        string          `json:"format"`
	IsPublic      bool            `json:"is_public"`
	CreatedBy     int64           `json:"created_by"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// ReportSchedule represents a scheduled report
type ReportSchedule struct {
	ID              int64           `gorm:"primaryKey" json:"id"`
	RestaurantID    int64           `json:"restaurant_id"`
	ReportID        int64           `json:"report_id"`
	CreatedBy       int64           `json:"created_by"`
	Frequency       string          `json:"frequency"` // "daily", "weekly", "monthly", "quarterly", "yearly"
	DayOfWeek       *int            `json:"day_of_week"` // 0-6 for weekly
	DayOfMonth      *int            `json:"day_of_month"` // 1-31 for monthly
	Time            string          `json:"time"` // HH:MM format
	Timezone        string          `json:"timezone"`
	Recipients      []string        `json:"recipients"` // Email addresses
	Format          string          `json:"format"` // pdf, excel, csv, html
	IncludeCharts   bool            `json:"include_charts"`
	IsActive        bool            `json:"is_active"`
	LastRun         *time.Time      `json:"last_run"`
	NextRun         *time.Time      `json:"next_run"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// ReportExecution tracks report generation
type ReportExecution struct {
	ID              int64           `gorm:"primaryKey" json:"id"`
	ReportID        int64           `json:"report_id"`
	TriggeredBy     string          `json:"triggered_by"` // "user", "schedule", "api"
	Status          string          `json:"status"` // "pending", "running", "completed", "failed"
	StartedAt       time.Time       `json:"started_at"`
	CompletedAt     *time.Time      `json:"completed_at"`
	Duration        *int64          `json:"duration"` // milliseconds
	RowsProcessed   int64           `json:"rows_processed"`
	ErrorMessage    *string         `json:"error_message"`
	CreatedAt       time.Time       `json:"created_at"`
}

// ReportDistribution tracks report sharing
type ReportDistribution struct {
	ID              int64           `gorm:"primaryKey" json:"id"`
	ReportID        int64           `json:"report_id"`
	DistributedTo   string          `json:"distributed_to"` // Email or user
	Method          string          `json:"method"` // "email", "download", "view"
	DistributedAt   time.Time       `json:"distributed_at"`
	Status          string          `json:"status"` // "pending", "sent", "failed"
	CreatedAt       time.Time       `json:"created_at"`
}

// ReportConfig stores reporting system configuration
type ReportConfig struct {
	ID                      int64           `gorm:"primaryKey" json:"id"`
	RestaurantID            int64           `json:"restaurant_id"`
	MaxReportRows           int64           `json:"max_report_rows"`
	MaxConcurrentGenerations int             `json:"max_concurrent_generations"`
	ReportRetentionDays     int             `json:"report_retention_days"`
	DefaultFormat           string          `json:"default_format"`
	EnableScheduling        bool            `json:"enable_scheduling"`
	EnableSharing           bool            `json:"enable_sharing"`
	EnableExport            bool            `json:"enable_export"`
	ExportFormats           []string        `json:"export_formats"` // pdf, excel, csv, json
	SmtpServer              string          `json:"smtp_server"`
	SmtpPort                int             `json:"smtp_port"`
	CreatedAt               time.Time       `json:"created_at"`
	UpdatedAt               time.Time       `json:"updated_at"`
}

// ReportMetadata tracks report statistics
type ReportMetadata struct {
	ID              int64           `gorm:"primaryKey" json:"id"`
	ReportID        int64           `json:"report_id"`
	TotalExecutions int             `json:"total_executions"`
	AvgExecutionTime int64           `json:"avg_execution_time"` // ms
	TotalRows       int64           `json:"total_rows"`
	TotalViews      int             `json:"total_views"`
	LastViewed      *time.Time      `json:"last_viewed"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// ReportData represents the actual data in a report
type ReportData struct {
	ID              int64           `gorm:"primaryKey" json:"id"`
	ReportID        int64           `json:"report_id"`
	ExecutionID     int64           `json:"execution_id"`
	Headers         []string        `json:"headers"`
	Rows            [][]interface{} `json:"rows"`
	Summary         map[string]interface{} `json:"summary"`
	Charts          []*ChartData    `json:"charts"`
}

// ChartData represents chart data for a dashboard widget
type ChartData struct {
	Title       string                   `json:"title"`
	ChartType   string                   `json:"chart_type"`
	Labels      []string                 `json:"labels"`
	Datasets    []map[string]interface{} `json:"datasets"`
	Options     map[string]interface{}   `json:"options"`
}
