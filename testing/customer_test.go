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

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	mid "github.com/labstack/echo/v4/middleware"
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

func TestGetAllProductCustomer(t *testing.T) {
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
	InsertDataProductToCart()
	e.GET("/customer/jwt/product-to-cart/1",
		func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			return c.JSON(http.StatusOK, token.Claims)
		})
	e.Use(mid.JWT([]byte(os.Getenv("JWT_SECRET"))))
	req := httptest.NewRequest(http.MethodGet, "/customer/jwt/product-to-cart/1", nil)
	req.Header.Set(echo.HeaderAuthorization, "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1hIjoiIiwidXNlcm5hbWUiOiJ0ZXN0In0.MBgzq4cnyZ9w-JC_Xji2Hss_-IEWbe-XbEI9Cg_qVT0")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.PostProductToCart(e.AcquireContext())) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
				var cart []model.Produk_Keranjang
				model.DB.Find(&cart)
				assert.Equal(t, tc.SizeData, len(cart))
			}
		})
	}
}

func TestGetCart(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
		SizeData     int
	}{
		{
			Name:         "sucessfully get cart",
			Path:         "/customer/jwt/cart",
			ExpectedCode: http.StatusOK,
			SizeData:     1,
		},
	}

	e := InitCustomerTestAPI()
	InsertDataCart()
	e.GET("/customer/jwt/cart",
		func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			return c.JSON(http.StatusOK, token.Claims)
		})
	e.Use(mid.JWT([]byte(os.Getenv("JWT_SECRET"))))
	req := httptest.NewRequest(http.MethodGet, "/customer/jwt/cart", nil)
	req.Header.Set(echo.HeaderAuthorization, "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1hIjoiIiwidXNlcm5hbWUiOiJ0ZXN0In0.MBgzq4cnyZ9w-JC_Xji2Hss_-IEWbe-XbEI9Cg_qVT0")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.GetCart(e.AcquireContext())) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
				var cart []model.Keranjang
				model.DB.Find(&cart)
				assert.Equal(t, tc.SizeData, len(cart))
			}
		})
	}
}

func TestUpdateProductFromCartPlus(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
		SizeData     int
	}{
		{
			Name:         "sucessfully update product from cart plus",
			Path:         "/customer/jwt/cart/plus/1",
			ExpectedCode: http.StatusOK,
			SizeData:     1,
		},
	}

	e := InitCustomerTestAPI()
	InsertDataProductToCart()
	e.GET("/customer/jwt/cart/plus/1",
		func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			return c.JSON(http.StatusOK, token.Claims)
		})
	e.Use(mid.JWT([]byte(os.Getenv("JWT_SECRET"))))
	req := httptest.NewRequest(http.MethodGet, "/customer/jwt/cart/plus/1", nil)
	req.Header.Set(echo.HeaderAuthorization, "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1hIjoiIiwidXNlcm5hbWUiOiJ0ZXN0In0.MBgzq4cnyZ9w-JC_Xji2Hss_-IEWbe-XbEI9Cg_qVT0")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.UpdateProductFromCartPlus(e.AcquireContext())) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
				var cart []model.Produk_Keranjang
				model.DB.Find(&cart)
				assert.Equal(t, tc.SizeData, len(cart))
			}
		})
	}
}

func TestUpdateProductFromCartMinus(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
		SizeData     int
	}{
		{
			Name:         "sucessfully update product from cart minus",
			Path:         "/customer/jwt/cart/minus/1",
			ExpectedCode: http.StatusOK,
			SizeData:     1,
		},
	}

	e := InitCustomerTestAPI()
	InsertDataProductToCart()
	e.GET("/customer/jwt/cart/minus/1",
		func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			return c.JSON(http.StatusOK, token.Claims)
		})
	e.Use(mid.JWT([]byte(os.Getenv("JWT_SECRET"))))
	req := httptest.NewRequest(http.MethodGet, "/customer/jwt/cart/minus/1", nil)
	req.Header.Set(echo.HeaderAuthorization, "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1hIjoiIiwidXNlcm5hbWUiOiJ0ZXN0In0.MBgzq4cnyZ9w-JC_Xji2Hss_-IEWbe-XbEI9Cg_qVT0")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.UpdateProductFromCartMinus(e.AcquireContext())) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
				var cart []model.Produk_Keranjang
				model.DB.Find(&cart)
				assert.Equal(t, tc.SizeData, len(cart))
			}
		})
	}
}

func TestCheckout(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
		SizeData     int
	}{
		{
			Name:         "sucessfully checkout",
			Path:         "/customer/jwt/cart/checkout",
			ExpectedCode: http.StatusOK,
			SizeData:     1,
		},
	}

	e := InitCustomerTestAPI()
	InsertDataProductToCart()
	e.GET("/customer/jwt/cart/checkout",
		func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			return c.JSON(http.StatusOK, token.Claims)
		})
	e.Use(mid.JWT([]byte(os.Getenv("JWT_SECRET"))))
	req := httptest.NewRequest(http.MethodGet, "/customer/jwt/cart/checkout", nil)
	req.Header.Set(echo.HeaderAuthorization, "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1hIjoiIiwidXNlcm5hbWUiOiJ0ZXN0In0.MBgzq4cnyZ9w-JC_Xji2Hss_-IEWbe-XbEI9Cg_qVT0")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.Checkout(e.AcquireContext())) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
				var cart []model.Produk_Keranjang
				model.DB.Find(&cart)
				assert.Equal(t, tc.SizeData, len(cart))
			}
		})
	}
}

func TestConfirmCheckout(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
	}{
		{
			Name:         "sucessfully get checkout confirm",
			Path:         "/customer/jwt/cart/checkout/confirm",
			ExpectedCode: http.StatusOK,
		},
	}

	e := InitCustomerTestAPI()
	JSON := `{"kurir":"test jek, "alamat":"Jl. Sukolilo No. 1, Surabaya"}`
	e.POST("/customer/jwt/cart/checkout/confirm",
		func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			return c.JSON(http.StatusOK, token.Claims)
		})
	e.Use(mid.JWT([]byte(os.Getenv("JWT_SECRET"))))
	req := httptest.NewRequest(http.MethodPost, "/customer/jwt/cart/checkout/confirm", strings.NewReader(JSON))
	req.Header.Set(echo.HeaderAuthorization, "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1hIjoiIiwidXNlcm5hbWUiOiJ0ZXN0In0.MBgzq4cnyZ9w-JC_Xji2Hss_-IEWbe-XbEI9Cg_qVT0")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.ConfirmCheckout(e.AcquireContext())) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
			}
		})
	}
}

func TestConfirmPayment(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
	}{
		{
			Name:         "sucessfully get payment confirm",
			Path:         "/customer/jwt/cart/checkout/payment/confirm",
			ExpectedCode: http.StatusOK,
		},
	}

	e := InitCustomerTestAPI()
	e.POST("/customer/jwt/cart/checkout/payment/confirm",
		func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			return c.JSON(http.StatusOK, token.Claims)
		})
	e.Use(mid.JWT([]byte(os.Getenv("JWT_SECRET"))))
	JSON := `{"bukti_pembayaran":"ini link bukti pembayaran test"}`
	req := httptest.NewRequest(http.MethodPost, "/customer/jwt/cart/checkout/payment/confirm", strings.NewReader(JSON))
	req.Header.Set(echo.HeaderAuthorization, "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1hIjoiIiwidXNlcm5hbWUiOiJ0ZXN0In0.MBgzq4cnyZ9w-JC_Xji2Hss_-IEWbe-XbEI9Cg_qVT0")
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.ConfirmPayment(e.AcquireContext())) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
			}
		})
	}
}

func TestGetHistory(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
	}{
		{
			Name:         "sucessfully get history",
			Path:         "/customer/jwt/history",
			ExpectedCode: http.StatusOK,
		},
	}

	e := InitCustomerTestAPI()
	e.GET("/customer/jwt/history",
		func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			return c.JSON(http.StatusOK, token.Claims)
		})
	e.Use(mid.JWT([]byte(os.Getenv("JWT_SECRET"))))
	req := httptest.NewRequest(http.MethodGet, "/customer/jwt/history", nil)
	req.Header.Set(echo.HeaderAuthorization, "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1hIjoiIiwidXNlcm5hbWUiOiJ0ZXN0In0.MBgzq4cnyZ9w-JC_Xji2Hss_-IEWbe-XbEI9Cg_qVT0")
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.GetHistory(e.AcquireContext())) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
			}
		})
	}
}

func TestGetHistoryDetail(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
	}{
		{
			Name:         "sucessfully get history detail",
			Path:         "/customer/jwt/history/1",
			ExpectedCode: http.StatusOK,
		},
	}

	e := InitCustomerTestAPI()
	e.GET("/customer/jwt/history/1",
		func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			return c.JSON(http.StatusOK, token.Claims)
		})
	e.Use(mid.JWT([]byte(os.Getenv("JWT_SECRET"))))
	req := httptest.NewRequest(http.MethodGet, "/customer/jwt/history/1", nil)
	req.Header.Set(echo.HeaderAuthorization, "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1hIjoiIiwidXNlcm5hbWUiOiJ0ZXN0In0.MBgzq4cnyZ9w-JC_Xji2Hss_-IEWbe-XbEI9Cg_qVT0")
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.GetHistoryDetail(e.AcquireContext())) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
			}
		})
	}
}

func TestPostFeedback(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
	}{
		{
			Name:         "sucessfully post feedback to product 1",
			Path:         "/customer/jwt/History/1/feedback/1",
			ExpectedCode: http.StatusOK,
		},
	}

	e := InitCustomerTestAPI()
	e.POST("/customer/jwt/History/1/feedback/1",
		func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			return c.JSON(http.StatusOK, token.Claims)
		})
	e.Use(mid.JWT([]byte(os.Getenv("JWT_SECRET"))))
	JSON := `{"isi_feedback":"feedback test for product 1","rating":5}`
	req := httptest.NewRequest(http.MethodPost, "/customer/jwt/History/1/feedback/1", strings.NewReader(JSON))
	req.Header.Set(echo.HeaderAuthorization, "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1hIjoiIiwidXNlcm5hbWUiOiJ0ZXN0In0.MBgzq4cnyZ9w-JC_Xji2Hss_-IEWbe-XbEI9Cg_qVT0")
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.PostFeedback(e.AcquireContext())) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
			}
		})
	}
}
