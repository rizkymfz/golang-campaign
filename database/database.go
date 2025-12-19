package database

import (
	"fmt"
	"log"

	"github.com/rizkymfz/golang-campaign/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DB() *gorm.DB {
	dbHost := config.GetEnv("DB_HOST", "")
	dbPort := config.GetEnv("DB_PORT", "")
	dbName := config.GetEnv("DB_NAME", "")
	dbUser := config.GetEnv("DB_USER", "")
	dbPass := config.GetEnv("DB_PASS", "")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}
