package main

import (
	"cmp"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/sourcegraph/conc/pool"
)

func main() {
	xmlFile := os.Args[1]

	file, err := os.Open(xmlFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	feedParser := gofeed.NewParser()
	extractor := NewExtractor()

	start := time.Now()

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
			entry, err := NewEntryFromItem(extractor, item)
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

	entries.Summary(os.Stdout)

	elapsed := time.Since(start)
	fmt.Printf(
		"%d entries processed in %d ms (%d jobs)\n",
		len(feed.Items),
		elapsed.Milliseconds(),
		nWorkers,
	)
}
