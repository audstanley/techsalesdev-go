package handlers

import "math/big"

// in redis Db: 2 (each product)
// Db: 7 (pcb)
// Db: 8 (wires)
// Db: 9 (diodes)
// Db: 10 (caps)
type Product struct {
	Name        string  `json:"name"`
	Cost        float64 `json:"cost"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	InStock     bool    `json:"inStock"`
	ProductId   string  `json:"productId"`
	OnSale      bool    `json:"onSale"`
	Iat         uint64  `json:"iat"`
	Category    string  `json:"category"`
}

type ProductForCart struct {
	Name      string  `json:"name"`
	Cost      float64 `json:"cost"`
	InStock   bool    `json:"inStock"`
	ProductId string  `json:"productId"`
	Amount    uint64  `json:"amount"`
}

// in redis Db: 2 (then assembled)
type ProductReturn struct {
	OnSale      []Product
	NewArrivals []Product
}

// in redis Db: 3
type EmailPendingUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Iat      uint64 `json:"iat"`
	Link     string `json:"link"`
}

// in redis Db: 1 (as string) and 4 (as a jwt)
type User struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	EmailVerification bool   `json:"emailVerification"`
}

type UserShipping struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Address1  string `json:"address1"`
	Address2  string `json:"address2"`
	City      string `json:"city"`
	State     string `json:"state"`
	Zip       uint32 `json:"zip"`
}

// in redis Db: 6
type FullUserSigningUp struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password1 string `json:"password1"`
	Password2 string `json:"password2"`
	Address1  string `json:"address1"`
	Address2  string `json:"address2"`
	City      string `json:"city"`
	State     string `json:"state"`
	Zip       uint32 `json:"zip"`
	Pending   bool   `json:"pending"`
}

// for request POST body /forgotPassword
type JustEmail struct {
	Email string `json:"email"`
}

// in redis Db: 5
type ForgotPasswordEmailStruct struct {
	Email string `json:"email"`
	Iat   uint64 `json:"iat"`
	Link  string `json:"link"`
}

// in redis Db: 5
type PasswordForgotten struct {
	Hash string `json:"hash"`
	Iat  uint64 `json:"iat"`
}

type PasswordChecker struct {
	Email             string `json:"email"`
	Password1         string `json:"password1"`
	Password2         string `json:"password2"`
	EmailVerification bool   `json:"emailVerification"`
}

// Db: 11
type Cart struct {
	Products         []ProductForCart `json:"products"`
	Total            float64          `json:"total"`
	ConfirmationCode string           `json:"confirmationCode"`
	ShipTo           UserShipping     `json:"shipTo"`
	Placed           bool             `json:"placed"`
	Shipped          bool             `json:"shipped"`
	Email            string           `json:"email"`
}

type Orders struct {
	Orders []Cart `json:"total"`
}

type EtheriumWallet struct {
	Private  string    `json:"private"`
	Public   string    `json:"public"`
	Pending  bool      `json:"pending"`
	Ballance big.Float `json:"ballance"`
	Email    string    `json:"email"`
}

type PrivatePaymentOrders struct {
	Orders []Cart           `json:"total"`
	Wallet []EtheriumWallet `json:"wallet"`
}
