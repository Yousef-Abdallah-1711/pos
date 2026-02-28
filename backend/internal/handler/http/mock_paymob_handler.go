package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	appconfig "pos-saas/internal/config"
	"strconv"
	"strings"
	"time"

	"pos-saas/internal/domain"
)

// MockPaymobHandler serves a fake hosted payment page for development mode.
type MockPaymobHandler struct {
	webhookClient *http.Client
}

// NewMockPaymobHandler creates a mock Paymob handler.
func NewMockPaymobHandler() *MockPaymobHandler {
	return &MockPaymobHandler{
		webhookClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type mockIframePageData struct {
	PaymentToken string
	AmountLabel  string
	Currency     string
	RedirectURL  string
}

type mockPaymentTokenData struct {
	OrderID     int64
	AmountCents int
	Currency    string
}

// ServeIframe handles GET /api/v1/mock/paymob/iframe?payment_token=...
func (h *MockPaymobHandler) ServeIframe(w http.ResponseWriter, r *http.Request) {
	paymentToken := strings.TrimSpace(r.URL.Query().Get("payment_token"))
	if paymentToken == "" {
		BadRequestResponse(w, "payment_token is required")
		return
	}

	redirectURL := strings.TrimSpace(r.URL.Query().Get("redirect_url"))
	if redirectURL == "" {
		redirectURL = defaultBillingRedirectURL(r)
	}
	redirectURL = normalizeBillingRedirectURL(redirectURL)
	if redirectURL == "" {
		redirectURL = defaultBillingRedirectURL(r)
	}

	tokenData, err := parseMockPaymentToken(paymentToken)
	if err != nil {
		log.Printf("[MockPaymob] WARN: token parsing failed: %v\n", err)
		// Keep mock page usable even with legacy token format.
		tokenData = mockPaymentTokenData{
			OrderID:     time.Now().UnixNano(),
			AmountCents: 50000,
			Currency:    "EGP",
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := mockIframeTemplate.Execute(w, mockIframePageData{
		PaymentToken: paymentToken,
		AmountLabel:  formatAmountLabel(tokenData.AmountCents),
		Currency:     tokenData.Currency,
		RedirectURL:  redirectURL,
	}); err != nil {
		log.Printf("[MockPaymob] ERROR: failed to render iframe: %v\n", err)
	}
}

// ProcessPayment handles POST /api/v1/mock/paymob/process and triggers webhook callback.
func (h *MockPaymobHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		BadRequestResponse(w, "invalid form payload")
		return
	}

	paymentToken := strings.TrimSpace(r.FormValue("payment_token"))
	if paymentToken == "" {
		BadRequestResponse(w, "payment_token is required")
		return
	}

	tokenData, err := parseMockPaymentToken(paymentToken)
	if err != nil {
		BadRequestResponse(w, "invalid payment_token format")
		return
	}

	simulate := strings.ToLower(strings.TrimSpace(r.FormValue("simulate")))
	isFailure := simulate == "fail" || simulate == "failed" || r.URL.Query().Get("fail") == "true"
	isSuccess := !isFailure

	payload := buildMockWebhookPayload(tokenData, isSuccess, r.FormValue("card_number"))
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		InternalErrorResponse(w, "failed to prepare mock webhook payload")
		return
	}

	webhookURL := buildWebhookURL(r)
	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, webhookURL, bytes.NewReader(payloadBytes))
	if err != nil {
		InternalErrorResponse(w, "failed to prepare webhook request")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.webhookClient.Do(req)
	if err != nil {
		log.Printf("[MockPaymob] ERROR: webhook callback failed: %v\n", err)
	} else {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		log.Printf("[MockPaymob] webhook callback completed - status=%d order_id=%d\n", resp.StatusCode, tokenData.OrderID)
	}

	redirectURL := strings.TrimSpace(r.FormValue("redirect_url"))
	if redirectURL == "" {
		redirectURL = defaultBillingRedirectURL(r)
	}
	redirectURL = normalizeBillingRedirectURL(redirectURL)
	if redirectURL == "" {
		redirectURL = defaultBillingRedirectURL(r)
	}
	target := withBillingStatus(redirectURL, map[bool]string{true: "success", false: "failed"}[isSuccess], tokenData.OrderID)

	// When submitted inside an iframe, notify parent window so the dashboard modal
	// can complete the flow without requiring full-page redirects.
	if strings.EqualFold(strings.TrimSpace(r.Header.Get("Sec-Fetch-Dest")), "iframe") {
		status := "failed"
		if isSuccess {
			status = "success"
		}

		escapedTarget := template.JSEscapeString(target)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprintf(w, `<!doctype html>
<html lang="en">
<body>
  <script>
    (function () {
      try {
        if (window.parent && window.parent !== window) {
          window.parent.postMessage({ type: 'payment', status: '%s' }, '*');
        }
      } catch (e) {}
      window.location.href = '%s';
    })();
  </script>
</body>
</html>`, status, escapedTarget)
		return
	}

	http.Redirect(w, r, target, http.StatusSeeOther)
}

func parseMockPaymentToken(token string) (mockPaymentTokenData, error) {
	token = strings.TrimSpace(token)
	parts := strings.Split(token, ".")
	if len(parts) >= 5 && parts[0] == "mockpk" {
		orderID, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil || orderID <= 0 {
			return mockPaymentTokenData{}, fmt.Errorf("invalid order id")
		}
		amountCents, err := strconv.Atoi(parts[2])
		if err != nil || amountCents <= 0 {
			return mockPaymentTokenData{}, fmt.Errorf("invalid amount")
		}
		currency := strings.ToUpper(strings.TrimSpace(parts[3]))
		if currency == "" {
			currency = "EGP"
		}
		return mockPaymentTokenData{
			OrderID:     orderID,
			AmountCents: amountCents,
			Currency:    currency,
		}, nil
	}

	// Backward compatibility for old mock token style.
	if strings.HasPrefix(token, "mock_payment_key_") {
		raw := strings.TrimPrefix(token, "mock_payment_key_")
		orderID, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || orderID <= 0 {
			orderID = time.Now().UnixNano()
		}
		return mockPaymentTokenData{
			OrderID:     orderID,
			AmountCents: 50000,
			Currency:    "EGP",
		}, nil
	}

	return mockPaymentTokenData{}, fmt.Errorf("unsupported token format")
}

func buildMockWebhookPayload(tokenData mockPaymentTokenData, success bool, cardNumber string) domain.PaymobWebhookPayload {
	payload := domain.PaymobWebhookPayload{
		Type: "TRANSACTION",
		HMAC: "mock_hmac",
	}

	payload.Obj.ID = time.Now().UnixNano()
	payload.Obj.AmountCents = tokenData.AmountCents
	payload.Obj.Success = success
	payload.Obj.Currency = tokenData.Currency
	payload.Obj.IsAuth = success
	payload.Obj.IsCapture = success
	payload.Obj.IsRefunded = false
	payload.Obj.IsVoided = false
	payload.Obj.Is3DSecure = false
	payload.Obj.IsStandalonePayment = true
	payload.Obj.HasParentTransaction = false
	payload.Obj.ErrorOccured = !success
	payload.Obj.Pending = false
	payload.Obj.IntegrationID = 0
	payload.Obj.Owner = 0
	payload.Obj.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	payload.Obj.Order.ID = tokenData.OrderID
	payload.Obj.Order.MerchantOrderID = fmt.Sprintf("mock_order_%d", tokenData.OrderID)
	payload.Obj.SourceData.Pan = sanitizeCardNumber(cardNumber)
	payload.Obj.SourceData.SubType = detectCardSubtype(payload.Obj.SourceData.Pan)
	payload.Obj.SourceData.Type = "card"

	return payload
}

func buildWebhookURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")); forwarded != "" {
		scheme = strings.ToLower(strings.Split(forwarded, ",")[0])
	}
	return fmt.Sprintf("%s://%s/api/v1/webhook/paymob?hmac=mock_hmac", scheme, r.Host)
}

func withBillingStatus(baseURL string, status string, orderID int64) string {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		baseURL = defaultBillingRedirectURL(nil)
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		u, _ = url.Parse(defaultBillingRedirectURL(nil))
	}

	query := u.Query()
	query.Set("status", status)
	query.Set("source", "mock_paymob")
	query.Set("order_id", strconv.FormatInt(orderID, 10))
	u.RawQuery = query.Encode()
	return u.String()
}

func defaultBillingRedirectURL(r *http.Request) string {
	dashboardBase, err := appconfig.ResolveRuntimeURL(
		[]string{"DASHBOARD_URL", "NEXT_PUBLIC_DASHBOARD_URL"},
		"http://localhost:3002",
		"dashboard URL",
	)
	if err != nil {
		log.Printf("[MockPaymob] dashboard URL configuration error: %v", err)
		return ""
	}

	locale := detectLocale(r)
	return fmt.Sprintf("%s/%s/dashboard/settings/billing", dashboardBase, locale)
}

func normalizeBillingRedirectURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
	}

	if parsed.IsAbs() {
		return parsed.String()
	}

	// Convert app-relative URLs into absolute dashboard URLs so opening the iframe
	// directly in a new tab does not redirect back to backend host.
	if strings.HasPrefix(parsed.Path, "/") {
		base, err := appconfig.ResolveRuntimeURL(
			[]string{"DASHBOARD_URL", "NEXT_PUBLIC_DASHBOARD_URL"},
			"http://localhost:3002",
			"dashboard URL",
		)
		if err != nil {
			log.Printf("[MockPaymob] dashboard URL configuration error: %v", err)
			return ""
		}
		return base + parsed.String()
	}

	return ""
}

func detectLocale(r *http.Request) string {
	if r != nil {
		if queryLocale := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("locale"))); queryLocale == "ar" || queryLocale == "en" {
			return queryLocale
		}

		if referer := strings.TrimSpace(r.Referer()); referer != "" {
			if parsed, err := url.Parse(referer); err == nil {
				parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
				if len(parts) > 0 {
					candidate := strings.ToLower(strings.TrimSpace(parts[0]))
					if candidate == "ar" || candidate == "en" {
						return candidate
					}
				}
			}
		}

		acceptLanguage := strings.ToLower(strings.TrimSpace(r.Header.Get("Accept-Language")))
		if strings.HasPrefix(acceptLanguage, "ar") {
			return "ar"
		}
	}

	return "en"
}

func formatAmountLabel(amountCents int) string {
	if amountCents <= 0 {
		return "0.00"
	}
	return fmt.Sprintf("%.2f", float64(amountCents)/100.0)
}

func sanitizeCardNumber(cardNumber string) string {
	var digits strings.Builder
	for _, ch := range cardNumber {
		if ch >= '0' && ch <= '9' {
			digits.WriteRune(ch)
		}
	}
	out := digits.String()
	if out == "" {
		return "5123456789011111"
	}
	return out
}

func detectCardSubtype(cardNumber string) string {
	if strings.HasPrefix(cardNumber, "4") {
		return "visa"
	}
	if strings.HasPrefix(cardNumber, "5") {
		return "mastercard"
	}
	return "card"
}

var mockIframeTemplate = template.Must(template.New("mock_paymob_iframe").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Mock Paymob Checkout</title>
  <style>
    body { margin: 0; font-family: Arial, sans-serif; background: #f4f6f8; color: #1f2937; }
    .card { max-width: 420px; margin: 24px auto; background: #fff; border-radius: 10px; padding: 20px; box-shadow: 0 8px 24px rgba(16,24,40,.08); }
    h1 { margin: 0 0 8px; font-size: 20px; }
    .sub { margin: 0 0 16px; color: #6b7280; font-size: 13px; }
    .amount { margin: 12px 0 16px; padding: 12px; border-radius: 8px; background: #f8fafc; font-weight: 600; }
    label { display: block; font-size: 12px; color: #4b5563; margin: 10px 0 6px; }
    input { width: 100%; box-sizing: border-box; border: 1px solid #d1d5db; border-radius: 8px; padding: 10px 12px; font-size: 14px; }
    .row { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
    .actions { margin-top: 16px; display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
    button { border: 0; border-radius: 8px; padding: 11px 12px; font-size: 14px; cursor: pointer; }
    .ok { background: #111827; color: #fff; }
    .fail { background: #fee2e2; color: #991b1b; }
  </style>
</head>
<body>
  <div class="card">
    <h1>Mock Paymob Checkout</h1>
    <p class="sub">Development mode only. No real payment is processed.</p>
    <div class="amount">Amount: {{.AmountLabel}} {{.Currency}}</div>
    <form method="post" action="/api/v1/mock/paymob/process">
      <input type="hidden" name="payment_token" value="{{.PaymentToken}}" />
      <input type="hidden" name="redirect_url" value="{{.RedirectURL}}" />

      <label>Card Number</label>
      <input name="card_number" value="5123456789011111" autocomplete="off" />

      <div class="row">
        <div>
          <label>Expiry</label>
          <input name="expiry" value="12/30" autocomplete="off" />
        </div>
        <div>
          <label>CVV</label>
          <input name="cvv" value="123" autocomplete="off" />
        </div>
      </div>

      <div class="actions">
        <button class="ok" type="submit" name="simulate" value="success">Pay</button>
        <button class="fail" type="submit" name="simulate" value="failed">Fail</button>
      </div>
    </form>
  </div>
</body>
</html>`))
