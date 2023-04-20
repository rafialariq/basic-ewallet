package usecase

import (
	"encoding/base64"
	"final_project_easycash/model"
	"final_project_easycash/repository"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strconv"
)

type UserUsecase interface {
	CheckProfile(id int) (model.User, error)
	EditProfile(updatedUserData *model.User) error
	EditPhotoProfile(id int, fileExt string, file *multipart.File) error
	UnregProfile(id int) error
}

type userUsecase struct {
	userRepo repository.UserRepo
	fileRepo repository.FileRepository
}

func (u *userUsecase) CheckProfile(id int) (model.User, error) {
	res, err := u.userRepo.GetUserById(id)
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
	return u.userRepo.UpdateUserById(updatedUserData)
}

func (u *userUsecase) EditPhotoProfile(id int, fileExt string, file *multipart.File) error {
	idString := strconv.Itoa(id)
	filePath, err := u.fileRepo.Save("user"+idString+"."+fileExt, file)
	if err != nil {
		return err
	}
	return u.userRepo.UpdatePhotoProfile(id, filePath)
}

func (u *userUsecase) UnregProfile(id int) error {
	return u.userRepo.DeleteUserById(id)
}

func NewUserUsecase(userRepo repository.UserRepo, fileRepo repository.FileRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		fileRepo: fileRepo,
	}
}
