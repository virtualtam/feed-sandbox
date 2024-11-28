package main

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/mmcdole/gofeed"
)

type Entry struct {
	Title       string
	Authors     []string
	Description string
	Content     string
	Phrases     []string
}

func NewEntryFromItem(extractor *Extractor, item *gofeed.Item) (Entry, error) {
	var authorNames []string

	for _, author := range item.Authors {
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
	var err error

	switch {
	case len(item.Description) > 20:
		phrases, err = extractor.ExtractKeyPhrases(item.Description)
		if err != nil {
			return Entry{}, err
		}

	case len(item.Content) > 20:
		phrases, err = extractor.ExtractKeyPhrases(item.Content)
		if err != nil {
			return Entry{}, err
		}
	}

	return Entry{
		Title:       item.Title,
		Authors:     authorNames,
		Description: item.Description,
		Content:     item.Content,
		Phrases:     phrases,
	}, nil
}

type Entries []Entry

func (es *Entries) Summary(output io.Writer) {
	w := tabwriter.NewWriter(output, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "Title\tAuthors\tDescLen\tConLen\tWords")

	for _, entry := range *es {
		fmt.Fprintf(
			w,
			"%s\t%s\t%d\t%d\t%s\n",
			entry.Title,
			strings.Join(entry.Authors, " & "),
			len(entry.Description),
			len(entry.Content),
			strings.Join(entry.Phrases, ", "),
		)
	}

	w.Flush()
}
