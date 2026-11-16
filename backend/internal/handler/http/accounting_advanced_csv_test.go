package http

import (
	"encoding/csv"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteCSVAttachmentAdvertisesCSVContract(t *testing.T) {
	response := httptest.NewRecorder()
	writeCSVAttachment(response, "chart-of-accounts-export.csv", [][]string{
		{"account_code", "account_name"},
		{"1000", "Cash"},
	})

	if response.Code != 200 {
		t.Fatalf("status = %d, want 200", response.Code)
	}
	if got := response.Header().Get("Content-Type"); got != "text/csv; charset=utf-8" {
		t.Fatalf("Content-Type = %q, want CSV media type", got)
	}
	if got := response.Header().Get("Content-Disposition"); got != `attachment; filename="chart-of-accounts-export.csv"` {
		t.Fatalf("Content-Disposition = %q, want CSV filename", got)
	}
	rows, err := csv.NewReader(strings.NewReader(response.Body.String())).ReadAll()
	if err != nil {
		t.Fatalf("response body is not valid CSV: %v", err)
	}
	if len(rows) != 2 || rows[1][0] != "1000" || rows[1][1] != "Cash" {
		t.Fatalf("unexpected CSV rows: %#v", rows)
	}
}
