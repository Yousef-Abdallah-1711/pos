package repository

import "testing"

func TestNumericTextToCents_PreservesDecimalPrecision(t *testing.T) {
	testCases := []struct {
		name  string
		value string
		want  int64
	}{
		{name: "whole amount", value: "1000", want: 100000},
		{name: "one fractional digit", value: "12.3", want: 1230},
		{name: "negative fractional amount", value: "-0.05", want: -5},
		{name: "explicit positive sign", value: "+7.09", want: 709},
		{name: "zero", value: "0.00", want: 0},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := numericTextToCents(testCase.value)
			if err != nil {
				t.Fatalf("numericTextToCents(%q): %v", testCase.value, err)
			}
			if got != testCase.want {
				t.Fatalf("numericTextToCents(%q) = %d, want %d", testCase.value, got, testCase.want)
			}
		})
	}
}

func TestNumericTextToCents_RejectsSubCentPrecision(t *testing.T) {
	if _, err := numericTextToCents("1.001"); err == nil {
		t.Fatal("numericTextToCents accepted a value below the bank-account schema's cent precision")
	}
}

func TestNumericTextToCents_RejectsValuesOutsideCentsRange(t *testing.T) {
	if _, err := numericTextToCents("92233720368547758.08"); err == nil {
		t.Fatal("numericTextToCents silently overflowed the supported signed-cent range")
	}
}
