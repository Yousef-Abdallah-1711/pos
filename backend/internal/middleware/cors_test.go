package middleware

import "testing"

func TestIsOriginAllowedSupportsSingleRestaurantSubdomain(t *testing.T) {
	corsOrigins = map[string]struct{}{
		"https://pos.rakmyat.online":                      {},
		"https://*.restaurants.pos.rakmyat.online": {},
	}
	tests := []struct {
		origin string
		want   bool
	}{
		{origin: "https://pos.rakmyat.online", want: true},
		{origin: "https://cafe.restaurants.pos.rakmyat.online", want: true},
		{origin: "https://a-b.restaurants.pos.rakmyat.online", want: true},
		{origin: "https://restaurants.pos.rakmyat.online", want: false},
		{origin: "https://a.b.restaurants.pos.rakmyat.online", want: false},
		{origin: "https://cafe.restaurants.pos.rakmyat.online.evil.test", want: false},
		{origin: "https://cafe-restaurants.pos.rakmyat.online", want: false},
		{origin: "http://cafe.restaurants.pos.rakmyat.online", want: false},
		{origin: "https://cafe.restaurants.pos.rakmyat.online:444", want: false},
		{origin: "https://user@cafe.restaurants.pos.rakmyat.online", want: false},
		{origin: "https://cafe.restaurants.pos.rakmyat.online/path", want: false},
	}
	for _, test := range tests {
		t.Run(test.origin, func(t *testing.T) {
			if got := isOriginAllowed(test.origin); got != test.want {
				t.Fatalf("isOriginAllowed(%q) = %t, want %t", test.origin, got, test.want)
			}
		})
	}
}
