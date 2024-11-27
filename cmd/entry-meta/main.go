package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	textrank "github.com/DavidBelicza/TextRank/v2"
	"github.com/DavidBelicza/TextRank/v2/convert"
	"github.com/DavidBelicza/TextRank/v2/parse"
	"github.com/DavidBelicza/TextRank/v2/rank"
	"github.com/jaytaylor/html2text"
	"github.com/mmcdole/gofeed"
)

var (
	maxLen = 100

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

func NewExtractor() Extractor {
	return Extractor{
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

	description = normalizeDescription(description)

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

func normalizeDescription(keyword string) string {
	keyword = strings.TrimSpace(keyword)
	keyword = keywordReplacer.Replace(keyword)

	return keyword
}

func main() {
	xmlFile := os.Args[1]

	file, err := os.Open(xmlFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fp := gofeed.NewParser()

	feed, err := fp.Parse(file)
	if err != nil {
		log.Fatal(err)
	}

	extractor := NewExtractor()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "Title\tAuthors\tDescLen\tConLen\tWords")

	for _, entry := range feed.Items {
		var authorNames []string

		for _, author := range entry.Authors {
			names := strings.Split(author.Name, " and ")

			for _, name := range names {
				name := strings.TrimSpace(name)
				if name == "" {
					continue
				}

				authorNames = append(authorNames, name)
			}
		}

		var phrases []string

		switch {
		case len(entry.Description) > 20:
			phrases, err = extractor.ExtractKeyPhrases(entry.Description)
			if err != nil {
				log.Fatal(err)
			}

		case len(entry.Content) > 20:
			phrases, err = extractor.ExtractKeyPhrases(entry.Content)
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Fprintf(
			w,
			"%s\t%s\t%d\t%d\t%s\n",
			entry.Title,
			strings.Join(authorNames, " & "),
			len(entry.Description),
			len(entry.Content),
			strings.Join(phrases, ", "),
		)
	}

	w.Flush()
}
