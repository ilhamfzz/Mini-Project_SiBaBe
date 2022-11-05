package service

import (
	"Mini-Project_SiBaBe/dto"
	"Mini-Project_SiBaBe/middleware"
	"Mini-Project_SiBaBe/model"
	"errors"

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
		return user, errors.New("Username sudah terdaftar")
	}
	return user, nil
}

func (cs *customerService) GetAllProduct(c echo.Context) ([]model.Produk, error) {
	var products []model.Produk
	err := cs.connection.Find(&products).Error
	if err != nil {
		return products, errors.New("Tidak ada produk")
	}
	return products, nil
}

func (cs *customerService) GetProductById(c echo.Context, id int) (model.Produk, error) {
	var product model.Produk
	err := cs.connection.Where("id = ?", id).First(&product).Error
	if err != nil {
		return product, errors.New("Produk tidak ditemukan")
	}
	return product, nil
}

func (cs *customerService) LoginUser(c echo.Context, user model.Customer) (dto.Login, error) {
	var userLogin model.Customer
	err := cs.connection.Where("username = ? AND password = ?", user.Username, user.Password).First(&userLogin).Error
	if err != nil {
		return dto.Login{}, errors.New("Username atau password salah")
	}

	var token string
	token, err = middleware.CreateToken(userLogin.Username, userLogin.Nama)
	if err != nil {
		return dto.Login{}, errors.New("Gagal membuat token")
	}
	result := dto.Login{
		Username: userLogin.Username,
		Nama:     userLogin.Nama,
		Token:    token,
	}
	return result, nil
}

func (cs *customerService) PostProductToCart(c echo.Context, id int) (model.Produk_Keranjang, error) {
	var product model.Produk
	err := cs.connection.Where("id = ?", id).First(&product).Error
	if err != nil {
		return model.Produk_Keranjang{}, errors.New("Produk tidak ditemukan")
	}

	var chart model.Keranjang
	chart.Username, err = middleware.ExtractTokenUsername(c)
	if err != nil {
		return model.Produk_Keranjang{}, errors.New("Gagal mendapatkan username")
	}
	var flag bool = true
	err = cs.connection.Where("username = ? AND status = ?", chart.Username, true).First(&chart).Error
	if err != nil {
		flag = false
		chart.TotalHarga = 0
		chart.Status = true
		err = cs.connection.Create(&chart).Error
		if err != nil {
			return model.Produk_Keranjang{}, errors.New("Gagal membuat keranjang")
		}

		err = cs.connection.Where("username = ? AND status = ?", chart.Username, true).First(&chart).Error
		if err != nil {
			return model.Produk_Keranjang{}, errors.New("Gagal mendapatkan keranjang")
		}
	}

	var productFromChart model.Produk_Keranjang
	if flag {
		err = cs.connection.Where("id_produk = ? AND id_keranjang = ?", product.ID, chart.ID).First(&productFromChart).Error
		if err != nil {
			return model.Produk_Keranjang{}, errors.New("Gagal menemukan produk pada keranjang")
		}
		productFromChart.JumlahProduk += 1
		productFromChart.TotalHarga += product.Harga
		err = cs.connection.Where("id_produk = ? AND id_keranjang = ?", product.ID, chart.ID).Updates(&productFromChart).Error
		if err != nil {
			return model.Produk_Keranjang{}, errors.New("Gagal menambahkan produk ke keranjang")
		}
		chart.TotalHarga += product.Harga
	} else {
		productFromChart.IdKeranjang = chart.ID
		productFromChart.IdProduk = product.ID
		productFromChart.JumlahProduk = 1
		productFromChart.TotalHarga = product.Harga
		err = cs.connection.Create(&productFromChart).Error
		if err != nil {
			return model.Produk_Keranjang{}, errors.New("Gagal menambahkan produk ke keranjang")
		}
		chart.TotalHarga = productFromChart.TotalHarga
	}

	err = cs.connection.Where("id = ?", chart.ID).Updates(&chart).Error
	if err != nil {
		return model.Produk_Keranjang{}, errors.New("Gagal mengupdate total harga keranjang")
	}
	return productFromChart, nil
}

func (cs *customerService) GetCart(c echo.Context) ([]model.Produk_Keranjang, error) {
	var chart model.Keranjang
	var err error
	chart.Username, err = middleware.ExtractTokenUsername(c)
	if err != nil {
		return []model.Produk_Keranjang{}, errors.New("Gagal mendapatkan username")
	}
	err = cs.connection.Where("username = ? AND status = ?", chart.Username, true).First(&chart).Error
	if err != nil {
		return []model.Produk_Keranjang{}, errors.New("Gagal Menemukan Keranjang")
	}

	var products []model.Produk_Keranjang
	err = cs.connection.Where("id_keranjang = ?", chart.ID).Find(&products).Error
	if err != nil {
		return []model.Produk_Keranjang{}, errors.New("Gagal mendapatkan produk pada keranjang")
	}
	return products, nil
}

func (cs *customerService) UpdateProductFromCartPlus(c echo.Context, id int) (model.Produk_Keranjang, error) {
	var chart model.Keranjang
	var err error
	chart.Username, err = middleware.ExtractTokenUsername(c)
	if err != nil {
		return model.Produk_Keranjang{}, err
	}
	err = cs.connection.Where("username = ? AND status = ?", chart.Username, true).First(&chart).Error
	if err != nil {
		return model.Produk_Keranjang{}, err
	}

	var productsChart []model.Produk_Keranjang
	err = cs.connection.Where("id_keranjang = ?", chart.ID).First(&productsChart).Error
	if err != nil {
		return model.Produk_Keranjang{}, err
	}

	var productFromChart model.Produk_Keranjang
	for i := 0; i < len(productsChart); i++ {
		if productsChart[i].IdProduk == uint(id) {
			productFromChart = productsChart[i]
			break
		}
	}

	var productModel model.Produk
	err = cs.connection.Where("id = ?", productFromChart.IdProduk).First(&productModel).Error
	if err != nil {
		return model.Produk_Keranjang{}, err
	}

	productFromChart.JumlahProduk += 1
	productFromChart.TotalHarga += productModel.Harga
	err = cs.connection.Where("id_keranjang = ? AND id_produk = ?", chart.ID, productFromChart.IdProduk).Updates(&productFromChart).Error
	if err != nil {
		return model.Produk_Keranjang{}, err
	}

	chart.TotalHarga += productModel.Harga
	err = cs.connection.Where("id = ?", chart.ID).Updates(&chart).Error
	if err != nil {
		return model.Produk_Keranjang{}, err
	}
	return productFromChart, nil
}

func (cs *customerService) UpdateProductFromCartMinus(c echo.Context, id int) (model.Produk_Keranjang, error) {
	var chart model.Keranjang
	var err error
	chart.Username, err = middleware.ExtractTokenUsername(c)
	if err != nil {
		return model.Produk_Keranjang{}, err
	}
	err = cs.connection.Where("username = ? AND status = ?", chart.Username, true).First(&chart).Error
	if err != nil {
		return model.Produk_Keranjang{}, err
	}

	var productsChart []model.Produk_Keranjang
	err = cs.connection.Where("id_keranjang = ?", chart.ID).First(&productsChart).Error
	if err != nil {
		return model.Produk_Keranjang{}, err
	}

	var productFromChart model.Produk_Keranjang
	for i := 0; i < len(productsChart); i++ {
		if productsChart[i].IdProduk == uint(id) {
			productFromChart = productsChart[i]
			break
		}
	}

	var productModel model.Produk
	err = cs.connection.Where("id = ?", productFromChart.IdProduk).First(&productModel).Error
	if err != nil {
		return model.Produk_Keranjang{}, err
	}

	productFromChart.JumlahProduk -= 1
	productFromChart.TotalHarga -= productModel.Harga
	err = cs.connection.Where("id_keranjang = ? AND id_produk = ?", chart.ID, productFromChart.IdProduk).Updates(&productFromChart).Error
	if err != nil {
		return model.Produk_Keranjang{}, err
	}

	chart.TotalHarga -= productModel.Harga
	err = cs.connection.Where("id = ?", chart.ID).Updates(&chart).Error
	if err != nil {
		return model.Produk_Keranjang{}, err
	}
	return productFromChart, nil
}
