package epub

import (
	"archive/zip"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tett23/narou_epub/src/novel"
)

type Epub struct {
	NCode string
	ID    string
}

var containerRoot = ""
var outDirectory = ""

const containerBodyDirectory = "body"

func init() {
	_, filename, _, _ := runtime.Caller(1)
	dir, _ := filepath.Abs(filepath.Dir(filename))

	containerRoot = filepath.Join(dir, "epub")
	outDirectory = filepath.Join(dir, "out")
}

func NewEpub(container *novel.Container) *Epub {
	return &Epub{
		NCode: container.NCode,
		ID:    id(container.NCode, time.Now()),
	}
}

func id(nCode string, t time.Time) string {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(1000)

	return fmt.Sprintf("%s-%d-%d", nCode, t.Unix(), r)
}

func (epub Epub) GenerateAll() error {
	container, err := novel.GetContainer(epub.NCode)
	if err != nil {
		return errors.Errorf("epub.GenerateByEpisodeNumber: container not found NCode: %s", epub.NCode)
	}

	episodeNumbers, err := container.GetAvailableEpisodeNumbers()
	if err != nil {
		return errors.Wrapf(err, "epub.GenerateAll: container.GetAvailableEpisodeNumbers NCode: %s", epub.NCode)
	}

	for _, episodeNumber := range episodeNumbers {
		if err = epub.addEpisodeFile(container, episodeNumber); err != nil {
			return errors.Wrapf(err, "epub.GenerateByAll: container.addEpisodeFile NCode: %s, EpisodeNumber: %d", epub.NCode, episodeNumber)
		}
	}

	if err = epub.generateEpub(); err != nil {
		return errors.Wrapf(err, "epub.GenerateAll: generateEpub NCode: %s", epub.NCode)
	}

	return nil
}

func (epub Epub) GenerateByEpisodeNumber(episodeNumber int) error {
	container, err := novel.GetContainer(epub.NCode)
	if err != nil {
		return errors.Errorf("epub.GenerateByEpisodeNumber: container not found NCode: %s", epub.NCode)
	}

	if err = epub.addEpisodeFile(container, episodeNumber); err != nil {
		return errors.Wrapf(err, "epub.GenerateByEpisodeNumber: container.addEpisodeFile NCode: %s, EpisodeNumber: %d", epub.NCode, episodeNumber)
	}

	if err = epub.generateEpub(); err != nil {
		return errors.Wrapf(err, "epub.GenerateByEpisodeNumber: generateEpub NCode: %s, EpisodeNumber: %d", epub.NCode, episodeNumber)
	}

	return nil
}

func (epub Epub) addEpisodeFile(container *novel.Container, episodeNumber int) error {
	if !container.IsExistEpisode(episodeNumber) {
		return errors.Errorf("epub.addEpisodeFile: episode not found NCode: %s, EpisodeNumber: %d", epub.NCode, episodeNumber)
	}

	content, err := container.GetEpisode(episodeNumber)
	if err != nil {
		return errors.Wrapf(err, "epub.addEpisodeFile: container.GetEpisode NCode: %s, EpisodeNumber: %d", epub.NCode, episodeNumber)
	}

	html, err := toHTML(content)
	if err != nil {
		return errors.Errorf("epub.addEpisodeFile: toHTML NCode: %s, EpisodeNumber: %d", epub.NCode, episodeNumber)
	}

	if err := epub.checkDirectory(); err != nil {
		return errors.Wrapf(err, "epub.addEpisodeFile: checkDirectory NCode: %s, EpisodeNumber: %d", epub.NCode, episodeNumber)
	}

	if err := ioutil.WriteFile(epub.episodeFilePath(episodeNumber), html, os.ModePerm); err != nil {
		return errors.Wrapf(err, "epub.addEpisodeFile: WriteFile NCode: %s, EpisodeNumber: %d", epub.NCode, episodeNumber)
	}

	return nil
}

func (epub Epub) generateEpub() error {
	dir := epub.containerDirectory()
	out := epub.OutputFileName()

	if err := createZip(dir, out); err != nil {
		return err
	}

	return nil
}

func (epub Epub) checkDirectory() error {
	containerDir := epub.containerDirectory()
	if !epub.isExistContainerDirectory() {
		if err := os.MkdirAll(containerDir, os.ModePerm); err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Join(containerDir, containerBodyDirectory), os.ModePerm); err != nil {
			return err
		}

		if err := os.MkdirAll(outDirectory, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func (epub Epub) OutputFileName() string {
	return filepath.Join(outDirectory, fmt.Sprintf("%s.epub", epub.ID))
}

func (epub Epub) isExistContainerDirectory() bool {
	stat, err := os.Stat(epub.containerDirectory())
	if err != nil {
		return false
	}

	return stat.IsDir()
}

func (epub Epub) containerDirectory() string {
	return filepath.Join(containerRoot, epub.ID)
}

func (epub Epub) episodeFilePath(episodeNumber int) string {
	filename := fmt.Sprintf("%04d.html", episodeNumber)

	return filepath.Join(containerRoot, epub.ID, containerBodyDirectory, filename)
}

func toHTML(body []byte) ([]byte, error) {
	lines := strings.Split(string(body), "\n")
	title := lines[0]
	lines = lines[2:]

	data := map[string]interface{}{
		"title": title,
		"lines": lines,
	}

	tmpl := template.Must(template.New("base").Parse(htmlTemplate))
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "base", &data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func createZip(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			// header.Name = strings.TrimPrefix(path, source)
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})
	if err != nil {
		return errors.Wrap(err, "createZip filepath.Walk")
	}

	return err
}
