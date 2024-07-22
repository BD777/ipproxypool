package crawlers

import (
	"testing"
)

func TestCrawler89IP_crawlPage(t *testing.T) {
	t.Logf("TestCrawler89IP_crawlPage")
	c := &Crawler89IP{}
	resp := c.Crawl()
	for i, item := range resp {
		t.Logf("item %d: %+v", i, item)
	}
}
