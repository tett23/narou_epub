package epub

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/tett23/narou_epub/src/novel"
)

type Epub struct {
	NCode string
	UUID  string

	episodes []*novel.Episode
	title    string
	author   string
}

var containerRoot string
var outDirectory string
var templateDirectory string
var tmpl *template.Template

func init() {
	_, filename, _, _ := runtime.Caller(1)
	dir, _ := filepath.Abs(filepath.Dir(filename))

	containerRoot = filepath.Join(dir, "tmp")
	outDirectory = filepath.Join(dir, "epub")
	templateDirectory = filepath.Join(dir, "epub_template")

	tmpl = template.Must(template.New("base").Parse(htmlTemplate + overviewTemplate + tocNcxTemplate + contentOpfTemplate))
}

func NewEpub(nCode, title, author string) *Epub {
	id := uuid.NewV4().String()
	return &Epub{
		NCode: nCode,
		UUID:  id,

		episodes: make([]*novel.Episode, 0),
		title:    title,
		author:   author,
	}
}

func (epub Epub) Name() string {
	return fmt.Sprintf("%s-%s", epub.NCode, epub.UUID)
}

func (epub Epub) GenerateAll() error {
	container, err := novel.GetContainer(epub.NCode)
	if err != nil {
		return errors.Errorf("epub.GenerateByEpisodeNumber: container not found NCode: %s", epub.NCode)
	}

	episodes, err := container.GetAvailableEpisodes()
	if err != nil {
		return errors.Wrapf(err, "epub.GenerateAll: container.GetAvailableEpisodeNumbers NCode: %s", epub.NCode)
	}

	for i := range episodes {
		if err = epub.addEpisodeFile(&episodes[i]); err != nil {
			return errors.Wrapf(err, "epub.GenerateByAll: container.addEpisodeFile NCode: %s, EpisodeNumber: %d", epub.NCode, episodes[i].EpisodeNumber)
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

	episode, err := container.GetEpisode(episodeNumber)
	if err != nil {
		return errors.Errorf("epub.GenerateByEpisodeNumber: episode not found NCode: %s", epub.NCode)
	}

	if err = epub.addEpisodeFile(episode); err != nil {
		return errors.Wrapf(err, "epub.GenerateByEpisodeNumber: container.addEpisodeFile NCode: %s, EpisodeNumber: %d", epub.NCode, episodeNumber)
	}

	if err = epub.generateEpub(); err != nil {
		return errors.Wrapf(err, "epub.GenerateByEpisodeNumber: generateEpub NCode: %s, EpisodeNumber: %d", epub.NCode, episodeNumber)
	}

	return nil
}

func (epub *Epub) addEpisodeFile(episode *novel.Episode) error {
	episodeNumber := episode.EpisodeNumber

	html, err := toHTML(episode)
	if err != nil {
		return errors.Errorf("epub.addEpisodeFile: toHTML NCode: %s, EpisodeNumber: %d", epub.NCode, episodeNumber)
	}

	if err := epub.checkDirectory(); err != nil {
		return errors.Wrapf(err, "epub.addEpisodeFile: checkDirectory NCode: %s, EpisodeNumber: %d", epub.NCode, episodeNumber)
	}

	path := filepath.Join(containerRoot, epub.Name(), episode.EpubPath())
	if err := ioutil.WriteFile(path, html, os.ModePerm); err != nil {
		return errors.Wrapf(err, "epub.addEpisodeFile: WriteFile NCode: %s, EpisodeNumber: %d", epub.NCode, episodeNumber)
	}

	epub.episodes = append(epub.episodes, episode)

	return nil
}

func (epub Epub) generateEpub() error {
	dir := epub.containerDirectory()
	out := epub.OutputFileName()

	if err := copyTemplateDirectory(dir); err != nil {
		return errors.Wrap(err, "copyTemplateDirectory")
	}

	if err := epub.createContentOpf(); err != nil {
		return errors.Wrap(err, "createContentOpf")
	}
	if err := epub.createTocNcx(); err != nil {
		return errors.Wrap(err, "createTocNcx")
	}
	if err := epub.createOverview(); err != nil {
		return errors.Wrap(err, "createOverview")
	}

	if err := createZip(dir, out); err != nil {
		return errors.Wrap(err, "createZip")
	}

	return nil
}

func (epub Epub) checkDirectory() error {
	containerDir := epub.containerDirectory()
	if !epub.isExistContainerDirectory() {
		if err := os.MkdirAll(containerDir, os.ModePerm); err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Join(containerDir, "body"), os.ModePerm); err != nil {
			return err
		}

		if err := os.MkdirAll(outDirectory, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func (epub Epub) OutputFileName() string {
	return filepath.Join(outDirectory, fmt.Sprintf("%s.epub", epub.Name()))
}

func (epub Epub) isExistContainerDirectory() bool {
	stat, err := os.Stat(epub.containerDirectory())
	if err != nil {
		return false
	}

	return stat.IsDir()
}

func (epub Epub) containerDirectory() string {
	return filepath.Join(containerRoot, epub.Name())
}

func toHTML(episode *novel.Episode) ([]byte, error) {
	data := map[string]interface{}{
		"title": episode.EpisodeTitle,
		"body":  strings.Split(episode.Body, "\n"),
	}
	if episode.Preface != "" {
		data["preface"] = strings.Split(episode.Preface, "\n")
	}
	if episode.Postscript != "" {
		data["postscript"] = strings.Split(episode.Postscript, "\n")
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "base", &data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type epubItem struct {
	Name  string
	Path  string
	Order int
}

func (epub Epub) indexItems() []epubItem {
	items := make([]epubItem, 0)
	items = append(items, epubItem{
		Name:  "overview",
		Path:  "body/overview.html",
		Order: 0,
	})

	for i, item := range epub.episodes {
		items = append(items, epubItem{
			Name:  item.EpisodeTitle,
			Path:  item.EpubPath(),
			Order: i + 1,
		})
	}

	return items
}

func (epub Epub) createContentOpf() error {

	data := map[string]interface{}{
		"title":  epub.title,
		"author": epub.author,
		"uuid":   epub.UUID,
		"date":   time.Now().Format("2006-01-02T15:04:05-07:00"),
		// 2017-12-18T23:32:49+00:00
		"items": epub.indexItems(),
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "opf", &data); err != nil {
		return err
	}

	filename := filepath.Join(epub.containerDirectory(), "content.opf")
	ioutil.WriteFile(filename, buf.Bytes(), os.ModePerm)

	return nil
}

func (epub Epub) createTocNcx() error {
	data := map[string]interface{}{
		"title":  epub.title,
		"author": epub.author,
		"uuid":   epub.UUID,
		"date":   time.Now().Format("2006-01-02T15:04:05-07:00"),
		// 2017-12-18T23:32:49+00:00
		"items": epub.indexItems(),
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "ncx", &data); err != nil {
		return err
	}

	filename := filepath.Join(epub.containerDirectory(), "toc.ncx")
	ioutil.WriteFile(filename, buf.Bytes(), os.ModePerm)

	return nil
}

func (epub Epub) createOverview() error {
	data := map[string]interface{}{
		"title":  epub.title,
		"nCode":  epub.NCode,
		"author": epub.author,
		"date":   time.Now().Format("2006-01-02T15:04:05-07:00"),
	}

	if len(epub.episodes) == 1 {
		data["episodeTitle"] = fmt.Sprintf("%d部分:%s", epub.episodes[0].EpisodeNumber, epub.episodes[0].EpisodeTitle)
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "overview", &data); err != nil {
		return err
	}

	filename := filepath.Join(epub.containerDirectory(), "body", "overview.html")
	ioutil.WriteFile(filename, buf.Bytes(), os.ModePerm)

	return nil
}

func createZip(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return errors.Wrap(err, "createZip os.Create")
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return errors.Wrap(err, "createZip os.Stat")
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "createZip filepath.Walk")
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return errors.Wrap(err, "createZip zip.FileInfoHeader")
		}

		if baseDir != "" {
			header.Name = strings.TrimPrefix(path, source)
		}
		header.Name = strings.TrimPrefix(header.Name, "/")

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return errors.Wrap(err, "createZip archive.createheader")
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return errors.Wrap(err, "createZip os.Open")
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		if err != nil {
			return errors.Wrap(err, "createZip io.Copy")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "createZip filepath.Walk")
	}

	return err
}

func copyTemplateDirectory(dest string) error {
	stat, err := os.Stat(templateDirectory)
	if err != nil {
		return errors.Wrap(err, "copyTemplateDirectory stat")
	}

	return dcopy(templateDirectory, dest, stat)
}

func copy(src, dest string, info os.FileInfo) error {
	if info.IsDir() {
		return dcopy(src, dest, info)
	}
	return fcopy(src, dest, info)
}

func fcopy(src, dest string, info os.FileInfo) error {

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = os.Chmod(f.Name(), info.Mode()); err != nil {
		return err
	}

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = io.Copy(f, s)
	return err
}

func dcopy(src, dest string, info os.FileInfo) error {

	if err := os.MkdirAll(dest, info.Mode()); err != nil {
		return err
	}

	infos, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, info := range infos {
		if err := copy(
			filepath.Join(src, info.Name()),
			filepath.Join(dest, info.Name()),
			info,
		); err != nil {
			return err
		}
	}

	return nil
}
