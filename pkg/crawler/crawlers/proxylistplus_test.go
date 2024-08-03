package crawlers

import (
	"testing"

	"github.com/levigross/grequests"
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

func TestCrawlerProxyListPlus_Detect(t *testing.T) {
	type fields struct {
		session *grequests.Session
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "TestCrawlerProxyListPlus_Detect",
			fields: fields{},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CrawlerProxyListPlus{
				session: tt.fields.session,
			}
			if got := c.Detect(); got != tt.want {
				t.Errorf("CrawlerProxyListPlus.Detect() = %v, want %v", got, tt.want)
			}
		})
	}
}
