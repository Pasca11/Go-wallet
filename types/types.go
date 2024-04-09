package types

type CreateAccountRequest struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Patronymic string `json:"patronymic,omitempty"`
	Password   string `json:"password"`
}

type LoginRequest struct {
	Wallet   int    `json:"wallet"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Wallet int64  `json:"wallet"`
	Token  string `json:"token"`
}

type TransferRequest struct {
	From   int `json:"from"`
	To     int `json:"to"`
	Amount int `json:"amount"`
}
