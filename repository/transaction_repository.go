package repository

import (
	"errors"
	"final_project_easycash/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type TransactionRepo interface {
	WithdrawBalance(sender string, receiver string, amount float64) error
	TransferBalance(sender string, receiver string, amount float64) error
	TopUpBalance(sender string, receiver string, amount float64) error
}

type transactionRepo struct {
	db *sqlx.DB
}

func (t *transactionRepo) WithdrawBalance(sender string, receiver string, amount float64) error {
	var balance float64
	var senderInDb model.TransactionCode
	var receiverInDb model.TransactionCode

	row := t.db.QueryRow(`SELECT balance FROM mst_user WHERE phone_number = $1`, sender)
	err := row.Scan(&balance)

	if err != nil {
		return err
	}

	if balance < amount {
		return errors.New("Balance is not sufficient")
	}

	row = t.db.QueryRow(`SELECT phone_number FROM mst_user WHERE phone_number = $1`, sender)
	err = row.Scan(&senderInDb.Code)

	if senderInDb.Code == "" {
		return errors.New("Sender number not found")
	}

	if err != nil {
		return err
	}

	row = t.db.QueryRow(`SELECT bank_number FROM mst_bank WHERE bank_number = $1`, receiver)
	err = row.Scan(&receiverInDb.Code)

	if receiverInDb.Code == "" {
		return errors.New("Receiver number not found")
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
	_, err = t.db.Exec(query, 1, senderInDb.Code, 3, amount, time.Now(), 2, receiverInDb.Code)

	if err != nil {
		_, err = t.db.Exec("ROLLBACK;")
		return err
	}

	query = "UPDATE mst_user SET balance = balance - $1 WHERE phone_number = $2;"
	_, err = t.db.Exec(query, amount, senderInDb.Code)

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

func (t *transactionRepo) TransferBalance(sender string, receiver string, amount float64) error {
	var balance float64
	var senderInDb model.TransactionCode
	var receiverInDb model.TransactionCode

	row := t.db.QueryRow(`SELECT balance FROM mst_user WHERE phone_number = $1`, sender)
	err := row.Scan(&balance)

	if err != nil {
		return err
	}

	if balance < amount {
		return errors.New("Balance is not sufficient")
	}

	row = t.db.QueryRow(`SELECT phone_number FROM mst_user WHERE phone_number = $1`, sender)
	err = row.Scan(&senderInDb.Code)

	if senderInDb.Code == "" {
		return errors.New("Sender number not found")
	}

	if err != nil {
		return err
	}

	row = t.db.QueryRow(`SELECT phone_number FROM mst_user WHERE phone_number = $1`, receiver)
	err = row.Scan(&receiverInDb.Code)

	if receiverInDb.Code == "" {
		return errors.New("Receiver number not found")
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
	_, err = t.db.Exec(query, 1, senderInDb.Code, 3, amount, time.Now(), 1, receiverInDb.Code)

	if err != nil {
		_, err = t.db.Exec("ROLLBACK;")
		return err
	}

	query = "UPDATE mst_user SET balance = balance - $1 WHERE phone_number = $2;"
	_, err = t.db.Exec(query, amount, senderInDb.Code)

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

func (t *transactionRepo) TopUpBalance(sender string, receiver string, amount float64) error {
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

func NewTransactionRepo(db *sqlx.DB) TransactionRepo {
	repo := new(transactionRepo)
	repo.db = db
	return repo
}