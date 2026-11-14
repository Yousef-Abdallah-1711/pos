// Package commerce owns the approved offer for future subscriptions.
package commerce

import (
	_ "embed"
	"encoding/json"
)

//go:embed offer.json
var encodedOffer []byte

type Offer struct {
	Version           string `json:"version"`
	Currency          string `json:"currency"`
	TrialDays         int    `json:"trial_days"`
	MonthlyPriceCents int    `json:"monthly_price_cents"`
	AnnualPriceCents  int    `json:"annual_price_cents"`
	MonthlyPlanID     string `json:"monthly_plan_id"`
	AnnualPlanID      string `json:"annual_plan_id"`
	CheckoutAvailable bool   `json:"checkout_available"`
}

// Current returns a value copy; callers cannot mutate the commercial policy.
func Current() Offer {
	var offer Offer
	if err := json.Unmarshal(encodedOffer, &offer); err != nil {
		panic("invalid embedded commercial offer")
	}
	return offer
}
