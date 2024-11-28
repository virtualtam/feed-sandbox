# feed-sandbox

A toy project to experiment with extracting and transforming metadata for Atom and RSS feeds
using Go libraries.

These experiments are conducted with the goal of improving full-text search features in
[SparkleMuffin](https://github.com/virtualtam/sparklemuffin).

## Configuration
### `feeds.csv`
This file contains a list of Atom or RSS feed URLs, that will be fetched and saved in the `xml`
directory for further usage.

Format:

```csv
feed_url
https://domain.tld/feed
https://blog.domain2.tld/rss
```

## Usage
### `fetch` - Download feed data
Download XML data for all feed URLs in `feeds.csv`:

```shell
$ go run ./cmd/fetch
```

The resulting files can be found in the `xml` directory.

### `feed-meta` - Print feed metadata

Print feed metadata for all downloaded XML files:

```shell
$ go run ./cmd/feed-meta
```

### `entry-meta` - Print entry metadata for a feed

Print entry metadata for a given feed:

```shell
$ go run ./cmd/entry-meta xml/myfeed.xml
```

## LICENSE
`feed-sandbox` is licensed under the MIT license.

## Credits
This toy project uses the following libraries (in order of appearance):

- [mmcdole/gofeed](https://github.com/mmcdole/gofeed) to parse Atom, RSS and JSON feeds
- [gosimple/slug](https://github.com/gosimple/slug) to turn feed titles into normalized slugs that can be used to build filenames
- [sourcegraph/conc](https://github.com/sourcegraph/conc) to allocate a worker pool to process feed entries concurrently
- [jaytaylor/html2text](https://github.com/jaytaylor/html2text) to convert raw entry descriptions to plain text
- [DavidBelicza/TextRank](https://github.com/DavidBelicza/TextRank) to extract key phrases from entry descriptions
