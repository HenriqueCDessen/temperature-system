package main

import (
	"context"
	"log"
	"net/http"

	"temperature-system/service-b/internal/handler"
	"temperature-system/service-b/internal/tracing"

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
	mux.Handle("/weather", otelhttp.NewHandler(http.HandlerFunc(handler.WeatherHandler), "weatherHandler"))

	log.Println("Service B running on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
