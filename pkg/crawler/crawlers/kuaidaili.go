package crawlers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/levigross/grequests"
	"github.com/sirupsen/logrus"
)

type CrawlerKuaiDaiLi struct {
	session *grequests.Session
}

func NewKuaiDaiLi() *CrawlerKuaiDaiLi {
	return &CrawlerKuaiDaiLi{}
}

func (c *CrawlerKuaiDaiLi) Name() string {
	return "kuaidaili"
}

func (c *CrawlerKuaiDaiLi) Crawl() <-chan IPProxyItem {
	const MaxPage = 7

	ch := make(chan IPProxyItem, 100)

	go func() {
		defer close(ch)
		subnames := []string{"inha", "fps"}
		for nameIdx, subname := range subnames {
			for page := 1; page <= MaxPage; page++ {
				if nameIdx != 0 || page != 1 {
					time.Sleep(time.Second * 3) // avoid anti-crawler
				}
				items, err := c.crawlPage(subname, page)
				if err != nil {
					logrus.Errorf("failed to crawl %s page %d: %v", subname, page, err)
					return
				}
				if len(items) == 0 {
					break
				}

				for _, item := range items {
					ch <- item
				}
			}
		}
	}()

	return ch
}

func (c *CrawlerKuaiDaiLi) Detect() bool {
	for page := 1; page <= 2; page++ {
		if page > 1 {
			time.Sleep(time.Second * 3) // avoid anti-crawler
		}
		resp, err := c.crawlPage("inha", page)
		if err != nil {
			logrus.Errorf("failed to detect kuaidaili: %v", err)
			return false
		}
		if len(resp) == 0 {
			logrus.Errorf("failed to detect kuaidaili: no items in page %d", page)
			return false
		}
		logrus.Infof("detected %d items in page %d", len(resp), page)
	}
	return true
}

func (c *CrawlerKuaiDaiLi) crawlPage(subname string, page int) ([]*KuaiDaiLiItem, error) {
	if c.session == nil {
		c.newSession()
	}

	logrus.Infof("[CrawlerKuaiDaiLi] start to crawl %s page %d", subname, page)

	resp, err := c.session.Get(fmt.Sprintf("https://www.kuaidaili.com/free/%s/%d/", subname, page), nil)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("failed to get page %d: %s", page, resp.String())
	}

	return parseKuaiDaiLi(resp.String()), nil
}

func (c *CrawlerKuaiDaiLi) newSession() {
	c.session = grequests.NewSession(&grequests.RequestOptions{
		UserAgent: browser.Chrome(),
	})
}

func parseKuaiDaiLi(html string) []*KuaiDaiLiItem {
	// reg match r"fpsList\s*=\s*(\[\{.*?\}\]);"
	reg := regexp.MustCompile(`fpsList\s*=\s*(\[\{.*?\}\]);`)
	matches := reg.FindStringSubmatch(html)
	if len(matches) != 2 {
		logrus.Errorf("failed to match fpsList")
		return nil
	}

	var items []*KuaiDaiLiItem

	if err := json.Unmarshal([]byte(matches[1]), &items); err != nil {
		logrus.Errorf("failed to unmarshal items: %v", err)
		return nil
	}

	return items
}

type KuaiDaiLiItem struct {
	IP            string `json:"ip"`
	LastCheckTime string `json:"last_check_time"`
	Port          string `json:"port"`
	Speed         int    `json:"speed"`
	Location      string `json:"location"`
}

func (i *KuaiDaiLiItem) GetSource() string {
	return "kuaidaili"
}

func (i *KuaiDaiLiItem) GetIP() string {
	return i.IP
}

func (i *KuaiDaiLiItem) GetPort() int {
	port, _ := strconv.Atoi(i.Port)
	return port
}

func (i *KuaiDaiLiItem) GetCountry() string {
	return ""
}

func (i *KuaiDaiLiItem) GetRegion() string {
	return i.Location
}

func (i *KuaiDaiLiItem) GetISP() string {
	return ""
}

func (i *KuaiDaiLiItem) GetUpdatedAt() int64 {
	// 2024-08-03 16:30:01
	fmt := "2006-01-02 15:04:05"
	tz := time.FixedZone("CST", 8*3600)
	t, err := time.ParseInLocation(fmt, i.LastCheckTime, tz)
	if err != nil {
		logrus.Errorf("failed to parse time %s: %v", i.LastCheckTime, err)
		return 0
	}

	return t.Unix()
}
