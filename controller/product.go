package controller

import (
	"errors"
	"net/http"
	"os"
	"pos-go/db"
	"pos-go/dto"
	"pos-go/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
}

func (p Product) FindAll(ctx *gin.Context) {
	categoryID := ctx.Query("catagoryId")
	status := ctx.Query("status")
	search := ctx.Query("search")
	var product []model.Product
	query := db.Conn.Preload("Category")
	if categoryID != "" {
		query = query.Where("category_id = ? ", &categoryID)
	}
	if status != "" {
		query = query.Where("status = ? ", &status)
	}
	if search != "" {
		query = query.Where("sku LIKE ? or name LIKE ? ", "%"+search+"%", "%"+search+"%")
	}
	query.Order("id desc").Find(&product)
	var result []dto.ReadProductResponse
	for _, product := range product {
		result = append(result, dto.ReadProductResponse{
			ID:     product.ID,
			Name:   product.Name,
			SKU:    product.SKU,
			Desc:   product.Desc,
			Price:  product.Price,
			Status: product.Status,
			Image:  product.Image,
			Category: dto.CategoryResponse{
				ID:   product.Category.ID,
				Name: product.Category.Name,
			},
		})
	}
	j := result
	if len(j) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "no item found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}
func (p Product) FindOne(ctx *gin.Context) {
	id := ctx.Param("id")
	var product []model.Product
	query := db.Conn.Preload("Category")
	query.First(&product, id)
	var result []dto.ReadProductResponse
	for _, product := range product {
		result = append(result, dto.ReadProductResponse{
			ID:     product.ID,
			Name:   product.Name,
			SKU:    product.SKU,
			Desc:   product.Desc,
			Price:  product.Price,
			Status: product.Status,
			Image:  product.Image,
			Category: dto.CategoryResponse{
				ID:   product.Category.ID,
				Name: product.Category.Name,
			},
		})
	}
	j := result
	if len(j) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "no item found"})
		return
	}
	ctx.JSON(http.StatusOK, result)
}
func (p Product) Create(ctx *gin.Context) {
	var form dto.ProductRequest
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uimage := form.Image

	// image, err := ctx.FormFile("image")
	image := []byte(uimage)

	imagePath := ""
	if image != nil {
		imagePath = "./uploads/product/" + uuid.New().String()
		os.WriteFile(imagePath, image, 0755)

	}
	// ctx.SaveUploadedFile(image, imagePath)

	product := model.Product{
		Name:       form.Name,
		SKU:        form.SKU,
		Desc:       form.Desc,
		Price:      form.Price,
		Status:     form.Status,
		Image:      imagePath,
		CategoryID: form.CategoryID,
	}
	if err := db.Conn.Create(&product).Error; err != nil {
		os.Remove(imagePath)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, dto.CreateOrUpdateProductResponse{
		ID:         product.ID,
		Name:       product.Name,
		SKU:        product.SKU,
		Desc:       product.Desc,
		Price:      product.Price,
		Status:     product.Status,
		CategoryID: product.CategoryID,
		Image:      product.Image,
	})

}
func (p Product) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var form dto.UpdateProductRequest
	//check if parametter is correcr
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var product model.Product
	//check if have the id the db or not
	if err := db.Conn.First(&product, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	image, err := ctx.FormFile("image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//check if image is null or not
	if image != nil {
		imagePath := "./uploads/product/" + uuid.New().String()
		ctx.SaveUploadedFile(image, imagePath)
		//remove the old image
		os.Remove(product.Image)
		product.Image = imagePath
	}
	if len(form.Name) > 1 {
		product.Name = form.Name
	}
	if len(form.SKU) > 1 {
		product.SKU = form.SKU
	}
	if len(form.Desc) > 1 {
		product.Desc = form.Desc
	}
	if form.Status == 1 || form.Status == 2 {
		product.Status = form.Status
	} else if form.Status == 0 {

	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input status"})
		return
	}

	if form.Price > 0 {
		product.Price = form.Price
	}
	if form.CategoryID > 0 {
		product.CategoryID = form.CategoryID
	}
	db.Conn.Save(&product)
	ctx.JSON(http.StatusOK, dto.CreateOrUpdateProductResponse{
		ID:         product.ID,
		Name:       product.Name,
		SKU:        product.SKU,
		Desc:       product.Desc,
		Price:      product.Price,
		Status:     product.Status,
		CategoryID: product.CategoryID,
		Image:      product.Image,
	})

}
func (p Product) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	db.Conn.Unscoped().Delete(&model.Product{}, id)
	ctx.JSON(http.StatusOK, gin.H{"deletedAt": time.Now()})
}
