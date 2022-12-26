package model

import "github.com/google/uuid"

type (
	// Lapak
	Lapak struct {
		UserID     uuid.UUID `db:"user_id" json:"user_id"`
		FullName   string    `db:"user_id" json:"full_name"`
		LapakID    uuid.UUID `db:"lapak_id" json:"lapak_id"`
		LapakName  string    `db:"lapak_name" json:"lapak_name"`
		Status     string    `db:"status" json:"status"`
		LocationID uuid.UUID `db:"location_id" json:"location_id"`
		Daerah     string    `db:"daerah" json:"daerah"`
	}

	// LapakUpdate
	LapakUpdate struct {
		Name   string `db:"name" json:"name"`
		Status string `db:"status" json:"status"`
	}

	// LapakStatus
	LapakUpdateStatus struct {
		Status string `db:"status" json:"status"`
	}
)
