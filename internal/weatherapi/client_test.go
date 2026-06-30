package weatherapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCurrentTempC(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("key") != "test-key" {
			t.Fatalf("unexpected key: %q", r.URL.Query().Get("key"))
		}
		if r.URL.Query().Get("q") != "Sao Paulo,Brazil" {
			t.Fatalf("unexpected q: %q", r.URL.Query().Get("q"))
		}
		if r.URL.Query().Get("aqi") != "no" {
			t.Fatalf("unexpected aqi: %q", r.URL.Query().Get("aqi"))
		}
		_, _ = w.Write([]byte(`{"current":{"temp_c":28.5}}`))
	}))
	defer server.Close()

	client := NewClientWithBaseURL(server.Client(), server.URL, "test-key")
	tempC, err := client.CurrentTempC(context.Background(), "Sao Paulo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tempC != 28.5 {
		t.Fatalf("expected 28.5, got %v", tempC)
	}
}
