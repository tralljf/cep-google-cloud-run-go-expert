package weatherapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const defaultBaseURL = "https://api.weatherapi.com/v1/current.json"

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

type response struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func NewClient(httpClient *http.Client, apiKey string) *Client {
	return NewClientWithBaseURL(httpClient, defaultBaseURL, apiKey)
}

func NewClientWithBaseURL(httpClient *http.Client, baseURL, apiKey string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{httpClient: httpClient, baseURL: baseURL, apiKey: apiKey}
}

func (c *Client) CurrentTempC(ctx context.Context, location string) (float64, error) {
	values := url.Values{}
	values.Set("key", c.apiKey)
	values.Set("q", location+",Brazil")
	values.Set("aqi", "no")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"?"+values.Encode(), nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return 0, fmt.Errorf("weatherapi returned status %d", resp.StatusCode)
	}

	var data response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	return data.Current.TempC, nil
}
