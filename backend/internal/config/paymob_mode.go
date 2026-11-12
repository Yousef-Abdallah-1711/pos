package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Paymob provider modes.
const (
	PaymobModeMock = "mock"
	PaymobModeLive = "live"
)

// ResolvePaymobMode decides which Paymob client the API runs with and fails
// closed on unsafe configuration:
//   - production requires PAYMOB_MODE=live explicitly (an empty value or "mock"
//     is an error, so production can never silently use the mock provider);
//   - live always requires PAYMOB_API_KEY, PAYMOB_HMAC_SECRET and a numeric
//     PAYMOB_IFRAME_ID;
//   - outside production an empty value or "mock" selects the mock provider;
//   - any other value is an error.
func ResolvePaymobMode() (string, error) {
	mode := strings.ToLower(strings.TrimSpace(os.Getenv("PAYMOB_MODE")))

	switch mode {
	case "", PaymobModeMock:
		if IsProductionMode() {
			return "", fmt.Errorf("PAYMOB_MODE must be set to %q in production (got %q): the mock payment provider is not allowed", PaymobModeLive, mode)
		}
		return PaymobModeMock, nil
	case PaymobModeLive:
		if strings.TrimSpace(os.Getenv("PAYMOB_API_KEY")) == "" {
			return "", fmt.Errorf("PAYMOB_API_KEY is required when PAYMOB_MODE=live")
		}
		if strings.TrimSpace(os.Getenv("PAYMOB_HMAC_SECRET")) == "" {
			return "", fmt.Errorf("PAYMOB_HMAC_SECRET is required when PAYMOB_MODE=live")
		}
		iframeID, err := strconv.Atoi(strings.TrimSpace(os.Getenv("PAYMOB_IFRAME_ID")))
		if err != nil || iframeID <= 0 {
			return "", fmt.Errorf("PAYMOB_IFRAME_ID must be a positive integer when PAYMOB_MODE=live")
		}
		return PaymobModeLive, nil
	default:
		return "", fmt.Errorf("unsupported PAYMOB_MODE %q (use %q or %q)", mode, PaymobModeLive, PaymobModeMock)
	}
}
