package main

import (
	"cmp"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/sourcegraph/conc/pool"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <feed.xml>\n", os.Args[0])
		os.Exit(1)
	}

	xmlFile := os.Args[1]
	file, err := os.Open(xmlFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	start := time.Now()

	// Parse feed
	feedParser := gofeed.NewParser()
	feed, err := feedParser.Parse(file)
	if err != nil {
		log.Fatal(err)
	}

	nWorkers := cmp.Or(runtime.NumCPU()/2, 2)
	workerPool := pool.New().WithErrors().WithMaxGoroutines(nWorkers)

	entries := make(Entries, len(feed.Items))
	var entriesLocker sync.Mutex

	for i, item := range feed.Items {
		workerPool.Go(func() error {
			entry, err := NewEntryFromItem(item)
			if err != nil {
				return err
			}

			entriesLocker.Lock()
			entries[i] = entry
			entriesLocker.Unlock()

			return nil
		})
	}

	if err := workerPool.Wait(); err != nil {
		log.Fatal(err)
	}

	entries.WriteInfo(os.Stdout)

	elapsed := time.Since(start)
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("Processed %d entries in %s\n", len(feed.Items), elapsed)
}
