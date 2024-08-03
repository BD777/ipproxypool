# ipproxypool
Version
[En](https://github.com/BD777/ipproxypool/blob/main/README.md)
|
[Zh](https://github.com/BD777/ipproxypool/blob/main/README_ZH.md)

Collect and maintain IP proxy data, and provide a RESTful API for access.

## Crawlers Status
| Name | Status |
| --- | --- |
| [89ip](https://www.89ip.cn) | [![Detect 89ip](https://github.com/BD777/ipproxypool/actions/workflows/detect_crawler_89ip.yml/badge.svg)](https://github.com/BD777/ipproxypool/actions/workflows/detect_crawler_89ip.yml) |
| [ip3366](http://www.ip3366.net/free) | [![Detect ip3366](https://github.com/BD777/ipproxypool/actions/workflows/detect_crawler_ip3366.yml/badge.svg)](https://github.com/BD777/ipproxypool/actions/workflows/detect_crawler_ip3366.yml) |
| [kuaidaili](https://www.kuaidaili.com/free) | [![Detect kuaidaili](https://github.com/BD777/ipproxypool/actions/workflows/detect_crawler_kuaidaili.yml/badge.svg)](https://github.com/BD777/ipproxypool/actions/workflows/detect_crawler_kuaidaili.yml) |

## Features
1. Crawl IP proxies from the web and store them in a SQLite database.
2. Periodically check if the proxies in the database are still valid and remove them if they are invalid.
3. Provide RESTful APIs for access.

### API
#### **GET** `/`
Query proxies.

| param | desc | sample |
| -- | -- | -- |
| type | proxy type<br/>1: transparent<br/>2: annoymous<br/>3: high annoymous | 3 |
| protocol | protocol type<br/>1: HTTP<br/>2: HTTPS<br/>3: Both | 3 |
| limit | limit count | 10 |

#### **GET** `/count`
Query total count of proxies.

#### **POST** `/delete`
Delete proxy from database.

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
`go run cmd/main.go`

