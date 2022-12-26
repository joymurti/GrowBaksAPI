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
	ILapakRepository interface {
		GetAllLapak(search string) ([]model.Lapak, error)
		GetAllLapakByLocation(search string, daerah string) ([]model.Lapak, error)
		GetLapakByID(id uuid.UUID) (*model.Lapak, error)
		UpdateByID(id uuid.UUID, nama_lapak string, status string) error
		UpdateStatusByID(id uuid.UUID, status string) error
		DeleteByID(id uuid.UUID) error
	}

	LapakRepository struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		DB      *pgx.Conn
	}
)

// GetAllLapak
func (lapr *LapakRepository) GetAllLapak(search string) ([]model.Lapak, error) {
	q := `SELECT u.id AS user_id, 
		u.full_name, 
		l.id as lapak_id,
		l.name as lapak_name,
		l.status, 
		loc.id as lokasi_id,
		loc.daerah 
		FROM "lapak" l 
		LEFT JOIN users u ON l.user_id = u.id 
		LEFT JOIN lokasi loc ON loc.id = l.location_id
	`

	criteria := ""
	criteria = ""

	if len(search) > 0 {
		criteria += " l.name ILIKE '%" + search + "%'"
	}

	cr := ""
	if len(criteria) > 0 {
		cr = " WHERE " + criteria
	}

	if cr != "" {
		q += cr
	}

	lapr.Logger.Info(fmt.Sprintf("Query : %s", q))

	rows, err := lapr.DB.Query(lapr.Context, q)

	if err != nil {
		return nil, err
	}

	var listData []model.Lapak
	for rows.Next() {
		data := &model.Lapak{}
		err := rows.Scan(&data.UserID,
			&data.FullName,
			&data.LapakID,
			&data.LapakName,
			&data.Status,
			&data.LocationID,
			&data.Daerah)
		if err != nil {
			lapr.Logger.Error(fmt.Errorf("LapakRepository.GetAllLapak rows.Next Scan ERROR %v MSG %s", err, err.Error()))
			return nil, err
		}

		listData = append(listData, *data)
	}

	return listData, nil
}

func (lapr *LapakRepository) GetAllLapakByLocation(search string, daerah string) ([]model.Lapak, error) {
	q := `SELECT u.id AS user_id, 
		u.full_name, 
		l.id as lapak_id,
		l.name as lapak_name,
		l.status, 
		loc.id as lokasi_id,
		loc.daerah 
		FROM "lapak" l 
		LEFT JOIN users u ON l.user_id = u.id 
		LEFT JOIN lokasi loc ON loc.id = l.location_id
	`

	criteria := ""
	criteria = ""

	if len(search) > 0 {
		criteria += " l.name ILIKE '%" + search + "%'"
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

	lapr.Logger.Info(fmt.Sprintf("Query : %s", q))

	rows, err := lapr.DB.Query(lapr.Context, q)

	if err != nil {
		return nil, err
	}

	var listData []model.Lapak
	for rows.Next() {
		data := &model.Lapak{}
		err := rows.Scan(&data.UserID,
			&data.FullName,
			&data.LapakID,
			&data.LapakName,
			&data.Status,
			&data.LocationID,
			&data.Daerah)
		if err != nil {
			lapr.Logger.Error(fmt.Errorf("LapakRepository.GetAllLapak rows.Next Scan ERROR %v MSG %s", err, err.Error()))
			return nil, err
		}

		listData = append(listData, *data)
	}

	return listData, nil
}

// GetLapakById
func (lapr *LapakRepository) GetLapakByID(id uuid.UUID) (*model.Lapak, error) {
	var lapak model.Lapak

	q := `SELECT u.id AS user_id, 
		u.full_name, 
		l.id as lapak_id,
		l.name as lapak_name,
		l.status, 
		loc.id as lokasi_id,
		loc.daerah 
		FROM "lapak" l 
		LEFT JOIN users u ON l.user_id = u.id 
		LEFT JOIN lokasi loc ON loc.id = l.location_id
		WHERE l.id = $1
	`

	row := lapr.DB.QueryRow(lapr.Context, q, id)
	err := row.Scan(&lapak.UserID,
		&lapak.FullName,
		&lapak.LapakID,
		&lapak.LapakName,
		&lapak.Status,
		&lapak.LocationID,
		&lapak.Daerah)
	if err != nil {
		if err == pgx.ErrNoRows {
			lapr.Logger.Info(fmt.Errorf("LapakRepository.GetLapakByID INFO : %v MSG : %s", err, err.Error()))
		} else {
			lapr.Logger.Error(fmt.Errorf("LapakRepository.GetLapakByID ERROR : %v MSG : %s", err, err.Error()))
		}

		return nil, err
	}

	return &lapak, nil
}

// Edit Lapak
func (lapr *LapakRepository) UpdateByID(id uuid.UUID, nama_lapak string, status string) error {
	q := ` UPDATE lapak
		SET name = $1,
		status = $2
	    WHERE id = $3
	`
	_, err := lapr.DB.Exec(lapr.Context, q, nama_lapak, status, id)
	if err != nil {
		lapr.Logger.Error(fmt.Errorf("LapakRepository.UpdateByID Lapak ERROR : %v MSG : %s", err, err.Error()))
		return err
	}

	return nil
}

// Edit Status Lapak
func (lapr *LapakRepository) UpdateStatusByID(id uuid.UUID, status string) error {
	q := ` UPDATE lapak
		SET status = $1
	    WHERE id = $2
	`
	_, err := lapr.DB.Exec(lapr.Context, q, status, id)
	if err != nil {
		lapr.Logger.Error(fmt.Errorf("LapakRepository.UpdateByStatusID Lapak ERROR : %v MSG : %s", err, err.Error()))
		return err
	}

	return nil
}

// Delete Lapak By ID
func (lapr *LapakRepository) DeleteByID(id uuid.UUID) error {
	q := `DELETE FROM lapak WHERE id = $1`

	_, err := lapr.DB.Exec(lapr.Context, q, id)
	if err != nil {
		lapr.Logger.Error(fmt.Errorf("LapakRepository.DeleteByID Exec ERROR %v MSG %s", err, err.Error()))
		return err
	}

	q2 := `DELETE FROM products WHERE lapak_id = $1`

	_, err = lapr.DB.Exec(lapr.Context, q2, id)
	if err != nil {
		lapr.Logger.Error(fmt.Errorf("LapakRepository.DeleteByID Exec ERROR %v MSG %s", err, err.Error()))
		return err
	}

	return nil
}
