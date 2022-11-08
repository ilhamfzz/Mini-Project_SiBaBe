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
	eJwt.GET("/products", controller.GetAllProduct)
	eJwt.GET("/product/:id", controller.GetProductById)
	eJwt.GET("/product-to-cart/:id", controller.PostProductToCart)
	eJwt.GET("/cart", controller.GetCart)
	eJwt.GET("/cart/plus/:id", controller.UpdateProductFromCartPlus)
	eJwt.GET("/cart/minus/:id", controller.UpdateProductFromCartMinus)
	eJwt.GET("/checkout", controller.Checkout)
	eJwt.POST("/checkout/confrim", controller.ConfirmCheckout)
	eJwt.POST("/checkout/confirm/payment", controller.ConfirmPayment)
	eJwt.GET("/history", controller.GetHistory)
	eJwt.GET("/history/:id", controller.GetHistoryDetail)
	eJwt.POST("/history/:id/feedback/:id_produk", controller.PostFeedback)

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
	eJwtAdmin.GET("report/monthly", controller.GetMonthlyReport)
	eJwtAdmin.POST("/production", controller.CreatePrduction)
	eJwtAdmin.GET("/order-list", controller.GetOrderList)
	eJwtAdmin.POST("/order-list/:id", controller.UpdateOrderStatus)

	return e
}
