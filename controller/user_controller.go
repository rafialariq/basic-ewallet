package controller

import (
	"encoding/json"
	"final_project_easycash/model"
	"final_project_easycash/usecase"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	usecase usecase.UserUsecase
}

func (c *UserController) CheckProfile(ctx *gin.Context) {
	claims, exists := ctx.Keys["claims"].(jwt.MapClaims)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
		return
	}

	usernameToken, ok := claims["username"].(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.Header("Content-Disposition", "attachment; filename=data.json")
	username := ctx.Param("username")

	if usernameToken != username {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	res, err := c.usecase.CheckProfile(username)

	if err != nil {
		if err.Error() == "Username not found" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *UserController) EditProfile(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	rawBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user model.User
	if err := json.Unmarshal(rawBody, &user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
		return
	}

	usernameToken, ok := claims.(jwt.MapClaims)["username"].(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
		return
	}

	if usernameToken != user.Username {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err = c.usecase.EditProfile(&user)

	if err != nil {
		if err.Error() == "your username is too short or too long" || err.Error() == "invalid password" || err.Error() == "invalid password" || err.Error() == "invalid email" || err.Error() == "invalid phone number" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "profile edited"})
}

func (c *UserController) EditPhotoProfile(ctx *gin.Context) {
	username := ctx.Param("username")

	file, fileHeader, err := ctx.Request.FormFile("photo")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileName := strings.Split(fileHeader.Filename, ".")
	if len(fileName) != 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileExt := fileName[1]
	if strings.ToLower(fileExt) != "jpg" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file extension"})
		return
	}

	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "missing claims"})
		return
	}

	usernameToken, ok := claims.(jwt.MapClaims)["username"].(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims"})
		return
	}

	if usernameToken != username {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err = c.usecase.EditPhotoProfile(username, fileExt, &file)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "photo profile successfuly edited"})
}

func (c *UserController) UnregProfile(ctx *gin.Context) {
	username := ctx.Param("username")

	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
		return
	}

	usernameToken, ok := claims.(jwt.MapClaims)["username"].(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
		return
	}

	if usernameToken != username {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := c.usecase.UnregProfile(username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "profile successfully deleted"})
}

func NewUserController(rg *gin.RouterGroup, u usecase.UserUsecase) *UserController {
	controller := UserController{
		usecase: u,
	}
	rg.GET("/profile/:username", controller.CheckProfile)
	rg.POST("/profile/edit", controller.EditProfile)
	rg.POST("/profile/edit/photo/:username", controller.EditPhotoProfile)
	rg.DELETE("/profile/:username", controller.UnregProfile)
	return &controller
}
