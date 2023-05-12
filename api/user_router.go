package api

import (
	"github.com/Paulo-Lopes-Estevao/observability-wtih-openTelemetry/metrics"
	"github.com/gorilla/mux"
)

type Router struct {
	IuserService IUserService
	IMetrics     metrics.IMetrics
	SeverMux     *mux.Router
}

func NewRouter(iuserService IUserService, serverMux *mux.Router, iMetrics metrics.IMetrics) *Router {
	return &Router{
		IuserService: iuserService,
		SeverMux:     serverMux,
		IMetrics:     iMetrics,
	}
}

func (r *Router) Init() {

	r.SeverMux.HandleFunc("/login", r.IuserService.Login)
	r.SeverMux.HandleFunc("/logout", r.IuserService.Logout)
	r.SeverMux.HandleFunc("/users", r.IuserService.GetUsers)
	r.SeverMux.HandleFunc("/user", r.IuserService.GetIdByUser)
	r.SeverMux.HandleFunc("/user/create", r.IuserService.CreateUser)
	r.SeverMux.HandleFunc("/user/update", r.IuserService.UpdateUser)
	r.SeverMux.HandleFunc("/user/delete", r.IuserService.DeleteUser)
	r.SeverMux.Use(r.IMetrics.PrometheusMiddleware)

}
