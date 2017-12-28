package novel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Episode struct {
	NCode         string    `json:"ncode"`
	EpisodeNumber int       `json:"episode_number"`
	EpisodeTitle  string    `json:"episode_title"`
	Body          string    `json:"body"`
	Preface       string    `json:"preface"`
	Postscript    string    `json:"postscript"`
	CrawledAt     time.Time `json:"crawled_at"`
}

const separator = "********************************************"

func NewEpisode(nCode string, episodeNumber int, body string) *Episode {
	ret := &Episode{
		NCode:         nCode,
		EpisodeNumber: episodeNumber,
	}

	ret.Parse(body)

	return ret
}

func NewEpisodeByJSONFile(filename string) (*Episode, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		errors.Wrapf(err, "NewEpisodeByJSONFile json file not found filename: %s", filename)
	}

	var ret Episode
	err = json.Unmarshal(bytes, &ret)

	return &ret, nil
}

func (episode Episode) Path() string {
	return episodeFilePath(episode.NCode, episode.EpisodeNumber)
}

func episodeFilePath(nCode string, episodeNumber int) string {
	filename := fmt.Sprintf("%04d.json", episodeNumber)

	return filepath.Join(containerRoot, nCode, containerBodyDirectory, filename)
}

func (episode Episode) EpubPath() string {
	return filepath.Join("body", fmt.Sprintf("section_%d.html", episode.EpisodeNumber))
}

func (episode *Episode) Parse(txt string) {
	var body string

	parts := strings.Split(txt, separator)
	if len(parts) == 3 {
		episode.Preface = parts[0]
		body = parts[1]
		episode.Postscript = parts[2]
	} else if len(parts) == 2 {
		if len(parts[0]) > len(parts[1]) {
			body = parts[0]
			episode.Postscript = parts[1]
		} else {
			episode.Preface = parts[0]
			body = parts[1]
		}
	} else if len(parts) == 1 {
		body = parts[0]
	}

	body = strings.TrimSpace(body)
	body = strings.Replace(body, "\r\n", "\n", -1)
	// regexpは否定先読みを使えない
	re := regexp.MustCompile("\n\n([^\n])")
	body = re.ReplaceAllString(body, "\n$1")
	re = regexp.MustCompile("\n{4,}\n*")
	body = re.ReplaceAllString(body, "\n\n\n\n")
	body = strings.Replace(body, " ", "　", -1)
	body = strings.Replace(body, "!", "！", -1)
	body = strings.Replace(body, "?", "？", -1)
	// 2文字ずつ取得して置換のほうが正しい結果が得られそうだけど、面倒
	body = strings.Replace(body, "！？", "⁉", -1)
	body = strings.Replace(body, "！？", "⁈", -1)
	body = strings.Replace(body, "！！", "‼", -1)
	body = strings.Replace(body, "？？", "⁇", -1)

	lines := strings.Split(body, "\n")

	episode.EpisodeTitle = lines[0]
	episode.Body = strings.Join(lines[1:], "\n")
}

func (episode *Episode) Write() error {
	if err := checkContainerDirectory(episode.NCode); err != nil {
		return errors.Wrap(err, "Episode.Write checkContainerDirectory")
	}

	fmt.Println("write", episode.Path())

	episode.CrawledAt = time.Now()

	bytes, err := json.Marshal(episode)
	if err != nil {
		return errors.Wrap(err, "Episode.Write json.Marshall")
	}

	outPath := episode.Path()
	if err := ioutil.WriteFile(outPath, bytes, os.ModePerm); err != nil {
		return errors.Wrap(err, "Episode.Write ioutil.WriteFile")
	}

	return nil
}

func (episode Episode) IsExistFile() bool {
	stat, err := os.Stat(episode.Path())
	if err != nil {
		return false
	}

	return !stat.IsDir()
}
