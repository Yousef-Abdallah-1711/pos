//go:build legacy
// +build legacy

package service

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"pos-saas/internal/repository"
)

// FeatureExtractionService handles extraction of ML features from user behavior
type FeatureExtractionService struct {
	behaviorRepo  repository.UserBehaviorRepository
	featureRepo   repository.UserFeatureRepository
	analyticsRepo repository.AnalyticsEventRepository
	componentRepo repository.ComponentRepository
	mu            sync.RWMutex
}

// NewFeatureExtractionService creates a new feature extraction service
func NewFeatureExtractionService(
	behaviorRepo repository.UserBehaviorRepository,
	featureRepo repository.UserFeatureRepository,
	analyticsRepo repository.AnalyticsEventRepository,
	componentRepo repository.ComponentRepository,
) *FeatureExtractionService {
	return &FeatureExtractionService{
		behaviorRepo:  behaviorRepo,
		featureRepo:   featureRepo,
		analyticsRepo: analyticsRepo,
		componentRepo: componentRepo,
	}
}

// ExtractBehaviorProfile builds a comprehensive user behavior profile from events
func (fes *FeatureExtractionService) ExtractBehaviorProfile(
	ctx context.Context,
	restaurantID int64,
	userID int64,
	lookbackDays int,
) (*models.UserBehaviorProfile, error) {
	// Get all analytics events for this user
	cutoffDate := time.Now().Add(-time.Duration(lookbackDays*24) * time.Hour)
	events, err := fes.analyticsRepo.GetEventsByUserAndRestaurant(
		ctx,
		userID,
		restaurantID,
		cutoffDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	// Initialize profile
	profile := &models.UserBehaviorProfile{
		RestaurantID:          restaurantID,
		UserID:                userID,
		ComponentViews:        make(map[string]int64),
		ComponentClicks:       make(map[string]int64),
		ComponentInteractions: make(map[string]int64),
		CategoryPreferences:   make(map[string]float64),
		TagPreferences:        make(map[string]float64),
		TimeSpentPerComponent: make(map[string]int64),
		PreferredThemes:       []string{},
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	// Process events
	sessionStart := time.Time{}
	var sessionEvents []*models.AnalyticsEvent

	for _, event := range events {
		// Group events into sessions (30-minute timeout)
		if sessionStart.IsZero() || event.Timestamp.Sub(sessionStart) > 30*time.Minute {
			if len(sessionEvents) > 0 {
				fes.processSession(ctx, profile, sessionEvents)
			}
			sessionStart = event.Timestamp
			sessionEvents = []*models.AnalyticsEvent{}
		}

		sessionEvents = append(sessionEvents, event)
		profile.TotalInteractions++
		profile.LastActiveTime = event.Timestamp

		// Extract component-level metrics
		if componentID, ok := event.Properties["component_id"].(float64); ok {
			compID := fmt.Sprintf("%d", int64(componentID))

			switch event.EventType {
			case models.EventTypeComponentLoad:
				profile.ComponentViews[compID]++

			case models.EventTypeClick:
				profile.ComponentClicks[compID]++
				profile.ComponentInteractions[compID]++

			case models.EventTypeFormSubmit:
				profile.ComponentInteractions[compID]++
			}
		}

		// Extract category preferences
		if category, ok := event.Properties["category"].(string); ok {
			profile.CategoryPreferences[category]++
		}

		// Extract tag preferences
		if tags, ok := event.Properties["tags"].([]interface{}); ok {
			for _, tag := range tags {
				if tagStr, ok := tag.(string); ok {
					profile.TagPreferences[tagStr]++
				}
			}
		}
	}

	// Process remaining session
	if len(sessionEvents) > 0 {
		fes.processSession(ctx, profile, sessionEvents)
	}

	// Normalize preferences to 0-1 range
	profile.CategoryPreferences = fes.normalizeScores(profile.CategoryPreferences)
	profile.TagPreferences = fes.normalizeScores(profile.TagPreferences)

	// Extract device and browser info from latest event
	if len(events) > 0 {
		lastEvent := events[len(events)-1]
		if deviceType, ok := lastEvent.Properties["device_type"].(string); ok {
			profile.DeviceType = deviceType
		}
		if browserType, ok := lastEvent.Properties["browser_type"].(string); ok {
			profile.BrowserType = browserType
		}
	}

	return profile, nil
}

// processSession extracts features from a session
func (fes *FeatureExtractionService) processSession(
	ctx context.Context,
	profile *models.UserBehaviorProfile,
	events []*models.AnalyticsEvent,
) {
	if len(events) == 0 {
		return
	}

	profile.SessionCount++

	// Calculate session duration
	if len(events) > 1 {
		duration := events[len(events)-1].Timestamp.Sub(events[0].Timestamp).Seconds()
		if duration > 0 {
			profile.AverageSessionDuration = int64((float64(profile.AverageSessionDuration) + duration) / 2)
		}
	}

	// Track time spent per component
	for i := 0; i < len(events)-1; i++ {
		event := events[i]
		nextEvent := events[i+1]

		if componentID, ok := event.Properties["component_id"].(float64); ok {
			compID := fmt.Sprintf("%d", int64(componentID))
			timeSpent := nextEvent.Timestamp.Sub(event.Timestamp).Milliseconds()

			if timeSpent > 0 && timeSpent < 600000 { // Cap at 10 minutes per component
				profile.TimeSpentPerComponent[compID] += timeSpent
			}
		}
	}
}

// normalizeScores normalizes scores to 0-1 range
func (fes *FeatureExtractionService) normalizeScores(
	scores map[string]float64,
) map[string]float64 {
	normalized := make(map[string]float64)

	// Find max score
	var maxScore float64
	for _, score := range scores {
		if score > maxScore {
			maxScore = score
		}
	}

	if maxScore == 0 {
		return normalized
	}

	// Normalize all scores
	for key, score := range scores {
		normalized[key] = score / maxScore
	}

	return normalized
}

// ComputeFeatures calculates ML features from behavior profile
func (fes *FeatureExtractionService) ComputeFeatures(
	ctx context.Context,
	profile *models.UserBehaviorProfile,
) *models.UserFeatures {
	features := &models.UserFeatures{
		UserBehaviorProfileID: profile.ID,
		RestaurantID:          profile.RestaurantID,
		UserID:                profile.UserID,
		ModelVersion:          "1.0.0",
		ComputedAt:            time.Now(),
		CreatedAt:             time.Now(),
	}

	if profile.TotalInteractions == 0 {
		return features
	}

	// Engagement Score: Based on total interactions
	// 0-1000 interactions = 0-1.0 score
	features.EngagementScore = math.Min(float64(profile.TotalInteractions)/1000.0, 1.0)

	// Recency Score: Based on last activity time
	timeSinceActive := time.Since(profile.LastActiveTime)
	switch {
	case timeSinceActive < 1*time.Hour:
		features.RecencyScore = 1.0
	case timeSinceActive < 24*time.Hour:
		features.RecencyScore = 0.8
	case timeSinceActive < 7*24*time.Hour:
		features.RecencyScore = 0.5
	case timeSinceActive < 30*24*time.Hour:
		features.RecencyScore = 0.2
	default:
		features.RecencyScore = 0.05
	}

	// Diversity Score: Based on variety of components viewed
	// Max 100 different components = score of 1.0
	features.DiversityScore = math.Min(float64(len(profile.ComponentViews))/100.0, 1.0)

	// Loyalty Score: Based on session frequency and consistency
	// 100+ sessions = score of 1.0
	features.LoyaltyScore = math.Min(float64(profile.SessionCount)/100.0, 1.0)

	// Consistency Score: Higher if user is active consistently
	// This is simplified - would benefit from time-series analysis
	consistencyScore := math.Min(
		(features.EngagementScore + features.RecencyScore + features.LoyaltyScore) / 3.0,
		1.0,
	)

	// Conversion Likelihood: Predicted probability of conversion
	// Higher engagement + consistency = higher conversion likelihood
	features.ConversionLikelihood = 0.3 + (consistencyScore * 0.6)

	// Churn Risk: Predicted probability of user becoming inactive
	// High recency and engagement = low churn risk
	features.ChurnRisk = 1.0 - ((features.RecencyScore + features.LoyaltyScore) / 2.0)

	// Preferred Hours: Hours when user is most active (simplified - would need hourly analysis)
	features.PreferredHours = []int{9, 10, 11, 14, 15, 16, 18, 19, 20}

	// Generate feature vector for ML models
	features.FeatureVector = []float64{
		features.EngagementScore,
		features.RecencyScore,
		features.DiversityScore,
		features.LoyaltyScore,
		features.ConversionLikelihood,
		features.ChurnRisk,
		consistencyScore,
		float64(profile.SessionCount) / 100.0,
		float64(profile.TotalInteractions) / 1000.0,
		float64(len(profile.ComponentViews)) / 100.0,
	}

	return features
}

// BatchExtractFeatures extracts features for multiple users
func (fes *FeatureExtractionService) BatchExtractFeatures(
	ctx context.Context,
	restaurantID int64,
	userIDs []int64,
	lookbackDays int,
) error {
	// Process users concurrently with limited goroutines
	semaphore := make(chan struct{}, 10) // Max 10 concurrent extractions

	var wg sync.WaitGroup
	errorChan := make(chan error, len(userIDs))

	for _, userID := range userIDs {
		wg.Add(1)

		go func(uid int64) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			select {
			case <-ctx.Done():
				errorChan <- ctx.Err()
				return
			default:
			}

			// Extract behavior profile
			profile, err := fes.ExtractBehaviorProfile(ctx, restaurantID, uid, lookbackDays)
			if err != nil {
				errorChan <- err
				return
			}

			// Save profile
			if err := fes.behaviorRepo.Upsert(ctx, profile); err != nil {
				errorChan <- err
				return
			}

			// Compute features
			features := fes.ComputeFeatures(ctx, profile)

			// Save features
			if err := fes.featureRepo.Upsert(ctx, features); err != nil {
				errorChan <- err
				return
			}
		}(userID)
	}

	wg.Wait()
	close(errorChan)

	// Check for errors
	for err := range errorChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// GetUserPreferences returns user preferences based on features
func (fes *FeatureExtractionService) GetUserPreferences(
	ctx context.Context,
	restaurantID int64,
	userID int64,
) (map[string]interface{}, error) {
	features, err := fes.featureRepo.GetByUserID(ctx, restaurantID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get features: %w", err)
	}

	if features == nil {
		return make(map[string]interface{}), nil
	}

	return map[string]interface{}{
		"engagement":    features.EngagementScore,
		"recency":       features.RecencyScore,
		"diversity":     features.DiversityScore,
		"loyalty":       features.LoyaltyScore,
		"conversion":    features.ConversionLikelihood,
		"churn_risk":    features.ChurnRisk,
		"hours":         features.PreferredHours,
		"model_version": features.ModelVersion,
	}, nil
}
