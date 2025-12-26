package logging

import (
	"log/slog"
	"os"
)

/*
Package logging configures structured logging.

We use slog (standard library) to avoid extra dependencies and to keep logs machine-readable
(JSON), which is a strong default for production services.
*/

// NewJSONLogger returns a JSON slog.Logger writing to stdout.
func NewJSONLogger(level slog.Level) *slog.Logger {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	return slog.New(h)
}
