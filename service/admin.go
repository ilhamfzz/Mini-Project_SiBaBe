package service

import (
	"Mini-Project_SiBaBe/dto"
	"Mini-Project_SiBaBe/middleware"
	"Mini-Project_SiBaBe/model"
	"errors"
	"strconv"
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

	// sort by id
	for i := 0; i < len(productsView); i++ {
		for j := i + 1; j < len(productsView); j++ {
			if productsView[i].Id > productsView[j].Id {
				productsView[i], productsView[j] = productsView[j], productsView[i]
			}
		}
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
		allReport     []model.Laporan_Keuangan
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

		monthlyReportTemp.Month = time.Month(i).String()
		monthlyReportTemp.Year = time.Now().Year()

		for _, report := range allReportView {
			temp := report.Date
			compareDate, err := time.Parse("2006-01-02", temp)
			if err != nil {
				return monthlyReport, errors.New("failed to parse date")
			}
			if compareDate.Month().String() == monthlyReportTemp.Month && compareDate.Year() == time.Now().Year() {
				singleReport = append(singleReport, report)
			}
		}

		monthlyReportTemp.Report = singleReport
		if len(singleReport) == 0 {
			monthlyReportTemp.Report = []model.Money_Report{}
		}
		monthlyReport = append(monthlyReport, monthlyReportTemp)
	}

	return monthlyReport, nil
}

func (as *adminService) CreateProduction(c echo.Context, production model.Production_Binding) (model.General_Production, error) {
	produksi := model.Produksi{
		Date:          production.Date,
		AdminUsername: middleware.ExtractTokenUsername(c),
		TotalBiaya:    production.TotalPrice,
		NamaBarang:    production.Name,
		Gambar:        production.Image,
	}
	err := as.connection.Create(&produksi).Error
	if err != nil {
		return model.General_Production{}, errors.New("failed to create production")
	}

	var report model.Laporan_Keuangan
	err = as.connection.Where("tanggal = ?", produksi.Date).Find(&report).Error
	if err != nil {
		return model.General_Production{}, errors.New("failed to get report")
	}
	if report.Tanggal != "" && (report.TotalPemasukan != 0 || report.TotalPengeluaran != 0) {
		report.TotalPengeluaran += produksi.TotalBiaya
		err = as.connection.Where("tanggal = ?", produksi.Date).Save(&report).Error
		if err != nil {
			return model.General_Production{}, errors.New("failed to update report")
		}
	} else {
		report.Tanggal = produksi.Date
		report.TotalPemasukan = 0
		report.TotalPengeluaran = produksi.TotalBiaya
		err = as.connection.Create(&report).Error
		if err != nil {
			return model.General_Production{}, errors.New("failed to create report")
		}
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

	err = as.connection.Where("status = ? OR status = ? OR status = ?", "Menunggu Validasi", "Terima", "Tolak").Find(&orders).Error
	if err != nil {
		return result, errors.New("failed to get order list")
	}

	for _, order := range orders {
		var (
			singleOrder       model.Order_List
			OrderDetail       []model.Order_Product_View
			singleOrderDetail model.Order_Product_View
			cartsDomain       []model.Produk_Keranjang
		)

		err := as.connection.Where("id_keranjang = ?", order.IdKeranjang).Find(&cartsDomain).Error
		if err != nil {
			return result, errors.New("failed get produk from each order list")
		}

		for _, cart := range cartsDomain {
			singleOrderDetail.ProductID = cart.IdProduk
			singleOrderDetail.Quantity = cart.JumlahProduk
			singleOrderDetail.TotalPrice = cart.TotalHarga
			var product model.Produk
			err = as.connection.Raw("SELECT * FROM produks WHERE id = ?", cart.IdProduk).Scan(&product).Error
			if err != nil {
				return result, errors.New("failed to get produk from each order list")
			}
			ProductOrderDetail := model.Product_View{
				Id:          product.ID,
				Name:        product.Nama,
				Image:       product.Gambar,
				Description: product.Deskripsi,
				Price:       product.Harga,
			}
			singleOrderDetail.Product = ProductOrderDetail
			OrderDetail = append(OrderDetail, singleOrderDetail)
		}

		// sort OrderDetail by product ID
		for i := 0; i < len(OrderDetail); i++ {
			for j := i + 1; j < len(OrderDetail); j++ {
				if OrderDetail[i].ProductID > OrderDetail[j].ProductID {
					OrderDetail[i], OrderDetail[j] = OrderDetail[j], OrderDetail[i]
				}
			}
		}

		singleOrder.Invoice = "P"
		for i := 0; i < 8-len(strconv.Itoa(int(order.ID))); i++ {
			singleOrder.Invoice += "0"
		}
		singleOrder.Invoice += strconv.Itoa(int(order.ID))
		singleOrder.OrderID = order.ID
		singleOrder.CartID = order.IdKeranjang
		singleOrder.FinalPrice = order.TotalHarga

		var user model.Customer
		err = as.connection.Where("username = ?", order.CustomerUsername).Find(&user).Error
		if err != nil {
			return result, errors.New("failed to get user from each order list")
		}

		singleOrder.Customer = user.Nama
		singleOrder.Phone = user.Telp

		addresLink := "https://www.google.co.id/maps/search/"
		var addresConv string
		for _, char := range user.Alamat {
			if char == ' ' {
				addresConv += "+"
			} else {
				addresConv += string(char)
			}
		}
		addresSuffix := "/@-7.279386,112.797269,14z"
		singleOrder.Address = addresLink + addresConv + addresSuffix
		singleOrder.Status = order.Status
		singleOrder.OrderList = OrderDetail

		result = append(result, singleOrder)
	}

	// sort by status Menunggu Validasi in the first
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Status > result[j].Status {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Status == result[j].Status && result[i].OrderID > result[j].OrderID {
				result[i], result[j] = result[j], result[i]
			}
		}
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

	order.Status = status.Status
	if status.Status == "terima" {
		order.Status = "Terima"
	}
	if status.Status != "Terima" {
		order.Status = "Tolak"
	}

	order.DiValidasiOleh = middleware.ExtractTokenUsername(c)
	err = as.connection.Save(&order).Error
	if err != nil {
		return model.General_Order{}, errors.New("failed to update order status")
	}

	if order.Status == "Terima" {
		var report model.Laporan_Keuangan
		date := order.UpdatedAt.Format("2006-01-02")
		err = as.connection.Where("tanggal = ?", date).Find(&report).Error
		if err != nil {
			return model.General_Order{}, errors.New("failed to get report")
		}
		if report.Tanggal != "" && (report.TotalPemasukan != 0 || report.TotalPengeluaran != 0) {
			report.TotalPemasukan += order.TotalHarga
			err = as.connection.Where("tanggal = ?", date).Save(&report).Error
			if err != nil {
				return model.General_Order{}, errors.New("failed to update report")
			}
		} else {
			report.Tanggal = date
			report.TotalPemasukan = order.TotalHarga
			report.TotalPengeluaran = 0
			err = as.connection.Create(&report).Error
			if err != nil {
				return model.General_Order{}, errors.New("failed to create report")
			}
		}

		var carts []model.Produk_Keranjang
		err = as.connection.Where("id_keranjang = ?", order.IdKeranjang).Find(&carts).Error
		if err != nil {
			return model.General_Order{}, errors.New("failed to get cart")
		}
		for _, cart := range carts {
			var product model.Produk
			err = as.connection.Where("id = ?", cart.IdProduk).First(&product).Error
			if err != nil {
				return model.General_Order{}, errors.New("failed to get product")
			}
			product.Stok = product.Stok - cart.JumlahProduk
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
		CartID:           order.IdKeranjang,
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
