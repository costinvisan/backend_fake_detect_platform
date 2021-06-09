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
	articles := getAllArticles()

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
		if article, err := getArticleByID(articleID); err == nil {
			// Call the render function with the title, article and the name of the
			// template
			render(c, gin.H{
				"title":   article.Title,
				"payload": article}, "article.html")

		} else {
			// If the article is not found, abort with an error
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
	content := c.PostForm("content")

	if a, err := createNewArticle(title, content); err == nil {
		// If the article is created successfully, show success message
		render(c, gin.H{
			"title":   "Submission Successful",
			"payload": a}, "submission-successful.html")
	} else {
		// if there was an error while creating the article, abort with an error
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func similaritiesArticle(c *gin.Context) {
	var body URLstring

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mainArticle := processArticleByUrl(body.Url)

	keyWordsCleanText := []string{mainArticle.WordCleanText, mainArticle.RankedPhraseCleanText}
	keyWordsTitle := []string{mainArticle.WordTitle}

	queryArray := append_words(keyWordsTitle, keyWordsCleanText)

	querySearchGoogle := strings.Join(queryArray, " ")

	fmt.Println(querySearchGoogle)

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
	fmt.Println(filteredResults)
	var otherArticlesFound []ArticleProccesed
	for _, value := range filteredResults {
		otherArticlesFound = append(otherArticlesFound, processArticleByUrl(value))
	}

	articleWiki, err := wiki_api(mainArticle.WordCleanText)
	if err != nil {
		fmt.Println(err)
	}

	wikiSimilarities := compare_articles(mainArticle.CleanText, articleWiki.Content)
	articleWiki.Similarities = wikiSimilarities

	var otherArticlesSimilarities [][]string
	for _, value := range otherArticlesFound {
		otherArticlesSimilarities = append(otherArticlesSimilarities, compare_articles(mainArticle.CleanText, value.CleanText))
	}

	for i := 0; i < len(otherArticlesSimilarities); i++ {
		otherArticlesFound[i].Similarities = otherArticlesSimilarities[i]
	}

	response := ResponseArticleCompared{
		QueryUsedForSearch: querySearchGoogle,
		MainArticle:        mainArticle,
		WikiArticle:        articleWiki,
		OtherArticlesFound: otherArticlesFound,
	}

	c.SecureJSON(http.StatusOK, response)
}
