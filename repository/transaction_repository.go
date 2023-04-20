package repository

import (
	"errors"
	"final_project_easycash/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type TransactionRepo interface {
	WithdrawBalance(sender string, receiver string, amount float64) error
	// TransferBalance(sender string, receiver string, amount float64) error
}

type transactionRepo struct {
	db *sqlx.DB
}

func (t *transactionRepo) WithdrawBalance(sender string, receiver string, amount float64) error {
	var senderInDb model.TransactionCode
	var receiverInDb model.TransactionCode

	row := t.db.QueryRow(`SELECT id, code FROM mst_transaction_codes WHERE code = $1`, sender)
	err := row.Scan(&senderInDb.Id, &senderInDb.Code)

	if senderInDb.Id == 0 {
		return errors.New("Sender number not found")
	}

	if err != nil {
		return err
	}

	row = t.db.QueryRow(`SELECT id, code FROM mst_transaction_codes WHERE code = $1`, receiver)
	err = row.Scan(&receiverInDb.Id, &receiverInDb.Code)

	if receiverInDb.Id == 0 {
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

// func (t *transactionRepo) TransferBalance(sender string, receiver string, amount float64) error {
// 	var senderInDb model.TransactionCode
// 	var receiverInDb model.TransactionCode

// 	row := t.db.QueryRow(`SELECT id, code FROM mst_transaction_code WHERE code = $1`, sender)
// 	err := row.Scan(&senderInDb.Id, &senderInDb.AccountType, &senderInDb.Code)

// 	if senderInDb.Id == 0 {
// 		return errors.New("Sender number not found")
// 	}

// 	if err != nil {
// 		return err
// 	}

// 	row = t.db.QueryRow(`SELECT id, code FROM mst_transaction_code WHERE code = $1`, receiver)
// 	err = row.Scan(&receiverInDb.Id, &receiverInDb.AccountType, &receiverInDb.Code)

// 	if receiverInDb.Id == 0 {
// 		return errors.New("Receiver number not found")
// 	}

// 	if err != nil {
// 		return err
// 	}

// 	query := "BEGIN; INSERT INTO trx_bills (sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id) VALUES ($1, $2, $3, $4, $5, $6, $7); UPDATE mst_users SET balance = balance - $8 WHERE phone_number = $9; UPDATE mst_users SET balance = balance - $10 WHERE phone_number = $11"
// 	_, err = t.db.Exec(query, &senderInDb.AccountType, &senderInDb.Code, "3", amount, time.Now(), &receiverInDb.AccountType, &receiverInDb.Code, amount, &senderInDb.Id, amount, &receiverInDb.Id)

// 	if err != nil {
// 		_, err = t.db.Exec("ROLLBACK;")
// 		return err
// 	}

// 	_, err = t.db.Exec("COMMIT;")
// 	if err != nil {
// 		_, err = t.db.Exec("ROLLBACK;")
// 		return err
// 	}

// 	return nil
// }

func NewTransactionRepo(db *sqlx.DB) TransactionRepo {
	repo := new(transactionRepo)
	repo.db = db
	return repo
}
