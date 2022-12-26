package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"github.com/wiormiw/GrowBaks/config"
	"github.com/wiormiw/GrowBaks/helper"
	"github.com/wiormiw/GrowBaks/model"
	"github.com/wiormiw/GrowBaks/repository"
)

type (
	// IUserService is an interface that has all the function to be implemented inside user service
	IUserService interface {
		GetAllUserSvc(search string) ([]model.ViewUserResponse, error)
		GetUserByIDSvc(id uuid.UUID) (*model.ViewUserResponse, error)
		GetUserProfileByIDSvc(id uuid.UUID) (*model.ViewUserProfileResponse, error)
		UpdateUserByIDSvc(id uuid.UUID, req model.UpdateUserRequest) error
		UpdateUserProfileByIDSvc(id uuid.UUID, req model.UpdateUserProfileRequest) error
		DeleteUserByIDSvc(id uuid.UUID) error
	}

	// UserService is an app user check struct that consists of all the dependencies needed for user service
	UserService struct {
		Context  context.Context
		Config   *config.Configuration
		Logger   *logrus.Logger
		UserRepo repository.IUserRepository
	}
)

// GetAllUserSvc service layer for getting all user
func (us *UserService) GetAllUserSvc(search string) ([]model.ViewUserResponse, error) {
	data, err := us.UserRepo.GetAll(search)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetUserByIDSvc service layer for get a user by id
func (us *UserService) GetUserByIDSvc(id uuid.UUID) (*model.ViewUserResponse, error) {
	data, err := us.UserRepo.GetByID(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, model.ErrUserNotFound
		}

		return nil, err
	}

	return data, nil
}

// GetUserProfileByIDSvc service layer for get a user profile by id
func (us *UserService) GetUserProfileByIDSvc(id uuid.UUID) (*model.ViewUserProfileResponse, error) {
	data, err := us.UserRepo.GetProfileByID(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, model.ErrUserNotFound
		}

		return nil, err
	}

	return data, nil
}

// UpdateUserByIDSvc service layer for update user by id
func (us *UserService) UpdateUserByIDSvc(id uuid.UUID, req model.UpdateUserRequest) error {
	_, err := us.UserRepo.GetByID(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.ErrUserNotFound
		}

		return err
	}

	_, err = us.UserRepo.GetByEmail(req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			err := validateUpdateUserRequest(&req)
			if err != nil {
				return err
			}

			req.Password, err = helper.HashPassword(req.Password)
			if err != nil {
				return err
			}

			err = us.UserRepo.UpdateByID(id, req.FullName, req.Email, req.Password)
			if err != nil {
				return err
			}

			return nil

		}

		return err
	}

	return model.ErrEmailExisted
}

func (us *UserService) UpdateUserProfileByIDSvc(id uuid.UUID, req model.UpdateUserProfileRequest) error {
	_, err := us.UserRepo.GetByID(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.ErrUserNotFound
		}

		return err
	}

	_, err = us.UserRepo.GetByEmail(req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = validateUpdateProfileUserRequest(&req)
			if err != nil {
				return err
			}

			err = us.UserRepo.UpdateProfileByID(id, req.Email, req.FullName, req.TanggalLahir, req.Gender, req.Telepon)
			if err != nil {
				return err
			}

			return nil

		}

		return err
	}

	return model.ErrEmailExisted
}

// DeleteUserByIDSvc service layer for delete user by id
func (us *UserService) DeleteUserByIDSvc(id uuid.UUID) error {
	_, err := us.UserRepo.GetByID(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.ErrUserNotFound
		}

		return err
	}

	err = us.UserRepo.DeleteByID(id)
	if err != nil {
		return err
	}

	return nil
}

// validateUpdateUserRequest responsible to validating update user
func validateUpdateUserRequest(req *model.UpdateUserRequest) error {
	if len(req.FullName) < 5 {
		return model.ErrInvalidRequest
	}

	if !model.IsAllowedEmailInput.MatchString(req.Email) {
		return model.ErrInvalidRequest
	}

	if len(req.Password) < 8 {
		return model.ErrInvalidRequest
	}

	return nil
}

func validateUpdateProfileUserRequest(req *model.UpdateUserProfileRequest) error {
	if len(req.FullName) < 5 {
		return model.ErrInvalidRequest
	}

	if !model.IsAllowedDateInput.MatchString(req.TanggalLahir) {
		return model.ErrInvalidRequest
	}

	if req.Gender != "laki-laki" && req.Gender != "perempuan" {
		return model.ErrInvalidRequest
	}

	if !model.IsAllowedTeleponInput.MatchString(req.Telepon) {
		return model.ErrInvalidRequest
	}

	if !model.IsAllowedEmailInput.MatchString(req.Email) {
		return model.ErrInvalidRequest
	}

	return nil
}
