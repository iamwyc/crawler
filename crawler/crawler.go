package crawler

import (
	"errors"
	"strings"
)

func HttpGet(uri string) (err error) {
	var (
		flag = false
	)
	if strings.Contains(uri, "collectui") {
		downloadCollectUI(uri)
		flag = true;
	}
	if strings.Contains(uri, "dribbble") {
		downloadDribbble(uri)
		flag = true
	}
	if !flag {
		return errors.New("暂不支持当前uri")
	}
	return nil
}
