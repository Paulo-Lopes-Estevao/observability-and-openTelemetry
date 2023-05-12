package tracing

import (
	"log"

	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type ZipkinOpenTel struct {
	ServiceName      string
	ServiceVersion   string
	ExporterEndpoint string
}

func NewZipkinOpenTel() *ZipkinOpenTel {
	return &ZipkinOpenTel{}
}

func (z *ZipkinOpenTel) Zipkin() *sdktrace.TracerProvider {

	exporter, err := zipkin.New(
		z.ExporterEndpoint,
		zipkin.WithSDKOptions(sdktrace.WithSampler(sdktrace.AlwaysSample())),
	)
	if err != nil {
		log.Fatal(err)
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("server-traces-http"),
		)),
	)

	return tp

}
