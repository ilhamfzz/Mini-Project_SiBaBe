package model

type Checkout_Binding struct {
	Courier string `json:"courier"`
	Address string `json:"address"`
}

type Payment_Binding struct {
	ProofOfPayment string `json:"proofOfPayment"`
}

type Production_Binding struct {
	ProductName string `json:"productName"`
	TotalPrice  uint   `json:"totalPrice"`
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
