package controller

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
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
	// IProductController is an interface that has all the function to be implemented inside product controller
	IProductController interface {
		UploadIMG(ctx *fiber.Ctx) error
		CreateProduct(ctx *fiber.Ctx) error
		ListProduct(ctx *fiber.Ctx) error
		ListProductByLapak(ctx *fiber.Ctx) error
		DetailProduct(ctx *fiber.Ctx) error
		// DetailUserProfile(ctx *fiber.Ctx) error
		// UpdateUser(ctx *fiber.Ctx) error
		// UpdateUserProfile(ctx *fiber.Ctx) error
		// DeleteUser(ctx *fiber.Ctx) error
	}

	// Product Controller is an app tag struct that consists of all the dependencies needed for user controller
	ProductController struct {
		Context    context.Context
		Config     *config.Configuration
		Logger     *logrus.Logger
		ProductSvc service.IProductService
	}
)

// Upload IMG Endpoint
func (pc *ProductController) UploadIMG(ctx *fiber.Ctx) error {
	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		pc.Logger.Error(fmt.Errorf("ERROR : %v MSG : %s", err, err.Error()))
		if errors.Is(err, model.ErrInvalidRequest) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
		}

		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	file, err := fileHeader.Open()

	if err != nil {
		pc.Logger.Error(fmt.Errorf("ERROR : %v MSG : %s", err, err.Error()))
		if errors.Is(err, model.ErrInvalidRequest) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
		}

		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	cldService, err := cloudinary.NewFromURL(pc.Config.Const.CloudinaryUrl)
	if err != nil {
		pc.Logger.Error(fmt.Errorf("ERROR : %v MSG : %s", err, err.Error()))
		if errors.Is(err, model.ErrInvalidRequest) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
		}

		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	context := context.Background()

	resp, err := cldService.Upload.Upload(context, file, uploader.UploadParams{
		Folder: "GrowBaks",
	})
	if err != nil {
		pc.Logger.Error(fmt.Errorf("ERROR : %v MSG : %s", err, err.Error()))
		if errors.Is(err, model.ErrInvalidRequest) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
		}

		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusCreated, nil, "Success Upload Image", resp.SecureURL)
}

// CreateProduct responsible to creating a product from controller layer
func (pc *ProductController) CreateProduct(ctx *fiber.Ctx) error {
	var prodReq model.CreateProductRequest

	id := ctx.Params("id", "")

	if id == "" {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, nil, model.ErrUserNotFound.Error(), nil)
	}

	lapakID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	if err := ctx.BodyParser(&prodReq); err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
	}

	err = pc.ProductSvc.CreateProductSvc(lapakID, prodReq)
	if err != nil {
		if errors.Is(err, model.ErrInvalidRequest) {
			return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, err, err.Error(), nil)
		}

		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusCreated, nil, "Success Create Product", nil)
}

// ListProduct responsible to getting all product from controller layer
func (pc *ProductController) ListProduct(ctx *fiber.Ctx) error {
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

	pc.Logger.Info(fmt.Sprintf("Value Test : %s ID, %s Search", userID, search))

	data, err = pc.ProductSvc.GetAllProductSvc(userID, search)
	if err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusOK, nil, "Success Getting all Product", data)
}

func (pc *ProductController) ListProductByLapak(ctx *fiber.Ctx) error {
	lapak_id := ctx.Params("lapak_id", "")

	if lapak_id == "" {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, nil, model.ErrUserNotFound.Error(), nil)
	}

	lapakID, err := uuid.Parse(lapak_id)
	if err != nil {
		return err
	}

	data, err := pc.ProductSvc.GetAllProductByLapakSvc(lapakID)
	if err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusOK, nil, "Success Getting all Product By Lapak", data)
}

func (pc *ProductController) DetailProduct(ctx *fiber.Ctx) error {
	id := ctx.Params("id", "")

	if id == "" {
		return helper.ResponseFormatter[any](ctx, fiber.StatusBadRequest, nil, model.ErrProductNotFound.Error(), nil)
	}

	productID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	data, err := pc.ProductSvc.GetProductByIDSvc(productID)
	if err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusInternalServerError, err, err.Error(), nil)
	}

	return helper.ResponseFormatter[any](ctx, fiber.StatusOK, nil, "Success Getting Product Detail", data)
}
