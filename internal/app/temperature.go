package app

type Temperature struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func NewTemperature(celsius float64) Temperature {
	return Temperature{
		TempC: celsius,
		TempF: CelsiusToFahrenheit(celsius),
		TempK: CelsiusToKelvin(celsius),
	}
}

func CelsiusToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

func CelsiusToKelvin(celsius float64) float64 {
	return celsius + 273.15
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
