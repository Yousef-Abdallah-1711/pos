package service

import (
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

// MockPaymobClient is used for local development and tests.
type MockPaymobClient struct{}

// NewMockPaymobClient creates a mock Paymob client.
func NewMockPaymobClient() PaymobClient {
	return &MockPaymobClient{}
}

// Authenticate returns a deterministic mock auth token.
func (c *MockPaymobClient) Authenticate() (string, error) {
	return fmt.Sprintf("mock_auth_token_%d", time.Now().UnixNano()), nil
}

// CreateOrder returns a synthetic positive order ID.
func (c *MockPaymobClient) CreateOrder(authToken string, amountCents int, currency string, merchantOrderID string) (int64, error) {
	// Keep IDs within JS-safe range for frontend/Postman compatibility.
	return time.Now().UnixMilli()*1000 + int64(rand.Intn(999)), nil
}

// GetPaymentKey returns a synthetic payment key.
func (c *MockPaymobClient) GetPaymentKey(authToken string, orderID int64, amountCents int, currency string, integrationID int, billingData map[string]string, expiration int) (string, error) {
	return fmt.Sprintf(
		"mockpk.%d.%d.%s.%d",
		orderID,
		amountCents,
		strings.ToUpper(strings.TrimSpace(currency)),
		time.Now().UnixNano(),
	), nil
}

// RefundTransaction always succeeds in mock mode.
func (c *MockPaymobClient) RefundTransaction(authToken string, transactionID int64, amountCents int) error {
	return nil
}

// GetIframeURL returns the mock iframe endpoint.
func (c *MockPaymobClient) GetIframeURL(paymentKey string) string {
	return fmt.Sprintf("/api/v1/mock/paymob/iframe?payment_token=%s", url.QueryEscape(paymentKey))
}

// VerifyHMAC always succeeds in mock mode.
func (c *MockPaymobClient) VerifyHMAC(payload []byte, receivedHMAC string) bool {
	return true
}
