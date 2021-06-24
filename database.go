package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB_user *gorm.DB
var DB_article *gorm.DB

func ConnectDataBase() {
	database_user, err := gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	database_user.AutoMigrate(&user{})

	DB_user = database_user

	database_article, err := gorm.Open(sqlite.Open("articles.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	database_article.AutoMigrate(&article{})

	DB_article = database_article
}
