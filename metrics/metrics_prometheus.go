package metrics

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
)

type IMetrics interface {
	PrometheusMiddleware(next http.Handler) http.Handler
}

type Metrics struct {
	Users          prometheus.Gauge
	CountTotalReqs prometheus.Counter
	UserInfo       *prometheus.GaugeVec
	TotalRequests  *prometheus.CounterVec
	ResponseTime   *prometheus.HistogramVec
	ResponseSize   prometheus.Histogram
	ResponseStatus *prometheus.CounterVec
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		Users: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "user_metrics",
			Name:      "connected_users",
			Help:      "Number of users",
		}),
		CountTotalReqs: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "user_metrics",
			Name:      "count_total_requests",
			Help:      "Total number of requests",
		}),
		UserInfo: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "user_metrics",
			Name:      "user_info",
			Help:      "User information",
		},
			[]string{"name", "email"},
		),
		TotalRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "user_metrics",
				Name:      "user_total_requests",
				Help:      "Total number of requests",
			},
			[]string{"path"},
		),
		ResponseTime: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "user_metrics",
			Name:      "user_response_time",
			Help:      "Response time",
		},
			[]string{"path", "status", "method"},
		),
		ResponseSize: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: "user_metrics",
			Name:      "user_response_size",
			Help:      "Response size",
		}),
		ResponseStatus: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "user_metrics",
			Name:      "user_response_status",
			Help:      "Response status",
		},
			[]string{"status"},
		),
	}
	reg.MustRegister(m.Users, m.CountTotalReqs, m.UserInfo, m.TotalRequests, m.ResponseTime, m.ResponseSize, m.ResponseStatus)
	return m
}

func (m *Metrics) PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		rw := NewResponseWriter(w)
		next.ServeHTTP(w, r)

		statusCode := rw.statusCode

		timer := prometheus.NewTimer(m.ResponseTime.WithLabelValues(path, strconv.Itoa(statusCode), r.Method))

		m.TotalRequests.WithLabelValues(path).Inc()

		m.CountTotalReqs.Inc()

		m.ResponseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()

		size := w.Header().Get("Content-Length")
		m.ResponseSize.Observe(float64(len(size)))

		timer.ObserveDuration()
	})
}
