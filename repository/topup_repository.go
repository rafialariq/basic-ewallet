package repository

import (
	"errors"
	"final_project_easycash/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type TopUpRepo interface {
	TopUpBalance(sender string, receiver string, amount float64) error
}

type topUpRepo struct {
	db *sqlx.DB
}

func (t *topUpRepo) TopUpBalance(sender string, receiver string, amount float64) error {
	var senderInDb model.TransactionCode
	var receiverInDb model.TransactionCode

	row := t.db.QueryRow(`SELECT phone_number FROM mst_user WHERE phone_number = $1`, receiver)
	err := row.Scan(&receiverInDb.Code)

	if receiverInDb.Code == "" {
		return errors.New("Receiver number not found")
	}

	if err != nil {
		return err
	}

	row = t.db.QueryRow(`SELECT bank_number FROM mst_bank WHERE bank_number = $1`, sender)
	err = row.Scan(&senderInDb.Code)

	if senderInDb.Code == "" {
		return errors.New("Sender number not found")
	}

	if err != nil {
		return err
	}

	query := "BEGIN;"
	_, err = t.db.Exec(query)

	if err != nil {
		return err
	}

	query = "INSERT INTO trx_bill (sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id) VALUES ($1, $2, $3, $4, $5, $6, $7);"
	_, err = t.db.Exec(query, 2, senderInDb.Code, 1, amount, time.Now(), 1, receiverInDb.Code)

	if err != nil {
		_, err = t.db.Exec("ROLLBACK;")
		return err
	}

	query = "UPDATE mst_user SET balance = balance + $1 WHERE phone_number = $2;"
	_, err = t.db.Exec(query, amount, receiverInDb.Code)

	if err != nil {
		_, err = t.db.Exec("ROLLBACK;")
		return err
	}

	_, err = t.db.Exec("COMMIT;")
	if err != nil {
		_, err = t.db.Exec("ROLLBACK;")
		return err
	}

	return nil
}

func NewTopUpRepo(db *sqlx.DB) TopUpRepo {
	repo := new(topUpRepo)
	repo.db = db
	return repo
}
