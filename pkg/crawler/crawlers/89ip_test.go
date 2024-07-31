package crawlers

import (
	"testing"
)

func TestCrawler89IP_crawlPage(t *testing.T) {
	t.Logf("TestCrawler89IP_crawlPage")
	c := &Crawler89IP{}
	ch := c.Crawl()
	for item := range ch {
		t.Logf("%+v, updatedAt:%d", item, item.GetUpdatedAt())
	}
}
