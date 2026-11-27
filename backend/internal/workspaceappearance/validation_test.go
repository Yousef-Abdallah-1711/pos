package workspaceappearance

import "testing"

func TestValidateOverrides(t *testing.T) {
	tests := []struct {
		name, input string
		valid       bool
	}{
		{"defaults", `{}`, true},
		{"valid color", `{"colors":{"light":{"primary":"#722116"}}}`, true},
		{"unknown key", `{"colors":{"light":{"unknown":"#FFFFFF"}}}`, false},
		{"unsafe CSS", `{"colors":{"light":{"primary":"url(evil)"}}}`, false},
		{"unknown font", `{"typography":{"english":"remote"}}`, false},
		{"unreadable", `{"colors":{"light":{"foreground":"#F8F2E8"}}}`, false},
		{"null root", `null`, false},
		{"extra document", `{} {}`, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateOverrides([]byte(tc.input))
			if (err == nil) != tc.valid {
				t.Fatalf("valid=%v: %v", tc.valid, err)
			}
		})
	}
}
