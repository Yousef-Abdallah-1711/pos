// Package testdb resolves the PostgreSQL URL used by database-backed tests and refuses to hand out
// anything that is not clearly a disposable test database.
package testdb

import (
	"net/url"
	"os"
	"strings"
	"testing"
)

// EnvVar names the environment variable holding the test database URL.
const EnvVar = "POS_TEST_DATABASE_URL"

// defaultURL is what the tests used before the URL became configurable.
const defaultURL = "postgres://test:test@localhost/pos_test?sslmode=disable"

// URL returns the test database URL. The database name must contain "test" and the host must not
// look like a production host, otherwise the test fails immediately: these tests create and drop tables.
func URL(t testing.TB) string {
	t.Helper()
	raw := strings.TrimSpace(os.Getenv(EnvVar))
	if raw == "" {
		raw = defaultURL
	}
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("%s is not a valid URL: %v", EnvVar, err)
	}
	dbName := strings.TrimPrefix(u.Path, "/")
	if !strings.Contains(strings.ToLower(dbName), "test") {
		t.Fatalf("refusing to run: database %q does not look like a test database (name must contain \"test\")", dbName)
	}
	host := strings.ToLower(u.Hostname())
	for _, bad := range []string{"prod", "production", "rds.amazonaws.com", "supabase", "neon.tech"} {
		if strings.Contains(host, bad) {
			t.Fatalf("refusing to run against host %q", host)
		}
	}
	if u.Query().Get("sslmode") == "" {
		q := u.Query()
		q.Set("sslmode", "disable")
		u.RawQuery = q.Encode()
	}
	return u.String()
}
