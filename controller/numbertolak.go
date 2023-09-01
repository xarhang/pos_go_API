package controller

import (
	"net/http"
	"pos-go/middleware"

	"github.com/gin-gonic/gin"
)

type NumberToLak struct {
	Number string `json:"number" binding:"required"`
}

func (n NumberToLak) ConvertToText(ctx *gin.Context) {
	var x = NumberToLak{}
	if err := ctx.ShouldBindJSON(&x); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error:": err.Error()})
		return
	}
	laoWord := middleware.ConvertToText(x.Number)
	ctx.JSON(http.StatusOK, gin.H{"lao_word": laoWord})
}
