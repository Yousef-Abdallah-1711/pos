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

// AnomalyDetectionService detects unusual patterns in metrics
type AnomalyDetectionService struct {
	analyticsRepo   repository.AnalyticsEventRepository
	predictionRepo  repository.PredictionRepository
	cacheService    *CacheService
}

// NewAnomalyDetectionService creates a new anomaly detection service
func NewAnomalyDetectionService(
	analyticsRepo repository.AnalyticsEventRepository,
	predictionRepo repository.PredictionRepository,
	cacheService *CacheService,
) *AnomalyDetectionService {
	return &AnomalyDetectionService{
		analyticsRepo:  analyticsRepo,
		predictionRepo: predictionRepo,
		cacheService:   cacheService,
	}
}

// DetectAnomalies detects anomalies in restaurant metrics
func (ads *AnomalyDetectionService) DetectAnomalies(
	ctx context.Context,
	restaurantID int64,
	metricType string, // "revenue", "traffic", "conversion", "engagement", "order_value"
	lookbackDays int,
	stdDevThreshold float64, // Number of standard deviations for anomaly
) ([]*models.AnomalyDetection, error) {
	// Get baseline metrics
	baseline, err := ads.getMetricBaseline(ctx, restaurantID, metricType, lookbackDays)
	if err != nil {
		return nil, fmt.Errorf("failed to get baseline: %w", err)
	}

	// Get current metrics
	currentMetrics, err := ads.analyticsRepo.GetMetrics(ctx, restaurantID, metricType, 7)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	anomalies := make([]*models.AnomalyDetection, 0)

	// Detect anomalies for each time point
	for _, metric := range currentMetrics {
		detectedValue := metric["value"].(float64)
		period := metric["period"].(string)

		// Check if metric is anomalous
		anomalyScore := ads.calculateAnomalyScore(
			detectedValue,
			baseline.Mean,
			baseline.StdDev,
			stdDevThreshold,
		)

		if anomalyScore > 0.5 { // Above threshold
			anomaly := &models.AnomalyDetection{
				RestaurantID:   restaurantID,
				MetricType:     metricType,
				AnomalyScore:   anomalyScore,
				DetectedValue:  detectedValue,
				ExpectedValue:  baseline.Mean,
				Deviation:      ((detectedValue - baseline.Mean) / baseline.Mean) * 100,
				AnomalyType:    ads.classifyAnomalyType(detectedValue, baseline),
				Severity:       ads.determineSeverity(anomalyScore, abs(((detectedValue-baseline.Mean)/baseline.Mean)*100)),
				Description:    ads.generateAnomalyDescription(metricType, detectedValue, baseline.Mean),
				AffectedPeriod: period,
				ActionRequired: anomalyScore > 0.75,
				Status:         "new",
				RootCauseHypotheses: ads.generateRootCauseHypotheses(metricType, detectedValue, baseline.Mean),
				DetectedAt:     time.Now(),
				CreatedAt:      time.Now(),
			}

			anomalies = append(anomalies, anomaly)
		}
	}

	// Save anomalies to database
	for _, anomaly := range anomalies {
		_ = ads.predictionRepo.CreateAnomaly(ctx, anomaly)
	}

	return anomalies, nil
}

// DetectUserAnomalies detects unusual user behavior
func (ads *AnomalyDetectionService) DetectUserAnomalies(
	ctx context.Context,
	restaurantID int64,
	userID int64,
) ([]*models.AnomalyDetection, error) {
	// Get user metrics
	userMetrics, err := ads.analyticsRepo.GetUserMetrics(ctx, restaurantID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user metrics: %w", err)
	}

	// Get baseline for comparison
	baseline, err := ads.getMetricBaseline(ctx, restaurantID, "engagement", 30)
	if err != nil {
		return nil, fmt.Errorf("failed to get baseline: %w", err)
	}

	anomalies := make([]*models.AnomalyDetection, 0)

	// Check if user engagement is anomalous
	userEngagement := userMetrics["engagement"].(float64)
	if userEngagement < baseline.Mean-2*baseline.StdDev {
		anomaly := &models.AnomalyDetection{
			RestaurantID:   restaurantID,
			UserID:         &userID,
			MetricType:     "engagement",
			AnomalyScore:   0.8,
			DetectedValue:  userEngagement,
			ExpectedValue:  baseline.Mean,
			Deviation:      ((userEngagement - baseline.Mean) / baseline.Mean) * 100,
			AnomalyType:    "drop",
			Severity:       "warning",
			Description:    fmt.Sprintf("User engagement dropped to %.2f (expected %.2f)", userEngagement, baseline.Mean),
			ActionRequired: true,
			Status:         "new",
			RootCauseHypotheses: []string{"user_churn_risk", "dissatisfaction", "changed_usage_pattern"},
			DetectedAt:     time.Now(),
			CreatedAt:      time.Now(),
		}

		anomalies = append(anomalies, anomaly)
	}

	return anomalies, nil
}

// calculateAnomalyScore calculates z-score based anomaly detection
func (ads *AnomalyDetectionService) calculateAnomalyScore(
	value float64,
	mean float64,
	stdDev float64,
	threshold float64,
) float64 {
	if stdDev == 0 {
		return 0
	}

	// Z-score calculation
	zScore := math.Abs((value - mean) / stdDev)

	// Convert to anomaly score (0-1)
	// If z-score equals threshold, anomaly score is 0.5
	// If z-score is 0, anomaly score is 0
	// If z-score >> threshold, anomaly score approaches 1.0
	anomalyScore := math.Atan(zScore-threshold) / math.Pi + 0.5

	// Clamp to 0-1
	if anomalyScore < 0 {
		anomalyScore = 0
	}
	if anomalyScore > 1 {
		anomalyScore = 1
	}

	return anomalyScore
}

// classifyAnomalyType classifies the type of anomaly detected
func (ads *AnomalyDetectionService) classifyAnomalyType(value float64, baseline *models.MetricBaseline) string {
	deviation := value - baseline.Mean

	if deviation > 2*baseline.StdDev {
		return "spike"
	}

	if deviation < -2*baseline.StdDev {
		return "drop"
	}

	// Check if it's outside normal range but within 2 std dev
	if value > baseline.P95 {
		return "spike"
	}

	if value < baseline.P5 {
		return "drop"
	}

	// Check for pattern shift (simplified)
	if math.Abs(deviation) > baseline.Mean*0.5 {
		return "trend_shift"
	}

	return "pattern_break"
}

// determineSeverity determines severity level of anomaly
func (ads *AnomalyDetectionService) determineSeverity(anomalyScore float64, percentageDeviation float64) string {
	absDeviation := math.Abs(percentageDeviation)

	if anomalyScore > 0.85 && absDeviation > 50 {
		return "critical"
	}

	if anomalyScore > 0.75 && absDeviation > 30 {
		return "warning"
	}

	if anomalyScore > 0.6 && absDeviation > 15 {
		return "info"
	}

	return "info"
}

// generateAnomalyDescription generates human-readable description
func (ads *AnomalyDetectionService) generateAnomalyDescription(
	metricType string,
	detectedValue float64,
	expectedValue float64,
) string {
	deviation := ((detectedValue - expectedValue) / expectedValue) * 100

	metricName := metricType
	direction := "decreased"

	if deviation > 0 {
		direction = "increased"
	}

	return fmt.Sprintf(
		"%s %s by %.1f%% (detected: %.2f, expected: %.2f)",
		metricName,
		direction,
		math.Abs(deviation),
		detectedValue,
		expectedValue,
	)
}

// generateRootCauseHypotheses generates possible root causes
func (ads *AnomalyDetectionService) generateRootCauseHypotheses(
	metricType string,
	detectedValue float64,
	expectedValue float64,
) []string {
	hypotheses := []string{}

	anomalouslyHigh := detectedValue > expectedValue*1.2
	anomalouslyLow := detectedValue < expectedValue*0.8

	switch metricType {
	case "revenue":
		if anomalouslyHigh {
			hypotheses = append(hypotheses, "special_promotion_active")
			hypotheses = append(hypotheses, "seasonal_spike")
			hypotheses = append(hypotheses, "marketing_campaign_success")
		} else if anomalouslyLow {
			hypotheses = append(hypotheses, "technical_issue")
			hypotheses = append(hypotheses, "competitor_activity")
			hypotheses = append(hypotheses, "seasonal_decline")
		}

	case "traffic":
		if anomalouslyHigh {
			hypotheses = append(hypotheses, "viral_content")
			hypotheses = append(hypotheses, "marketing_activity")
			hypotheses = append(hypotheses, "network_issue_resolved")
		} else if anomalouslyLow {
			hypotheses = append(hypotheses, "system_outage")
			hypotheses = append(hypotheses, "dns_issue")
			hypotheses = append(hypotheses, "traffic_redirected")
		}

	case "conversion":
		if anomalouslyHigh {
			hypotheses = append(hypotheses, "improved_ui")
			hypotheses = append(hypotheses, "successful_ab_test")
		} else if anomalouslyLow {
			hypotheses = append(hypotheses, "checkout_issue")
			hypotheses = append(hypotheses, "payment_processor_issue")
			hypotheses = append(hypotheses, "user_experience_degradation")
		}

	case "engagement":
		if anomalouslyLow {
			hypotheses = append(hypotheses, "user_churn")
			hypotheses = append(hypotheses, "feature_removed")
			hypotheses = append(hypotheses, "competitor_launch")
		}
	}

	return hypotheses
}

// getMetricBaseline gets or calculates baseline for a metric
func (ads *AnomalyDetectionService) getMetricBaseline(
	ctx context.Context,
	restaurantID int64,
	metricType string,
	lookbackDays int,
) (*models.MetricBaseline, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("baseline:%d:%s:%d", restaurantID, metricType, lookbackDays)
	if cached, ok := ads.cacheService.Get(cacheKey); ok {
		if baseline, ok := cached.(*models.MetricBaseline); ok {
			return baseline, nil
		}
	}

	// Get historical data
	historicalMetrics, err := ads.analyticsRepo.GetMetrics(ctx, restaurantID, metricType, lookbackDays)
	if err != nil {
		return nil, err
	}

	// Calculate statistics
	values := make([]float64, len(historicalMetrics))
	for i, metric := range historicalMetrics {
		values[i] = metric["value"].(float64)
	}

	baseline := ads.calculateBaselineStatistics(values, restaurantID, metricType)

	// Cache for 7 days
	ads.cacheService.Set(cacheKey, baseline, 604800, []string{"baseline", "metric"})

	// Save to database
	_ = ads.predictionRepo.CreateMetricBaseline(ctx, baseline)

	return baseline, nil
}

// calculateBaselineStatistics calculates statistical baseline
func (ads *AnomalyDetectionService) calculateBaselineStatistics(
	values []float64,
	restaurantID int64,
	metricType string,
) *models.MetricBaseline {
	if len(values) == 0 {
		return &models.MetricBaseline{}
	}

	// Sort for percentile calculation
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	// Calculate mean
	mean := 0.0
	for _, v := range values {
		mean += v
	}
	mean /= float64(len(values))

	// Calculate std dev
	variance := 0.0
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(len(values))
	stdDev := math.Sqrt(variance)

	// Get percentiles
	p5Index := int(float64(len(sorted)) * 0.05)
	p95Index := int(float64(len(sorted)) * 0.95)

	p5 := sorted[p5Index]
	p95 := sorted[p95Index]
	median := sorted[len(sorted)/2]
	minVal := sorted[0]
	maxVal := sorted[len(sorted)-1]

	return &models.MetricBaseline{
		RestaurantID: restaurantID,
		MetricType:   metricType,
		Mean:         mean,
		StdDev:       stdDev,
		Median:       median,
		P5:           p5,
		P95:          p95,
		Min:          minVal,
		Max:          maxVal,
		DataPoints:   len(values),
		CalculatedAt: time.Now(),
		CreatedAt:    time.Now(),
	}
}

// AcknowledgeAnomaly marks an anomaly as acknowledged
func (ads *AnomalyDetectionService) AcknowledgeAnomaly(
	ctx context.Context,
	anomalyID int64,
) error {
	anomaly, err := ads.predictionRepo.GetAnomaly(ctx, anomalyID)
	if err != nil {
		return err
	}

	anomaly.Status = "acknowledged"
	now := time.Now()
	anomaly.Investigation = &models.Investigation{
		Status:         "pending",
		InvestigatedAt: &now,
	}

	return ads.predictionRepo.UpdateAnomaly(ctx, anomaly)
}

// InvestigateAnomaly investigates and resolves anomaly
func (ads *AnomalyDetectionService) InvestigateAnomaly(
	ctx context.Context,
	anomalyID int64,
	investigatedBy int64,
	rootCause string,
	resolution string,
	notes string,
) error {
	anomaly, err := ads.predictionRepo.GetAnomaly(ctx, anomalyID)
	if err != nil {
		return err
	}

	now := time.Now()
	anomaly.Status = "investigated"
	anomaly.Investigation = &models.Investigation{
		Status:         "completed",
		InvestigatedBy: investigatedBy,
		InvestigatedAt: &now,
		RootCause:      rootCause,
		Resolution:     resolution,
		Notes:          notes,
	}

	return ads.predictionRepo.UpdateAnomaly(ctx, anomaly)
}

// GetRecentAnomalies returns recent anomalies
func (ads *AnomalyDetectionService) GetRecentAnomalies(
	ctx context.Context,
	restaurantID int64,
	limit int,
	severity string, // empty for all, or "critical", "warning", "info"
) ([]*models.AnomalyDetection, error) {
	anomalies, err := ads.predictionRepo.GetAnomalies(ctx, restaurantID, limit)
	if err != nil {
		return nil, err
	}

	// Filter by severity if specified
	if severity != "" {
		filtered := make([]*models.AnomalyDetection, 0)
		for _, anomaly := range anomalies {
			if anomaly.Severity == severity {
				filtered = append(filtered, anomaly)
			}
		}
		return filtered, nil
	}

	return anomalies, nil
}

// GetAnomalyStats returns statistics about anomalies
func (ads *AnomalyDetectionService) GetAnomalyStats(
	ctx context.Context,
	restaurantID int64,
	days int,
) (map[string]interface{}, error) {
	anomalies, err := ads.predictionRepo.GetAnomaliesSince(ctx, restaurantID, time.Now().AddDate(0, 0, -days))
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_anomalies":    len(anomalies),
		"by_severity":        make(map[string]int),
		"by_type":            make(map[string]int),
		"by_metric":          make(map[string]int),
		"unresolved_count":   0,
		"action_required":    0,
	}

	bySeverity := stats["by_severity"].(map[string]int)
	byType := stats["by_type"].(map[string]int)
	byMetric := stats["by_metric"].(map[string]int)

	for _, anomaly := range anomalies {
		bySeverity[anomaly.Severity]++
		byType[anomaly.AnomalyType]++
		byMetric[anomaly.MetricType]++

		if anomaly.Status == "new" {
			stats["unresolved_count"] = stats["unresolved_count"].(int) + 1
		}

		if anomaly.ActionRequired {
			stats["action_required"] = stats["action_required"].(int) + 1
		}
	}

	return stats, nil
}

// abs returns absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
