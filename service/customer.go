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

func (cs *customerService) CreateUser(c echo.Context, user model.General_Customer) (model.General_Customer, error) {
	userDomain := model.Customer{
		Username:  user.Username,
		Password:  user.Password,
		Nama:      user.Name,
		Umur:      user.Age,
		Email:     user.Email,
		Telp:      user.Phone,
		Alamat:    user.Address,
		CreatedAt: time.Now(),
	}

	err := cs.connection.Create(&userDomain).Error
	if err != nil {
		return user, errors.New("username sudah terdaftar")
	}
	user.CreatedAt = userDomain.CreatedAt
	return user, nil
}

func (cs *customerService) GetAllProduct(c echo.Context) ([]model.Product_View_Integrated, error) {
	var products []model.Produk
	err := cs.connection.Find(&products).Error
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
		err = cs.connection.Find(&temp_review, "id_produk = ?", p.ID).Error
		if err != nil {
			productView.Reviews = nil
		} else {
			for _, r := range temp_review {
				var review_view model.Review_View
				var temp_feedback_pemesanan model.Feedback_Pemesanan
				err = cs.connection.Find(&temp_feedback_pemesanan, "id_feedback = ?", r.ID).Error
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

func (cs *customerService) GetProductById(c echo.Context, id int) (model.Detail_Product_View, error) {
	var product model.Produk
	err := cs.connection.Where("id = ?", id).First(&product).Error
	if err != nil {
		return model.Detail_Product_View{}, errors.New("produk tidak ditemukan")
	}

	var productView model.Detail_Product_View
	productView.Id = product.ID
	productView.Name = product.Nama
	productView.Image = product.Gambar
	productView.Description = product.Deskripsi
	productView.Price = product.Harga
	productView.Stock = product.Stok

	var feedback []model.Feedback
	err = cs.connection.Where("id_produk = ?", id).Find(&feedback).Error
	if err != nil {
		return model.Detail_Product_View{}, errors.New("feedback tidak ditemukan")
	}

	var feedbackView []model.Feedback_Full_View
	for _, f := range feedback {
		var FeedbackViewTemp model.Feedback_View
		FeedbackViewTemp.Id = f.ID
		FeedbackViewTemp.CreatedAt = f.CreatedAt
		FeedbackViewTemp.ProductID = f.IdProduk
		FeedbackViewTemp.Comment = f.IsiFeedback
		FeedbackViewTemp.Rating = f.Rating
		var feedbackPemesanan model.Feedback_Pemesanan
		err = cs.connection.Where("id_feedback = ?", f.ID).First(&feedbackPemesanan).Error
		if err != nil {
			return model.Detail_Product_View{}, errors.New("feedback pemesanan tidak ditemukan")
		}
		feedbackView = append(feedbackView, model.Feedback_Full_View{
			Username:  feedbackPemesanan.Username,
			ProductId: uint(id),
			Feedback:  FeedbackViewTemp,
		})
	}

	productView.FeedbackList = feedbackView

	return productView, nil
}

func (cs *customerService) LoginUser(c echo.Context, user model.Login_Binding) (dto.Login, error) {
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
		Name:     userLogin.Nama,
		Token:    token,
	}
	return result, nil
}

func (cs *customerService) CreateChart(c echo.Context) (model.General_Chart, error) {
	var chart model.Keranjang
	chart.Username = middleware.ExtractTokenUsername(c)
	chart.TotalHarga = 0
	chart.Status = "Belum Checkout"
	err := cs.connection.Create(&chart).Error
	if err != nil {
		return model.General_Chart{}, errors.New("gagal membuat keranjang")
	}
	return model.General_Chart{
		Id:         chart.ID,
		CreatedAt:  chart.CreatedAt,
		UpdatedAt:  chart.UpdatedAt,
		Username:   chart.Username,
		TotalPrice: chart.TotalHarga,
		Status:     chart.Status,
	}, nil
}

func (cs *customerService) PostProductToCart(c echo.Context, id int) (model.General_Product_Chart, error) {
	var product model.Produk
	err := cs.connection.Where("id = ?", id).First(&product).Error
	if err != nil {
		return model.General_Product_Chart{}, errors.New("produk tidak ditemukan")
	}

	var chart model.Keranjang
	_ = cs.connection.Where("username = ? AND status = ?", middleware.ExtractTokenUsername(c), "Belum Checkout").Find(&chart).Error
	if chart.Username == "" && chart.TotalHarga == 0 && chart.Status == "" {
		_, err = cs.CreateChart(c)
		if err != nil {
			return model.General_Product_Chart{}, errors.New("gagal membuat keranjang")
		}
		err = cs.connection.Where("username = ? AND status = ?", middleware.ExtractTokenUsername(c), "Belum Checkout").Find(&chart).Error
		if err != nil {
			return model.General_Product_Chart{}, errors.New("gagal mendapatkan keranjang")
		}
	}

	var productFromChart model.Produk_Keranjang
	err = cs.connection.Where("id_produk = ? AND id_keranjang = ?", product.ID, chart.ID).Find(&productFromChart).Error
	if err != nil {
		return model.General_Product_Chart{}, errors.New("produk tidak ditemukan di keranjang")
	}
	if productFromChart.IdProduk == 0 && productFromChart.IdKeranjang == 0 && productFromChart.JumlahProduk == 0 && productFromChart.TotalHarga == 0 {
		productFromChart.IdProduk = product.ID
		productFromChart.IdKeranjang = chart.ID
		productFromChart.JumlahProduk = 1
		productFromChart.TotalHarga = product.Harga
		err = cs.connection.Create(&productFromChart).Error
		if err != nil {
			return model.General_Product_Chart{}, errors.New("gagal menambahkan produk ke keranjang baru")
		}
	} else {
		productFromChart.JumlahProduk = productFromChart.JumlahProduk + 1
		productFromChart.TotalHarga = productFromChart.TotalHarga + product.Harga
		err = cs.connection.Where("id_produk = ? AND id_keranjang = ?", product.ID, chart.ID).Updates(&productFromChart).Error
		if err != nil {
			return model.General_Product_Chart{}, errors.New("gagal menambahkan produk ke keranjang lama")
		}
	}

	chart.TotalHarga = chart.TotalHarga + product.Harga
	err = cs.connection.Save(&chart).Error
	if err != nil {
		return model.General_Product_Chart{}, errors.New("gagal update total harga keranjang")
	}

	return model.General_Product_Chart{
		ChartID:    productFromChart.IdKeranjang,
		ProductID:  productFromChart.IdProduk,
		Quantity:   productFromChart.JumlahProduk,
		TotalPrice: productFromChart.TotalHarga,
	}, nil
}

func (cs *customerService) GetCart(c echo.Context) (model.Chart_View, error) {
	var chart model.Keranjang
	err := cs.connection.Where("username = ? AND status = ?", middleware.ExtractTokenUsername(c), "Belum Checkout").Find(&chart).Error
	if err != nil {
		return model.Chart_View{}, errors.New("keranjang tidak ditemukan")
	}
	if chart.ID == 0 && chart.Username == "" && chart.TotalHarga == 0 && chart.Status == "" {
		return model.Chart_View{}, errors.New("tidak ada barang di keranjang")
	}
	var jumlah_barang uint = 0
	var result_product_Chart_view []model.Product_Chart_View
	var productFromChart []model.Produk_Keranjang
	err = cs.connection.Where("id_keranjang = ?", chart.ID).Find(&productFromChart).Error
	if err != nil {
		return model.Chart_View{}, errors.New("tidak ada barang di keranjang")
	}
	for _, produk_keranjang := range productFromChart {
		result_product_Chart_view = append(result_product_Chart_view, model.Product_Chart_View{
			CartID:     produk_keranjang.IdKeranjang,
			ProductID:  produk_keranjang.IdProduk,
			Quantity:   produk_keranjang.JumlahProduk,
			TotalPrice: produk_keranjang.TotalHarga,
		})
		jumlah_barang = jumlah_barang + produk_keranjang.JumlahProduk
	}

	var result_product_view []model.Product_View
	for _, v := range result_product_Chart_view {
		var product model.Produk
		err = cs.connection.Where("id = ?", v.ProductID).First(&product).Error
		if err != nil {
			return model.Chart_View{}, errors.New("gagal mendapatkan barang di keranjang")
		}
		result_product_view = append(result_product_view, model.Product_View{
			Id:          product.ID,
			Name:        product.Nama,
			Image:       product.Gambar,
			Description: product.Deskripsi,
			Price:       product.Harga,
		})
	}

	for i, v := range result_product_Chart_view {
		for j, v2 := range result_product_view {
			if v.ProductID == v2.Id {
				result_product_Chart_view[i].Product = result_product_view[j]
			}
		}
	}

	// sort result_product_Chart_view by product id
	for i := 0; i < len(result_product_Chart_view); i++ {
		for j := 0; j < len(result_product_Chart_view)-1; j++ {
			if result_product_Chart_view[j].ProductID > result_product_Chart_view[j+1].ProductID {
				temp := result_product_Chart_view[j]
				result_product_Chart_view[j] = result_product_Chart_view[j+1]
				result_product_Chart_view[j+1] = temp
			}
		}
	}

	result := model.Chart_View{
		Id:         chart.ID,
		Username:   chart.Username,
		TotalQty:   jumlah_barang,
		TotalPrice: chart.TotalHarga,
		Product:    result_product_Chart_view,
	}

	return result, nil
}

func (cs *customerService) UpdateProductFromCartPlus(c echo.Context, id int) (model.General_Product_Chart, error) {
	var product model.Produk
	err := cs.connection.Where("id = ?", id).First(&product).Error
	if err != nil {
		return model.General_Product_Chart{}, errors.New("produk tidak ditemukan")
	}

	var chart model.Keranjang
	err = cs.connection.Where("username = ? AND status = ?", middleware.ExtractTokenUsername(c), "Belum Checkout").First(&chart).Error
	if err != nil {
		return model.General_Product_Chart{}, errors.New("keranjang tidak ditemukan")
	}

	var productFromChart model.Produk_Keranjang
	err = cs.connection.Where("id_produk = ? AND id_keranjang = ?", product.ID, chart.ID).First(&productFromChart).Error
	if err != nil {
		return model.General_Product_Chart{}, errors.New("produk tidak ditemukan di keranjang")
	}

	productFromChart.JumlahProduk = productFromChart.JumlahProduk + 1
	productFromChart.TotalHarga = productFromChart.TotalHarga + product.Harga
	err = cs.connection.Where("id_produk = ? AND id_keranjang = ?", product.ID, chart.ID).Updates(&productFromChart).Error
	if err != nil {
		return model.General_Product_Chart{}, errors.New("gagal menambahkan produk ke keranjang")
	}

	chart.TotalHarga = chart.TotalHarga + product.Harga
	err = cs.connection.Save(&chart).Error
	if err != nil {
		return model.General_Product_Chart{}, errors.New("gagal update total harga keranjang")
	}

	return model.General_Product_Chart{
		ChartID:    productFromChart.IdKeranjang,
		ProductID:  productFromChart.IdProduk,
		Quantity:   productFromChart.JumlahProduk,
		TotalPrice: productFromChart.TotalHarga,
	}, nil
}

func (cs *customerService) UpdateProductFromCartMinus(c echo.Context, id int) (model.General_Product_Chart, error) {
	var product model.Produk
	err := cs.connection.Where("id = ?", id).First(&product).Error
	if err != nil {
		return model.General_Product_Chart{}, errors.New("produk tidak ditemukan")
	}

	var chart model.Keranjang
	err = cs.connection.Where("username = ? AND status = ?", middleware.ExtractTokenUsername(c), "Belum Checkout").First(&chart).Error
	if err != nil {
		return model.General_Product_Chart{}, errors.New("keranjang tidak ditemukan")
	}

	var productFromChart model.Produk_Keranjang
	err = cs.connection.Where("id_produk = ? AND id_keranjang = ?", product.ID, chart.ID).First(&productFromChart).Error
	if err != nil {
		return model.General_Product_Chart{}, errors.New("produk tidak ditemukan di keranjang")
	}

	productFromChart.JumlahProduk = productFromChart.JumlahProduk - 1
	productFromChart.TotalHarga = productFromChart.TotalHarga - product.Harga
	err = cs.connection.Where("id_produk = ? AND id_keranjang = ?", product.ID, chart.ID).Updates(&productFromChart).Error
	if err != nil {
		return model.General_Product_Chart{}, errors.New("gagal mengurangi produk dari keranjang")
	}

	chart.TotalHarga = chart.TotalHarga - product.Harga
	err = cs.connection.Save(&chart).Error
	if err != nil {
		return model.General_Product_Chart{}, errors.New("gagal update total harga keranjang")
	}

	if productFromChart.JumlahProduk == 0 {
		err = cs.connection.Where("id_produk = ? AND id_keranjang = ?", product.ID, chart.ID).Delete(&productFromChart).Error
		if err != nil {
			return model.General_Product_Chart{}, errors.New("gagal menghapus produk dari keranjang")
		}
		return model.General_Product_Chart{}, errors.New("produk berhasil dihapus dari keranjang")
	}

	return model.General_Product_Chart{
		ChartID:    productFromChart.IdKeranjang,
		ProductID:  productFromChart.IdProduk,
		Quantity:   productFromChart.JumlahProduk,
		TotalPrice: productFromChart.TotalHarga,
	}, nil
}

func (cs *customerService) DeleteProductFromCart(c echo.Context, id int) error {
	var product model.Produk
	err := cs.connection.Where("id = ?", id).First(&product).Error
	if err != nil {
		return errors.New("produk tidak ditemukan")
	}

	var chart model.Keranjang
	err = cs.connection.Where("username = ? AND status = ?", middleware.ExtractTokenUsername(c), "Belum Checkout").First(&chart).Error
	if err != nil {
		return errors.New("keranjang tidak ditemukan")
	}

	var productFromChart model.Produk_Keranjang
	err = cs.connection.Where("id_produk = ? AND id_keranjang = ?", product.ID, chart.ID).First(&productFromChart).Error
	if err != nil {
		return errors.New("produk tidak ditemukan di keranjang")
	}

	chart.TotalHarga = chart.TotalHarga - productFromChart.TotalHarga
	err = cs.connection.Save(&chart).Error
	if err != nil {
		return errors.New("gagal update total harga keranjang")
	}

	err = cs.connection.Where("id_produk = ? AND id_keranjang = ?", product.ID, chart.ID).Delete(&productFromChart).Error
	if err != nil {
		return errors.New("gagal menghapus produk dari keranjang")
	}

	return nil
}

func (cs *customerService) Checkout(c echo.Context) (model.Chart_View, error) {
	chart, err := cs.GetCart(c)
	if err != nil {
		return model.Chart_View{}, errors.New("keranjang tidak ditemukan")
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
	checkout.Address = checkout_data.Address
	if checkout_data.Address == "" {
		checkout.Address = customer_data.Alamat
	}
	checkout.Courier = checkout_data.Courier
	checkout.ShippingCost = 0
	checkout.FinalPrice = chart.TotalPrice
	checkout.Chart = chart

	var pemesanan model.Pemesanan
	pemesanan.IdKeranjang = chart.Id
	pemesanan.CustomerUsername = chart.Username
	pemesanan.JumlahBarang = chart.TotalQty
	pemesanan.TotalHarga = checkout.FinalPrice
	pemesanan.Status = "Belum Dibayar"
	pemesanan.Alamat = checkout.Address
	pemesanan.Kurir = checkout.Courier
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
	err := cs.connection.Where("customer_username = ? AND status = ?", middleware.ExtractTokenUsername(c), "Belum Dibayar").Find(&pemesanan).Error
	if err != nil {
		return errors.New("pemesanan tidak ditemukan")
	}
	if pemesanan.IdKeranjang == 0 && pemesanan.CustomerUsername == "" {
		return errors.New("pemesanan tidak ditemukan")
	}

	pemesanan.BuktiPembayaran = payment_data.ProofOfPayment
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

	var OrdersDomain []model.Order_View
	for _, order := range pemesanan {
		var orderDomain model.Order_View
		orderDomain.Id = order.ID
		orderDomain.CreatedAt = order.CreatedAt
		orderDomain.CartID = order.IdKeranjang
		orderDomain.CustomerUsername = order.CustomerUsername
		orderDomain.TotalQty = order.JumlahBarang
		orderDomain.TotalPrice = order.TotalHarga
		orderDomain.Status = order.Status
		orderDomain.Address = order.Alamat
		orderDomain.Courier = order.Kurir
		orderDomain.ProofOfPayment = order.BuktiPembayaran
		orderDomain.ValidatedBy = order.DiValidasiOleh
		OrdersDomain = append(OrdersDomain, orderDomain)
	}

	return model.History_View{Order: OrdersDomain}, nil
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
	result.OrderID = pemesanan.ID
	result.CartID = keranjang.ID
	result.Status = pemesanan.Status
	result.Address = pemesanan.Alamat
	result.Courier = pemesanan.Kurir

	var result_product_view []model.Product_View
	for _, produk_data := range produk {
		result_product_view = append(result_product_view, model.Product_View{
			Id:          produk_data.ID,
			Name:        produk_data.Nama,
			Image:       produk_data.Gambar,
			Description: produk_data.Deskripsi,
			Price:       produk_data.Harga,
		})
	}

	var result_produk_pemesanan_view []model.Product_Order_View
	for i, produk_keranjang := range produk_keranjang {
		result_produk_pemesanan_view = append(result_produk_pemesanan_view, model.Product_Order_View{
			TotalQty:   produk_keranjang.JumlahProduk,
			TotalPrice: produk_keranjang.TotalHarga,
			Product:    result_product_view[i],
		})
	}

	result.Product = result_produk_pemesanan_view
	return result, nil
}

func (cs *customerService) CreateFeedbackPemesanan(c echo.Context, id uint) error {
	var pemesanan model.Pemesanan
	err := cs.connection.Where("id = ? AND customer_username = ?", id, middleware.ExtractTokenUsername(c)).First(&pemesanan).Error
	if err != nil {
		return errors.New("pemesanan tidak ditemukan")
	}

	var feedback_Id []model.Feedback
	err = cs.connection.Find(&feedback_Id).Error
	if err != nil {
		return errors.New("gagal mendapatkan id feedback")
	}

	var feedback model.Feedback_Pemesanan
	if len(feedback_Id) == 0 {
		feedback.IdFeedback = 1
	} else {
		feedback.IdFeedback = feedback_Id[len(feedback_Id)-1].ID + 1
	}
	feedback.IdPemesanan = pemesanan.ID
	feedback.Username = middleware.ExtractTokenUsername(c)
	feedback.Tanggal = time.Now()
	err = cs.connection.Create(&feedback).Error
	if err != nil {
		return errors.New("gagal membuat feedback pemesanan")
	}
	return nil
}

func (cs *customerService) PostFeedback(c echo.Context, feedback_data model.Feedback_Binding) (model.Feedback_Full_View, error) {
	feedbackDomain := model.Feedback{
		IdProduk:    feedback_data.ProductId,
		IsiFeedback: feedback_data.Feedback,
		Rating:      feedback_data.Rating,
	}

	err := cs.connection.Create(&feedbackDomain).Error
	if err != nil {
		return model.Feedback_Full_View{}, errors.New("gagal membuat feedback")
	}

	feedback_view := model.Feedback_View{
		Id:        feedbackDomain.ID,
		CreatedAt: feedbackDomain.CreatedAt,
		ProductID: feedbackDomain.IdProduk,
		Comment:   feedbackDomain.IsiFeedback,
		Rating:    feedbackDomain.Rating,
	}

	var result model.Feedback_Full_View
	result.Username = middleware.ExtractTokenUsername(c)
	result.ProductId = feedbackDomain.IdProduk
	result.Feedback = feedback_view

	return result, nil
}
