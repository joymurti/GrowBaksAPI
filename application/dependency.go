package application

import (
	"github.com/wiormiw/GrowBaks/controller"
	"github.com/wiormiw/GrowBaks/repository"
	"github.com/wiormiw/GrowBaks/service"
)

// Dependency can contain anything that will provide data for controller layer
type Dependency struct {
	HealthCheckController controller.IHealthCheckController
	AuthController        controller.IAuthController
	LapakController       controller.ILapakController
	UserController        controller.IUserController
	ProductController     controller.IProductController
	PemesananController   controller.IPemesananController
}

// SetupDependencyInjection is a function to set up dependencies
func SetupDependencyInjection(app *App) *Dependency {
	return &Dependency{
		HealthCheckController: setupHealthCheckDependency(app),
		AuthController:        setupAuthDependency(app),
		LapakController:       setupLapakDependency(app),
		UserController:        setupUserDependency(app),
		ProductController:     setupProductDependency(app),
		PemesananController:   setupPemesananDependency(app),
	}
}

// setupHealthCheckDependency is a function to set up dependencies to be used inside health check controller layer
func setupHealthCheckDependency(app *App) *controller.HealthCheckController {
	healthCheckRepo := &repository.HealthCheckRepository{
		Context: app.Context,
		Config:  app.Config,
		Logger:  app.Logger,
		DB:      app.DB,
	}

	healthCheckSvc := &service.HealthCheckService{
		Context:         app.Context,
		Config:          app.Config,
		Logger:          app.Logger,
		HealthCheckRepo: healthCheckRepo,
	}

	healthCheckCtrl := &controller.HealthCheckController{
		Context:        app.Context,
		Config:         app.Config,
		Logger:         app.Logger,
		HealthCheckSvc: healthCheckSvc,
	}

	return healthCheckCtrl
}

// setupAuthDependency is a function to set up dependencies to be used inside auth controller layer
func setupAuthDependency(app *App) *controller.AuthController {
	// init requester

	authRepo := &repository.AuthRepository{
		Context: app.Context,
		Config:  app.Config,
		Logger:  app.Logger,
		DB:      app.DB,
	}

	roleRepo := &repository.RoleRepository{
		Context: app.Context,
		Config:  app.Config,
		Logger:  app.Logger,
		DB:      app.DB,
	}

	locationRepo := &repository.LocationRepository{
		Context: app.Context,
		Config:  app.Config,
		Logger:  app.Logger,
		DB:      app.DB,
	}

	lapakRepo := &repository.LapakRepository{
		Context: app.Context,
		Config:  app.Config,
		Logger:  app.Logger,
		DB:      app.DB,
	}

	authSvc := &service.AuthService{
		Context:      app.Context,
		Config:       app.Config,
		Logger:       app.Logger,
		AuthRepo:     authRepo,
		RoleRepo:     roleRepo,
		LocationRepo: locationRepo,
		LapakRepo:    lapakRepo,
	}

	authCtrl := &controller.AuthController{
		Context: app.Context,
		Config:  app.Config,
		Logger:  app.Logger,
		AuthSvc: authSvc,
	}

	return authCtrl
}

func setupLapakDependency(app *App) *controller.LapakController {
	lapakRepo := &repository.LapakRepository{
		Context: app.Context,
		Config:  app.Config,
		Logger:  app.Logger,
		DB:      app.DB,
	}

	userRepo := &repository.UserRepository{
		Context: app.Context,
		Config:  app.Config,
		Logger:  app.Logger,
		DB:      app.DB,
	}

	lapakSvc := &service.LapakService{
		Context:   app.Context,
		Config:    app.Config,
		Logger:    app.Logger,
		LapakRepo: lapakRepo,
		UserRepo:  userRepo,
	}

	lapakCtrl := &controller.LapakController{
		Context:  app.Context,
		Config:   app.Config,
		Logger:   app.Logger,
		LapakSvc: lapakSvc,
	}

	return lapakCtrl
}

// setupUserDependency is a function to set up dependencies to be used inside user controller layer
func setupUserDependency(app *App) *controller.UserController {
	userRepo := &repository.UserRepository{
		Context: app.Context,
		Config:  app.Config,
		Logger:  app.Logger,
		DB:      app.DB,
	}

	userSvc := &service.UserService{
		Context:  app.Context,
		Config:   app.Config,
		Logger:   app.Logger,
		UserRepo: userRepo,
	}

	userCtrl := &controller.UserController{
		Context: app.Context,
		Config:  app.Config,
		Logger:  app.Logger,
		UserSvc: userSvc,
	}

	return userCtrl
}

// setupProductDependency is a function to set up dependencies to be used inside product controller layer
func setupProductDependency(app *App) *controller.ProductController {
	productRepo := &repository.ProductRepository{
		Context: app.Context,
		Config:  app.Config,
		Logger:  app.Logger,
		DB:      app.DB,
	}

	userRepo := &repository.UserRepository{
		Context: app.Context,
		Config:  app.Config,
		Logger:  app.Logger,
		DB:      app.DB,
	}

	productSvc := &service.ProductService{
		Context:     app.Context,
		Config:      app.Config,
		Logger:      app.Logger,
		ProductRepo: productRepo,
		UserRepo:    userRepo,
	}

	productCtrl := &controller.ProductController{
		Context:    app.Context,
		Config:     app.Config,
		Logger:     app.Logger,
		ProductSvc: productSvc,
	}

	return productCtrl
}

// setupPemesananDependency is a function to set up dependencies to be used inside pemesanan controller layer
func setupPemesananDependency(app *App) *controller.PemesananController {
	pemesananRepo := &repository.PemesananRepository{
		Context: app.Context,
		Config:  app.Config,
		Logger:  app.Logger,
		DB:      app.DB,
	}

	productRepo := &repository.ProductRepository{
		Context: app.Context,
		Config:  app.Config,
		Logger:  app.Logger,
		DB:      app.DB,
	}

	userRepo := &repository.UserRepository{
		Context: app.Context,
		Config:  app.Config,
		Logger:  app.Logger,
		DB:      app.DB,
	}

	pemesananSvc := &service.PemesananService{
		Context:       app.Context,
		Config:        app.Config,
		Logger:        app.Logger,
		PemesananRepo: pemesananRepo,
		UserRepo:      userRepo,
		ProductRepo:   productRepo,
	}

	pemesananCtrl := &controller.PemesananController{
		Context:      app.Context,
		Config:       app.Config,
		Logger:       app.Logger,
		PemesananSvc: pemesananSvc,
	}

	return pemesananCtrl
}
