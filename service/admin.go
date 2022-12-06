package service

import (
	"Mini-Project_SiBaBe/dto"
	"Mini-Project_SiBaBe/middleware"
	"Mini-Project_SiBaBe/model"
	"errors"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type adminService struct {
	connection *gorm.DB
}

func NewAdminService(db *gorm.DB) AdminSvc {
	return &adminService{
		connection: db,
	}
}

func (as *adminService) LoginAdmin(c echo.Context, admin model.Login_Binding) (dto.Login, error) {
	var (
		adminDomain model.Admin
		login       dto.Login
		err         error
	)

	err = as.connection.Where("username = ? AND password = ?", admin.Username, admin.Password).First(&adminDomain).Error
	if err != nil {
		return login, errors.New("username or Password is wrong")
	}

	login.Token, err = middleware.CreateToken(admin.Username, admin.Password)
	if err != nil {
		return login, errors.New("failed to create token")
	}

	login.Username = adminDomain.Username
	login.Name = adminDomain.Nama

	return login, nil
}

func (as *adminService) CreateProduct(c echo.Context, product model.Product_Binding) (model.General_Product, error) {
	productDomain := model.Produk{
		Nama:      product.Name,
		Gambar:    product.Image,
		Stok:      product.Stock,
		Deskripsi: product.Description,
		Harga:     product.Price,
	}
	err := as.connection.Create(&productDomain).Error
	if err != nil {
		return model.General_Product{}, errors.New("failed to create product")
	}

	return model.General_Product{
		Id:          productDomain.ID,
		CreatedAt:   productDomain.CreatedAt,
		UpdatedAt:   productDomain.UpdatedAt,
		Name:        productDomain.Nama,
		Image:       productDomain.Gambar,
		Stock:       productDomain.Stok,
		Description: productDomain.Deskripsi,
		Price:       productDomain.Harga,
	}, nil
}

func (as *adminService) GetAllProduct(c echo.Context) ([]model.Product_View_Integrated, error) {
	var products []model.Produk
	err := as.connection.Find(&products).Error
	if err != nil {
		return []model.Product_View_Integrated{}, errors.New("produk tidak ditemukan")
	}

	var productsView []model.Product_View_Integrated
	for _, p := range products {
		var productView model.Product_View_Integrated
		productView.Id = p.ID
		productView.Name = p.Nama
		productView.Price = p.Harga
		productView.Stock = p.Stok
		productView.Image = p.Gambar
		productView.Description = p.Deskripsi
		var reviews_view []model.Review_View
		var temp_review []model.Feedback
		err = as.connection.Find(&temp_review, "id_produk = ?", p.ID).Error
		if err != nil {
			productView.Reviews = nil
		} else {
			for _, r := range temp_review {
				var review_view model.Review_View
				var temp_feedback_pemesanan model.Feedback_Pemesanan
				err = as.connection.Find(&temp_feedback_pemesanan, "id_feedback = ?", r.ID).Error
				if err != nil {
					productView.Reviews = nil
				} else {
					review_view.Username = temp_feedback_pemesanan.Username
					review_view.Feedback = r.IsiFeedback
					review_view.Rating = r.Rating
					reviews_view = append(reviews_view, review_view)
				}
			}
			productView.Reviews = reviews_view
		}
		productsView = append(productsView, productView)
	}

	return productsView, nil
}

func (as *adminService) UpdateProduct(c echo.Context, id int, product model.Product_Binding) (model.General_Product, error) {
	var (
		updatedProduct model.Produk
		err            error
	)

	err = as.connection.Where("id = ?", id).First(&updatedProduct).Error
	if err != nil {
		return model.General_Product{}, errors.New("failed to get product")
	}

	if product.Name != "" {
		updatedProduct.Nama = product.Name
	}
	if product.Image != "" {
		updatedProduct.Gambar = product.Image
	}
	if product.Stock != 0 {
		updatedProduct.Stok = product.Stock
	}
	if product.Description != "" {
		updatedProduct.Deskripsi = product.Description
	}
	if product.Price != 0 {
		updatedProduct.Harga = product.Price
	}

	err = as.connection.Save(&updatedProduct).Error
	if err != nil {
		return model.General_Product{}, errors.New("failed to update product")
	}

	return model.General_Product{
		Id:          updatedProduct.ID,
		CreatedAt:   updatedProduct.CreatedAt,
		UpdatedAt:   updatedProduct.UpdatedAt,
		Name:        updatedProduct.Nama,
		Image:       updatedProduct.Gambar,
		Stock:       updatedProduct.Stok,
		Description: updatedProduct.Deskripsi,
		Price:       updatedProduct.Harga,
	}, nil
}

func (as *adminService) DeleteProduct(c echo.Context, id int) error {
	var (
		product model.Produk
		err     error
	)

	err = as.connection.Where("id = ?", id).First(&product).Error
	if err != nil {
		return errors.New("failed to get product")
	}

	err = as.connection.Delete(&product).Error
	if err != nil {
		return errors.New("failed to delete product")
	}

	return nil
}

func (as *adminService) GetMonthlyReport(c echo.Context) ([]model.Monthly_Report_View, error) {
	var (
		monthlyReport []model.Monthly_Report_View
		allReport     []model.Laporan_Keuangann
		allReportView []model.Money_Report
		err           error
	)

	err = as.connection.Find(&allReport).Error
	if err != nil {
		return monthlyReport, errors.New("no report found this year")
	}

	for _, r := range allReport {
		var reportView model.Money_Report
		reportView.Date = r.Tanggal
		reportView.Income = r.TotalPemasukan
		reportView.Expense = r.TotalPengeluaran
		allReportView = append(allReportView, reportView)
	}

	for i := 1; i <= 12; i++ {
		var (
			monthlyReportTemp model.Monthly_Report_View
			singleReport      []model.Money_Report
		)

		for _, report := range allReportView {
			if report.Date.Month() == time.Month(i) && report.Date.Year() == time.Now().Year() {
				monthlyReportTemp.Month = report.Date.Month().String()
				monthlyReportTemp.Year = report.Date.Year()
				singleReport = append(singleReport, report)
			}
		}

		monthlyReportTemp.Report = singleReport
		monthlyReport = append(monthlyReport, monthlyReportTemp)
	}

	return monthlyReport, nil
}

func (as *adminService) CreateProduction(c echo.Context, production model.Production_Binding) (model.General_Production, error) {
	produksi := model.Produksi{
		AdminUsername: middleware.ExtractTokenUsername(c),
		TotalBiaya:    production.TotalPrice,
	}
	err := as.connection.Create(&produksi).Error
	if err != nil {
		return model.General_Production{}, errors.New("failed to create production")
	}

	var reportDomain model.Laporan_Keuangann
	reportDomain.Tanggal = produksi.CreatedAt
	reportDomain.TotalPemasukan = 0
	reportDomain.TotalPengeluaran = produksi.TotalBiaya
	err = as.connection.Create(&reportDomain).Error
	if err != nil {
		return model.General_Production{}, errors.New("failed to create report")
	}

	return model.General_Production{
		Id:            produksi.ID,
		CreatedAt:     produksi.CreatedAt,
		UpdatedAt:     produksi.UpdatedAt,
		AdminUsername: produksi.AdminUsername,
		TotalPrice:    produksi.TotalBiaya,
	}, nil
}

func (as *adminService) GetOrderList(c echo.Context) ([]model.Order_List, error) {
	var (
		result []model.Order_List
		orders []model.Pemesanan
		err    error
	)

	err = as.connection.Where("status = ?", "Menunggu Validasi").Find(&orders).Error
	if err != nil {
		return result, errors.New("failed to get order list")
	}

	for _, order := range orders {
		var (
			singleOrder       model.Order_List
			singleOrderDetail []model.Product_View
			chartsDomain      []model.Produk_Keranjang
			charts            []model.General_Product_Chart
		)

		err := as.connection.Where("id_keranjang = ?", order.IdKeranjang).Find(&chartsDomain).Error
		if err != nil {
			return result, errors.New("failed get produk from each order list")
		}

		for _, chart := range chartsDomain {
			chartTemp := model.General_Product_Chart{
				ChartID:    chart.IdKeranjang,
				ProductID:  chart.IdProduk,
				Quantity:   chart.JumlahProduk,
				TotalPrice: chart.TotalHarga,
			}
			charts = append(charts, chartTemp)
		}

		var products []model.Produk
		for _, chart := range charts {
			var product model.Produk
			err := as.connection.Where("id = ?", chart.ProductID).Find(&product).Error
			if err != nil {
				return result, errors.New("failed to get produk from each order list")
			}
			products = append(products, product)
		}

		for _, product := range products {
			var singleProduct model.Product_View
			singleProduct.Id = product.ID
			singleProduct.Name = product.Nama
			singleProduct.Image = product.Gambar
			singleProduct.Description = product.Deskripsi
			singleProduct.Price = product.Harga

			singleOrderDetail = append(singleOrderDetail, singleProduct)
		}

		singleOrder.OrderID = order.ID
		singleOrder.CartID = order.IdKeranjang
		singleOrder.Status = order.Status
		singleOrder.OrderList = singleOrderDetail

		result = append(result, singleOrder)
	}

	return result, nil
}

func (as *adminService) UpdateOrderStatus(c echo.Context, id int, status model.Update_Order_Status_Binding) (model.General_Order, error) {
	var (
		order model.Pemesanan
		err   error
	)

	err = as.connection.Where("id = ?", id).First(&order).Error
	if err != nil {
		return model.General_Order{}, errors.New("failed to get order")
	}

	if status.Status == "terima" {
		order.Status = "Terima"
	}
	order.Status = status.Status
	if status.Status != "Terima" {
		order.Status = "Tolak"
	}

	order.DiValidasiOleh = middleware.ExtractTokenUsername(c)
	err = as.connection.Save(&order).Error
	if err != nil {
		return model.General_Order{}, errors.New("failed to update order status")
	}

	if order.Status == "Terima" {
		var report model.Laporan_Keuangann
		report.Tanggal = order.UpdatedAt
		report.TotalPemasukan = order.TotalHarga
		report.TotalPengeluaran = 0
		err = as.connection.Create(&report).Error
		if err != nil {
			return model.General_Order{}, errors.New("failed to create report")
		}

		var charts []model.Produk_Keranjang
		err = as.connection.Where("id_keranjang = ?", order.IdKeranjang).Find(&charts).Error
		if err != nil {
			return model.General_Order{}, errors.New("failed to get chart")
		}
		for _, chart := range charts {
			var product model.Produk
			err = as.connection.Where("id = ?", chart.IdProduk).First(&product).Error
			if err != nil {
				return model.General_Order{}, errors.New("failed to get product")
			}
			product.Stok = product.Stok - chart.JumlahProduk
			err = as.connection.Save(&product).Error
			if err != nil {
				return model.General_Order{}, errors.New("failed to update product stock")
			}
		}
	}

	admin_choice := model.Admin_Pemesanan{
		IdPemesanan:         order.ID,
		UsernameAdmin:       middleware.ExtractTokenUsername(c),
		UpdateStatusOrderTo: order.Status,
		TanggalValidasi:     order.UpdatedAt,
	}
	err = as.connection.Create(&admin_choice).Error
	if err != nil {
		return model.General_Order{}, errors.New("failed to create log admin choice")
	}

	return model.General_Order{
		Id:               order.ID,
		CreatedAt:        order.CreatedAt,
		UpdatedAt:        order.UpdatedAt,
		ChartID:          order.IdKeranjang,
		CustomerUsername: order.CustomerUsername,
		TotalQty:         order.JumlahBarang,
		TotalPrice:       order.TotalHarga,
		Status:           order.Status,
		Address:          order.Alamat,
		Courier:          order.Kurir,
		ProofOfPayment:   order.BuktiPembayaran,
		ValidatedBy:      order.DiValidasiOleh,
	}, nil
}
