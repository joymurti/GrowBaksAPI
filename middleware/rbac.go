package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wiormiw/GrowBaks/helper"
	"github.com/wiormiw/GrowBaks/model"
	"github.com/wiormiw/GrowBaks/util"
)

// SuperAdminOnlyMiddleware responsible to ensure authorized role is admin
func SuperAdminOnlyMiddleware(ctx *fiber.Ctx) error {
	data := ctx.Locals(model.KeyJWTValidAccess)

	extData, err := util.ExtractPayloadJWT(data)
	if err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	if extData.RoleName == "Admin" {
		return ctx.Next()
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusForbidden, model.ErrForbiddenAccess, model.ErrForbiddenAccess.Error(), nil)
}

// SellerOnlyMiddleware responsible to ensure authorized role is penjual
func SellerOnlyMiddleware(ctx *fiber.Ctx) error {
	data := ctx.Locals(model.KeyJWTValidAccess)

	extData, err := util.ExtractPayloadJWT(data)
	if err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	if extData.RoleName == "Penjual" {
		return ctx.Next()
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusForbidden, model.ErrForbiddenAccess, model.ErrForbiddenAccess.Error(), nil)
}

// AdminSellerOnlyMiddleware responsible to ensure authorized role is penjual & admin
func AdminSellerOnlyMiddleware(ctx *fiber.Ctx) error {
	data := ctx.Locals(model.KeyJWTValidAccess)

	extData, err := util.ExtractPayloadJWT(data)
	if err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	if extData.RoleName == "Penjual" || extData.RoleName == "Admin" {
		return ctx.Next()
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusForbidden, model.ErrForbiddenAccess, model.ErrForbiddenAccess.Error(), nil)
}

// AdminCustomerOnlyMiddleware responsible to ensure authorized role is pembeli & admin
func AdminCustomerOnlyMiddleware(ctx *fiber.Ctx) error {
	data := ctx.Locals(model.KeyJWTValidAccess)

	extData, err := util.ExtractPayloadJWT(data)
	if err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	if extData.RoleName == "Pembeli" || extData.RoleName == "Admin" {
		return ctx.Next()
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusForbidden, model.ErrForbiddenAccess, model.ErrForbiddenAccess.Error(), nil)
}

// CustomerSellerOnlyMiddleware responsible to ensure authorized role is pembeli & admin
func CustomerSellerOnlyMiddleware(ctx *fiber.Ctx) error {
	data := ctx.Locals(model.KeyJWTValidAccess)

	extData, err := util.ExtractPayloadJWT(data)
	if err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	if extData.RoleName == "Pembeli" || extData.RoleName == "Penjual" {
		return ctx.Next()
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusForbidden, model.ErrForbiddenAccess, model.ErrForbiddenAccess.Error(), nil)
}

func CustomerOnlyMiddleware(ctx *fiber.Ctx) error {
	data := ctx.Locals(model.KeyJWTValidAccess)

	extData, err := util.ExtractPayloadJWT(data)
	if err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	if extData.RoleName == "Pembeli" {
		return ctx.Next()
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusForbidden, model.ErrForbiddenAccess, model.ErrForbiddenAccess.Error(), nil)
}
