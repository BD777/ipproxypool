package crawlers

type IPProxyItem interface {
	GetSource() string
	GetIP() string
	GetPort() int
	GetCountry() string
	GetRegion() string
	GetISP() string
	GetUpdatedAt() int64
}

type IPProxyCrawler interface {
	// Name returns the name of the crawler.
	Name() string
	// Crawl returns a channel that yields IPProxyItem.
	Crawl() <-chan IPProxyItem
	// Detect returns true if the crawler is able to fetch proxies.
	Detect() bool
}

var Crawlers = []IPProxyCrawler{
	NewCrawler89IP(),
	NewCralwerIP3399(),
	NewKuaiDaiLi(),
	NewCrawlerProxyListPlus(),
	NewCrawlerZdaye(),
}
