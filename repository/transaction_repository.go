package repository

import (
	"errors"
	"final_project_easycash/model"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type TransactionRepo interface {
	TransferMoney(sender string, receiver string, amount float64) error
	WithdrawBalance(sender string, receiver string, amount float64) error
	TransferBalance(sender string, receiver string, amount float64) error
	TopUpBalance(sender string, receiver string, amount float64) error
}

type transactionRepo struct {
	db *sqlx.DB
}

func (t *transactionRepo) TransferMoney(sender string, receiver string, amount float64) error {
	var balance float64
	var senderInDb model.User
	var merchantInDb model.Merchant

	row := t.db.QueryRow(`SELECT balance FROM mst_user WHERE phone_number = $1`, sender)
	err := row.Scan(&balance)

	if err != nil {
		return err
	}

	if balance < amount {
		return errors.New("Balance is not sufficient")
	}

	row = t.db.QueryRow(`SELECT phone_number FROM mst_user WHERE phone_number = $1`, sender)
	err = row.Scan(&senderInDb.PhoneNumber)

	if err != nil {
		return err
	}

	row = t.db.QueryRow(`SELECT merchantcode FROM mst_merchant WHERE merchantcode = $1`, receiver)
	err = row.Scan(&merchantInDb.MerchantCode)

	if err != nil {
		fmt.Println("receiver error")
		return err
	}

	query := "BEGIN;"
	_, err = t.db.Exec(query)

	if err != nil {
		return err
	}

	query = "INSERT INTO trx_bill (sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id) VALUES ($1, $2, $3, $4, $5, $6, $7);"
	_, err = t.db.Exec(query, 1, senderInDb.PhoneNumber, 2, amount, time.Now(), 3, merchantInDb.MerchantCode)

	if err != nil {
		_, err = t.db.Exec("ROLLBACK;")
		return errors.New("transaction failed")
	}

	query = "UPDATE mst_user SET balance = balance - $1 WHERE phone_number = $2;"
	_, err = t.db.Exec(query, amount, senderInDb.PhoneNumber)

	if err != nil {
		_, err = t.db.Exec("ROLLBACK;")
		return errors.New("transaction failed")
	}

	query = "UPDATE mst_merchant SET amount = amount + $1 WHERE merchantcode = $2;"
	_, err = t.db.Exec(query, amount, merchantInDb.MerchantCode)

	if err != nil {
		_, err = t.db.Exec("ROLLBACK;")
		return errors.New("transaction failed")
	}

	_, err = t.db.Exec("COMMIT;")
	if err != nil {
		_, err = t.db.Exec("ROLLBACK;")
		return errors.New("transaction failed")
	}

	return nil
}

func (t *transactionRepo) WithdrawBalance(sender string, receiver string, amount float64) error {
	senderType := 1
	receiverType := 2
	transactionType := 3
	statusType := 1
	var balance float64
	var senderInDb model.User
	var receiverInDb model.Bank

	row := t.db.QueryRow(`SELECT balance FROM mst_user WHERE phone_number = $1`, sender)
	err := row.Scan(&balance)

	if err != nil {
		return err
	}

	if balance < amount {
		return errors.New("Balance is not sufficient")
	}

	row = t.db.QueryRow(`SELECT phone_number FROM mst_user WHERE phone_number = $1`, sender)
	err = row.Scan(&senderInDb.PhoneNumber)

	if err != nil {
		return err
	}

	row = t.db.QueryRow(`SELECT bank_number FROM mst_bank WHERE bank_number = $1`, receiver)
	err = row.Scan(&receiverInDb.BankNumber)

	if err != nil {
		return err
	}

	query := "BEGIN;"
	_, err = t.db.Exec(query)

	if err != nil {
		return err
	}

	query = "INSERT INTO trx_bill (sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id, status_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);"
	_, err = t.db.Exec(query, senderType, senderInDb.PhoneNumber, transactionType, amount, time.Now(), receiverType, receiverInDb.BankNumber, statusType)

	if err != nil {
		_, err = t.db.Exec("ROLLBACK;")
		return errors.New("Transaction failed 1")
	}

	query = "UPDATE mst_user SET balance = balance - $1 WHERE phone_number = $2;"
	_, err = t.db.Exec(query, amount, senderInDb.PhoneNumber)

	if err != nil {
		_, err = t.db.Exec("ROLLBACK;")
		return errors.New("Transaction failed 2")
	}

	_, err = t.db.Exec("COMMIT;")
	if err != nil {
		_, err = t.db.Exec("ROLLBACK;")
		return errors.New("Transaction failed 3")
	}

	return nil
}

func (t *transactionRepo) TransferBalance(sender string, receiver string, amount float64) error {
	senderType := 1
	receiverType := 1
	transactionType := 3
	statusType := 1
	var balance float64
	var senderInDb model.User
	var receiverInDb model.User

	row := t.db.QueryRow(`SELECT balance FROM mst_user WHERE phone_number = $1`, sender)
	err := row.Scan(&balance)

	if err != nil {
		return err
	}

	if balance < amount {
		return errors.New("Balance is not sufficient")
	}

	row = t.db.QueryRow(`SELECT phone_number FROM mst_user WHERE phone_number = $1`, sender)
	err = row.Scan(&senderInDb.PhoneNumber)

	if err != nil {
		return err
	}

	row = t.db.QueryRow(`SELECT phone_number FROM mst_user WHERE phone_number = $1`, receiver)
	err = row.Scan(&receiverInDb.PhoneNumber)

	if err != nil {
		return err
	}

	query := "BEGIN;"
	_, err = t.db.Exec(query)

	if err != nil {
		return err
	}

	query = "INSERT INTO trx_bill (sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id, status_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);"
	_, err = t.db.Exec(query, senderType, senderInDb.PhoneNumber, transactionType, amount, time.Now(), receiverType, receiverInDb.PhoneNumber, statusType)

	if err != nil {
		_, err = t.db.Exec("ROLLBACK;")
		return errors.New("Transaction failed")
	}

	query = "UPDATE mst_user SET balance = balance - $1 WHERE phone_number = $2;"
	_, err = t.db.Exec(query, amount, senderInDb.PhoneNumber)

	if err != nil {
		_, err = t.db.Exec("ROLLBACK;")
		return errors.New("Transaction failed")
	}

	query = "UPDATE mst_user SET balance = balance + $1 WHERE phone_number = $2;"
	_, err = t.db.Exec(query, amount, receiverInDb.PhoneNumber)

	if err != nil {
		_, err = t.db.Exec("ROLLBACK;")
		return errors.New("Transaction failed")
	}

	_, err = t.db.Exec("COMMIT;")
	if err != nil {
		_, err = t.db.Exec("ROLLBACK;")
		return errors.New("Transaction failed")
	}

	return nil
}

func (t *transactionRepo) TopUpBalance(sender string, receiver string, amount float64) error {
	senderType := 2
	receiverType := 1
	transactionType := 1
	statusType := 1
	var senderInDb model.Bank
	var receiverInDb model.User

	row := t.db.QueryRow(`SELECT phone_number FROM mst_user WHERE phone_number = $1`, receiver)
	err := row.Scan(&receiverInDb.PhoneNumber)

	if err != nil {
		return err
	}

	row = t.db.QueryRow(`SELECT bank_number FROM mst_bank WHERE bank_number = $1`, sender)
	err = row.Scan(&senderInDb.BankNumber)

	if err != nil {
		return err
	}

	query := "BEGIN;"
	_, err = t.db.Exec(query)

	if err != nil {
		return err
	}

	query = "INSERT INTO trx_bill (sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id, status_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);"
	_, err = t.db.Exec(query, senderType, senderInDb.BankNumber, transactionType, amount, time.Now(), receiverType, receiverInDb.PhoneNumber, statusType)

	if err != nil {
		_, err = t.db.Exec("ROLLBACK;")
		return errors.New("Transaction failed")
	}

	query = "UPDATE mst_user SET balance = balance + $1 WHERE phone_number = $2;"
	_, err = t.db.Exec(query, amount, receiverInDb.PhoneNumber)

	if err != nil {
		_, err = t.db.Exec("ROLLBACK;")
		return errors.New("Transaction failed")
	}

	_, err = t.db.Exec("COMMIT;")
	if err != nil {
		_, err = t.db.Exec("ROLLBACK;")
		return errors.New("Transaction failed")
	}

	return nil
}

func NewTransactionRepo(db *sqlx.DB) TransactionRepo {
	repo := new(transactionRepo)
	repo.db = db
	return repo
}
