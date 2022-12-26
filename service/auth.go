package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/wiormiw/GrowBaks/config"
	"github.com/wiormiw/GrowBaks/helper"
	"github.com/wiormiw/GrowBaks/model"
	"github.com/wiormiw/GrowBaks/repository"
	"github.com/wiormiw/GrowBaks/util"

	"github.com/sirupsen/logrus"
)

type (
	// IAuthService is an interface that has all the function to be implemented inside auth service
	IAuthService interface {
		Create(user model.CreateUserRequest) error
		LoginSvc(user model.LoginRequest) (*model.SuccessLoginResponse, error)
	}

	// AuthService is an app auth struct that consists of all the dependencies needed for auth service
	AuthService struct {
		Context      context.Context
		Config       *config.Configuration
		Logger       *logrus.Logger
		AuthRepo     repository.IAuthRepository
		RoleRepo     repository.IRoleRepository
		LocationRepo repository.ILocationRepository
		LapakRepo    repository.ILapakRepository
	}
)

// Create service layer for handling create a user
func (as *AuthService) Create(user model.CreateUserRequest) error {
	err := validateRegisterAuthorRequest(&user)
	if err != nil {
		as.Logger.Error(fmt.Sprintf("Error 4: %s", err))
		return err
	}

	_, err = as.AuthRepo.GetUserByEmail(user.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			as.Logger.Error(fmt.Sprintf("Error 5: %s", err))
			user.Password, err = helper.HashPassword(user.Password)
			if err != nil {
				as.Logger.Error(fmt.Sprintf("Error 6: %s", err))
			}

			// TO DO
			// 1. Cek IsSeller true/false, true = get role id dengan name penjual, false = name pembeli
			// 2. Get RoleID, assign ke model auth.CreateUserRequest
			if user.IsSeller {
				role, err := as.RoleRepo.GetRoleByName(model.RoleSeller)
				if err != nil {
					as.Logger.Error(fmt.Sprintf("Error 7: %s", err))
				}
				user.RoleID = role.ID
			} else {
				role, err := as.RoleRepo.GetRoleByName(model.RoleCustomer)
				if err != nil {
					as.Logger.Error(fmt.Sprintf("Error 8: %s", err))
				}
				user.RoleID = role.ID
			}

			locationID, err := uuid.Parse(user.LocationID)
			if err != nil {
				as.Logger.Error(fmt.Sprintf("Error 9: %s", err))
				return err
			}

			_, err = as.LocationRepo.GetLocationByID(locationID)
			if err != nil {
				as.Logger.Error(fmt.Sprintf("Error 10: %s", err))
				return err
			}

			err = as.AuthRepo.CreateUser(user, user.IsSeller)
			if err != nil {
				as.Logger.Error(fmt.Sprintf("Error 11: %s", err))
				return err
			}

			return nil
		}
		as.Logger.Error(fmt.Sprintf("Error 12: %s", err))
		return err
	}

	return model.ErrEmailExisted
}

// LoginSvc service layer for handling user login (author,superadmin)
func (as *AuthService) LoginSvc(user model.LoginRequest) (*model.SuccessLoginResponse, error) {
	err := validateLoginRequest(&user)
	if err != nil {
		return nil, err
	}

	getUser, err := as.AuthRepo.GetUserByEmail(user.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, model.ErrUserNotFound
		}

		return nil, err
	}

	if !helper.CheckPasswordHash(getUser.Password, user.Password) {
		return nil, model.ErrInvalidPassword
	}

	jwt, err := util.BuildJWT(as.Config, getUser)
	if err != nil {
		as.Logger.Error(fmt.Errorf("AuthService.BuildJWT ERROR : %v MSG : %s", err, err.Error()))
		return nil, err
	}

	return &model.SuccessLoginResponse{
		AccessToken: jwt,
		ExpiredAt:   time.Now().Add(util.JWTExpire),
		Role:        getUser.RoleName,
	}, nil
}

// validateRegisterAuthorRequest responsible to validating request register author
func validateRegisterAuthorRequest(user *model.CreateUserRequest) error {
	if len(user.FullName) < 5 || len(user.Password) < 5 {
		return model.ErrInvalidRequest
	}

	if !model.IsAllowedEmailInput.MatchString(user.Email) {
		return model.ErrInvalidRequest
	}

	if !model.IsAllowedTeleponInput.MatchString(user.Telepon) {
		return model.ErrInvalidRequest
	}

	return nil
}

// validateLoginRequest responsible to validating request login data
func validateLoginRequest(user *model.LoginRequest) error {
	if !model.IsAllowedEmailInput.MatchString(user.Email) {
		return model.ErrInvalidRequest
	}

	return nil
}
