package utils

import "strings"

var ISP = []string{
	"电信",
	"联通",
	"移动",
	"铁通",
	"教育网",
	"长城宽带",
	"阿里云",
	"腾讯云",
	"华为云",
	"亚马逊云",
	"电信云",
	"天翼云",
}

func SplitRegionISP(region string) (string, string) {
	for _, isp := range ISP {
		if i := strings.Index(region, isp); i != -1 {
			return region[:i], region[i:]
		}
	}
	return region, ""
}
