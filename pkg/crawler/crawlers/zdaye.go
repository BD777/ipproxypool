package crawlers

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/BD777/ipproxypool/pkg/utils"
	"github.com/BD777/ipproxypool/pkg/utils/htmlparser"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/levigross/grequests"
	"github.com/sirupsen/logrus"
)

type CrawlerZdaye struct {
	session *grequests.Session
}

func NewCrawlerZdaye() *CrawlerZdaye {
	return &CrawlerZdaye{}
}

func (c *CrawlerZdaye) Name() string {
	return "zdaye"
}

func (c *CrawlerZdaye) Crawl() <-chan IPProxyItem {
	const MaxPage = 7

	ch := make(chan IPProxyItem, 100)

	go func() {
		defer close(ch)
		for page := 1; page <= MaxPage; page++ {
			if page > 1 {
				time.Sleep(time.Second * 10) // avoid anti-crawler
			}

			items, err := c.crawlPage(page)
			if err != nil {
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

func (c *CrawlerZdaye) Detect() bool {
	if c.session == nil {
		c.newSession()
	}

	// It seems that non-Chinese IP addresses will be blocked
	// use KuaiDaiLi to get Chinese IP addresses
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

		for page := 1; page <= 2; page++ {
			if page > 1 {
				time.Sleep(time.Second * 10) // avoid anti-crawler
			}

			resp, err := c.crawlPage(page)
			if err != nil {
				logrus.Errorf("failed to detect zdaye with proxy %s:%d: %v", proxy.GetIP(), proxy.GetPort(), err)
				continue
			}
			if len(resp) == 0 {
				logrus.Errorf("failed to detect zdaye with proxy %s:%d: no items in page %d", proxy.GetIP(), proxy.GetPort(), page)
				continue
			}
			logrus.Infof("detected %d items in page %d with proxy %s:%d", len(resp), page, proxy.GetIP(), proxy.GetPort())
		}
		logrus.Infof("proxy %s:%d is able to fetch proxies", proxy.GetIP(), proxy.GetPort())
		return true
	}

	logrus.Infof("no proxy is able to fetch proxies")
	return false
}

func (c *CrawlerZdaye) newSession() {
	c.session = grequests.NewSession(&grequests.RequestOptions{
		UserAgent: browser.Chrome(),
		Headers: map[string]string{
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
			"Accept-Encoding": "gzip, deflate, br, zstd",
			"Accept-Language": "zh-CN,zh;q=0.9",
			"Connection":      "keep-alive",
			"Host":            "www.zdaye.com",
		},
	})
}

func (c *CrawlerZdaye) crawlPage(page int) ([]*ZdayeItem, error) {
	logrus.Infof("[CrawlerZdaye] start to crawl page %d", page)

	if c.session == nil {
		c.newSession()
	}

	var url string
	if page == 1 {
		url = "https://www.zdaye.com/free/"
	} else {
		url = fmt.Sprintf("https://www.zdaye.com/free/%d/", page)
	}
	resp, err := c.session.Get(url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to crawl page %d: %w", page, err)
	}

	if resp.StatusCode != 200 {
		// request url
		logrus.Infof("url %s", resp.RawResponse.Request.URL)
		// request headers
		logrus.Infof("headers %v", resp.RawResponse.Request.Header)
		// response headers
		logrus.Infof("headers %v", resp.RawResponse.Header)
		// response status code
		logrus.Infof("status code %d", resp.StatusCode)

		logrus.Infof("response %s", utils.GBK2UTF8(resp.Bytes()))
	}

	items, err := parseZdaye(resp.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse zdaye: %w", err)
	}
	return items, nil
}

func parseZdaye(html string) ([]*ZdayeItem, error) {
	resp := &ZdayeResponse{}
	err := htmlparser.ParseHTML(html, resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}

	return resp.List, nil
}

type ZdayeResponse struct {
	List []*ZdayeItem `xpath:"//table[@id='ipc']/tbody/tr"`
}

type ZdayeItem struct {
	IP       string `xpath:"td[1]/text()"`
	Port     string `xpath:"td[2]/text()"`
	Location string `xpath:"td[4]/text()"`
}

func (i *ZdayeItem) GetSource() string {
	return "zdaye"
}

func (i *ZdayeItem) GetIP() string {
	return strings.TrimSpace(i.IP)
}

func (i *ZdayeItem) GetPort() int {
	port, err := strconv.Atoi(strings.TrimSpace(i.Port))
	if err != nil {
		logrus.Errorf("failed to convert port %s to int: %v", i.Port, err)
		return 0
	}
	return port
}

func (i *ZdayeItem) GetCountry() string {
	return ""
}

func (i *ZdayeItem) GetRegion() string {
	resp := strings.Split(i.Location, " ")
	if len(resp) > 0 {
		return resp[0]
	}
	return ""
}

func (i *ZdayeItem) GetISP() string {
	resp := strings.Split(i.Location, " ")
	if len(resp) > 1 {
		return resp[1]
	}
	return ""
}

func (i *ZdayeItem) GetUpdatedAt() int64 {
	return 0
}
