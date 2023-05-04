package repository

import (
	// "errors"
	// "errors"
	"final_project_easycash/model"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type HistoryRepo interface {
	GetHistoryByUser(user model.User) ([]model.Bill, error)
	GetHistoryWithAccountFilter(user model.User, accountTypeId int) ([]model.Bill, error)
	GetHistoryWithTypeFilter(user model.User, typeId string) ([]model.Bill, error)
	GetHistoryWithAmountFilter(user model.User, moreThan, lessThan float64) ([]model.Bill, error)
}

type historyRepo struct {
	db *sqlx.DB
}

func (h *historyRepo) GetHistoryByUser(user model.User) ([]model.Bill, error) {
	var historyList []model.Bill

	query := "SELECT id, id_transaction, sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id, status FROM trx_bill WHERE (sender_id = $1 OR destination_id = $1);"
	rows, err := h.db.Query(query, &user.PhoneNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var history model.Bill
		err := rows.Scan(&history.Id, &history.TransactionId, &history.SenderTypeId, &history.SenderId, &history.TypeId, &history.Amount, &history.Date, &history.DestinationTypeId, &history.DestinationId, &history.Status)

		if err != nil {
			return nil, err
		}

		historyList = append(historyList, history)
	}

	return historyList, nil
}

func (h *historyRepo) GetHistoryWithAccountFilter(user model.User, accountTypeId int) ([]model.Bill, error) {
	var historyList []model.Bill

	query := "SELECT id, id_transaction, sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id, status FROM trx_bill WHERE (sender_id = $1 OR destination_id = $1) AND (sender_type_id = $2 OR destination_type_id = $2);"
	rows, err := h.db.Query(query, &user.PhoneNumber, &accountTypeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var history model.Bill
		err := rows.Scan(&history.Id, &history.TransactionId, &history.SenderTypeId, &history.SenderId, &history.TypeId, &history.Amount, &history.Date, &history.DestinationTypeId, &history.DestinationId, &history.Status)

		if err != nil {
			return nil, err
		}

		historyList = append(historyList, history)
	}

	return historyList, nil
}

func (h *historyRepo) GetHistoryWithTypeFilter(user model.User, typeId string) ([]model.Bill, error) {
	var historyList []model.Bill

	query := "SELECT id, id_transaction, sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id, status FROM trx_bill WHERE (sender_id = $1 OR destination_id = $1) AND type_id = $2;"
	rows, err := h.db.Query(query, &user.PhoneNumber, &typeId)
	if err != nil {
		return historyList, err
	}
	defer rows.Close()

	for rows.Next() {
		var history model.Bill
		err := rows.Scan(&history.Id, &history.TransactionId, &history.SenderTypeId, &history.SenderId, &history.TypeId, &history.Amount, &history.Date, &history.DestinationTypeId, &history.DestinationId, &history.Status)

		if err != nil {
			return historyList, err
		}

		historyList = append(historyList, history)
	}

	return historyList, nil
}

func (h *historyRepo) GetHistoryWithAmountFilter(user model.User, moreThan, lessThan float64) ([]model.Bill, error) {
	var historyList []model.Bill

	query := "SELECT id, id_transaction, sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id, status FROM trx_bill WHERE (sender_id = $1 OR destination_id = $1) AND amount >= $2 AND amount <= $3;"
	rows, err := h.db.Query(query, &user.PhoneNumber, &moreThan, &lessThan)
	if err != nil {
		return historyList, err
	}
	defer rows.Close()

	for rows.Next() {
		var history model.Bill
		err := rows.Scan(&history.Id, &history.TransactionId, &history.SenderTypeId, &history.SenderId, &history.TypeId, &history.Amount, &history.Date, &history.DestinationTypeId, &history.DestinationId, &history.Status)

		if err != nil {
			return historyList, err
		}

		historyList = append(historyList, history)
	}

	return historyList, nil
}

func NewHistoryRepo(db *sqlx.DB) HistoryRepo {
	repo := new(historyRepo)
	repo.db = db
	return repo
}
