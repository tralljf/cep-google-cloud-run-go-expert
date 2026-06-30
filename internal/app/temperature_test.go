package app

import "testing"

func TestTemperatureConversions(t *testing.T) {
	temp := NewTemperature(28.5)

	if temp.TempC != 28.5 {
		t.Fatalf("expected celsius 28.5, got %v", temp.TempC)
	}
	if temp.TempF != 83.30000000000001 {
		t.Fatalf("expected fahrenheit 83.30000000000001, got %v", temp.TempF)
	}
	if temp.TempK != 301.65 {
		t.Fatalf("expected kelvin 301.65, got %v", temp.TempK)
	}
}

func TestValidZipcode(t *testing.T) {
	tests := []struct {
		name    string
		zipcode string
		want    bool
	}{
		{name: "valid", zipcode: "01001000", want: true},
		{name: "too short", zipcode: "0100100", want: false},
		{name: "too long", zipcode: "010010000", want: false},
		{name: "with dash", zipcode: "01001-000", want: false},
		{name: "with letters", zipcode: "abcdefgh", want: false},
		{name: "empty", zipcode: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidZipcode(tt.zipcode); got != tt.want {
				t.Fatalf("ValidZipcode(%q) = %v, want %v", tt.zipcode, got, tt.want)
			}
		})
	}
}
