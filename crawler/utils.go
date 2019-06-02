package crawler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func downloadResource(url string, dir string, name string, header map[string]string, c chan string) (err error) {
	fmt.Printf("正在下载:%v \n", url)
	defer func() {
		c <- url
	}()
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("user-agent", "User-Agent:Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50")

	if header != nil {
		for k, v := range header {
			request.Header.Set(k, v);
		}
	}

	h := &http.Client{}
	response, err := h.Do(request)
	if err != nil {
		return err
	} else {
		b, _ := ioutil.ReadAll(response.Body)
		ioutil.WriteFile(dir+"/"+name, b, os.ModePerm)
	}
	return nil
}
func download(url string, dir string, name string) (err error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("user-agent", "User-Agent:Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50")

	h := &http.Client{}
	response, err := h.Do(request)

	if err != nil {
		return err
	} else {
		if response.StatusCode != 200 {
			return errors.New("响应吗不是200")
		}
		b, _ := ioutil.ReadAll(response.Body)
		ioutil.WriteFile(dir+"/"+name, b, os.ModePerm)
	}
	return nil
}
func randDir() string {
	s := "./" + time.Now().Format("20060102150405")
	os.Mkdir(s, os.ModePerm)
	return s
}
