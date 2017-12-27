package novel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/tett23/narou_epub/src/config"
)

func GetFeed(ch chan *config.CrawlData) error {
	conf, err := config.GetConfig()
	if err != nil {
		return err
	}

	crawlData, err := getDetail(conf.NCodes)
	if err != nil {
		return err
	}

	latestCrawlData, err := config.GetCrawlData()
	if err != nil {
		return err
	}
	fmt.Printf("latestCrawlData %+v\n", latestCrawlData)
	fmt.Printf("crawlData %+v\n", crawlData)

	for _, current := range crawlData {
		ok := false
		var latest config.CrawlData
		for i, d := range latestCrawlData {
			if current.NCode == d.NCode {
				ok = true
				latest = latestCrawlData[i]
				break
			}
		}

		if !ok {
			if current.GeneralAllNo == 1 {
				break
			}

			latest = current
			latest.GeneralAllNo = latest.GeneralAllNo - 1
		}

		if current.GeneralAllNo <= latest.GeneralAllNo {
			break
		}

		fmt.Println(current, latest)
		for i := latest.GeneralAllNo; i < current.GeneralAllNo; i++ {
			// channelでクロールに回す
			item := current
			item.GeneralAllNo = i

			ch <- &item
		}
	}

	return nil
}

const detailAPI = "https://api.syosetu.com/novelapi/api/"

func getDetail(nCodes []string) ([]config.CrawlData, error) {
	parameters := map[string]string{
		// "gzip": "5",
		"out":   "json",
		"ncode": strings.Join(nCodes, "-"),
		"lim":   strconv.Itoa(len(nCodes)),
	}

	params := []string{}
	for k, v := range parameters {
		params = append(params, fmt.Sprintf("%s=%s", k, v))
	}

	endpoint := fmt.Sprintf("%s?%s", detailAPI, strings.Join(params, "&"))
	fmt.Printf("endpoint %s\n", endpoint)
	res, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.Errorf("request error %s", endpoint)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// reader, err := zlib.NewReader(bytes.NewBuffer(body))
	// if err != nil {
	// 	panic(err)
	// 	return nil, err
	// }
	// defer reader.Close()

	// var body []byte
	// if _, err = reader.Read(body); err != nil {
	// 	return nil, err
	// }

	var tmp []interface{}
	if err = json.Unmarshal(body, &tmp); err != nil {
		return nil, err
	}

	var details []config.CrawlData
	for _, item := range tmp[1:] {
		b, err := json.Marshal(item)
		if err != nil {
			return nil, err
		}
		fmt.Printf("detail %#v\n", string(b))

		var detail config.CrawlData
		err = json.Unmarshal(b, &detail)
		if err != nil {
			return nil, err
		}

		num, err := nCodeNumber(detail.NCode)
		if err != nil {
			return nil, err
		}
		detail.NCodeNumber = num

		details = append(details, detail)
	}

	fmt.Printf("details %+v\n", details)

	return details, nil
}

func nCodeNumber(nCode string) (int, error) {
	// (('c'.codepoints[0]-97) * 259974) + (('n'.codepoints[0]-97) * 9999) + 1337
	nCode = strings.ToLower(nCode)
	a := (int(nCode[5]) - 97) * 259974
	b := (int(nCode[6]) - 97) * 9999

	num, err := strconv.Atoi(nCode[1:5])
	if err != nil {
		return -1, err
	}

	return a + b + num, nil
}

// func getFeed(userID int) {
//
// }
