package controller

import (
	"Mini-Project_SiBaBe/dto"
	"Mini-Project_SiBaBe/model"
	"Mini-Project_SiBaBe/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

var adminService service.AdminSvc

func NewAdminController(service service.AdminSvc) {
	adminService = service
}

func LoginAdmin(c echo.Context) error {
	admin := model.Admin{}
	if err := c.Bind(&admin); err != nil {
		return c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Failed to process request", err))
	}

	result, err := adminService.LoginAdmin(c, admin)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to login admin", err))
	}
	return c.JSON(http.StatusOK, dto.BuildResponse("Success login admin", result))
}

func CreateProduct(c echo.Context) error {
	product := model.Produk{}
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Failed to process request", err))
	}

	result, err := adminService.CreateProduct(c, product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to create product", err))
	}
	return c.JSON(http.StatusCreated, dto.BuildResponse("Success create product", result))
}

func GetAllProductAdmin(c echo.Context) error {
	result, err := adminService.GetAllProduct(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to get all product", err))
	}
	return c.JSON(http.StatusOK, dto.BuildResponse("Success get all product", result))
}

func UpdateProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Failed to process request", err))
	}

	product := model.Produk{}
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Failed to process request", err))
	}

	result, err := adminService.UpdateProduct(c, id, product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to update product", err))
	}
	return c.JSON(http.StatusOK, dto.BuildResponse("Success update product", result))
}

func DeleteProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Failed to process request", err))
	}

	result, err := adminService.DeleteProduct(c, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to delete product", err))
	}
	return c.JSON(http.StatusOK, dto.BuildResponse("Success delete product", result))
}

func GetMonthlyReport(c echo.Context) error {
	result, err := adminService.GetMonthlyReport(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to get monthly report", err))
	}
	return c.JSON(http.StatusOK, dto.BuildResponse("Success get monthly report", result))
}

func GetYearlyReport(c echo.Context) error {
	result, err := adminService.GetYearlyReport(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to get yearly report", err))
	}
	return c.JSON(http.StatusOK, dto.BuildResponse("Success get yearly report", result))
}

func CreatePrduction(c echo.Context) error {
	production := model.Produksi{}
	if err := c.Bind(&production); err != nil {
		return c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Failed to process request", err))
	}

	result, err := adminService.CreateProduction(c, production)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to create production", err))
	}
	return c.JSON(http.StatusCreated, dto.BuildResponse("Success create production", result))
}

func GetOrderList(c echo.Context) error {
	result, err := adminService.GetOrderList(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to get order list", err))
	}
	return c.JSON(http.StatusOK, dto.BuildResponse("Success get order list", result))
}

func UpdateOrderStatus(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Failed to process request", err))
	}

	var status bool
	status, err = strconv.ParseBool(c.Param("status"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Failed to process request", err))
	}

	result, err := adminService.UpdateOrderStatus(c, id, status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to update order status", err))
	}
	return c.JSON(http.StatusOK, dto.BuildResponse("Success update order status", result))
}