package db

import (
	"context"
	"database/sql"
	"time"
)

/*
Package db contains database wiring only (connection pools, health checks).

We register pgx as the database/sql driver to keep our application interface
portable and test-friendly.
*/

// Open opens a Postgres *sql.DB using DATABASE_URL.
func Open(dbURL string) (*sql.DB, error) {
	// pgx uses the connection string directly.
	return sql.Open("pgx", dbURL)
}

// Ping checks DB connectivity with a timeout.
func Ping(ctx context.Context, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return db.PingContext(ctx)
}
