package main

import "github.com/jinzhu/gorm"

var DB *gorm.DB

func ConnectDataBase() {
	database, err := gorm.Open("sqlite3", "test.db")

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&user{}, &ArticleProccesed{})

	DB = database
}
