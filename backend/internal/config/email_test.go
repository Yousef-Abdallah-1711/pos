package config

import "testing"

func TestLoadEmailConfigAllowsSMTPToBeDisabled(t *testing.T) {
	for _, key := range []string{"SMTP_HOST", "SMTP_PORT", "SMTP_USERNAME", "SMTP_PASSWORD", "SMTP_FROM_EMAIL", "SMTP_FROM_NAME"} {
		t.Setenv(key, "")
	}
	config, err := loadEmailConfig()
	if err != nil {
		t.Fatal(err)
	}
	if config.SMTPConfigured() || config.SMTPPort != 587 {
		t.Fatalf("unexpected disabled SMTP config: %+v", config)
	}
}

func TestLoadEmailConfigRequiresCompleteCredentials(t *testing.T) {
	t.Setenv("SMTP_HOST", "smtp.example.test")
	t.Setenv("SMTP_PORT", "587")
	t.Setenv("SMTP_USERNAME", "mailer")
	t.Setenv("SMTP_PASSWORD", "synthetic-secret")
	t.Setenv("SMTP_FROM_EMAIL", "mailer@example.test")
	t.Setenv("SMTP_FROM_NAME", "SERVIO")
	config, err := loadEmailConfig()
	if err != nil {
		t.Fatal(err)
	}
	if !config.SMTPConfigured() || config.SMTPHost != "smtp.example.test" || config.SMTPPort != 587 {
		t.Fatalf("SMTP config was not loaded: %+v", config)
	}
}

func TestLoadEmailConfigRejectsPartialCredentialsAndInvalidPort(t *testing.T) {
	for _, key := range []string{"SMTP_HOST", "SMTP_PORT", "SMTP_USERNAME", "SMTP_PASSWORD", "SMTP_FROM_EMAIL", "SMTP_FROM_NAME"} {
		t.Setenv(key, "")
	}
	t.Setenv("SMTP_HOST", "smtp.example.test")
	if _, err := loadEmailConfig(); err == nil {
		t.Fatal("expected partial SMTP config to be rejected")
	}
	t.Setenv("SMTP_USERNAME", "mailer")
	t.Setenv("SMTP_PASSWORD", "synthetic-secret")
	t.Setenv("SMTP_FROM_EMAIL", "mailer@example.test")
	t.Setenv("SMTP_PORT", "65536")
	if _, err := loadEmailConfig(); err == nil {
		t.Fatal("expected invalid SMTP port to be rejected")
	}
}
