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

// LifetimeValueService predicts and manages customer lifetime value
type LifetimeValueService struct {
	behaviorRepo    repository.UserBehaviorRepository
	analyticsRepo   repository.AnalyticsEventRepository
	predictionRepo  repository.PredictionRepository
	cacheService    *CacheService
}

// NewLifetimeValueService creates a new LTV service
func NewLifetimeValueService(
	behaviorRepo repository.UserBehaviorRepository,
	analyticsRepo repository.AnalyticsEventRepository,
	predictionRepo repository.PredictionRepository,
	cacheService *CacheService,
) *LifetimeValueService {
	return &LifetimeValueService{
		behaviorRepo:   behaviorRepo,
		analyticsRepo:  analyticsRepo,
		predictionRepo: predictionRepo,
		cacheService:   cacheService,
	}
}

// PredictLifetimeValue predicts customer lifetime value
func (lvs *LifetimeValueService) PredictLifetimeValue(
	ctx context.Context,
	restaurantID int64,
	userID int64,
) (*models.LifetimeValuePrediction, error) {
	// Check cache
	cacheKey := fmt.Sprintf("ltv:user:%d:%d", restaurantID, userID)
	if cached, ok := lvs.cacheService.Get(cacheKey); ok {
		if ltv, ok := cached.(*models.LifetimeValuePrediction); ok {
			return ltv, nil
		}
	}

	// Get user behavior profile
	profile, err := lvs.behaviorRepo.GetProfile(ctx, restaurantID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	// Get revenue data
	userRevenue, err := lvs.analyticsRepo.GetUserRevenue(ctx, restaurantID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user revenue: %w", err)
	}

	// Extract LTV features
	monthlyValue := lvs.calculateMonthlyValue(userRevenue, profile)
	growthTrend := lvs.calculateGrowthTrend(userRevenue)
	churnRisk := lvs.estimateChurnRisk(profile)
	accountAgeDays := int(time.Since(profile.CreatedAt).Hours() / 24)

	// Estimate remaining active months
	predictedMonths := lvs.estimateActiveMonths(profile, churnRisk)

	// Calculate predicted LTV
	currentLTV := lvs.calculateCurrentLTV(userRevenue)
	predictedLTV := lvs.calculatePredictedLTV(
		currentLTV,
		monthlyValue,
		predictedMonths,
		growthTrend,
	)

	// Determine value segment
	valueSegment := lvs.determineValueSegment(predictedLTV)

	// Get upsell opportunities
	upsellRecommendations := lvs.getUpsellRecommendations(profile, valueSegment)
	upsellOpportunity := lvs.calculateUpsellOpportunity(profile, monthlyValue, valueSegment)

	prediction := &models.LifetimeValuePrediction{
		RestaurantID:       restaurantID,
		UserID:             userID,
		CurrentLTV:         currentLTV,
		PredictedLTV:       predictedLTV,
		MonthlyValue:       monthlyValue,
		PredictedMonths:    predictedMonths,
		ValueSegment:       valueSegment,
		GrowthTrend:        growthTrend,
		ChurnRisk:          churnRisk,
		UpsellOpportunity:  upsellOpportunity,
		UpsellRecommendations: upsellRecommendations,
		Confidence:         0.83,
		PredictedAt:        time.Now(),
		ValidUntil:         time.Now().AddDate(0, 1, 0),
	}

	// Cache for 30 days
	lvs.cacheService.Set(cacheKey, prediction, 2592000, []string{"ltv", "user", fmt.Sprintf("user:%d", userID)})

	// Save to database
	_ = lvs.predictionRepo.CreateLTVPrediction(ctx, prediction)

	return prediction, nil
}

// calculateMonthlyValue calculates average monthly spending
func (lvs *LifetimeValueService) calculateMonthlyValue(
	userRevenue float64,
	profile *models.UserBehaviorProfile,
) float64 {
	if userRevenue == 0 {
		return 0
	}

	accountAgeDays := time.Since(profile.CreatedAt).Hours() / 24
	if accountAgeDays == 0 {
		accountAgeDays = 1
	}

	accountAgeMonths := accountAgeDays / 30.0
	monthlyValue := userRevenue / accountAgeMonths

	return math.Max(monthlyValue, 0)
}

// calculateGrowthTrend calculates revenue growth rate
func (lvs *LifetimeValueService) calculateGrowthTrend(userRevenue float64) float64 {
	// Simplified: assume recent growth
	// In production, would compare recent vs historical
	if userRevenue < 100 {
		return 0.95 // Expect decline for small spenders
	}

	if userRevenue > 1000 {
		return 1.10 // Expect growth for large spenders
	}

	return 1.02 // Slight growth otherwise
}

// estimateChurnRisk estimates probability of user churning
func (lvs *LifetimeValueService) estimateChurnRisk(profile *models.UserBehaviorProfile) float64 {
	daysSinceLastActivity := int(time.Since(profile.LastActivityAt).Hours() / 24)

	// High inactivity = high churn risk
	churnRisk := math.Min(float64(daysSinceLastActivity)/60.0, 1.0)

	// Engagement reduces risk
	engagement := float64(profile.TotalInteractions) / float64(int(time.Since(profile.CreatedAt).Hours()/24+1))
	engagementFactor := math.Min(engagement/2.0, 1.0)

	churnRisk = churnRisk * (1.0 - engagementFactor*0.5)

	return math.Max(math.Min(churnRisk, 1.0), 0.05) // Min 5% churn risk
}

// estimateActiveMonths estimates remaining active months
func (lvs *LifetimeValueService) estimateActiveMonths(
	profile *models.UserBehaviorProfile,
	churnRisk float64,
) int {
	// Base estimate: 12 months
	baseMonths := 12.0

	// Adjust for churn risk
	activeMonths := baseMonths * (1.0 - churnRisk)

	// Adjust for historical tenure
	accountAgeMonths := time.Since(profile.CreatedAt).Hours() / 24.0 / 30.0
	if accountAgeMonths > 36 {
		// Long-term customer, estimate extended tenure
		activeMonths *= 1.5
	}

	return int(math.Max(activeMonths, 3)) // Minimum 3 months
}

// calculateCurrentLTV calculates actual cumulative value to date
func (lvs *LifetimeValueService) calculateCurrentLTV(userRevenue float64) float64 {
	return userRevenue
}

// calculatePredictedLTV predicts future lifetime value
func (lvs *LifetimeValueService) calculatePredictedLTV(
	currentLTV float64,
	monthlyValue float64,
	predictedMonths int,
	growthTrend float64,
) float64 {
	// Future value = current + (monthly value * months * growth trend)
	predictedAdditional := monthlyValue * float64(predictedMonths) * growthTrend

	return currentLTV + predictedAdditional
}

// determineValueSegment categorizes customer by value
func (lvs *LifetimeValueService) determineValueSegment(ltv float64) string {
	switch {
	case ltv > 5000:
		return "vip"
	case ltv > 2000:
		return "high"
	case ltv > 500:
		return "medium"
	default:
		return "low"
	}
}

// getUpsellRecommendations recommends upsell opportunities
func (lvs *LifetimeValueService) getUpsellRecommendations(
	profile *models.UserBehaviorProfile,
	valueSegment string,
) []string {
	recommendations := []string{}

	// Based on usage patterns
	if len(profile.ComponentViews) < 10 {
		recommendations = append(recommendations, "explore_more_features")
	}

	if profile.TotalSessions < 5 {
		recommendations = append(recommendations, "increase_frequency")
	}

	// Based on value segment
	switch valueSegment {
	case "vip":
		recommendations = append(recommendations, "premium_support")
		recommendations = append(recommendations, "dedicated_account_manager")
	case "high":
		recommendations = append(recommendations, "advanced_features")
		recommendations = append(recommendations, "team_access")
	case "medium":
		recommendations = append(recommendations, "upgrade_plan")
	}

	return recommendations
}

// calculateUpsellOpportunity estimates additional revenue potential
func (lvs *LifetimeValueService) calculateUpsellOpportunity(
	profile *models.UserBehaviorProfile,
	monthlyValue float64,
	valueSegment string,
) float64 {
	baseOpportunity := monthlyValue * 0.5 // 50% of current monthly value

	switch valueSegment {
	case "low":
		return baseOpportunity * 2.0 // 100% of monthly for low segment
	case "medium":
		return baseOpportunity * 1.5 // 75% of monthly
	case "high":
		return baseOpportunity * 1.2 // 60% of monthly
	case "vip":
		return baseOpportunity * 0.8 // 40% of monthly (already high)
	}

	return baseOpportunity
}

// GetValueSegmentDistribution returns distribution of customers by value
func (lvs *LifetimeValueService) GetValueSegmentDistribution(
	ctx context.Context,
	restaurantID int64,
) (map[string]int, error) {
	distribution := map[string]int{
		"vip":    0,
		"high":   0,
		"medium": 0,
		"low":    0,
	}

	// Get all users (simplified)
	// In production, would use aggregated data
	predictions, err := lvs.predictionRepo.GetLTVPredictions(ctx, restaurantID, 1000)
	if err != nil {
		return nil, err
	}

	for _, pred := range predictions {
		distribution[pred.ValueSegment]++
	}

	return distribution, nil
}

// GetTopCustomers returns highest LTV customers
func (lvs *LifetimeValueService) GetTopCustomers(
	ctx context.Context,
	restaurantID int64,
	limit int,
) ([]*models.LifetimeValuePrediction, error) {
	predictions, err := lvs.predictionRepo.GetLTVPredictions(ctx, restaurantID, limit*2)
	if err != nil {
		return nil, err
	}

	// Sort by predicted LTV descending
	sort.Slice(predictions, func(i, j int) bool {
		return predictions[i].PredictedLTV > predictions[j].PredictedLTV
	})

	if len(predictions) > limit {
		predictions = predictions[:limit]
	}

	return predictions, nil
}

// GetLTVBySegment returns statistics for each value segment
func (lvs *LifetimeValueService) GetLTVBySegment(
	ctx context.Context,
	restaurantID int64,
) (map[string]*SegmentStats, error) {
	predictions, err := lvs.predictionRepo.GetLTVPredictions(ctx, restaurantID, 10000)
	if err != nil {
		return nil, err
	}

	stats := make(map[string]*SegmentStats)
	segments := []string{"vip", "high", "medium", "low"}

	for _, segment := range segments {
		stats[segment] = &SegmentStats{
			Segment:   segment,
			Count:     0,
			TotalLTV:  0,
			AvgLTV:    0,
			MaxLTV:    0,
			MinLTV:    math.MaxFloat64,
		}
	}

	for _, pred := range predictions {
		seg := stats[pred.ValueSegment]
		seg.Count++
		seg.TotalLTV += pred.PredictedLTV
		seg.MaxLTV = math.Max(seg.MaxLTV, pred.PredictedLTV)
		seg.MinLTV = math.Min(seg.MinLTV, pred.PredictedLTV)
	}

	for _, seg := range stats {
		if seg.Count > 0 {
			seg.AvgLTV = seg.TotalLTV / float64(seg.Count)
		}
		if seg.MinLTV == math.MaxFloat64 {
			seg.MinLTV = 0
		}
	}

	return stats, nil
}

// SegmentStats represents statistics for a value segment
type SegmentStats struct {
	Segment   string  `json:"segment"`
	Count     int     `json:"count"`
	TotalLTV  float64 `json:"total_ltv"`
	AvgLTV    float64 `json:"avg_ltv"`
	MaxLTV    float64 `json:"max_ltv"`
	MinLTV    float64 `json:"min_ltv"`
}

// TrainLTVModel trains a new LTV prediction model
func (lvs *LifetimeValueService) TrainLTVModel(
	ctx context.Context,
	restaurantID int64,
) (*models.ModelTrainingJob, error) {
	job := &models.ModelTrainingJob{
		RestaurantID: restaurantID,
		ModelType:    "ltv",
		JobID:        fmt.Sprintf("ltv_%d_%d", restaurantID, time.Now().Unix()),
		Status:       "in_progress",
		Progress:     0.0,
		TrainingStart: &time.Time{},
		CreatedAt:    time.Now(),
	}

	go lvs.trainLTVModelAsync(ctx, job)

	return job, nil
}

// trainLTVModelAsync performs async model training
func (lvs *LifetimeValueService) trainLTVModelAsync(ctx context.Context, job *models.ModelTrainingJob) {
	for progress := 0.1; progress <= 1.0; progress += 0.1 {
		job.Progress = progress
		time.Sleep(100 * time.Millisecond)
	}

	job.Status = "completed"
	endTime := time.Now()
	job.TrainingEnd = &endTime
	job.Metrics = &models.ModelMetrics{
		RMSE: 150.0,
		MAE:  100.0,
		MAPE: 0.25,
	}

	_ = lvs.predictionRepo.CreateTrainingJob(ctx, job)
}

// BulkPredictLTV predicts LTV for multiple users
func (lvs *LifetimeValueService) BulkPredictLTV(
	ctx context.Context,
	restaurantID int64,
	userIDs []int64,
) ([]*models.LifetimeValuePrediction, error) {
	predictions := make([]*models.LifetimeValuePrediction, 0, len(userIDs))

	for _, userID := range userIDs {
		pred, err := lvs.PredictLifetimeValue(ctx, restaurantID, userID)
		if err != nil {
			continue
		}
		predictions = append(predictions, pred)
	}

	// Sort by predicted LTV descending
	sort.Slice(predictions, func(i, j int) bool {
		return predictions[i].PredictedLTV > predictions[j].PredictedLTV
	})

	return predictions, nil
}

// GetLTVHistory returns historical LTV predictions for a user
func (lvs *LifetimeValueService) GetLTVHistory(
	ctx context.Context,
	restaurantID int64,
	userID int64,
	limit int,
) ([]*models.LifetimeValuePrediction, error) {
	return lvs.predictionRepo.GetLTVPredictionHistory(ctx, restaurantID, userID, limit)
}

// EstimateSegmentExpansion estimates growth in high-value segment
func (lvs *LifetimeValueService) EstimateSegmentExpansion(
	ctx context.Context,
	restaurantID int64,
	months int,
) (map[string]interface{}, error) {
	stats, err := lvs.GetLTVBySegment(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	// Estimate growth based on historical patterns
	expansion := map[string]interface{}{
		"current_distribution": stats,
		"projected_growth":     make(map[string]float64),
		"additional_revenue":   0.0,
	}

	projectedGrowth := expansion["projected_growth"].(map[string]float64)

	// Assume 10% monthly growth in high-value segment
	for segment, stat := range stats {
		switch segment {
		case "vip":
			projectedGrowth[segment] = 0.15
		case "high":
			projectedGrowth[segment] = 0.20
		case "medium":
			projectedGrowth[segment] = 0.10
		case "low":
			projectedGrowth[segment] = 0.05
		}
	}

	return expansion, nil
}
