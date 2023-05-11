package usecase

import (
	"strconv"
	"time"

	"final_project_easycash/model"
	"final_project_easycash/repository"
	"final_project_easycash/utils"

	"github.com/dgrijalva/jwt-go"
)

type LoginService interface {
	UserLogin(user model.User) (bool, string)
}

type loginService struct {
	loginRepo repository.LoginRepo
}

func (l *loginService) UserLogin(user model.User) (bool, string) {
	recUser, res := l.loginRepo.FindUser(user)

	authDuration, _ := strconv.Atoi(utils.DotEnv("AUTH_DURATION", ".env"))

	if recUser {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = user.Username
		claims["exp"] = time.Now().Add(time.Minute * time.Duration(authDuration)).Unix()

		tokenString, err := token.SignedString([]byte(utils.DotEnv("TOKEN_KEY", ".env")))
		if err != nil {
			return false, "failed to generate token"
		}

		return true, tokenString
	} else {
		return false, res
	}
}

func NewLoginService(loginRepo repository.LoginRepo) LoginService {
	return &loginService{
		loginRepo: loginRepo,
	}
}
