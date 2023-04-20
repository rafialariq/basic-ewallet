package repository

import (
	"errors"
	"final_project_easycash/model"

	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	GetUserById(id int) (model.User, error)
	UpdateUserById(updatedUserData *model.User) error
	UpdatePhotoProfile(id int, filePath string) error
	DeleteUserById(id int) error
}

type userRepo struct {
	db *sqlx.DB
}

func (u *userRepo) GetUserById(id int) (model.User, error) {
	var user model.User
	row := u.db.QueryRow(`SELECT id, username, password, email, phone_number, photo_profile, balance FROM mst_user WHERE id = $1`, id)
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.PhoneNumber, &user.PhotoProfile, &user.Balance)

	if user.Id == 0 {
		return model.User{}, errors.New("User ID not found")
	}

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u *userRepo) UpdateUserById(updatedUserData *model.User) error {
	_, err := u.GetUserById(updatedUserData.Id)

	if err != nil {
		return err
	}

	query := `UPDATE mst_user SET username = $1, password = $2, email = $3, phone_number = $4 WHERE id = $5`
	_, err = u.db.Exec(query, &updatedUserData.Username, &updatedUserData.Password, &updatedUserData.Email, &updatedUserData.PhoneNumber, &updatedUserData.Id)

	if err != nil {
		return err
	}

	return nil
}

func (u *userRepo) UpdatePhotoProfile(id int, filePath string) error {
	_, err := u.GetUserById(id)

	if err != nil {
		return err
	}

	query := `UPDATE mst_user SET photo_profile = $1 WHERE id = $2`
	_, err = u.db.Exec(query, &filePath, &id)

	if err != nil {
		return err
	}

	return nil
}

func (u *userRepo) DeleteUserById(id int) error {
	_, err := u.GetUserById(id)

	if err != nil {
		return err
	}

	query := "DELETE FROM mst_users WHERE id = $1"

	_, err = u.db.Exec(query, id)

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
