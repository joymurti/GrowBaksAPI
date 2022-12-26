package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"github.com/wiormiw/GrowBaks/config"
	"github.com/wiormiw/GrowBaks/model"
	"github.com/wiormiw/GrowBaks/repository"
)

type (
	ILapakService interface {
		GetAllLapakSvc(search string) ([]model.Lapak, error)
		GetAllLapakByLocationSvc(id uuid.UUID, search string) ([]model.Lapak, error)
		GetLapakByIDSvc(id uuid.UUID) (*model.Lapak, error)
		UpdateLapakSvc(id uuid.UUID, req model.LapakUpdate) error
		UpdateLapakStatusSvc(id uuid.UUID, req model.LapakUpdateStatus) error
		DeleteLapakSvc(id uuid.UUID) error
	}

	// LapakService is an app tag struct that consists of all the dependencies needed for lapak service
	LapakService struct {
		Context   context.Context
		Config    *config.Configuration
		Logger    *logrus.Logger
		LapakRepo repository.ILapakRepository
		UserRepo  repository.IUserRepository
	}
)

// GetAllLapakSvc service layer for getting all lapak
func (laps *LapakService) GetAllLapakSvc(search string) ([]model.Lapak, error) {
	data, err := laps.LapakRepo.GetAllLapak(search)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetAllLapakSvc service layer for getting all lapak
func (laps *LapakService) GetAllLapakByLocationSvc(id uuid.UUID, search string) ([]model.Lapak, error) {
	userProfile, err := laps.UserRepo.GetProfileByID(id)
	if err != nil {
		return nil, err
	}

	userDaerah := userProfile.Daerah

	data, err := laps.LapakRepo.GetAllLapakByLocation(search, userDaerah)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetAllLapakByIDSvc service layer for getting all lapak By Id
func (laps *LapakService) GetLapakByIDSvc(id uuid.UUID) (*model.Lapak, error) {
	data, err := laps.LapakRepo.GetLapakByID(id)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// UpdateTagSvc service layer for updating a tag by id
func (laps *LapakService) UpdateLapakSvc(id uuid.UUID, req model.LapakUpdate) error {
	_, err := laps.LapakRepo.GetLapakByID(id)

	if err != nil {
		if err == pgx.ErrNoRows {
			return model.ErrLapakNotFound
		}

		return err
	}

	err = validateUpdateLapakRequest(&req)
	if err != nil {
		return err
	}

	err = laps.LapakRepo.UpdateByID(id, req.Name, req.Status)
	if err != nil {
		return err
	}

	return nil
}

func (laps *LapakService) UpdateLapakStatusSvc(id uuid.UUID, req model.LapakUpdateStatus) error {
	_, err := laps.LapakRepo.GetLapakByID(id)

	if err != nil {
		if err == pgx.ErrNoRows {
			return model.ErrLapakNotFound
		}

		return err
	}

	err = validateUpdateLapakStatusRequest(&req)
	if err != nil {
		return err
	}

	err = laps.LapakRepo.UpdateStatusByID(id, req.Status)
	if err != nil {
		return err
	}

	return nil
}

// DeleteLapakSvc service layer for deleting a lapak by id
func (laps *LapakService) DeleteLapakSvc(id uuid.UUID) error {
	_, err := laps.LapakRepo.GetLapakByID(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.ErrLapakNotFound
		}

		return err
	}

	err = laps.LapakRepo.DeleteByID(id)
	if err != nil {
		return err
	}

	return nil
}

// validateUpdate
func validateUpdateLapakRequest(lapak *model.LapakUpdate) error {
	if lapak.Name == "" || lapak.Status == "" {
		return model.ErrInvalidRequest
	}

	if lapak.Name != "" {
		if len(lapak.Name) < 5 {
			return model.ErrInvalidRequest
		}
	}

	if lapak.Status != "" {
		if lapak.Status != "open" && lapak.Status != "closed" {
			return model.ErrInvalidRequest
		}
	}

	return nil
}

func validateUpdateLapakStatusRequest(lapak *model.LapakUpdateStatus) error {
	if lapak.Status == "" {
		return model.ErrInvalidRequest
	}
	if lapak.Status != "" {
		if lapak.Status != "open" && lapak.Status != "closed" {
			return model.ErrInvalidRequest
		}
	}

	return nil
}
