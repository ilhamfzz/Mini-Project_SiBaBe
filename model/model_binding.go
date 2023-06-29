package model

type Checkout_Binding struct {
	Courier string `json:"courier"`
	Address string `json:"address"`
}

type Payment_Binding struct {
	Invoice        string `json:"invoice"`
	ProofOfPayment string `json:"proofOfPayment"`
}

type Production_Binding struct {
	Date       string `json:"date"` // format: YYYY-MM-DD
	Name       string `json:"name"`
	TotalPrice uint   `json:"totalPrice"`
	Image      string `json:"image"`
}

type Login_Binding struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Feedback_Binding struct {
	ProductId uint   `json:"productId"`
	Feedback  string `json:"feedback"`
	Rating    uint   `json:"rating"`
}

type Product_Binding struct {
	Name        string `json:"name"`
	Image       string `json:"image"`
	Stock       uint   `json:"stock"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
}
