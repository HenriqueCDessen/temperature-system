package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

type request struct {
	CEP string `json:"cep"`
}

func TemperatureHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("service-a").Start(r.Context(), "handleTemperature")
	defer span.End()

	var input request
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || !isValidCEP(input.CEP) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	serviceBUrl := os.Getenv("SERVICE_B_URL")
	if serviceBUrl == "" {
		serviceBUrl = "http://localhost:8081/weather"
	}

	payload, _ := json.Marshal(input)
	req, err := http.NewRequestWithContext(ctx, "POST", serviceBUrl, bytes.NewReader(payload))
	if err != nil {
		http.Error(w, "failed to create request to service B", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "failed to reach service B", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func isValidCEP(cep string) bool {
	cep = strings.TrimSpace(cep)
	match, _ := regexp.MatchString(`^\d{8}$`, cep)
	println("[Validando CEP]", cep, "-> válido?", match)
	return match
}
