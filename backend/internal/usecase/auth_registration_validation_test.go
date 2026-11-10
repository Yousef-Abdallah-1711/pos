package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"

	"pos-saas/internal/domain"
)

func TestRegisterRejectsInvalidFieldsBeforeRepositoryAccess(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		slug  string
		field string
	}{
		{name: "19 character invalid phone", phone: "+165055501234567890", slug: "cafe-test", field: "phone"},
		{name: "20 character invalid phone", phone: "+1650555012345678901", slug: "cafe-test", field: "phone"},
		{name: "21 character invalid phone", phone: "+16505550123456789012", slug: "cafe-test", field: "phone"},
		{name: "reserved subdomain", phone: "", slug: "admin", field: "restaurant_slug"},
		{name: "invalid subdomain syntax", phone: "", slug: "bad--slug", field: "restaurant_slug"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := &domain.RegisterRequest{
				Name: "Synthetic Owner", Email: "owner@example.test", Password: "synthetic-password", RestaurantName: "Synthetic Cafe",
				RestaurantSlug: test.slug, Phone: test.phone,
			}
			response, err := (&AuthUseCase{}).RegisterContext(context.Background(), request)
			var validation *domain.ValidationErrors
			if response != nil || !errors.As(err, &validation) || validation.Fields[test.field] == "" {
				t.Fatalf("RegisterContext() response=%v error=%v; want %q validation error", response, err, test.field)
			}
		})
	}
}

func TestRegistrationRejectsUnicodeDigitsInPhone(t *testing.T) {
	request := &domain.RegisterRequest{Phone: "+١٦٥٠٥٥٥٠١٢٣"}
	err := request.NormalizeAndValidate()
	if err == nil || !strings.Contains(err.Error(), "highlighted fields") {
		t.Fatalf("NormalizeAndValidate() = %v, want a safe field validation error", err)
	}
}
