package viacep

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"cep-google-cloud-run/internal/app"
)

func TestFindCitySuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/01001000/json/" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"localidade":"Sao Paulo"}`))
	}))
	defer server.Close()

	client := NewClientWithBaseURL(server.Client(), server.URL)
	city, err := client.FindCity(context.Background(), "01001000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if city != "Sao Paulo" {
		t.Fatalf("expected Sao Paulo, got %q", city)
	}
}

func TestFindCityNotFound(t *testing.T) {
	payloads := map[string]string{
		"string erro": `{"erro":"true"}`,
		"bool erro":   `{"erro":true}`,
	}

	for name, payload := range payloads {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(payload))
			}))
			defer server.Close()

			client := NewClientWithBaseURL(server.Client(), server.URL)
			_, err := client.FindCity(context.Background(), "01001003")
			if !errors.Is(err, app.ErrZipcodeNotFound) {
				t.Fatalf("expected ErrZipcodeNotFound, got %v", err)
			}
		})
	}
}

func TestFindCityErroFalse(t *testing.T) {
	payloads := map[string]string{
		"string false": `{"localidade":"Sao Paulo","erro":"false"}`,
		"bool false":   `{"localidade":"Sao Paulo","erro":false}`,
	}

	for name, payload := range payloads {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(payload))
			}))
			defer server.Close()

			client := NewClientWithBaseURL(server.Client(), server.URL)
			city, err := client.FindCity(context.Background(), "01001000")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if city != "Sao Paulo" {
				t.Fatalf("expected Sao Paulo, got %q", city)
			}
		})
	}
}
