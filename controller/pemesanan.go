package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wiormiw/GrowBaks/config"
	"github.com/wiormiw/GrowBaks/helper"
	"github.com/wiormiw/GrowBaks/model"
	"github.com/wiormiw/GrowBaks/service"
	"github.com/wiormiw/GrowBaks/util"
)

type (
	// IPemesananController is an interface that has all the function to be implemented inside product controller
	IPemesananController interface {
		CreatePemesanan(ctx *fiber.Ctx) error
		ListPemesanan(ctx *fiber.Ctx) error
		ListPemesananPribadi(ctx *fiber.Ctx) error
		DetailPemesanan(ctx *fiber.Ctx) error
		UpdatePemesanan(ctx *fiber.Ctx) error
		DeletePemesanan(ctx *fiber.Ctx) error
	}

	// Product Controller is an app tag struct that consists of all the dependencies needed for user controller
	PemesananController struct {
		Context      context.Context
		Config       *config.Configuration
		Logger       *logrus.Logger
		PemesananSvc service.IPemesananService
	}
)

// CreatePemesanan responsible to creating a pemesanans from controller layer
func (pemc *PemesananController) CreatePemesanan(ctx *fiber.Ctx) error {
	var pemReq model.CreatePemesananRequest

	product_id := ctx.Params("product_id", "")

	if product_id == "" {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, nil, model.ErrUserNotFound.Error(), nil)
	}

	productID, err := uuid.Parse(product_id)
	if err != nil {
		return err
	}

	data := ctx.Locals(model.KeyJWTValidAccess)
	extData, err := util.ExtractPayloadJWT(data)
	if err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	userID, err := uuid.Parse(extData.UserID)
	if err != nil {
		return err
	}

	if err := ctx.BodyParser(&pemReq); err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
	}

	err = pemc.PemesananSvc.CreatePemesananSvc(userID, productID, pemReq)
	if err != nil {
		if errors.Is(err, model.ErrInvalidRequest) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
		}

		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusCreated, nil, "Success Create Pemesanan", nil)
}

// ListPemesananPribadi responsible to getting all pemesanan pribadi from controller layer
func (pemc *PemesananController) ListPemesananPribadi(ctx *fiber.Ctx) error {
	data := ctx.Locals(model.KeyJWTValidAccess)
	extData, err := util.ExtractPayloadJWT(data)
	if err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	userID, err := uuid.Parse(extData.UserID)
	if err != nil {
		return err
	}

	data, err = pemc.PemesananSvc.GetAllPemesananPribadiSvc(userID)
	if err != nil {
		pemc.Logger.Info(fmt.Sprintf("Error: %s", err))
		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusOK, nil, "Success Getting all Pemesanan", data)
}

// ListPemesanan responsible to getting all pemesanan from controller layer
func (pemc *PemesananController) ListPemesanan(ctx *fiber.Ctx) error {
	data, err := pemc.PemesananSvc.GetAllPemesananSvc()
	if err != nil {
		pemc.Logger.Info(fmt.Sprintf("Error: %s", err))
		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusOK, nil, "Success Getting all Pemesanan", data)
}

func (pemc *PemesananController) DetailPemesanan(ctx *fiber.Ctx) error {
	id := ctx.Params("id", "")

	if id == "" {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, nil, model.ErrPemesananNotFound.Error(), nil)
	}

	pemesananID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	data, err := pemc.PemesananSvc.GetPemesananByIDSvc(pemesananID)
	if err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusOK, nil, "Success Getting Pemesanan Detail", data)
}

// UpdateLapak responsible to updating a lapak by id from controller layer
func (pemc *PemesananController) UpdatePemesanan(ctx *fiber.Ctx) error {
	var pemesananUpdate model.UpdatePemesananRequest

	if err := ctx.BodyParser(&pemesananUpdate); err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
	}

	id := ctx.Params("id", "")

	pemesananID, err := uuid.Parse(id)
	if err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
	}

	err = pemc.PemesananSvc.UpdatePemesananStatusSvc(pemesananID, pemesananUpdate)
	if err != nil {
		if errors.Is(err, model.ErrTagNotFound) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusNotFound, err, err.Error(), nil)
		}

		if errors.Is(err, model.ErrInvalidRequest) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
		}

		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusOK, nil, "Success Update Pemesanan Status", nil)
}

func (pemc *PemesananController) DeletePemesanan(ctx *fiber.Ctx) error {
	id := ctx.Params("id", "")

	if id == "" {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, nil, model.ErrUserNotFound.Error(), nil)
	}

	pemesananID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	err = pemc.PemesananSvc.DeletePemesananSvc(pemesananID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusNotFound, err, err.Error(), nil)
		}

		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusOK, nil, "Success Delete Pemesanan", nil)
}
