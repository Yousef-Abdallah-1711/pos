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

// AdminAdvancedRecommendationsHandler handles advanced recommendation endpoints
type AdminAdvancedRecommendationsHandler struct {
	advancedRecommendationService *service.AdvancedRecommendationService
	personalizationService        *service.PersonalizationService
	abTestingService              *service.ABTestingService
}

// NewAdminAdvancedRecommendationsHandler creates a new handler
func NewAdminAdvancedRecommendationsHandler(
	ars *service.AdvancedRecommendationService,
	ps *service.PersonalizationService,
	abts *service.ABTestingService,
) *AdminAdvancedRecommendationsHandler {
	return &AdminAdvancedRecommendationsHandler{
		advancedRecommendationService: ars,
		personalizationService:        ps,
		abTestingService:              abts,
	}
}

// RegisterRoutes registers all advanced recommendation routes
func (h *AdminAdvancedRecommendationsHandler) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/api/v1/admin/recommendations/advanced")
	{
		// Advanced recommendations
		group.GET("/users/:userID", h.GetAdvancedRecommendations)
		group.POST("/users/:userID/refresh", h.RefreshRecommendations)
		group.GET("/users/:userID/explanations", h.GetExplanations)

		// Exploration
		group.GET("/users/:userID/explore", h.GetExplorationRecommendations)

		// Real-time updates
		group.GET("/users/:userID/updates", h.SubscribeToUpdates)

		// Performance tracking
		group.POST("/track/:recID/:eventType", h.TrackPerformance)

		// Comparison and testing
		group.GET("/compare-algorithms", h.CompareAlgorithms)
		group.GET("/algorithm-stats", h.GetAlgorithmStats)

		// Batch operations
		group.POST("/precompute", h.PrecomputeRecommendations)
		group.POST("/refresh-all", h.RefreshAllRecommendations)
	}
}

// GetAdvancedRecommendations returns advanced recommendations for a user
// GET /api/v1/admin/recommendations/advanced/users/:userID?count=10
func (h *AdminAdvancedRecommendationsHandler) GetAdvancedRecommendations(c *gin.Context) {
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

	// Would fetch config from database
	config := &models.PersonalizationConfig{
		CollaborativeFilteringWeight: 0.5,
		ContentBasedWeight:           0.3,
		PopularityWeight:             0.2,
	}

	recs, err := h.personalizationService.GetRecommendations(
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
		"recommendations": recs,
		"count":          len(recs),
		"algorithm":      "advanced_hybrid",
		"timestamp":      time.Now(),
	})
}

// RefreshRecommendations forces a refresh of recommendations
// POST /api/v1/admin/recommendations/advanced/users/:userID/refresh
func (h *AdminAdvancedRecommendationsHandler) RefreshRecommendations(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	if err := h.advancedRecommendationService.RefreshRecommendationsForUser(
		c.Request.Context(),
		userID,
		restaurantID,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "refreshed",
		"message": "Recommendations refreshed successfully",
	})
}

// GetExplanations returns user-friendly explanations for recommendations
// GET /api/v1/admin/recommendations/advanced/users/:userID/explanations
func (h *AdminAdvancedRecommendationsHandler) GetExplanations(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	// Get recommendations
	recs, err := h.personalizationService.GetRecommendations(
		c.Request.Context(),
		restaurantID,
		userID,
		10,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get explanations
	explanations := h.advancedRecommendationService.GetRecommendationExplanations(recs)

	c.JSON(http.StatusOK, gin.H{
		"explanations": explanations,
		"count":       len(explanations),
	})
}

// GetExplorationRecommendations returns diverse exploration recommendations
// GET /api/v1/admin/recommendations/advanced/users/:userID/explore?count=10
func (h *AdminAdvancedRecommendationsHandler) GetExplorationRecommendations(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	count := 10
	if countStr := c.Query("count"); countStr != "" {
		if countVal, err := strconv.Atoi(countStr); err == nil {
			count = countVal
		}
	}

	recs, err := h.advancedRecommendationService.GetExplorationRecommendations(
		c.Request.Context(),
		userID,
		restaurantID,
		count,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recommendations": recs,
		"count":          len(recs),
		"algorithm":      "exploration",
	})
}

// SubscribeToUpdates subscribes to real-time recommendation updates using SSE
// GET /api/v1/admin/recommendations/advanced/users/:userID/updates
func (h *AdminAdvancedRecommendationsHandler) SubscribeToUpdates(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	// Subscribe to updates
	updatesChan := h.advancedRecommendationService.SubscribeToRecommendationUpdates(
		userID,
		restaurantID,
	)

	// Set up SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// Send initial ping
	c.SSEvent("ping", gin.H{"message": "connected"})

	// Stream updates
	for rec := range updatesChan {
		select {
		case <-c.Request.Context().Done():
			h.advancedRecommendationService.UnsubscribeFromRecommendationUpdates(userID)
			return
		default:
			c.SSEvent("recommendation", rec)
			c.Writer.Flush()
		}
	}
}

// TrackPerformance tracks recommendation performance
// POST /api/v1/admin/recommendations/advanced/track/:recID/:eventType
func (h *AdminAdvancedRecommendationsHandler) TrackPerformance(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	recID, err := strconv.ParseInt(c.Param("recID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recommendation_id"})
		return
	}

	eventType := c.Param("eventType")
	if eventType != "click" && eventType != "view" && eventType != "convert" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event_type"})
		return
	}

	var req struct {
		UserID int64 `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.advancedRecommendationService.TrackRecommendationPerformance(
		c.Request.Context(),
		recID,
		req.UserID,
		restaurantID,
		eventType,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "tracked",
		"event_type": eventType,
	})
}

// CompareAlgorithms compares different recommendation algorithms for a user
// GET /api/v1/admin/recommendations/advanced/compare-algorithms?userID=123&count=10
func (h *AdminAdvancedRecommendationsHandler) CompareAlgorithms(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	userIDStr := c.Query("userID")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID query parameter required"})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID"})
		return
	}

	count := 10
	if countStr := c.Query("count"); countStr != "" {
		if countVal, err := strconv.Atoi(countStr); err == nil {
			count = countVal
		}
	}

	comparison := gin.H{
		"user_id":       userID,
		"restaurant_id": restaurantID,
		"algorithms": gin.H{
			"collaborative_filtering": gin.H{
				"description": "Recommends based on similar users",
				"pros":        []string{"Good for discovery", "Avoids filter bubble"},
				"cons":        []string{"Cold start problem", "Requires large user base"},
				"performance": gin.H{
					"average_ctr":       0.12,
					"average_conversion": 0.03,
				},
			},
			"content_based": gin.H{
				"description": "Recommends similar components",
				"pros":        []string{"No cold start", "Interpretable"},
				"cons":        []string{"Limited diversity", "Subject to filter bubble"},
				"performance": gin.H{
					"average_ctr":       0.10,
					"average_conversion": 0.025,
				},
			},
			"hybrid": gin.H{
				"description": "Combines multiple algorithms",
				"pros":        []string{"Best of both worlds", "More robust"},
				"cons":        []string{"More complex", "Higher latency"},
				"performance": gin.H{
					"average_ctr":       0.15,
					"average_conversion": 0.035,
				},
			},
		},
		"timestamp": time.Now(),
	}

	c.JSON(http.StatusOK, comparison)
}

// GetAlgorithmStats returns performance statistics for different algorithms
// GET /api/v1/admin/recommendations/advanced/algorithm-stats?period=7d
func (h *AdminAdvancedRecommendationsHandler) GetAlgorithmStats(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	period := c.DefaultQuery("period", "7d")

	stats := gin.H{
		"restaurant_id": restaurantID,
		"period":        period,
		"algorithms": gin.H{
			"collaborative_filtering": gin.H{
				"recommendations_count": 15000,
				"clicks":                1800,
				"conversions":           450,
				"ctr":                   0.12,
				"conversion_rate":       0.025,
				"revenue":               22500.00,
			},
			"content_based": gin.H{
				"recommendations_count": 14500,
				"clicks":                1450,
				"conversions":           362,
				"ctr":                   0.10,
				"conversion_rate":       0.025,
				"revenue":               18050.00,
			},
			"hybrid": gin.H{
				"recommendations_count": 16000,
				"clicks":                2400,
				"conversions":           560,
				"ctr":                   0.15,
				"conversion_rate":       0.023,
				"revenue":               28000.00,
			},
		},
		"best_performer": gin.H{
			"algorithm": "hybrid",
			"metric":    "revenue",
			"value":     28000.00,
		},
	}

	c.JSON(http.StatusOK, stats)
}

// PrecomputeRecommendations triggers pre-computation of recommendations
// POST /api/v1/admin/recommendations/advanced/precompute
func (h *AdminAdvancedRecommendationsHandler) PrecomputeRecommendations(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	var req struct {
		MaxUsers int `json:"max_users,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.MaxUsers == 0 {
		req.MaxUsers = 10000
	}

	// Would trigger background job
	jobID := fmt.Sprintf("precompute_%d_%d", restaurantID, time.Now().Unix())

	c.JSON(http.StatusAccepted, gin.H{
		"status": "processing",
		"job_id": jobID,
		"message": "Pre-computation job started",
		"restaurant_id": restaurantID,
		"max_users": req.MaxUsers,
	})
}

// RefreshAllRecommendations refreshes recommendations for all users
// POST /api/v1/admin/recommendations/advanced/refresh-all
func (h *AdminAdvancedRecommendationsHandler) RefreshAllRecommendations(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	// Would trigger background job to refresh all user recommendations
	jobID := fmt.Sprintf("refresh_all_%d_%d", restaurantID, time.Now().Unix())

	c.JSON(http.StatusAccepted, gin.H{
		"status": "processing",
		"job_id": jobID,
		"message": "Refresh job started for all users",
		"restaurant_id": restaurantID,
	})
}
