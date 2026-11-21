package usecase

import (
	"testing"
	"time"
)

func TestCurrentNumericTextToCents_PreservesSignsAndPrecision(t *testing.T) {
	testCases := []struct {
		value string
		want  int64
	}{
		{value: "50.00", want: 5000},
		{value: "12.3", want: 1230},
		{value: "-40.05", want: -4005},
		{value: "0.00", want: 0},
	}
	for _, testCase := range testCases {
		got, err := currentNumericTextToCents(testCase.value)
		if err != nil || got != testCase.want {
			t.Fatalf("currentNumericTextToCents(%q) = %d, %v; want %d", testCase.value, got, err, testCase.want)
		}
	}
}

func TestFinancialCalendarDateUsesSuppliedCalendarDay(t *testing.T) {
	input := time.Date(2025, 1, 1, 23, 30, 0, 0, time.FixedZone("UTC-2", -2*60*60))
	got := financialCalendarDate(input)
	if got.Year() != 2025 || got.Month() != time.January || got.Day() != 1 || got.Hour() != 0 || got.Location() != time.UTC {
		t.Fatalf("financialCalendarDate(%v) = %v", input, got)
	}
}
