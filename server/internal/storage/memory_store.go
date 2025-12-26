package storage

import (
	"context"
	"sync"
	"time"

	"github.com/cgrente/sequence-insights-platform/server/internal/models"
	"github.com/google/uuid"
)

/*
MemoryStore is an in-memory implementation used in unit tests.

Do not use in production. It is intentionally simple and not optimized.
*/

type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]models.Sequence
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: make(map[string]models.Sequence)}
}

func (m *MemoryStore) Health(ctx context.Context) error { return nil }

func (m *MemoryStore) CreateSequence(ctx context.Context, seq models.Sequence) (models.Sequence, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	seq.ID = uuid.NewString()
	seq.CreatedAt = time.Now().UTC()

	m.data[seq.ID] = seq
	return seq, nil
}

func (m *MemoryStore) GetSequence(ctx context.Context, id string) (models.Sequence, bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	seq, ok := m.data[id]
	return seq, ok, nil
}

func (m *MemoryStore) MarkProcessed(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	seq, ok := m.data[id]
	if !ok {
		return nil
	}
	seq.Processed = true
	m.data[id] = seq
	return nil
}
