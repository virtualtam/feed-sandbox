package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gosimple/slug"
	"github.com/mmcdole/gofeed"
)

const (
	xmlDir = "xml"
)

type Client struct {
	httpClient *http.Client
	feedParser *gofeed.Parser
}

func NewClient() *Client {
	return &Client{
		httpClient: http.DefaultClient,
		feedParser: gofeed.NewParser(),
	}
}

func (c *Client) Fetch(feedURL string) error {
	if err := os.MkdirAll(xmlDir, 0755); err != nil {
		return err
	}

	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		ce := resp.Body.Close()
		if ce != nil {
			err = ce
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	r := bytes.NewReader(body)

	parsedFeed, err := c.feedParser.Parse(r)
	if err != nil {
		return err
	}

	fileSlug := slug.Make(parsedFeed.Title)
	filePath := filepath.Join(
		xmlDir,
		fmt.Sprintf("%s.xml", fileSlug),
	)

	return os.WriteFile(filePath, body, 0644)
}
