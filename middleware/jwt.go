package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/wiormiw/GrowBaks/helper"
	"github.com/wiormiw/GrowBaks/model"
	"github.com/wiormiw/GrowBaks/util"
)

// ValidateJWTMiddleware responsible to validating jwt in header each request
func ValidateJWTMiddleware(ctx *fiber.Ctx) error {
	// validate JWT coming from request, if valid decode into a struct
	decodedPayload, err := util.ValidateJWT(ctx)
	if err != nil {
		return helper.ResponseFormatter[any](ctx, fiber.StatusUnauthorized, err, fmt.Sprintf("Unauthorized access, reason : %s", err.Error()), nil)
	}

	// pass decoded payload into ctx.Locals()
	ctx.Locals(model.KeyJWTValidAccess, decodedPayload)

	// going to next handler..
	return ctx.Next()
}
