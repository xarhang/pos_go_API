package main

import (
	"log"
	"os"
	"pos-go/db"
	"pos-go/route"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading env file")
		return
	}
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode("release")
		gin.SetMode(gin.ReleaseMode)
	}
	// gin.SetMode("release")
	dbType := os.Getenv("DB_TYPE")
	// Auth
	dsn := os.Getenv("DATABASE_DSN")
	a, _ := gormadapter.NewAdapter(dbType, dsn, true)
	e, _ := casbin.NewEnforcer("config/rbac_model.conf", a)
	e.LoadPolicy()

	//Database
	db.ConnectDB()
	//MigrateDB
	db.Migrate()
	//Set Path Permission
	os.MkdirAll("uploads/product", 0755)
	r := gin.New()
	r.SetTrustedProxies(nil)
	r.Static("/uploads", "./uploads")
	//Routing
	route.ServeRoutes(r)
	//Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "81"
	}
	r.Run(":" + port)
}
