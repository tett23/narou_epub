package novel

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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

func (container Container) NCodeNumber() (int, error) {
	return nCodeNumber(container.NCode)
}

func (container Container) Write(item *config.CrawlData, body []byte) error {
	if err := container.checkDirectory(); err != nil {
		return err
	}

	fmt.Println("write", container.containerDirectory())

	if err := ioutil.WriteFile(container.episodeFilePath(item), body, os.ModePerm); err != nil {
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

func (container Container) episodeFilePath(item *config.CrawlData) string {
	filename := fmt.Sprintf("%04d.txt", item.GeneralAllNo)

	return filepath.Join(containerRoot, container.NCode, containerBodyDirectory, filename)
}

func toHTML(body []byte, container *Container) ([]byte, error) {
	var ret []byte
	tmpl, err := template.New("").Parse(htmlTemplate)
	if err != nil {
		return ret, err
	}

	lines := strings.Split(string(body), "\n")

	data := map[string]interface{}{
		"title": "",
		"lines": lines,
	}

	var buf io.ReadWriter
	if err := tmpl.Execute(buf, data); err != nil {
		return ret, err
	}

	if _, err := buf.Read(ret); err != nil {
		return ret, err
	}

	return ret, nil
}
