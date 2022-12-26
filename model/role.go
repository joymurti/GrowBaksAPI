package model

import "github.com/google/uuid"

type (
	// Role
	Role struct {
		ID   uuid.UUID `db:"id" json:"-"`
		Name string    `db:"name" json:"name"`
	}
)
