package controller

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"pos-go/db"
	"pos-go/dto"
	"pos-go/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
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
			rans := uuid.NewString()
			hash := md5.Sum([]byte(user.Username + rans))
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"aud":     user.Username,
				"sub":     fmt.Sprintf("%v", user.ID),
				"tij":     hex.EncodeToString(hash[:]),
				"iss":     "API-Gateway",
				"iat":     time.Now().Unix(),
				"nbf":     time.Now().Unix(),
				"exp":     time.Now().Add(time.Minute * 10).Unix(),
				"rule_id": user.RuleID,
				"avatar":  user.Avatar,
			})

			tokenString, err := token.SignedString(hamacSimpleSecret)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			// fmt.Println("here.........................?", err)
			// fmt.Println(tokenString)
			refreshToken := jwt.New(jwt.SigningMethodHS256)
			rtClaims := refreshToken.Claims.(jwt.MapClaims)
			rtClaims["sub"] = fmt.Sprintf("%v", user.ID)
			rtClaims["aud"] = user.Username
			rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
			rt, err2 := refreshToken.SignedString([]byte("secret"))
			if err2 != nil {
				ctx.JSON(http.StatusOK, gin.H{"accessToken": tokenString, "userData": result})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"accessToken": tokenString, "refresh_token": rt, "userData": result})
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
