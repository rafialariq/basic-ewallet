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
}

type transactionRepo struct {
	db *sqlx.DB
}

var (
	ErrBillNotFound        = errors.New("bill not found")
	ErrBillPaid            = errors.New("bill has already been paid")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

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

	if senderInDb.PhoneNumber == "" {
		return errors.New("Sender number not found")
	}

	if err != nil {
		return err
	}

	row = t.db.QueryRow(`SELECT merchantcode FROM mst_merchant WHERE merchantcode = $1`, receiver)
	err = row.Scan(&merchantInDb.MerchantCode)

	if merchantInDb.MerchantCode == "" {
		return errors.New("Merchant not found")
	}

	if err != nil {
		fmt.Println("receiver error")
		return err
	}

	query := "BEGIN;"
	_, err = t.db.Exec(query)

	if err != nil {
		return err
	}

	query = "INSERT INTO trx_bill (sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id, status) VALUES ($1, $2, $3, $4, $5, $6, $7,$8);"
	_, err = t.db.Exec(query, 1, senderInDb.PhoneNumber, 2, amount, time.Now(), 3, merchantInDb.MerchantCode, 2)

	if err != nil {
		fmt.Println("aaa")
		_, err = t.db.Exec("ROLLBACK;")
		return errors.New("transaction failed")
	}

	query = "UPDATE mst_user SET balance = balance - $1 WHERE phone_number = $2;"
	_, err = t.db.Exec(query, amount, senderInDb.PhoneNumber)

	if err != nil {
		fmt.Println("bbb")
		_, err = t.db.Exec("ROLLBACK;")
		return errors.New("transaction failed")
	}

	query = "UPDATE mst_merchant SET amount = amount + $1 WHERE merchantcode = $2;"
	_, err = t.db.Exec(query, amount, merchantInDb.MerchantCode)

	if err != nil {
		fmt.Println("ccc")
		_, err = t.db.Exec("ROLLBACK;")
		return errors.New("transaction failed")
	}

	_, err = t.db.Exec("COMMIT;")
	if err != nil {
		fmt.Println("ddd")
		_, err = t.db.Exec("ROLLBACK;")
		return errors.New("transaction failed")
	}

	return nil
}

func NewTransactionRepo(db *sqlx.DB) TransactionRepo {
	repo := new(transactionRepo)
	repo.db = db
	return repo
}
