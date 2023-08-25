package controller

import (
	"net/http"
	"os"
	"pos-go/db"
	"pos-go/dto"
	"pos-go/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct{}

var hamacSimpleSecret []byte

func (auth Auth) Signin(ctx *gin.Context) {

	var userData dto.SigninRequest
	if err := ctx.ShouldBindJSON(&userData); err != nil {
		if err.Error() == "EOF" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"error": "please input request body"}})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user model.User

	db.Conn.Preload("Status").Preload("Rule").First(&user, "username=?", userData.Username)
	var result dto.SigninResponse
	password_db := ""
	result = dto.SigninResponse{
		Username: user.Username,
		Fullname: user.Fullname,
		// RuleID:   user.RuleID,
		// StatusID: user.StatusID,
		Avatar: user.Avatar,
		Rule:   dto.RuleResponse{ID: user.Status.ID, RuleName: user.Rule.RuleName},
		Status: dto.StatusResponse{ID: user.Status.ID, StatusName: user.Status.StatusName},
	}
	password_db = user.Password
	if result.Username != "" {
		// ctx.JSON(http.StatusOK, gin.H{"db": password_db, "byte": []byte(password_db), "db_convert": string([]byte(password_db)), "input": []byte(userData.Password)})

		check := bcrypt.CompareHashAndPassword([]byte(password_db), []byte(userData.Password))

		if check != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "login failed, please check username & password"})
			return
		} else {
			hamacSimpleSecret = []byte(os.Getenv("JWT_SECRETKEY"))
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"username": user.Username,
				"exp":      time.Now().Add(time.Minute * 10).Unix(),
				"rule_id":  user.RuleID,
				"avatar":   user.Avatar,
			})

			tokenString, err := token.SignedString(hamacSimpleSecret)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			// fmt.Println("here.........................?", err)
			// fmt.Println(tokenString)

			ctx.JSON(http.StatusOK, gin.H{"accessToken": tokenString, "userData": result})
			return

		}

	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "login failed, please check username & password"})
		return
	}

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
		StatusID: userData.StatusID,
		RuleID:   userData.RuleID,
	}
	if err := db.Conn.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, &userData)
}
