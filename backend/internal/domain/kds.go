package domain

import "time"

// KDSStatus represents a kitchen queue status.
type KDSStatus string

const (
	KDSStatusNew        KDSStatus = "new"
	KDSStatusInProgress KDSStatus = "in_progress"
	KDSStatusReady      KDSStatus = "ready"
	KDSStatusCompleted  KDSStatus = "completed"
	KDSStatusRecalled   KDSStatus = "recalled"
)

// KDSOrder represents a single queue entry displayed in KDS.
type KDSOrder struct {
	ID             int64      `json:"id"`
	TenantID       int64      `json:"tenant_id"`
	RestaurantID   int64      `json:"restaurant_id"`
	OrderID        int64      `json:"order_id"`
	OrderNumber    string     `json:"order_number"`
	OrderSource    string     `json:"order_source,omitempty"`
	TableID        *int64     `json:"table_id,omitempty"`
	TableNumber    *int       `json:"table_number,omitempty"`
	TableName      string     `json:"table_name,omitempty"`
	CustomerName   string     `json:"customer_name"`
	Status         KDSStatus  `json:"status"`
	Priority       int        `json:"priority"` // 0 = normal, 1 = rush
	StationID      string     `json:"station_id"`
	Items          []KDSItem  `json:"items"`
	Notes          string     `json:"notes,omitempty"`
	ReceivedAt     time.Time  `json:"received_at"`
	StartedAt      *time.Time `json:"started_at,omitempty"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	BumpedBy       *int64     `json:"bumped_by,omitempty"`
	ElapsedSeconds int64      `json:"elapsed_seconds"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// KDSItem stores denormalized order item details for kitchen use.
type KDSItem struct {
	ProductID           int64  `json:"product_id,omitempty"`
	ProductName         string `json:"product_name"`
	Quantity            int    `json:"quantity"`
	SpecialInstructions string `json:"special_instructions,omitempty"`
	AddOns              string `json:"add_ons,omitempty"`
}

// KDSSettings stores display and alert preferences.
type KDSSettings struct {
	ID                  int64     `json:"id"`
	TenantID            int64     `json:"tenant_id"`
	RestaurantID        int64     `json:"restaurant_id"`
	AutoBumpMinutes     int       `json:"auto_bump_minutes"`
	AlertSoundEnabled   bool      `json:"alert_sound_enabled"`
	ShowAnalyticsPanel  bool      `json:"show_analytics_panel"`
	DisplayDensity      string    `json:"display_density"` // compact, comfortable, spacious
	ColorWarningMinutes int       `json:"color_warning_minutes"`
	ColorDangerMinutes  int       `json:"color_danger_minutes"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// KDSStats summarizes queue performance.
type KDSStats struct {
	ActiveOrders   int64   `json:"active_orders"`
	AvgPrepMinutes float64 `json:"avg_prep_minutes"`
	CompletedToday int64   `json:"completed_today"`
	NewOrders      int64   `json:"new_orders"`
	InProgress     int64   `json:"in_progress"`
	Ready          int64   `json:"ready"`
}

// KDSPrepTrendPoint represents one day in KDS preparation trend analytics.
type KDSPrepTrendPoint struct {
	Date            string  `json:"date"`
	AvgPrepMinutes  float64 `json:"avg_prep_minutes"`
	CompletedOrders int64   `json:"completed_orders"`
}

// KDSHourLoadPoint represents order volume for one hour bucket (0-23).
type KDSHourLoadPoint struct {
	Hour       int   `json:"hour"`
	OrderCount int64 `json:"order_count"`
}

// KDSAnalytics contains aggregate preparation analytics used by operations.
type KDSAnalytics struct {
	AveragePrepMinutes   float64             `json:"average_prep_minutes"`
	TotalCompletedOrders int64               `json:"total_completed_orders"`
	PrepTimeTrend        []KDSPrepTrendPoint `json:"prep_time_trend"`
	BusiestHours         []KDSHourLoadPoint  `json:"busiest_hours"`
}

// KDSFilters represents queue fetch filters.
type KDSFilters struct {
	Statuses         []KDSStatus `json:"statuses,omitempty"`
	StationID        string      `json:"station_id,omitempty"`
	IncludeCompleted bool        `json:"include_completed,omitempty"`
}

// ValidKDSStatus checks whether a KDS status is recognized.
func ValidKDSStatus(status string) bool {
	switch KDSStatus(status) {
	case KDSStatusNew, KDSStatusInProgress, KDSStatusReady, KDSStatusCompleted, KDSStatusRecalled:
		return true
	default:
		return false
	}
}

// CanTransitionKDSStatus checks whether a KDS state change is valid.
func CanTransitionKDSStatus(oldStatus, newStatus KDSStatus) bool {
	validTransitions := map[KDSStatus][]KDSStatus{
		KDSStatusNew:        {KDSStatusInProgress},
		KDSStatusInProgress: {KDSStatusReady},
		KDSStatusReady:      {KDSStatusCompleted},
		KDSStatusCompleted:  {KDSStatusRecalled},
		KDSStatusRecalled:   {KDSStatusReady},
	}

	transitions, ok := validTransitions[oldStatus]
	if !ok {
		return false
	}

	for _, status := range transitions {
		if status == newStatus {
			return true
		}
	}

	return false
}
