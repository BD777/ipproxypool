package validator

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/BD777/ipproxypool/pkg/models"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/levigross/grequests"
	"github.com/sirupsen/logrus"
)

const (
	httpURL  = "http://httpbin.org/get"
	httpsURL = "https://httpbin.org/get"
)

func CheckProxy(ip string, port int) (models.ProtocolType, models.ProxyType, bool) {
	for retry := 0; retry < 3; retry++ {
		httpType, httpOk := checkHTTPProxy(httpURL, ip, port)
		httpsType, httpsOk := checkHTTPProxy(httpsURL, ip, port)

		if httpOk && httpsOk {
			return models.Both, httpType, true
		} else if httpOk {
			return models.HTTP, httpType, true
		} else if httpsOk {
			return models.HTTPS, httpsType, true
		}

		time.Sleep(time.Millisecond * 100)
	}

	return 0, 0, false
}

func checkHTTPProxy(testURL, ip string, port int) (models.ProxyType, bool) {
	httpProxy, err := url.Parse("http://" + ip + ":" + strconv.Itoa(port))
	if err != nil {
		logrus.Errorf("failed to parse http proxy: %v", err)
		return 0, false
	}
	httpsProxy, err := url.Parse("https://" + ip + ":" + strconv.Itoa(port))
	if err != nil {
		logrus.Errorf("failed to parse https proxy: %v", err)
		return 0, false
	}

	resp, err := grequests.Get(testURL, &grequests.RequestOptions{
		UserAgent: browser.Chrome(),
		Proxies: map[string]*url.URL{
			"http":  httpProxy,
			"https": httpsProxy,
		},
		RequestTimeout: time.Second * 5,
	})
	if err != nil {
		// logrus.Errorf("failed to check proxy %s:%d %v", ip, port, err)
		return 0, false
	}

	type response struct {
		Origin  string `json:"origin"`
		Headers struct {
			ProxyConnection string `json:"Proxy-Connection"`
		}
	}

	var respData response
	err = resp.JSON(&respData)
	if err != nil {
		// logrus.Errorf("failed to parse response: %v", err)
		return 0, false
	}

	if strings.Contains(respData.Origin, ",") {
		return models.TransparentProxy, true
	} else if respData.Headers.ProxyConnection != "" {
		return models.AnonymousProxy, true
	} else {
		return models.HighAnonymityProxy, true
	}
}
