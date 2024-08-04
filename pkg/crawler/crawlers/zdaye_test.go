package crawlers

import (
	"testing"

	"github.com/levigross/grequests"
)

func TestCralwerZdaye_Crawl(t *testing.T) {
	t.Logf("TestCralwerZdaye_Crawl")
	c := &CrawlerZdaye{}
	ch := c.Crawl()
	for item := range ch {
		t.Logf("source:%s, ip:%s, port:%d, country:%s, region:%s, isp:%s, updatedAt:%d",
			item.GetSource(), item.GetIP(), item.GetPort(), item.GetCountry(), item.GetRegion(), item.GetISP(), item.GetUpdatedAt())
	}
}

func TestCrawlerZdaye_Detect(t *testing.T) {
	type fields struct {
		session *grequests.Session
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "TestCrawlerZdaye_Detect",
			fields: fields{},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CrawlerZdaye{
				session: tt.fields.session,
			}
			if got := c.Detect(); got != tt.want {
				t.Errorf("CrawlerZdaye.Detect() = %v, want %v", got, tt.want)
			}
		})
	}
}
