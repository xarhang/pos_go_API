package controller

import (
	"errors"
	"net/http"
	"time"

	"pos-go/db"
	"pos-go/dto"
	"pos-go/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Category struct{}

func (c Category) FindAll(ctx *gin.Context) {
	var categories []model.Category
	db.Conn.Find(&categories)
	var result []dto.CategoryResponse
	for _, category := range categories {
		result = append(result, dto.CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
		})
	}
	j := result
	if len(j) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "no item found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}
func (c Category) FindOne(ctx *gin.Context) {
	id := ctx.Param("id")
	var category model.Category
	if err := db.Conn.First(&category, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	})
}
func (c Category) Create(ctx *gin.Context) {

	var form dto.CategoryRequest

	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	category := model.Category{
		Name: form.Name,
	}
	if err := db.Conn.Create(&category).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	})
}
func (c Category) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var form dto.CategoryRequest
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	var category model.Category
	if err := db.Conn.First(&category, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	category.Name = form.Name
	db.Conn.Save(&category)
	ctx.JSON(http.StatusOK, dto.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	})
}
func (c Category) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	var category model.Category
	if err := db.Conn.First(&category, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	db.Conn.Unscoped().Delete(&model.Category{}, id)
	ctx.JSON(http.StatusOK, gin.H{"deletedAt": time.Now()})
}
