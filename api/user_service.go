package api

import (
	"encoding/json"
	"github.com/Paulo-Lopes-Estevao/observability-wtih-openTelemetry/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"
)

type IUserService interface {
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	GetUsers(w http.ResponseWriter, r *http.Request)
	GetIdByUser(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type UserService struct {
	UserMetrics *metrics.Metrics
}

func NewUserService(m *metrics.Metrics) *UserService {
	return &UserService{
		UserMetrics: m,
	}
}

var tracer trace.Tracer

func init() {
	tracer = otel.Tracer("user_tracing")
}

func (s *UserService) Login(w http.ResponseWriter, r *http.Request) {
	s.UserMetrics.Users.Inc()
}

func (s *UserService) Logout(w http.ResponseWriter, r *http.Request) {
	s.UserMetrics.Users.Dec()
}

func (s *UserService) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := baggage.ContextWithoutBaggage(r.Context())
	ctx, span := tracer.Start(ctx, "get-users")
	defer span.End()

	users, err := GetUsers()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (s *UserService) GetIdByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := baggage.ContextWithoutBaggage(r.Context())
	ctx, span := tracer.Start(ctx, "get-id-by-user")
	defer span.End()

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	user, err := GetUser(id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	s.UserMetrics.UserInfo.With(prometheus.Labels{"name": user.Name, "email": user.Email}).Set(1)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (s *UserService) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := baggage.ContextWithoutBaggage(r.Context())
	ctx, span := tracer.Start(ctx, "create-user")
	defer span.End()

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	user, err = CreateUser(user)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (s *UserService) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := baggage.ContextWithoutBaggage(r.Context())
	ctx, span := tracer.Start(ctx, "update-user")
	defer span.End()

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	user, err = UpdateUser(user)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (s *UserService) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := baggage.ContextWithoutBaggage(r.Context())
	ctx, span := tracer.Start(ctx, "delete-user")
	defer span.End()

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	message, err := DeleteUser(id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message)
}
