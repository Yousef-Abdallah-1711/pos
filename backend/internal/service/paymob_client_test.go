package service

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	"pos-saas/internal/domain"
)

func TestMockPaymobClientAuthenticateReturnsToken(t *testing.T) {
	client := NewMockPaymobClient()

	token, err := client.Authenticate()
	if err != nil {
		t.Fatalf("Authenticate returned error: %v", err)
	}
	if !strings.HasPrefix(token, "mock_auth_token_") {
		t.Fatalf("expected mock auth token prefix, got: %s", token)
	}
}

func TestMockPaymobClientCreateOrderReturnsPositiveID(t *testing.T) {
	client := NewMockPaymobClient()

	orderID, err := client.CreateOrder("token", 50000, "EGP", "tenant_1_monthly_1")
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}
	if orderID <= 0 {
		t.Fatalf("expected positive order id, got: %d", orderID)
	}
}

func TestMockPaymobClientGetPaymentKeyReturnsNonEmptyToken(t *testing.T) {
	client := NewMockPaymobClient()

	key, err := client.GetPaymentKey("token", 123, 50000, "EGP", 1, map[string]string{"email": "t@example.com"}, 3600)
	if err != nil {
		t.Fatalf("GetPaymentKey returned error: %v", err)
	}
	if strings.TrimSpace(key) == "" {
		t.Fatal("expected non-empty payment key")
	}
}

func TestMockPaymobClientVerifyHMACAlwaysTrue(t *testing.T) {
	client := NewMockPaymobClient()

	if !client.VerifyHMAC([]byte(`{"obj":{"id":1}}`), "anything") {
		t.Fatal("expected mock VerifyHMAC to always return true")
	}
}

func TestRealPaymobClientVerifyHMACAcceptsValidDigest(t *testing.T) {
	const secret = "unit_test_secret"
	client := NewRealPaymobClient("api", 1, secret, "https://accept.paymob.com/api")
	realClient, ok := client.(*RealPaymobClient)
	if !ok {
		t.Fatal("expected *RealPaymobClient")
	}

	payload := fixedWebhookPayload()
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	message := buildPaymobHMACMessage(&payload)
	mac := hmac.New(sha512.New, []byte(secret))
	_, _ = mac.Write([]byte(message))
	expectedDigest := hex.EncodeToString(mac.Sum(nil))

	if !realClient.VerifyHMAC(raw, expectedDigest) {
		t.Fatal("expected VerifyHMAC to accept valid digest")
	}
}

func TestRealPaymobClientVerifyHMACRejectsTamperedDigest(t *testing.T) {
	const secret = "unit_test_secret"
	client := NewRealPaymobClient("api", 1, secret, "https://accept.paymob.com/api")
	realClient, ok := client.(*RealPaymobClient)
	if !ok {
		t.Fatal("expected *RealPaymobClient")
	}

	payload := fixedWebhookPayload()
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	message := buildPaymobHMACMessage(&payload)
	mac := hmac.New(sha512.New, []byte(secret))
	_, _ = mac.Write([]byte(message))
	validDigest := hex.EncodeToString(mac.Sum(nil))

	tampered := validDigest[:len(validDigest)-1] + "0"
	if tampered == validDigest {
		tampered = validDigest[:len(validDigest)-1] + "1"
	}

	if realClient.VerifyHMAC(raw, tampered) {
		t.Fatal("expected VerifyHMAC to reject tampered digest")
	}
}

func fixedWebhookPayload() domain.PaymobWebhookPayload {
	var payload domain.PaymobWebhookPayload
	payload.Type = "TRANSACTION"
	payload.HMAC = "unused"

	payload.Obj.ID = 987654321
	payload.Obj.AmountCents = 50000
	payload.Obj.Success = true
	payload.Obj.Currency = "EGP"
	payload.Obj.IsAuth = true
	payload.Obj.IsCapture = true
	payload.Obj.IsRefunded = false
	payload.Obj.IsVoided = false
	payload.Obj.Is3DSecure = false
	payload.Obj.IsStandalonePayment = true
	payload.Obj.HasParentTransaction = false
	payload.Obj.ErrorOccured = false
	payload.Obj.Pending = false
	payload.Obj.IntegrationID = 12345
	payload.Obj.Owner = 67890
	payload.Obj.CreatedAt = "2026-02-21T12:34:56+02:00"
	payload.Obj.Order.ID = 11223344
	payload.Obj.Order.MerchantOrderID = "tenant_1_monthly_11223344"
	payload.Obj.SourceData.Pan = "512345******1111"
	payload.Obj.SourceData.SubType = "mastercard"
	payload.Obj.SourceData.Type = "card"

	return payload
}
