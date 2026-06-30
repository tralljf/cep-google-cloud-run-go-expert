package app

import (
	"context"
	"errors"
	"testing"
)

type fakeZipcodeClient struct {
	city string
	err  error
}

func (f fakeZipcodeClient) FindCity(context.Context, string) (string, error) {
	return f.city, f.err
}

type fakeWeatherClient struct {
	tempC float64
	err   error
}

func (f fakeWeatherClient) CurrentTempC(context.Context, string) (float64, error) {
	return f.tempC, f.err
}

func TestServiceWeatherByZipcodeSuccess(t *testing.T) {
	service := NewService(fakeZipcodeClient{city: "Sao Paulo"}, fakeWeatherClient{tempC: 20})

	temp, err := service.WeatherByZipcode(context.Background(), "01001000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if temp.TempC != 20 || temp.TempF != 68 || temp.TempK != 293.15 {
		t.Fatalf("unexpected temperature: %+v", temp)
	}
}

func TestServiceWeatherByZipcodeInvalidZipcode(t *testing.T) {
	service := NewService(fakeZipcodeClient{}, fakeWeatherClient{})

	_, err := service.WeatherByZipcode(context.Background(), "01001-000")
	if !errors.Is(err, ErrInvalidZipcode) {
		t.Fatalf("expected ErrInvalidZipcode, got %v", err)
	}
}
