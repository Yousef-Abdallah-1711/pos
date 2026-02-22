//go:build legacy
// +build legacy

package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"pos-saas/internal/service"
)

// CollaborationWebSocketHandler handles WebSocket connections for real-time collaboration
type CollaborationWebSocketHandler struct {
	sessionManager   *service.CollaborationSessionManager
	otService        *service.OperationalTransformationService
	commentService   *service.CollaborationCommentService
	upgrader         websocket.Upgrader
}

// NewCollaborationWebSocketHandler creates a new WebSocket handler
func NewCollaborationWebSocketHandler(
	sessionManager *service.CollaborationSessionManager,
	otService *service.OperationalTransformationService,
	commentService *service.CollaborationCommentService,
) *CollaborationWebSocketHandler {
	return &CollaborationWebSocketHandler{
		sessionManager: sessionManager,
		otService:      otService,
		commentService: commentService,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // In production, validate origin
			},
		},
	}
}

// RegisterRoutes registers WebSocket routes
func (h *CollaborationWebSocketHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/ws/collaborate/:resourceID", h.HandleWebSocket)
	router.POST("/api/v1/collaboration/sessions", h.CreateSession)
	router.GET("/api/v1/collaboration/sessions/:sessionID", h.GetSession)
	router.POST("/api/v1/collaboration/sessions/:sessionID/leave", h.LeaveSession)
	router.GET("/api/v1/collaboration/comments/:resourceID", h.GetComments)
	router.POST("/api/v1/collaboration/comments", h.AddComment)
	router.POST("/api/v1/collaboration/comments/:commentID/resolve", h.ResolveComment)
}

// HandleWebSocket handles WebSocket connections
func (h *CollaborationWebSocketHandler) HandleWebSocket(c *gin.Context) {
	resourceIDStr := c.Param("resourceID")
	resourceID, err := strconv.ParseInt(resourceIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource_id"})
		return
	}

	resourceType := c.DefaultQuery("type", "component")
	restaurantID := c.GetInt64("restaurant_id")
	userID := c.GetInt64("user_id")
	username := c.GetString("username")
	sessionIDQuery := c.Query("session_id")

	// Upgrade connection
	ws, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer ws.Close()

	// Get or create session
	sessionID := sessionIDQuery
	if sessionID == "" {
		var createErr error
		sessionID, createErr = h.sessionManager.CreateSession(
			context.Background(),
			restaurantID,
			resourceID,
			resourceType,
			"",
		)
		if createErr != nil {
			ws.WriteJSON(map[string]interface{}{
				"type":  "error",
				"error": "failed to create session",
			})
			return
		}
	}

	// Join session
	updatesChan, err := h.sessionManager.JoinSession(
		context.Background(),
		sessionID,
		userID,
		username,
		generateUserColor(userID),
	)
	if err != nil {
		ws.WriteJSON(map[string]interface{}{
			"type":  "error",
			"error": err.Error(),
		})
		return
	}

	// Send session info
	ws.WriteJSON(map[string]interface{}{
		"type":       "session_info",
		"session_id": sessionID,
		"user_id":    userID,
		"username":   username,
	})

	// Handle messages
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go h.readPump(ws, sessionID, resourceID, resourceType, userID, username, ctx)
	h.writePump(ws, sessionID, updatesChan, ctx)
}

// readPump reads messages from WebSocket
func (h *CollaborationWebSocketHandler) readPump(
	ws *websocket.Conn,
	sessionID string,
	resourceID int64,
	resourceType string,
	userID int64,
	username string,
	ctx context.Context,
) {
	defer ws.Close()

	ws.SetReadDeadline(time.Now().Add(60 * time.Second))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var message map[string]interface{}
		err := ws.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			return
		}

		msgType := message["type"].(string)

		switch msgType {
		case "edit":
			h.handleEdit(ctx, ws, sessionID, resourceID, resourceType, userID, message)

		case "cursor":
			h.handleCursor(ctx, ws, sessionID, userID, message)

		case "comment":
			h.handleComment(ctx, ws, sessionID, resourceID, resourceType, userID, username, message)

		case "ping":
			ws.WriteJSON(map[string]string{"type": "pong"})

		default:
			log.Printf("Unknown message type: %s", msgType)
		}
	}
}

// handleEdit processes edit operations
func (h *CollaborationWebSocketHandler) handleEdit(
	ctx context.Context,
	ws *websocket.Conn,
	sessionID string,
	resourceID int64,
	resourceType string,
	userID int64,
	message map[string]interface{},
) {
	operation := &service.Operation{
		Type:     message["operation"].(string),
		Position: int32(message["position"].(float64)),
		Content:  message["content"].(string),
		Length:   int32(message["length"].(float64)),
	}

	// Apply operation
	session, _ := h.sessionManager.GetSession(ctx, sessionID)
	newDoc, err := h.otService.ApplyOperation(ctx, sessionID, userID, operation, session.Document)
	if err != nil {
		ws.WriteJSON(map[string]interface{}{
			"type":  "error",
			"error": err.Error(),
		})
		return
	}

	// Update session document
	session.Document = newDoc
	session.VersionNumber++

	// Broadcast edit
	h.broadcastEdit(sessionID, userID, operation, session.VersionNumber)
}

// handleCursor processes cursor position updates
func (h *CollaborationWebSocketHandler) handleCursor(
	ctx context.Context,
	ws *websocket.Conn,
	sessionID string,
	userID int64,
	message map[string]interface{},
) {
	position := int32(message["position"].(float64))
	line := int32(message["line"].(float64))
	column := int32(message["column"].(float64))

	h.sessionManager.UpdateCursor(ctx, sessionID, userID, position, line, column)
}

// handleComment processes comments
func (h *CollaborationWebSocketHandler) handleComment(
	ctx context.Context,
	ws *websocket.Conn,
	sessionID string,
	resourceID int64,
	resourceType string,
	userID int64,
	username string,
	message map[string]interface{},
) {
	position := int32(message["position"].(float64))
	lineNumber := int32(message["line_number"].(float64))
	text := message["text"].(string)

	comment, err := h.commentService.AddComment(
		ctx,
		sessionID,
		resourceID,
		resourceType,
		userID,
		username,
		position,
		lineNumber,
		text,
	)
	if err != nil {
		ws.WriteJSON(map[string]interface{}{
			"type":  "error",
			"error": err.Error(),
		})
		return
	}

	// Broadcast comment
	h.broadcastComment(sessionID, comment)
}

// writePump writes messages to WebSocket
func (h *CollaborationWebSocketHandler) writePump(
	ws *websocket.Conn,
	sessionID string,
	updatesChan <-chan interface{},
	ctx context.Context,
) {
	ticker := time.NewTicker(54 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-updatesChan:
			ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := ws.WriteJSON(message); err != nil {
				return
			}

		case <-ticker.C:
			ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

// broadcastEdit broadcasts an edit to all users
func (h *CollaborationWebSocketHandler) broadcastEdit(
	sessionID string,
	userID int64,
	operation *service.Operation,
	version int64,
) {
	message := map[string]interface{}{
		"type":      "edit",
		"user_id":   userID,
		"operation": operation.Type,
		"position":  operation.Position,
		"content":   operation.Content,
		"length":    operation.Length,
		"version":   version,
		"timestamp": time.Now(),
	}

	h.sessionManager.broadcastToSession(sessionID, message, userID)
}

// broadcastComment broadcasts a comment to all users
func (h *CollaborationWebSocketHandler) broadcastComment(
	sessionID string,
	comment *models.Comment,
) {
	message := map[string]interface{}{
		"type":    "comment",
		"comment": comment,
	}

	h.sessionManager.broadcastToSession(sessionID, message, 0)
}

// CreateSession creates a new collaboration session
// POST /api/v1/collaboration/sessions
func (h *CollaborationWebSocketHandler) CreateSession(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	var req struct {
		ResourceID   int64  `json:"resource_id" binding:"required"`
		ResourceType string `json:"resource_type" binding:"required"`
		InitialDoc   string `json:"initial_doc,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionID, err := h.sessionManager.CreateSession(
		c.Request.Context(),
		restaurantID,
		req.ResourceID,
		req.ResourceType,
		req.InitialDoc,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"session_id": sessionID,
		"ws_url":     fmt.Sprintf("ws://localhost:8080/ws/collaborate/%d?session_id=%s&type=%s", req.ResourceID, sessionID, req.ResourceType),
	})
}

// GetSession retrieves session information
// GET /api/v1/collaboration/sessions/:sessionID
func (h *CollaborationWebSocketHandler) GetSession(c *gin.Context) {
	sessionID := c.Param("sessionID")

	session, err := h.sessionManager.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	participants, _ := h.sessionManager.GetParticipants(sessionID)

	c.JSON(http.StatusOK, gin.H{
		"session_id":   session.ID,
		"resource_id":  session.ResourceID,
		"resource_type": session.ResourceType,
		"participants": participants,
		"version":      session.VersionNumber,
	})
}

// LeaveSession removes user from session
// POST /api/v1/collaboration/sessions/:sessionID/leave
func (h *CollaborationWebSocketHandler) LeaveSession(c *gin.Context) {
	sessionID := c.Param("sessionID")
	userID := c.GetInt64("user_id")

	if err := h.sessionManager.LeaveSession(c.Request.Context(), sessionID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "left"})
}

// GetComments retrieves comments for a resource
// GET /api/v1/collaboration/comments/:resourceID
func (h *CollaborationWebSocketHandler) GetComments(c *gin.Context) {
	resourceIDStr := c.Param("resourceID")
	resourceID, _ := strconv.ParseInt(resourceIDStr, 10, 64)
	resourceType := c.DefaultQuery("type", "component")
	unresolvedOnly := c.DefaultQuery("unresolved", "false") == "true"

	comments, err := h.commentService.GetComments(
		c.Request.Context(),
		resourceID,
		resourceType,
		unresolvedOnly,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"count":    len(comments),
	})
}

// AddComment adds a new comment
// POST /api/v1/collaboration/comments
func (h *CollaborationWebSocketHandler) AddComment(c *gin.Context) {
	userID := c.GetInt64("user_id")
	username := c.GetString("username")

	var req struct {
		SessionID    string `json:"session_id" binding:"required"`
		ResourceID   int64  `json:"resource_id" binding:"required"`
		ResourceType string `json:"resource_type" binding:"required"`
		Position     int32  `json:"position" binding:"required"`
		LineNumber   int32  `json:"line_number"`
		Text         string `json:"text" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := h.commentService.AddComment(
		c.Request.Context(),
		req.SessionID,
		req.ResourceID,
		req.ResourceType,
		userID,
		username,
		req.Position,
		req.LineNumber,
		req.Text,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// ResolveComment marks a comment as resolved
// POST /api/v1/collaboration/comments/:commentID/resolve
func (h *CollaborationWebSocketHandler) ResolveComment(c *gin.Context) {
	commentIDStr := c.Param("commentID")
	commentID, _ := strconv.ParseInt(commentIDStr, 10, 64)
	userID := c.GetInt64("user_id")

	if err := h.commentService.ResolveComment(c.Request.Context(), commentID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "resolved"})
}

// Helper function to generate user colors
func generateUserColor(userID int64) string {
	colors := []string{
		"#FF6B6B", "#4ECDC4", "#45B7D1", "#FFA07A", "#98D8C8",
		"#F7DC6F", "#BB8FCE", "#85C1E2", "#F8B739", "#52C9A8",
	}
	return colors[userID%int64(len(colors))]
}
