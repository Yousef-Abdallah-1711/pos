package domain

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"unicode"

	"github.com/nyaruka/phonenumbers"
)

const (
	RegistrationNameMinRunes     = 2
	RegistrationNameMaxRunes     = 255
	RegistrationEmailMaxRunes    = 255
	RegistrationPasswordMinRunes = 6
	RegistrationPasswordMaxBytes = 72
)

var (
	ErrEmailAlreadyExists        = errors.New("email already exists")
	ErrSubdomainAlreadyExists    = errors.New("subdomain already exists")
	ErrRestaurantNameConflict    = errors.New("restaurant name already exists")
	ErrAuthDependencyUnavailable = errors.New("authentication dependency unavailable")
	ErrTokenGeneration           = errors.New("authentication token generation failed")
)

// ValidationErrors contains safe, user-correctable field messages.
type ValidationErrors struct {
	Fields map[string]string
}

func (e *ValidationErrors) Error() string { return "Please review the highlighted fields." }

// NormalizeAndValidate applies the registration contract before persistence.
// Name limits count Unicode code points; phone is stored as canonical E.164.
func (r *RegisterRequest) NormalizeAndValidate() error {
	fields := make(map[string]string)
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.RestaurantName = strings.TrimSpace(r.RestaurantName)
	r.RestaurantSlug = strings.ToLower(strings.TrimSpace(r.RestaurantSlug))

	validateRequiredName(fields, "name", r.Name)
	validateRequiredName(fields, "restaurant_name", r.RestaurantName)
	if r.Email == "" {
		fields["email"] = "Email is required."
	} else if runeLength(r.Email) > RegistrationEmailMaxRunes {
		fields["email"] = "Email must be 255 characters or fewer."
	} else if parsed, err := mail.ParseAddress(r.Email); err != nil || parsed.Address != r.Email {
		fields["email"] = "Enter a valid email address."
	}
	if r.Password == "" {
		fields["password"] = "Password is required."
	} else if runeLength(r.Password) < RegistrationPasswordMinRunes {
		fields["password"] = "Password must be at least 6 characters."
	} else if len([]byte(r.Password)) > RegistrationPasswordMaxBytes {
		fields["password"] = "Password must be 72 bytes or fewer."
	}
	phone, err := NormalizeInternationalPhone(r.Phone)
	if err != nil {
		fields["phone"] = "Enter a valid international phone number, including its + country code."
	} else {
		r.Phone = phone
	}
	if len(fields) > 0 {
		return &ValidationErrors{Fields: fields}
	}
	return nil
}

// NormalizeAndValidate applies the login request's safe syntax and size rules.
// Credential mismatches remain a generic authentication failure after this
// user-correctable request validation has passed.
func (r *LoginRequest) NormalizeAndValidate() error {
	if r == nil {
		return &ValidationErrors{Fields: map[string]string{"email": "Email is required.", "password": "Password is required."}}
	}
	fields := make(map[string]string)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	if r.Email == "" {
		fields["email"] = "Email is required."
	} else if runeLength(r.Email) > RegistrationEmailMaxRunes {
		fields["email"] = "Email must be 255 characters or fewer."
	} else if parsed, err := mail.ParseAddress(r.Email); err != nil || parsed.Address != r.Email {
		fields["email"] = "Enter a valid email address."
	}
	if r.Password == "" {
		fields["password"] = "Password is required."
	} else if len([]byte(r.Password)) > RegistrationPasswordMaxBytes {
		fields["password"] = "Password must be 72 bytes or fewer."
	}
	if len(fields) > 0 {
		return &ValidationErrors{Fields: fields}
	}
	return nil
}

func validateRequiredName(fields map[string]string, field, value string) {
	length := runeLength(value)
	if length == 0 {
		fields[field] = "This field is required."
	} else if length < RegistrationNameMinRunes {
		fields[field] = "Enter at least 2 characters."
	} else if length > RegistrationNameMaxRunes {
		fields[field] = "Enter 255 characters or fewer."
	}
}

func runeLength(value string) int { return len([]rune(value)) }

// NormalizeInternationalPhone accepts formatted international input only when
// the resulting number is valid without an implicit country/region guess.
func NormalizeInternationalPhone(input string) (string, error) {
	value := strings.TrimSpace(input)
	if value == "" {
		return "", nil
	}
	var normalized strings.Builder
	for index, char := range value {
		switch {
		case char >= '0' && char <= '9':
			normalized.WriteRune(char)
		case char == '+' && index == 0:
			normalized.WriteRune(char)
		case unicode.IsSpace(char), char == '-', char == '(', char == ')', char == '.':
			// Common presentation separators are discarded before parsing.
		default:
			return "", fmt.Errorf("unsupported phone character")
		}
	}
	canonical := normalized.String()
	if !strings.HasPrefix(canonical, "+") {
		return "", fmt.Errorf("international phone must start with +")
	}
	number, err := phonenumbers.Parse(canonical, "ZZ")
	if err != nil || !phonenumbers.IsValidNumber(number) {
		return "", fmt.Errorf("invalid international phone number")
	}
	formatted := phonenumbers.Format(number, phonenumbers.E164)
	if len(strings.TrimPrefix(formatted, "+")) > 15 {
		return "", fmt.Errorf("international phone exceeds E.164 digit limit")
	}
	return formatted, nil
}
