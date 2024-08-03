package crawler

import (
	"time"

	"github.com/BD777/ipproxypool/pkg/config"
	"github.com/BD777/ipproxypool/pkg/crawler/crawlers"
	"github.com/BD777/ipproxypool/pkg/models"
	"github.com/BD777/ipproxypool/pkg/utils/concurrent"
	"github.com/BD777/ipproxypool/pkg/validator"
	"github.com/sirupsen/logrus"
)

func Run(cfg *config.Config) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("panic: %v", err)
		}
	}()

	go gatherFromCralers(cfg)
	go detectDBProxies(cfg)
}

func stat(prefix string, result <-chan bool) {
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	total, succ, fail := 0, 0, 0

	for {
		select {
		case res, ok := <-result:
			if !ok {
				logrus.Infof("[%s] quit: total:%d succ:%d fail:%d", prefix, total, succ, fail)
				return
			}
			total++
			if res {
				succ++
			} else {
				fail++
			}
		case <-ticker.C:
			logrus.Infof("[%s] total:%d succ:%d fail:%d", prefix, total, succ, fail)
		}
	}
}

func gatherFromCralers(cfg *config.Config) {
	for {
		cnt := models.GetIPProxyCount()
		if cnt < cfg.CountThreshold {
			logrus.Infof("total %d in db, less than threshold %d, start crawling...", cnt, cfg.CountThreshold)

			ch := make(chan models.IPProxy, 100) // Create a new channel for each crawl
			result := make(chan bool, 100)

			go stat("crawl", result)

			go func() {
				concurrent.ExecChan(ch, 50, func(item models.IPProxy) (struct{}, error) {
					validateAndSave(item, result)
					return struct{}{}, nil
				})
				close(result)
			}()

			func() {
				defer close(ch)
				concurrent.Exec(crawlers.Crawlers, 5, func(c crawlers.IPProxyCrawler) (struct{}, error) {
					logrus.Infof("start to crawl %s", c.Name())
					for item := range c.Crawl() {
						ch <- models.IPProxy{
							IP:        item.GetIP(),
							Port:      item.GetPort(),
							Country:   item.GetCountry(),
							Region:    item.GetRegion(),
							ISP:       item.GetISP(),
							Source:    item.GetSource(),
							UpdatedAt: item.GetUpdatedAt(),
						}
					}
					return struct{}{}, nil
				})
			}()
		}

		time.Sleep(time.Second * time.Duration(cfg.CrawlSleepSec))
	}
}

func validateAndSave(item models.IPProxy, result chan<- bool) {
	// logrus.Infof("validateAndSave start %s:%d", item.IP, item.Port)
	proxyProtocol, proxyType, ok := validator.CheckProxy(item.IP, item.Port)
	if !ok {
		// logrus.Infof("proxy %s:%d is invalid", item.IP, item.Port)
		models.DeleteIPProxy(&models.DeleteIPProxyRequest{IP: item.IP, Port: item.Port})
		result <- false
		return
	}

	item.Protocol = proxyProtocol
	item.Type = proxyType
	item.UpdatedAt = time.Now().Unix()

	if err := models.UpsertIPProxy(&item); err != nil {
		logrus.Errorf("failed to upsert proxy %s:%d: %v", item.IP, item.Port, err)
		result <- false
	} else {
		logrus.Debugf("upsert proxy %s:%d successfully", item.IP, item.Port)
		result <- true
	}
}

func detectDBProxies(cfg *config.Config) {
	for {
		result := make(chan bool, 20)

		go stat("detect", result)

		func() {
			proxies := models.ListIPProxy(&models.ListIPProxyRequest{})
			logrus.Infof("start to detect %d proxies", len(proxies))
			concurrent.Exec(proxies, 20, func(item *models.IPProxy) (struct{}, error) {
				validateAndSave(*item, result)
				return struct{}{}, nil
			})
			close(result)
		}()

		time.Sleep(time.Second * time.Duration(cfg.DetectSleepSec))
	}
}
