//go:build legacy
// +build legacy

package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// AnalyticsEvent represents a single tracked event
type AnalyticsEvent struct {
	ID            int64             `json:"id"`
	EventID       string            `json:"event_id" db:"event_id"` // UUID for deduplication
	RestaurantID  int64             `json:"restaurant_id"`
	EventType     string            `json:"event_type"`
	EventName     string            `json:"event_name"`
	Category      string            `json:"category"`
	Action        string            `json:"action"`
	Label         string            `json:"label,omitempty"`
	Value         *int64            `json:"value,omitempty"`
	Properties    map[string]interface{} `json:"properties"`
	UserID        *int64            `json:"user_id,omitempty"`
	SessionID     string            `json:"session_id,omitempty"`
	URL           string            `json:"url,omitempty"`
	Referrer      string            `json:"referrer,omitempty"`
	UserAgent     string            `json:"user_agent,omitempty"`
	IPAddress     string            `json:"ip_address,omitempty"`
	Country       string            `json:"country,omitempty"`
	Region        string            `json:"region,omitempty"`
	City          string            `json:"city,omitempty"`
	Timestamp     time.Time         `json:"timestamp"`
	CreatedAt     time.Time         `json:"created_at"`
}

// Event type constants
const (
	EventTypePageView     = "page_view"
	EventTypeClick        = "click"
	EventTypeScroll       = "scroll"
	EventTypeFormSubmit   = "form_submit"
	EventTypeError        = "error"
	EventTypePerformance  = "performance"
	EventTypeCustom       = "custom"
	EventTypeComponentLoad = "component_load"
	EventTypeComponentError = "component_error"
	EventTypeAssetLoad    = "asset_load"
	EventTypeAPICall      = "api_call"
)

// EventCategory constants
const (
	CategoryUI         = "ui"
	CategoryUser       = "user"
	CategorySystem     = "system"
	CategoryPerformance = "performance"
	CategoryError      = "error"
	CategoryAsset      = "asset"
	CategoryComponent  = "component"
)

// PageViewEvent tracks page views
type PageViewEvent struct {
	URL           string
	Title         string
	Referrer      string
	ScreenWidth   int
	ScreenHeight  int
	ViewportWidth int
	ViewportHeight int
	TimeOnPage    int64 // milliseconds
}

// ClickEvent tracks user clicks
type ClickEvent struct {
	ElementID    string
	ElementClass string
	ElementText  string
	XCoordinate  int
	YCoordinate  int
}

// FormSubmitEvent tracks form submissions
type FormSubmitEvent struct {
	FormID      string
	FormName    string
	FieldCount  int
	SubmitTime  int64 // milliseconds
	HasErrors   bool
	ErrorCount  int
}

// ErrorEvent tracks JavaScript/API errors
type ErrorEvent struct {
	ErrorType    string
	ErrorMessage string
	ErrorStack   string
	LineNumber   int
	ColumnNumber int
	SourceURL    string
	Severity     string // "info", "warning", "error", "critical"
}

// PerformanceEvent tracks performance metrics
type PerformanceEvent struct {
	MetricName       string
	MetricValue      float64
	MetricUnit       string
	NavigationStart  int64 // milliseconds since navigation started
	DOMContentLoaded int64
	LoadComplete     int64
	FirstPaint       int64
	FirstContentfulPaint int64
	LargestContentfulPaint int64
	CumulativeLayoutShift  float64
	TimeToInteractive int64
}

// ComponentLoadEvent tracks component load events
type ComponentLoadEvent struct {
	ComponentID   int64
	ComponentName string
	ComponentType string
	LoadTime      int64 // milliseconds
	Success       bool
	ErrorMessage  string
}

// AssetLoadEvent tracks asset load events
type AssetLoadEvent struct {
	AssetID      int64
	AssetType    string
	AssetSize    int64
	LoadTime     int64 // milliseconds
	Success      bool
	ErrorMessage string
	FormatUsed   string
}

// APICallEvent tracks API calls
type APICallEvent struct {
	Method       string
	Endpoint     string
	StatusCode   int
	ResponseTime int64 // milliseconds
	RequestSize  int64
	ResponseSize int64
	Success      bool
	ErrorMessage string
}

// AnalyticsEventBatch represents a batch of events for bulk insertion
type AnalyticsEventBatch struct {
	Events    []*AnalyticsEvent
	BatchID   string
	Timestamp time.Time
}

// AnalyticsAggregation represents aggregated analytics data
type AnalyticsAggregation struct {
	ID           int64                      `json:"id"`
	RestaurantID int64                      `json:"restaurant_id"`
	EventType    string                     `json:"event_type"`
	Period       string                     `json:"period"` // "hour", "day", "week", "month"
	PeriodStart  time.Time                  `json:"period_start"`
	PeriodEnd    time.Time                  `json:"period_end"`
	EventCount   int64                      `json:"event_count"`
	UniqueUsers  int64                      `json:"unique_users"`
	Properties   map[string]interface{}     `json:"properties"`
	CreatedAt    time.Time                  `json:"created_at"`
	UpdatedAt    time.Time                  `json:"updated_at"`
}

// UserAnalyticsProfile aggregates user behavior
type UserAnalyticsProfile struct {
	ID                int64
	UserID            int64
	RestaurantID      int64
	FirstSeen         time.Time
	LastSeen          time.Time
	TotalSessions     int64
	TotalEvents       int64
	TotalPageViews    int64
	TotalErrors       int64
	AverageSessionDuration int64 // seconds
	BrowserType       string
	BrowserVersion    string
	OSType            string
	OSVersion         string
	DeviceType        string
	Country           string
	Region            string
	City              string
	Properties        map[string]interface{}
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// SessionAnalytics tracks session-level analytics
type SessionAnalytics struct {
	ID                int64
	SessionID         string
	RestaurantID      int64
	UserID            *int64
	StartTime         time.Time
	EndTime           *time.Time
	Duration          int64 // seconds
	PageViewCount     int
	ClickCount        int
	ErrorCount        int
	EventCount        int
	EntryURL          string
	ExitURL           string
	ReferrerURL       string
	DeviceType        string
	BrowserType       string
	Country           string
	ConversionValue   *float64
	Properties        map[string]interface{}
	CreatedAt         time.Time
}

// AnalyticsFilter represents filtering options for analytics queries
type AnalyticsFilter struct {
	EventType      string
	Category       string
	DateFrom       time.Time
	DateTo         time.Time
	UserID         *int64
	SessionID      string
	Country        string
	DeviceType     string
	BrowserType    string
	MinValue       *int64
	MaxValue       *int64
	Properties     map[string]interface{}
	Limit          int
	Offset         int
}

// AnalyticsReport represents a generated analytics report
type AnalyticsReport struct {
	ID              int64
	RestaurantID    int64
	ReportType      string
	Title           string
	Description     string
	DateFrom        time.Time
	DateTo          time.Time
	Metrics         map[string]interface{}
	Dimensions      map[string]interface{}
	Data            []map[string]interface{}
	GeneratedAt     time.Time
	ExpiresAt       *time.Time
	IsPublic        bool
	AccessToken     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Privacy-related models
type UserPrivacySettings struct {
	ID               int64
	UserID           *int64
	RestaurantID     int64
	AllowAnalytics   bool
	AllowTracking    bool
	AllowCookies     bool
	AllowPersonalization bool
	OptOutCategories []string
	DataRetentionDays int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// ConsentRecord tracks user consent for analytics
type ConsentRecord struct {
	ID               int64
	UserID           *int64
	RestaurantID     int64
	ConsentType      string // "analytics", "marketing", "functional"
	Granted          bool
	Version          string
	IPAddress        string
	UserAgent        string
	Timestamp        time.Time
}

// Implement database driver interfaces for JSON storage
func (p UserAnalyticsProfile) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (a AnalyticsAggregation) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a AnalyticsReport) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// EventDimensions represents event attributes for dimensional analysis
type EventDimensions struct {
	EventType       string
	Category        string
	Country         string
	Region          string
	City            string
	DeviceType      string
	BrowserType     string
	BrowserVersion  string
	OSType          string
	PageTitle       string
	PagePath        string
}

// EventMetrics represents event values for metric analysis
type EventMetrics struct {
	EventCount         int64
	UniqueUsers        int64
	UniqueSessions     int64
	AverageValue       float64
	MinValue           int64
	MaxValue           int64
	TotalValue         int64
	ConversionRate     float64
	BounceRate         float64
	AverageSessionTime int64
}

// FunnelStep represents a step in conversion funnel
type FunnelStep struct {
	StepName     string
	EventType    string
	EventName    string
	UserCount    int64
	SessionCount int64
	Conversions  int64
	Dropoffs     int64
}

// FunnelAnalysis represents conversion funnel analysis
type FunnelAnalysis struct {
	FunnelName string
	Steps      []FunnelStep
	TotalUsers int64
	Completion float64 // percentage
}
