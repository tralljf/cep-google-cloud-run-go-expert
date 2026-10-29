package app

import (
	"encoding/json"
	"testing"
)

func TestTemperatureConversions(t *testing.T) {
	tests := []struct {
		name    string
		celsius float64
		wantC   float64
		wantF   float64
		wantK   float64
	}{
		{name: "kelvin floating point noise", celsius: 20.4, wantC: 20.4, wantF: 68.72, wantK: 293.55},
		{name: "fahrenheit floating point noise", celsius: 18.5, wantC: 18.5, wantF: 65.3, wantK: 291.65},
		{name: "both scales with noise", celsius: 28.5, wantC: 28.5, wantF: 83.3, wantK: 301.65},
		{name: "celsius with many decimals", celsius: 21.456, wantC: 21.46, wantF: 70.62, wantK: 294.61},
		{name: "negative", celsius: -3.3, wantC: -3.3, wantF: 26.06, wantK: 269.85},
		{name: "zero", celsius: 0, wantC: 0, wantF: 32, wantK: 273.15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp := NewTemperature(tt.celsius)

			if temp.TempC != tt.wantC {
				t.Fatalf("expected celsius %v, got %v", tt.wantC, temp.TempC)
			}
			if temp.TempF != tt.wantF {
				t.Fatalf("expected fahrenheit %v, got %v", tt.wantF, temp.TempF)
			}
			if temp.TempK != tt.wantK {
				t.Fatalf("expected kelvin %v, got %v", tt.wantK, temp.TempK)
			}
		})
	}
}

func TestTemperatureJSONHasNoFloatingPointTail(t *testing.T) {
	tests := []struct {
		celsius float64
		want    string
	}{
		{celsius: 20.4, want: `{"temp_C":20.4,"temp_F":68.72,"temp_K":293.55}`},
		{celsius: 18.5, want: `{"temp_C":18.5,"temp_F":65.3,"temp_K":291.65}`},
		{celsius: 28.5, want: `{"temp_C":28.5,"temp_F":83.3,"temp_K":301.65}`},
	}

	for _, tt := range tests {
		body, err := json.Marshal(NewTemperature(tt.celsius))
		if err != nil {
			t.Fatalf("marshal temperature: %v", err)
		}
		if string(body) != tt.want {
			t.Fatalf("NewTemperature(%v) JSON = %s, want %s", tt.celsius, body, tt.want)
		}
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
