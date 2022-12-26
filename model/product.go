package model

import (
	"time"

	"github.com/google/uuid"
)

type (
	// Product
	CreateProductRequest struct {
		ProductName     string    `db:"name" json:"product_name"`
		Stock           int       `db:"stok" json:"stok"`
		ProductCategory string    `db:"product_kategori" json:"product_kategori"`
		ProductImg      string    `db:"product_img" json:"product_img"`
		LapakID         uuid.UUID `db:"lapak_id" json:"-"`
	}

	// ProductUpdate
	UpdateProductRequest struct {
		ProductName     string `db:"name" json:"product_name"`
		Stock           int    `db:"stok" json:"stok"`
		ProductCategory string `db:"product_kategori" json:"product_kategori"`
		ProductImg      string `db:"product_img" json:"product_img" form:"image"`
	}

	// Get All Product Model
	GetAllProductRequest struct {
		ID              uuid.UUID  `db:"product_id" json:"product_id"`
		ProductName     string     `db:"product_name" json:"product_name"`
		Stock           int        `db:"stok" json:"stok"`
		ProductCategory string     `db:"product_kategori" json:"product_kategori"`
		ProductImg      string     `db:"product_img" json:"product_img"`
		CreatedAt       *time.Time `db:"created_at" json:"created_at"`
		UpdatedAt       *time.Time `db:"updated_at" json:"updated_at,omitempty"`
		LapakID         uuid.UUID  `db:"lapak_id" json:"lapak_id"`
		LapakName       string     `db:"lapak_name" json:"lapak_name"`
		LocationID      uuid.UUID  `db:"location_id" json:"location_id"`
		Daerah          string     `db:"daerah" json:"daerah"`
	}
)
