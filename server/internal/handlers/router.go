package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/cgrente/sequence-insights-platform/server/internal/auth"
	"github.com/cgrente/sequence-insights-platform/server/internal/jobs"
	"github.com/cgrente/sequence-insights-platform/server/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

/*
router.go wires HTTP routes + middleware.

We keep this file "composition only":
- no business logic
- no DB logic
- only glue
*/

type RouterParams struct {
	Log    *slog.Logger
	Store  storage.Store
	Queue  *jobs.Queue
	APIKey string
}

func NewRouter(p RouterParams) http.Handler {
	r := chi.NewRouter()

	// Middlewares
	logger := httplog.NewLogger("sequence-insights-api", httplog.Options{
		LogLevel:        slog.LevelInfo,
		JSON:            true,
		Concise:         true,
		RequestHeaders:  true,
		TimeFieldFormat: time.RFC3339Nano,
	})
	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))

	// Auth (optional)
	r.Use(auth.BearerToken(p.APIKey))

	// Routes
	r.Get("/health", HealthHandler{Store: p.Store}.handle)

	r.Route("/v1/sequences", func(sr chi.Router) {
		sr.Mount("/", SequencesHandler{Store: p.Store, Queue: p.Queue}.Routes())
	})

	return r
}
