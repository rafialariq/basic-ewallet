package repository

import (
	"errors"
	"final_project_easycash/model"
	"final_project_easycash/utils"

	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	GetUserById(username string) (model.User, error)
	UpdateUserById(updatedUserData *model.User) error
	UpdatePhotoProfile(username string, filePath string) error
	DeleteUserById(username string) error
}

type userRepo struct {
	db *sqlx.DB
}

func (u *userRepo) GetUserById(username string) (model.User, error) {
	var user model.User
	row := u.db.QueryRow(`SELECT id, username, password, email, phone_number, photo_profile, balance FROM mst_user WHERE username = $1`, username)
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.PhoneNumber, &user.PhotoProfile, &user.Balance)

	if user.Username == "" {
		return model.User{}, errors.New("Username not found")
	}

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u *userRepo) UpdateUserById(updatedUserData *model.User) error {
	hashedPassword := utils.PasswordHashing(updatedUserData.Password)

	query := `UPDATE mst_user SET password = $1, email = $2, phone_number = $3 WHERE username = $4`
	_, err := u.db.Exec(query, &hashedPassword, &updatedUserData.Email, &updatedUserData.PhoneNumber, &updatedUserData.Username)

	if err != nil {
		return err
	}

	return nil
}

func (u *userRepo) UpdatePhotoProfile(username string, filePath string) error {
	_, err := u.GetUserById(username)

	if err != nil {
		return err
	}

	query := `UPDATE mst_user SET photo_profile = $1 WHERE username = $2`
	_, err = u.db.Exec(query, &filePath, &username)

	if err != nil {
		return err
	}

	return nil
}

func (u *userRepo) DeleteUserById(username string) error {
	_, err := u.GetUserById(username)

	if err != nil {
		return err
	}

	query := "DELETE FROM mst_user WHERE username = $1"

	_, err = u.db.Exec(query, username)

	if err != nil {
		return err
	}

	return nil
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	repo := new(userRepo)
	repo.db = db
	return repo
}
