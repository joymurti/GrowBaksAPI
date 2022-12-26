package controller

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
	ILapakController interface {
		ListLapak(ctx *fiber.Ctx) error
		ListLapakByLocation(ctx *fiber.Ctx) error
		DetailLapak(ctx *fiber.Ctx) error
		UpdateLapak(ctx *fiber.Ctx) error
		UpdateLapakByStatus(ctx *fiber.Ctx) error
		DeleteLapak(ctx *fiber.Ctx) error
	}

	// TagController is an app tag struct that consists of all the dependencies needed for tag controller
	LapakController struct {
		Context  context.Context
		Config   *config.Configuration
		Logger   *logrus.Logger
		LapakSvc service.ILapakService
	}
)

// ListLapak responsible to getting all lapak from controller layer
func (lapc *LapakController) ListLapak(ctx *fiber.Ctx) error {
	var (
		search = ""
	)

	if sr := ctx.Query("s", ""); sr != "" {
		if len(sr) > 0 {
			search = strings.Trim(sr, " ")
		}
	}

	data, err := lapc.LapakSvc.GetAllLapakSvc(search)
	if err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusOK, nil, "Success Getting all Lapak", data)
}

// ListLapak responsible to getting all lapak from controller layer by location
func (lapc *LapakController) ListLapakByLocation(ctx *fiber.Ctx) error {
	var (
		search = ""
	)

	if sr := ctx.Query("s", ""); sr != "" {
		if len(sr) > 0 {
			search = strings.Trim(sr, " ")
		}
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

	lapc.Logger.Info(fmt.Sprintf("Value Test : %s ID, %s Search", userID, search))

	data, err = lapc.LapakSvc.GetAllLapakByLocationSvc(userID, search)
	if err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusOK, nil, "Success Getting all Lapak By Location", data)
}

// DetailLapak responsible to getting one lapak from controller layer
func (lapc *LapakController) DetailLapak(ctx *fiber.Ctx) error {
	id := ctx.Params("id", "")

	if id == "" {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, nil, model.ErrLapakNotFound.Error(), nil)
	}

	lapakID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	data, err := lapc.LapakSvc.GetLapakByIDSvc(lapakID)
	if err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusOK, nil, "Success Getting Lapak Detail", data)
}

// UpdateLapak responsible to updating a lapak by id from controller layer
func (lapc *LapakController) UpdateLapak(ctx *fiber.Ctx) error {
	var lapakUpdate model.LapakUpdate

	if err := ctx.BodyParser(&lapakUpdate); err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
	}

	id := ctx.Params("id", "")

	lapakID, err := uuid.Parse(id)
	if err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
	}

	err = lapc.LapakSvc.UpdateLapakSvc(lapakID, lapakUpdate)
	if err != nil {
		if errors.Is(err, model.ErrTagNotFound) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusNotFound, err, err.Error(), nil)
		}

		if errors.Is(err, model.ErrInvalidRequest) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
		}

		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusOK, nil, "Success Update Lapak", nil)
}

// UpdateLapakByStatus responsible to updating a lapak by id from controller layer
func (lapc *LapakController) UpdateLapakByStatus(ctx *fiber.Ctx) error {
	var lapakUpdate model.LapakUpdateStatus

	if err := ctx.BodyParser(&lapakUpdate); err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
	}

	id := ctx.Params("id", "")

	lapakID, err := uuid.Parse(id)
	if err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
	}

	err = lapc.LapakSvc.UpdateLapakStatusSvc(lapakID, lapakUpdate)
	if err != nil {
		if errors.Is(err, model.ErrTagNotFound) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusNotFound, err, err.Error(), nil)
		}

		if errors.Is(err, model.ErrInvalidRequest) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
		}

		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusOK, nil, "Success Update Lapak", nil)
}

// DeleteLapak responsible to deleting a lapak by id from controller layer
func (lapc *LapakController) DeleteLapak(ctx *fiber.Ctx) error {
	id := ctx.Params("id", "")

	lapakID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	err = lapc.LapakSvc.DeleteLapakSvc(lapakID)
	if err != nil {
		if errors.Is(err, model.ErrTagNotFound) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusNotFound, err, err.Error(), nil)
		}

		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusOK, nil, "Success Deleting Tags", nil)
}
