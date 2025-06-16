package rss

import (
	"aggregator/pkg/posts"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
	Categories  []string `xml:"category"`
}

func Parse(url string) ([]posts.Post, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, resp.Status)
	}

	var feed Feed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("XML decode error: %w", err)
	}

	if len(feed.Channel.Items) == 0 {
		return nil, errors.New("no items found in RSS feed")
	}

	var result []posts.Post
	for _, item := range feed.Channel.Items {
		if item.Link == "" || item.Title == "" {
			continue
		}

		pubTime, err := parseDate(item.PubDate)
		if err != nil {
			continue
		}

		cleanDescription := html.UnescapeString(item.Description)
		cleanDescription = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(cleanDescription, "")

		result = append(result, posts.Post{
			Title:   item.Title,
			Content: cleanDescription,
			Link:    item.Link,
			PubTime: pubTime.Unix(),
		})
	}

	if len(result) == 0 {
		return nil, errors.New("no valid posts found after parsing")
	}

	return result, nil
}

func parseDate(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	formats := []string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC822,
		time.RFC822Z,
		time.RFC3339,
		"Mon, 2 Jan 2006 15:04:05 -0700",
	}

	for _, format := range formats {
		t, err := time.Parse(format, dateStr)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognized date format: %s", dateStr)
}

func ParseURL(url string, posts chan<- []posts.Post, errs chan<- error, period int) {
	for {
		news, err := Parse(url)
		if err != nil {
			errs <- fmt.Errorf("%s: %w", url, err)
			time.Sleep(time.Duration(period) * time.Minute)
			continue
		}
		posts <- news
		time.Sleep(time.Duration(period) * time.Minute)
	}
}
