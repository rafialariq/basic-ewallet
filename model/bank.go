package model

type Bank struct {
	Id         int    `json:"id"`
	BankNumber string `json:"bank_number"`
	Name       string `json:"name"`
}
