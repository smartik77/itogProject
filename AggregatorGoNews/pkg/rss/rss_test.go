package rss

import (
	"testing"
)

func TestParse(t *testing.T) {
	feed, err := Parse("https://habr.com/ru/rss/best/daily/?fl=ru")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(feed) == 0 {
		t.Fatal("No posts parsed")
	}

	t.Logf("Got %d posts", len(feed))
	for _, p := range feed[:3] {
		t.Logf("Post: %s\nLink: %s\nPubDate: %d\n",
			p.Title, p.Link, p.PubTime)
	}
}
