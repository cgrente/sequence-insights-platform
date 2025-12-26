package jobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/cgrente/sequence-insights-platform/server/internal/storage"
)

/*
Package jobs runs lightweight background work in-process.

This keeps the demo self-contained (no Kafka/Redis queue) while still showing:
- decoupled async execution
- back-pressure via buffered channels
- graceful shutdown
*/

type Job struct {
	SequenceID string
}

type Queue struct {
	log   *slog.Logger
	store storage.Store

	jobs chan Job
}

func NewQueue(log *slog.Logger, store storage.Store, size int) *Queue {
	return &Queue{
		log:   log,
		store: store,
		jobs:  make(chan Job, size),
	}
}

func (q *Queue) Enqueue(job Job) bool {
	select {
	case q.jobs <- job:
		return true
	default:
		return false
	}
}

func (q *Queue) RunWorkers(ctx context.Context, workers int) {
	if workers < 1 {
		workers = 1
	}
	for i := 0; i < workers; i++ {
		go q.worker(ctx, i+1)
	}
}

func (q *Queue) worker(ctx context.Context, n int) {
	q.log.Info("worker started", "worker", n)
	for {
		select {
		case <-ctx.Done():
			q.log.Info("worker stopping", "worker", n)
			return
		case job := <-q.jobs:
			// Simulate post-processing: in a real system this could be heavier analytics.
			time.Sleep(50 * time.Millisecond)

			if err := q.store.MarkProcessed(ctx, job.SequenceID); err != nil {
				q.log.Error("failed to mark processed", "worker", n, "sequence_id", job.SequenceID, "err", err)
				continue
			}
			q.log.Info("sequence processed", "worker", n, "sequence_id", job.SequenceID)
		}
	}
}
