package crawlers

import (
	"testing"
)

func TestCralwerProxyListPlus_Crawl(t *testing.T) {
	t.Logf("TestCralwerProxyListPlus_Crawl")
	c := &CrawlerProxyListPlus{}
	ch := c.Crawl()
	for item := range ch {
		t.Logf("source:%s, ip:%s, port:%d, country:%s, region:%s, isp:%s, updatedAt:%d",
			item.GetSource(), item.GetIP(), item.GetPort(), item.GetCountry(), item.GetRegion(), item.GetISP(), item.GetUpdatedAt())
	}
}
