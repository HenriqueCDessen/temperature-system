package handler

import (
	"encoding/json"
	"net/http"

	"temperature-system/service-b/internal/service"

	"go.opentelemetry.io/otel"
)

type CEPRequest struct {
	CEP string `json:"cep"`
}

type WeatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("service-b").Start(r.Context(), "WeatherHandler")
	defer span.End()

	var req CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	city, err := service.FetchCityByCEP(ctx, req.CEP)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	result, err := service.FetchWeatherByCity(ctx, city)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := WeatherResponse{
		City:  city,
		TempC: result.TempC,
		TempF: result.TempC*1.8 + 32,
		TempK: result.TempC + 273,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
