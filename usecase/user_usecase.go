package usecase

import (
	"encoding/base64"
	"errors"
	"final_project_easycash/model"
	"final_project_easycash/repository"
	"final_project_easycash/utils"
	"io/ioutil"
	"mime/multipart"
	"os"
)

type UserUsecase interface {
	CheckProfile(username string) (model.User, error)
	EditProfile(updatedUserData *model.User) error
	EditPhotoProfile(username string, fileExt string, file *multipart.File) error
	UnregProfile(username string) error
}

type userUsecase struct {
	userRepo repository.UserRepo
	fileRepo repository.FileRepository
}

func (u *userUsecase) CheckProfile(username string) (model.User, error) {
	res, err := u.userRepo.GetUserById(username)
	if res.PhotoProfile != "-" && res.PhotoProfile != "" {
		file, err := os.Open(res.PhotoProfile)
		if err != nil {
			return model.User{}, err
		}
		defer file.Close()

		img, err := ioutil.ReadAll(file)
		if err != nil {
			return model.User{}, err
		}
		res.PhotoProfile = base64.StdEncoding.EncodeToString(img)
	}
	return res, err
}

func (u *userUsecase) EditProfile(updatedUserData *model.User) error {
	if utils.ValidateEmail(updatedUserData.Email) {
		return errors.New("invalid email")
	} else if !utils.ValidatePhoneNumber(updatedUserData.PhoneNumber) {
		return errors.New("invalid phone number")
	}
	return u.userRepo.UpdateUserById(updatedUserData)
}

func (u *userUsecase) EditPhotoProfile(username string, fileExt string, file *multipart.File) error {
	filePath, err := u.fileRepo.Save("user_"+username+"."+fileExt, file)
	if err != nil {
		return err
	}
	return u.userRepo.UpdatePhotoProfile(username, filePath)
}

func (u *userUsecase) UnregProfile(username string) error {
	return u.userRepo.DeleteUserById(username)
}

func NewUserUsecase(userRepo repository.UserRepo, fileRepo repository.FileRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		fileRepo: fileRepo,
	}
}
