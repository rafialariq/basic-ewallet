package repository

import (
	"errors"
	"final_project_easycash/model"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type HistoryRepo interface {
	GetAllByUser(user model.User) ([]model.Bill, error)
	GetByAccountType(user model.User, senderTypeId int) ([]model.Bill, error)
	GetByType(user model.User, typeId string) ([]model.Bill, error)
	GetByAmount(user model.User, moreThan, lessThan float64) ([]model.Bill, error)
}

type historyRepo struct {
	db *sqlx.DB
}

func (h *historyRepo) GetAllByUser(user model.User) ([]model.Bill, error) {
	var historyList []model.Bill

	query := "SELECT id, sender_type_id, sender_id, type_id, amount, destination_type_id, destination_id FROM trx_bill WHERE sender_id = $1 OR destination_id = $2;"
	rows, err := h.db.Query(query, &user.PhoneNumber, &user.PhoneNumber)
	if err != nil {
		return historyList, errors.New("gagal 1")
	}
	defer rows.Close()

	for rows.Next() {
		var history model.Bill
		err := rows.Scan(&history.Id, &history.SenderTypeId, &history.SenderId, &history.TypeId, &history.Amount, &history.DestinationTypeId, &history.DestinationId)

		if err != nil {
			return historyList, errors.New("gagal 2")
		}

		historyList = append(historyList, history)
	}

	return historyList, nil
}

func (h *historyRepo) GetByAccountType(user model.User, accountTypeId int) ([]model.Bill, error) {
	var historyList []model.Bill

	query := "SELECT id, sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id FROM trx_bill WHERE sender_id = $1 OR destination_id = $1 AND sender_type_id = $2 OR destination_type_id = $2;"
	rows, err := h.db.Query(query, &user.PhoneNumber, &accountTypeId)
	if err != nil {
		return historyList, err
	}
	defer rows.Close()

	for rows.Next() {
		var history model.Bill
		err := rows.Scan(&history.Id, &history.SenderTypeId, &history.SenderId, &history.TypeId, &history.Amount, &history.Date, &history.DestinationTypeId, &history.DestinationId)

		if err != nil {
			return historyList, err
		}

		historyList = append(historyList, history)
	}

	return historyList, nil
}

func (h *historyRepo) GetByType(user model.User, typeId string) ([]model.Bill, error) {
	var historyList []model.Bill

	query := "SELECT id, sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id FROM trx_bill WHERE sender_id = $1 OR destination_id = $1 AND type_id = $2;"
	rows, err := h.db.Query(query, &user.PhoneNumber, &typeId)
	if err != nil {
		return historyList, err
	}
	defer rows.Close()

	for rows.Next() {
		var history model.Bill
		err := rows.Scan(&history.Id, &history.SenderTypeId, &history.SenderId, &history.TypeId, &history.Amount, &history.Date, &history.DestinationTypeId, &history.DestinationId)

		if err != nil {
			return historyList, err
		}

		historyList = append(historyList, history)
	}

	return historyList, nil
}

func (h *historyRepo) GetByAmount(user model.User, moreThan, lessThan float64) ([]model.Bill, error) {
	var historyList []model.Bill

	query := "SELECT id, sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id FROM trx_bill WHERE sender_id = $1 OR destination_id = $1 AND amount > $2 AND amount < $3;"
	rows, err := h.db.Query(query, &user.PhoneNumber, &moreThan, &lessThan)
	if err != nil {
		return historyList, err
	}
	defer rows.Close()

	for rows.Next() {
		var history model.Bill
		err := rows.Scan(&history.Id, &history.SenderTypeId, &history.SenderId, &history.TypeId, &history.Amount, &history.Date, &history.DestinationTypeId, &history.DestinationId)

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
