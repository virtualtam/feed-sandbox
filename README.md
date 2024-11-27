# feed-sandbox

A sandbox project to experiment with extracting and transforming metadata for Atom and RSS feeds
using Go libraries.

## Usage
### `feeds.csv`
This file contains a list of Atom or RSS feed URLs, that will be fetched and saved in the `xml`
directory for further usage.

Format:


```csv
feed_url
https://domain.tld/feed
https://blog.domain2.tld/rss
```

### `fetch` - Download feed data

```shell
$ go run ./cmd/fetch
```

### `feed-meta` - Print feed metadata

```shell
$ go run ./cmd/feed-meta
```

### `entry-meta` - Print entry metadata for a feed

```shell
$ go run ./cmd/entry-meta xml/myfeed.xml
```

## LICENSE
`feed-sandbox` is licensed under the MIT license.
