package main

import (
	"fmt"
	"strings"

	textrank "github.com/DavidBelicza/TextRank/v2"
	anyascii "github.com/anyascii/go"
)

var (
	textRankLanguage  = textrank.NewDefaultLanguage()
	textRankAlgorithm = textrank.NewDefaultAlgorithm()
	textRankRule      = textrank.NewDefaultRule()

	keywordReplacer = strings.NewReplacer(
		"/", " ",
		"-", " ",
		":", " ",
		"-", " ",
	)

	nItems = 10
)

func TextRankPhrases(text string) ([]string, []string, error) {
	text = anyascii.Transliterate(text)
	text = keywordReplacer.Replace(text)

	tr := textrank.NewTextRank()
	tr.Populate(text, textRankLanguage, textRankRule)
	tr.Ranking(textRankAlgorithm)

	// extract phrases
	rankedPhrases := textrank.FindPhrases(tr)

	nPhrases := nItems
	if len(rankedPhrases) < nItems {
		nPhrases = len(rankedPhrases)
	}

	phrases := make([]string, nPhrases)

	for i := 0; i < nPhrases; i++ {
		phrases[i] = fmt.Sprintf("%s %s", rankedPhrases[i].Left, rankedPhrases[i].Right)
	}

	// extract single words
	rankedWords := textrank.FindSingleWords(tr)

	nWords := nItems
	if len(rankedWords) < nItems {
		nWords = len(rankedWords)
	}

	words := make([]string, nWords)

	for i := 0; i < nWords; i++ {
		words[i] = rankedWords[i].Word
	}

	return phrases, words, nil
}
