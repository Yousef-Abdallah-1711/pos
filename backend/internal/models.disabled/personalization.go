//go:build legacy
// +build legacy

package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// UserBehaviorProfile stores comprehensive user behavior data for ML
type UserBehaviorProfile struct {
	ID                 int64                      `json:"id" db:"id"`
	RestaurantID       int64                      `json:"restaurant_id" db:"restaurant_id"`
	UserID             int64                      `json:"user_id" db:"user_id"`
	ComponentViews     map[string]int64           `json:"component_views" db:"component_views"`         // componentID -> view count
	ComponentClicks    map[string]int64           `json:"component_clicks" db:"component_clicks"`       // componentID -> click count
	ComponentInteractions map[string]int64        `json:"component_interactions" db:"component_interactions"` // componentID -> interaction count
	CategoryPreferences map[string]float64        `json:"category_preferences" db:"category_preferences"` // category -> preference score (0-1)
	TagPreferences     map[string]float64        `json:"tag_preferences" db:"tag_preferences"`         // tag -> preference score (0-1)
	TimeSpentPerComponent map[string]int64       `json:"time_spent_per_component" db:"time_spent_per_component"` // componentID -> milliseconds
	LastActiveTime     time.Time                 `json:"last_active_time" db:"last_active_time"`
	SessionCount       int64                     `json:"session_count" db:"session_count"`
	TotalInteractions  int64                     `json:"total_interactions" db:"total_interactions"`
	AverageSessionDuration int64                 `json:"average_session_duration" db:"average_session_duration"` // seconds
	PreferredThemes    []string                  `json:"preferred_themes" db:"preferred_themes"`
	DeviceType         string                    `json:"device_type" db:"device_type"` // mobile, tablet, desktop
	BrowserType        string                    `json:"browser_type" db:"browser_type"`
	UpdatedAt          time.Time                 `json:"updated_at" db:"updated_at"`
	CreatedAt          time.Time                 `json:"created_at" db:"created_at"`
}

// Scan implements sql.Scanner interface for JSON deserialization
func (ubp *UserBehaviorProfile) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), ubp)
}

// Value implements driver.Valuer interface for JSON serialization
func (ubp UserBehaviorProfile) Value() (driver.Value, error) {
	return json.Marshal(ubp)
}

// UserFeatures represents extracted features for ML model
type UserFeatures struct {
	ID                    int64     `json:"id" db:"id"`
	UserBehaviorProfileID int64     `json:"user_behavior_profile_id" db:"user_behavior_profile_id"`
	RestaurantID          int64     `json:"restaurant_id" db:"restaurant_id"`
	UserID                int64     `json:"user_id" db:"user_id"`
	EngagementScore       float64   `json:"engagement_score" db:"engagement_score"`             // 0-1: based on interactions
	RecencyScore          float64   `json:"recency_score" db:"recency_score"`                 // 0-1: recent activity weight
	DiversityScore        float64   `json:"diversity_score" db:"diversity_score"`             // 0-1: category/tag diversity
	LoyaltyScore          float64   `json:"loyalty_score" db:"loyalty_score"`                 // 0-1: session frequency
	ConversionLikelihood  float64   `json:"conversion_likelihood" db:"conversion_likelihood"` // 0-1: predicted conversion
	ChurnRisk             float64   `json:"churn_risk" db:"churn_risk"`                       // 0-1: risk of becoming inactive
	PreferredHours        []int     `json:"preferred_hours" db:"preferred_hours"`             // hours 0-23 when user is active
	FeatureVector         []float64 `json:"feature_vector" db:"feature_vector"`               // normalized feature values
	ModelVersion          string    `json:"model_version" db:"model_version"`                 // version of feature extraction pipeline
	ComputedAt            time.Time `json:"computed_at" db:"computed_at"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
}

// ContentRecommendation represents a single recommendation
type ContentRecommendation struct {
	ID                   int64            `json:"id" db:"id"`
	RestaurantID         int64            `json:"restaurant_id" db:"restaurant_id"`
	UserID               int64            `json:"user_id" db:"user_id"`
	ComponentID          int64            `json:"component_id" db:"component_id"`
	RecommendationType   string           `json:"recommendation_type" db:"recommendation_type"` // collaborative, content_based, hybrid, trending, trending
	RelevanceScore       float64          `json:"relevance_score" db:"relevance_score"`         // 0-1: confidence in recommendation
	Reason               string           `json:"reason" db:"reason"`                           // "Similar to viewed", "Popular in your category", etc.
	Position             int32            `json:"position" db:"position"`                       // rank in recommendation list
	ExposureCount        int64            `json:"exposure_count" db:"exposure_count"`           // how many times shown
	ClickCount           int64            `json:"click_count" db:"click_count"`                 // how many times clicked
	ConversionValue      *float64         `json:"conversion_value" db:"conversion_value"`       // revenue if converted
	Properties           map[string]interface{} `json:"properties" db:"properties"`
	IsActive             bool             `json:"is_active" db:"is_active"`
	ExpiresAt            time.Time        `json:"expires_at" db:"expires_at"`
	CreatedAt            time.Time        `json:"created_at" db:"created_at"`
}

// Scan implements sql.Scanner interface
func (cr *ContentRecommendation) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), cr)
}

// Value implements driver.Valuer interface
func (cr ContentRecommendation) Value() (driver.Value, error) {
	return json.Marshal(cr)
}

// MLModelMetadata stores information about trained ML models
type MLModelMetadata struct {
	ID                 int64                  `json:"id" db:"id"`
	RestaurantID       int64                  `json:"restaurant_id" db:"restaurant_id"`
	ModelType          string                 `json:"model_type" db:"model_type"`               // personalization, churn_prediction, etc.
	ModelVersion       string                 `json:"model_version" db:"model_version"`         // semantic version
	Algorithm          string                 `json:"algorithm" db:"algorithm"`                 // collaborative_filtering, content_based, random_forest, etc.
	TrainingDataSize   int64                  `json:"training_data_size" db:"training_data_size"`
	Accuracy           float64                `json:"accuracy" db:"accuracy"`                   // 0-1: model accuracy on test set
	Precision          float64                `json:"precision" db:"precision"`                 // 0-1: precision metric
	Recall             float64                `json:"recall" db:"recall"`                       // 0-1: recall metric
	F1Score            float64                `json:"f1_score" db:"f1_score"`                   // 0-1: F1 score
	Hyperparameters    map[string]interface{} `json:"hyperparameters" db:"hyperparameters"`
	TrainingCompleted  time.Time              `json:"training_completed" db:"training_completed"`
	IsActive           bool                   `json:"is_active" db:"is_active"`
	DeployedAt         *time.Time             `json:"deployed_at" db:"deployed_at"`
	CreatedAt          time.Time              `json:"created_at" db:"created_at"`
}

// Scan implements sql.Scanner interface
func (mm *MLModelMetadata) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), mm)
}

// Value implements driver.Valuer interface
func (mm MLModelMetadata) Value() (driver.Value, error) {
	return json.Marshal(mm)
}

// PersonalizationSession tracks a user's personalization session
type PersonalizationSession struct {
	ID                 int64     `json:"id" db:"id"`
	RestaurantID       int64     `json:"restaurant_id" db:"restaurant_id"`
	UserID             int64     `json:"user_id" db:"user_id"`
	SessionID          string    `json:"session_id" db:"session_id"`
	RecommendationsShown int32   `json:"recommendations_shown" db:"recommendations_shown"`
	RecommendationsClicked int32 `json:"recommendations_clicked" db:"recommendations_clicked"`
	ClickThroughRate   float64   `json:"click_through_rate" db:"click_through_rate"`
	ConversionValue    *float64  `json:"conversion_value" db:"conversion_value"`
	StartedAt          time.Time `json:"started_at" db:"started_at"`
	EndedAt            *time.Time `json:"ended_at" db:"ended_at"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

// FeatureExtractionJob represents an async feature extraction task
type FeatureExtractionJob struct {
	ID                 int64                  `json:"id" db:"id"`
	RestaurantID       int64                  `json:"restaurant_id" db:"restaurant_id"`
	UserID             *int64                 `json:"user_id" db:"user_id"`                     // null = all users
	JobType            string                 `json:"job_type" db:"job_type"`                   // incremental, full_rebuild
	Status             string                 `json:"status" db:"status"`                       // pending, running, completed, failed
	TotalUsers         int64                  `json:"total_users" db:"total_users"`
	ProcessedUsers     int64                  `json:"processed_users" db:"processed_users"`
	FailedUsers        int64                  `json:"failed_users" db:"failed_users"`
	ErrorMessage       *string                `json:"error_message" db:"error_message"`
	Progress           float64                `json:"progress" db:"progress"`                   // 0-1
	StartedAt          time.Time              `json:"started_at" db:"started_at"`
	CompletedAt        *time.Time             `json:"completed_at" db:"completed_at"`
	CreatedAt          time.Time              `json:"created_at" db:"created_at"`
}

// RecommendationSession represents an A/B test variant for recommendations
type RecommendationAlgorithmVariant struct {
	ID                       int64                  `json:"id" db:"id"`
	RestaurantID             int64                  `json:"restaurant_id" db:"restaurant_id"`
	ExperimentID             string                 `json:"experiment_id" db:"experiment_id"`       // Link to A/B test
	VariantName              string                 `json:"variant_name" db:"variant_name"`
	AlgorithmType            string                 `json:"algorithm_type" db:"algorithm_type"`     // collaborative, content_based, hybrid
	Weight                   float64                `json:"weight" db:"weight"`                     // 0-1: percentage of traffic
	Parameters               map[string]interface{} `json:"parameters" db:"parameters"`
	IsControl                bool                   `json:"is_control" db:"is_control"`
	MetricsSnapshot          map[string]float64     `json:"metrics_snapshot" db:"metrics_snapshot"` // CTR, conversion, etc.
	StartedAt                time.Time              `json:"started_at" db:"started_at"`
	EndedAt                  *time.Time             `json:"ended_at" db:"ended_at"`
	CreatedAt                time.Time              `json:"created_at" db:"created_at"`
}

// RecommendationMetrics tracks real-time metrics for recommendations
type RecommendationMetrics struct {
	ID                      int64     `json:"id" db:"id"`
	RestaurantID            int64     `json:"restaurant_id" db:"restaurant_id"`
	PeriodStart             time.Time `json:"period_start" db:"period_start"`
	PeriodEnd               time.Time `json:"period_end" db:"period_end"`
	TotalRecommendations    int64     `json:"total_recommendations" db:"total_recommendations"`
	TotalClicks             int64     `json:"total_clicks" db:"total_clicks"`
	ClickThroughRate        float64   `json:"click_through_rate" db:"click_through_rate"`
	TotalConversions        int64     `json:"total_conversions" db:"total_conversions"`
	ConversionRate          float64   `json:"conversion_rate" db:"conversion_rate"`
	TotalConversionValue    float64   `json:"total_conversion_value" db:"total_conversion_value"`
	AverageRelevanceScore   float64   `json:"average_relevance_score" db:"average_relevance_score"`
	TopRecommendedComponents []int64  `json:"top_recommended_components" db:"top_recommended_components"`
	CreatedAt               time.Time `json:"created_at" db:"created_at"`
}

// PersonalizationConfig stores configuration for personalization engine
type PersonalizationConfig struct {
	ID                           int64   `json:"id" db:"id"`
	RestaurantID                 int64   `json:"restaurant_id" db:"restaurant_id"`
	IsEnabled                    bool    `json:"is_enabled" db:"is_enabled"`
	RecommendationsPerUser       int32   `json:"recommendations_per_user" db:"recommendations_per_user"`
	MinimumInteractionsForProfile int32   `json:"minimum_interactions_for_profile" db:"minimum_interactions_for_profile"`
	FeatureExtractionInterval    int32   `json:"feature_extraction_interval" db:"feature_extraction_interval"` // seconds
	RecommendationRefreshInterval int32  `json:"recommendation_refresh_interval" db:"recommendation_refresh_interval"` // seconds
	CollaborativeFilteringWeight float64 `json:"collaborative_filtering_weight" db:"collaborative_filtering_weight"`   // 0-1
	ContentBasedWeight           float64 `json:"content_based_weight" db:"content_based_weight"`                       // 0-1
	PopularityWeight             float64 `json:"popularity_weight" db:"popularity_weight"`                             // 0-1
	DiversificationEnabled       bool    `json:"diversification_enabled" db:"diversification_enabled"`
	ColdStartStrategy            string  `json:"cold_start_strategy" db:"cold_start_strategy"`                         // popular, random, hybrid
	CreatedAt                    time.Time `json:"created_at" db:"created_at"`
}
