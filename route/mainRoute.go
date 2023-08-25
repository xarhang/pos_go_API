package route

import (
	"log"
	"os"
	"pos-go/controller"

	"github.com/gin-gonic/gin"

	"net/http"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/joho/godotenv"
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

func ServeRoutes(r *gin.Engine) {
	dbType := "mysql"
	if os.Getenv("APP_ENV") == "production" {
		dbType = "postgres"
	} else {
		if err := godotenv.Load(); err != nil {
			log.Fatal("error loading env file")
		}
	}
	// Auth
	dsn := os.Getenv("DATABASE_DSN")
	adapter, _ := gormadapter.NewAdapter(dbType, dsn, true)
	e, _ := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	e.LoadModel()
	e.LoadPolicy()

	r.Use(CORS())
	r.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to pos-go",
		})
	})
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "404", "message": "not found"})
	})

	//AuthorizationRoute
	AuthorizationRoute(r)

	//UserRoute
	UserRoute(r, adapter)
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
