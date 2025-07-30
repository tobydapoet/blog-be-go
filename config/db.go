package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	database := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	port := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, database,
	)

	fmt.Println(dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error when connect DB!", err)
	}
	DB = db
	log.Println(("connect success!"))
}
