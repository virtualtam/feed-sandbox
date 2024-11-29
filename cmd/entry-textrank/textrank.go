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
)

func TextRankPhrases(text string) ([]string, error) {
	text = anyascii.Transliterate(text)
	text = keywordReplacer.Replace(text)

	tr := textrank.NewTextRank()
	tr.Populate(text, textRankLanguage, textRankRule)
	tr.Ranking(textRankAlgorithm)

	rankedPhrases := textrank.FindPhrases(tr)

	phrases := make([]string, len(rankedPhrases))

	for i, rankedPhrase := range rankedPhrases {
		phrases[i] = fmt.Sprintf("%s %s", rankedPhrase.Left, rankedPhrase.Right)
	}

	nItems := 10
	if len(phrases) < nItems {
		nItems = len(phrases)
	}

	return phrases[:nItems], nil
}
