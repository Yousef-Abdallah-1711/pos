package ws

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	kdsWriteWait      = 10 * time.Second
	kdsPongWait       = 60 * time.Second
	kdsPingPeriod     = (kdsPongWait * 9) / 10 // 54 seconds
	kdsMaxMessageSize = 1024
)

type kdsBroadcastMessage struct {
	tenantID     int64
	restaurantID int64
	payload      []byte
}

type kdsEnvelope struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// KDSClient represents one websocket connection in a tenant+restaurant room.
type KDSClient struct {
	hub            *KDSHub
	conn           *websocket.Conn
	send           chan []byte
	tenantID       int64
	restaurantID   int64
	userID         int64
	onDisconnect   func()
	disconnectOnce sync.Once
}

func newKDSClient(
	hub *KDSHub,
	conn *websocket.Conn,
	tenantID, restaurantID, userID int64,
	onDisconnect func(),
) *KDSClient {
	return &KDSClient{
		hub:          hub,
		conn:         conn,
		send:         make(chan []byte, 32),
		tenantID:     tenantID,
		restaurantID: restaurantID,
		userID:       userID,
		onDisconnect: onDisconnect,
	}
}

func (c *KDSClient) roomKey() string {
	return buildKDSRoomKey(c.tenantID, c.restaurantID)
}

// ReadPump only keeps connection liveness and unregisters on disconnect.
func (c *KDSClient) ReadPump() {
	defer func() {
		c.disconnectOnce.Do(func() {
			if c.onDisconnect != nil {
				c.onDisconnect()
			}
		})
		c.hub.Unregister(c)
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(kdsMaxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(kdsPongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(kdsPongWait))
	})

	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[KDS WS] unexpected close error: %v", err)
			}
			return
		}
	}
}

// WritePump writes outbound messages and periodic ping frames.
func (c *KDSClient) WritePump() {
	ticker := time.NewTicker(kdsPingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(kdsWriteWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			if _, err := writer.Write(message); err != nil {
				_ = writer.Close()
				return
			}

			queued := len(c.send)
			for i := 0; i < queued; i++ {
				if _, err := writer.Write([]byte("\n")); err != nil {
					_ = writer.Close()
					return
				}
				next := <-c.send
				if _, err := writer.Write(next); err != nil {
					_ = writer.Close()
					return
				}
			}

			if err := writer.Close(); err != nil {
				return
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(kdsWriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// KDSHub manages websocket clients by tenant+restaurant room.
type KDSHub struct {
	rooms      map[string]map[*KDSClient]struct{}
	register   chan *KDSClient
	unregister chan *KDSClient
	broadcast  chan kdsBroadcastMessage
	shutdown   chan struct{}
	done       chan struct{}
	stopOnce   sync.Once
}

func NewKDSHub() *KDSHub {
	return &KDSHub{
		rooms:      make(map[string]map[*KDSClient]struct{}),
		register:   make(chan *KDSClient, 128),
		unregister: make(chan *KDSClient, 128),
		broadcast:  make(chan kdsBroadcastMessage, 512),
		shutdown:   make(chan struct{}),
		done:       make(chan struct{}),
	}
}

func buildKDSRoomKey(tenantID, restaurantID int64) string {
	return "tenant:" + strconv.FormatInt(tenantID, 10) + ":restaurant:" + strconv.FormatInt(restaurantID, 10)
}

// Run starts the hub loop.
func (h *KDSHub) Run() {
	defer close(h.done)

	for {
		select {
		case client := <-h.register:
			room := client.roomKey()
			if h.rooms[room] == nil {
				h.rooms[room] = make(map[*KDSClient]struct{})
			}
			h.rooms[room][client] = struct{}{}

		case client := <-h.unregister:
			h.removeClient(client)

		case message := <-h.broadcast:
			room := buildKDSRoomKey(message.tenantID, message.restaurantID)
			clients := h.rooms[room]
			for client := range clients {
				select {
				case client.send <- message.payload:
				default:
					// Slow consumer, drop connection to protect hub throughput.
					h.removeClient(client)
				}
			}

		case <-h.shutdown:
			for room, clients := range h.rooms {
				for client := range clients {
					delete(clients, client)
					close(client.send)
					_ = client.conn.Close()
				}
				delete(h.rooms, room)
			}
			return
		}
	}
}

func (h *KDSHub) removeClient(client *KDSClient) {
	room := client.roomKey()
	clients, ok := h.rooms[room]
	if !ok {
		return
	}
	if _, exists := clients[client]; !exists {
		return
	}

	delete(clients, client)
	close(client.send)

	if len(clients) == 0 {
		delete(h.rooms, room)
	}
}

func (h *KDSHub) Register(client *KDSClient) {
	select {
	case h.register <- client:
	case <-h.done:
		_ = client.conn.Close()
	}
}

func (h *KDSHub) Unregister(client *KDSClient) {
	select {
	case h.unregister <- client:
	case <-h.done:
	}
}

// BroadcastKDSUpdate implements usecase.KDSBroadcaster.
func (h *KDSHub) BroadcastKDSUpdate(ctx context.Context, tenantID, restaurantID int64, messageType string, payload interface{}) error {
	if tenantID <= 0 || restaurantID <= 0 {
		return errors.New("invalid tenant context for broadcast")
	}

	envelope := kdsEnvelope{
		Type:      messageType,
		Data:      payload,
		Timestamp: time.Now().UTC(),
	}

	encoded, err := json.Marshal(envelope)
	if err != nil {
		return err
	}

	message := kdsBroadcastMessage{
		tenantID:     tenantID,
		restaurantID: restaurantID,
		payload:      encoded,
	}

	select {
	case h.broadcast <- message:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return errors.New("kds broadcast queue is full")
	}
}

// Shutdown gracefully stops the hub loop and all connected clients.
func (h *KDSHub) Shutdown(ctx context.Context) error {
	h.stopOnce.Do(func() {
		close(h.shutdown)
	})

	select {
	case <-h.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
