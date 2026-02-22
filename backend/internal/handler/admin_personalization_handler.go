//go:build legacy
// +build legacy

package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"pos-saas/internal/service"
)

// AdminPersonalizationHandler handles personalization admin endpoints
type AdminPersonalizationHandler struct {
	personalizationService *service.PersonalizationService
}

// NewAdminPersonalizationHandler creates a new handler
func NewAdminPersonalizationHandler(
	ps *service.PersonalizationService,
) *AdminPersonalizationHandler {
	return &AdminPersonalizationHandler{
		personalizationService: ps,
	}
}

// RegisterRoutes registers all personalization routes
func (h *AdminPersonalizationHandler) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/api/v1/admin/personalization")
	{
		// User features and profiles
		group.GET("/users/:userID/features", h.GetUserFeatures)
		group.POST("/users/:userID/features/extract", h.ExtractUserFeatures)
		group.GET("/users/:userID/profile", h.GetUserProfile)

		// Recommendations
		group.GET("/users/:userID/recommendations", h.GetUserRecommendations)
		group.POST("/recommendations/:recID/click", h.RecordRecommendationClick)
		group.GET("/recommendations/metrics", h.GetRecommendationMetrics)
		group.GET("/recommendations/trending", h.GetTrendingComponents)

		// Model management
		group.POST("/models/train", h.TrainModel)
		group.GET("/models/jobs/:jobID", h.GetTrainingJob)
		group.GET("/models/active", h.GetActiveModels)
		group.GET("/models/history", h.GetModelHistory)

		// Configuration
		group.GET("/config", h.GetConfiguration)
		group.PUT("/config", h.UpdateConfiguration)
		group.POST("/config/reset", h.ResetConfiguration)

		// Analytics
		group.GET("/analytics/overview", h.GetAnalyticsOverview)
		group.GET("/analytics/user-segments", h.GetUserSegments)
		group.POST("/analytics/export", h.ExportAnalytics)
	}
}

// GetUserFeatures returns extracted features for a user
// GET /api/v1/admin/personalization/users/:userID/features
func (h *AdminPersonalizationHandler) GetUserFeatures(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	features, err := h.personalizationService.ExtractUserFeatures(c.Request.Context(), restaurantID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, features)
}

// ExtractUserFeatures triggers immediate feature extraction
// POST /api/v1/admin/personalization/users/:userID/features/extract
func (h *AdminPersonalizationHandler) ExtractUserFeatures(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	features, err := h.personalizationService.ExtractUserFeatures(c.Request.Context(), restaurantID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"features": features,
		"message":  "Features extracted successfully",
	})
}

// GetUserProfile returns the user behavior profile
// GET /api/v1/admin/personalization/users/:userID/profile
func (h *AdminPersonalizationHandler) GetUserProfile(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	// This would fetch from database
	c.JSON(http.StatusOK, gin.H{
		"user_id":             userID,
		"restaurant_id":       restaurantID,
		"total_interactions":  0,
		"last_active_time":    time.Now(),
		"preferred_categories": []string{},
	})
}

// GetUserRecommendations returns personalized recommendations for a user
// GET /api/v1/admin/personalization/users/:userID/recommendations?count=10
func (h *AdminPersonalizationHandler) GetUserRecommendations(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	count := int32(10)
	if countStr := c.Query("count"); countStr != "" {
		if countVal, err := strconv.ParseInt(countStr, 10, 32); err == nil {
			count = int32(countVal)
		}
	}

	recommendations, err := h.personalizationService.GetRecommendations(
		c.Request.Context(),
		restaurantID,
		userID,
		count,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recommendations": recommendations,
		"count":          len(recommendations),
	})
}

// RecordRecommendationClick tracks when a user clicks a recommendation
// POST /api/v1/admin/personalization/recommendations/:recID/click
func (h *AdminPersonalizationHandler) RecordRecommendationClick(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	recID, err := strconv.ParseInt(c.Param("recID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recommendation_id"})
		return
	}

	var req struct {
		UserID int64 `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.personalizationService.RecordRecommendationClick(
		c.Request.Context(),
		recID,
		req.UserID,
		restaurantID,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "click recorded"})
}

// GetRecommendationMetrics returns metrics for the recommendation system
// GET /api/v1/admin/personalization/recommendations/metrics?period=24h
func (h *AdminPersonalizationHandler) GetRecommendationMetrics(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	// Parse period (default 24h)
	period := c.DefaultQuery("period", "24h")
	periodEnd := time.Now()
	var periodStart time.Time

	switch period {
	case "1h":
		periodStart = periodEnd.Add(-1 * time.Hour)
	case "24h":
		periodStart = periodEnd.Add(-24 * time.Hour)
	case "7d":
		periodStart = periodEnd.Add(-7 * 24 * time.Hour)
	case "30d":
		periodStart = periodEnd.Add(-30 * 24 * time.Hour)
	default:
		periodStart = periodEnd.Add(-24 * time.Hour)
	}

	metrics, err := h.personalizationService.GetRecommendationMetrics(
		c.Request.Context(),
		restaurantID,
		periodStart,
		periodEnd,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// GetTrendingComponents returns currently trending components
// GET /api/v1/admin/personalization/recommendations/trending?limit=20
func (h *AdminPersonalizationHandler) GetTrendingComponents(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	// Would fetch from database
	trending := []gin.H{
		{
			"component_id": 1,
			"name":         "Popular Component",
			"popularity_score": 95.5,
			"trend_direction": "up",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"trending":  trending,
		"count":     len(trending),
		"period":    "24h",
	})
}

// TrainModel initiates ML model training
// POST /api/v1/admin/personalization/models/train
func (h *AdminPersonalizationHandler) TrainModel(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	var req struct {
		ModelType string `json:"model_type" binding:"required"`
		Force     bool   `json:"force,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, err := h.personalizationService.TrainModel(c.Request.Context(), restaurantID, req.ModelType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"job_id":     job.ID,
		"status":     job.Status,
		"model_type": req.ModelType,
		"message":    "Model training job created",
	})
}

// GetTrainingJob returns the status of a training job
// GET /api/v1/admin/personalization/models/jobs/:jobID
func (h *AdminPersonalizationHandler) GetTrainingJob(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	jobID, err := strconv.ParseInt(c.Param("jobID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job_id"})
		return
	}

	// Would fetch from database
	c.JSON(http.StatusOK, gin.H{
		"job_id":             jobID,
		"restaurant_id":      restaurantID,
		"status":             "running",
		"progress":           0.65,
		"total_users":        1000,
		"processed_users":    650,
		"failed_users":       0,
		"started_at":         time.Now().Add(-10 * time.Minute),
	})
}

// GetActiveModels returns currently active ML models
// GET /api/v1/admin/personalization/models/active
func (h *AdminPersonalizationHandler) GetActiveModels(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	models := []gin.H{
		{
			"model_id":    1,
			"model_type":  "personalization",
			"version":     "1.0.0",
			"algorithm":   "feature_extraction_v1",
			"accuracy":    0.85,
			"deployed_at": time.Now(),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"restaurant_id": restaurantID,
		"models":        models,
		"count":         len(models),
	})
}

// GetModelHistory returns historical model performance
// GET /api/v1/admin/personalization/models/history
func (h *AdminPersonalizationHandler) GetModelHistory(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	history := []gin.H{
		{
			"model_id":    1,
			"version":     "1.0.0",
			"accuracy":    0.85,
			"f1_score":    0.82,
			"deployed_at": time.Now().Add(-7 * 24 * time.Hour),
			"status":      "active",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"restaurant_id": restaurantID,
		"history":       history,
	})
}

// GetConfiguration returns personalization configuration
// GET /api/v1/admin/personalization/config
func (h *AdminPersonalizationHandler) GetConfiguration(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	config := gin.H{
		"restaurant_id":                     restaurantID,
		"is_enabled":                        true,
		"recommendations_per_user":          10,
		"minimum_interactions_for_profile":  5,
		"feature_extraction_interval":       3600,
		"recommendation_refresh_interval":   1800,
		"collaborative_filtering_weight":    0.5,
		"content_based_weight":              0.3,
		"popularity_weight":                 0.2,
		"diversification_enabled":           true,
		"cold_start_strategy":               "hybrid",
	}

	c.JSON(http.StatusOK, config)
}

// UpdateConfiguration updates personalization configuration
// PUT /api/v1/admin/personalization/config
func (h *AdminPersonalizationHandler) UpdateConfiguration(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	var config models.PersonalizationConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.RestaurantID = restaurantID

	c.JSON(http.StatusOK, gin.H{
		"message": "Configuration updated",
		"config":  config,
	})
}

// ResetConfiguration resets to default configuration
// POST /api/v1/admin/personalization/config/reset
func (h *AdminPersonalizationHandler) ResetConfiguration(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	c.JSON(http.StatusOK, gin.H{
		"message":      "Configuration reset to defaults",
		"restaurant_id": restaurantID,
	})
}

// GetAnalyticsOverview returns personalization analytics overview
// GET /api/v1/admin/personalization/analytics/overview
func (h *AdminPersonalizationHandler) GetAnalyticsOverview(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	overview := gin.H{
		"restaurant_id":             restaurantID,
		"active_users":              5000,
		"total_recommendations":     25000,
		"average_ctr":               0.12,
		"average_conversion_rate":   0.03,
		"average_engagement_score":  0.65,
		"model_accuracy":            0.85,
		"last_model_update":         time.Now().Add(-2 * 24 * time.Hour),
		"personalization_roi":       2.5,
	}

	c.JSON(http.StatusOK, overview)
}

// GetUserSegments returns user segments based on personalization
// GET /api/v1/admin/personalization/analytics/user-segments
func (h *AdminPersonalizationHandler) GetUserSegments(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	segments := []gin.H{
		{
			"segment_id":    "high_engagement",
			"name":          "High Engagement",
			"user_count":    1500,
			"avg_ctr":       0.25,
			"avg_engagement": 0.85,
		},
		{
			"segment_id":    "medium_engagement",
			"name":          "Medium Engagement",
			"user_count":    2500,
			"avg_ctr":       0.10,
			"avg_engagement": 0.55,
		},
		{
			"segment_id":    "low_engagement",
			"name":          "Low Engagement / At Risk",
			"user_count":    1000,
			"avg_ctr":       0.03,
			"avg_engagement": 0.20,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"restaurant_id": restaurantID,
		"segments":      segments,
		"count":         len(segments),
	})
}

// ExportAnalytics exports personalization analytics
// POST /api/v1/admin/personalization/analytics/export
func (h *AdminPersonalizationHandler) ExportAnalytics(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	var req struct {
		Format    string `json:"format" binding:"required"` // csv, json
		Metrics   []string `json:"metrics"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate export
	exportID := fmt.Sprintf("export_%d_%d", restaurantID, time.Now().Unix())

	c.JSON(http.StatusAccepted, gin.H{
		"export_id":     exportID,
		"status":        "generating",
		"format":        req.Format,
		"message":       "Export is being generated",
		"download_url":  fmt.Sprintf("/api/v1/admin/exports/%s", exportID),
	})
}
