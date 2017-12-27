package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/tett23/narou_epub/src/config"
	"github.com/tett23/narou_epub/src/epub"
	"github.com/tett23/narou_epub/src/novel"
)

func main() {
	ch := make(chan *config.CrawlData, 1)
	err := novel.GetFeed(ch)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case item := <-ch:
			fmt.Println("ch", item)
			go func() {
				if err := crawl(item); err != nil {
					fmt.Println("channel err crawl", err)
				}
			}()
		}

		time.Sleep(3 * time.Second)
	}
}

func crawl(item *config.CrawlData) error {
	episodeNumber := item.GeneralAllNo
	params := url.Values{
		"no":      {strconv.Itoa(episodeNumber)},
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

	container := novel.NewContainer(item.NCode)
	if err = container.Write(item, body); err != nil {
		return err
	}

	e := epub.NewEpub(container)
	if err = e.GenerateByEpisodeNumber(episodeNumber); err != nil {
		return err
	}
	fmt.Printf("write epub file %s\n", e.OutputFileName())

	return nil
}
