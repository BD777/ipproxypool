package crawlers

import (
	"testing"
)

func TestCralwerIP3399_Crawl(t *testing.T) {
	t.Logf("TestCralwerIP3399_Crawl")
	c := &CralwerIP3399{}
	ch := c.Crawl()
	for item := range ch {
		t.Logf("%+v", item)
	}
}
