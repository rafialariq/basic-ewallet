package repository

import (
	"database/sql"
	"errors"
	"final_project_easycash/model"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type TransactionRepo interface {
	TransferMoney(sender string, receiver string, amount float64) error
	SplitBill(sender string, receiver []string, amount []float64) error
	PayBill(receiver string, idTransaction string) error
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

func (t *transactionRepo) SplitBill(sender string, receiver []string, amount []float64) error {
	var balance float64
	var senderInDb model.User
	var receiverInDb model.User

	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT balance FROM mst_user WHERE phone_number = $1`, sender)
	err = row.Scan(&balance)
	if err != nil {
		return err
	}

	totalAmount := 0.0
	for _, amount := range amount {
		totalAmount += amount
	}

	if balance < totalAmount {
		return errors.New("Balance is not sufficient")
	}

	row = tx.QueryRow(`SELECT phone_number FROM mst_user WHERE phone_number = $1`, sender)
	err = row.Scan(&senderInDb.PhoneNumber)
	if senderInDb.PhoneNumber == "" {
		return errors.New("Sender number not found")
	}
	if err != nil {
		return err
	}

	for i, receiver := range receiver {
		row = tx.QueryRow(`SELECT phone_number FROM mst_user WHERE phone_number = $1`, receiver)
		err = row.Scan(&receiverInDb.PhoneNumber)
		if receiverInDb.PhoneNumber == "" {
			return errors.New(fmt.Sprintf("Receiver number at index %d not found", i))
		}
		if err != nil {
			return err
		}

		query := "INSERT INTO trx_bill (sender_type_id, sender_id, type_id, amount, date, destination_type_id, destination_id, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);"
		_, err = tx.Exec(query, 1, senderInDb.PhoneNumber, 4, amount[i], time.Now(), 1, receiverInDb.PhoneNumber, 1)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (t *transactionRepo) PayBill(receiver string, id_transaction string) error {
	var billAmount float64
	var senderInDb model.User
	var receiverInDb model.User
	var status int

	row := t.db.QueryRow(`SELECT amount, destination_id, status, sender_id FROM trx_bill WHERE id_transaction = $1`, id_transaction)
	err := row.Scan(&billAmount, &receiverInDb.PhoneNumber, &status, &senderInDb.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrBillNotFound
		}
		return err
	}

	if status == 2 {
		return ErrBillPaid
	}

	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Mendapatkan saldo penerima tagihan
	var receiverBalance float64
	row = tx.QueryRow(`SELECT balance FROM mst_user WHERE phone_number = $1`, receiverInDb.PhoneNumber)
	err = row.Scan(&receiverBalance)
	if err != nil {
		return err
	}

	// Jika saldo penerima kurang dari jumlah tagihan
	if receiverBalance < billAmount {
		return ErrInsufficientBalance
	}

	// Mengurangi saldo penerima sebesar jumlah tagihan
	query := `UPDATE mst_user SET balance = balance - $1 WHERE phone_number = $2`
	_, err = tx.Exec(query, billAmount, receiverInDb.PhoneNumber)
	if err != nil {
		return err
	}

	// Menambah saldo pengirim sebesar jumlah tagihan
	query = `UPDATE mst_user SET balance = balance + $1 WHERE phone_number = $2`
	_, err = tx.Exec(query, billAmount, senderInDb.PhoneNumber)
	if err != nil {
		return err
	}

	// Mengubah status tagihan menjadi "paid"
	query = `UPDATE trx_bill SET status = $1 WHERE id_transaction = $2`
	_, err = tx.Exec(query, 2, id_transaction)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func NewTransactionRepo(db *sqlx.DB) TransactionRepo {
	repo := new(transactionRepo)
	repo.db = db
	return repo
}
