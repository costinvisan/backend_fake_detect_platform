// models.article.go

package main

import (
	"gorm.io/gorm"
)

type article struct {
	gorm.Model
	Title   string `json:"title"`
	Url     string `json:"url"`
	Rating  int    `json:"rating"`
	User_id int    `json:"user_id"`
}

// Return a list of all the articles
func getAllArticles() []article {
	var articleList []article
	DB_article.Find(&articleList)
	return articleList
}

// Fetch an article based on the ID supplied
func getArticleByID(id int) article {
	var article article
	DB_article.Find(&article, id)
	return article
}

func getArticleByUserID(id int) []article {
	var articles []article
	DB_article.Where("user_id = ?", id).Find(&articles)
	return articles
}

// Create a new article with the title and content provided
func createNewArticle(title, url string, rating, user_id int) (*article, error) {
	// Set the ID of a new article to one more than the number of articles
	a := article{Title: title, Url: url, Rating: rating, User_id: user_id}

	// Add the article to the list of articles
	DB_article.Create(&a)
	return &a, nil
}
