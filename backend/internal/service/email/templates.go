package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"regexp"
	"strings"
)

type Locale string

const (
	LocaleEnglish Locale = "en"
	LocaleArabic  Locale = "ar"
	LogoURL              = "https://pos.rakmyat.online/brand/png/poos-wordmark-480.png"
)

type TemplateKind string

const (
	TemplateVerification   TemplateKind = "verification"
	TemplateInvitation     TemplateKind = "invitation"
	TemplatePasswordReset  TemplateKind = "password_reset"
	TemplateSecurityNotice TemplateKind = "security_notice"
)

type TemplateData struct {
	Code           string
	TenantName     string
	RoleName       string
	ActionPath     string
	Detail         string
	ExpiresMinutes int
}

type renderedTemplate struct {
	Direction   string
	Locale      string
	LogoURL     string
	Title       string
	Greeting    string
	Body        string
	Action      string
	ActionURL   string
	Code        string
	Tenant      string
	Role        string
	Detail      string
	Expiry      string
	Footer      string
	TenantLabel string
	RoleLabel   string
}

type localizedCopy struct {
	subject   string
	title     string
	greeting  string
	body      string
	action    string
	footer    string
	codeLabel string
	tenant    string
	role      string
	detail    string
}

var otpPattern = regexp.MustCompile(`^[0-9]{6}$`)
var emailTemplates = template.Must(template.New("email").Parse(emailTemplateHTML))

const emailTemplateHTML = `<!doctype html>
<html lang="{{.Locale}}" dir="{{.Direction}}">
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<style>
@media (prefers-color-scheme: dark) {
  body, .page { background-color:#111827 !important; }
  .card { background-color:#1f2937 !important; color:#f9fafb !important; }
  .muted { color:#d1d5db !important; }
}
</style></head>
<body style="margin:0;background:#f3f4f6;font-family:Arial,Helvetica,sans-serif;color:#111827;">
<table role="presentation" width="100%" cellspacing="0" cellpadding="0" class="page" style="background:#f3f4f6;padding:24px 12px;">
<tr><td align="center">
<table role="presentation" width="100%" cellspacing="0" cellpadding="0" class="card" style="max-width:560px;background:#ffffff;border-radius:12px;padding:28px;">
<tr><td align="center" style="padding:4px 0 24px;"><img src="{{.LogoURL}}" width="180" alt="SERVIO" style="display:block;width:180px;max-width:100%;height:auto;border:0;"></td></tr>
<tr><td><h1 style="font-size:22px;line-height:1.35;margin:0 0 16px;text-align:center;">{{.Title}}</h1></td></tr>
<tr><td style="font-size:16px;line-height:1.65;"><p style="margin:0 0 12px;">{{.Greeting}}</p><p class="muted" style="margin:0 0 18px;color:#4b5563;">{{.Body}}</p></td></tr>
{{if .Tenant}}<tr><td style="font-size:15px;line-height:1.6;padding:4px 0;"><strong>{{.TenantLabel}}:</strong> {{.Tenant}}</td></tr>{{end}}
{{if .Role}}<tr><td style="font-size:15px;line-height:1.6;padding:4px 0 16px;"><strong>{{.RoleLabel}}:</strong> {{.Role}}</td></tr>{{end}}
{{if .Code}}<tr><td align="center" style="padding:8px 0 20px;"><span style="display:inline-block;border-radius:8px;background:#eef2ff;color:#312e81;font-size:28px;font-weight:700;letter-spacing:8px;padding:12px 20px;">{{.Code}}</span><p class="muted" style="font-size:13px;color:#6b7280;margin:10px 0 0;">{{.Expiry}}</p></td></tr>{{end}}
{{if .ActionURL}}<tr><td align="center" style="padding:8px 0 20px;"><a href="{{.ActionURL}}" style="display:inline-block;background:#4338ca;color:#ffffff;text-decoration:none;border-radius:8px;padding:13px 22px;font-size:15px;font-weight:700;">{{.Action}}</a></td></tr>{{end}}
{{if .Detail}}<tr><td class="muted" style="font-size:14px;line-height:1.6;color:#4b5563;padding:0 0 16px;">{{.Detail}}</td></tr>{{end}}
<tr><td class="muted" style="font-size:12px;line-height:1.6;color:#6b7280;border-top:1px solid #e5e7eb;padding-top:18px;">{{.Footer}}</td></tr>
</table></td></tr></table></body></html>`

var localized = map[Locale]map[TemplateKind]localizedCopy{
	LocaleEnglish: {
		TemplateVerification:   {"Verify your SERVIO email", "Verify your email", "Hello,", "Use this one-time code to verify your email address.", "", "If you did not request this, you can ignore this message.", "Verification code", "Tenant", "Role", ""},
		TemplateInvitation:     {"Your SERVIO invitation", "You are invited to SERVIO", "Hello,", "You have been invited to join this restaurant workspace.", "Accept invitation", "If you were not expecting this invitation, you can ignore this message.", "", "Tenant", "Role", ""},
		TemplatePasswordReset:  {"Reset your SERVIO password", "Reset your password", "Hello,", "Use the secure button below to choose a new password.", "Reset password", "If you did not request a password reset, ignore this message and your password will remain unchanged.", "", "", "", ""},
		TemplateSecurityNotice: {"SERVIO security notice", "Account security notice", "Hello,", "A security-related change was made to your account.", "", "If you do not recognize this activity, contact your administrator.", "", "", "", "Details"},
	},
	LocaleArabic: {
		TemplateVerification:   {"تحقق من بريدك الإلكتروني في SERVIO", "تحقق من بريدك الإلكتروني", "مرحباً،", "استخدم الرمز لمرة واحدة للتحقق من عنوان بريدك الإلكتروني.", "", "إذا لم تطلب ذلك، يمكنك تجاهل هذه الرسالة.", "رمز التحقق", "المطعم", "الدور", ""},
		TemplateInvitation:     {"دعوة للانضمام إلى SERVIO", "تمت دعوتك إلى SERVIO", "مرحباً،", "تمت دعوتك للانضمام إلى مساحة عمل المطعم.", "قبول الدعوة", "إذا لم تكن تتوقع هذه الدعوة، يمكنك تجاهل الرسالة.", "", "المطعم", "الدور", ""},
		TemplatePasswordReset:  {"إعادة تعيين كلمة مرور SERVIO", "إعادة تعيين كلمة المرور", "مرحباً،", "استخدم الزر الآمن أدناه لاختيار كلمة مرور جديدة.", "إعادة تعيين كلمة المرور", "إذا لم تطلب إعادة التعيين، تجاهل الرسالة وستبقى كلمة المرور كما هي.", "", "", "", ""},
		TemplateSecurityNotice: {"تنبيه أمان SERVIO", "تنبيه أمان الحساب", "مرحباً،", "تم إجراء تغيير أمني على حسابك.", "", "إذا لم تتعرف على هذا النشاط، تواصل مع مسؤول النظام.", "", "", "", "التفاصيل"},
	},
}

func Render(kind TemplateKind, locale Locale, data TemplateData, appBaseURL string) (Message, error) {
	copySet, ok := localized[locale][kind]
	if !ok {
		return Message{}, fmt.Errorf("unsupported email template or locale")
	}
	if err := validateTemplateData(kind, data); err != nil {
		return Message{}, err
	}
	baseURL, err := validateBaseURL(appBaseURL)
	if err != nil {
		return Message{}, err
	}
	content := renderedTemplate{
		Direction: "ltr", LogoURL: LogoURL, Title: copySet.title, Greeting: copySet.greeting,
		Body: copySet.body, Action: copySet.action, Code: data.Code, Tenant: data.TenantName,
		Role: data.RoleName, Detail: data.Detail, Expiry: expiryText(locale, data.ExpiresMinutes),
		Footer: copySet.footer,
	}
	if locale == LocaleArabic {
		content.Direction = "rtl"
	}
	content.Locale = string(locale)
	if data.ActionPath != "" {
		content.ActionURL = strings.TrimRight(baseURL, "/") + data.ActionPath
	}
	content.TenantLabel, content.RoleLabel = copySet.tenant, copySet.role
	var htmlBody bytes.Buffer
	if err := emailTemplates.Execute(&htmlBody, content); err != nil {
		return Message{}, fmt.Errorf("render email template: %w", err)
	}
	return Message{
		Subject: copySet.subject, HTMLBody: htmlBody.String(), TextBody: plainText(locale, copySet, data, content.ActionURL),
	}, nil
}

func validateTemplateData(kind TemplateKind, data TemplateData) error {
	switch kind {
	case TemplateVerification:
		if !otpPattern.MatchString(data.Code) || data.ExpiresMinutes != 10 {
			return fmt.Errorf("email verification requires a six-digit code with a 10-minute expiry")
		}
	case TemplateInvitation:
		if strings.TrimSpace(data.TenantName) == "" || strings.TrimSpace(data.RoleName) == "" || data.ActionPath == "" {
			return fmt.Errorf("invitation requires a tenant, role, and acceptance path")
		}
	case TemplatePasswordReset:
		if data.ActionPath == "" || data.ExpiresMinutes < 1 || data.ExpiresMinutes > 60 {
			return fmt.Errorf("password reset requires a bounded expiry and reset path")
		}
	case TemplateSecurityNotice:
		if strings.TrimSpace(data.Detail) == "" {
			return fmt.Errorf("security notice details are required")
		}
	default:
		return fmt.Errorf("unsupported email template")
	}
	if strings.ContainsAny(data.TenantName+data.RoleName+data.Detail, "\r\n") {
		return fmt.Errorf("email template content contains an invalid line break")
	}
	if data.ActionPath != "" && (strings.HasPrefix(data.ActionPath, "//") || !strings.HasPrefix(data.ActionPath, "/") || strings.ContainsAny(data.ActionPath, "\\\r\n")) {
		return fmt.Errorf("email action must be a same-site absolute path")
	}
	return nil
}

func validateBaseURL(value string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(value))
	if err != nil || parsed.Scheme != "https" || parsed.Hostname() == "" || parsed.User != nil || parsed.Fragment != "" || parsed.RawQuery != "" {
		return "", fmt.Errorf("email application base URL must be an HTTPS origin")
	}
	if parsed.Path != "" && parsed.Path != "/" {
		return "", fmt.Errorf("email application base URL must not include a path")
	}
	return parsed.Scheme + "://" + parsed.Host, nil
}

func expiryText(locale Locale, minutes int) string {
	if locale == LocaleArabic {
		return fmt.Sprintf("تنتهي صلاحية الرمز خلال %d دقائق.", minutes)
	}
	return fmt.Sprintf("This code expires in %d minutes.", minutes)
}

func plainText(locale Locale, copySet localizedCopy, data TemplateData, actionURL string) string {
	lines := []string{copySet.greeting, copySet.body}
	if data.TenantName != "" {
		lines = append(lines, copySet.tenant+": "+data.TenantName)
	}
	if data.RoleName != "" {
		lines = append(lines, copySet.role+": "+data.RoleName)
	}
	if data.Code != "" {
		lines = append(lines, copySet.codeLabel+": "+data.Code, expiryText(locale, data.ExpiresMinutes))
	}
	if actionURL != "" {
		lines = append(lines, copySet.action+": "+actionURL)
	}
	if data.Detail != "" {
		lines = append(lines, copySet.detail+": "+data.Detail)
	}
	lines = append(lines, copySet.footer)
	return strings.Join(lines, "\n\n")
}
