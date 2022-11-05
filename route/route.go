package route

import (
	"Mini-Project_SiBaBe/controller"
	m "Mini-Project_SiBaBe/middleware"
	"Mini-Project_SiBaBe/service"
	"os"

	"github.com/labstack/echo/v4"
	mid "github.com/labstack/echo/v4/middleware"
)

func New(customerSvc service.CustomerSvc, adminSvc service.AdminSvc) *echo.Echo {
	controller.NewCustomerController(customerSvc)
	controller.NewAdminController(adminSvc)

	e := echo.New()
	m.LogMiddleware(e)

	// Routing withouth JWT at Customer Page
	eCust := e.Group("/customer")
	eCust.POST("/register", controller.CreateUser)
	eCust.GET("/products", controller.GetAllProduct)
	eCust.GET("/product/:id", controller.GetProductById)
	eCust.POST("/login", controller.LoginUser)

	eJwt := eCust.Group("/jwt")
	eJwt.Use(mid.JWT([]byte(os.Getenv("SECRET_JWT"))))
	// Routing User with JWT
	eJwt.POST("/product/:id", controller.PostProductToCart)
	eJwt.GET("/cart", controller.GetCart)
	eJwt.GET("/chart/product/:id", controller.UpdateProductFromCartPlus)
	eJwt.GET("/chart/product/:id", controller.UpdateProductFromCartMinus)
	eJwt.GET("/checkout", controller.Checkout)
	eJwt.GET("/checkout/confrim", controller.ConfirmCheckout)

	eJwt.GET("/history", controller.GetHistory)
	eJwt.GET("/history/:id", controller.GetHistoryById)
	eJwt.POST("Feedback/:id", controller.PostFeedback)

	// Routing withouth JWT at Admin Page
	eAdmin := e.Group("/admin")
	eAdmin.POST("/login", controller.LoginAdmin)

	eJwtAdmin := eAdmin.Group("/jwt")
	eJwtAdmin.Use(mid.JWT([]byte(os.Getenv("SECRET_JWT"))))
	// Routing Admin with JWT
	eJwtAdmin.POST("/product", controller.CreateProduct)
	eJwtAdmin.GET("/products", controller.GetAllProductAdmin)
	eJwtAdmin.PUT("/product/:id", controller.UpdateProduct)
	eJwtAdmin.DELETE("/product/:id", controller.DeleteProduct)
	eJwtAdmin.GET("laporan/bulanan", controller.GetMonthlyReport)
	eJwtAdmin.GET("laporan/tahunan", controller.GetYearlyReport)
	eJwtAdmin.POST("/produksi", controller.CreatePrduction)
	eJwtAdmin.GET("/daftar-pemesanan", controller.GetOrderList)
	eJwtAdmin.GET("/daftar-pemesanan/:id/:status", controller.UpdateOrderStatus)

	return e
}
