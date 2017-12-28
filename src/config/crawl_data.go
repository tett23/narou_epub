package config

import (
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type CrawlData struct {
	NCode         string `json:"ncode"`
	NCodeNumber   int    `json:"ncode_number"`
	UserID        int    `json:"userid"`
	GeneralAllNo  int    `json:"general_all_no"`
	End           int    `json:"end"`
	IsStop        int    `json:"isstop"`
	Writer        string `json:"writer"`
	Title         string `json:"title"`
	GeneralLastUp string `json:"general_lastup"`
}

func GetCrawlData() ([]CrawlData, error) {
	data, err := ioutil.ReadFile("./config/crawl_data.yml")
	if err != nil {
		return []CrawlData{}, nil
	}

	var ret []CrawlData
	err = yaml.Unmarshal([]byte(data), &ret)
	if err != nil {
		panic(err)
	}

	return ret, err
}

func (crawlData CrawlData) LastUpdatedAt() (time.Time, error) {
	t, err := time.Parse("2006-01-02 15:04:05", crawlData.GeneralLastUp)

	return t, err
}
