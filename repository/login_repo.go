package repository

import (
	"log"

	"final_project_easycash/model"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type LoginRepo interface {
	FindUser(recUser model.User) (bool, string)
}

type loginRepo struct {
	db *sqlx.DB
}

func (l *loginRepo) FindUser(recUser model.User) (bool, string) {
	var resUser model.User
	query := "SELECT username, password FROM mst_user WHERE username = $1;"
	row := l.db.QueryRow(query, recUser.Username)

	if err := row.Scan(&resUser.Username, &resUser.Password); err != nil {
		log.Println(err)
		return false, "user not found"
	}

	err := bcrypt.CompareHashAndPassword([]byte(resUser.Password), []byte(recUser.Password))
	if err != nil {
		log.Println(err)
		return false, "invalid password"
	}

	return true, "successfully login"

}

func NewLoginRepo(db *sqlx.DB) LoginRepo {
	repo := new(loginRepo)
	repo.db = db
	return repo
}
