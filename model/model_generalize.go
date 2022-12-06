package model

import (
	"time"
)

type General_Customer struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Age       uint      `json:"age"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
}

type General_Chart struct {
	Id         uint      `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Username   string    `json:"username"`
	TotalPrice uint      `json:"totalPrice"`
	Status     string    `json:"status"`
}

type General_Product_Chart struct {
	ChartID    uint `json:"chartId"`
	ProductID  uint `json:"productId"`
	Quantity   uint `json:"quantity"`
	TotalPrice uint `json:"totalPrice"`
}

type General_Product struct {
	Id          uint      `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Name        string    `json:"name"`
	Image       string    `json:"image"`
	Stock       uint      `json:"stock"`
	Description string    `json:"description"`
	Price       uint      `json:"price"`
}

type General_Production struct {
	Id            uint      `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	AdminUsername string    `json:"adminUsername"`
	TotalPrice    uint      `json:"totalPrice"`
}

type General_Order struct {
	Id               uint      `json:"id"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	ChartID          uint      `json:"chartId"`
	CustomerUsername string    `json:"customerUsername"`
	TotalQty         uint      `json:"totalQty"`
	TotalPrice       uint      `json:"totalPrice"`
	Status           string    `json:"status"`
	Address          string    `json:"address"`
	Courier          string    `json:"courier"`
	ProofOfPayment   string    `json:"proofOfPayment"`
	ValidatedBy      string    `json:"validatedBy"`
}
