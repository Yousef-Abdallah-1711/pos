package email

import (
	"strings"
	"testing"
)

func TestRenderVerificationEmailInEnglishAndArabic(t *testing.T) {
	data := TemplateData{Code: "482901", ExpiresMinutes: 10}
	for _, locale := range []Locale{LocaleEnglish, LocaleArabic} {
		message, err := Render(TemplateVerification, locale, data, "https://pos.app.rakmyat.online")
		if err != nil {
			t.Fatalf("Render(%q): %v", locale, err)
		}
		if !strings.Contains(message.HTMLBody, `src="`+LogoURL+`"`) {
			t.Fatalf("email does not include the official logo URL")
		}
		if !strings.Contains(message.HTMLBody, `dir="ltr"`) && locale == LocaleEnglish {
			t.Fatalf("English email is not left-to-right")
		}
		if !strings.Contains(message.HTMLBody, `dir="rtl"`) && locale == LocaleArabic {
			t.Fatalf("Arabic email is not right-to-left")
		}
		if !strings.Contains(message.TextBody, "482901") || !strings.Contains(message.TextBody, "10") {
			t.Fatalf("plain-text alternative omits code or expiry")
		}
	}
}

func TestRenderInvitationEscapesDynamicContentAndUsesSameSiteURL(t *testing.T) {
	message, err := Render(TemplateInvitation, LocaleEnglish, TemplateData{
		TenantName: `<img src=x onerror=alert(1)>`, RoleName: `Cashier & Manager`, ActionPath: "/en/auth/invitations/accept?token=synthetic",
	}, "https://pos.app.rakmyat.online")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(message.HTMLBody, `<img src=x onerror`) || !strings.Contains(message.HTMLBody, "&lt;img") {
		t.Fatalf("dynamic tenant content was not HTML-escaped")
	}
	if !strings.Contains(message.HTMLBody, "Cashier &amp; Manager") {
		t.Fatalf("dynamic role content was not HTML-escaped")
	}
	if !strings.Contains(message.HTMLBody, `https://pos.app.rakmyat.online/en/auth/invitations/accept?token=synthetic`) {
		t.Fatalf("invitation link was not built on the configured application origin")
	}
	if !strings.Contains(message.TextBody, `https://pos.app.rakmyat.online/en/auth/invitations/accept?token=synthetic`) {
		t.Fatalf("plain-text alternative does not include the full invitation link")
	}
}

func TestRenderRejectsInvalidOTPAndUntrustedURLs(t *testing.T) {
	validURL := "https://pos.app.rakmyat.online"
	if _, err := Render(TemplateVerification, LocaleEnglish, TemplateData{Code: "12345x", ExpiresMinutes: 10}, validURL); err == nil {
		t.Fatal("expected malformed verification code to be rejected")
	}
	if _, err := Render(TemplateVerification, LocaleEnglish, TemplateData{Code: "123456", ExpiresMinutes: 9}, validURL); err == nil {
		t.Fatal("expected non-10-minute verification expiry to be rejected")
	}
	if _, err := Render(TemplateInvitation, LocaleEnglish, TemplateData{TenantName: "Cafe", RoleName: "Owner", ActionPath: "//evil.example"}, validURL); err == nil {
		t.Fatal("expected protocol-relative action URL to be rejected")
	}
	if _, err := Render(TemplateInvitation, LocaleEnglish, TemplateData{TenantName: "Cafe", RoleName: "Owner", ActionPath: "/accept"}, "http://pos.app.rakmyat.online"); err == nil {
		t.Fatal("expected non-TLS application URL to be rejected")
	}
}

func TestRenderPasswordResetAndSecurityNotice(t *testing.T) {
	reset, err := Render(TemplatePasswordReset, LocaleArabic, TemplateData{
		ActionPath: "/ar/auth/reset?token=synthetic", ExpiresMinutes: 20,
	}, "https://pos.app.rakmyat.online")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(reset.HTMLBody, "20") || !strings.Contains(reset.HTMLBody, "إعادة تعيين كلمة المرور") {
		t.Fatal("Arabic reset email is missing its expiry or localized action")
	}
	notice, err := Render(TemplateSecurityNotice, LocaleEnglish, TemplateData{Detail: "Password changed"}, "https://pos.app.rakmyat.online")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(notice.TextBody, "Password changed") {
		t.Fatal("security notice omitted its detail")
	}
}
