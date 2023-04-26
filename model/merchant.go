package model

type Merchant struct {
	Id           int     `json:"id"`
	MerchantCode string  `json:"merchantcode"`
	Name         string  `json:"name"`
	Amount       float64 `json:"amount"`
}
