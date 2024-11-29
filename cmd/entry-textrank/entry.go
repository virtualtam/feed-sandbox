package main

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/jaytaylor/html2text"
	"github.com/mmcdole/gofeed"
)

var (
	authorSeparatorRegexp = regexp.MustCompile("(,| and )")
	paragraphSplitRegexp  = regexp.MustCompile(`\n\s*\n`)

	html2textOptions = html2text.Options{
		OmitLinks: true,
		TextOnly:  true,
	}
)

type Entry struct {
	URL         string
	Title       string
	Authors     []string
	Description string
	Content     string

	// Computed fields
	Summary            string
	DescriptionPhrases []string
	ContentPhrases     []string
	SummaryPhrases     []string
}

func NewEntryFromItem(item *gofeed.Item) (Entry, error) {
	var err error

	var authorNames []string

	for _, author := range item.Authors {
		names := authorSeparatorRegexp.Split(author.Name, -1)

		for _, name := range names {
			name := strings.TrimSpace(name)
			if name == "" {
				continue
			}

			authorNames = append(authorNames, name)
		}
	}

	var description, content string

	if item.Description != "" {
		description, err = html2text.FromString(item.Description, html2textOptions)
		if err != nil {
			return Entry{}, err
		}

		description = strings.TrimSpace(description)
	}

	if item.Content != "" {
		content, err = html2text.FromString(item.Content, html2textOptions)
		if err != nil {
			return Entry{}, err
		}

		content = strings.TrimSpace(content)
	}

	entry := Entry{
		URL:         item.Link,
		Title:       item.Title,
		Authors:     authorNames,
		Description: description,
		Content:     content,
	}

	entry.Summarize()
	entry.ComputePhrases()

	return entry, nil
}

// buildSummaryFromParagraphs builds a summary from multiple paragraphs
// keeping it under maxLength characters
func buildSummaryFromParagraphs(paragraphs []string, maxLength int) string {
	var summary strings.Builder

	for i, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		// Add separator between paragraphs
		if summary.Len() > 0 {
			summary.WriteString(" ")
		}

		// Check if adding this paragraph would exceed maxLength
		if summary.Len()+len(p) > maxLength {
			// If this is the first paragraph, take what we can
			if i == 0 {
				return p[:maxLength-3] + "..."
			}
			// Otherwise stop here
			break
		}

		summary.WriteString(p)
	}

	return summary.String()
}

func (e *Entry) Summarize() {
	const (
		shortLength = 100 // Length to consider text "short enough" as is
		maxLength   = 200 // Maximum length for multi-paragraph summary
	)

	// early return: nothing to summarize
	if e.Description == "" && e.Content == "" {
		return
	}

	// Try Description first
	if e.Description != "" {
		if len(e.Description) <= shortLength {
			e.Summary = e.Description
			return
		}

		paragraphs := paragraphSplitRegexp.Split(e.Description, -1)
		if len(paragraphs) > 0 {
			e.Summary = buildSummaryFromParagraphs(paragraphs, maxLength)
			return
		}
	}

	// Fall back to Content if Description didn't yield a summary
	if e.Content != "" {
		if len(e.Content) <= shortLength {
			e.Summary = e.Content
			return
		}

		paragraphs := paragraphSplitRegexp.Split(e.Content, -1)
		if len(paragraphs) > 0 {
			e.Summary = buildSummaryFromParagraphs(paragraphs, maxLength)
			return
		}
	}
}

func (e *Entry) ComputePhrases() error {
	var err error

	if e.Description != "" {
		e.DescriptionPhrases, err = TextRankPhrases(e.Description)
		if err != nil {
			return err
		}
	}

	if e.Content != "" {
		e.ContentPhrases, err = TextRankPhrases(e.Content)
		if err != nil {
			return err
		}
	}

	if e.Summary != "" {
		e.SummaryPhrases, err = TextRankPhrases(e.Summary)
		if err != nil {
			return err
		}
	}

	return nil
}

func writePhrases(w io.Writer, label string, phrases []string) {
	if len(phrases) == 0 {
		return
	}
	fmt.Fprintf(w, "%s phrases:\t%s\n", label, strings.Join(phrases, ", "))
}

func (e *Entry) WriteInfo(output io.Writer) {
	fmt.Fprintln(output, strings.Repeat("-", 80))

	fmt.Fprintf(output, "Title:\t\t%s\n", e.Title)
	fmt.Fprintf(output, "URL:\t\t%s\n", e.URL)

	if len(e.Authors) > 0 {
		fmt.Fprintf(output, "Authors:\t%s\n", strings.Join(e.Authors, ", "))
	}

	if e.Summary != "" {
		fmt.Fprintf(output, "\nSummary (%d chars): %s\n", len(e.Summary), e.Summary)
	}

	if e.Description != "" {
		fmt.Fprintf(output, "\nDescription (%d chars)\n", len(e.Description))
		// fmt.Fprintln(output, e.Description)
	}

	if e.Content != "" {
		fmt.Fprintf(output, "\nContent (%d chars)\n", len(e.Content))
		// fmt.Fprintln(output, e.Content)
	}

	fmt.Fprintln(output)

	writePhrases(output, "Summary", e.SummaryPhrases)
	writePhrases(output, "Description", e.DescriptionPhrases)
	writePhrases(output, "Content", e.ContentPhrases)

	if e.Description == "" && e.Content == "" {
		fmt.Fprintln(output, "No text content available")
	}
}

type Entries []Entry

func (es *Entries) WriteInfo(output io.Writer) {
	for _, entry := range *es {
		entry.WriteInfo(output)
	}
}
