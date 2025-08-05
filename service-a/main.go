package main

import (
	"context"
	"log"
	"net/http"

	"temperature-system/service-a/internal/handler"
	"temperature-system/service-a/internal/tracing"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	ctx := context.Background()
	tp, err := tracing.SetupTracerProvider(ctx)
	if err != nil {
		log.Fatalf("failed to initialize tracer provider: %v", err)
	}
	defer tp.Shutdown(ctx)

	mux := http.NewServeMux()
	mux.Handle("/temperature", otelhttp.NewHandler(http.HandlerFunc(handler.TemperatureHandler), "temperatureHandler"))

	log.Println("Service A running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
