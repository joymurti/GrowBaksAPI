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
	// IAuthRepository is an interface that has all the function to be implemented inside auth repository
	ILocationRepository interface {
		GetAllLocation() ([]model.Location, error)
		GetLocationByID(id uuid.UUID) (*model.Location, error)
	}

	// AuthRepository is an app auth struct that consists of all the dependencies needed for auth repository
	LocationRepository struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		DB      *pgx.Conn
	}
)

// GetAllLocation
func (lr *LocationRepository) GetAllLocation() ([]model.Location, error) {
	q := `SELECT 
		id,
		daerah
		FROM "lokasi"
	`

	rows, err := lr.DB.Query(lr.Context, q)
	if err != nil {
		lr.Logger.Error(fmt.Errorf("LocationRepository.GetAll Query ERROR %v MSG %s", err, err.Error()))
		return nil, err
	}

	var listData []model.Location
	for rows.Next() {
		data := &model.Location{}
		err := rows.Scan(&data.ID, &data.Daerah)
		if err != nil {
			lr.Logger.Error(fmt.Errorf("LocationRepository.GetAll rows.Next Scan ERROR %v MSG %s", err, err.Error()))
			return nil, err
		}

		listData = append(listData, *data)
	}

	return listData, nil
}

// GetLocationById
func (lr *LocationRepository) GetLocationByID(id uuid.UUID) (*model.Location, error) {
	var location model.Location

	q := `SELECT 
		id,
		daerah
		FROM "lokasi" 
		WHERE id = $1
	`

	row := lr.DB.QueryRow(lr.Context, q, id)
	err := row.Scan(&location.ID, &location.Daerah)
	if err != nil {
		if err == pgx.ErrNoRows {
			lr.Logger.Info(fmt.Errorf("LocationRepository.GetLocationByID INFO : %v MSG : %s", err, err.Error()))
		} else {
			lr.Logger.Error(fmt.Errorf("LocationRepository.GetLocationByID ERROR : %v MSG : %s", err, err.Error()))
		}

		return nil, err
	}

	return &location, nil
}
