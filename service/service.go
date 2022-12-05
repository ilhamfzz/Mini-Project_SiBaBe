package service

import (
	"Mini-Project_SiBaBe/dto"
	"Mini-Project_SiBaBe/model"

	"github.com/labstack/echo/v4"
)

type CustomerSvc interface {
	CreateUser(c echo.Context, user model.Customer) (model.Customer, error)
	GetAllProduct(c echo.Context) ([]model.Produk_View_Integrated, error)
	GetProductById(c echo.Context, id int) (model.Detail_Produk_View, error)
	LoginUser(c echo.Context, user model.Customer) (dto.Login, error)
	CreateChart(c echo.Context) (model.Keranjang, error)
	PostProductToCart(c echo.Context, id int) (model.Produk_Keranjang, error)
	GetCart(c echo.Context) (model.Keranjang_View, error)
	UpdateProductFromCartPlus(c echo.Context, id int) (model.Produk_Keranjang, error)
	UpdateProductFromCartMinus(c echo.Context, id int) (model.Produk_Keranjang, error)
	Checkout(c echo.Context) (model.Keranjang_View, error)
	ConfirmCheckout(c echo.Context, checkout model.Checkout_Binding) (model.Checkout, error)
	ConfirmPayment(c echo.Context, payment model.Payment_Binding) error
	GetHistory(c echo.Context) (model.History_View, error)
	GetHistoryDetail(c echo.Context, id int) (model.Detail_History_View, error)
	CreateFeedbackPemesanan(c echo.Context, id uint) error
	PostFeedback(c echo.Context, feedback model.Feedback) (model.Feedback_View, error)
}

type AdminSvc interface {
	LoginAdmin(c echo.Context, admin model.Admin) (dto.Login, error)
	CreateProduct(c echo.Context, product model.Produk) (model.Produk, error)
	GetAllProduct(c echo.Context) ([]model.Produk_View_Integrated, error)
	UpdateProduct(c echo.Context, id int, product model.Produk) (model.Produk, error)
	DeleteProduct(c echo.Context, id int) (model.Produk, error)
	GetMonthlyReport(c echo.Context) ([]model.Laporan_Bulanan_View, error)
	CreateProduction(c echo.Context, production model.Produksi_Binding) (model.Produksi, error)
	GetOrderList(c echo.Context) ([]model.Daftar_Pemesanan, error)
	UpdateOrderStatus(c echo.Context, id int, status model.Update_Order_Status_Binding) (model.Pemesanan, error)
}
