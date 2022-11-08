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

func (as *adminService) LoginAdmin(c echo.Context, admin model.Admin) (dto.Login, error) {
	var (
		login dto.Login
		err   error
	)

	if admin.Username != "admin" || admin.Password != "admin" {
		err = as.connection.Where("username = ? AND password = ?", admin.Username, admin.Password).First(&admin).Error
		if err != nil {
			return login, errors.New("username or Password is wrong")
		}
	}

	login.Token, err = middleware.CreateToken(admin.Username, admin.Password)
	if err != nil {
		return login, errors.New("failed to create token")
	}

	login.Username = admin.Username
	login.Nama = admin.Nama

	return login, nil
}

func (as *adminService) CreateProduct(c echo.Context, product model.Produk) (model.Produk, error) {
	err := as.connection.Create(&product).Error
	if err != nil {
		return product, errors.New("failed to create product")
	}

	return product, nil
}

func (as *adminService) GetAllProduct(c echo.Context) ([]model.Produk, error) {
	var (
		products []model.Produk
		err      error
	)

	err = as.connection.Find(&products).Error
	if err != nil {
		return products, errors.New("failed to get all product")
	}

	return products, nil
}

func (as *adminService) UpdateProduct(c echo.Context, id int, product model.Produk) (model.Produk, error) {
	var (
		updatedProduct model.Produk
		err            error
	)

	err = as.connection.Where("id = ?", id).First(&updatedProduct).Error
	if err != nil {
		return updatedProduct, errors.New("failed to get product")
	}

	if product.Nama != "" {
		updatedProduct.Nama = product.Nama
	}
	if product.Gambar != "" {
		updatedProduct.Gambar = product.Gambar
	}
	updatedProduct.Stok = product.Stok
	if product.Deskripsi != "" {
		updatedProduct.Deskripsi = product.Deskripsi
	}
	if product.Harga != 0 {
		updatedProduct.Harga = product.Harga
	}

	err = as.connection.Save(&updatedProduct).Error
	if err != nil {
		return updatedProduct, errors.New("failed to update product")
	}

	return updatedProduct, nil
}

func (as *adminService) DeleteProduct(c echo.Context, id int) (model.Produk, error) {
	var (
		product model.Produk
		err     error
	)

	err = as.connection.Where("id = ?", id).First(&product).Error
	if err != nil {
		return product, errors.New("failed to get product")
	}

	err = as.connection.Delete(&product).Error
	if err != nil {
		return product, errors.New("failed to delete product")
	}

	return product, nil
}

func (as *adminService) GetMonthlyReport(c echo.Context) ([]model.Laporan_Bulanan_View, error) {
	var (
		monthlyReport  []model.Laporan_Bulanan_View
		pemesanan      []model.Pemesanan
		produksi       []model.Produksi
		flag_pemesanan bool = false
		flag_produksi  bool = false
	)

	if err := as.connection.Where("status = ?", "Selesai").Find(&pemesanan); err != nil {
		flag_pemesanan = true
	}
	if err := as.connection.Find(&produksi); err != nil {
		flag_produksi = true
	}

	if flag_pemesanan && flag_produksi {
		return monthlyReport, errors.New("no order and prodction data")
	}

	for i := 0; i < 12; i++ {
		var (
			monthReport model.Laporan_Bulanan_View
			Report      []model.Laporan_Keuangann
			month       time.Month
		)

		month = time.Month(i + 1)
		monthReport.Bulan = month.String()
		monthReport.Tahun = time.Now().Year()

		if flag_pemesanan && !flag_produksi {
			var singleReport model.Laporan_Keuangann
			for _, produksi := range produksi {
				if produksi.CreatedAt.Month() == month {
					singleReport.TotalPengeluaran = singleReport.TotalPengeluaran + produksi.TotalBiaya
				}
			}
			Report = append(Report, singleReport)
		}
		if !flag_pemesanan {
			for _, pemesanan := range pemesanan {
				if pemesanan.UpdatedAt.Month() == month {
					var singleReport model.Laporan_Keuangann
					singleReport.Tanggal = pemesanan.UpdatedAt
					singleReport.TotalPemasukan = pemesanan.TotalHarga
					singleReport.TotalPengeluaran = 0
					for _, produksi := range produksi {
						if produksi.CreatedAt.Year() == pemesanan.UpdatedAt.Year() && produksi.CreatedAt.Month() == pemesanan.UpdatedAt.Month() && produksi.CreatedAt.Day() == pemesanan.UpdatedAt.Day() {
							singleReport.TotalPengeluaran = produksi.TotalBiaya
						}
					}
					Report = append(Report, singleReport)
				}
			}
		}

		monthReport.Laporan = Report
		monthlyReport = append(monthlyReport, monthReport)
	}

	return monthlyReport, nil
}

func (as *adminService) CreateProduction(c echo.Context, production model.Produksi_Binding) (model.Produksi, error) {
	produksi := model.Produksi{
		AdminUsername: middleware.ExtractTokenUsername(c),
		JumlahBarang:  production.JumlahBarang,
		TotalBiaya:    production.TotalBiaya,
	}
	err := as.connection.Create(&produksi).Error
	if err != nil {
		return model.Produksi{}, errors.New("failed to create production")
	}

	return produksi, nil
}

func (as *adminService) GetOrderList(c echo.Context) ([]model.Daftar_Pemesanan, error) {
	var (
		result []model.Daftar_Pemesanan
		orders []model.Pemesanan
		err    error
	)

	err = as.connection.Where("status = ?", "Menunggu Validasi").Find(&orders).Error
	if err != nil {
		return result, errors.New("failed to get order list")
	}

	for _, order := range orders {
		var (
			singleOrder       model.Daftar_Pemesanan
			singleOrderDetail []model.Produk_View
			charts            []model.Produk_Keranjang
		)

		err := as.connection.Where("id = ?", order.IdKeranjang).Find(&charts).Error
		if err != nil {
			return result, errors.New("failed to produk from each order list")
		}

		var products []model.Produk
		for _, chart := range charts {
			var product model.Produk
			err := as.connection.Where("id = ?", chart.IdProduk).Find(&product).Error
			if err != nil {
				return result, errors.New("failed to get produk from each order list")
			}
			products = append(products, product)
		}

		for _, product := range products {
			var singleProduct model.Produk_View
			singleProduct.Id = product.ID
			singleProduct.Nama = product.Nama
			singleProduct.Gambar = product.Gambar
			singleProduct.Deskripsi = product.Deskripsi
			singleProduct.Harga = product.Harga

			singleOrderDetail = append(singleOrderDetail, singleProduct)
		}

		singleOrder.IdPemesanan = order.ID
		singleOrder.IdKeranjang = order.IdKeranjang
		singleOrder.Status = order.Status
		singleOrder.Daftar_Pemesanan = singleOrderDetail

		result = append(result, singleOrder)
	}

	return result, nil
}

func (as *adminService) UpdateOrderStatus(c echo.Context, id int, status model.Update_Order_Status_Binding) (model.Pemesanan, error) {
	var (
		order model.Pemesanan
		err   error
	)

	err = as.connection.Where("id = ?", id).First(&order).Error
	if err != nil {
		return order, errors.New("failed to get order")
	}

	order.Status = status.Status
	if status.Status != "Terima" {
		order.Status = "Tolak"
	}
	err = as.connection.Save(&order).Error
	if err != nil {
		return order, errors.New("failed to update order status")
	}

	return order, nil
}
