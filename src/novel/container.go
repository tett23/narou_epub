package novel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/pkg/errors"
)

type Container struct {
	NCode     string    `json:"n_code"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	UserID    int       `json:"user_id"`
	UpdatedAt time.Time `json:"updated_at"`

	episodes []Episode
}

var containerRoot = ""

const containerBodyDirectory = "body"
const containerOutDirectory = "out"

func init() {
	_, filename, _, _ := runtime.Caller(1)
	dir, _ := filepath.Abs(filepath.Dir(filename))

	containerRoot = filepath.Join(dir, "containers")
}

func NewContainer(nCode, title, author string, userID int) *Container {
	return &Container{
		NCode:  nCode,
		Title:  title,
		Author: author,
		UserID: userID,
	}
}

func GetContainer(nCode string) (*Container, error) {
	ret := Container{
		NCode: nCode,
	}

	dir := containerDirectory(nCode)
	stat, err := os.Stat(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "GetContainer not found NCode: %s", nCode)
	}
	if !stat.IsDir() {
		return nil, errors.Wrapf(err, "GetContainer not found NCode: %s", nCode)
	}

	if err = ret.loadDirectory(); err != nil {
		return nil, errors.Wrapf(err, "GetContainer loadDirectory %s", nCode)
	}

	return &ret, nil
}

func (container Container) NCodeNumber() (int, error) {
	return nCodeNumber(container.NCode)
}

func (container Container) Write() error {
	if err := checkContainerDirectory(container.NCode); err != nil {
		errors.Wrap(err, "Container.Write checkContainerDirectory")
	}

	fmt.Println("write", container.Path())

	container.UpdatedAt = time.Now()

	bytes, err := json.Marshal(container)
	if err != nil {
		return errors.Wrap(err, "Container.Write json.Marshall")
	}

	outPath := container.Path()
	if err := ioutil.WriteFile(outPath, bytes, os.ModePerm); err != nil {
		return errors.Wrap(err, "Container.Write ioutil.WriteFile")
	}

	return nil
}

func (container Container) Path() string {
	return filepath.Join(containerDirectory(container.NCode), "container.json")
}

func checkContainerDirectory(nCode string) error {
	containerDir := containerDirectory(nCode)
	if !isExistContainerDirectory(nCode) {
		if err := os.MkdirAll(containerDir, os.ModePerm); err != nil {
			return err
		}

		if err := os.MkdirAll(bodyDirectory(nCode), os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func isExistContainerDirectory(nCode string) bool {
	stat, err := os.Stat(containerDirectory(nCode))
	if err != nil {
		return false
	}

	return stat.IsDir()
}

func containerDirectory(nCode string) string {
	return filepath.Join(containerRoot, nCode)
}

func bodyDirectory(nCode string) string {
	return filepath.Join(containerDirectory(nCode), containerBodyDirectory)
}

func (container Container) IsExistEpisode(episodeNumber int) bool {
	return Episode{NCode: container.NCode, EpisodeNumber: episodeNumber}.IsExistFile()
}

func (container *Container) GetAvailableEpisodes() ([]Episode, error) {
	err := container.loadDirectory()
	if err != nil {
		return nil, errors.Wrap(err, "Container.GetAvailableEpisodes loadDirectory")
	}

	return container.episodes, nil
}

func (container *Container) loadDirectory() error {
	dir := bodyDirectory(container.NCode)
	stat, err := os.Stat(dir)
	if err != nil {
		return errors.Wrap(err, "novel.loadDirecory stat error")
	}

	if !stat.IsDir() {
		return errors.Errorf("novel.loadDirecory stat.IsDir %s", dir)
	}

	episodes := make([]Episode, 0)
	err = filepath.Walk(dir, func(path string, stat os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "filepath.Walk")
		}
		if stat.IsDir() {
			return nil
		}

		episode, err := NewEpisodeByJSONFile(path)
		if err != nil {
			return errors.Wrapf(err, "filepath.Walk NewEpisodeByJSONFile %s", path)
		}
		episodes = append(episodes, *episode)

		return nil
	})

	container.episodes = episodes

	return nil

}

func (container Container) GetEpisode(episodeNumber int) (*Episode, error) {
	episode, err := NewEpisodeByJSONFile(episodeFilePath(container.NCode, episodeNumber))
	if err != nil {
		return nil, errors.Errorf("GetEpisode error: episode file not found NCode: %s, EpisodeNumber: %d", container.NCode, episodeNumber)
	}

	return episode, nil
}
