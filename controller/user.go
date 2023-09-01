package controller

import (
	"net/http"
	"pos-go/db"
	"pos-go/dto"
	"pos-go/model"

	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt/v5"
)

type User struct {
}

func (u User) FindAll(ctx *gin.Context) {
	username := ctx.MustGet("aud").(string)

	var user []model.User
	if Rule := ctx.MustGet("rule_id").(int); Rule == 1 {
		db.Conn.Preload("Status").Preload("Rule").Find(&user)
	} else {
		db.Conn.Preload("Status").Preload("Rule").Find(&user, "username=?", username)
	}

	result := []dto.SigninResponse{}
	for _, u := range user {
		result = append(result, dto.SigninResponse{
			Username: u.Username,
			Fullname: u.Fullname,
			// RuleID:   u.Rule.ID,
			// StatusID: u.StatusID,
			Avatar: u.Avatar,
			Status: dto.StatusResponse{
				ID:         u.Status.ID,
				StatusName: u.Status.StatusName,
			},
			Rule: dto.RuleResponse{
				ID:       u.Rule.ID,
				RuleName: u.Rule.RuleName,
			},
		})
	}
	ctx.JSON(http.StatusOK, result)
}

func (u User) FindOne(ctx *gin.Context) {
	id := ctx.Param("id")
	var user model.User
	db.Conn.Find(&user, id)
	ctx.JSON(http.StatusOK, &user)
}
