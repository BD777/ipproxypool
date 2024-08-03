package config

type Config struct {
	Mode string // "debug", "release"
	Host string // default "0.0.0.0"
	Port int    // default 9002

	CountThreshold int64 // if less than threshold, start crawler
	CrawlSleepSec  int64 // sleep seconds between crawls
	DetectSleepSec int64 // sleep seconds between detects proxy from db
}

var DefaultConfig = &Config{
	Mode:           "release",
	Host:           "0.0.0.0",
	Port:           9002,
	CountThreshold: 100,
	CrawlSleepSec:  30,
	DetectSleepSec: 300,
}
