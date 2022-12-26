package controller

import (
	"context"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wiormiw/GrowBaks/config"
	"github.com/wiormiw/GrowBaks/helper"
	"github.com/wiormiw/GrowBaks/model"
	"github.com/wiormiw/GrowBaks/service"
)

type (
	// IUserController is an interface that has all the function to be implemented inside user controller
	IUserController interface {
		ListUser(ctx *fiber.Ctx) error
		DetailUser(ctx *fiber.Ctx) error
		DetailUserProfile(ctx *fiber.Ctx) error
		UpdateUser(ctx *fiber.Ctx) error
		UpdateUserProfile(ctx *fiber.Ctx) error
		DeleteUser(ctx *fiber.Ctx) error
	}

	// UserController is an app tag struct that consists of all the dependencies needed for user controller
	UserController struct {
		Context  context.Context
		Config   *config.Configuration
		Logger   *logrus.Logger
		UserSvc  service.IUserService
		LapakSvc service.ILapakService
	}
)

// ListUser responsible to getting all user from controller layer
func (uc *UserController) ListUser(ctx *fiber.Ctx) error {
	var (
		search = ""
	)

	if sr := ctx.Query("s", ""); sr != "" {
		if len(sr) > 0 {
			search = strings.Trim(sr, " ")
		}
	}

	data, err := uc.UserSvc.GetAllUserSvc(search)
	if err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusOK, nil, "Success Getting all Users", data)
}

// DetailUser responsible to get a user by id from controller layer
func (uc *UserController) DetailUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id", "")

	if id == "" {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, nil, model.ErrUserNotFound.Error(), nil)
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	data, err := uc.UserSvc.GetUserByIDSvc(userID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusNotFound, err, err.Error(), nil)
		}

		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusOK, nil, "Success Get User", data)
}

// DetailUserProfile responsible to get detail user profile from controller layer
func (uc *UserController) DetailUserProfile(ctx *fiber.Ctx) error {
	id := ctx.Params("id", "")

	if id == "" {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, nil, model.ErrUserNotFound.Error(), nil)
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	data, err := uc.UserSvc.GetUserProfileByIDSvc(userID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusNotFound, err, err.Error(), nil)
		}

		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusOK, nil, "Success Get User Profile", data)
}

// UpdateUser responsible to update a user by id from controller layer
func (uc *UserController) UpdateUser(ctx *fiber.Ctx) error {
	var userReq model.UpdateUserRequest

	id := ctx.Params("id", "")

	if id == "" {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, nil, model.ErrLapakNotFound.Error(), nil)
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	if err := ctx.BodyParser(&userReq); err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
	}

	err = uc.UserSvc.UpdateUserByIDSvc(userID, userReq)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusNotFound, err, err.Error(), nil)
		}

		if errors.Is(err, model.ErrInvalidRequest) || errors.Is(err, model.ErrEmailExisted) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
		}

		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusOK, nil, "Success Update User", nil)
}

// UpdateUserProfile responsible to update a user by id from controller layer
func (uc *UserController) UpdateUserProfile(ctx *fiber.Ctx) error {
	var userReq model.UpdateUserProfileRequest

	id := ctx.Params("id", "")

	if id == "" {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, nil, model.ErrLapakNotFound.Error(), nil)
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	if err := ctx.BodyParser(&userReq); err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
	}

	err = uc.UserSvc.UpdateUserProfileByIDSvc(userID, userReq)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusNotFound, err, err.Error(), nil)
		}

		if errors.Is(err, model.ErrInvalidRequest) || errors.Is(err, model.ErrEmailExisted) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
		}

		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusOK, nil, "Success Update User Profile", nil)
}

// DeleteUser responsible to delete a user by id from controller layer
func (uc *UserController) DeleteUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id", "")

	if id == "" {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, nil, model.ErrUserNotFound.Error(), nil)
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	err = uc.UserSvc.DeleteUserByIDSvc(userID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusNotFound, err, err.Error(), nil)
		}

		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusOK, nil, "Success Delete User", nil)
}
