package main

import (
	"log"
	"sort"
	"strings"

	textrank "github.com/DavidBelicza/TextRank"
	goose "github.com/advancedlogic/GoOse"
	"github.com/jdkato/prose/v2"
)

type ArticleProccesed struct {
	Url                   string
	Title                 string
	CleanText             string
	WordCleanText         string
	WordTitle             string
	RankedPhraseCleanText string
	RankedSentences       []string
	Similarities          []Similarity
}

type ResponseArticleCompared struct {
	QueryUsedForSearch string
	MainArticle        ArticleProccesed
	WikiArticle        WikiArticle
	OtherArticlesFound []ArticleProccesed
}

func append_words(a []string, b []string) []string {

	check := make(map[string]int)
	d := append(a, b...)
	res := make([]string, 0)
	for _, val := range d {
		check[val] = 1
	}

	for letter := range check {
		res = append(res, letter)
	}

	return res
}

func uniqueNonEmptyElementsOf(s []Similarity) []Similarity {
	unique := make(map[string]bool, len(s))
	var us []Similarity
	for _, elem := range s {
		if len(elem.OtherArticle) != 0 {
			if !unique[elem.OtherArticle] {
				us = append(us, elem)
				unique[elem.OtherArticle] = true
			}
		}
	}
	return us
}

func processArticleByUrl(url string) ArticleProccesed {
	g := goose.New()
	article, _ := g.ExtractFromURL(url)

	doc, err := prose.NewDocument(strings.ReplaceAll(article.CleanedText, "\n", ""))
	if err != nil {
		log.Fatal(err)
	}
	var sents []string
	for _, value := range doc.Sentences() {
		sents = append(sents, value.Text)
	}

	cleanText := strings.Join(sents, "")

	title := article.Title
	// TextRank object
	tr1 := textrank.NewTextRank()
	tr2 := textrank.NewTextRank()
	// Default Rule for parsing.
	rule := textrank.NewDefaultRule()
	// Default Language for filtering stop words.
	language := textrank.NewDefaultLanguage()
	// Default algorithm for ranking text.
	algorithm := textrank.NewChainAlgorithm()

	// Add text.
	tr1.Populate(cleanText, language, rule)
	tr2.Populate(title, language, rule)

	// Run the ranking.
	tr1.Ranking(algorithm)
	tr2.Ranking(algorithm)

	// Get all phrases order by weight.
	rankedPhraseCleanText := textrank.FindPhrases(tr1)[0]

	// Get all words order by weight.
	wordCleanText := textrank.FindSingleWords(tr1)[0]
	wordTitle := textrank.FindSingleWords(tr2)[0]

	sentences := textrank.FindSentencesByRelationWeight(tr1, 10)

	var rankedSentences []string
	for _, value := range sentences {
		rankedSentences = append(rankedSentences, value.Value)
	}

	return ArticleProccesed{
		Url:           url,
		Title:         title,
		CleanText:     cleanText,
		WordCleanText: wordCleanText.Word,
		WordTitle:     wordTitle.Word,
		RankedPhraseCleanText: rankedPhraseCleanText.Right + " " +
			rankedPhraseCleanText.Left,
		RankedSentences: rankedSentences,
	}
}

func compare_articles(mainArticle string, toCompareArticle string) []Similarity {
	mainArticleArray := strings.Split(mainArticle, ".")
	toCompareArticleArray := strings.Split(toCompareArticle, ".")

	var similarities []Similarity
	for _, value := range mainArticleArray {
		similarities = append(similarities, compare_articles_text(toCompareArticleArray, value)...)
	}
	map_sim := make(map[string]Similarity)
	for _, value := range mainArticleArray {
		for _, sim := range similarities {
			if value == sim.MainArticle {
				if _, ok := map_sim[value]; ok {
					if sim.Simi < map_sim[value].Simi {
						map_sim[value] = sim
					}
				} else {
					map_sim[value] = sim
				}
			}
		}
	}
	var values_map []Similarity
	for _, val := range map_sim {
		values_map = append(values_map, val)
	}
	sort.SliceStable(values_map, func(i, j int) bool {
		return values_map[i].Simi < values_map[j].Simi
	})
	uniq := uniqueNonEmptyElementsOf(values_map)
	return uniq
}
