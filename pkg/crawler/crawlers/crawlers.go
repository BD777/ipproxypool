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
	Crawl() <-chan IPProxyItem
}
