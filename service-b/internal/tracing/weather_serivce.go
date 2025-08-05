package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"go.opentelemetry.io/otel"
)

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type WeatherResult struct {
	TempC float64
}

func FetchCityByCEP(ctx context.Context, cep string) (string, error) {
	ctx, span := otel.Tracer("service-b").Start(ctx, "FetchCityByCEP")
	defer span.End()

	resp, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep))
	if err != nil || resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch CEP")
	}
	defer resp.Body.Close()

	var data struct {
		Localidade string `json:"localidade"`
		Uf         string `json:"uf"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil || data.Localidade == "" {
		return "", fmt.Errorf("invalid response from viacep")
	}

	return fmt.Sprintf("%s,%s", data.Localidade, data.Uf), nil
}

func FetchWeatherByCity(ctx context.Context, location string) (*WeatherResult, error) {
	ctx, span := otel.Tracer("service-b").Start(ctx, "FetchWeatherByCity")
	defer span.End()

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("missing WEATHER_API_KEY")
	}

	query := url.QueryEscape(location)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, query)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("weatherapi error")
	}
	defer resp.Body.Close()

	var data WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &WeatherResult{TempC: data.Current.TempC}, nil
}
