package model

import (
	"time"

	"github.com/google/uuid"
)

type (
	// Pemesanan
	CreatePemesananRequest struct {
		Status string `db:"status" json:"status"`
		QTY    int    `db:"qty" json:"qty"`
	}

	UpdatePemesananRequest struct {
		Status string `db:"status" json:"status"`
	}

	PemesananResponse struct {
		ID            uuid.UUID  `db:"id" json:"id_pemesanan"`
		PemesananName string     `db:"name" json:"pemesanan_name"`
		Status        string     `db:"status" json:"status"`
		UserID        string     `db:"user_id" json:"user_id"`
		ProductID     string     `db:"product_id" json:"product_id"`
		QTY           int        `db:"qty" json:"qty"`
		CreatedAt     *time.Time `db:"created_at" json:"created_at"`
		UpdatedAt     *time.Time `db:"created_at" json:"updated_at"`
	}

	/*
		id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
		name VARCHAR NOT NULL,
		status pemesanan_status NOT NULL,
		user_id uuid NOT NULL,
		product_id uuid NOT NULL,
		qty INTEGER NOT NULL,
		created_at TIMESTAMPTZ DEFAULT now(),
		updated_at TIMESTAMPTZ
	*/
)
