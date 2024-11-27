package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/mmcdole/gofeed"
)

var (
	xmlDir = "xml"
)

func main() {
	files, err := os.ReadDir(xmlDir)
	if err != nil {
		log.Fatal(err)
	}

	fp := gofeed.NewParser()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "File\tTitle\tDescription")

	for _, fileEntry := range files {
		file, err := os.Open(filepath.Join(xmlDir, fileEntry.Name()))
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		feed, err := fp.Parse(file)
		if err != nil {
			log.Println(fileEntry.Name())
			log.Fatal(err)
		}

		fmt.Fprintf(
			w,
			"%s\t%s\t%s\n",
			fileEntry.Name(),
			feed.Title,
			feed.Description,
		)
	}

	w.Flush()
}
