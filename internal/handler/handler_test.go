package handler

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"cep-google-cloud-run/internal/app"
)

type zipcodeClient struct {
	city string
	err  error
}

func (z zipcodeClient) FindCity(context.Context, string) (string, error) {
	return z.city, z.err
}

type weatherClient struct {
	tempC float64
	err   error
}

func (w weatherClient) CurrentTempC(context.Context, string) (float64, error) {
	return w.tempC, w.err
}

func TestWeatherByZipcodeSuccess(t *testing.T) {
	server := New(app.NewService(zipcodeClient{city: "Sao Paulo"}, weatherClient{tempC: 28.5}))

	req := httptest.NewRequest(http.MethodGet, "/weather/01001000", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d with body %q", rec.Code, rec.Body.String())
	}

	var got app.Temperature
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if !closeTo(got.TempC, 28.5) || !closeTo(got.TempF, 83.3) || !closeTo(got.TempK, 301.65) {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestWeatherByZipcodeInvalidZipcode(t *testing.T) {
	server := New(app.NewService(zipcodeClient{}, weatherClient{}))

	req := httptest.NewRequest(http.MethodGet, "/weather/01001-000", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", rec.Code)
	}
	if rec.Body.String() != "invalid zipcode" {
		t.Fatalf("unexpected body %q", rec.Body.String())
	}
}

func TestWeatherByZipcodeNotFound(t *testing.T) {
	server := New(app.NewService(zipcodeClient{err: app.ErrZipcodeNotFound}, weatherClient{}))

	req := httptest.NewRequest(http.MethodGet, "/weather/99999999", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
	if rec.Body.String() != "can not find zipcode" {
		t.Fatalf("unexpected body %q", rec.Body.String())
	}
}

func TestWeatherByZipcodeWeatherProviderError(t *testing.T) {
	server := New(app.NewService(zipcodeClient{city: "Sao Paulo"}, weatherClient{err: errors.New("provider down")}))

	req := httptest.NewRequest(http.MethodGet, "/weather/01001000", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
}

func closeTo(got, want float64) bool {
	return math.Abs(got-want) < 0.000001
}
