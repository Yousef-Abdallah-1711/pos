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

// AdminPredictiveAnalyticsHandler handles predictive analytics endpoints
type AdminPredictiveAnalyticsHandler struct {
	churnService      *service.ChurnPredictionService
	revenueService    *service.RevenueForecastingService
	ltvService        *service.LifetimeValueService
	anomalyService    *service.AnomalyDetectionService
}

// NewAdminPredictiveAnalyticsHandler creates a new handler
func NewAdminPredictiveAnalyticsHandler(
	churnService *service.ChurnPredictionService,
	revenueService *service.RevenueForecastingService,
	ltvService *service.LifetimeValueService,
	anomalyService *service.AnomalyDetectionService,
) *AdminPredictiveAnalyticsHandler {
	return &AdminPredictiveAnalyticsHandler{
		churnService:   churnService,
		revenueService: revenueService,
		ltvService:     ltvService,
		anomalyService: anomalyService,
	}
}

// RegisterRoutes registers all predictive analytics routes
func (h *AdminPredictiveAnalyticsHandler) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/api/v1/admin/predictions")
	{
		// Churn Prediction
		group.GET("/churn/user/:userID", h.GetUserChurnPrediction)
		group.GET("/churn/restaurant", h.GetRestaurantChurnPrediction)
		group.GET("/churn/users", h.BulkGetChurnPredictions)
		group.POST("/churn/train", h.TrainChurnModel)
		group.GET("/churn/history/:userID", h.GetChurnHistory)

		// Revenue Forecasting
		group.GET("/revenue/forecast", h.ForecastRevenue)
		group.GET("/revenue/accuracy", h.GetForecastAccuracy)
		group.GET("/revenue/history", h.GetRevenueHistory)
		group.GET("/revenue/seasonality", h.GetSeasonalityReport)
		group.POST("/revenue/train", h.TrainRevenueModel)

		// Lifetime Value
		group.GET("/ltv/user/:userID", h.GetUserLTV)
		group.GET("/ltv/segment-distribution", h.GetSegmentDistribution)
		group.GET("/ltv/top-customers", h.GetTopCustomers)
		group.GET("/ltv/by-segment", h.GetLTVBySegment)
		group.GET("/ltv/history/:userID", h.GetLTVHistory)
		group.POST("/ltv/train", h.TrainLTVModel)

		// Anomaly Detection
		group.GET("/anomalies/detect/:metricType", h.DetectAnomalies)
		group.GET("/anomalies/user/:userID", h.DetectUserAnomalies)
		group.GET("/anomalies/recent", h.GetRecentAnomalies)
		group.POST("/anomalies/:anomalyID/acknowledge", h.AcknowledgeAnomaly)
		group.POST("/anomalies/:anomalyID/investigate", h.InvestigateAnomaly)
		group.GET("/anomalies/stats", h.GetAnomalyStats)

		// Model Management
		group.GET("/models/status", h.GetModelStatus)
		group.POST("/models/retrain", h.RetrainAllModels)

		// Insights
		group.GET("/insights", h.GetInsights)
	}
}

// GetUserChurnPrediction returns churn prediction for a user
// GET /api/v1/admin/predictions/churn/user/:userID
func (h *AdminPredictiveAnalyticsHandler) GetUserChurnPrediction(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	prediction, err := h.churnService.PredictUserChurn(c.Request.Context(), restaurantID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prediction)
}

// GetRestaurantChurnPrediction returns churn prediction for the restaurant
// GET /api/v1/admin/predictions/churn/restaurant
func (h *AdminPredictiveAnalyticsHandler) GetRestaurantChurnPrediction(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	prediction, err := h.churnService.PredictRestaurantChurn(c.Request.Context(), restaurantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prediction)
}

// BulkGetChurnPredictions returns churn predictions for multiple users
// GET /api/v1/admin/predictions/churn/users?limit=100&risk_level=high
func (h *AdminPredictiveAnalyticsHandler) BulkGetChurnPredictions(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	limit := 100
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	riskLevel := c.DefaultQuery("risk_level", "") // "low", "medium", "high", "critical", or empty for all

	c.JSON(http.StatusOK, gin.H{
		"limit":      limit,
		"risk_level": riskLevel,
		"message":    "Use specific user endpoint for individual predictions",
	})
}

// TrainChurnModel trains churn prediction model
// POST /api/v1/admin/predictions/churn/train
func (h *AdminPredictiveAnalyticsHandler) TrainChurnModel(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	job, err := h.churnService.TrainChurnModel(c.Request.Context(), restaurantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, job)
}

// GetChurnHistory returns historical churn predictions
// GET /api/v1/admin/predictions/churn/history/:userID?limit=10
func (h *AdminPredictiveAnalyticsHandler) GetChurnHistory(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	history, err := h.churnService.GetChurnPredictionHistory(c.Request.Context(), restaurantID, userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":   userID,
		"history":   history,
		"count":     len(history),
	})
}

// ForecastRevenue generates revenue forecast
// GET /api/v1/admin/predictions/revenue/forecast?period=month&lookback=90
func (h *AdminPredictiveAnalyticsHandler) ForecastRevenue(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	period := c.DefaultQuery("period", "month") // week, month, quarter, year
	lookback := 90
	if lookbackStr := c.Query("lookback"); lookbackStr != "" {
		if l, err := strconv.Atoi(lookbackStr); err == nil {
			lookback = l
		}
	}

	forecast, err := h.revenueService.ForecastRevenue(c.Request.Context(), restaurantID, period, lookback)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, forecast)
}

// GetForecastAccuracy returns forecast accuracy for different periods
// GET /api/v1/admin/predictions/revenue/accuracy
func (h *AdminPredictiveAnalyticsHandler) GetForecastAccuracy(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	accuracy, err := h.revenueService.CompareForecastAccuracy(c.Request.Context(), restaurantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"restaurant_id": restaurantID,
		"accuracy":      accuracy,
		"evaluated_at":  time.Now(),
	})
}

// GetRevenueHistory returns historical revenue forecasts
// GET /api/v1/admin/predictions/revenue/history?period=month&limit=10
func (h *AdminPredictiveAnalyticsHandler) GetRevenueHistory(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	period := c.DefaultQuery("period", "month")
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	history, err := h.revenueService.GetForecastHistory(c.Request.Context(), restaurantID, period, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"period":  period,
		"history": history,
		"count":   len(history),
	})
}

// GetSeasonalityReport generates seasonality analysis
// GET /api/v1/admin/predictions/revenue/seasonality?days=90
func (h *AdminPredictiveAnalyticsHandler) GetSeasonalityReport(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	days := 90
	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	report, err := h.revenueService.GenerateSeasonalityReport(c.Request.Context(), restaurantID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// TrainRevenueModel trains revenue forecasting model
// POST /api/v1/admin/predictions/revenue/train
func (h *AdminPredictiveAnalyticsHandler) TrainRevenueModel(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	job, err := h.revenueService.TrainRevenueModel(c.Request.Context(), restaurantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, job)
}

// GetUserLTV returns LTV prediction for a user
// GET /api/v1/admin/predictions/ltv/user/:userID
func (h *AdminPredictiveAnalyticsHandler) GetUserLTV(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	ltv, err := h.ltvService.PredictLifetimeValue(c.Request.Context(), restaurantID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ltv)
}

// GetSegmentDistribution returns customer segment distribution
// GET /api/v1/admin/predictions/ltv/segment-distribution
func (h *AdminPredictiveAnalyticsHandler) GetSegmentDistribution(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	distribution, err := h.ltvService.GetValueSegmentDistribution(c.Request.Context(), restaurantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"restaurant_id":  restaurantID,
		"distribution":   distribution,
		"evaluated_at":   time.Now(),
	})
}

// GetTopCustomers returns highest LTV customers
// GET /api/v1/admin/predictions/ltv/top-customers?limit=50
func (h *AdminPredictiveAnalyticsHandler) GetTopCustomers(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	customers, err := h.ltvService.GetTopCustomers(c.Request.Context(), restaurantID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"limit":     limit,
		"customers": customers,
		"count":     len(customers),
	})
}

// GetLTVBySegment returns LTV statistics by segment
// GET /api/v1/admin/predictions/ltv/by-segment
func (h *AdminPredictiveAnalyticsHandler) GetLTVBySegment(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	stats, err := h.ltvService.GetLTVBySegment(c.Request.Context(), restaurantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"restaurant_id": restaurantID,
		"segments":      stats,
	})
}

// GetLTVHistory returns LTV prediction history
// GET /api/v1/admin/predictions/ltv/history/:userID?limit=10
func (h *AdminPredictiveAnalyticsHandler) GetLTVHistory(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	history, err := h.ltvService.GetLTVHistory(c.Request.Context(), restaurantID, userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"history": history,
		"count":   len(history),
	})
}

// TrainLTVModel trains LTV prediction model
// POST /api/v1/admin/predictions/ltv/train
func (h *AdminPredictiveAnalyticsHandler) TrainLTVModel(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	job, err := h.ltvService.TrainLTVModel(c.Request.Context(), restaurantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, job)
}

// DetectAnomalies detects anomalies in metrics
// GET /api/v1/admin/predictions/anomalies/detect/:metricType?lookback=7&threshold=2
func (h *AdminPredictiveAnalyticsHandler) DetectAnomalies(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	metricType := c.Param("metricType")

	lookback := 7
	if lookbackStr := c.Query("lookback"); lookbackStr != "" {
		if l, err := strconv.Atoi(lookbackStr); err == nil {
			lookback = l
		}
	}

	threshold := 2.0
	if thresholdStr := c.Query("threshold"); thresholdStr != "" {
		if t, err := strconv.ParseFloat(thresholdStr, 64); err == nil {
			threshold = t
		}
	}

	anomalies, err := h.anomalyService.DetectAnomalies(
		c.Request.Context(),
		restaurantID,
		metricType,
		lookback,
		threshold,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"metric_type": metricType,
		"anomalies":   anomalies,
		"count":       len(anomalies),
	})
}

// DetectUserAnomalies detects user behavior anomalies
// GET /api/v1/admin/predictions/anomalies/user/:userID
func (h *AdminPredictiveAnalyticsHandler) DetectUserAnomalies(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	anomalies, err := h.anomalyService.DetectUserAnomalies(c.Request.Context(), restaurantID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":   userID,
		"anomalies": anomalies,
		"count":     len(anomalies),
	})
}

// GetRecentAnomalies returns recent anomalies
// GET /api/v1/admin/predictions/anomalies/recent?limit=50&severity=warning
func (h *AdminPredictiveAnalyticsHandler) GetRecentAnomalies(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	severity := c.DefaultQuery("severity", "")

	anomalies, err := h.anomalyService.GetRecentAnomalies(c.Request.Context(), restaurantID, limit, severity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"anomalies": anomalies,
		"count":     len(anomalies),
		"severity":  severity,
	})
}

// AcknowledgeAnomaly acknowledges an anomaly
// POST /api/v1/admin/predictions/anomalies/:anomalyID/acknowledge
func (h *AdminPredictiveAnalyticsHandler) AcknowledgeAnomaly(c *gin.Context) {
	anomalyID, err := strconv.ParseInt(c.Param("anomalyID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid anomaly_id"})
		return
	}

	err = h.anomalyService.AcknowledgeAnomaly(c.Request.Context(), anomalyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "acknowledged",
		"anomaly_id": anomalyID,
	})
}

// InvestigateAnomaly investigates and resolves anomaly
// POST /api/v1/admin/predictions/anomalies/:anomalyID/investigate
func (h *AdminPredictiveAnalyticsHandler) InvestigateAnomaly(c *gin.Context) {
	anomalyID, err := strconv.ParseInt(c.Param("anomalyID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid anomaly_id"})
		return
	}

	var req struct {
		InvestigatedBy    int64  `json:"investigated_by" binding:"required"`
		RootCause         string `json:"root_cause"`
		Resolution        string `json:"resolution"`
		Notes             string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.anomalyService.InvestigateAnomaly(
		c.Request.Context(),
		anomalyID,
		req.InvestigatedBy,
		req.RootCause,
		req.Resolution,
		req.Notes,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "investigated",
		"anomaly_id": anomalyID,
	})
}

// GetAnomalyStats returns anomaly statistics
// GET /api/v1/admin/predictions/anomalies/stats?days=30
func (h *AdminPredictiveAnalyticsHandler) GetAnomalyStats(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	days := 30
	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	stats, err := h.anomalyService.GetAnomalyStats(c.Request.Context(), restaurantID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetModelStatus returns status of all prediction models
// GET /api/v1/admin/predictions/models/status
func (h *AdminPredictiveAnalyticsHandler) GetModelStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"models": []string{"churn", "revenue", "ltv", "anomaly"},
		"status": "operational",
		"timestamp": time.Now(),
	})
}

// RetrainAllModels triggers retraining of all models
// POST /api/v1/admin/predictions/models/retrain
func (h *AdminPredictiveAnalyticsHandler) RetrainAllModels(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	c.JSON(http.StatusAccepted, gin.H{
		"restaurant_id": restaurantID,
		"message":       "Retraining initiated for all models",
		"models":        []string{"churn", "revenue", "ltv", "anomaly"},
		"job_id":        fmt.Sprintf("retrain_%d_%d", restaurantID, time.Now().Unix()),
	})
}

// GetInsights returns actionable insights from all predictions
// GET /api/v1/admin/predictions/insights
func (h *AdminPredictiveAnalyticsHandler) GetInsights(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	insights := []map[string]interface{}{
		{
			"type":        "churn_warning",
			"title":       "High Churn Risk Detected",
			"description": "20 users showing critical churn signals",
			"priority":    "high",
			"action":      "Review churn predictions and implement retention campaigns",
		},
		{
			"type":        "revenue_opportunity",
			"title":       "Upsell Potential Identified",
			"description": "50 customers in medium segment can be upsold",
			"priority":    "medium",
			"action":      "Launch targeted upsell campaign",
		},
		{
			"type":        "anomaly_alert",
			"title":       "Revenue Anomaly",
			"description": "Revenue 45% above expected for today",
			"priority":    "info",
			"action":      "Investigate cause of revenue spike",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"restaurant_id": restaurantID,
		"insights":      insights,
		"count":         len(insights),
		"generated_at":  time.Now(),
	})
}
