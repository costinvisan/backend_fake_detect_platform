package main

import (
	"github.com/masatana/go-textdistance"
)

type Similarity struct {
	Simi         int
	MainArticle  string
	OtherArticle string
}

func compare_articles_text(mainArticle []string, toCompare string) []Similarity {
	var importantSimilarities []Similarity
	// oc := metrics.NewOverlapCoefficient()
	// oc.CaseSensitive = false
	// oc.NgramSize = 3
	for i := 0; i < len(mainArticle); i++ {
		sim := textdistance.DamerauLevenshteinDistance(mainArticle[i], toCompare)
		importantSimilarities =
			append(importantSimilarities,
				Similarity{
					Simi:         sim,
					MainArticle:  toCompare,
					OtherArticle: mainArticle[i],
				})

	}

	return importantSimilarities
}
