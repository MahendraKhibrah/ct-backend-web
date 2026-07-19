//go:build wireinject
// +build wireinject

package Route

import (
	"ct-backend/Controller"
	"ct-backend/Middleware"
	"ct-backend/Repository"
	"ct-backend/Services"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func AuthDI(db *gorm.DB) *Controller.AuthController {
	panic(wire.Build(wire.NewSet(
		Repository.AuthRepositoryProvider,
		Services.AuthServiceProvider,
		Controller.AuthControllerProvider,
		Services.JwtServiceProvider,

		wire.Bind(new(Controller.IAuthController), new(*Controller.AuthController)),
		wire.Bind(new(Services.IAuthService), new(*Services.AuthService)),
		wire.Bind(new(Repository.IAuthRepository), new(*Repository.AuthRepository)),
		wire.Bind(new(Services.IJwtService), new(*Services.JwtService)),
	),
	))
	return &Controller.AuthController{}
}

func ProductDI(db *gorm.DB) *Controller.ProductController {
	panic(wire.Build(wire.NewSet(
		Repository.ProductRepositoryProvider,
		Services.ProductServiceProvider,
		Controller.ProductControllerProvider,

		wire.Bind(new(Controller.IProductController), new(*Controller.ProductController)),
		wire.Bind(new(Services.IProductService), new(*Services.ProductService)),
		wire.Bind(new(Repository.IProductRepository), new(*Repository.ProductRepository)),
	),
	))
	return &Controller.ProductController{}
}

func CommonMiddlewareDI() *Middleware.CommonMiddleware {
	panic(wire.Build(wire.NewSet(
		Middleware.CommonMiddlewareProvider,
		Services.JwtServiceProvider,

		wire.Bind(new(Services.IJwtService), new(*Services.JwtService)),
	),
	))
	return &Middleware.CommonMiddleware{}
}

func PurchaseDI(db *gorm.DB) *Controller.PurchaseController {
	panic(wire.Build(wire.NewSet(
		Repository.ProductRepositoryProvider,
		Repository.PurchaseRepositoryProvider,
		Services.PurchaseServiceProvider,
		Controller.PurchaseControllerProvider,

		wire.Bind(new(Controller.IPurchaseController), new(*Controller.PurchaseController)),
		wire.Bind(new(Services.IPurchaseService), new(*Services.PurchaseService)),
		wire.Bind(new(Repository.IProductRepository), new(*Repository.ProductRepository)),
		wire.Bind(new(Repository.IPurchaseRepository), new(*Repository.PurchaseRepository)),
	),
	))
	return &Controller.PurchaseController{}
}

func InvoiceDI(db *gorm.DB) *Controller.InvoiceController {
	panic(wire.Build(wire.NewSet(
		Repository.ProductRepositoryProvider,
		Repository.InvoiceRepositoryProvider,
		Services.InvoiceServiceProvider,
		Services.MinioServiceProvider,
		Controller.InvoiceControllerProvider,
		Repository.DeliveryRepositoryProvider,

		wire.Bind(new(Controller.IInvoiceController), new(*Controller.InvoiceController)),
		wire.Bind(new(Services.IStorageService), new(*Services.MinioService)),
		wire.Bind(new(Services.IInvoiceService), new(*Services.InvoiceService)),
		wire.Bind(new(Repository.IInvoiceRepository), new(*Repository.InvoiceRepository)),
		wire.Bind(new(Repository.IProductRepository), new(*Repository.ProductRepository)),
		wire.Bind(new(Repository.IDeliveryRepository), new(*Repository.DeliveryRepository)),
	),
	))
	return &Controller.InvoiceController{}
}

func ClientDI(db *gorm.DB) *Controller.ClientController {
	panic(wire.Build(wire.NewSet(
		Repository.ClientRepositoryProvider,
		Services.ClientServiceProvider,
		Controller.ClientControllerProvider,

		wire.Bind(new(Controller.IClientController), new(*Controller.ClientController)),
		wire.Bind(new(Services.IClientService), new(*Services.ClientService)),
		wire.Bind(new(Repository.IClientRepository), new(*Repository.ClientRepository)),
	),
	))
	return &Controller.ClientController{}
}

func DeliveryDI(db *gorm.DB) *Controller.DeliveryController {
	panic(wire.Build(wire.NewSet(
		Repository.InvoiceRepositoryProvider,
		Repository.DeliveryRepositoryProvider,
		Services.DeliveryServiceProvider,
		Controller.DeliveryControllerProvider,

		wire.Bind(new(Controller.IDeliveryController), new(*Controller.DeliveryController)),
		wire.Bind(new(Services.IDeliveryService), new(*Services.DeliveryService)),
		wire.Bind(new(Repository.IDeliveryRepository), new(*Repository.DeliveryRepository)),
		wire.Bind(new(Repository.IInvoiceRepository), new(*Repository.InvoiceRepository)),
	),
	))
	return &Controller.DeliveryController{}
}

func UserDI(db *gorm.DB) *Controller.UserController {
	panic(wire.Build(wire.NewSet(
		Repository.UserRepositoryProvider,
		Services.UserServiceProvider,
		Controller.UserControllerProvider,

		wire.Bind(new(Controller.IUserController), new(*Controller.UserController)),
		wire.Bind(new(Services.IUserService), new(*Services.UserService)),
		wire.Bind(new(Repository.IUserRepository), new(*Repository.UserRepository)),
	),
	))
	return &Controller.UserController{}
}

func ReceiptDI(db *gorm.DB) *Controller.ReceiptController {
	panic(wire.Build(wire.NewSet(
		Repository.InvoiceRepositoryProvider,
		Repository.ReceiptRepositoryProvider,
		Services.ReceiptServiceProvider,
		Controller.ReceiptControllerProvider,

		wire.Bind(new(Controller.IReceiptController), new(*Controller.ReceiptController)),
		wire.Bind(new(Services.IReceiptService), new(*Services.ReceiptService)),
		wire.Bind(new(Repository.IReceiptRepository), new(*Repository.ReceiptRepository)),
		wire.Bind(new(Repository.IInvoiceRepository), new(*Repository.InvoiceRepository)),
	),
	))
	return &Controller.ReceiptController{}
}

func DashboardDI(db *gorm.DB) *Controller.DashboardController {
	panic(wire.Build(wire.NewSet(
		Repository.DashboardRepositoryProvider,
		Services.DashboardServiceProvider,
		Controller.DashboardControllerProvider,

		wire.Bind(new(Controller.IDashboardController), new(*Controller.DashboardController)),
		wire.Bind(new(Services.IDashboardService), new(*Services.DashboardService)),
		wire.Bind(new(Repository.IDashboardRepository), new(*Repository.DashboardRepository)),
	),
	))
	return &Controller.DashboardController{}
}
