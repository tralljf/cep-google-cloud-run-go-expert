package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"cep-google-cloud-run/internal/app"
	"cep-google-cloud-run/internal/handler"
	"cep-google-cloud-run/internal/viacep"
	"cep-google-cloud-run/internal/weatherapi"
)

func main() {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("WEATHER_API_KEY is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	httpClient := &http.Client{Timeout: 5 * time.Second}
	service := app.NewService(
		viacep.NewClient(httpClient),
		weatherapi.NewClient(httpClient, apiKey),
	)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler.New(service),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Printf("listening on :%s", port)
	log.Fatal(server.ListenAndServe())
}
