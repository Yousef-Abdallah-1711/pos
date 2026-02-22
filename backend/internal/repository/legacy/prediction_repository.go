//go:build legacy
// +build legacy

package repository

import (
	"context"
	"time"
)

// PredictionRepository handles prediction model data persistence
type PredictionRepository interface {
	// Churn Prediction
	CreateChurnPrediction(ctx context.Context, prediction *models.ChurnPrediction) error
	GetChurnPrediction(ctx context.Context, id int64) (*models.ChurnPrediction, error)
	GetChurnPredictions(ctx context.Context, restaurantID int64, userID *int64, limit int) ([]*models.ChurnPrediction, error)
	UpdateChurnPrediction(ctx context.Context, prediction *models.ChurnPrediction) error

	// Revenue Forecasting
	CreateRevenueForecast(ctx context.Context, forecast *models.RevenueForecast) error
	GetRevenueForecast(ctx context.Context, id int64) (*models.RevenueForecast, error)
	GetRevenueForecasts(ctx context.Context, restaurantID int64, period string, limit int) ([]*models.RevenueForecast, error)
	UpdateRevenueForecast(ctx context.Context, forecast *models.RevenueForecast) error

	// Daily Forecasts
	CreateDailyForecast(ctx context.Context, forecast *models.DailyForecast) error
	GetDailyForecasts(ctx context.Context, forecastID int64) ([]*models.DailyForecast, error)

	// Lifetime Value
	CreateLTVPrediction(ctx context.Context, prediction *models.LifetimeValuePrediction) error
	GetLTVPrediction(ctx context.Context, id int64) (*models.LifetimeValuePrediction, error)
	GetLTVPredictions(ctx context.Context, restaurantID int64, limit int) ([]*models.LifetimeValuePrediction, error)
	GetLTVPredictionHistory(ctx context.Context, restaurantID int64, userID int64, limit int) ([]*models.LifetimeValuePrediction, error)
	UpdateLTVPrediction(ctx context.Context, prediction *models.LifetimeValuePrediction) error

	// Anomaly Detection
	CreateAnomaly(ctx context.Context, anomaly *models.AnomalyDetection) error
	GetAnomaly(ctx context.Context, id int64) (*models.AnomalyDetection, error)
	GetAnomalies(ctx context.Context, restaurantID int64, limit int) ([]*models.AnomalyDetection, error)
	GetAnomaliesSince(ctx context.Context, restaurantID int64, since time.Time) ([]*models.AnomalyDetection, error)
	UpdateAnomaly(ctx context.Context, anomaly *models.AnomalyDetection) error

	// Metric Baselines
	CreateMetricBaseline(ctx context.Context, baseline *models.MetricBaseline) error
	GetMetricBaseline(ctx context.Context, restaurantID int64, metricType string) (*models.MetricBaseline, error)
	UpdateMetricBaseline(ctx context.Context, baseline *models.MetricBaseline) error

	// Model Training Jobs
	CreateTrainingJob(ctx context.Context, job *models.ModelTrainingJob) error
	GetTrainingJob(ctx context.Context, jobID string) (*models.ModelTrainingJob, error)
	UpdateTrainingJob(ctx context.Context, job *models.ModelTrainingJob) error

	// Insights
	CreateInsight(ctx context.Context, insight *models.PredictionInsight) error
	GetInsights(ctx context.Context, restaurantID int64, status string, limit int) ([]*models.PredictionInsight, error)
	UpdateInsight(ctx context.Context, insight *models.PredictionInsight) error

	// Configuration
	GetPredictionConfig(ctx context.Context, restaurantID int64) (*models.PredictionConfig, error)
	UpdatePredictionConfig(ctx context.Context, config *models.PredictionConfig) error
	CreatePredictionConfig(ctx context.Context, config *models.PredictionConfig) error

	// Exports
	CreateExport(ctx context.Context, export *models.PredictionExport) error
	GetExport(ctx context.Context, id int64) (*models.PredictionExport, error)
	ListExports(ctx context.Context, restaurantID int64, limit int) ([]*models.PredictionExport, error)
	UpdateExport(ctx context.Context, export *models.PredictionExport) error
}
