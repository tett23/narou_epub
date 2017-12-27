package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/tett23/narou_epub/src/config"
	"github.com/tett23/narou_epub/src/novel"
)

func main() {
	ch := make(chan config.CrawlData, 1)
	err := novel.GetFeed(ch)
	if err != nil {
		panic(err)
	}

	for {
		fmt.Println(ch)
		select {
		case item := <-ch:
			fmt.Println("ch", item)
			go crawl(item)
		}

		time.Sleep(3 * time.Second)
	}
}

func crawl(item config.CrawlData) error {
	params := url.Values{
		"no":      {strconv.Itoa(item.GeneralAllNo)},
		"hankaku": {"0"},
		"code":    {"utf-8"},
		"kaigyo":  {"crlf"},
	}

	endpoint := fmt.Sprintf("https://ncode.syosetu.com/txtdownload/dlstart/ncode/%d/", item.NCodeNumber)
	// fmt.Printf("crawl endpoint %s\n", endpoint)
	res, err := http.PostForm(endpoint, params)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// fmt.Printf("crawl body %s\n", body)

	if err = write(item, string(body)); err != nil {
		return err
	}

	return nil
}

func write(item config.CrawlData, body string) error {
	return nil
}
