package repository

import (
	"fmt"
	"log"

	"final_project_easycash/model"
	"final_project_easycash/utils"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type RegisterRepo interface {
	UserRegister(newUser *model.User) (bool, string)
	RegisterValidate(newUser *model.User) bool
}

type registerRepo struct {
	db *sqlx.DB
}

func (r *registerRepo) UserRegister(newUser *model.User) (bool, string) {
	hashedPassword := utils.PasswordHashing(newUser.Password)

	query := "INSERT INTO mst_user (username, email, phone_number, password) VALUES ($1, $2, $3, $4);"
	_, err := r.db.Exec(query, &newUser.Username, &newUser.Email, &newUser.PhoneNumber, &hashedPassword)
	if err != nil {
		log.Println(err)
		return false, "failed to create user"
	}

	return true, "user created successfully"
}

func (r *registerRepo) RegisterValidate(recUser *model.User) bool {
	var resUser model.User

	query := "SELECT username, phone_number FROM mst_user WHERE username = $1 OR phone_number = $2;"
	row := r.db.QueryRow(query, &recUser.Username, &recUser.PhoneNumber)

	if err := row.Scan(&resUser.Username, &resUser.PhoneNumber); err != nil {
		fmt.Println("errornya: ", err)
	}

	if recUser.Username == resUser.Username || recUser.PhoneNumber == resUser.PhoneNumber {
		return true
	}

	return false

}

func NewRegisterRepo(db *sqlx.DB) RegisterRepo {
	repo := new(registerRepo)
	repo.db = db
	return repo
}
