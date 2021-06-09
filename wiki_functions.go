package main

import (
	"encoding/json"
	"strings"

	"cgt.name/pkg/go-mwclient"
	"github.com/jmoiron/jsonq"
)

type WikiArticle struct {
	Title        string
	Content      string
	Similarities []string
}

func wiki_api(toSearch string) (WikiArticle, error) {
	client, err := mwclient.New("https://en.wikipedia.org/w/api.php", "myWikibot")
	if err != nil {
		return WikiArticle{}, err
	}

	queryForSearch := map[string]string{
		"action":   "query",
		"format":   "json",
		"list":     "search",
		"srsearch": toSearch,
	}

	body, err := client.Get(queryForSearch)
	if err != nil {
		return WikiArticle{}, err
	}

	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(body.String()))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)

	queryObject, err := jq.Object("query")
	jq = jsonq.NewQuery(queryObject)

	searchArray, err := jq.Array("search")
	jq = jsonq.NewQuery(searchArray[0])

	pageTitle, err := jq.String("title")
	if err != nil {
		return WikiArticle{}, err
	}

	queryForReadPage := map[string]string{
		"action":        "query",
		"prop":          "extracts",
		"exsentences":   "10",
		"exlimit":       "1",
		"titles":        pageTitle,
		"explaintext":   "1",
		"formatversion": "2",
	}

	body, err = client.Get(queryForReadPage)
	if err != nil {
		return WikiArticle{}, err
	}

	dec = json.NewDecoder(strings.NewReader(body.String()))
	dec.Decode(&data)
	jq = jsonq.NewQuery(data)

	queryObject, err = jq.Object("query")
	jq = jsonq.NewQuery(queryObject)

	pagesArray, err := jq.Array("pages")
	jq = jsonq.NewQuery(pagesArray[0])

	title, err := jq.String("title")
	content, err := jq.String("extract")

	return WikiArticle{Title: title, Content: content}, nil
}
