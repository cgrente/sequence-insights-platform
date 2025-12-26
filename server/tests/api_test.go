package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"log/slog"

	"github.com/cgrente/sequence-insights-platform/server/internal/handlers"
	"github.com/cgrente/sequence-insights-platform/server/internal/jobs"
	"github.com/cgrente/sequence-insights-platform/server/internal/logging"
	"github.com/cgrente/sequence-insights-platform/server/internal/storage"
)

func TestHealth(t *testing.T) {
	store := storage.NewMemoryStore()
	r := handlers.NewRouter(handlers.RouterParams{
		Log:    logging.NewJSONLogger(slog.LevelInfo),
		Store:  store,
		Queue:  jobs.NewQueue(logging.NewJSONLogger(slog.LevelInfo), store, 10),
		APIKey: "",
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestIngestAndGet(t *testing.T) {
	store := storage.NewMemoryStore()
	queue := jobs.NewQueue(logging.NewJSONLogger(slog.LevelInfo), store, 10)

	r := handlers.NewRouter(handlers.RouterParams{
		Log:    logging.NewJSONLogger(slog.LevelInfo),
		Store:  store,
		Queue:  queue,
		APIKey: "",
	})

	body, _ := json.Marshal(map[string]any{"values": []int{3, -1, 1, 14}})
	req := httptest.NewRequest(http.MethodPost, "/v1/sequences/ingest", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp struct {
		Sequence struct {
			ID string `json:"id"`
		} `json:"sequence"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	if resp.Sequence.ID == "" {
		t.Fatalf("expected id")
	}

	getReq := httptest.NewRequest(http.MethodGet, "/v1/sequences/"+resp.Sequence.ID, nil)
	getRR := httptest.NewRecorder()
	r.ServeHTTP(getRR, getReq)

	if getRR.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", getRR.Code, getRR.Body.String())
	}
}
