package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type CrawlData struct {
	NCode        string `json:"ncode"`
	NCodeNumber  int    `json:"ncode_number"`
	UserID       int    `json:"userid"`
	GeneralAllNo int    `json:"general_all_no"`
	End          int    `json:"end"`
	IsStop       int    `json:"isstop"`
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
