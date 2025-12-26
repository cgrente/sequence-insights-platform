package storage

import (
	"context"

	"github.com/cgrente/sequence-insights-platform/server/internal/models"
)

// Store abstracts persistence. Keep handlers/services unaware of the database.
type Store interface {
	Health(ctx context.Context) error

	CreateSequence(ctx context.Context, seq models.Sequence) (models.Sequence, error)
	GetSequence(ctx context.Context, id string) (models.Sequence, bool, error)

	MarkProcessed(ctx context.Context, id string) error
}
