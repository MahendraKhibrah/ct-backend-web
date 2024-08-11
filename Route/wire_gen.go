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

func InvoiceDI(db *gorm.DB) *Controller.InvoiceController {
	invoiceRepository := Repository.InvoiceRepositoryProvider(db)
	productRepository := Repository.ProductRepositoryProvider(db)
	invoiceService := Services.InvoiceServiceProvider(invoiceRepository, productRepository)
	invoiceController := Controller.InvoiceControllerProvider(invoiceService)
	return invoiceController
}

func ClientDI(db *gorm.DB) *Controller.ClientController {
	clientRepository := Repository.ClientRepositoryProvider(db)
	clientService := Services.ClientServiceProvider(clientRepository)
	clientController := Controller.ClientControllerProvider(clientService)
	return clientController
}

func DeliveryDI(db *gorm.DB) *Controller.DeliveryController {
	deliveryRepository := Repository.DeliveryRepositoryProvider(db)
	invoiceRepository := Repository.InvoiceRepositoryProvider(db)
	deliveryService := Services.DeliveryServiceProvider(deliveryRepository, invoiceRepository)
	deliveryController := Controller.DeliveryControllerProvider(deliveryService)
	return deliveryController
}

func UserDI(db *gorm.DB) *Controller.UserController {
	userRepository := Repository.UserRepositoryProvider(db)
	userService := Services.UserServiceProvider(userRepository)
	userController := Controller.UserControllerProvider(userService)
	return userController
}
