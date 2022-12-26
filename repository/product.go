package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"github.com/wiormiw/GrowBaks/config"
	"github.com/wiormiw/GrowBaks/model"
)

type (
	IProductRepository interface {
		Create(req model.CreateProductRequest) error
		GetAllProduct(search string, daerah string) ([]model.GetAllProductRequest, error)
		GetAllProductByLapak(lapak_id uuid.UUID) ([]model.GetAllProductRequest, error)
		GetProductByID(id uuid.UUID) (*model.GetAllProductRequest, error)
		// UpdateByID(id uuid.UUID, nama_produk string, kategori_produk string) error
		// DeleteByID(id uuid.UUID) error
	}

	ProductRepository struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		DB      *pgx.Conn
	}
)

// Create Product
func (pr *ProductRepository) Create(req model.CreateProductRequest) error {
	q := `INSERT INTO products (
		"name",
		stok,
		product_kategori,
		product_img,
		lapak_id
	) 
	VALUES ($1,$2,$3,$4,$5)
	`

	_, err := pr.DB.Exec(pr.Context, q, req.ProductName, req.Stock, req.ProductCategory, req.ProductImg, req.LapakID)
	if err != nil {
		pr.Logger.Error(fmt.Errorf("ProductRepository.Create Exec ERROR %v MSG %s", err, err.Error()))
		return err
	}

	return nil
}

// GetAllProduct
func (pr *ProductRepository) GetAllProduct(search string, daerah string) ([]model.GetAllProductRequest, error) {

	q := `SELECT p.id AS product_id,
		p.name as product_name,
		p.stok,
		p.product_kategori,
		p.product_img,
		p.created_at,
		p.updated_at,
		l.id as lapak_id,
		l.name as lapak_name,
		loc.id as location_id,
		loc.daerah as daerah
		FROM "products" p
		LEFT JOIN lapak l ON p.lapak_id = l.id
		LEFT JOIN lokasi loc on l.location_id = loc.id
	`

	criteria := ""
	criteria = ""

	if len(search) > 0 {
		criteria += " p.name ILIKE '%" + search + "%'"
	}

	cr := ""
	filterDaerah := ""
	if len(criteria) > 0 {
		cr = " WHERE " + criteria
	}

	if len(criteria) > 0 && daerah != "" {
		filterDaerah = " AND daerah = '" + daerah + "'"
	} else if len(criteria) < 1 && daerah != "" {
		filterDaerah = " WHERE daerah = '" + daerah + "'"
	}

	if cr != "" {
		q += cr
	}

	if filterDaerah != "" {
		q += filterDaerah
	}

	pr.Logger.Info(fmt.Sprintf("Query : %s", q))

	rows, err := pr.DB.Query(pr.Context, q)

	if err != nil {
		return nil, err
	}

	var listData []model.GetAllProductRequest
	for rows.Next() {
		data := &model.GetAllProductRequest{}
		err := rows.Scan(&data.ID,
			&data.ProductName,
			&data.Stock,
			&data.ProductCategory,
			&data.ProductImg,
			&data.CreatedAt,
			&data.UpdatedAt,
			&data.LapakID,
			&data.LapakName,
			&data.LocationID,
			&data.Daerah)
		if err != nil {
			pr.Logger.Error(fmt.Errorf("ProductRepository.GetAllProduct rows.Next Scan ERROR %v MSG %s", err, err.Error()))
			return nil, err
		}

		listData = append(listData, *data)
	}

	return listData, nil
}

func (pr *ProductRepository) GetAllProductByLapak(lapak_id uuid.UUID) ([]model.GetAllProductRequest, error) {
	q := `SELECT p.id AS product_id,
		p.name as product_name,
		p.stok,
		p.product_kategori,
		p.product_img,
		p.created_at,
		p.updated_at,
		l.id as lapak_id,
		l.name as lapak_name,
		loc.id as location_id,
		loc.daerah as daerah
		FROM "products" p
		LEFT JOIN lapak l ON p.lapak_id = l.id
		LEFT JOIN lokasi loc on l.location_id = loc.id
		WHERE lapak_id = $1
	`

	pr.Logger.Info(fmt.Sprintf("Query : %s", q))

	rows, err := pr.DB.Query(pr.Context, q, lapak_id)

	if err != nil {
		return nil, err
	}

	var listData []model.GetAllProductRequest
	for rows.Next() {
		data := &model.GetAllProductRequest{}
		err := rows.Scan(&data.ID,
			&data.ProductName,
			&data.Stock,
			&data.ProductCategory,
			&data.ProductImg,
			&data.CreatedAt,
			&data.UpdatedAt,
			&data.LapakID,
			&data.LapakName,
			&data.LocationID,
			&data.Daerah)
		if err != nil {
			pr.Logger.Error(fmt.Errorf("ProductRepository.GetAllProductByLapak rows.Next Scan ERROR %v MSG %s", err, err.Error()))
			return nil, err
		}

		listData = append(listData, *data)
	}

	return listData, nil
}

// GetLapakById
func (pr *ProductRepository) GetProductByID(id uuid.UUID) (*model.GetAllProductRequest, error) {
	var product model.GetAllProductRequest

	q := `SELECT p.id AS product_id,
		p.name as product_name,
		p.stok,
		p.product_kategori,
		p.product_img,
		p.created_at,
		p.updated_at,
		l.id as lapak_id,
		l.name as lapak_name,
		loc.id as location_id,
		loc.daerah as daerah
		FROM "products" p
		LEFT JOIN lapak l ON p.lapak_id = l.id
		LEFT JOIN lokasi loc on l.location_id = loc.id
		WHERE p.id = $1
	`

	row := pr.DB.QueryRow(pr.Context, q, id)
	err := row.Scan(&product.ID,
		&product.ProductName,
		&product.Stock,
		&product.ProductCategory,
		&product.ProductImg,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.LapakID,
		&product.LapakName,
		&product.LocationID,
		&product.Daerah)
	if err != nil {
		if err == pgx.ErrNoRows {
			pr.Logger.Info(fmt.Errorf("LapakRepository.GetLapakByID INFO : %v MSG : %s", err, err.Error()))
		} else {
			pr.Logger.Error(fmt.Errorf("LapakRepository.GetLapakByID ERROR : %v MSG : %s", err, err.Error()))
		}

		return nil, err
	}

	return &product, nil
}

// // Edit Lapak
// func (lapr *LapakRepository) UpdateByID(id uuid.UUID, nama_lapak string, status string) error {
// 	q := ` UPDATE lapak
// 		SET name = $1,
// 		status = $2
// 	    WHERE id = $3
// 	`
// 	_, err := lapr.DB.Exec(lapr.Context, q, nama_lapak, status, id)
// 	if err != nil {
// 		lapr.Logger.Error(fmt.Errorf("LapakRepository.UpdateByID Lapak ERROR : %v MSG : %s", err, err.Error()))
// 		return err
// 	}

// 	return nil
// }

// // Edit Status Lapak
// func (lapr *LapakRepository) UpdateStatusByID(id uuid.UUID, status string) error {
// 	q := ` UPDATE lapak
// 		SET status = $1
// 	    WHERE id = $2
// 	`
// 	_, err := lapr.DB.Exec(lapr.Context, q, status, id)
// 	if err != nil {
// 		lapr.Logger.Error(fmt.Errorf("LapakRepository.UpdateByStatusID Lapak ERROR : %v MSG : %s", err, err.Error()))
// 		return err
// 	}

// 	return nil
// }

// // Delete Lapak By ID
// func (lapr *LapakRepository) DeleteByID(id uuid.UUID) error {
// 	q := `DELETE FROM lapak WHERE id = $1`

// 	_, err := lapr.DB.Exec(lapr.Context, q, id)
// 	if err != nil {
// 		lapr.Logger.Error(fmt.Errorf("LapakRepository.DeleteByID Exec ERROR %v MSG %s", err, err.Error()))
// 		return err
// 	}

// 	return nil
// }
