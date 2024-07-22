// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package Route

import (
	"ct-backend/Controller"
	"ct-backend/Middleware"
	"ct-backend/Repository"
	"ct-backend/Services"
	"gorm.io/gorm"
)

// Injectors from Wire.go:

func AuthDI(db *gorm.DB) *Controller.AuthController {
	authRepository := Repository.AuthRepositoryProvider(db)
	jwtService := Services.JwtServiceProvider()
	authService := Services.AuthServiceProvider(authRepository, jwtService)
	authController := Controller.AuthControllerProvider(authService)
	return authController
}

func ProductDI(db *gorm.DB) *Controller.ProductController {
	productRepository := Repository.ProductRepositoryProvider(db)
	productService := Services.ProductServiceProvider(productRepository)
	productController := Controller.ProductControllerProvider(productService)
	return productController
}

func CommonMiddlewareDI() *Middleware.CommonMiddleware {
	jwtService := Services.JwtServiceProvider()
	commonMiddleware := Middleware.CommonMiddlewareProvider(jwtService)
	return commonMiddleware
}

func PurchaseDI(db *gorm.DB) *Controller.PurchaseController {
	purchaseRepository := Repository.PurchaseRepositoryProvider(db)
	productRepository := Repository.ProductRepositoryProvider(db)
	purchaseService := Services.PurchaseServiceProvider(purchaseRepository, productRepository)
	purchaseController := Controller.PurchaseControllerProvider(purchaseService)
	return purchaseController
}
