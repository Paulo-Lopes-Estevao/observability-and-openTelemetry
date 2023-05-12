package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

var RequestCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "metrics",
	Name:      "total_requests",
	Help:      "Total number of requests",
})

func TotalRequestCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Incrementa o contador de requisições
		RequestCounter.Inc()

		// Chama o próximo manipulador
		next.ServeHTTP(w, r)
	})
}
