package app

import (
	"context"
	"errors"
)

var (
	ErrInvalidZipcode  = errors.New("invalid zipcode")
	ErrZipcodeNotFound = errors.New("can not find zipcode")
)

type ZipcodeClient interface {
	FindCity(ctx context.Context, zipcode string) (string, error)
}

type WeatherClient interface {
	CurrentTempC(ctx context.Context, location string) (float64, error)
}

type Service struct {
	zipcodes ZipcodeClient
	weather  WeatherClient
}

func NewService(zipcodes ZipcodeClient, weather WeatherClient) *Service {
	return &Service{zipcodes: zipcodes, weather: weather}
}

func (s *Service) WeatherByZipcode(ctx context.Context, zipcode string) (Temperature, error) {
	if !ValidZipcode(zipcode) {
		return Temperature{}, ErrInvalidZipcode
	}

	city, err := s.zipcodes.FindCity(ctx, zipcode)
	if err != nil {
		return Temperature{}, err
	}

	tempC, err := s.weather.CurrentTempC(ctx, city)
	if err != nil {
		return Temperature{}, err
	}

	return NewTemperature(tempC), nil
}
