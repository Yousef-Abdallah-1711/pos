package commerce

import (
	"bytes"
	"os"
	"testing"
)

func TestPublishedOfferMatchesPolicy(t *testing.T) {
	frontend, err := os.ReadFile("../../../frontend/apps/landing/src/content/commercial-offer.json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(frontend, encodedOffer) {
		t.Fatal("frontend offer is stale; run scripts/sync-commercial-offer.py")
	}
	offer := Current()
	if offer.MonthlyPriceCents != 75000 || offer.AnnualPriceCents != 500000 || offer.TrialDays != 14 || offer.Currency != "EGP" || offer.CheckoutAvailable {
		t.Fatalf("unexpected commercial offer: %+v", offer)
	}
}
