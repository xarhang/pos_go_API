package main

import (
	"log"
	"os"
	"pos-go/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode("release")
	} else {
		if err := godotenv.Load(); err != nil {
			log.Fatal("error loading env file")
		}
	}

	db.ConnectDB()
	db.Migrate()
	os.MkdirAll("uploads/product", 0755)
	r := gin.Default()
	r.Static("/uploads", "./uploads")
	serveRoutes(r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "81"
	}
	r.Run(":" + port)
	//test
}
