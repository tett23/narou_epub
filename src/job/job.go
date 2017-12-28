package job

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tett23/narou_epub/src/config"
	"github.com/tett23/narou_epub/src/epub"
	"github.com/tett23/narou_epub/src/novel"
)

type Job struct {
	JobType       JobType
	NCode         string
	EpisodeNumber int
}

type JobType int

const (
	JobTypeNone JobType = iota
	JobTypeFetchLatestEpisode
	JobTypeFetchEpisode
	JobTypeFetchAll
	JobTypeBuildLatestEpisode
	JobTypeBuildEpisode
	JobTypeBuildAll
	JobTypeSendToKindleLatest
	JobTypeSendToKindleEpisode
	JobTypeSendToKindleAll
)

var queue chan *Job

func init() {
	queue = make(chan *Job, 1)
}

func Enqueue(t JobType, nCode string, episodeNumber int) {
	queue <- &Job{
		JobType:       t,
		NCode:         nCode,
		EpisodeNumber: episodeNumber,
	}
}

func ProcessJobQueue() {
	for {
		select {
		case item := <-queue:
			fmt.Println("receive job", item)
			go func() {
				if err := processJob(item); err != nil {
					fmt.Println("ProcessJobQueue err", err)
				}
			}()
		}

		time.Sleep(5 * time.Second)
	}
}

func processJob(job *Job) error {
	var err error
	switch job.JobType {
	case JobTypeFetchLatestEpisode:
		err = fetchLatestEpisode(job)
	case JobTypeFetchEpisode:
		err = fetchEpisode(job)
	case JobTypeFetchAll:
		err = fetchAll(job)
	case JobTypeBuildLatestEpisode:
		err = buildLatestEpisode(job)
	case JobTypeBuildEpisode:
		err = buildEpisode(job)
	case JobTypeBuildAll:
		err = buildAll(job)
	case JobTypeSendToKindleLatest:
		err = sendToKindleLatest(job)
	case JobTypeSendToKindleEpisode:
		err = sendToKindleEpisode(job)
	case JobTypeSendToKindleAll:
		err = sendToKindleAll(job)
	case JobTypeNone:
		fallthrough
	default:
		return nil
	}

	return err
}

func fetchLatestEpisode(job *Job) error {
	data, err := getDetail([]string{job.NCode})
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	d := data[0]

	var container *novel.Container
	if container, err = novel.GetContainer(d.NCode); err != nil {
		container = novel.NewContainer(d.NCode, d.Title, d.Writer, d.UserID)
		Enqueue(JobTypeFetchAll, d.NCode, -1)
	}
	container.GeneralAllNo = d.GeneralAllNo

	updatedAt, err := d.LastUpdatedAt()
	if err != nil {
		return err
	}
	if container.UpdatedAt.Unix() >= updatedAt.Unix() {
		return nil
	}

	if err = container.Write(); err != nil {
		return err
	}

	Enqueue(JobTypeFetchEpisode, job.NCode, d.GeneralAllNo)
	Enqueue(JobTypeBuildLatestEpisode, d.NCode, d.GeneralAllNo)
	Enqueue(JobTypeBuildAll, d.NCode, d.GeneralAllNo)

	return nil
}

func fetchEpisode(job *Job) error {
	episode, err := crawl(job.NCode, job.EpisodeNumber)
	if err != nil {
		return err
	}

	if err = episode.Write(); err != nil {
		return err
	}

	return nil
}

func fetchAll(job *Job) error {
	data, err := getDetail([]string{job.NCode})
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	d := data[0]
	for i := 1; i <= data[0].GeneralAllNo; i++ {
		Enqueue(JobTypeFetchEpisode, job.NCode, i)
	}

	container := novel.NewContainer(d.NCode, d.Title, d.Writer, d.UserID)
	container.GeneralAllNo = d.GeneralAllNo
	if err = container.Write(); err != nil {
		return err
	}

	Enqueue(JobTypeBuildAll, d.NCode, -1)

	return nil
}

func buildLatestEpisode(job *Job) error {
	container, err := novel.GetContainer(job.NCode)
	if err != nil {
		errors.Wrapf(err, "buildLatestEpisode GetContainer %s", job.NCode)
	}

	episode, err := container.LatestEpisode()
	if err != nil {
		errors.Wrapf(err, "buildLatestEpisode container.LatestEpisode %s", job.NCode)
	}

	buildEpisode(&Job{
		JobType:       JobTypeBuildEpisode,
		NCode:         job.NCode,
		EpisodeNumber: episode.EpisodeNumber,
	})

	return nil
}

func buildEpisode(job *Job) error {
	container, err := novel.GetContainer(job.NCode)
	if err != nil {
		errors.Wrapf(err, "buildLatestEpisode GetContainer %s", job.NCode)
	}

	episode, err := container.GetEpisode(job.EpisodeNumber)
	if err != nil {
		errors.Wrapf(err, "buildLatestEpisode container.LatestEpisode %s", job.NCode)
	}

	e := epub.NewEpub(job.NCode, container.Title, container.Author)
	if err = e.GenerateByEpisodeNumber(episode.EpisodeNumber); err != nil {
		errors.Wrapf(err, "buildLatestEpisode GenerateByEpisodeNumber %s", job.NCode)
	}

	fmt.Printf("write epub file %s\n", e.OutputFileName())
	return nil
}

func buildAll(job *Job) error {
	container, err := novel.GetContainer(job.NCode)
	if err != nil {
		errors.Wrapf(err, "buildLatestEpisode GetContainer %s", job.NCode)
	}

	e := epub.NewEpub(job.NCode, container.Title, container.Author)
	if err = e.GenerateAll(); err != nil {
		errors.Wrapf(err, "buildLatestEpisode GenerateAll %s", job.NCode)
	}

	fmt.Printf("write epub file %s\n", e.OutputFileName())

	return nil
}

func sendToKindleLatest(job *Job) error {
	return nil
}

func sendToKindleEpisode(job *Job) error {
	return nil
}

func sendToKindleAll(job *Job) error {
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

func crawl(nCode string, episodeNumber int) (*novel.Episode, error) {
	params := url.Values{
		"no":      {strconv.Itoa(episodeNumber)},
		"hankaku": {"0"},
		"code":    {"utf-8"},
		"kaigyo":  {"crlf"},
	}

	nCodeNumber, err := nCodeNumber(nCode)
	if err != nil {
		return nil, errors.Wrapf(err, "job.crawl nCodeNumber NCode: %s, EpisodeNumber: %d", nCode, episodeNumber)
	}

	endpoint := fmt.Sprintf("https://ncode.syosetu.com/txtdownload/dlstart/ncode/%d/", nCodeNumber)
	res, err := http.PostForm(endpoint, params)
	if err != nil {
		return nil, errors.Wrapf(err, "job.crawl http.PostForm NCode: %s, EpisodeNumber: %d", nCode, episodeNumber)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "job.crawl iouti.ReadAll NCode: %s, EpisodeNumber: %d", nCode, episodeNumber)
	}

	episode := novel.NewEpisode(nCode, episodeNumber, string(body))

	return episode, nil
}
