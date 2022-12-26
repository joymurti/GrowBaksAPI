package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wiormiw/GrowBaks/config"
	"github.com/wiormiw/GrowBaks/model"
	"github.com/wiormiw/GrowBaks/repository"
)

type (
	IProductService interface {
		CreateProductSvc(id uuid.UUID, req model.CreateProductRequest) error
		GetAllProductSvc(id uuid.UUID, search string) ([]model.GetAllProductRequest, error)
		GetAllProductByLapakSvc(lapak_id uuid.UUID) ([]model.GetAllProductRequest, error)
		GetProductByIDSvc(id uuid.UUID) (*model.GetAllProductRequest, error)
		// GetLapakByIDSvc(id uuid.UUID) (*model.Lapak, error)
		// UpdateLapakSvc(id uuid.UUID, req model.LapakUpdate) error
		// UpdateLapakStatusSvc(id uuid.UUID, req model.LapakUpdateStatus) error
		// DeleteLapakSvc(id uuid.UUID) error
	}

	// ProductService is an app tag struct that consists of all the dependencies needed for lapak service
	ProductService struct {
		Context     context.Context
		Config      *config.Configuration
		Logger      *logrus.Logger
		ProductRepo repository.IProductRepository
		UserRepo    repository.IUserRepository
	}
)

// Create Product Service
func (ps *ProductService) CreateProductSvc(id uuid.UUID, req model.CreateProductRequest) error {
	req.LapakID = id

	err := validateCreateProductRequest(&req)
	if err != nil {
		ps.Logger.Error(fmt.Errorf("ERROR : %v MSG : %s", err, err.Error()))
		return err
	}

	err = ps.ProductRepo.Create(req)
	if err != nil {
		ps.Logger.Error(fmt.Errorf("ERROR : %v MSG : %s", err, err.Error()))
		return err
	}

	return nil
}

// GetAllProductSvc service layer for getting all lapak
func (ps *ProductService) GetAllProductSvc(id uuid.UUID, search string) ([]model.GetAllProductRequest, error) {

	userProfile, err := ps.UserRepo.GetProfileByID(id)
	if err != nil {
		return nil, err
	}

	userDaerah := userProfile.Daerah

	ps.Logger.Info(fmt.Sprintf("Value Test : %#+v UProfile", userProfile))

	data, err := ps.ProductRepo.GetAllProduct(search, userDaerah)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (ps *ProductService) GetAllProductByLapakSvc(lapak_id uuid.UUID) ([]model.GetAllProductRequest, error) {

	data, err := ps.ProductRepo.GetAllProductByLapak(lapak_id)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetAllLapakByIDSvc service layer for getting all lapak By Id
func (ps *ProductService) GetProductByIDSvc(id uuid.UUID) (*model.GetAllProductRequest, error) {
	data, err := ps.ProductRepo.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// validateCreateProductRequest responsible to validating create product request
func validateCreateProductRequest(req *model.CreateProductRequest) error {
	if req.ProductName == "" || req.ProductCategory == "" || req.ProductImg == "" || req.Stock < 1 {
		return model.ErrInvalidRequest
	}

	if len(req.ProductName) < 5 {
		return model.ErrInvalidRequest
	}

	if req.ProductCategory != "makanan" && req.ProductCategory != "minuman" {
		return model.ErrInvalidRequest
	}

	return nil
}
