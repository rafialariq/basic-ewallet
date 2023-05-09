package usecase

import (
	"time"

	"final_project_easycash/model"
	"final_project_easycash/repository"
	"final_project_easycash/utils"

	"github.com/dgrijalva/jwt-go"
)

type RegisterService interface {
	UserSignup(newUser *model.User) (bool, string)
}

type registerService struct {
	registerRepo repository.RegisterRepo
}

func (r *registerService) UserSignup(newUser *model.User) (bool, string) {
	if !utils.IsUsernameValid(newUser.Username) {
		return false, "your username is too short or too long"
	} else if !utils.IsPasswordValid(newUser.Password) {
		return false, "invalid password"
	} else if !utils.IsEmailValid(newUser.Email) {
		return false, "invalid email"
	} else if !utils.IsPhoneNumberValid(newUser.PhoneNumber) {
		return false, "invalid phone number"
	} else if r.registerRepo.RegisterValidate(newUser) {
		return false, "user already exist"
	}

	newUser.Password = utils.PasswordHashing(newUser.Password)

	user, res := r.registerRepo.UserRegister(newUser)

	if user {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = newUser.Username
		claims["exp"] = time.Now().Add(time.Minute * 5).Unix() // hardcode

		tokenString, err := token.SignedString([]byte("secretkey")) // hardcode
		if err != nil {
			return false, "failed to generate token"
		}

		return true, tokenString
	} else {
		return false, res
	}
}

func NewRegisterService(registerRepo repository.RegisterRepo) RegisterService {
	return &registerService{
		registerRepo: registerRepo,
	}
}
