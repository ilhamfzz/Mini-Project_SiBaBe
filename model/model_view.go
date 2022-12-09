package model

import (
	"time"
)

type Order_List struct {
	OrderID   uint           `json:"orderId"`
	CartID    uint           `json:"cartId"`
	Status    string         `json:"status"`
	OrderList []Product_View `json:"orderList"`
}

type Checkout struct {
	Invoice      string    `json:"invoice"`
	Address      string    `json:"address"`
	Courier      string    `json:"courier"`
	ShippingCost uint      `json:"shippingCost"`
	FinalPrice   uint      `json:"finalPrice"`
	Cart         Cart_View `json:"cart"`
}

type Cart_View struct {
	Id         uint                `json:"id"`
	Username   string              `json:"username"`
	TotalQty   uint                `json:"totalQty"`
	TotalPrice uint                `json:"totalPrice"`
	Product    []Product_Cart_View `json:"product"`
}

type Product_Cart_View struct {
	CartID     uint         `json:"cartId"`
	ProductID  uint         `json:"productId"`
	Quantity   uint         `json:"quantity"`
	TotalPrice uint         `json:"totalPrice"`
	Product    Product_View `json:"product"`
}

type Product_View struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
}

type Detail_Product_View struct {
	Id           uint                 `json:"id"`
	Name         string               `json:"name"`
	Image        string               `json:"image"`
	Description  string               `json:"description"`
	Price        uint                 `json:"price"`
	Stock        uint                 `json:"stock"`
	FeedbackList []Feedback_Full_View `json:"feedbackList"`
}

type Feedback_Full_View struct {
	Username  string        `json:"username"`
	ProductId uint          `json:"productId"`
	Feedback  Feedback_View `json:"feedback"`
}

type Feedback_View struct {
	Id        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	ProductID uint      `json:"productId"`
	Comment   string    `json:"comment"`
	Rating    uint      `json:"rating"`
}

type History_View struct {
	Order []Order_View `json:"order"`
}

type Order_View struct {
	Invoice          string    `json:"invoice"`
	Id               uint      `json:"id"`
	CreatedAt        time.Time `json:"createdAt"`
	CartID           uint      `json:"cartId"`
	CustomerUsername string    `json:"customerUsername"`
	TotalQty         uint      `json:"totalQty"`
	TotalPrice       uint      `json:"totalPrice"`
	Status           string    `json:"status"`
	Address          string    `json:"address"`
	Courier          string    `json:"courier"`
	ProofOfPayment   string    `json:"proofOfPayment"`
	ValidatedBy      string    `json:"validationBy"`
}

type Detail_History_View struct {
	Invoice string               `json:"invoice"`
	OrderID uint                 `json:"orderId"`
	CartID  uint                 `json:"cartId"`
	Status  string               `json:"status"`
	Address string               `json:"address"`
	Courier string               `json:"courier"`
	Product []Product_Order_View `json:"product"`
}

type Product_Order_View struct {
	TotalQty   uint         `json:"totalQty"`
	TotalPrice uint         `json:"totalPrice"`
	Product    Product_View `json:"product"`
}

type Monthly_Report_View struct {
	Month  string         `json:"month"`
	Year   int            `json:"year"`
	Report []Money_Report `json:"report"`
}

type Money_Report struct {
	Date    time.Time `json:"date"`
	Income  uint      `json:"income"`
	Expense uint      `json:"expense"`
}

type Update_Order_Status_Binding struct {
	Status string `json:"status"`
}

type Product_View_Integrated struct {
	Id          uint          `json:"id"`
	Name        string        `json:"name"`
	Price       uint          `json:"price"`
	Stock       uint          `json:"stock"`
	Image       string        `json:"image"`
	Description string        `json:"description"`
	Reviews     []Review_View `json:"reviews"`
}

type Review_View struct {
	Username string `json:"username"`
	Feedback string `json:"feedback"`
	Rating   uint   `json:"rating"`
}
