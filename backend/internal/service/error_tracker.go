package service

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
)

// ErrorTracker captures runtime errors/messages for observability.
type ErrorTracker interface {
	CaptureError(err error, context map[string]interface{})
	CaptureMessage(msg string, level string)
}

// NoopTracker is used when no tracker is configured.
type NoopTracker struct{}

// NewNoopTracker creates a no-op tracker.
func NewNoopTracker() ErrorTracker {
	return &NoopTracker{}
}

func (n *NoopTracker) CaptureError(err error, context map[string]interface{}) {}

func (n *NoopTracker) CaptureMessage(msg string, level string) {}

// SentryTracker implements ErrorTracker using sentry-go SDK.
type SentryTracker struct{}

// NewSentryTracker initializes sentry-go with DSN and returns a tracker.
func NewSentryTracker(dsn string) (ErrorTracker, error) {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              strings.TrimSpace(dsn),
		AttachStacktrace: true,
	}); err != nil {
		return nil, err
	}

	return &SentryTracker{}, nil
}

// NewErrorTrackerFromEnv creates Sentry tracker when SENTRY_DSN exists, else no-op.
func NewErrorTrackerFromEnv() ErrorTracker {
	dsn := strings.TrimSpace(os.Getenv("SENTRY_DSN"))
	if dsn == "" {
		return NewNoopTracker()
	}

	tracker, err := NewSentryTracker(dsn)
	if err != nil {
		log.Printf("[ErrorTracker] WARN: failed to initialize sentry tracker: %v", err)
		return NewNoopTracker()
	}

	log.Printf("[ErrorTracker] Sentry tracker initialized")
	return tracker
}

func (s *SentryTracker) CaptureError(err error, context map[string]interface{}) {
	if err == nil {
		return
	}

	sentry.WithScope(func(scope *sentry.Scope) {
		for key, value := range context {
			scope.SetExtra(key, value)
		}
		sentry.CaptureException(err)
	})
	sentry.Flush(2 * time.Second)
}

func (s *SentryTracker) CaptureMessage(msg string, level string) {
	trimmed := strings.TrimSpace(msg)
	if trimmed == "" {
		return
	}

	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(parseSentryLevel(level))
		sentry.CaptureMessage(trimmed)
	})
	sentry.Flush(2 * time.Second)
}

func parseSentryLevel(level string) sentry.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return sentry.LevelDebug
	case "warning", "warn":
		return sentry.LevelWarning
	case "error":
		return sentry.LevelError
	case "fatal":
		return sentry.LevelFatal
	default:
		return sentry.LevelInfo
	}
}
