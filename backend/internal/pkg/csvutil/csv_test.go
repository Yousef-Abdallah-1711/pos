package csvutil

import "testing"

func TestSafeTextCellPreventsSpreadsheetFormulaInjection(t *testing.T) {
	tests := []struct {
		name, input, expected string
	}{
		{name: "equals", input: `=HYPERLINK("https://example.invalid")`, expected: `'=HYPERLINK("https://example.invalid")`},
		{name: "plus after whitespace", input: "  +SUM(1,1)", expected: "'  +SUM(1,1)"},
		{name: "equals after tab", input: "\t=1+1", expected: "'\t=1+1"},
		{name: "minus", input: "-1+1", expected: "'-1+1"},
		{name: "at", input: "@SUM(A1:A2)", expected: "'@SUM(A1:A2)"},
		{name: "ordinary text", input: "Fresh salad", expected: "Fresh salad"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if actual := SafeTextCell(test.input); actual != test.expected {
				t.Fatalf("SafeTextCell(%q)=%q, want %q", test.input, actual, test.expected)
			}
		})
	}
}
