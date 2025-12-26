package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/cgrente/sequence-insights-platform/server/internal/httpx"
	"github.com/cgrente/sequence-insights-platform/server/internal/storage"
)

type HealthHandler struct {
	Store storage.Store
}

func (h HealthHandler) Routes() http.Handler {
	return http.HandlerFunc(h.handle)
}

func (h HealthHandler) handle(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	err := h.Store.Health(ctx)
	if err != nil {
		httpx.WriteJSON(w, http.StatusServiceUnavailable, map[string]any{
			"status": "down",
			"error":  err.Error(),
		})
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}
