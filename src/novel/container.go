package novel

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
	"github.com/tett23/narou_epub/src/config"
)

type Container struct {
	NCode string
}

var containerRoot = ""

const containerBodyDirectory = "body"
const containerOutDirectory = "out"

func init() {
	_, filename, _, _ := runtime.Caller(1)
	dir, _ := filepath.Abs(filepath.Dir(filename))

	containerRoot = filepath.Join(dir, "containers")
}

func NewContainer(nCode string) *Container {
	return &Container{
		NCode: nCode,
	}
}

func GetContainer(nCode string) (*Container, error) {
	ret := Container{
		NCode: nCode,
	}

	dir := ret.containerDirectory()
	stat, err := os.Stat(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "GetContainer not found NCode: %s", nCode)
	}
	if !stat.IsDir() {
		return nil, errors.Wrapf(err, "GetContainer not found NCode: %s", nCode)
	}

	return &ret, nil
}

func (container Container) NCodeNumber() (int, error) {
	return nCodeNumber(container.NCode)
}

func (container Container) Write(item *config.CrawlData, body []byte) error {
	if err := container.checkDirectory(); err != nil {
		return err
	}

	fmt.Println("write", container.containerDirectory())

	if err := ioutil.WriteFile(container.episodeFilePath(item.GeneralAllNo), body, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (container Container) checkDirectory() error {
	containerDir := container.containerDirectory()
	if !container.isExistContainerDirectory() {
		if err := os.MkdirAll(containerDir, os.ModePerm); err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Join(containerDir, containerBodyDirectory), os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func (container Container) isExistContainerDirectory() bool {
	stat, err := os.Stat(container.containerDirectory())
	if err != nil {
		return false
	}

	return stat.IsDir()
}

func (container Container) containerDirectory() string {
	return filepath.Join(containerRoot, container.NCode)
}

func (container Container) episodeFilePath(episodeNumber int) string {
	filename := fmt.Sprintf("%04d.txt", episodeNumber)

	return filepath.Join(containerRoot, container.NCode, containerBodyDirectory, filename)
}

func (container Container) IsExistEpisode(episodeNumber int) bool {
	stat, err := os.Stat(container.episodeFilePath(episodeNumber))
	if err != nil {
		return false
	}

	return !stat.IsDir()
}

func (container Container) GetAvailableEpisodeNumbers() ([]int, error) {
	return nil, nil
}

func (container Container) GetEpisode(episodeNumber int) ([]byte, error) {
	if !container.isExistContainerDirectory() {
		return nil, errors.Errorf("GetEpisode error: episode file not found NCode: %s, EpisodeNumber: %d", container.NCode, episodeNumber)
	}

	p := container.episodeFilePath(episodeNumber)
	ret, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("GetEpisode ioutil.ReadFile(%d)", episodeNumber))
	}

	return ret, nil
}
