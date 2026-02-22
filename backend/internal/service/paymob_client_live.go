package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"pos-saas/internal/domain"
)

var paymobHMACFieldOrder = []string{
	"amount_cents",
	"created_at",
	"currency",
	"error_occured",
	"has_parent_transaction",
	"id",
	"integration_id",
	"is_3d_secure",
	"is_auth",
	"is_capture",
	"is_refunded",
	"is_standalone_payment",
	"is_voided",
	"order.id",
	"owner",
	"pending",
	"source_data.pan",
	"source_data.sub_type",
	"source_data.type",
	"success",
}

// RealPaymobClient calls Paymob APIs directly.
type RealPaymobClient struct {
	apiKey     string
	iframeID   int
	hmacSecret string
	baseURL    string
	httpClient *http.Client
}

// NewRealPaymobClient creates a real Paymob API client.
func NewRealPaymobClient(apiKey string, iframeID int, hmacSecret string, baseURL string) PaymobClient {
	base := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if base == "" {
		base = "https://accept.paymob.com/api"
	}

	return &RealPaymobClient{
		apiKey:     apiKey,
		iframeID:   iframeID,
		hmacSecret: hmacSecret,
		baseURL:    base,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Authenticate obtains a Paymob auth token.
func (c *RealPaymobClient) Authenticate() (string, error) {
	payload := map[string]string{
		"api_key": c.apiKey,
	}

	var resp domain.PaymobAuthResponse
	if err := c.postJSON(c.baseURL+"/auth/tokens", payload, nil, &resp); err != nil {
		return "", err
	}

	if strings.TrimSpace(resp.Token) == "" {
		return "", fmt.Errorf("paymob auth token missing in response")
	}

	return resp.Token, nil
}

// CreateOrder registers an order on Paymob.
func (c *RealPaymobClient) CreateOrder(authToken string, amountCents int, currency string, merchantOrderID string) (int64, error) {
	payload := map[string]interface{}{
		"auth_token": authToken,
		// NOTE: This integration currently sends Paymob boolean flags as strings.
		// If live API rejects type, switch to native bool (false/true).
		"delivery_needed":   "false",
		"amount_cents":      strconv.Itoa(amountCents),
		"currency":          strings.ToUpper(currency),
		"merchant_order_id": merchantOrderID,
		"items": []map[string]string{
			{
				"name":         "POS SaaS Monthly Subscription",
				"amount_cents": strconv.Itoa(amountCents),
				"quantity":     "1",
			},
		},
	}

	var resp domain.PaymobOrderResponse
	if err := c.postJSON(c.baseURL+"/ecommerce/orders", payload, nil, &resp); err != nil {
		return 0, err
	}

	if resp.ID <= 0 {
		return 0, fmt.Errorf("invalid paymob order id")
	}

	return resp.ID, nil
}

// GetPaymentKey creates a payment key for iframe checkout.
func (c *RealPaymobClient) GetPaymentKey(authToken string, orderID int64, amountCents int, currency string, integrationID int, billingData map[string]string, expiration int) (string, error) {
	payload := map[string]interface{}{
		"auth_token":   authToken,
		"amount_cents": strconv.Itoa(amountCents),
		"expiration":   expiration,
		"order_id":     orderID,
		"billing_data": map[string]string{
			"email":           getMapValueOrDefault(billingData, "email", "unknown@example.com"),
			"first_name":      getMapValueOrDefault(billingData, "first_name", "NA"),
			"last_name":       getMapValueOrDefault(billingData, "last_name", "NA"),
			"phone_number":    getMapValueOrDefault(billingData, "phone_number", "+200000000000"),
			"apartment":       getMapValueOrDefault(billingData, "apartment", "NA"),
			"floor":           getMapValueOrDefault(billingData, "floor", "NA"),
			"building":        getMapValueOrDefault(billingData, "building", "NA"),
			"street":          getMapValueOrDefault(billingData, "street", "NA"),
			"city":            getMapValueOrDefault(billingData, "city", "Cairo"),
			"state":           getMapValueOrDefault(billingData, "state", "Cairo"),
			"country":         getMapValueOrDefault(billingData, "country", "EG"),
			"shipping_method": getMapValueOrDefault(billingData, "shipping_method", "NA"),
			"postal_code":     getMapValueOrDefault(billingData, "postal_code", "NA"),
		},
		"currency":       strings.ToUpper(currency),
		"integration_id": integrationID,
		// NOTE: Keep as string for compatibility with current Paymob payload style.
		"lock_order_when_paid": "true",
	}

	var resp domain.PaymobPaymentKeyResponse
	if err := c.postJSON(c.baseURL+"/acceptance/payment_keys", payload, nil, &resp); err != nil {
		return "", err
	}

	if strings.TrimSpace(resp.Token) == "" {
		return "", fmt.Errorf("paymob payment key missing in response")
	}

	return resp.Token, nil
}

// RefundTransaction requests a refund for a captured transaction.
func (c *RealPaymobClient) RefundTransaction(authToken string, transactionID int64, amountCents int) error {
	if strings.TrimSpace(authToken) == "" {
		return fmt.Errorf("paymob auth token is required")
	}
	if transactionID <= 0 {
		return fmt.Errorf("invalid paymob transaction id")
	}
	if amountCents <= 0 {
		return fmt.Errorf("invalid refund amount")
	}

	payload := map[string]interface{}{
		"auth_token":     authToken,
		"transaction_id": transactionID,
		"amount_cents":   amountCents,
	}

	return c.postJSON(c.baseURL+"/acceptance/void_refund/refund", payload, nil, nil)
}

// GetIframeURL returns the hosted Paymob iframe URL.
func (c *RealPaymobClient) GetIframeURL(paymentKey string) string {
	base := strings.TrimRight(c.baseURL, "/")
	return fmt.Sprintf("%s/acceptance/iframes/%d?payment_token=%s", base, c.iframeID, url.QueryEscape(paymentKey))
}

// VerifyHMAC validates webhook authenticity using Paymob's required field order.
func (c *RealPaymobClient) VerifyHMAC(payload []byte, receivedHMAC string) bool {
	if strings.TrimSpace(receivedHMAC) == "" || strings.TrimSpace(c.hmacSecret) == "" {
		return false
	}

	var webhook domain.PaymobWebhookPayload
	if err := json.Unmarshal(payload, &webhook); err != nil {
		return false
	}

	concatenated := buildPaymobHMACMessage(&webhook)
	mac := hmac.New(sha512.New, []byte(c.hmacSecret))
	_, _ = mac.Write([]byte(concatenated))
	computed := hex.EncodeToString(mac.Sum(nil))

	expected := []byte(strings.ToLower(strings.TrimSpace(computed)))
	actual := []byte(strings.ToLower(strings.TrimSpace(receivedHMAC)))
	return subtle.ConstantTimeCompare(expected, actual) == 1
}

func (c *RealPaymobClient) postJSON(url string, payload interface{}, headers map[string]string, out interface{}) error {
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("paymob request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("paymob request failed (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			return fmt.Errorf("failed to decode paymob response: %w", err)
		}
	}

	return nil
}

func getMapValueOrDefault(data map[string]string, key string, fallback string) string {
	if data == nil {
		return fallback
	}
	value := strings.TrimSpace(data[key])
	if value == "" {
		return fallback
	}
	return value
}

func buildPaymobHMACMessage(webhook *domain.PaymobWebhookPayload) string {
	if webhook == nil {
		return ""
	}

	obj := webhook.Obj
	orderedValues := map[string]string{
		"amount_cents":           strconv.Itoa(obj.AmountCents),
		"created_at":             obj.CreatedAt,
		"currency":               obj.Currency,
		"error_occured":          strconv.FormatBool(obj.ErrorOccured),
		"has_parent_transaction": strconv.FormatBool(obj.HasParentTransaction),
		"id":                     strconv.FormatInt(obj.ID, 10),
		"integration_id":         strconv.FormatInt(obj.IntegrationID, 10),
		"is_3d_secure":           strconv.FormatBool(obj.Is3DSecure),
		"is_auth":                strconv.FormatBool(obj.IsAuth),
		"is_capture":             strconv.FormatBool(obj.IsCapture),
		"is_refunded":            strconv.FormatBool(obj.IsRefunded),
		"is_standalone_payment":  strconv.FormatBool(obj.IsStandalonePayment),
		"is_voided":              strconv.FormatBool(obj.IsVoided),
		"order.id":               strconv.FormatInt(obj.Order.ID, 10),
		"owner":                  strconv.FormatInt(obj.Owner, 10),
		"pending":                strconv.FormatBool(obj.Pending),
		"source_data.pan":        obj.SourceData.Pan,
		"source_data.sub_type":   obj.SourceData.SubType,
		"source_data.type":       obj.SourceData.Type,
		"success":                strconv.FormatBool(obj.Success),
	}

	parts := make([]string, 0, len(paymobHMACFieldOrder))
	for _, key := range paymobHMACFieldOrder {
		parts = append(parts, orderedValues[key])
	}
	return strings.Join(parts, "")
}
