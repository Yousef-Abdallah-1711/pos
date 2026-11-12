package http

import (
	"fmt"
	"strings"
)

type stringLengthRule struct {
	Name  string
	Value string
	Max   int
}

func validateStringLengths(rules ...stringLengthRule) error {
	for _, rule := range rules {
		if rule.Max <= 0 {
			continue
		}
		if len([]rune(strings.TrimSpace(rule.Value))) > rule.Max {
			return fmt.Errorf("%s exceeds maximum length of %d", rule.Name, rule.Max)
		}
	}
	return nil
}
