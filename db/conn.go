package db

import (
	"log"
	"os"

	"pos-go/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// "gorm.io/driver/postgres"
var Conn *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if err != nil {
		log.Fatal("Cannot Connect to the database " + err.Error())
	}
	Conn = db
}

func Migrate() {
	Conn.AutoMigrate(
		&model.Category{},
		&model.Product{},
		&model.Order{},
		&model.OrderItem{},
		&model.User{},
	)
}
