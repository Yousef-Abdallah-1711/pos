//go:build legacy
// +build legacy

package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// ChurnPrediction predicts likelihood of user/restaurant becoming inactive
type ChurnPrediction struct {
	ID               int64     `gorm:"primaryKey" json:"id"`
	RestaurantID     int64     `json:"restaurant_id"`
	UserID           *int64    `json:"user_id"` // NULL for restaurant-level churn
	ChurnScore       float64   `json:"churn_score"` // 0.0-1.0, higher = more likely to churn
	ChurnProbability float64   `json:"churn_probability"` // 0.0-1.0
	RiskLevel        string    `json:"risk_level"` // "low", "medium", "high", "critical"
	RiskFactors      []string  `json:"risk_factors"` // ["low_engagement", "declining_usage", ...]
	DaysToChurn      int       `json:"days_to_churn"` // Estimated days until churn
	Confidence       float64   `json:"confidence"` // 0.0-1.0
	RecommendedActions []string `json:"recommended_actions"`
	PredictedAt      time.Time `json:"predicted_at"`
	ValidUntil       time.Time `json:"valid_until"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// RevenueForecast predicts future revenue
type RevenueForecast struct {
	ID             int64           `gorm:"primaryKey" json:"id"`
	RestaurantID   int64           `json:"restaurant_id"`
	ForecastPeriod string          `json:"forecast_period"` // "week", "month", "quarter", "year"
	StartDate      time.Time       `json:"start_date"`
	EndDate        time.Time       `json:"end_date"`
	BaseForecast   float64         `json:"base_forecast"` // Base prediction
	OptimisticForecast float64     `json:"optimistic_forecast"` // +1 std dev
	PessimisticForecast float64    `json:"pessimistic_forecast"` // -1 std dev
	Confidence     float64         `json:"confidence"` // 0.0-1.0
	ActualRevenue  *float64        `json:"actual_revenue"` // Filled in after period ends
	Variance       *float64        `json:"variance"` // |actual - base| / base
	DailyForecasts []*DailyForecast `json:"daily_forecasts" gorm:"foreignKey:ForecastID"`
	Drivers        []*RevenueDriver `json:"drivers"`
	PredictedAt    time.Time       `json:"predicted_at"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// DailyForecast represents daily revenue projection
type DailyForecast struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	ForecastID  int64     `json:"forecast_id"`
	Date        time.Time `json:"date"`
	Forecast    float64   `json:"forecast"`
	Actual      *float64  `json:"actual"`
	IsWeekend   bool      `json:"is_weekend"`
	IsHoliday   bool      `json:"is_holiday"`
	CreatedAt   time.Time `json:"created_at"`
}

// RevenueDriver represents factors contributing to revenue forecast
type RevenueDriver struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	ForecastID int64     `json:"forecast_id"`
	Factor     string    `json:"factor"` // "seasonality", "trend", "promotion", "competition", etc.
	Impact     float64   `json:"impact"` // Percentage impact on forecast
	Description string   `json:"description"`
	CreatedAt  time.Time `json:"created_at"`
}

// LifetimeValuePrediction predicts customer lifetime value
type LifetimeValuePrediction struct {
	ID               int64           `gorm:"primaryKey" json:"id"`
	RestaurantID     int64           `json:"restaurant_id"`
	UserID           int64           `json:"user_id"`
	CurrentLTV       float64         `json:"current_ltv"` // Actual cumulative value to date
	PredictedLTV     float64         `json:"predicted_ltv"` // Predicted total lifetime value
	MonthlyValue     float64         `json:"monthly_value"` // Average monthly spending
	PredictedMonths  int             `json:"predicted_months"` // Estimated active months remaining
	ValueSegment     string          `json:"value_segment"` // "low", "medium", "high", "vip"
	GrowthTrend      float64         `json:"growth_trend"` // Growth rate (e.g., 1.05 = 5% growth)
	ChurnRisk        float64         `json:"churn_risk"` // 0.0-1.0
	UpsellOpportunity float64        `json:"upsell_opportunity"` // Potential additional value
	UpsellRecommendations []string   `json:"upsell_recommendations"`
	Confidence       float64         `json:"confidence"` // 0.0-1.0
	PredictedAt      time.Time       `json:"predicted_at"`
	ValidUntil       time.Time       `json:"valid_until"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

// AnomalyDetection detects unusual patterns in metrics
type AnomalyDetection struct {
	ID            int64       `gorm:"primaryKey" json:"id"`
	RestaurantID  int64       `json:"restaurant_id"`
	UserID        *int64      `json:"user_id"`
	MetricType    string      `json:"metric_type"` // "revenue", "traffic", "conversion", "engagement", etc.
	AnomalyScore  float64     `json:"anomaly_score"` // 0.0-1.0, higher = more anomalous
	Severity      string      `json:"severity"` // "info", "warning", "critical"
	Description   string      `json:"description"`
	DetectedValue float64     `json:"detected_value"`
	ExpectedValue float64     `json:"expected_value"`
	Deviation     float64     `json:"deviation"` // (detected - expected) / expected * 100
	AnomalyType   string      `json:"anomaly_type"` // "spike", "drop", "trend_shift", "pattern_break"
	RootCauseHypotheses []string `json:"root_cause_hypotheses"`
	AffectedPeriod string     `json:"affected_period"` // "2024-01-15T10:00Z to 2024-01-15T11:00Z"
	ActionRequired bool       `json:"action_required"`
	Status        string      `json:"status"` // "new", "acknowledged", "investigated", "resolved"
	Investigation *Investigation `json:"investigation"`
	DetectedAt    time.Time  `json:"detected_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// Investigation tracks anomaly investigation status
type Investigation struct {
	Status      string `json:"status"` // "pending", "in_progress", "completed"
	InvestigatedBy int64 `json:"investigated_by"`
	InvestigatedAt *time.Time `json:"investigated_at"`
	RootCause   string `json:"root_cause"`
	Resolution  string `json:"resolution"`
	Notes       string `json:"notes"`
}

// ModelTrainingJob tracks model training progress
type ModelTrainingJob struct {
	ID             int64           `gorm:"primaryKey" json:"id"`
	RestaurantID   int64           `json:"restaurant_id"`
	ModelType      string          `json:"model_type"` // "churn", "revenue", "ltv", "anomaly"
	JobID          string          `json:"job_id"`
	Status         string          `json:"status"` // "queued", "in_progress", "completed", "failed"
	Progress       float64         `json:"progress"` // 0.0-1.0
	TrainingStart  *time.Time      `json:"training_start"`
	TrainingEnd    *time.Time      `json:"training_end"`
	DataPoints     int64           `json:"data_points"`
	Features       []string        `json:"features"`
	Metrics        *ModelMetrics   `json:"metrics"`
	ErrorMessage   *string         `json:"error_message"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// ModelMetrics contains model performance metrics
type ModelMetrics struct {
	Accuracy    float64 `json:"accuracy"`
	Precision   float64 `json:"precision"`
	Recall      float64 `json:"recall"`
	F1Score     float64 `json:"f1_score"`
	AUC         float64 `json:"auc"` // Area under ROC curve
	RMSE        float64 `json:"rmse"` // Root mean squared error
	MAE         float64 `json:"mae"` // Mean absolute error
	MAPE        float64 `json:"mape"` // Mean absolute percentage error
	FeatureImportance map[string]float64 `json:"feature_importance"`
}

// Scan implements the sql.Scanner interface
func (m *ModelMetrics) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return ErrInvalidData
	}
	return json.Unmarshal(bytes, &m)
}

// Value implements the driver.Valuer interface
func (m ModelMetrics) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan implements the sql.Scanner interface
func (i *Investigation) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return ErrInvalidData
	}
	return json.Unmarshal(bytes, &i)
}

// Value implements the driver.Valuer interface
func (i Investigation) Value() (driver.Value, error) {
	return json.Marshal(i)
}

// PredictionConfig stores configuration for prediction models
type PredictionConfig struct {
	ID                      int64           `gorm:"primaryKey" json:"id"`
	RestaurantID            int64           `json:"restaurant_id"`
	ChurnThreshold          float64         `json:"churn_threshold"` // Score above this triggers alert
	AnomalyThreshold        float64         `json:"anomaly_threshold"` // Deviation % threshold
	AnomalyStdDev           float64         `json:"anomaly_std_dev"` // Standard deviations for anomaly
	RevenueForecastDays     int             `json:"revenue_forecast_days"` // Days to forecast ahead
	LTVWindow               int             `json:"ltv_window"` // Days of history for LTV calc
	MinDataPointsForModel   int             `json:"min_data_points_for_model"` // Min data before training
	EnableChurnPrediction   bool            `json:"enable_churn_prediction"`
	EnableRevenueForecasting bool           `json:"enable_revenue_forecasting"`
	EnableLTVPrediction     bool            `json:"enable_ltv_prediction"`
	EnableAnomalyDetection  bool            `json:"enable_anomaly_detection"`
	AutomaticRetraining     bool            `json:"automatic_retraining"`
	RetrainingIntervalDays  int             `json:"retraining_interval_days"`
	CreatedAt               time.Time       `json:"created_at"`
	UpdatedAt               time.Time       `json:"updated_at"`
}

// PredictionInsight combines multiple predictions for actionable insight
type PredictionInsight struct {
	ID              int64           `gorm:"primaryKey" json:"id"`
	RestaurantID    int64           `json:"restaurant_id"`
	UserID          *int64          `json:"user_id"`
	InsightType     string          `json:"insight_type"` // "churn_risk", "upsell_opportunity", "revenue_concern", "anomaly_alert"
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	Priority        string          `json:"priority"` // "low", "medium", "high", "critical"
	Confidence      float64         `json:"confidence"` // 0.0-1.0
	RecommendedActions []string     `json:"recommended_actions"`
	SupportingData  map[string]interface{} `json:"supporting_data"`
	Status          string          `json:"status"` // "new", "acknowledged", "acted_upon", "dismissed"
	AcknowledgedAt  *time.Time      `json:"acknowledged_at"`
	ActedUponAt     *time.Time      `json:"acted_upon_at"`
	CreatedAt       time.Time       `json:"created_at"`
	ExpiresAt       time.Time       `json:"expires_at"`
}

// MetricBaseline stores baseline metrics for anomaly detection
type MetricBaseline struct {
	ID              int64     `gorm:"primaryKey" json:"id"`
	RestaurantID    int64     `json:"restaurant_id"`
	MetricType      string    `json:"metric_type"`
	DayOfWeek       *int      `json:"day_of_week"` // 0-6, NULL for overall
	HourOfDay       *int      `json:"hour_of_day"` // 0-23, NULL for daily
	Mean            float64   `json:"mean"`
	StdDev          float64   `json:"std_dev"`
	Median          float64   `json:"median"`
	P5              float64   `json:"p5"` // 5th percentile
	P95             float64   `json:"p95"` // 95th percentile
	Min             float64   `json:"min"`
	Max             float64   `json:"max"`
	DataPoints      int       `json:"data_points"`
	CalculatedAt    time.Time `json:"calculated_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// PredictionExport represents exported prediction data
type PredictionExport struct {
	ID              int64     `gorm:"primaryKey" json:"id"`
	RestaurantID    int64     `json:"restaurant_id"`
	ExportType      string    `json:"export_type"` // "churn", "revenue", "ltv", "anomalies", "all"
	Format          string    `json:"format"` // "csv", "json", "excel"
	Status          string    `json:"status"` // "generating", "ready", "expired"
	FileURL         string    `json:"file_url"`
	DataRowCount    int       `json:"data_row_count"`
	GeneratedAt     time.Time `json:"generated_at"`
	ExpiresAt       time.Time `json:"expires_at"`
	CreatedAt       time.Time `json:"created_at"`
}
