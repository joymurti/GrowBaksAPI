package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"github.com/wiormiw/GrowBaks/config"
	"github.com/wiormiw/GrowBaks/model"
	"github.com/wiormiw/GrowBaks/repository"
)

type (
	IPemesananService interface {
		CreatePemesananSvc(id uuid.UUID, product_id uuid.UUID, req model.CreatePemesananRequest) error
		GetAllPemesananSvc() ([]model.PemesananResponse, error)
		GetAllPemesananPribadiSvc(id uuid.UUID) ([]model.PemesananResponse, error)
		GetPemesananByIDSvc(id uuid.UUID) (*model.PemesananResponse, error)
		UpdatePemesananStatusSvc(id uuid.UUID, req model.UpdatePemesananRequest) error
		DeletePemesananSvc(id uuid.UUID) error
	}

	// PemesananService is an app tag struct that consists of all the dependencies needed for pemesanan service
	PemesananService struct {
		Context       context.Context
		Config        *config.Configuration
		Logger        *logrus.Logger
		PemesananRepo repository.IPemesananRepository
		UserRepo      repository.IUserRepository
		ProductRepo   repository.IProductRepository
	}
)

// Create Product Service
func (pems *PemesananService) CreatePemesananSvc(id uuid.UUID, product_id uuid.UUID, req model.CreatePemesananRequest) error {

	product, err := pems.ProductRepo.GetProductByID(product_id)
	if err != nil {
		pems.Logger.Error(fmt.Errorf("ERROR : %v MSG : %s", err, err.Error()))
		return err
	}

	if product.Stock >= req.QTY {
		err = validateCreatePemesananRequest(&req)
		if err != nil {
			pems.Logger.Error(fmt.Errorf("ERROR : %v MSG : %s", err, err.Error()))
			return err
		}

		err = pems.PemesananRepo.CreatePemesanan(id, product_id, product.Stock, req, product)
		if err != nil {
			pems.Logger.Error(fmt.Errorf("ERROR : %v MSG : %s", err, err.Error()))
			return err
		}
	} else {
		return model.ErrInvalidRequest
	}

	return nil
}

func (pems *PemesananService) GetAllPemesananSvc() ([]model.PemesananResponse, error) {
	data, err := pems.PemesananRepo.GetAllPemesanan()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (pems *PemesananService) GetAllPemesananPribadiSvc(id uuid.UUID) ([]model.PemesananResponse, error) {
	data, err := pems.PemesananRepo.GetAllPemesananPribadi(id)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (pems *PemesananService) GetPemesananByIDSvc(id uuid.UUID) (*model.PemesananResponse, error) {
	data, err := pems.PemesananRepo.GetPemesananByID(id)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (pems *PemesananService) UpdatePemesananStatusSvc(id uuid.UUID, req model.UpdatePemesananRequest) error {
	_, err := pems.PemesananRepo.GetPemesananByID(id)

	if err != nil {
		if err == pgx.ErrNoRows {
			return model.ErrLapakNotFound
		}

		return err
	}

	err = validateUpdatePemesananStatus(&req)
	if err != nil {
		return err
	}

	err = pems.PemesananRepo.UpdatePemesanan(id, req.Status)
	if err != nil {
		return err
	}

	return nil
}

// Delete Pemesanan
func (pems *PemesananService) DeletePemesananSvc(id uuid.UUID) error {
	pemesanan, err := pems.PemesananRepo.GetPemesananByID(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.ErrPemesananNotFound
		}

		return err
	}

	productID, err := uuid.Parse(pemesanan.ProductID)
	if err != nil {
		pems.Logger.Error(fmt.Sprintf("Error 9: %s", err))
		return err
	}

	product, err := pems.ProductRepo.GetProductByID(productID)

	err = pems.PemesananRepo.DeleteByID(pemesanan, product)
	if err != nil {
		return err
	}

	return nil
}

// validateCreateProductRequest responsible to validating create product request
func validateCreatePemesananRequest(req *model.CreatePemesananRequest) error {
	if req.Status == "" || req.QTY < 1 {
		return model.ErrInvalidRequest
	}

	return nil
}

func validateUpdatePemesananStatus(req *model.UpdatePemesananRequest) error {
	if req.Status == "" {
		return model.ErrInvalidRequest
	}

	return nil
}
