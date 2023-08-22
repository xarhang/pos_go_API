package controller

import (
	"net/http"
	"pos-go/db"
	"pos-go/dto"
	"pos-go/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct{}

func (auth Auth) Signin(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (r Auth) Register(ctx *gin.Context) {
	var userData dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&userData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pwdEncrypted, _ := bcrypt.GenerateFromPassword([]byte(userData.Password), 10)
	userData.Password = string(pwdEncrypted)
	user := model.User{
		Username: userData.Username,
		Password: userData.Password,
		Fullname: userData.Fullname,
		Avatar:   userData.Avatar,
		Status:   userData.Status,
		Rule:     userData.Rule,
	}
	if err := db.Conn.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, &userData)
}
