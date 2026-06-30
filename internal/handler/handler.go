package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"cep-google-cloud-run/internal/app"
)

type service interface {
	WeatherByZipcode(ctx context.Context, zipcode string) (app.Temperature, error)
}

type Handler struct {
	service service
}

func New(service service) http.Handler {
	h := &Handler{service: service}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /weather/{zipcode}", h.weatherByZipcode)
	mux.HandleFunc("GET /healthz", h.healthz)
	mux.HandleFunc("/", notFound)

	return mux
}

func (h *Handler) weatherByZipcode(w http.ResponseWriter, r *http.Request) {
	zipcode := strings.TrimSpace(r.PathValue("zipcode"))

	temperature, err := h.service.WeatherByZipcode(r.Context(), zipcode)
	if err != nil {
		switch {
		case errors.Is(err, app.ErrInvalidZipcode):
			writeText(w, http.StatusUnprocessableEntity, "invalid zipcode")
		case errors.Is(err, app.ErrZipcodeNotFound):
			writeText(w, http.StatusNotFound, "can not find zipcode")
		default:
			writeText(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(temperature)
}

func (h *Handler) healthz(w http.ResponseWriter, _ *http.Request) {
	writeText(w, http.StatusOK, "ok")
}

func notFound(w http.ResponseWriter, _ *http.Request) {
	writeText(w, http.StatusNotFound, "not found")
}

func writeText(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(message))
}
