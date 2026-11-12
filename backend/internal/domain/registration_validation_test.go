package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNormalizeAndValidateRegistration(t *testing.T) {
	tests := []struct {
		name       string
		request    RegisterRequest
		wantPhone  string
		wantFields []string
	}{
		{
			name:      "trims and normalizes international input",
			request:   RegisterRequest{Name: "  يوسف عبد الله ", Email: " USER@Example.COM ", Password: "123456", RestaurantName: " مطعم النيل ", RestaurantSlug: "nile-food", Phone: "+1 (650) 555-0123"},
			wantPhone: "+16505550123",
		},
		{
			name:    "optional and whitespace-only phones become empty",
			request: RegisterRequest{Name: "Jane Doe", Email: "jane@example.com", Password: "123456", RestaurantName: "Cafe North", Phone: " \t\n "},
		},
		{
			name:       "overlong unicode names are rejected by code point count",
			request:    RegisterRequest{Name: strings.Repeat("名", 256), Email: "jane@example.com", Password: "123456", RestaurantName: "Cafe North"},
			wantFields: []string{"name"},
		},
		{
			name:       "short unicode names are rejected",
			request:    RegisterRequest{Name: "名", Email: "jane@example.com", Password: "123456", RestaurantName: "Cafe North"},
			wantFields: []string{"name"},
		},
		{
			name:       "whitespace-only names are rejected",
			request:    RegisterRequest{Name: "   ", Email: "jane@example.com", Password: "123456", RestaurantName: "Cafe North"},
			wantFields: []string{"name"},
		},
		{
			name:       "invalid email is rejected",
			request:    RegisterRequest{Name: "Jane Doe", Email: "not-an-email", Password: "123456", RestaurantName: "Cafe North"},
			wantFields: []string{"email"},
		},
		{
			name:       "password exceeding bcrypt byte limit is rejected",
			request:    RegisterRequest{Name: "Jane Doe", Email: "jane@example.com", Password: strings.Repeat("é", 37), RestaurantName: "Cafe North"},
			wantFields: []string{"password"},
		},
		{
			name:      "valid international phone is accepted even when presentation exceeds varchar length",
			request:   RegisterRequest{Name: "Jane Doe", Email: "jane@example.com", Password: "123456", RestaurantName: "Cafe North", Phone: "+1       (650) 555-0123"},
			wantPhone: "+16505550123",
		},
		{
			name:       "local phone without country context is rejected",
			request:    RegisterRequest{Name: "Jane Doe", Email: "jane@example.com", Password: "123456", RestaurantName: "Cafe North", Phone: "6505550123"},
			wantFields: []string{"phone"},
		},
		{
			name:       "invalid phone is rejected",
			request:    RegisterRequest{Name: "Jane Doe", Email: "jane@example.com", Password: "123456", RestaurantName: "Cafe North", Phone: "+999123"},
			wantFields: []string{"phone"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := test.request
			err := request.NormalizeAndValidate()
			if len(test.wantFields) == 0 {
				if err != nil {
					t.Fatalf("NormalizeAndValidate() error = %v", err)
				}
				if request.Phone != test.wantPhone {
					t.Errorf("normalized phone = %q, want %q", request.Phone, test.wantPhone)
				}
				if request.Name != strings.TrimSpace(test.request.Name) {
					t.Errorf("normalized name = %q", request.Name)
				}
				if request.Email != strings.ToLower(strings.TrimSpace(test.request.Email)) {
					t.Errorf("normalized email = %q", request.Email)
				}
				return
			}
			var validationErrors *ValidationErrors
			if !errors.As(err, &validationErrors) {
				t.Fatalf("error = %v, want ValidationErrors", err)
			}
			for _, field := range test.wantFields {
				if validationErrors.Fields[field] == "" {
					t.Errorf("missing validation error for %q: %#v", field, validationErrors.Fields)
				}
			}
		})
	}
}

func TestNormalizeInternationalPhoneBoundaries(t *testing.T) {
	valid := map[string]string{
		"+16505550123":            "+16505550123",
		"+1 (650) 555-0123":       "+16505550123",
		"+44 20 7031 3000":        "+442070313000",
		"+20 2 1234 5678":         "+20212345678",
		"+1       (650) 555-0123": "+16505550123",
	}
	for input, expected := range valid {
		got, err := NormalizeInternationalPhone(input)
		if err != nil || got != expected {
			t.Errorf("NormalizeInternationalPhone(%q) = %q, %v; want %q", input, got, err, expected)
		}
	}

	for _, input := range []string{"6505550123", "+999123", "+1 abc 6505550123", "+١٦٥٠٥٥٥٠١٢٣"} {
		got, err := NormalizeInternationalPhone(input)
		if err == nil {
			t.Errorf("NormalizeInternationalPhone(%q) = %q, want error", input, got)
		}
	}
	for _, input := range []string{"", "  "} {
		got, err := NormalizeInternationalPhone(input)
		if err != nil || got != "" {
			t.Errorf("NormalizeInternationalPhone(%q) = %q, %v; want optional empty phone", input, got, err)
		}
	}
}

func TestRegistrationPhoneRawStringLengthsDoNotBypassCanonicalContract(t *testing.T) {
	for _, input := range []string{
		"+165055501234567890",   // 19 characters, invalid numbering plan length
		"+1650555012345678901",  // 20 characters
		"+16505550123456789012", // 21 characters
	} {
		if _, err := NormalizeInternationalPhone(input); err == nil {
			t.Errorf("phone %q unexpectedly passed international number validation", input)
		}
	}
}

func TestLoginRequestNormalizeAndValidate(t *testing.T) {
	tests := []struct {
		name  string
		input LoginRequest
		field string
	}{
		{name: "trim and normalize email", input: LoginRequest{Email: " USER@Example.COM ", Password: "secret"}},
		{name: "required email", input: LoginRequest{Password: "secret"}, field: "email"},
		{name: "email format", input: LoginRequest{Email: "invalid", Password: "secret"}, field: "email"},
		{name: "required password", input: LoginRequest{Email: "user@example.com"}, field: "password"},
		{name: "bcrypt byte limit", input: LoginRequest{Email: "user@example.com", Password: strings.Repeat("é", 37)}, field: "password"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.input.NormalizeAndValidate()
			if test.field == "" {
				if err != nil || test.input.Email != "user@example.com" {
					t.Fatalf("NormalizeAndValidate() = %v, email %q", err, test.input.Email)
				}
				return
			}
			var validation *ValidationErrors
			if !errors.As(err, &validation) || validation.Fields[test.field] == "" {
				t.Fatalf("NormalizeAndValidate() = %v, want %q field error", err, test.field)
			}
		})
	}
}
