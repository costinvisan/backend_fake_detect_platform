// handlers.article.go

package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	googlesearch "github.com/rocketlaunchr/google-search"
)

type URLstring struct {
	Url string `json:"url" binding:"required"`
}

func showIndexPage(c *gin.Context) {
	loggedInInterface, _ := c.Get("is_logged_in")
	loggedIn := loggedInInterface.(bool)
	var articles []article
	if loggedIn {
		articles = getAllArticles()
	} else {
		articles = []article{}
	}
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title":   "Home Page",
		"payload": articles}, "index.html")
}

func showArticleCreationPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title": "Create New Article"}, "create-article.html")
}

func getArticle(c *gin.Context) {
	// Check if the article ID is valid
	if articleID, err := strconv.Atoi(c.Param("article_id")); err == nil {
		// Check if the article exists
		if art := getArticleByID(articleID); (art != article{}) {
			// Call the render function with the title, article and the name of the
			// template
			render(c, gin.H{
				"title":   art.Title,
				"payload": art}, "article.html")

		} else {
			// If the article is not found, abort with an error
			c.AbortWithError(http.StatusNotFound, err)
		}

	} else {
		// If an invalid article ID is specified in the URL, abort with an error
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func getArticleByUser(c *gin.Context) {
	// Check if the article ID is valid
	if userID, err := strconv.Atoi(c.Param("user_id")); err == nil {
		fmt.Println(userID)
		// Check if the article exists
		if art := getArticleByUserID(userID); len(art) != 0 {
			// Call the render function with the title, article and the name of the
			// template
			c.SecureJSON(http.StatusOK, art)
		} else {
			// If the article is not found, abort with an error
			fmt.Println(art)
			c.AbortWithError(http.StatusNotFound, err)
		}

	} else {
		// If an invalid article ID is specified in the URL, abort with an error
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func createArticle(c *gin.Context) {
	// Obtain the POSTed title and content values
	title := c.PostForm("title")
	url := c.PostForm("url")
	rating, err := strconv.Atoi(c.PostForm("rating"))
	if err == nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	user_id, err := strconv.Atoi(c.PostForm("user_id"))
	if err == nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	var body article
	c.ShouldBindJSON(&body)
	if (article{} == body) {
		if a, err := createNewArticle(title, url, rating, user_id); err == nil {
			// If the article is created successfully, show success message
			render(c, gin.H{
				"title":   "Submission Successful",
				"payload": a}, "submission-successful.html")
		} else {
			// if there was an error while creating the article, abort with an error
			c.AbortWithStatus(http.StatusBadRequest)
		}
	} else {
		fmt.Println("rating", body.Rating)
		fmt.Println("user_id", body.User_id)
		if a, err := createNewArticle(body.Title, body.Url, body.Rating, body.User_id); err == nil {
			// If the article is created successfully, show success message
			c.SecureJSON(http.StatusOK, a)
		} else {
			// if there was an error while creating the article, abort with an error
			c.AbortWithStatus(http.StatusBadRequest)
		}
	}
}

func similaritiesArticle(c *gin.Context) {
	var body URLstring
	//URL := c.PostForm("url")
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mainArticle := processArticleByUrl(body.Url)

	querySearchGoogle := strings.Join(append_words(
		[]string{mainArticle.WordCleanText, mainArticle.RankedPhraseCleanText},
		[]string{mainArticle.WordTitle}), " ")

	var filteredResults []string
	searchResults, _ := googlesearch.Search(c, querySearchGoogle)
	for _, searchResult := range searchResults {
		u, err := url.ParseRequestURI(searchResult.URL)
		if err == nil &&
			strings.Contains(u.String(), ".html") &&
			u.String() != body.Url {
			filteredResults = append(filteredResults, u.String())
		}
	}

	var otherArticlesFound []ArticleProccesed
	for _, value := range filteredResults {
		article := processArticleByUrl(value)
		if len(article.CleanText) > 2000 {
			otherArticlesFound = append(otherArticlesFound, article)
		}
	}

	articleWiki, err := wiki_api(mainArticle.WordCleanText)
	if err != nil {
		fmt.Println(err)
	}

	articleWiki.Similarities = compare_articles(mainArticle.CleanText, articleWiki.Content)

	for i := 0; i < len(otherArticlesFound); i++ {
		otherArticlesFound[i].Similarities = compare_articles(mainArticle.CleanText, otherArticlesFound[i].CleanText)
	}

	response := ResponseArticleCompared{
		QueryUsedForSearch: querySearchGoogle,
		MainArticle:        mainArticle,
		WikiArticle:        articleWiki,
		OtherArticlesFound: otherArticlesFound,
	}

	c.SecureJSON(http.StatusOK, response)
}
