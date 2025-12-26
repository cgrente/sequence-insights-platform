package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/cgrente/sequence-insights-platform/server/internal/models"
)

/*
PostgresStore implements Store backed by PostgreSQL.

Schema is defined in server/migrations/001_init.sql.

Notes:
- Use context on all DB calls.
- Keep SQL explicit; avoid ORM for small services.
*/

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (p *PostgresStore) Health(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

func (p *PostgresStore) CreateSequence(ctx context.Context, seq models.Sequence) (models.Sequence, error) {
	valuesJSON, err := json.Marshal(seq.Values)
	if err != nil {
		return models.Sequence{}, err
	}

	row := p.db.QueryRowContext(ctx, `
		INSERT INTO sequences (values_json, count, sum_fourth_powers_non_positive, min_value, max_value, processed)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`, valuesJSON, seq.Count, seq.SumFourthPowersNonPos, seq.Min, seq.Max, seq.Processed)

	if err := row.Scan(&seq.ID, &seq.CreatedAt); err != nil {
		return models.Sequence{}, err
	}
	return seq, nil
}

func (p *PostgresStore) GetSequence(ctx context.Context, id string) (models.Sequence, bool, error) {
	var seq models.Sequence
	var valuesJSON []byte
	var createdAt time.Time

	err := p.db.QueryRowContext(ctx, `
		SELECT id, created_at, values_json, count, sum_fourth_powers_non_positive, min_value, max_value, processed
		FROM sequences
		WHERE id = $1
	`, id).Scan(&seq.ID, &createdAt, &valuesJSON, &seq.Count, &seq.SumFourthPowersNonPos, &seq.Min, &seq.Max, &seq.Processed)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Sequence{}, false, nil
	}
	if err != nil {
		return models.Sequence{}, false, err
	}

	seq.CreatedAt = createdAt.UTC()
	if err := json.Unmarshal(valuesJSON, &seq.Values); err != nil {
		return models.Sequence{}, false, err
	}

	return seq, true, nil
}

func (p *PostgresStore) MarkProcessed(ctx context.Context, id string) error {
	_, err := p.db.ExecContext(ctx, `UPDATE sequences SET processed = TRUE WHERE id = $1`, id)
	return err
}
