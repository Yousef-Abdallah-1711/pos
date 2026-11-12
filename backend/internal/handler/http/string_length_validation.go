package http

import (
	"fmt"
	"unicode/utf8"
)

type stringLengthRule struct {
	Name  string
	Value string
	Max   int
}

func validateStringLengths(rules ...stringLengthRule) error {
	for _, rule := range rules {
		if rule.Max <= 0 {
			return fmt.Errorf("invalid maximum length for %s", rule.Name)
		}
		if utf8.RuneCountInString(rule.Value) > rule.Max {
			return fmt.Errorf("%s exceeds the maximum length of %d characters", rule.Name, rule.Max)
		}
	}
	return nil
}
