package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wiormiw/GrowBaks/config"
	"github.com/wiormiw/GrowBaks/helper"
	"github.com/wiormiw/GrowBaks/model"
	"github.com/wiormiw/GrowBaks/service"
)

type (
	// IAuthController is an interface that has all the function to be implemented inside auth controller
	IAuthController interface {
		RegisterAuthor(ctx *fiber.Ctx) error
		Login(ctx *fiber.Ctx) error
	}

	// AuthController is an app auth struct that consists of all the dependencies needed for auth controller
	AuthController struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		AuthSvc service.IAuthService
	}
)

// RegisterAuthor responsible to registering data author from controller layer
func (ac *AuthController) RegisterAuthor(ctx *fiber.Ctx) error {
	var regAuthorReq model.CreateUserRequest

	if err := ctx.BodyParser(&regAuthorReq); err != nil {
		ac.Logger.Error(fmt.Sprintf("Error 1: %s", err))
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, model.ErrFailedParseBody.Error(), nil)
	}

	err := ac.AuthSvc.Create(regAuthorReq)
	if err != nil {
		if errors.Is(err, model.ErrInvalidRequest) || errors.Is(err, model.ErrEmailExisted) {
			ac.Logger.Error(fmt.Sprintf("Error 2: %s", err))
			return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
		}
		ac.Logger.Error(fmt.Sprintf("Error 3: %s", err))
		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusCreated, nil, "Successfully register", nil)
}

// Login responsible to log-in data from controller layer
func (ac *AuthController) Login(ctx *fiber.Ctx) error {
	var loginReq model.LoginRequest

	if err := ctx.BodyParser(&loginReq); err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, model.ErrFailedParseBody.Error(), nil)
	}

	jwt, err := ac.AuthSvc.LoginSvc(loginReq)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusNotFound, err, err.Error(), nil)
		}

		if errors.Is(err, model.ErrInvalidRequest) || errors.Is(err, model.ErrInvalidPassword) || errors.Is(err, model.ErrMismatchLogin) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
		}

		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusOK, nil, "Success Login", jwt)
}
