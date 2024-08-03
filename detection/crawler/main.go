package main

import (
	"flag"
	"os"

	"github.com/BD777/ipproxypool/pkg/crawler/crawlers"
	"github.com/sirupsen/logrus"
)

func main() {
	// parse the command line arguments: -n name string
	name := flag.String("n", "", "crawler name")
	flag.Parse()

	for _, c := range crawlers.Crawlers {
		if c.Name() == *name {
			if !c.Detect() {
				logrus.Infof("crawler %s is not able to fetch proxies", c.Name())
				os.Exit(1)
			} else {
				logrus.Infof("crawler %s is able to fetch proxies", c.Name())
				os.Exit(0)
			}
		}
	}

	logrus.Errorf("crawler %s not found", *name)
	os.Exit(1)
}
