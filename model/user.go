package model

import (
	"time"

	"github.com/google/uuid"
)

type (
	// UpdateUserRequest consist data of updating a user
	UpdateUserRequest struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	UpdateUserProfileRequest struct {
		FullName     string `json:"full_name,omitempty"`
		TanggalLahir string `json:"tanggal_lahir,omitempty"`
		Gender       string `json:"gender,omitempty"`
		Email        string `json:"email,omitempty"`
		Telepon      string `json:"telepon,omitempty"`
	}

	// ViewUserResponse consist data of user
	ViewUserResponse struct {
		ID        uuid.UUID  `db:"id" json:"id,omitempty"`
		FullName  string     `db:"full_name" json:"full_name,omitempty"`
		Email     string     `db:"email" json:"email,omitempty"`
		CreatedAt *time.Time `db:"created_at" json:"created_at,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
	}

	ViewUserProfileResponse struct {
		UserID       uuid.UUID  `db:"user_id" json:"user_id,omitempty"`
		FullName     string     `db:"full_name" json:"full_name,omitempty"`
		Email        string     `db:"email" json:"email,omitempty"`
		CreatedAt    *time.Time `db:"created_at" json:"created_at,omitempty"`
		UpdatedAt    *time.Time `db:"updated_at" json:"updated_at,omitempty"`
		Telepon      string     `db:"telepon" json:"telepon,omitempty"`
		Gender       string     `db:"gender" json:"gender,omitempty"`
		TanggalLahir *time.Time `db:"tanggal_lahir" json:"tanggal_lahir,omitempty"`
		LocationID   uuid.UUID  `db:"location_id" json:"location_id,omitempty"`
		Daerah       string     `db:"daerah" json:"daerah,omitempty"`
	}
)
