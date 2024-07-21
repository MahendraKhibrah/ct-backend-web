// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package Route

import (
	"ct-backend/Controller"
	"ct-backend/Repository"
	"ct-backend/Services"
	"gorm.io/gorm"
)

// Injectors from Wire.go:

func AuthDI(db *gorm.DB) *Controller.AuthController {
	authRepository := Repository.AuthRepositoryProvider(db)
	authService := Services.AuthServiceProvider(authRepository)
	authController := Controller.AuthControllerProvider(authService)
	return authController
}
