package ws

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func TestKDSHubBroadcastRoutesMessageToMatchingRoom(t *testing.T) {
	t.Parallel()

	hub := NewKDSHub()
	go hub.Run()

	sameRoom := &KDSClient{
		hub:          hub,
		send:         make(chan []byte, 4),
		tenantID:     1,
		restaurantID: 10,
	}
	otherRoom := &KDSClient{
		hub:          hub,
		send:         make(chan []byte, 4),
		tenantID:     1,
		restaurantID: 11,
	}

	hub.Register(sameRoom)
	hub.Register(otherRoom)

	if err := hub.BroadcastKDSUpdate(context.Background(), 1, 10, "order_updated", map[string]interface{}{"id": 77}); err != nil {
		t.Fatalf("BroadcastKDSUpdate returned error: %v", err)
	}

	select {
	case payload := <-sameRoom.send:
		var envelope kdsEnvelope
		if err := json.Unmarshal(payload, &envelope); err != nil {
			t.Fatalf("failed to decode broadcast payload: %v", err)
		}
		if envelope.Type != "order_updated" {
			t.Fatalf("expected message type order_updated, got %s", envelope.Type)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for matching-room broadcast payload")
	}

	select {
	case <-otherRoom.send:
		t.Fatal("did not expect payload for non-matching room")
	case <-time.After(250 * time.Millisecond):
		// expected: no payload
	}

	// Unregister before shutdown to avoid nil websocket connection close on synthetic clients.
	hub.Unregister(sameRoom)
	hub.Unregister(otherRoom)
	time.Sleep(50 * time.Millisecond)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := hub.Shutdown(shutdownCtx); err != nil {
		t.Fatalf("Shutdown returned error: %v", err)
	}
}

func TestKDSHubBroadcastValidationAndQueueBackpressure(t *testing.T) {
	t.Parallel()

	hub := NewKDSHub()

	if err := hub.BroadcastKDSUpdate(context.Background(), 0, 1, "order_new", map[string]interface{}{}); err == nil {
		t.Fatal("expected invalid tenant context error")
	}

	for i := 0; i < cap(hub.broadcast); i++ {
		hub.broadcast <- kdsBroadcastMessage{
			tenantID:     1,
			restaurantID: 1,
			payload:      []byte(`{"type":"filled"}`),
		}
	}

	err := hub.BroadcastKDSUpdate(context.Background(), 1, 1, "order_new", map[string]interface{}{"id": 1})
	if err == nil {
		t.Fatal("expected queue full error when broadcast channel is saturated")
	}
}
