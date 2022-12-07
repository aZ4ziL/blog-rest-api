package models

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Mendeklarasikan variabel untuk gorm.DB
var db *gorm.DB

func init() {
	log.Println("migrating all models...")
	db = GetDB()
	db.AutoMigrate(
		&Category{},
		&Article{},
		&User{},
		&Comment{},
		&SubComment{},
	)
}

func GetDB() *gorm.DB {
	dsn := "root:root@tcp(127.0.0.1:3306)/blog-rest-api?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}
