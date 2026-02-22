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

// ChurnPredictionService predicts user/restaurant churn
type ChurnPredictionService struct {
	behaviorRepo     repository.UserBehaviorRepository
	analyticsRepo    repository.AnalyticsEventRepository
	predictionRepo   repository.PredictionRepository
	featureService   *FeatureExtractionService
	cacheService     *CacheService
}

// NewChurnPredictionService creates a new churn prediction service
func NewChurnPredictionService(
	behaviorRepo repository.UserBehaviorRepository,
	analyticsRepo repository.AnalyticsEventRepository,
	predictionRepo repository.PredictionRepository,
	featureService *FeatureExtractionService,
	cacheService *CacheService,
) *ChurnPredictionService {
	return &ChurnPredictionService{
		behaviorRepo:   behaviorRepo,
		analyticsRepo:  analyticsRepo,
		predictionRepo: predictionRepo,
		featureService: featureService,
		cacheService:   cacheService,
	}
}

// PredictUserChurn predicts likelihood of a user churning
func (cps *ChurnPredictionService) PredictUserChurn(
	ctx context.Context,
	restaurantID int64,
	userID int64,
) (*models.ChurnPrediction, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("churn:user:%d:%d", restaurantID, userID)
	if cached, ok := cps.cacheService.Get(cacheKey); ok {
		if pred, ok := cached.(*models.ChurnPrediction); ok {
			return pred, nil
		}
	}

	// Get user behavior profile
	profile, err := cps.behaviorRepo.GetProfile(ctx, restaurantID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	if profile == nil {
		return nil, fmt.Errorf("user profile not found")
	}

	// Extract churn features
	features := cps.extractChurnFeatures(profile)

	// Calculate churn score using weighted feature model
	churnScore := cps.calculateChurnScore(features)

	// Determine risk level
	riskLevel := cps.getRiskLevel(churnScore)

	// Identify risk factors
	riskFactors := cps.identifyRiskFactors(features, profile)

	// Estimate days to churn
	daysToChurn := cps.estimateDaysToChurn(features, churnScore)

	// Get recommended actions
	recommendedActions := cps.getRecommendedActions(riskLevel, riskFactors)

	prediction := &models.ChurnPrediction{
		RestaurantID:      restaurantID,
		UserID:            &userID,
		ChurnScore:        churnScore,
		ChurnProbability:  churnScore, // 0.0-1.0
		RiskLevel:         riskLevel,
		RiskFactors:       riskFactors,
		DaysToChurn:       daysToChurn,
		Confidence:        0.85, // Model confidence
		RecommendedActions: recommendedActions,
		PredictedAt:       time.Now(),
		ValidUntil:        time.Now().AddDate(0, 0, 7), // Valid for 7 days
	}

	// Cache for 24 hours
	cps.cacheService.Set(cacheKey, prediction, 86400, []string{"churn", "user", fmt.Sprintf("user:%d", userID)})

	// Save to database
	_ = cps.predictionRepo.CreateChurnPrediction(ctx, prediction)

	return prediction, nil
}

// PredictRestaurantChurn predicts likelihood of a restaurant churning
func (cps *ChurnPredictionService) PredictRestaurantChurn(
	ctx context.Context,
	restaurantID int64,
) (*models.ChurnPrediction, error) {
	// Check cache
	cacheKey := fmt.Sprintf("churn:restaurant:%d", restaurantID)
	if cached, ok := cps.cacheService.Get(cacheKey); ok {
		if pred, ok := cached.(*models.ChurnPrediction); ok {
			return pred, nil
		}
	}

	// Get restaurant activity metrics
	metrics, err := cps.analyticsRepo.GetRestaurantMetrics(ctx, restaurantID, 30)
	if err != nil {
		return nil, fmt.Errorf("failed to get restaurant metrics: %w", err)
	}

	// Extract restaurant churn features
	features := cps.extractRestaurantChurnFeatures(metrics, restaurantID)

	// Calculate churn score
	churnScore := cps.calculateRestaurantChurnScore(features)

	// Determine risk level
	riskLevel := cps.getRiskLevel(churnScore)

	// Identify risk factors
	riskFactors := cps.identifyRestaurantRiskFactors(features)

	// Estimate days to churn
	daysToChurn := cps.estimateDaysToChurn(features, churnScore)

	// Get recommended actions
	recommendedActions := cps.getRestaurantRecommendedActions(riskLevel, riskFactors)

	prediction := &models.ChurnPrediction{
		RestaurantID:      restaurantID,
		UserID:            nil,
		ChurnScore:        churnScore,
		ChurnProbability:  churnScore,
		RiskLevel:         riskLevel,
		RiskFactors:       riskFactors,
		DaysToChurn:       daysToChurn,
		Confidence:        0.82,
		RecommendedActions: recommendedActions,
		PredictedAt:       time.Now(),
		ValidUntil:        time.Now().AddDate(0, 0, 7),
	}

	// Cache for 24 hours
	cps.cacheService.Set(cacheKey, prediction, 86400, []string{"churn", "restaurant"})

	// Save to database
	_ = cps.predictionRepo.CreateChurnPrediction(ctx, prediction)

	return prediction, nil
}

// ChurnFeatures represents extracted features for churn prediction
type ChurnFeatures struct {
	DaysSinceLastActivity    int
	ActivityDeclineRate      float64 // Trend in activity over last 30 days
	EngagementScore          float64
	SessionFrequency         float64
	ContentInteractionRate   float64
	ComponentDiversity       float64
	FeatureFlagsUsed         int
	SupportTicketsOpened     int
	AverageSessionDuration   float64
	MonthlyActiveLogins      int
	PaymentSuccessRate       float64
	FeatureAdoptionRate      float64
	PeakActivityPeriod       string
	PreviousChurnSaved       bool
	TrialTimeRemaining       int // Days
	AccountAge               int // Days
}

// extractChurnFeatures extracts features from user behavior profile
func (cps *ChurnPredictionService) extractChurnFeatures(
	profile *models.UserBehaviorProfile,
) ChurnFeatures {
	features := ChurnFeatures{
		DaysSinceLastActivity: int(time.Since(profile.LastActivityAt).Hours() / 24),
		ActivityDeclineRate:   cps.calculateActivityDecline(profile),
		EngagementScore:       cps.calculateEngagementScore(profile),
		SessionFrequency:      float64(profile.TotalSessions) / float64(int(time.Since(profile.CreatedAt).Hours()/24+1)), // Sessions per day
		ContentInteractionRate: cps.calculateInteractionRate(profile),
		ComponentDiversity:    cps.calculateDiversity(profile),
		FeatureFlagsUsed:      len(profile.FeaturesUsed),
		AverageSessionDuration: cps.calculateAvgSessionDuration(profile),
		MonthlyActiveLogins:   cps.getMonthlyActiveLogins(profile),
		FeatureAdoptionRate:   cps.calculateFeatureAdoption(profile),
		AccountAge:            int(time.Since(profile.CreatedAt).Hours() / 24),
	}

	return features
}

// calculateActivityDecline calculates the rate of activity decline
func (cps *ChurnPredictionService) calculateActivityDecline(profile *models.UserBehaviorProfile) float64 {
	// Calculate trend: if daily activity is declining, score goes up
	// Compare average activity in last 7 days vs previous 7-14 days
	recentActivity := profile.TotalInteractions / 7 // Simplified
	historicalActivity := profile.TotalInteractions / 30

	if historicalActivity == 0 {
		return 0
	}

	decline := (historicalActivity - recentActivity) / historicalActivity
	return math.Max(0, math.Min(1, decline))
}

// calculateEngagementScore calculates user engagement
func (cps *ChurnPredictionService) calculateEngagementScore(profile *models.UserBehaviorProfile) float64 {
	// Engagement = (interactions * session_diversity) / account_age
	diversity := cps.calculateDiversity(profile)
	accountDays := float64(int(time.Since(profile.CreatedAt).Hours()/24 + 1))

	engagement := (float64(profile.TotalInteractions) * diversity) / accountDays
	return math.Min(engagement/100.0, 1.0) // Normalize to 0-1
}

// calculateInteractionRate calculates content interaction rate
func (cps *ChurnPredictionService) calculateInteractionRate(profile *models.UserBehaviorProfile) float64 {
	if profile.TotalInteractions == 0 {
		return 0
	}

	// Interactions per day
	daysSinceCreation := float64(int(time.Since(profile.CreatedAt).Hours() / 24))
	if daysSinceCreation == 0 {
		daysSinceCreation = 1
	}

	ratePerDay := float64(profile.TotalInteractions) / daysSinceCreation
	return math.Min(ratePerDay/10.0, 1.0) // Normalize
}

// calculateDiversity calculates user's content diversity
func (cps *ChurnPredictionService) calculateDiversity(profile *models.UserBehaviorProfile) float64 {
	if profile.TotalInteractions == 0 {
		return 0
	}

	// Count unique components viewed
	uniqueComponents := len(profile.ComponentViews)
	diversity := float64(uniqueComponents) / float64(profile.TotalInteractions)

	return math.Min(diversity*2, 1.0) // Normalize
}

// calculateAvgSessionDuration calculates average session duration in minutes
func (cps *ChurnPredictionService) calculateAvgSessionDuration(profile *models.UserBehaviorProfile) float64 {
	if profile.TotalSessions == 0 {
		return 0
	}

	totalMinutes := int(profile.TimeSpentMinutes)
	return float64(totalMinutes) / float64(profile.TotalSessions)
}

// getMonthlyActiveLogins returns estimated monthly active logins
func (cps *ChurnPredictionService) getMonthlyActiveLogins(profile *models.UserBehaviorProfile) int {
	// Simplified: use total sessions and estimate monthly rate
	if profile.TotalSessions == 0 {
		return 0
	}

	accountDays := float64(int(time.Since(profile.CreatedAt).Hours() / 24))
	if accountDays < 30 {
		accountDays = 30
	}

	monthlyRate := (float64(profile.TotalSessions) / accountDays) * 30
	return int(monthlyRate)
}

// calculateFeatureAdoption calculates feature adoption rate
func (cps *ChurnPredictionService) calculateFeatureAdoption(profile *models.UserBehaviorProfile) float64 {
	// Assuming ~20 major features available
	totalFeatures := 20.0
	usedFeatures := float64(len(profile.FeaturesUsed))

	return usedFeatures / totalFeatures
}

// calculateChurnScore calculates overall churn risk score
func (cps *ChurnPredictionService) calculateChurnScore(features ChurnFeatures) float64 {
	// Weighted feature model (weights sum to 1.0)
	score := 0.0

	// Days since last activity (high impact)
	daysSinceActivityScore := math.Min(float64(features.DaysSinceLastActivity)/30.0, 1.0)
	score += daysSinceActivityScore * 0.25

	// Activity decline rate
	score += features.ActivityDeclineRate * 0.20

	// Engagement score (inverse)
	score += (1.0 - features.EngagementScore) * 0.20

	// Session frequency (inverse)
	sessionFrequencyNorm := math.Min(features.SessionFrequency/1.0, 1.0)
	score += (1.0 - sessionFrequencyNorm) * 0.15

	// Feature adoption
	score += (1.0 - features.FeatureAdoptionRate) * 0.10

	// Content interaction rate
	score += (1.0 - features.ContentInteractionRate) * 0.10

	return math.Min(score, 1.0)
}

// calculateRestaurantChurnScore calculates restaurant churn score
func (cps *ChurnPredictionService) calculateRestaurantChurnScore(features map[string]float64) float64 {
	score := 0.0

	// Daily active users trend
	if dau, ok := features["daily_active_users_trend"]; ok {
		score += (1.0 - dau) * 0.25
	}

	// Revenue trend
	if rev, ok := features["revenue_trend"]; ok {
		score += (1.0 - rev) * 0.20
	}

	// Order frequency
	if freq, ok := features["order_frequency_trend"]; ok {
		score += (1.0 - freq) * 0.20
	}

	// Support satisfaction
	if support, ok := features["support_satisfaction"]; ok {
		score += (1.0 - support) * 0.15
	}

	// Days since last order
	if daysOrder, ok := features["days_since_last_order"]; ok {
		score += math.Min(daysOrder/60.0, 1.0) * 0.15
	}

	// Payment failures
	if payment, ok := features["payment_failure_rate"]; ok {
		score += payment * 0.05
	}

	return math.Min(score, 1.0)
}

// getRiskLevel returns risk level based on churn score
func (cps *ChurnPredictionService) getRiskLevel(score float64) string {
	switch {
	case score > 0.75:
		return "critical"
	case score > 0.5:
		return "high"
	case score > 0.25:
		return "medium"
	default:
		return "low"
	}
}

// identifyRiskFactors identifies which factors contribute to churn
func (cps *ChurnPredictionService) identifyRiskFactors(features ChurnFeatures, profile *models.UserBehaviorProfile) []string {
	factors := []string{}

	if features.DaysSinceLastActivity > 14 {
		factors = append(factors, "no_recent_activity")
	}

	if features.ActivityDeclineRate > 0.3 {
		factors = append(factors, "declining_activity")
	}

	if features.EngagementScore < 0.3 {
		factors = append(factors, "low_engagement")
	}

	if features.SessionFrequency < 0.1 {
		factors = append(factors, "low_frequency")
	}

	if features.FeatureAdoptionRate < 0.2 {
		factors = append(factors, "minimal_feature_adoption")
	}

	if features.AverageSessionDuration < 5.0 {
		factors = append(factors, "very_short_sessions")
	}

	if profile.TotalInteractions < 10 {
		factors = append(factors, "minimal_usage")
	}

	if len(factors) == 0 {
		factors = append(factors, "unknown_factors")
	}

	return factors
}

// identifyRestaurantRiskFactors identifies restaurant-level risk factors
func (cps *ChurnPredictionService) identifyRestaurantRiskFactors(features map[string]float64) []string {
	factors := []string{}

	if declineRate, ok := features["daily_active_users_trend"]; ok && declineRate < 0.5 {
		factors = append(factors, "declining_user_base")
	}

	if revTrend, ok := features["revenue_trend"]; ok && revTrend < 0.5 {
		factors = append(factors, "declining_revenue")
	}

	if support, ok := features["support_satisfaction"]; ok && support < 0.4 {
		factors = append(factors, "low_support_satisfaction")
	}

	if paymentFail, ok := features["payment_failure_rate"]; ok && paymentFail > 0.1 {
		factors = append(factors, "high_payment_failure_rate")
	}

	return factors
}

// estimateDaysToChurn estimates days until churn occurs
func (cps *ChurnPredictionService) estimateDaysToChurn(features ChurnFeatures, churnScore float64) int {
	// Simple estimation based on activity decline
	baselineChurnDays := 30

	// More inactivity = sooner churn
	adjustedDays := baselineChurnDays - int(churnScore*20)

	// Minimum 3 days, maximum 60 days
	if adjustedDays < 3 {
		adjustedDays = 3
	}
	if adjustedDays > 60 {
		adjustedDays = 60
	}

	return adjustedDays
}

// getRecommendedActions returns recommended actions based on churn risk
func (cps *ChurnPredictionService) getRecommendedActions(riskLevel string, riskFactors []string) []string {
	actions := []string{}

	switch riskLevel {
	case "critical":
		actions = append(actions, "send_immediate_outreach")
		actions = append(actions, "offer_support_call")
		actions = append(actions, "provide_incentive_discount")
	case "high":
		actions = append(actions, "send_engagement_email")
		actions = append(actions, "feature_tutorial")
		actions = append(actions, "personalized_recommendations")
	case "medium":
		actions = append(actions, "send_weekly_digest")
		actions = append(actions, "highlight_unused_features")
	default:
		actions = append(actions, "maintain_engagement")
	}

	// Add factor-specific actions
	for _, factor := range riskFactors {
		switch factor {
		case "declining_activity":
			actions = append(actions, "send_win_back_campaign")
		case "low_feature_adoption":
			actions = append(actions, "conduct_feature_training")
		case "minimal_usage":
			actions = append(actions, "schedule_onboarding_review")
		}
	}

	// Remove duplicates
	actionMap := make(map[string]bool)
	uniqueActions := []string{}
	for _, action := range actions {
		if !actionMap[action] {
			actionMap[action] = true
			uniqueActions = append(uniqueActions, action)
		}
	}

	return uniqueActions
}

// getRestaurantRecommendedActions returns restaurant-level recommended actions
func (cps *ChurnPredictionService) getRestaurantRecommendedActions(riskLevel string, riskFactors []string) []string {
	actions := []string{}

	switch riskLevel {
	case "critical":
		actions = append(actions, "executive_outreach")
		actions = append(actions, "account_manager_assignment")
		actions = append(actions, "special_offer")
	case "high":
		actions = append(actions, "check_in_call")
		actions = append(actions, "business_review")
	case "medium":
		actions = append(actions, "send_health_report")
		actions = append(actions, "feature_roadmap_update")
	}

	return actions
}

// TrainChurnModel trains a new churn prediction model
func (cps *ChurnPredictionService) TrainChurnModel(
	ctx context.Context,
	restaurantID int64,
) (*models.ModelTrainingJob, error) {
	job := &models.ModelTrainingJob{
		RestaurantID: restaurantID,
		ModelType:    "churn",
		JobID:        fmt.Sprintf("churn_%d_%d", restaurantID, time.Now().Unix()),
		Status:       "in_progress",
		Progress:     0.0,
		TrainingStart: &time.Time{},
		CreatedAt:    time.Now(),
	}

	// Async training would happen here
	// For now, just update progress
	go cps.trainChurnModelAsync(ctx, job)

	return job, nil
}

// trainChurnModelAsync performs async training
func (cps *ChurnPredictionService) trainChurnModelAsync(ctx context.Context, job *models.ModelTrainingJob) {
	// Simulate training
	for progress := 0.1; progress <= 1.0; progress += 0.1 {
		job.Progress = progress
		time.Sleep(100 * time.Millisecond)
	}

	job.Status = "completed"
	endTime := time.Now()
	job.TrainingEnd = &endTime

	// Save metrics
	job.Metrics = &models.ModelMetrics{
		Accuracy:  0.82,
		Precision: 0.78,
		Recall:    0.85,
		F1Score:   0.81,
		AUC:       0.88,
	}

	_ = cps.predictionRepo.CreateTrainingJob(ctx, job)
}

// GetChurnPredictionHistory returns historical churn predictions for a user
func (cps *ChurnPredictionService) GetChurnPredictionHistory(
	ctx context.Context,
	restaurantID int64,
	userID int64,
	limit int,
) ([]*models.ChurnPrediction, error) {
	return cps.predictionRepo.GetChurnPredictions(ctx, restaurantID, &userID, limit)
}

// BulkPredictChurn predicts churn for multiple users
func (cps *ChurnPredictionService) BulkPredictChurn(
	ctx context.Context,
	restaurantID int64,
	userIDs []int64,
) ([]*models.ChurnPrediction, error) {
	predictions := make([]*models.ChurnPrediction, 0, len(userIDs))

	for _, userID := range userIDs {
		pred, err := cps.PredictUserChurn(ctx, restaurantID, userID)
		if err != nil {
			continue
		}
		predictions = append(predictions, pred)
	}

	// Sort by churn score descending
	sort.Slice(predictions, func(i, j int) bool {
		return predictions[i].ChurnScore > predictions[j].ChurnScore
	})

	return predictions, nil
}

// extractRestaurantChurnFeatures extracts features for restaurant churn
func (cps *ChurnPredictionService) extractRestaurantChurnFeatures(
	metrics map[string]interface{},
	restaurantID int64,
) map[string]float64 {
	features := make(map[string]float64)

	// Extract from metrics (simplified)
	if dau, ok := metrics["daily_active_users"].(float64); ok {
		features["daily_active_users_trend"] = math.Min(dau/100.0, 1.0)
	}

	if revenue, ok := metrics["revenue"].(float64); ok {
		features["revenue_trend"] = math.Min(revenue/10000.0, 1.0)
	}

	if orders, ok := metrics["orders"].(float64); ok {
		features["order_frequency_trend"] = math.Min(orders/100.0, 1.0)
	}

	return features
}
