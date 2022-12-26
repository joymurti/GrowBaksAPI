package infrastructure

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/wiormiw/GrowBaks/application"
	"github.com/wiormiw/GrowBaks/helper"
	m "github.com/wiormiw/GrowBaks/middleware"
)

// ServeHTTP is wrapper function to start the apps infra in HTTP mode
func ServeHTTP(app *application.App) *fiber.App {
	//call setup router
	setupRouter(app)

	return app.Application
}

// setupRouter is function to manage all routings
func setupRouter(app *application.App) {
	var dep = application.SetupDependencyInjection(app)

	v1 := app.Application.Group("/v1", func(ctx *fiber.Ctx) error {
		ctx.Set("Version", "v1")
		return ctx.Next()
	})

	v1.Options("/*", helper.OptionsHandler)
	v1.Get("/health", dep.HealthCheckController.HealthCheck)

	// AUTH SECTION
	{
		v1.Post("/auth/register", dep.AuthController.RegisterAuthor)
		v1.Post("/auth/login", limiter.New(limiter.Config{
			Expiration: 15 * time.Minute,
			Max:        5,
			LimitReached: func(ctx *fiber.Ctx) error {
				return helper.ResponseFormatter[any](ctx, fiber.StatusTooManyRequests, nil, "Login Attempts already reached the limit, tell our super admin about resetting a password, Thank you", nil)
			},
		}), dep.AuthController.Login)
	}

	// E-COMMERCE SECTION
	{
		// LAPAK SECTION
		v1.Get("/lapak", m.ValidateJWTMiddleware, m.SuperAdminOnlyMiddleware, dep.LapakController.ListLapak)
		v1.Get("/lapak/location", m.ValidateJWTMiddleware, m.CustomerOnlyMiddleware, dep.LapakController.ListLapakByLocation)
		v1.Get("/lapak/:id", m.ValidateJWTMiddleware, dep.LapakController.DetailLapak)
		v1.Put("/lapak/:id", m.ValidateJWTMiddleware, m.AdminSellerOnlyMiddleware, dep.LapakController.UpdateLapak)
		v1.Put("/lapak/:id/status", m.ValidateJWTMiddleware, m.AdminSellerOnlyMiddleware, dep.LapakController.UpdateLapakByStatus)
		v1.Delete("/lapak/:id", m.ValidateJWTMiddleware, m.SuperAdminOnlyMiddleware, dep.LapakController.DeleteLapak)

		// PRODUCT SECTION
		v1.Post("lapak/:id/product/upload", m.ValidateJWTMiddleware, m.AdminSellerOnlyMiddleware, dep.ProductController.UploadIMG)
		v1.Post("lapak/:id/product", m.ValidateJWTMiddleware, m.AdminSellerOnlyMiddleware, dep.ProductController.CreateProduct)
		v1.Get("product/", m.ValidateJWTMiddleware, dep.ProductController.ListProduct)
		v1.Get("product/:id/", m.ValidateJWTMiddleware, dep.ProductController.DetailProduct)
		v1.Get("lapak/:lapak_id/product/", m.ValidateJWTMiddleware, dep.ProductController.ListProductByLapak)

		// PEMESANAN SECTION
		v1.Post("pemesanan/:product_id", m.ValidateJWTMiddleware, m.CustomerOnlyMiddleware, dep.PemesananController.CreatePemesanan)
		v1.Get("pemesanan/", m.ValidateJWTMiddleware, m.SuperAdminOnlyMiddleware, dep.PemesananController.ListPemesanan)
		v1.Get("pemesanan/self", m.ValidateJWTMiddleware, m.CustomerSellerOnlyMiddleware, dep.PemesananController.ListPemesananPribadi)
		v1.Get("pemesanan/:id", m.ValidateJWTMiddleware, dep.PemesananController.DetailPemesanan)
		v1.Put("pemesanan/:id", m.ValidateJWTMiddleware, m.CustomerSellerOnlyMiddleware, dep.PemesananController.UpdatePemesanan)
		v1.Delete("pemesanan/:id", m.ValidateJWTMiddleware, m.CustomerSellerOnlyMiddleware, dep.PemesananController.DeletePemesanan)
	}

	// USER SECTION
	{
		v1.Get("/users", m.ValidateJWTMiddleware, m.SuperAdminOnlyMiddleware, dep.UserController.ListUser)
		v1.Get("/users/:id", m.ValidateJWTMiddleware, m.SuperAdminOnlyMiddleware, dep.UserController.DetailUser)
		v1.Get("/users/:id/profile", m.ValidateJWTMiddleware, dep.UserController.DetailUserProfile)
		v1.Put("/users/:id", m.ValidateJWTMiddleware, dep.UserController.UpdateUser)
		v1.Put("/users/:id/profile", m.ValidateJWTMiddleware, dep.UserController.UpdateUserProfile)
		v1.Delete("/users/:id", m.ValidateJWTMiddleware, m.SuperAdminOnlyMiddleware, dep.UserController.DeleteUser)
	}

	// handler for route not found
	app.Application.Use(func(c *fiber.Ctx) error {
		return helper.ResponseFormatter[any](c, fiber.StatusNotFound, nil, "Route not found", nil)
	})

}
