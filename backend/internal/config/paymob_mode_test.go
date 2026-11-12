package config

import "testing"

func setPaymobEnv(t *testing.T, env, mode string) {
	t.Helper()
	clearRuntimeEnv(t)
	t.Setenv("ENV", env)
	t.Setenv("PAYMOB_MODE", mode)
	t.Setenv("PAYMOB_API_KEY", "")
	t.Setenv("PAYMOB_HMAC_SECRET", "")
	t.Setenv("PAYMOB_IFRAME_ID", "")
}

// SEC-G: production must never silently fall back to the mock payment provider.
func TestResolvePaymobMode_ProductionEmptyFailsClosed(t *testing.T) {
	setPaymobEnv(t, "production", "")
	if mode, err := ResolvePaymobMode(); err == nil {
		t.Fatalf("empty PAYMOB_MODE in production must error, got mode=%q", mode)
	}
}

func TestResolvePaymobMode_ProductionMockIsRejected(t *testing.T) {
	setPaymobEnv(t, "production", "mock")
	if _, err := ResolvePaymobMode(); err == nil {
		t.Fatal("mock mode in production must error")
	}
}

func TestResolvePaymobMode_ProductionLiveRequiresCredentials(t *testing.T) {
	setPaymobEnv(t, "production", "live")
	if _, err := ResolvePaymobMode(); err == nil {
		t.Fatal("live mode without credentials must error")
	}

	t.Setenv("PAYMOB_API_KEY", "key")
	t.Setenv("PAYMOB_HMAC_SECRET", "secret")
	if _, err := ResolvePaymobMode(); err == nil {
		t.Fatal("live mode without iframe id must error")
	}

	t.Setenv("PAYMOB_IFRAME_ID", "not-a-number")
	if _, err := ResolvePaymobMode(); err == nil {
		t.Fatal("non-numeric iframe id must error")
	}

	t.Setenv("PAYMOB_IFRAME_ID", "12345")
	mode, err := ResolvePaymobMode()
	if err != nil || mode != PaymobModeLive {
		t.Fatalf("fully configured live mode must pass, got %q err=%v", mode, err)
	}
}

func TestResolvePaymobMode_UnknownValueIsRejectedEverywhere(t *testing.T) {
	for _, env := range []string{"development", "production", "test"} {
		setPaymobEnv(t, env, "sandboxish")
		if _, err := ResolvePaymobMode(); err == nil {
			t.Fatalf("unknown mode must error in %s", env)
		}
	}
}

func TestResolvePaymobMode_MockStillWorksOutsideProduction(t *testing.T) {
	for _, env := range []string{"development", "dev", "local", "test", ""} {
		for _, mode := range []string{"", "mock", " MOCK "} {
			setPaymobEnv(t, env, mode)
			got, err := ResolvePaymobMode()
			if err != nil || got != PaymobModeMock {
				t.Fatalf("env=%q mode=%q: want mock, got %q err=%v", env, mode, got, err)
			}
		}
	}
}

func TestResolvePaymobMode_LiveOutsideProductionNeedsCredentialsToo(t *testing.T) {
	setPaymobEnv(t, "development", "live")
	if _, err := ResolvePaymobMode(); err == nil {
		t.Fatal("live mode must always require credentials")
	}
}
