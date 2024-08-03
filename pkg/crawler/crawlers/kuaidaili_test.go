package crawlers

import (
	"testing"
)

func TestCralweKuaiDaiLi_Crawl(t *testing.T) {
	t.Logf("TestCralweKuaiDaiLi_Crawl")
	c := &CrawlerKuaiDaiLi{}
	ch := c.Crawl()
	for item := range ch {
		t.Logf("source:%s, ip:%s, port:%d, country:%s, region:%s, isp:%s, updatedAt:%d",
			item.GetSource(), item.GetIP(), item.GetPort(), item.GetCountry(), item.GetRegion(), item.GetISP(), item.GetUpdatedAt())
	}
}
