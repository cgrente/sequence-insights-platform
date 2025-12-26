package models

import "time"

// Sequence represents an ingested numeric sequence and computed metrics.
type Sequence struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`

	Values []int64 `json:"values"`

	// Metrics (computed)
	Count                 int   `json:"count"`
	SumFourthPowersNonPos int64 `json:"sum_fourth_powers_non_positive"`
	Min                   int64 `json:"min"`
	Max                   int64 `json:"max"`

	Processed bool `json:"processed"`
}
