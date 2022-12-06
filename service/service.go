package service

import (
	"Mini-Project_SiBaBe/dto"
	"Mini-Project_SiBaBe/model"

	"github.com/labstack/echo/v4"
)

type CustomerSvc interface {
	CreateUser(c echo.Context, user model.General_Customer) (model.General_Customer, error)
	GetAllProduct(c echo.Context) ([]model.Product_View_Integrated, error)
	GetProductById(c echo.Context, id int) (model.Detail_Product_View, error)
	LoginUser(c echo.Context, user model.Login_Binding) (dto.Login, error)
	CreateChart(c echo.Context) (model.General_Chart, error)
	PostProductToCart(c echo.Context, id int) (model.General_Product_Chart, error)
	GetCart(c echo.Context) (model.Chart_View, error)
	UpdateProductFromCartPlus(c echo.Context, id int) (model.General_Product_Chart, error)
	UpdateProductFromCartMinus(c echo.Context, id int) (model.General_Product_Chart, error)
	DeleteProductFromCart(c echo.Context, id int) error
	Checkout(c echo.Context) (model.Chart_View, error)
	ConfirmCheckout(c echo.Context, checkout model.Checkout_Binding) (model.Checkout, error)
	ConfirmPayment(c echo.Context, payment model.Payment_Binding) error
	GetHistory(c echo.Context) (model.History_View, error)
	GetHistoryDetail(c echo.Context, id int) (model.Detail_History_View, error)
	CreateFeedbackPemesanan(c echo.Context, id uint) error
	PostFeedback(c echo.Context, feedback_data model.Feedback_Binding) (model.Feedback_Full_View, error)
}

type AdminSvc interface {
	LoginAdmin(c echo.Context, admin model.Login_Binding) (dto.Login, error)
	CreateProduct(c echo.Context, product model.Product_Binding) (model.General_Product, error)
	GetAllProduct(c echo.Context) ([]model.Product_View_Integrated, error)
	UpdateProduct(c echo.Context, id int, product model.Product_Binding) (model.General_Product, error)
	DeleteProduct(c echo.Context, id int) error
	GetMonthlyReport(c echo.Context) ([]model.Monthly_Report_View, error)
	CreateProduction(c echo.Context, production model.Production_Binding) (model.General_Production, error)
	GetOrderList(c echo.Context) ([]model.Order_List, error)
	UpdateOrderStatus(c echo.Context, id int, status model.Update_Order_Status_Binding) (model.General_Order, error)
}
