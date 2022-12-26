package model

import "github.com/google/uuid"

type (
	// Location
	Location struct {
		ID     uuid.UUID `db:"id" json:"-"`
		Daerah string    `db:"daerah" json:"daerah"`
	}
)
