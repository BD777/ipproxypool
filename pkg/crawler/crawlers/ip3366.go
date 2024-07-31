package crawlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/BD777/ipproxypool/pkg/utils"
	"github.com/BD777/ipproxypool/pkg/utils/htmlparser"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/levigross/grequests"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type CralwerIP3399 struct {
	session *grequests.Session
}

func NewCralwerIP3399() *CralwerIP3399 {
	return &CralwerIP3399{}
}

func (c *CralwerIP3399) Name() string {
	return "ip3366"
}

func (c *CralwerIP3399) Crawl() <-chan IPProxyItem {
	const MaxPage = 7

	ch := make(chan IPProxyItem, 100)

	go func() {
		defer close(ch)
		for page := 1; page <= MaxPage; page++ {
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

			time.Sleep(time.Second * 3) // take a break
		}
	}()

	return ch
}

func (c *CralwerIP3399) newSession() {
	c.session = grequests.NewSession(&grequests.RequestOptions{
		UserAgent: browser.Chrome(),
	})
}

func (c *CralwerIP3399) crawlPage(page int) ([]IPProxyItem, error) {
	url := fmt.Sprintf("http://www.ip3366.net/free/?stype=1&page=%d", page)

	if c.session == nil {
		c.newSession()
	}

	logrus.Infof("[CralwerIP3399] start to crawl page %d", page)

	httpResp, err := c.session.Get(url, nil)
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get response: %d", httpResp.StatusCode)
	}

	resp := &CralwerIP3399Response{}
	reader := transform.NewReader(bytes.NewReader(httpResp.Bytes()), simplifiedchinese.GB18030.NewDecoder())
	respBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	err = htmlparser.ParseHTML(string(respBytes), resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}

	return resp.Items(), nil
}

type CralwerIP3399Response struct {
	List []*CralwerIP3399Item `xpath:"//table[@class='table table-bordered table-striped']/tbody/tr"`
}

func (c *CralwerIP3399Response) Items() []IPProxyItem {
	items := make([]IPProxyItem, 0, len(c.List))
	for _, item := range c.List {
		item.Clean()
		items = append(items, item)
	}
	return items
}

type CralwerIP3399Item struct {
	IP        string `xpath:"td[1]/text()"`
	Port      int    `xpath:"td[2]/text()"`
	Area      string `xpath:"td[5]/text()"`
	UpdatedAt string `xpath:"td[7]/text()"`
	Region    string
	ISP       string
}

func (c *CralwerIP3399Item) Clean() {
	c.IP = strings.TrimSpace(c.IP)
	c.Area = strings.TrimSpace(c.Area)
	c.UpdatedAt = strings.TrimSpace(c.UpdatedAt)

	parts := strings.Split(c.Area, "_")
	if len(parts) > 1 {
		c.Region, c.ISP = utils.SplitRegionISP(parts[1])
	} else {
		c.Region, c.ISP = utils.SplitRegionISP(c.Area)
	}
}

func (c *CralwerIP3399Item) GetSource() string {
	return "ip3366"
}

func (c *CralwerIP3399Item) GetIP() string {
	return c.IP
}

func (c *CralwerIP3399Item) GetPort() int {
	return c.Port
}

func (c *CralwerIP3399Item) GetCountry() string {
	return "China"
}

func (c *CralwerIP3399Item) GetRegion() string {
	return c.Region
}

func (c *CralwerIP3399Item) GetISP() string {
	return c.ISP
}

func (c *CralwerIP3399Item) GetUpdatedAt() int64 {
	format := "2006/1/2 15:4:5"
	updatedAt := strings.TrimSpace(c.UpdatedAt)
	tz := time.FixedZone("CST", 8*3600)
	t, err := time.ParseInLocation(format, updatedAt, tz)
	if err != nil {
		logrus.Errorf("failed to parse time %s: %v", updatedAt, err)
		return time.Now().Unix()
	}
	return t.Unix()
}
