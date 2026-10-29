package viacep

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cep-google-cloud-run/internal/app"
)

const defaultBaseURL = "https://viacep.com.br/ws"

type Client struct {
	httpClient *http.Client
	baseURL    string
}

type response struct {
	Localidade string          `json:"localidade"`
	Erro       json.RawMessage `json:"erro"`
}

func (r response) notFound() bool {
	erro := bytes.TrimSpace(r.Erro)
	if len(erro) == 0 {
		return false
	}

	switch string(erro) {
	case "false", `"false"`, "null":
		return false
	}

	return true
}

func NewClient(httpClient *http.Client) *Client {
	return NewClientWithBaseURL(httpClient, defaultBaseURL)
}

func NewClientWithBaseURL(httpClient *http.Client, baseURL string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{httpClient: httpClient, baseURL: baseURL}
}

func (c *Client) FindCity(ctx context.Context, zipcode string) (string, error) {
	url := fmt.Sprintf("%s/%s/json/", c.baseURL, zipcode)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", app.ErrZipcodeNotFound
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("viacep returned status %d", resp.StatusCode)
	}

	var data response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	if data.notFound() || data.Localidade == "" {
		return "", app.ErrZipcodeNotFound
	}

	return data.Localidade, nil
}
