package main

import (
	"fmt"
	"strings"

	textrank "github.com/DavidBelicza/TextRank/v2"
	"github.com/DavidBelicza/TextRank/v2/convert"
	"github.com/DavidBelicza/TextRank/v2/parse"
	"github.com/DavidBelicza/TextRank/v2/rank"
	"github.com/jaytaylor/html2text"
)

var (
	keywordReplacer = strings.NewReplacer(
		"/", " ",
		".", " ",
		"-", " ",
		":", " ",
	)
)

type Extractor struct {
	html2textOptions html2text.Options

	textRankLanguage  convert.Language
	textRankRule      parse.Rule
	textRankAlgorithm rank.Algorithm
}

func NewExtractor() *Extractor {
	return &Extractor{
		html2textOptions: html2text.Options{
			OmitLinks: true,
			TextOnly:  true,
		},
		textRankLanguage:  textrank.NewDefaultLanguage(),
		textRankAlgorithm: textrank.NewDefaultAlgorithm(),
		textRankRule:      textrank.NewDefaultRule(),
	}
}

func (e *Extractor) ExtractKeyPhrases(htmlDescription string) ([]string, error) {
	description, err := html2text.FromString(htmlDescription, e.html2textOptions)
	if err != nil {
		return []string{}, err
	}

	description = normalizeText(description)

	tr := textrank.NewTextRank()
	tr.Populate(description, e.textRankLanguage, e.textRankRule)
	tr.Ranking(e.textRankAlgorithm)

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

func normalizeText(text string) string {
	text = strings.TrimSpace(text)
	text = keywordReplacer.Replace(text)

	return text
}