package api

import (
	"github.com/Paulo-Lopes-Estevao/observability-wtih-openTelemetry/metrics"
	"github.com/Paulo-Lopes-Estevao/observability-wtih-openTelemetry/tracing"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"log"
	"net/http"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() {

	mux := mux.NewRouter()

	mux.Use(metrics.TotalRequestCounter)

	// Create a metrics registry. This registry is used by the Prometheus
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
		metrics.RequestCounter,
	)
	m := metrics.NewMetrics(reg)

	// Create a new OpenTelemetry tracing provider.
	opt := tracing.NewZipkinOpenTel()
	opt.ExporterEndpoint = "http://zipkin:9411/api/v2/spans"
	opt.ServiceName = "ServerHTTP"
	tp := opt.Zipkin()
	otel.SetTracerProvider(tp)

	router := NewRouter(NewUserService(m), mux, m)
	router.Init()

	handler := otelhttp.NewHandler(mux, "server",
		otelhttp.WithTracerProvider(tp),
	)

	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	mux.Handle("/metrics", promHandler)

	log.Printf("Starting server on port %s...\n", ":8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
