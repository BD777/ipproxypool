package utils

import (
	"bytes"
	"io"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

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

func GBK2UTF8(b []byte) string {
	reader := transform.NewReader(bytes.NewReader(b), simplifiedchinese.GB18030.NewDecoder())
	respBytes, _ := io.ReadAll(reader)
	return string(respBytes)
}
