package main

import (
	"strings"
)

var (
	keywordReplacer = strings.NewReplacer(
		"/", " ",
		"-", " ",
		":", " ",
		"-", " ",
	)
)

func normalizeText(text string) string {
	text = strings.TrimSpace(text)
	text = keywordReplacer.Replace(text)

	return text
}
