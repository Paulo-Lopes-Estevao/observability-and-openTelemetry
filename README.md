# Observability and OpenTelemetry

## Introduction
Observability is a measure of how well internal states of a system can be inferred from knowledge of its external outputs. It helps to understand the health of the system and to debug issues. 
OpenTelemetry is a collection of tools, APIs, and SDKs used to instrument, generate, collect, and export telemetry data (metrics, logs, and traces) for analysis in order to understand your software's performance and behavior. 

This repository contains a sample application that is instrumented with OpenTelemetry and the telemetry data is collected and exported to Zipkin and Prometheus using OpenTelemetry Collector. The metrics are visualized using Grafana.

## Prerequisites
- [Docker](https://www.docker.com/products/docker-desktop)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Getting Started
1. Clone this repository
2. Run `docker-compose up` to start the application
3. Open `http://localhost:9411` to view the traces in Zipkin
4. Open `http://localhost:3000` to view the metrics in Grafana
6. Open `http://localhost:8080` to view the application

## References
- [OpenTelemetry](https://opentelemetry.io/)
- [OpenTelemetry Go](https://opentelemetry.io/docs/go/)
- [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/)
- [Grafana](https://grafana.com/)
- [Zipkin](https://zipkin.io/)
- [Prometheus](https://prometheus.io/)

## License
This project is licensed under the GNU AFFERO GENERAL PUBLIC LICENSE - see the [LICENSE](LICENSE) file for details.