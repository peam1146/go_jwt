package databases

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"main.go/models"
)

var DB *gorm.DB

// init user table
func InitDatabase() {
	db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = db
	DB.AutoMigrate(&models.User{})
}
