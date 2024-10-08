package crawlers

import (
	"testing"
)

func TestCralwerKuaiDaiLi_Crawl(t *testing.T) {
	t.Logf("TestCralwerKuaiDaiLi_Crawl")
	c := &CrawlerKuaiDaiLi{}
	ch := c.Crawl()
	for item := range ch {
		t.Logf("source:%s, ip:%s, port:%d, country:%s, region:%s, isp:%s, updatedAt:%d",
			item.GetSource(), item.GetIP(), item.GetPort(), item.GetCountry(), item.GetRegion(), item.GetISP(), item.GetUpdatedAt())
	}
}
