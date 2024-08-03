package crawlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/BD777/ipproxypool/pkg/utils/htmlparser"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/levigross/grequests"
	"github.com/robertkrimen/otto"
	"github.com/sirupsen/logrus"
)

type Crawler89IP struct {
	session *grequests.Session
}

func NewCrawler89IP() *Crawler89IP {
	return &Crawler89IP{}
}

func (c *Crawler89IP) Name() string {
	return "89ip"
}

func (c *Crawler89IP) Crawl() <-chan IPProxyItem {
	const MaxPage = 100

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

func (c *Crawler89IP) Detect() bool {
	for page := 1; page <= 2; page++ {
		resp, err := c.crawlPage(page)
		if err != nil {
			logrus.Errorf("failed to detect 89ip: %v", err)
			return false
		}
		if len(resp) == 0 {
			logrus.Errorf("failed to detect 89ip: no items in page %d", page)
			return false
		}
	}
	return true
}

func (c *Crawler89IP) newSession() {
	c.session = grequests.NewSession(&grequests.RequestOptions{
		UserAgent: browser.Chrome(),
	})
}

func (c *Crawler89IP) crawlPage(page int) ([]IPProxyItem, error) {
	url := fmt.Sprintf("https://www.89ip.cn/index_%d.html", page)

	if c.session == nil {
		c.newSession()
	}

	logrus.Infof("[Crawler89IP] start to crawl page %d", page)

	httpResp, err := c.session.Get(url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}
	// if http code == 521, retry
	if httpResp.StatusCode == 521 {
		s := httpResp.String()
		// re.match r'(function\s+.*)</script>' and extract $1
		re := regexp.MustCompile(`(function\s+.*)(</script>)`)
		matches := re.FindStringSubmatch(s)
		if len(matches) < 2 {
			return nil, fmt.Errorf("failed to extract js function")
		}
		js := matches[1]
		js = strings.Replace(js, `eval("qo=eval;qo(po);")`, "return po", -1)

		// re.match "setTimeout(\".*?\((\d+)\)\"")" and extract $1 as function name and $2 as param
		re = regexp.MustCompile(`setTimeout\("(.*?)\((\d+)\)"`)
		matches = re.FindStringSubmatch(s)
		if len(matches) < 3 {
			return nil, fmt.Errorf("failed to extract function name and param")
		}
		funcName := matches[1]
		param := matches[2]

		vm := otto.New()
		if _, err := vm.Run(js); err != nil {
			return nil, fmt.Errorf("failed to run js: %w", err)
		}

		if val, err := vm.Call(funcName, nil, param); err == nil {
			logrus.Infof("result: %s", val)
			// re.match r'https_ydclearance=(.*?);' and extract $1
			re = regexp.MustCompile(`https_ydclearance=(.*?);`)
			matches = re.FindStringSubmatch(val.String())
			if len(matches) < 2 {
				return nil, fmt.Errorf("failed to extract https_ydclearance")
			}
			ydclearance := matches[1]

			// add ydclearance to cookies in session
			c.session.HTTPClient.Jar.SetCookies(httpResp.RawResponse.Request.URL, []*http.Cookie{
				{
					Name:  "https_ydclearance",
					Value: ydclearance,
				},
			})

			logrus.Infof("set cookie: https_ydclearance=%s", ydclearance)
		}

		httpResp, err = c.session.Get(url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get response: %w", err)
		}
	}

	if httpResp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", httpResp.StatusCode)
	}

	resp := &Crawler89IPResponse{}
	err = htmlparser.ParseHTML(httpResp.String(), resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}

	return resp.Items(), nil
}

type Crawler89IPResponse struct {
	List []*Crawler89IPItem `xpath:"//table[@class='layui-table']/tbody/tr"`
}

func (c *Crawler89IPResponse) Items() []IPProxyItem {
	items := make([]IPProxyItem, 0, len(c.List))
	for _, item := range c.List {
		item.Clean()
		items = append(items, item)
	}
	return items
}

type Crawler89IPItem struct {
	IP        string `xpath:"td[1]"`
	Port      int    `xpath:"td[2]"`
	Region    string `xpath:"td[3]"`
	ISP       string `xpath:"td[4]"`
	UpdatedAt string `xpath:"td[5]"`
}

func (c *Crawler89IPItem) Clean() {
	// trimSpace for all fields
	c.IP = strings.TrimSpace(c.IP)
	c.Region = strings.TrimSpace(c.Region)
	c.ISP = strings.TrimSpace(c.ISP)
	c.UpdatedAt = strings.TrimSpace(c.UpdatedAt)
}

func (c *Crawler89IPItem) GetSource() string {
	return "89ip"
}

func (c *Crawler89IPItem) GetIP() string {
	return c.IP
}

func (c *Crawler89IPItem) GetPort() int {
	return c.Port
}

func (c *Crawler89IPItem) GetCountry() string {
	return "China"
}

func (c *Crawler89IPItem) GetRegion() string {
	return c.Region
}

func (c *Crawler89IPItem) GetISP() string {
	return c.ISP
}

func (c *Crawler89IPItem) GetUpdatedAt() int64 {
	format := "2006/01/02 15:04:05"
	updatedAt := strings.Trim(c.UpdatedAt, " ")
	tz := time.FixedZone("CST", 8*3600)
	t, _ := time.ParseInLocation(format, updatedAt, tz)
	return t.Unix()
}
