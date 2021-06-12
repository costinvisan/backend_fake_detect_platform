package main

import (
	"fmt"

	"github.com/james-bowman/nlp"
	"github.com/james-bowman/nlp/measures/pairwise"
	"gonum.org/v1/gonum/mat"
)

func compare_articles_text(mainArticle []string, toCompare string) []string {

	vectoriser := nlp.NewCountVectoriser(stopWords...)
	transformer := nlp.NewTfidfTransformer()

	// set k (the number of dimensions following truncation) to 4
	reducer := nlp.NewTruncatedSVD(4)

	lsiPipeline := nlp.NewPipeline(vectoriser, transformer, reducer)

	// Transform the corpus into an LSI fitting the model to the documents in the process
	lsi, err := lsiPipeline.FitTransform(mainArticle...)
	if err != nil {
		fmt.Printf("Failed to process documents because %v", err)
	}

	// run the query through the same pipeline that was fitted to the corpus and
	// to project it into the same dimensional space
	queryVector, err := lsiPipeline.Transform(toCompare)
	if err != nil {
		fmt.Printf("Failed to process documents because %v", err)
	}

	// iterate over document feature vectors (columns) in the LSI matrix and compare
	// with the query vector for similarity.  Similarity is determined by the difference
	// between the angles of the vectors known as the cosine similarity
	highestSimilarity := 0.992
	var similarity float64
	_, docs := lsi.Dims()
	var importantSimilarities []string
	for i := 0; i < docs; i++ {
		similarity = pairwise.CosineSimilarity(queryVector.(mat.ColViewer).ColView(0), lsi.(mat.ColViewer).ColView(i))
		if similarity > highestSimilarity {
			importantSimilarities = append(importantSimilarities, mainArticle[i])
		}
	}

	return importantSimilarities
}
