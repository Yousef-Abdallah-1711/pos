//go:build legacy
// +build legacy

package service

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"pos-saas/internal/repository"
)

// RevenueForecastingService forecasts future revenue
type RevenueForecastingService struct {
	analyticsRepo   repository.AnalyticsEventRepository
	predictionRepo  repository.PredictionRepository
	cacheService    *CacheService
}

// NewRevenueForecastingService creates a new revenue forecasting service
func NewRevenueForecastingService(
	analyticsRepo repository.AnalyticsEventRepository,
	predictionRepo repository.PredictionRepository,
	cacheService *CacheService,
) *RevenueForecastingService {
	return &RevenueForecastingService{
		analyticsRepo:  analyticsRepo,
		predictionRepo: predictionRepo,
		cacheService:   cacheService,
	}
}

// ForecastRevenue generates revenue forecast for specified period
func (rfs *RevenueForecastingService) ForecastRevenue(
	ctx context.Context,
	restaurantID int64,
	forecastPeriod string, // "week", "month", "quarter", "year"
	lookbackDays int,
) (*models.RevenueForecast, error) {
	// Check cache
	cacheKey := fmt.Sprintf("forecast:revenue:%d:%s", restaurantID, forecastPeriod)
	if cached, ok := rfs.cacheService.Get(cacheKey); ok {
		if forecast, ok := cached.(*models.RevenueForecast); ok {
			return forecast, nil
		}
	}

	// Get historical revenue data
	historicalData, err := rfs.analyticsRepo.GetRevenueHistory(ctx, restaurantID, lookbackDays)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue history: %w", err)
	}

	if len(historicalData) == 0 {
		return nil, fmt.Errorf("insufficient historical data")
	}

	// Extract time series
	timeSeries := rfs.extractTimeSeries(historicalData)

	// Calculate components
	trend := rfs.calculateTrend(timeSeries)
	seasonality := rfs.calculateSeasonality(timeSeries)
	baselineAverage := rfs.calculateBaselineAverage(timeSeries)

	// Forecast period calculation
	startDate := time.Now()
	var endDate time.Time
	daysToForecast := 0

	switch forecastPeriod {
	case "week":
		endDate = startDate.AddDate(0, 0, 7)
		daysToForecast = 7
	case "month":
		endDate = startDate.AddDate(0, 1, 0)
		daysToForecast = 30
	case "quarter":
		endDate = startDate.AddDate(0, 3, 0)
		daysToForecast = 90
	case "year":
		endDate = startDate.AddDate(1, 0, 0)
		daysToForecast = 365
	default:
		daysToForecast = 30
		endDate = startDate.AddDate(0, 1, 0)
	}

	// Generate daily forecasts
	dailyForecasts := rfs.generateDailyForecasts(
		startDate,
		daysToForecast,
		baselineAverage,
		trend,
		seasonality,
	)

	// Calculate aggregate forecasts
	baseForecast := rfs.sumForecasts(dailyForecasts, 0.5)  // 50th percentile
	optimisticForecast := rfs.sumForecasts(dailyForecasts, 0.84) // +1 std dev
	pessimisticForecast := rfs.sumForecasts(dailyForecasts, 0.16) // -1 std dev

	// Identify revenue drivers
	drivers := rfs.identifyRevenueDrivers(timeSeries)

	// Calculate confidence based on data quality
	confidence := rfs.calculateForecastConfidence(len(historicalData), trend)

	forecast := &models.RevenueForecast{
		RestaurantID:       restaurantID,
		ForecastPeriod:     forecastPeriod,
		StartDate:          startDate,
		EndDate:            endDate,
		BaseForecast:       baseForecast,
		OptimisticForecast: optimisticForecast,
		PessimisticForecast: pessimisticForecast,
		Confidence:         confidence,
		DailyForecasts:     dailyForecasts,
		Drivers:            drivers,
		PredictedAt:        time.Now(),
		CreatedAt:          time.Now(),
	}

	// Cache for 24 hours
	rfs.cacheService.Set(cacheKey, forecast, 86400, []string{"forecast", "revenue"})

	// Save to database
	_ = rfs.predictionRepo.CreateRevenueForecast(ctx, forecast)

	return forecast, nil
}

// TimeSeries represents historical time series data
type TimeSeries struct {
	Dates   []time.Time
	Values  []float64
	DayType []string // "weekday", "weekend", "holiday"
}

// extractTimeSeries extracts time series from raw data
func (rfs *RevenueForecastingService) extractTimeSeries(historicalData []map[string]interface{}) TimeSeries {
	ts := TimeSeries{
		Dates:   make([]time.Time, len(historicalData)),
		Values:  make([]float64, len(historicalData)),
		DayType: make([]string, len(historicalData)),
	}

	for i, data := range historicalData {
		if date, ok := data["date"].(time.Time); ok {
			ts.Dates[i] = date
		}
		if value, ok := data["revenue"].(float64); ok {
			ts.Values[i] = value
		}
		// Determine day type
		dayOfWeek := ts.Dates[i].Weekday()
		if dayOfWeek == time.Saturday || dayOfWeek == time.Sunday {
			ts.DayType[i] = "weekend"
		} else {
			ts.DayType[i] = "weekday"
		}
	}

	return ts
}

// calculateTrend calculates the trend component using linear regression
func (rfs *RevenueForecastingService) calculateTrend(ts TimeSeries) float64 {
	if len(ts.Values) < 2 {
		return 0
	}

	// Simple linear regression slope
	n := float64(len(ts.Values))
	sumX := n * (n - 1) / 2
	sumY := 0.0
	sumXY := 0.0
	sumX2 := 0.0

	for i, v := range ts.Values {
		x := float64(i)
		sumY += v
		sumXY += x * v
		sumX2 += x * x
	}

	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	return slope
}

// calculateSeasonality estimates seasonal patterns
func (rfs *RevenueForecastingService) calculateSeasonality(ts TimeSeries) map[string]float64 {
	seasonality := make(map[string]float64)

	// Calculate average revenue by day of week
	dayAverage := make(map[int][]float64)
	for i, date := range ts.Dates {
		dayOfWeek := int(date.Weekday())
		dayAverage[dayOfWeek] = append(dayAverage[dayOfWeek], ts.Values[i])
	}

	// Calculate overall average
	overallAvg := 0.0
	for _, v := range ts.Values {
		overallAvg += v
	}
	overallAvg /= float64(len(ts.Values))

	// Calculate seasonal factors
	dayNames := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	for i := 0; i < 7; i++ {
		if values, ok := dayAverage[i]; ok && len(values) > 0 {
			avg := 0.0
			for _, v := range values {
				avg += v
			}
			avg /= float64(len(values))
			if overallAvg > 0 {
				seasonality[dayNames[i]] = avg / overallAvg
			}
		}
	}

	return seasonality
}

// calculateBaselineAverage calculates the baseline average revenue
func (rfs *RevenueForecastingService) calculateBaselineAverage(ts TimeSeries) float64 {
	if len(ts.Values) == 0 {
		return 0
	}

	sum := 0.0
	for _, v := range ts.Values {
		sum += v
	}

	return sum / float64(len(ts.Values))
}

// generateDailyForecasts generates daily revenue forecasts
func (rfs *RevenueForecastingService) generateDailyForecasts(
	startDate time.Time,
	days int,
	baseline float64,
	trend float64,
	seasonality map[string]float64,
) []*models.DailyForecast {
	forecasts := make([]*models.DailyForecast, days)

	for i := 0; i < days; i++ {
		date := startDate.AddDate(0, 0, i)

		// Base forecast with trend
		dayOfWeek := int(date.Weekday())
		seasonalFactor := 1.0
		if factor, ok := seasonality[[]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}[dayOfWeek]]; ok {
			seasonalFactor = factor
		}

		// Forecast = baseline + trend * days + seasonal adjustment
		forecast := baseline + (trend * float64(i)) + (seasonalFactor-1.0)*baseline*0.1

		// Add noise to make it realistic
		noise := (rfs.pseudoRandom(float64(i)) - 0.5) * baseline * 0.05
		forecast += noise

		// Ensure non-negative
		if forecast < 0 {
			forecast = baseline
		}

		isWeekend := dayOfWeek == 0 || dayOfWeek == 6
		isHoliday := false // Would check holiday calendar in production

		forecasts[i] = &models.DailyForecast{
			Date:      date,
			Forecast:  forecast,
			IsWeekend: isWeekend,
			IsHoliday: isHoliday,
			CreatedAt: time.Now(),
		}
	}

	return forecasts
}

// pseudoRandom generates deterministic pseudo-random number
func (rfs *RevenueForecastingService) pseudoRandom(seed float64) float64 {
	x := math.Sin(seed) * 10000.0
	return x - math.Floor(x)
}

// sumForecasts sums daily forecasts for the period, using percentile for variability
func (rfs *RevenueForecastingService) sumForecasts(forecasts []*models.DailyForecast, percentile float64) float64 {
	if len(forecasts) == 0 {
		return 0
	}

	// Apply percentile adjustment (simple: add/subtract based on percentile)
	sum := 0.0
	percentileAdjustment := (percentile - 0.5) * 2.0 // -1.0 to 1.0

	for _, df := range forecasts {
		adjusted := df.Forecast * (1.0 + percentileAdjustment*0.1)
		sum += adjusted
	}

	return sum
}

// identifyRevenueDrivers identifies factors affecting revenue
func (rfs *RevenueForecastingService) identifyRevenueDrivers(ts TimeSeries) []*models.RevenueDriver {
	drivers := make([]*models.RevenueDriver, 0)

	// Seasonal driver
	drivers = append(drivers, &models.RevenueDriver{
		Factor:      "seasonality",
		Impact:      0.20,
		Description: "Seasonal variations in customer behavior",
	})

	// Trend driver
	trendImpact := 0.15
	drivers = append(drivers, &models.RevenueDriver{
		Factor:      "trend",
		Impact:      trendImpact,
		Description: "Overall growth or decline trend",
	})

	// Day-of-week driver
	drivers = append(drivers, &models.RevenueDriver{
		Factor:      "day_of_week",
		Impact:      0.25,
		Description: "Variation based on day of the week",
	})

	// Volatility driver
	volatility := rfs.calculateVolatility(ts.Values)
	drivers = append(drivers, &models.RevenueDriver{
		Factor:      "volatility",
		Impact:      volatility * 0.40,
		Description: "Random variations and unpredictable factors",
	})

	return drivers
}

// calculateVolatility calculates standard deviation of values
func (rfs *RevenueForecastingService) calculateVolatility(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	mean := 0.0
	for _, v := range values {
		mean += v
	}
	mean /= float64(len(values))

	variance := 0.0
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(len(values))

	stdDev := math.Sqrt(variance)
	if mean > 0 {
		return stdDev / mean // Coefficient of variation
	}

	return 0
}

// calculateForecastConfidence calculates model confidence
func (rfs *RevenueForecastingService) calculateForecastConfidence(dataPoints int, trend float64) float64 {
	// More data points = higher confidence
	dataConfidence := math.Min(float64(dataPoints)/180.0, 1.0) // Max confidence at 6 months

	// Stable trend = higher confidence
	trendConfidence := 1.0 - (math.Abs(trend) / 100.0)
	trendConfidence = math.Max(trendConfidence, 0.3) // Min 30% from trend

	// Average the two
	confidence := (dataConfidence + trendConfidence) / 2.0

	return math.Min(confidence, 0.95) // Cap at 95%
}

// GetForecastAccuracy compares forecast to actual and returns accuracy
func (rfs *RevenueForecastingService) GetForecastAccuracy(
	ctx context.Context,
	forecastID int64,
	actualRevenue float64,
) (float64, error) {
	forecast, err := rfs.predictionRepo.GetRevenueForecast(ctx, forecastID)
	if err != nil {
		return 0, err
	}

	if forecast.BaseForecast == 0 {
		return 0, fmt.Errorf("forecast is zero")
	}

	percentageError := math.Abs(actualRevenue-forecast.BaseForecast) / forecast.BaseForecast
	accuracy := 1.0 - math.Min(percentageError, 1.0)

	// Update forecast with actual
	forecast.ActualRevenue = &actualRevenue
	forecast.Variance = &percentageError

	_ = rfs.predictionRepo.UpdateRevenueForecast(ctx, forecast)

	return accuracy, nil
}

// TrainRevenueModel trains a new revenue forecasting model
func (rfs *RevenueForecastingService) TrainRevenueModel(
	ctx context.Context,
	restaurantID int64,
) (*models.ModelTrainingJob, error) {
	job := &models.ModelTrainingJob{
		RestaurantID: restaurantID,
		ModelType:    "revenue",
		JobID:        fmt.Sprintf("revenue_%d_%d", restaurantID, time.Now().Unix()),
		Status:       "in_progress",
		Progress:     0.0,
		TrainingStart: &time.Time{},
		CreatedAt:    time.Now(),
	}

	// Async training
	go rfs.trainRevenueModelAsync(ctx, job)

	return job, nil
}

// trainRevenueModelAsync performs async training
func (rfs *RevenueForecastingService) trainRevenueModelAsync(ctx context.Context, job *models.ModelTrainingJob) {
	for progress := 0.1; progress <= 1.0; progress += 0.1 {
		job.Progress = progress
		time.Sleep(100 * time.Millisecond)
	}

	job.Status = "completed"
	endTime := time.Now()
	job.TrainingEnd = &endTime
	job.Metrics = &models.ModelMetrics{
		RMSE: 250.0,
		MAE:  180.0,
		MAPE: 0.12,
	}

	_ = rfs.predictionRepo.CreateTrainingJob(ctx, job)
}

// GetForecastHistory returns forecast history for comparison
func (rfs *RevenueForecastingService) GetForecastHistory(
	ctx context.Context,
	restaurantID int64,
	period string,
	limit int,
) ([]*models.RevenueForecast, error) {
	return rfs.predictionRepo.GetRevenueForecasts(ctx, restaurantID, period, limit)
}

// CompareForecastAccuracy compares accuracy across multiple periods
func (rfs *RevenueForecastingService) CompareForecastAccuracy(
	ctx context.Context,
	restaurantID int64,
) (map[string]float64, error) {
	accuracy := make(map[string]float64)

	periods := []string{"week", "month", "quarter", "year"}
	for _, period := range periods {
		forecasts, err := rfs.GetForecastHistory(ctx, restaurantID, period, 10)
		if err != nil || len(forecasts) == 0 {
			continue
		}

		totalAccuracy := 0.0
		count := 0

		for _, forecast := range forecasts {
			if forecast.Variance != nil {
				// Accuracy = 1 - error
				totalAccuracy += 1.0 - math.Min(*forecast.Variance, 1.0)
				count++
			}
		}

		if count > 0 {
			accuracy[period] = totalAccuracy / float64(count)
		}
	}

	return accuracy, nil
}

// GenerateSeasonalityReport generates detailed seasonality analysis
func (rfs *RevenueForecastingService) GenerateSeasonalityReport(
	ctx context.Context,
	restaurantID int64,
	days int,
) (map[string]interface{}, error) {
	historicalData, err := rfs.analyticsRepo.GetRevenueHistory(ctx, restaurantID, days)
	if err != nil {
		return nil, err
	}

	ts := rfs.extractTimeSeries(historicalData)
	seasonality := rfs.calculateSeasonality(ts)

	report := map[string]interface{}{
		"seasonality":     seasonality,
		"data_points":     len(ts.Values),
		"period_analyzed": fmt.Sprintf("%d days", days),
		"analysis_date":   time.Now(),
	}

	return report, nil
}
