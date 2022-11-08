package testing

import (
	"Mini-Project_SiBaBe/config"
	"Mini-Project_SiBaBe/controller"
	"Mini-Project_SiBaBe/model"
	"Mini-Project_SiBaBe/service"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

func InitCustomerTestAPI() *echo.Echo {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	db := config.InitDatabaseTest()
	Svc := service.NewCustomerService(db)
	controller.NewCustomerController(Svc)
	e := echo.New()
	return e
}

func TestCreateUser(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
		SizeData     int
	}{
		{
			Name:         "sucessfully create user",
			Path:         "/customer/register",
			ExpectedCode: http.StatusCreated,
			SizeData:     1,
		},
	}

	e := InitCustomerTestAPI()
	userJSON := `{"nama":"test", "password" : "testpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/customer/register", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.CreateUser(c)) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
				var user []model.Customer
				model.DB.Find(&user)
				assert.Equal(t, tc.SizeData, len(user))
			}
		})
	}
}

func InsertDataProduct() {
	product := model.Produk{
		Nama:      "test",
		Gambar:    "test link",
		Stok:      10,
		Deskripsi: "test",
		Harga:     10000,
	}
	err := model.DB.Save(&product).Error
	if err != nil {
		panic(err)
	}
}

func TestGetAllProduct(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
		SizeData     int
	}{
		{
			Name:         "sucessfully get all product",
			Path:         "/customer/products",
			ExpectedCode: http.StatusOK,
			SizeData:     1,
		},
	}

	e := InitCustomerTestAPI()
	InsertDataProduct()
	req := httptest.NewRequest(http.MethodGet, "/customer/products", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.GetAllProduct(c)) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
				var products []model.Produk
				model.DB.Find(&products)
				assert.Equal(t, tc.SizeData, len(products))
			}
		})
	}
}

func TestGetProductByID(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
		SizeData     int
	}{
		{
			Name:         "sucessfully get product by id",
			Path:         "/customer/products/1",
			ExpectedCode: http.StatusOK,
			SizeData:     1,
		},
	}

	e := InitCustomerTestAPI()
	InsertDataProduct()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/customer/products/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.GetProductById(c)) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
				var product []model.Produk
				model.DB.Find(&product)
				assert.Equal(t, tc.SizeData, len(product))
			}
		})
	}
}

func InsertDataCustomer() {
	customer := model.Customer{
		Username: "test",
		Password: "testpassword",
	}
	err := model.DB.Create(&customer).Error
	if err != nil {
		panic(err)
	}
}

func TestLoginUser(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
		SizeData     int
	}{
		{
			Name:         "sucessfully login user",
			Path:         "/customer/login",
			ExpectedCode: http.StatusOK,
			SizeData:     1,
		},
	}

	e := InitCustomerTestAPI()
	InsertDataCustomer()
	userJSON := `{"username":"test", "password" : "testpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/customer/login", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.LoginUser(c)) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
				var user []model.Customer
				model.DB.Find(&user)
				assert.Equal(t, tc.SizeData, len(user))
			}
		})
	}
}

func InsertDataCart() {
	cart := model.Keranjang{
		Username:   "test",
		TotalHarga: 10000,
		Status:     "Belum Checkout",
	}
	err := model.DB.Create(&cart).Error
	if err != nil {
		panic(err)
	}
}

func InsertDataProductToCart() {
	cart_product := model.Produk_Keranjang{
		IdKeranjang:  1,
		IdProduk:     1,
		JumlahProduk: 1,
		TotalHarga:   10000,
	}
	err := model.DB.Create(&cart_product).Error
	if err != nil {
		panic(err)
	}
}

// kendala mulai dari sini
func TestPostProductToCart(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
		SizeData     int
	}{
		{
			Name:         "sucessfully post product to cart",
			Path:         "/customer/jwt/product-to-cart/1",
			ExpectedCode: http.StatusOK,
			SizeData:     1,
		},
	}

	e := InitCustomerTestAPI()
	InsertDataCustomer()
	InsertDataProduct()
	InsertDataCart()
	InsertDataProductToCart()
	e.Use(middleware.JWT([]byte(os.Getenv("JWT_SECRET"))))
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1hIjoiIiwidXNlcm5hbWUiOiJ0ZXN0In0.MBgzq4cnyZ9w-JC_Xji2Hss_-IEWbe-XbEI9Cg_qVT0"
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/customer/jwt/product-to-cart/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.PostProductToCart(c)) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
				var cart []model.Produk_Keranjang
				model.DB.Find(&cart)
				assert.Equal(t, tc.SizeData, len(cart))
			}
		})
	}
}
