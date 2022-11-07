package controller

import (
	"Mini-Project_SiBaBe/dto"
	"Mini-Project_SiBaBe/model"
	"Mini-Project_SiBaBe/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

var customerService service.CustomerSvc

func NewCustomerController(service service.CustomerSvc) {
	customerService = service
}

func CreateUser(c echo.Context) error {
	user := model.Customer{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Failed to process request", err))
	}

	result, err := customerService.CreateUser(c, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to create user", err))
	}
	return c.JSON(http.StatusCreated, dto.BuildResponse("Success create user", result))
}

func GetAllProduct(c echo.Context) error {
	result, err := customerService.GetAllProduct(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to get all product", err))
	}
	return c.JSON(http.StatusOK, dto.BuildResponse("Success get all product", result))
}

func GetProductById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Failed to process request", err))
	}

	result, err := customerService.GetProductById(c, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to get product by id", err))
	}
	return c.JSON(http.StatusOK, dto.BuildResponse("Success get product by id", result))
}

func LoginUser(c echo.Context) error {
	user := model.Customer{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Failed to process request", err))
	}

	result, err := customerService.LoginUser(c, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to login user", err))
	}
	return c.JSON(http.StatusOK, dto.BuildResponse("Success login user", result))
}

// function with customer authentication
func PostProductToCart(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Failed to process request", err))
	}

	result, err := customerService.PostProductToCart(c, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to post product to cart", err))
	}
	return c.JSON(http.StatusOK, dto.BuildResponse("Success post product to cart", result))
}

func GetCart(c echo.Context) error {
	result, err := customerService.GetCart(c)
	if err != nil {
		if err.Error() == "tidak ada barang di keranjang" {
			return c.JSON(http.StatusOK, dto.BuildResponse("Success get cart", err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to get cart", err))
	}
	return c.JSON(http.StatusOK, dto.BuildResponse("Success get cart", result))
}

func UpdateProductFromCartPlus(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Failed to process request", err))
	}

	result, err := customerService.UpdateProductFromCartPlus(c, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to update product from cart plus", err))
	}
	return c.JSON(http.StatusOK, dto.BuildResponse("Success update product from cart plus", result))
}

func UpdateProductFromCartMinus(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Failed to process request", err))
	}

	result, err := customerService.UpdateProductFromCartMinus(c, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to update product from cart minus", err))
	}
	return c.JSON(http.StatusOK, dto.BuildResponse("Success update product from cart minus", result))
}

func Checkout(c echo.Context) error {
	result, err := customerService.Checkout(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to checkout", err))
	}
	return c.JSON(http.StatusOK, dto.BuildResponse("Success checkout", result))
}

func ConfirmCheckout(c echo.Context) error {
	var checkout model.Checkout_Binding
	if err := c.Bind(&checkout); err != nil {
		return c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Failed to process request", err))
	}
	result, err := customerService.ConfirmCheckout(c, checkout)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to confirm checkout", err))
	}
	return c.JSON(http.StatusOK, dto.BuildResponse("Success confirm checkout", result))
}

func ConfirmPayment(c echo.Context) error {
	var payment model.Payment_Binding
	if err := c.Bind(&payment); err != nil {
		return c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Failed to process request", err))
	}
	err := customerService.ConfirmPayment(c, payment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to confirm payment", err))
	}
	return c.JSON(http.StatusOK, dto.BuildResponse("Success send payment confirmation", nil))
}

func GetHistory(c echo.Context) error {
	result, err := customerService.GetHistory(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to get history", err))
	}
	return c.JSON(http.StatusOK, dto.BuildResponse("Success get history", result))
}

func GetHistoryDetail(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Failed to process request", err))
	}

	result, err := customerService.GetHistoryDetail(c, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to get history by id", err))
	}
	return c.JSON(http.StatusOK, dto.BuildResponse("Success get history by id", result))
}

func PostFeedback(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Failed to process request", err))
	}

	err = customerService.CreatFeedbackPemesanan(c, uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to create feedback", err))
	}

	id_produk, err := strconv.Atoi(c.Param("id_produk"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Failed to process request", err))
	}

	var feedback model.Feedback
	if err := c.Bind(&feedback); err != nil {
		return c.JSON(http.StatusBadRequest, dto.BuildErrorResponse("Failed to process request", err))
	}

	feedback.IdProduk = uint(id_produk)
	var result model.Feedback_View
	result, err = customerService.PostFeedback(c, feedback)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.BuildErrorResponse("Failed to post feedback", err))
	}
	return c.JSON(http.StatusOK, dto.BuildResponse("Success post feedback", result))
}
