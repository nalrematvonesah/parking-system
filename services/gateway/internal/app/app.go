package app

import (
	"context"
	"net/http"
	"time"

	"gateway-service/internal/auth"
	"gateway-service/internal/clients"
	"gateway-service/internal/config"
	parkinghandler "gateway-service/internal/handler/parking"
	sessionhandler "gateway-service/internal/handler/session"
	userhandler "gateway-service/internal/handler/user"
	mw "gateway-service/internal/middleware"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type App struct {
	server  *http.Server
	clients *clients.Clients
}

func New(cfg config.Config, log *zap.Logger) (*App, error) {
	c, err := clients.Dial(
		cfg.UserAddr,
		cfg.ParkingAddr,
		cfg.SessionAddr,
	)
	if err != nil {
		return nil, err
	}

	jwtMgr := auth.New(
		cfg.JWTSecret,
		cfg.JWTTTL,
	)

	// handlers
	uh := userhandler.New(c.User, jwtMgr)
	ph := parkinghandler.New(c.Parking)
	sh := sessionhandler.New(c.Session)

	r := chi.NewRouter()

	// middleware
	r.Use(mw.CORS(cfg.AllowedOrigins))
	r.Use(mw.Recovery(log))
	r.Use(mw.AccessLog(log))

	// =========================
	// PUBLIC ROUTES
	// =========================

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.Post("/auth/register", uh.Register)
	r.Post("/auth/login", uh.Login)

	// IMPORTANT:
	// availability should be public
	r.Get("/slots/available", ph.Available)

	// =========================
	// PROTECTED ROUTES
	// =========================

	r.Group(func(r chi.Router) {

		r.Use(mw.JWTAuth(jwtMgr))

		// auth
		r.Post("/auth/logout", uh.Logout)

		// vehicles
		r.Post("/vehicles", uh.AddVehicle)
		r.Delete("/vehicles", uh.DeleteVehicle)
		r.Get("/vehicles", uh.ListVehicles)

		// sessions
		r.Post("/sessions/start", sh.Start)
		r.Get("/sessions/active", sh.Active)


		r.Post("/sessions/{id}/end", sh.End)

		r.Get("/sessions/{id}", sh.Get)
		r.Get("/sessions/{id}/price", sh.Price)
	})

	server := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	return &App{
		server:  server,
		clients: c,
	}, nil
}
func (a *App) Run() error {
	return a.server.ListenAndServe()
}

func (a *App) Stop(ctx context.Context) {
	_ = a.server.Shutdown(ctx)
	a.clients.Close()
}
