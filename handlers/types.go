package handlers

// in redis Db: 2
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

// in redis Db: 2
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

// in redis Db: 4
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
