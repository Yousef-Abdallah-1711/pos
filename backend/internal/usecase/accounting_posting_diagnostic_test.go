package usecase

import (
	"bytes"
	"errors"
	"log"
	"strings"
	"testing"
)

func TestLogAccountingPostingFailureIncludesScopeWithoutErrorPayload(t *testing.T) {
	var output bytes.Buffer
	previousWriter := log.Writer()
	previousFlags := log.Flags()
	previousPrefix := log.Prefix()
	log.SetOutput(&output)
	log.SetFlags(0)
	t.Cleanup(func() {
		log.SetOutput(previousWriter)
		log.SetFlags(previousFlags)
		log.SetPrefix(previousPrefix)
	})

	tests := []struct {
		operation       string
		failureCategory string
	}{
		{operation: "sales_post", failureCategory: "posting_failure"},
		{operation: "order_reload", failureCategory: "order_reload_failure"},
	}
	for _, test := range tests {
		t.Run(test.operation, func(t *testing.T) {
			output.Reset()
			logAccountingPostingFailure("paid_order", test.operation, 12, 34, 56, errors.New("provider token=must-not-appear"))
			message := output.String()
			for _, field := range []string{
				"source_event_type=\"paid_order\"", "operation=\"" + test.operation + "\"", "tenant_id=12",
				"restaurant_id=34", "source_event_id=56", "failure_category=" + test.failureCategory,
				"recovery_state=not_durable", "correlation_id=unavailable", "error_type=*errors.errorString",
			} {
				if !strings.Contains(message, field) {
					t.Errorf("diagnostic missing %q: %s", field, message)
				}
			}
			if strings.Contains(message, "must-not-appear") || strings.Contains(message, "provider token") {
				t.Fatalf("diagnostic exposed error payload: %s", message)
			}
		})
	}
}
