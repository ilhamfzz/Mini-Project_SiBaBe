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

func InitAdminTestAPI() *echo.Echo {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	db := config.InitDatabaseTest()
	Svc := service.NewAdminService(db)
	controller.NewAdminController(Svc)
	e := echo.New()
	return e
}

func InsertDataAdmin() {
	Admin := model.Admin{
		Username: "admin",
		Password: "adminpassword",
		Nama:     "admin",
	}
	err := model.DB.Create(&Admin).Error
	if err != nil {
		panic(err)
	}
}

func TestAdminLogin(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
		SizeData     int
	}{
		{
			Name:         "sucessfully create user",
			Path:         "/customer/register",
			ExpectedCode: http.StatusOK,
			SizeData:     1,
		},
	}

	e := InitAdminTestAPI()
	InsertDataAdmin()
	AdminJSON := `{"username":"admin","password":"adminpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(AdminJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.LoginAdmin(c)) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
				var admin []model.Admin
				model.DB.Find(&admin)
				assert.Equal(t, tc.SizeData, len(admin))
			}
		})
	}
}

func TestCreateProduct(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
	}{
		{
			Name:         "sucessfully create product",
			Path:         "/admin/jwt/product",
			ExpectedCode: http.StatusOK,
		},
	}

	e := InitAdminTestAPI()
	InsertDataAdmin()
	e.POST("/admin/jwt/product",
		func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			return c.JSON(http.StatusOK, token.Claims)
		})
	e.Use(mid.JWT([]byte(os.Getenv("JWT_SECRET"))))
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1hIjoiYWRtaW4iLCJ1c2VybmFtZSI6ImFkbWluIn0.I-y1F2xt6BYhBXO-qZpHrMdW0KZeBa4A8tdUG5kNlGg"
	req := httptest.NewRequest(http.MethodPost, "/admin/jwt/product", nil)
	req.Header.Set(echo.HeaderAuthorization, "bearer "+token)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.CreateProduct(e.AcquireContext())) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
			}
		})
	}
}

func TestGetAllProductAdmin(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
		SizeData     int
	}{
		{
			Name:         "sucessfully get all product",
			Path:         "/admin/jwt/products",
			ExpectedCode: http.StatusOK,
			SizeData:     1,
		},
	}

	e := InitAdminTestAPI()
	InsertDataProduct()
	e.GET("/admin/jwt/products",
		func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			return c.JSON(http.StatusOK, token.Claims)
		})
	e.Use(mid.JWT([]byte(os.Getenv("JWT_SECRET"))))
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1hIjoiYWRtaW4iLCJ1c2VybmFtZSI6ImFkbWluIn0.I-y1F2xt6BYhBXO-qZpHrMdW0KZeBa4A8tdUG5kNlGg"
	req := httptest.NewRequest(http.MethodGet, "/admin/jwt/products", nil)
	req.Header.Set(echo.HeaderAuthorization, "bearer "+token)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.GetAllProductAdmin(e.AcquireContext())) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
				var product []model.Produk
				model.DB.Find(&product)
				assert.Equal(t, tc.SizeData, len(product))
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
		SizeData     int
	}{
		{
			Name:         "sucessfully update product",
			Path:         "/admin/jwt/product/1",
			ExpectedCode: http.StatusOK,
			SizeData:     1,
		},
	}

	e := InitAdminTestAPI()
	InsertDataProduct()
	e.PUT("/admin/jwt/product/1",
		func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			return c.JSON(http.StatusOK, token.Claims)
		})
	e.Use(mid.JWT([]byte(os.Getenv("JWT_SECRET"))))
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1hIjoiYWRtaW4iLCJ1c2VybmFtZSI6ImFkbWluIn0.I-y1F2xt6BYhBXO-qZpHrMdW0KZeBa4A8tdUG5kNlGg"
	JSON := `{"nama":"test","deskripsi":"test","harga":10000,"stok":100,"gambar":"test"}`
	req := httptest.NewRequest(http.MethodPut, "/admin/jwt/product/1", strings.NewReader(JSON))
	req.Header.Set(echo.HeaderAuthorization, "bearer "+token)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.UpdateProduct(e.AcquireContext())) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
				var product []model.Produk
				model.DB.Find(&product)
				assert.Equal(t, tc.SizeData, len(product))
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
	}{
		{
			Name:         "sucessfully delete product",
			Path:         "/admin/jwt/product/1",
			ExpectedCode: http.StatusOK,
		},
	}

	e := InitAdminTestAPI()
	InsertDataProduct()
	e.DELETE("/admin/jwt/product/1",
		func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			return c.JSON(http.StatusOK, token.Claims)
		})
	e.Use(mid.JWT([]byte(os.Getenv("JWT_SECRET"))))
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1hIjoiYWRtaW4iLCJ1c2VybmFtZSI6ImFkbWluIn0.I-y1F2xt6BYhBXO-qZpHrMdW0KZeBa4A8tdUG5kNlGg"
	req := httptest.NewRequest(http.MethodDelete, "/admin/jwt/product/1", nil)
	req.Header.Set(echo.HeaderAuthorization, "bearer "+token)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.DeleteProduct(e.AcquireContext())) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
			}
		})
	}
}

func TestGetMonthlyReport(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
	}{
		{
			Name:         "sucessfully get monthly report",
			Path:         "/admin/jwt/report/monthly",
			ExpectedCode: http.StatusOK,
		},
	}

	e := InitAdminTestAPI()
	e.GET("/admin/jwt/report/monthly",
		func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			return c.JSON(http.StatusOK, token.Claims)
		})
	e.Use(mid.JWT([]byte(os.Getenv("JWT_SECRET"))))
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1hIjoiYWRtaW4iLCJ1c2VybmFtZSI6ImFkbWluIn0.I-y1F2xt6BYhBXO-qZpHrMdW0KZeBa4A8tdUG5kNlGg"
	req := httptest.NewRequest(http.MethodGet, "/admin/jwt/report/monthly", nil)
	req.Header.Set(echo.HeaderAuthorization, "bearer "+token)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.GetMonthlyReport(e.AcquireContext())) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
			}
		})
	}
}

func TestCreateProduction(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
		SizeData     int
	}{
		{
			Name:         "sucessfully create production",
			Path:         "/admin/jwt/production",
			ExpectedCode: http.StatusOK,
			SizeData:     1,
		},
	}

	e := InitAdminTestAPI()
	e.POST("/admin/jwt/production",
		func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			return c.JSON(http.StatusOK, token.Claims)
		})
	e.Use(mid.JWT([]byte(os.Getenv("JWT_SECRET"))))
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1hIjoiYWRtaW4iLCJ1c2VybmFtZSI6ImFkbWluIn0.I-y1F2xt6BYhBXO-qZpHrMdW0KZeBa4A8tdUG5kNlGg"
	JSON := `{"nama_produk":"kiw produk nih","total_biaya":100000}`
	req := httptest.NewRequest(http.MethodPost, "/admin/jwt/production", strings.NewReader(JSON))
	req.Header.Set(echo.HeaderAuthorization, "bearer "+token)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.CreateProduction(e.AcquireContext())) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
				var production []model.Produksi
				model.DB.Find(&production)
				assert.Equal(t, tc.SizeData, len(production))
			}
		})
	}
}

func TestGetOrderList(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
	}{
		{
			Name:         "sucessfully get order list",
			Path:         "/admin/jwt/order-list",
			ExpectedCode: http.StatusOK,
		},
	}

	e := InitAdminTestAPI()
	e.GET("/admin/jwt/order-list",
		func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			return c.JSON(http.StatusOK, token.Claims)
		})
	e.Use(mid.JWT([]byte(os.Getenv("JWT_SECRET"))))
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1hIjoiYWRtaW4iLCJ1c2VybmFtZSI6ImFkbWluIn0.I-y1F2xt6BYhBXO-qZpHrMdW0KZeBa4A8tdUG5kNlGg"
	req := httptest.NewRequest(http.MethodGet, "/admin/jwt/order-list", nil)
	req.Header.Set(echo.HeaderAuthorization, "bearer "+token)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.GetOrderList(e.AcquireContext())) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
			}
		})
	}
}

func TestHandleOrderPayment(t *testing.T) {
	testCase := []struct {
		Name         string
		Path         string
		ExpectedCode int
	}{
		{
			Name:         "sucessfully handle order payment",
			Path:         "/admin/jwt/order-list/1",
			ExpectedCode: http.StatusOK,
		},
	}

	e := InitAdminTestAPI()
	e.POST("/admin/jwt/order-list/1",
		func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			return c.JSON(http.StatusOK, token.Claims)
		})
	e.Use(mid.JWT([]byte(os.Getenv("JWT_SECRET"))))
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1hIjoiYWRtaW4iLCJ1c2VybmFtZSI6ImFkbWluIn0.I-y1F2xt6BYhBXO-qZpHrMdW0KZeBa4A8tdUG5kNlGg"
	JSON := `{"status":"Terima"}`
	req := httptest.NewRequest(http.MethodPost, "/admin/jwt/order-list/1", strings.NewReader(JSON))
	req.Header.Set(echo.HeaderAuthorization, "bearer "+token)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			if assert.NoError(t, controller.UpdateOrderStatus(e.AcquireContext())) {
				assert.Equal(t, tc.ExpectedCode, rec.Code)
			}
		})
	}
}
