package main

import (
	"pos-go/controller"

	"github.com/gin-gonic/gin"

	"net/http"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, OPTIONS, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func serveRoutes(r *gin.Engine) {
	r.Use(CORS())
	r.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to pos-go",
		})
	})
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "404", "message": "not found"})
	})

	//authentication
	authController := controller.Auth{}
	authGroup := r.Group("/auth")
	authGroup.POST("/register", authController.Register)

	//product
	productController := controller.Product{}
	productGroup := r.Group("/products")
	productGroup.GET("", productController.FindAll)
	productGroup.GET("/:id", productController.FindOne)
	productGroup.POST("", productController.Create)
	productGroup.PATCH("/:id", productController.Update)
	productGroup.DELETE("/:id", productController.Delete)

	categoryController := controller.Category{}
	categoryGroup := r.Group("/categories")
	categoryGroup.GET("", categoryController.FindAll)
	categoryGroup.GET("/:id", categoryController.FindOne)
	categoryGroup.POST("", categoryController.Create)
	categoryGroup.PATCH("/:id", categoryController.Update)
	categoryGroup.DELETE("/:id", categoryController.Delete)

	oderController := controller.Order{}
	oderGroup := r.Group("/orders")
	oderGroup.GET("", oderController.FindAll)
	oderGroup.GET("/:id", oderController.FindOne)
	oderGroup.POST("", oderController.Create)
	// oderGroup.PATCH("/:id", oderController.Update)
	// oderGroup.DELETE("/:id", oderController.Delete)
}
