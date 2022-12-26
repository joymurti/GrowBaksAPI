package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"github.com/wiormiw/GrowBaks/config"
	"github.com/wiormiw/GrowBaks/model"
)

type (
	IPemesananRepository interface {
		CreatePemesanan(id uuid.UUID, product_id uuid.UUID, old_stok int, req model.CreatePemesananRequest, product *model.GetAllProductRequest) error
		GetAllPemesanan() ([]model.PemesananResponse, error)
		GetAllPemesananPribadi(id uuid.UUID) ([]model.PemesananResponse, error)
		GetPemesananByID(id uuid.UUID) (*model.PemesananResponse, error)
		UpdatePemesanan(id uuid.UUID, status string) error
		DeleteByID(pemesanan *model.PemesananResponse, product *model.GetAllProductRequest) error
	}

	PemesananRepository struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		DB      *pgx.Conn
	}
)

func (pemr *PemesananRepository) CreatePemesanan(id uuid.UUID, product_id uuid.UUID, old_stok int, req model.CreatePemesananRequest, product *model.GetAllProductRequest) error {

	q := `INSERT INTO pemesanan (
		"name",
		status,
		user_id,
		product_id,
		qty
	) 
	VALUES ($1,$2,$3,$4,$5) RETURNING id
	`

	tx, err := pemr.DB.Begin(pemr.Context)
	if err != nil {
		pemr.Logger.Error(fmt.Sprintf("Error 13: %s", err))
		return err
	}

	var pemesananID uuid.UUID
	sQTY := strconv.Itoa(req.QTY)
	pemesananName := "Pesanan:  " + sQTY + " " + product.ProductName

	err = tx.QueryRow(pemr.Context, q, pemesananName, req.Status, id, product_id, req.QTY).Scan(&pemesananID)
	if err != nil {
		err = tx.Rollback(pemr.Context)
		if err != nil {
			pemr.Logger.Error(fmt.Errorf("PemesananRepository.CreatePemesanan.QueryRow.Scan Rollback ERROR %v MSG %s", err, err.Error()))
			return err
		}

		pemr.Logger.Error(fmt.Errorf("PemesananRepository.CreatePemesanan.QueryRow Scan ERROR %v MSG %s", err, err.Error()))
		return err
	}

	q2 := `UPDATE "products" set stok = $1 WHERE id = $2`
	old_stok -= req.QTY
	_, err = tx.Exec(pemr.Context, q2, old_stok, product_id)
	if err != nil {
		err = tx.Rollback(pemr.Context)
		if err != nil {
			pemr.Logger.Error(fmt.Errorf("PemesananRepository.ChangeStock.Exec Rollback ERROR %v MSG %s", err, err.Error()))
			return err
		}

		pemr.Logger.Error(fmt.Errorf("PemesananRepository.ChangeStock.Exec ERROR %v MSG %s", err, err.Error()))
		return err
	}

	err = tx.Commit(pemr.Context)
	if err != nil {
		pemr.Logger.Error(fmt.Errorf("PemesananRepository.CreatePemesanan Commit ERROR %v MSG %s", err, err.Error()))
		return err
	}

	return nil
}

// GetAllPemesanan
func (pemr *PemesananRepository) GetAllPemesanan() ([]model.PemesananResponse, error) {
	q := `SELECT pem.id, 
		pem.name,
		pem.status, 
		pem.qty,
		pem.created_at,
		pem.updated_at,
		u.id as user_id,
		p.id as product_id
		FROM pemesanan pem
		LEFT JOIN users u ON pem.user_id = u.id
		LEFT JOIN products p on pem.product_id = p.id 
	`

	rows, err := pemr.DB.Query(pemr.Context, q)

	if err != nil {
		return nil, err
	}

	var listData []model.PemesananResponse

	for rows.Next() {
		pemesanan := &model.PemesananResponse{}
		err := rows.Scan(
			&pemesanan.ID,
			&pemesanan.PemesananName,
			&pemesanan.Status,
			&pemesanan.QTY,
			&pemesanan.CreatedAt,
			&pemesanan.UpdatedAt,
			&pemesanan.UserID,
			&pemesanan.ProductID)

		if err != nil {
			pemr.Logger.Error(fmt.Errorf("PemesananRepository.GetAllPemesanan rows.Next Scan ERROR %v MSG %s", err, err.Error()))
			return nil, err
		}

		listData = append(listData, *pemesanan)
	}

	return listData, nil
}

func (pemr *PemesananRepository) GetAllPemesananPribadi(id uuid.UUID) ([]model.PemesananResponse, error) {
	q := `SELECT pem.id, 
		pem.name,
		pem.status, 
		pem.qty,
		pem.created_at,
		pem.updated_at,
		u.id as user_id,
		p.id as product_id
		FROM pemesanan pem
		LEFT JOIN users u ON pem.user_id = u.id
		LEFT JOIN products p on pem.product_id = p.id 
		WHERE user_id = $1
	`

	rows, err := pemr.DB.Query(pemr.Context, q, id)

	if err != nil {
		return nil, err
	}

	var listData []model.PemesananResponse

	for rows.Next() {
		pemesanan := &model.PemesananResponse{}
		err := rows.Scan(
			&pemesanan.ID,
			&pemesanan.PemesananName,
			&pemesanan.Status,
			&pemesanan.QTY,
			&pemesanan.CreatedAt,
			&pemesanan.UpdatedAt,
			&pemesanan.UserID,
			&pemesanan.ProductID)

		if err != nil {
			pemr.Logger.Error(fmt.Errorf("PemesananRepository.GetAllPemesanan rows.Next Scan ERROR %v MSG %s", err, err.Error()))
			return nil, err
		}

		listData = append(listData, *pemesanan)
	}

	return listData, nil
}

// Edit Pemesanan
func (pemr *PemesananRepository) UpdatePemesanan(id uuid.UUID, status string) error {
	q := ` UPDATE pemesanan
		SET status = $1
		WHERE id = $2
	`
	_, err := pemr.DB.Exec(pemr.Context, q, status, id)
	if err != nil {
		pemr.Logger.Error(fmt.Errorf("PemesananRepository.UpdateByPemesanan ERROR : %v MSG : %s", err, err.Error()))
		return err
	}

	return nil
}

func (pemr *PemesananRepository) GetPemesananByID(id uuid.UUID) (*model.PemesananResponse, error) {
	var pemesanan model.PemesananResponse

	q := `SELECT pem.id, 
		pem.name,
		pem.status, 
		pem.qty,
		pem.created_at,
		pem.updated_at,
		u.id as user_id,
		p.id as product_id
		FROM pemesanan pem
		LEFT JOIN users u ON pem.user_id = u.id
		LEFT JOIN products p on pem.product_id = p.id 
		WHERE pem.id = $1
	`

	row := pemr.DB.QueryRow(pemr.Context, q, id)
	err := row.Scan(
		&pemesanan.ID,
		&pemesanan.PemesananName,
		&pemesanan.Status,
		&pemesanan.QTY,
		&pemesanan.CreatedAt,
		&pemesanan.UpdatedAt,
		&pemesanan.UserID,
		&pemesanan.ProductID,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			pemr.Logger.Info(fmt.Errorf("PemesananRepository.GetLapakByID INFO : %v MSG : %s", err, err.Error()))
		} else {
			pemr.Logger.Error(fmt.Errorf("PemesananRepository.GetLapakByID ERROR : %v MSG : %s", err, err.Error()))
		}

		return nil, err
	}

	return &pemesanan, nil
}

// Delete Pemesanan
func (pemr *PemesananRepository) DeleteByID(pemesanan *model.PemesananResponse, product *model.GetAllProductRequest) error {
	q := `DELETE FROM "pemesanan" WHERE id = $1`

	_, err := pemr.DB.Exec(pemr.Context, q, pemesanan.ID)
	if err != nil {
		pemr.Logger.Error(fmt.Errorf("PemesananRepository.DeleteByID Exec Pemesanan ERROR %v MSG %s", err, err.Error()))
		return err
	}

	newStock := product.Stock + pemesanan.QTY

	q2 := `UPDATE products SET stok = $1 WHERE id = $2`
	_, err = pemr.DB.Exec(pemr.Context, q2, newStock, product.ID)
	if err != nil {
		pemr.Logger.Error(fmt.Errorf("PemesananRepository.DeleteByID Exec Pemesanan ERROR %v MSG %s", err, err.Error()))
		return err
	}

	return nil
}
