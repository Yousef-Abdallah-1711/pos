package domain

import "testing"

func TestPublicTrackingToken_RoundTripAndUniqueness(t *testing.T) {
	raw1, hash1, err := GeneratePublicTrackingToken()
	if err != nil {
		t.Fatal(err)
	}
	raw2, hash2, _ := GeneratePublicTrackingToken()
	if raw1 == raw2 || hash1 == hash2 {
		t.Fatal("tokens must be unique")
	}
	if len(raw1) < 40 {
		t.Fatalf("token too short to be unguessable: %d chars", len(raw1))
	}
	if hash1 == raw1 {
		t.Fatal("the stored value must be a hash, not the raw token")
	}
	if !VerifyPublicTrackingToken(raw1, hash1) {
		t.Fatal("the issued token must verify against its hash")
	}
}

func TestPublicTrackingToken_RejectsWrongEmptyAndForeignTokens(t *testing.T) {
	raw1, hash1, _ := GeneratePublicTrackingToken()
	raw2, hash2, _ := GeneratePublicTrackingToken()
	cases := []struct {
		name        string
		raw, stored string
	}{
		{"other order's token", raw2, hash1},
		{"empty token", "", hash1},
		{"empty stored hash (legacy order)", raw1, ""},
		{"both empty", "", ""},
		{"raw token used as hash", raw1, raw1},
		{"sequential-id guess", "1", hash1},
	}
	for _, c := range cases {
		if VerifyPublicTrackingToken(c.raw, c.stored) {
			t.Errorf("%s must not verify", c.name)
		}
	}
	if !VerifyPublicTrackingToken(raw2, hash2) {
		t.Fatal("sanity: own token verifies")
	}
}

func TestPhonesMatch(t *testing.T) {
	if !PhonesMatch("+20 100-000 0000", "201000000000") {
		t.Error("formatting differences must not matter")
	}
	if PhonesMatch("+201000000000", "+201000000001") {
		t.Error("different numbers must not match")
	}
	for _, weak := range [][2]string{{"", ""}, {"123", "123"}, {"abc", "abc"}} {
		if PhonesMatch(weak[0], weak[1]) {
			t.Errorf("weak/empty numbers must never match: %q", weak[0])
		}
	}
}
