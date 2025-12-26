package handlers

import (
	"net/http"
	"strings"

	"github.com/cgrente/sequence-insights-platform/server/internal/httpx"
	"github.com/cgrente/sequence-insights-platform/server/internal/jobs"
	"github.com/cgrente/sequence-insights-platform/server/internal/models"
	"github.com/cgrente/sequence-insights-platform/server/internal/services"
	"github.com/cgrente/sequence-insights-platform/server/internal/storage"
	"github.com/go-chi/chi/v5"
)

type SequencesHandler struct {
	Store storage.Store
	Queue *jobs.Queue
}

type ingestRequest struct {
	Values []int64 `json:"values"`
}

type ingestResponse struct {
	Sequence models.Sequence `json:"sequence"`
	Queued   bool            `json:"queued"`
}

func (h SequencesHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/ingest", h.ingest)
	r.Get("/{id}", h.getByID)

	return r
}

func (h SequencesHandler) ingest(w http.ResponseWriter, r *http.Request) {
	var req ingestRequest
	if err := httpx.ReadJSON(r, &req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	if len(req.Values) == 0 {
		http.Error(w, "values must not be empty", http.StatusBadRequest)
		return
	}

	count, sumFourth, min, max := services.ComputeMetrics(req.Values)

	seq := models.Sequence{
		Values:                req.Values,
		Count:                 count,
		SumFourthPowersNonPos: sumFourth,
		Min:                   min,
		Max:                   max,
		Processed:             false,
	}

	created, err := h.Store.CreateSequence(r.Context(), seq)
	if err != nil {
		http.Error(w, "failed to persist sequence", http.StatusInternalServerError)
		return
	}

	queued := false
	if h.Queue != nil {
		queued = h.Queue.Enqueue(jobs.Job{SequenceID: created.ID})
	}

	httpx.WriteJSON(w, http.StatusCreated, ingestResponse{
		Sequence: created,
		Queued:   queued,
	})
}

func (h SequencesHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	seq, ok, err := h.Store.GetSequence(r.Context(), id)
	if err != nil {
		http.Error(w, "failed to load sequence", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"sequence": seq})
}
