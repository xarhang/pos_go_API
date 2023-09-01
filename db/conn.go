package db

import (
	"log"
	"os"

	"pos-go/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// "gorm.io/driver/postgres"
var Conn *gorm.DB

func ConnectDB() {
	if os.Getenv("DATABASE_DSN") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("error loading env file in connect DB")
			return
		}
	}
	dsn := os.Getenv("DATABASE_DSN")
	//check if Production or not
	if env := os.Getenv("DB_TYPE"); env == "postgres" {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Error),
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
		})

		if err != nil {
			log.Fatal("Cannot Connect to the database postgres " + err.Error())
		}
		Conn = db
	} else {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Error),
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
		})

		if err != nil {
			log.Fatal("Cannot Connect to the database mysql " + err.Error())
		}
		Conn = db
	}

}

func Migrate() {
	Conn.AutoMigrate(
		&model.Category{},
		&model.Product{},
		&model.Order{},
		&model.OrderItem{},
		&model.User{},
		&model.Status{},
	)
}
