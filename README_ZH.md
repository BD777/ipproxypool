# ipproxypool
Version
[En](https://github.com/BD777/ipproxypool/blob/main/README.md)
|
[Zh](https://github.com/BD777/ipproxypool/blob/main/README_ZH.md)


爬取并维护IP代理，提供RESTful API用于读写

## Crawlers Status
| Name | Status |
| --- | --- |
| [89ip](https://www.89ip.cn) | [![Detect Crawler 89ip](https://github.com/BD777/ipproxypool/actions/workflows/detect_crawler_89ip.yml/badge.svg)](https://github.com/BD777/ipproxypool/actions/workflows/detect_crawler_89ip.yml) |
| [ip3366](http://www.ip3366.net/free) | [![Detect Crawler ip3366](https://github.com/BD777/ipproxypool/actions/workflows/detect_crawler_ip3366.yml/badge.svg)](https://github.com/BD777/ipproxypool/actions/workflows/detect_crawler_ip3366.yml) |
| [kuaidaili](https://www.kuaidaili.com/free) | [![Detect kuaidaili](https://github.com/BD777/ipproxypool/actions/workflows/detect_crawler_kuaidaili.yml/badge.svg)](https://github.com/BD777/ipproxypool/actions/workflows/detect_crawler_kuaidaili.yml) |


## Features
1. 爬取IP代理存储到SQLite中；
2. 定期检测数据库中的代理，将不可用的移除；
3. RESTful API；

### API
#### **GET** `/`
查询代理。

| param | desc | sample |
| -- | -- | -- |
| type | proxy type<br/>1: 透明<br/>2: 匿名<br/>3: 高匿 | 3 |
| protocol | protocol type<br/>1: HTTP<br/>2: HTTPS<br/>3: Both | 3 |
| limit | limit count | 10 |

#### **GET** `/count`
查询db中的总数。

#### **POST** `/delete`
删掉一个。

| param | desc | sample |
| -- | -- | -- |
| ip | ip to delete | 10.10.10.10 |
| port | port | 8080 |

## Config
```go
type Config struct {
	Mode string // "debug", "release"
	Host string // default "0.0.0.0"
	Port int    // default 9002

	CountThreshold int64 // if less than threshold, start crawler
	CrawlSleepSec  int64 // sleep seconds between crawls
	DetectSleepSec int64 // sleep seconds between detects proxy from db
}
```

## Run
`go mod tidy`
`go run cmd/main.go`

