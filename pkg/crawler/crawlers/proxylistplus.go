package crawlers

import (
	"fmt"
	"net/url"
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
	if c.session == nil {
		c.newSession()
	}

	k := NewKuaiDaiLi()
	proxies := k.Crawl()

	for proxy := range proxies {
		httpProxy, err := url.Parse(fmt.Sprintf("http://%s:%d", proxy.GetIP(), proxy.GetPort()))
		if err != nil {
			logrus.Errorf("failed to parse http proxy: %v", err)
			continue
		}
		httpsProxy, err := url.Parse(fmt.Sprintf("https://%s:%d", proxy.GetIP(), proxy.GetPort()))
		if err != nil {
			logrus.Errorf("failed to parse https proxy: %v", err)
			continue
		}

		c.session.RequestOptions.Proxies = map[string]*url.URL{
			"http":  httpProxy,
			"https": httpsProxy,
		}
		c.session.RequestOptions.DialTimeout = time.Second * 5

		resp, err := c.getWith("https://httpbin.org/ip", nil)
		if err != nil {
			logrus.Errorf("Failed to get IP via proxy: %v", err)
		} else {
			logrus.Infof("Response from IP check: %v", resp.String())
		}

		ok := func() bool {
			for page := 1; page <= 2; page++ {
				if page > 1 {
					time.Sleep(time.Second * 3) // avoid anti-crawler
				}

				resp, err := c.crawlPage(page)
				if err != nil {
					logrus.Errorf("failed to detect proxylistplus with proxy %s:%d: %v", proxy.GetIP(), proxy.GetPort(), err)
					return false
				}
				if len(resp) == 0 {
					logrus.Errorf("failed to detect proxylistplus with proxy %s:%d: no items in page %d", proxy.GetIP(), proxy.GetPort(), page)
					return false
				}
				logrus.Infof("detected %d items in page %d with proxy %s:%d", len(resp), page, proxy.GetIP(), proxy.GetPort())
			}
			logrus.Infof("proxy %s:%d is able to fetch proxies", proxy.GetIP(), proxy.GetPort())
			return true
		}()
		if !ok {
			logrus.Infof("proxy %s:%d is not able to fetch proxies", proxy.GetIP(), proxy.GetPort())
		} else {
			return true
		}
	}

	logrus.Infof("no proxy is able to fetch proxies")
	return false
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
	httpResp, err := c.getWith(url, nil)
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

func (c *CrawlerProxyListPlus) getWith(url string, ro *grequests.RequestOptions) (*grequests.Response, error) {
	if c.session == nil {
		c.newSession()
	}

	// it seems that grequest.session.Get is not working with proxy settings
	var getFunc func(url string, ro *grequests.RequestOptions) (*grequests.Response, error)

	if ro == nil {
		ro = &grequests.RequestOptions{}
	}
	if len(ro.Proxies) == 0 && len(c.session.RequestOptions.Proxies) > 0 {
		logrus.Infof("use session proxies: %v", c.session.RequestOptions.Proxies)
		ro.Proxies = c.session.RequestOptions.Proxies
		getFunc = grequests.Get
	} else {
		getFunc = c.session.Get
	}

	return getFunc(url, ro)
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
