package app

import "math"

type Temperature struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func NewTemperature(celsius float64) Temperature {
	return Temperature{
		TempC: roundTwoDecimals(celsius),
		TempF: roundTwoDecimals(CelsiusToFahrenheit(celsius)),
		TempK: roundTwoDecimals(CelsiusToKelvin(celsius)),
	}
}

func CelsiusToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

func CelsiusToKelvin(celsius float64) float64 {
	return celsius + 273.15
}

func roundTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}

func ValidZipcode(zipcode string) bool {
	if len(zipcode) != 8 {
		return false
	}

	for _, char := range zipcode {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}
