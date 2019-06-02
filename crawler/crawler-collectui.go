package crawler

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func downloadCollectUI(uri string) (err error) {
	uris, err := parseCollectUIDoc(uri)
	if (err != nil) {
		return err
	}
	dir := randDir()
	var (
		ch  = make(chan string, len(uris))
		now = 0
	)
	for index := 0; index < len(uris); index++ {
		name := parseCollectUIName(dir, uris[index])
		go downloadResource(uris[index], dir, name, nil, ch);
		time.Sleep(time.Second * 1)
	}

	for {
		finishUri := <-ch
		fmt.Printf("已下载(%d/%d):%v \n", now, len(uris), finishUri)
		now++
		if (now == len(uris)) {
			break
		}
	}
	close(ch)
	fmt.Println("下载完成 (100%):" + uri)
	return
}
func parseCollectUIName(dir, uri string) (name string) {
	split := strings.Split(uri, "/")
	name = split[len(split)-1]
	name = examPath(dir, name, ".png")
	return name + ".png"
}
func examPath(dir, name, suffix string) (string) {
	for {
		_, err := os.Stat(dir + "/" + name + suffix)
		if err == nil {
			name = name + strconv.Itoa(rand.Int())
			continue
		}
		if os.IsNotExist(err) {
			return name
		}
	}
}
func parseCollectUIDoc(uri string) (uris []string, err error) {
	res, err := http.Get(uri)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("请求%v 失败：%d", uri, err))
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("请求%v 状态码为：%d", uri, res.StatusCode))
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	doc.Find(".designs .design").Each(func(i int, s *goquery.Selection) {
		val, exists := s.Find("a").Attr("href")
		if exists {
			uris = append(uris, val)
		}
	})
	return uris, nil;
}
