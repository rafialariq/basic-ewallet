package model

import "time"

type Bill struct {
	Id                int       `json:"id"`
	TransactionId     string    `json:"id_transaction"`
	SenderTypeId      int       `json:"sender_type_id"`
	SenderId          string    `json:"sender_id"`
	TypeId            int       `json:"type_id"`
	Amount            float64   `json:"amount"`
	Date              time.Time `json:"date"`
	DestinationTypeId int       `json:"destination_type_id"`
	DestinationId     string    `json:"destination_id"`
	Status            int       `json:"status"`
}
