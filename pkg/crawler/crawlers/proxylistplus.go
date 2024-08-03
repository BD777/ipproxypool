package crawlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/BD777/ipproxypool/pkg/utils/htmlparser"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/levigross/grequests"
	"github.com/sirupsen/logrus"
)

type CrawlerProxyListPlus struct {
	session *grequests.Session
}

func NewCrawlerProxyListPlus() *CrawlerProxyListPlus {
	return &CrawlerProxyListPlus{}
}

func (c *CrawlerProxyListPlus) Name() string {
	return "proxylistplus"
}

func (c *CrawlerProxyListPlus) Crawl() <-chan IPProxyItem {
	const MaxPage = 6

	ch := make(chan IPProxyItem, 100)

	go func() {
		defer close(ch)
		for page := 1; page <= MaxPage; page++ {
			if page > 1 {
				time.Sleep(time.Second * 3) // avoid anti-crawler
			}
			items, err := c.crawlPage(page)
			if err != nil {
				logrus.Errorf("failed to crawl page %d: %v", page, err)
				return
			}
			if len(items) == 0 {
				break
			}

			for _, item := range items {
				ch <- item
			}
		}
	}()

	return ch
}

func (c *CrawlerProxyListPlus) Detect() bool {
	for page := 1; page <= 2; page++ {
		if page > 1 {
			time.Sleep(time.Second * 3) // avoid anti-crawler
		}

		resp, err := c.crawlPage(page)
		if err != nil {
			logrus.Errorf("failed to detect proxylistplus: %v", err)
			return false
		}
		if len(resp) == 0 {
			logrus.Errorf("failed to detect proxylistplus: no items in page %d", page)
			return false
		}
		logrus.Infof("detected %d items in page %d", len(resp), page)
	}
	return true
}

func (c *CrawlerProxyListPlus) newSession() {
	c.session = grequests.NewSession(&grequests.RequestOptions{
		UserAgent: browser.Chrome(),
	})
}

func (c *CrawlerProxyListPlus) crawlPage(page int) ([]*ProxyListPlusItem, error) {
	if c.session == nil {
		c.newSession()
	}

	logrus.Infof("[CrawlerProxyListPlus] start to crawl page %d", page)

	url := fmt.Sprintf("https://list.proxylistplus.com/Fresh-HTTP-Proxy-List-%d", page)
	httpResp, err := c.session.Get(url, nil)
	if err != nil {
		return nil, err
	}
	if !httpResp.Ok {
		return nil, err
	}

	resp := &ProxyListPlusResponse{}
	err = htmlparser.ParseHTML(httpResp.String(), resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}

	return resp.List, nil
}

type ProxyListPlusResponse struct {
	List []*ProxyListPlusItem `xpath:"//table[@class='bg']/tbody/tr[@class='cells']"`
}

type ProxyListPlusItem struct {
	IP      string `xpath:"td[2]/text()"`
	Port    string `xpath:"td[3]/text()"`
	Type    string `xpath:"td[4]/text()"`
	Country string `xpath:"td[5]/text()"`
}

func (p *ProxyListPlusItem) GetSource() string {
	return "proxylistplus"
}

func (p *ProxyListPlusItem) GetIP() string {
	return strings.TrimSpace(p.IP)
}

func (p *ProxyListPlusItem) GetPort() int {
	port, err := strconv.Atoi(strings.TrimSpace(p.Port))
	if err != nil {
		logrus.Errorf("failed to convert port %s to int: %v", p.Port, err)
		return 0
	}
	return port
}

func (p *ProxyListPlusItem) GetCountry() string {
	return p.Country
}

func (p *ProxyListPlusItem) GetRegion() string {
	return ""
}

func (p *ProxyListPlusItem) GetISP() string {
	return ""
}

func (p *ProxyListPlusItem) GetUpdatedAt() int64 {
	return 0
}
