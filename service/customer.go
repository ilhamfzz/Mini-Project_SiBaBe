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

type customerService struct {
	connection *gorm.DB
}

func NewCustomerService(db *gorm.DB) CustomerSvc {
	return &customerService{
		connection: db,
	}
}

func (cs *customerService) CreateUser(c echo.Context, user model.Customer) (model.Customer, error) {
	err := cs.connection.Create(&user).Error
	if err != nil {
		return user, errors.New("username sudah terdaftar")
	}
	return user, nil
}

func (cs *customerService) GetAllProduct(c echo.Context) ([]model.Produk, error) {
	var products []model.Produk
	err := cs.connection.Find(&products).Error
	if err != nil {
		return products, errors.New("tidak ada produk")
	}
	return products, nil
}

func (cs *customerService) GetProductById(c echo.Context, id int) (model.Detail_Produk_View, error) {
	var product model.Produk
	err := cs.connection.Where("id = ?", id).First(&product).Error
	if err != nil {
		return model.Detail_Produk_View{}, errors.New("produk tidak ditemukan")
	}

	var productView model.Detail_Produk_View
	productView.Id = product.ID
	productView.Nama = product.Nama
	productView.Gambar = product.Gambar
	productView.Deskripsi = product.Deskripsi
	productView.Harga = product.Harga
	productView.Stok = product.Stok

	var feedback []model.Feedback
	err = cs.connection.Where("id_produk = ?", id).Find(&feedback).Error
	if err != nil {
		return model.Detail_Produk_View{}, errors.New("feedback tidak ditemukan")
	}

	var feedbackView []model.Feedback_View
	for _, f := range feedback {
		var feedbackPemesanan model.Feedback_Pemesanan
		err = cs.connection.Where("id_feedback = ?", f.ID).First(&feedbackPemesanan).Error
		if err != nil {
			return model.Detail_Produk_View{}, errors.New("feedback pemesanan tidak ditemukan")
		}
		feedbackView = append(feedbackView, model.Feedback_View{
			Username: feedbackPemesanan.Username,
			IdProduk: uint(id),
		})
	}

	for _, f := range feedbackView {
		for _, f2 := range feedback {
			if f2.IdProduk == f.IdProduk {
				f.Feedback = f2
			}
		}
	}

	productView.DaftarFeedback = feedbackView

	return productView, nil
}

func (cs *customerService) LoginUser(c echo.Context, user model.Customer) (dto.Login, error) {
	var userLogin model.Customer
	err := cs.connection.Where("username = ? AND password = ?", user.Username, user.Password).First(&userLogin).Error
	if err != nil {
		return dto.Login{}, errors.New("username atau password salah")
	}

	var token string
	token, err = middleware.CreateToken(userLogin.Username, userLogin.Nama)
	if err != nil {
		return dto.Login{}, errors.New("gagal membuat token")
	}
	result := dto.Login{
		Username: userLogin.Username,
		Nama:     userLogin.Nama,
		Token:    token,
	}
	return result, nil
}

func (cs *customerService) CreateChart(c echo.Context) (model.Keranjang, error) {
	var chart model.Keranjang
	chart.Username = middleware.ExtractTokenUsername(c)
	chart.TotalHarga = 0
	chart.Status = "Belum Checkout"
	err := cs.connection.Create(&chart).Error
	if err != nil {
		return chart, errors.New("gagal membuat keranjang")
	}
	return chart, nil
}

func (cs *customerService) PostProductToCart(c echo.Context, id int) (model.Produk_Keranjang, error) {
	var product model.Produk
	err := cs.connection.Where("id = ?", id).First(&product).Error
	if err != nil {
		return model.Produk_Keranjang{}, errors.New("produk tidak ditemukan")
	}

	var chart model.Keranjang
	err = cs.connection.Where("username = ? AND status = ?", middleware.ExtractTokenUsername(c), "Belum Checkout").First(&chart).Error
	if err != nil {
		chart, err = cs.CreateChart(c)
		if err != nil {
			return model.Produk_Keranjang{}, errors.New("gagal membuat keranjang")
		}
	}

	var productFromChart model.Produk_Keranjang
	err = cs.connection.Where("id_produk = ? AND id_keranjang = ?", product.ID, chart.ID).First(&productFromChart).Error
	if err != nil {
		productFromChart.IdProduk = product.ID
		productFromChart.IdKeranjang = chart.ID
		productFromChart.JumlahProduk = 1
		productFromChart.TotalHarga = product.Harga
		err = cs.connection.Create(&productFromChart).Error
		if err != nil {
			return model.Produk_Keranjang{}, errors.New("gagal menambahkan produk ke keranjang baru")
		}
	} else {
		productFromChart.JumlahProduk = productFromChart.JumlahProduk + 1
		productFromChart.TotalHarga = productFromChart.TotalHarga + product.Harga
		err = cs.connection.Where("id_produk = ? AND id_keranjang = ?", product.ID, chart.ID).Updates(&productFromChart).Error
		if err != nil {
			return model.Produk_Keranjang{}, errors.New("gagal menambahkan produk ke keranjang lama")
		}
	}

	chart.TotalHarga = chart.TotalHarga + product.Harga
	err = cs.connection.Save(&chart).Error
	if err != nil {
		return model.Produk_Keranjang{}, errors.New("gagal update total harga keranjang")
	}

	return productFromChart, nil
}

func (cs *customerService) GetCart(c echo.Context) (model.Keranjang_View, error) {
	var chart model.Keranjang
	err := cs.connection.Where("username = ? AND status = ?", middleware.ExtractTokenUsername(c), "Belum Checkout").First(&chart).Error
	if err != nil {
		return model.Keranjang_View{}, errors.New("tidak ada barang di keranjang")
	}
	var jumlah_barang uint = 0
	var result_produk_keranjang_view []model.Produk_Keranjang_View
	var productFromChart []model.Produk_Keranjang
	err = cs.connection.Where("id_keranjang = ?", chart.ID).Find(&productFromChart).Error
	if err != nil {
		return model.Keranjang_View{}, errors.New("tidak ada barang di keranjang")
	}
	for _, produk_keranjang := range productFromChart {
		result_produk_keranjang_view = append(result_produk_keranjang_view, model.Produk_Keranjang_View{
			IdKeranjang:  produk_keranjang.IdKeranjang,
			IdProduk:     produk_keranjang.IdProduk,
			JumlahProduk: produk_keranjang.JumlahProduk,
			TotalHarga:   produk_keranjang.TotalHarga,
		})
		jumlah_barang = jumlah_barang + produk_keranjang.JumlahProduk
	}

	var result_produk_view []model.Produk_View
	for _, v := range result_produk_keranjang_view {
		var product model.Produk
		err = cs.connection.Where("id = ?", v.IdProduk).First(&product).Error
		if err != nil {
			return model.Keranjang_View{}, errors.New("gagal mendapatkan barang di keranjang")
		}
		result_produk_view = append(result_produk_view, model.Produk_View{
			Id:        product.ID,
			Nama:      product.Nama,
			Gambar:    product.Gambar,
			Deskripsi: product.Deskripsi,
			Harga:     product.Harga,
		})
	}

	for i, v := range result_produk_keranjang_view {
		for j, v2 := range result_produk_view {
			if v.IdProduk == v2.Id {
				result_produk_keranjang_view[i].Produk = result_produk_view[j]
			}
		}
	}

	result := model.Keranjang_View{
		Id:           chart.ID,
		Username:     chart.Username,
		JumlahBarang: jumlah_barang,
		TotalHarga:   chart.TotalHarga,
		Produk:       result_produk_keranjang_view,
	}

	return result, nil
}

func (cs *customerService) UpdateProductFromCartPlus(c echo.Context, id int) (model.Produk_Keranjang, error) {
	var product model.Produk
	err := cs.connection.Where("id = ?", id).First(&product).Error
	if err != nil {
		return model.Produk_Keranjang{}, errors.New("produk tidak ditemukan")
	}

	var chart model.Keranjang
	err = cs.connection.Where("username = ? AND status = ?", middleware.ExtractTokenUsername(c), "Belum Checkout").First(&chart).Error
	if err != nil {
		return model.Produk_Keranjang{}, errors.New("keranjang tidak ditemukan")
	}

	var productFromChart model.Produk_Keranjang
	err = cs.connection.Where("id_produk = ? AND id_keranjang = ?", product.ID, chart.ID).First(&productFromChart).Error
	if err != nil {
		return model.Produk_Keranjang{}, errors.New("produk tidak ditemukan di keranjang")
	}

	productFromChart.JumlahProduk = productFromChart.JumlahProduk + 1
	productFromChart.TotalHarga = productFromChart.TotalHarga + product.Harga
	err = cs.connection.Where("id_produk = ? AND id_keranjang = ?", product.ID, chart.ID).Updates(&productFromChart).Error
	if err != nil {
		return model.Produk_Keranjang{}, errors.New("gagal menambahkan produk ke keranjang")
	}

	chart.TotalHarga = chart.TotalHarga + product.Harga
	err = cs.connection.Save(&chart).Error
	if err != nil {
		return model.Produk_Keranjang{}, errors.New("gagal update total harga keranjang")
	}

	return productFromChart, nil
}

func (cs *customerService) UpdateProductFromCartMinus(c echo.Context, id int) (model.Produk_Keranjang, error) {
	var product model.Produk
	err := cs.connection.Where("id = ?", id).First(&product).Error
	if err != nil {
		return model.Produk_Keranjang{}, errors.New("produk tidak ditemukan")
	}

	var chart model.Keranjang
	err = cs.connection.Where("username = ? AND status = ?", middleware.ExtractTokenUsername(c), "Belum Checkout").First(&chart).Error
	if err != nil {
		return model.Produk_Keranjang{}, errors.New("keranjang tidak ditemukan")
	}

	var productFromChart model.Produk_Keranjang
	err = cs.connection.Where("id_produk = ? AND id_keranjang = ?", product.ID, chart.ID).First(&productFromChart).Error
	if err != nil {
		return model.Produk_Keranjang{}, errors.New("produk tidak ditemukan di keranjang")
	}

	if productFromChart.JumlahProduk == 1 {
		err = cs.connection.Delete(&productFromChart).Error
		if err != nil {
			return model.Produk_Keranjang{}, errors.New("gagal menghapus produk dari keranjang")
		}
	} else {
		productFromChart.JumlahProduk = productFromChart.JumlahProduk - 1
		productFromChart.TotalHarga = productFromChart.TotalHarga - product.Harga
		err = cs.connection.Where("id_produk = ? AND id_keranjang = ?", product.ID, chart.ID).Updates(&productFromChart).Error
		if err != nil {
			return model.Produk_Keranjang{}, errors.New("gagal mengurangi produk dari keranjang")
		}
	}

	chart.TotalHarga = chart.TotalHarga - product.Harga
	err = cs.connection.Save(&chart).Error
	if err != nil {
		return model.Produk_Keranjang{}, errors.New("gagal update total harga keranjang")
	}

	return productFromChart, nil
}

func (cs *customerService) Checkout(c echo.Context) (model.Keranjang_View, error) {
	chart, err := cs.GetCart(c)
	if err != nil {
		return model.Keranjang_View{}, errors.New("keranjang tidak ditemukan")
	}

	return chart, nil
}

func (cs *customerService) ConfirmCheckout(c echo.Context, checkout_data model.Checkout_Binding) (model.Checkout, error) {
	chart, err := cs.GetCart(c)
	if err != nil {
		return model.Checkout{}, errors.New("keranjang tidak ditemukan")
	}

	var updateChart model.Keranjang
	err = cs.connection.Where("username = ? AND status = ?", middleware.ExtractTokenUsername(c), "Belum Checkout").First(&updateChart).Error
	if err != nil {
		return model.Checkout{}, errors.New("keranjang tidak ditemukan")
	}

	updateChart.Status = "Telah Checkout"
	err = cs.connection.Save(&updateChart).Error
	if err != nil {
		return model.Checkout{}, errors.New("gagal update status keranjang")
	}

	var customer_data model.Customer
	err = cs.connection.Where("username = ?", middleware.ExtractTokenUsername(c)).First(&customer_data).Error
	if err != nil {
		return model.Checkout{}, errors.New("customer tidak ditemukan")
	}

	var checkout model.Checkout
	if checkout_data.Alamat == "" {
		checkout.Alamat = customer_data.Alamat
	}
	checkout.Kurir = checkout_data.Kurir
	checkout.Total_Harga = chart.TotalHarga + 25000
	checkout.Keranjang = chart

	var pemesanan model.Pemesanan
	pemesanan.IdKeranjang = chart.Id
	pemesanan.CustomerUsername = chart.Username
	pemesanan.JumlahBarang = chart.JumlahBarang
	pemesanan.TotalHarga = checkout.Total_Harga
	pemesanan.Status = "Belum Dibayar"
	pemesanan.Alamat = checkout.Alamat
	pemesanan.Kurir = checkout.Kurir
	pemesanan.BuktiPembayaran = ""
	pemesanan.DiValidasiOleh = ""
	err = cs.connection.Create(&pemesanan).Error
	if err != nil {
		return model.Checkout{}, errors.New("gagal membuat pemesanan")
	}

	var admin_pemesanan model.Admin_Pemesanan
	admin_pemesanan.IdPemesanan = pemesanan.ID
	err = cs.connection.Create(&admin_pemesanan).Error
	if err != nil {
		return model.Checkout{}, errors.New("gagal membuat permintaan pemesanan kepada admin")
	}
	return checkout, nil
}

func (cs *customerService) ConfirmPayment(c echo.Context, payment_data model.Payment_Binding) error {
	var pemesanan model.Pemesanan
	err := cs.connection.Where("username = ? AND status = ?", middleware.ExtractTokenUsername(c), "Belum Dibayar").First(&pemesanan).Error
	if err != nil {
		return errors.New("pemesanan tidak ditemukan")
	}

	pemesanan.BuktiPembayaran = payment_data.BuktiPembayaran
	pemesanan.Status = "Menunggu Validasi"
	err = cs.connection.Save(&pemesanan).Error
	if err != nil {
		return errors.New("gagal update status pemesanan")
	}

	return nil
}

func (cs *customerService) GetHistory(c echo.Context) (model.History_View, error) {
	var pemesanan []model.Pemesanan
	err := cs.connection.Where("customer_username = ?", middleware.ExtractTokenUsername(c)).Find(&pemesanan).Error
	if err != nil {
		return model.History_View{}, errors.New("pemesanan tidak ditemukan")
	}

	return model.History_View{Pemesanan: pemesanan}, nil
}

func (cs *customerService) GetHistoryDetail(c echo.Context, id_pemesanan int) (model.Detail_History_View, error) {
	var pemesanan model.Pemesanan
	err := cs.connection.Where("id = ? AND customer_username = ?", id_pemesanan, middleware.ExtractTokenUsername(c)).First(&pemesanan).Error
	if err != nil {
		return model.Detail_History_View{}, errors.New("pemesanan tidak ditemukan")
	}

	var keranjang model.Keranjang
	err = cs.connection.Where("username = ? AND status = ?", middleware.ExtractTokenUsername(c), "Telah Checkout").First(&keranjang).Error
	if err != nil {
		return model.Detail_History_View{}, errors.New("keranjang tidak ditemukan")
	}

	var produk_keranjang []model.Produk_Keranjang
	err = cs.connection.Where("id_keranjang = ?", keranjang.ID).Find(&produk_keranjang).Error
	if err != nil {
		return model.Detail_History_View{}, errors.New("produk keranjang tidak ditemukan")
	}

	var produk []model.Produk
	for _, produk_keranjang := range produk_keranjang {
		var produk_data model.Produk
		err = cs.connection.Where("id = ?", produk_keranjang.IdProduk).First(&produk_data).Error
		if err != nil {
			return model.Detail_History_View{}, errors.New("produk tidak ditemukan")
		}
		produk = append(produk, produk_data)
	}

	var result model.Detail_History_View
	result.IdPemesanan = pemesanan.ID
	result.IdKeranjang = keranjang.ID
	result.Status = pemesanan.Status
	result.Alamat = pemesanan.Alamat
	result.Kurir = pemesanan.Kurir

	var result_produk_view []model.Produk_View
	for _, produk_data := range produk {
		result_produk_view = append(result_produk_view, model.Produk_View{
			Id:        produk_data.ID,
			Nama:      produk_data.Nama,
			Gambar:    produk_data.Gambar,
			Deskripsi: produk_data.Deskripsi,
			Harga:     produk_data.Harga,
		})
	}

	var result_produk_pemesanan_view []model.Produk_Pemesanan_View
	for i, produk_keranjang := range produk_keranjang {
		result_produk_pemesanan_view = append(result_produk_pemesanan_view, model.Produk_Pemesanan_View{
			JumlahProduk: produk_keranjang.JumlahProduk,
			TotalHarga:   produk_keranjang.TotalHarga,
			Produk:       result_produk_view[i],
		})
	}

	result.Produk = result_produk_pemesanan_view
	return result, nil
}

func (cs *customerService) CreatFeedbackPemesanan(c echo.Context, id uint) error {
	var pemesanan model.Pemesanan
	err := cs.connection.Where("id = ? AND customer_username = ?", id, middleware.ExtractTokenUsername(c)).First(&pemesanan).Error
	if err != nil {
		return errors.New("pemesanan tidak ditemukan")
	}

	var feedback model.Feedback_Pemesanan
	feedback.IdPemesanan = pemesanan.ID
	feedback.Tanggal = time.Now()
	err = cs.connection.Create(&feedback).Error
	if err != nil {
		return errors.New("gagal membuat feedback pemesanan")
	}
	return nil
}

func (cs *customerService) PostFeedback(c echo.Context, feedback_data model.Feedback) (model.Feedback_View, error) {
	err := cs.connection.Create(&feedback_data).Error
	if err != nil {
		return model.Feedback_View{}, errors.New("gagal membuat feedback")
	}

	var feedback_view model.Feedback_View
	feedback_view.Username = middleware.ExtractTokenUsername(c)
	feedback_view.Feedback = feedback_data
	return feedback_view, nil
}
