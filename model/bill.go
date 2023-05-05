package model

import (
	"time"
)

type Bill struct {
	Id                int       `json:"id"`
	Id_transaction    string    `json:"id_transaction"`
	SenderTypeId      int       `json:"sender_type_id"`
	SenderId          string    `json:"sender_id"`
	TypeId            string    `json:"type_id"`
	Amount            float64   `json:"amount"`
	Date              time.Time `json:"date"`
	DestinationTypeId int       `json:"destination_type_id"`
	DestinationId     string    `json:"destination_id"`
	Status            string    `json:"status"`
}

/*func (b *Bill) GetDestinationId() []string {
	fmt.Println(b)
	if b.DestinationId == nil {
		return []string{}
	}

	if ids, ok := b.DestinationId.([]string); ok {
		return ids
	}

	return []string{b.DestinationId.(string)}
}

func (b *Bill) GetAmount() []float64 {
	if b.Amount == nil {
		return []float64{}
	}

	if amt, ok := b.Amount.([]float64); ok {
		return amt
	}

	return []float64{b.Amount.(float64)}
}*/
