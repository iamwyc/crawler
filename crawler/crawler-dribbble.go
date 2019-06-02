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

func downloadDribbble(uri string) (err error) {
	uris, err := parseDribbbleDoc(uri)
	if err != nil {
		return err
	}
	dir := randDir()
	var (
		ch  = make(chan string, len(uris))
		now = 0
	)
	for index := 0; index < len(uris); index++ {
		name := parseDribbbleName(dir, uris[index])
		go downloadDribbbleResource(uris[index], dir, name, nil, ch);
		time.Sleep(time.Microsecond * 50)
	}

	for {
		finishUri := <-ch
		fmt.Printf("已下载(%d/%d):%v \n", now, len(uris), finishUri)
		now++
		if now == len(uris) {
			break
		}
	}
	close(ch)
	fmt.Println("下载完成 (100%):" + uri)
	return
}
func parseDribbbleName(dir, uri string) (name string) {
	split := strings.Split(uri, "/")
	name = split[len(split)-1]
	return examDribbblePath(dir, name)
}
func examDribbblePath(dir, name string) (string) {
	split := strings.Split(name, ".")
	var (
		prefix string
		suffix string
	)
	if len(split) == 2 {
		prefix = split[0]
		suffix = split[1]
	} else if len(split) == 3 {
		prefix = split[0] + "." + split[1]
		suffix = split[2]
	}
	for {
		_, err := os.Stat(dir + "/" + prefix + "." + suffix)
		if err == nil {
			prefix = prefix + strconv.Itoa(rand.Int())
			continue
		}
		if os.IsNotExist(err) {
			return prefix + "." + suffix
		}
	}
}
func parseDribbbleDoc(uri string) (uris []string, err error) {
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

	doc.Find("#main ol li").Each(func(index int, s *goquery.Selection) {

		resource := s.Find(".dribbble")
		if resource.HasClass("video") {
			val, exists := resource.Attr("data-video-teaser-xlarge")
			if exists {
				uris = append(uris, val)
			} else {
				val, exists = resource.Attr("data-video-teaser-large")
				if exists {
					uris = append(uris, val)
				} else {
					val, exists = resource.Attr("data-video-teaser-small")
					if exists {
						uris = append(uris, val)
					}
				}
			}
			fmt.Printf("video=%v\n", val)
			return
		}
		if resource.HasClass("gif") {
			gif := resource.Find("source")
			val, exists := gif.Attr("srcset")
			if exists {
				fmt.Printf("gif=%v\n", val)
				uris = append(uris, val)
				return
			}
		}

		img := resource.Find("source")
		val, exists := img.Attr("srcset")
		if exists {
			fmt.Printf("image=%v\n", val)
			uris = append(uris, val)
			return
		}
	})

	fmt.Printf("总资源数 %v\n", len(uris))
	return uris, nil;
}

var (
	oneTimes  = "1x"
	twoTimes  = "2x"
	fourtimes = "4x"
)

func downloadDribbbleResource(url string, dir string, name string, header map[string]string, c chan string) (err error) {
	var (
		uri    string
		suffix string
		prefix string
		index  = 2
	)
	defer func() {
		c <- uri
	}()
	prefixInde := strings.LastIndex(url, "_")
	suffixInde := strings.LastIndex(url, ".")
	if prefixInde < 0 || suffixInde < 0 {
		uri = url
	} else {
		prefix = string(url[:prefixInde+1])
		suffix = string(url[suffixInde:])
		if strings.Contains(suffix, "mp4") {
			index = 1;
			uri = url;
		} else {
			uri = prefix + strconv.Itoa(4) + "x" + suffix;
		}
	}
	for ; index > 0; index = index / 2 {
		fmt.Printf("正在下载:%v\n", uri)
		err := download(uri, dir, name)
		if err != nil {
			if strings.Compare(err.Error(), "响应吗不是200") == 0 {
				fmt.Printf("不存在下载地址:%v\n", uri)
				uri = prefix + strconv.Itoa(index) + "x" + suffix;
				continue
			}
		}
		break
	}
	return
}
