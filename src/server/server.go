package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/tett23/narou_epub/src/job"
	"github.com/tett23/narou_epub/src/novel"
)

func Start(host string, port int) {
	e := echo.New()

	e.Renderer = &tmpl

	e.GET("/", func(c echo.Context) error {
		containers, err := novel.GetContainers()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}

		data := struct {
			Containers []novel.Container
		}{
			Containers: containers,
		}

		return c.Render(http.StatusOK, "index", data)
	})

	e.GET("/containers/:nCode", func(c echo.Context) error {
		nCode := c.Param("nCode")
		container, err := novel.GetContainer(nCode)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		}

		data := struct {
			Container *novel.Container
		}{
			Container: container,
		}

		return c.Render(http.StatusOK, "container", data)
	})

	e.POST("/containers/:nCode/fetch", func(c echo.Context) error {
		nCode := c.Param("nCode")

		job.Enqueue(job.JobTypeFetchAll, nCode, -1)

		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/containers/%s", nCode))
	})

	e.POST("/containers/:nCode/build", func(c echo.Context) error {
		nCode := c.Param("nCode")

		job.Enqueue(job.JobTypeBuildAll, nCode, -1)

		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/containers/%s", nCode))
	})

	e.POST("/containers/:nCode/publish", func(c echo.Context) error {
		nCode := c.Param("nCode")

		job.Enqueue(job.JobTypeSendToKindleAll, nCode, -1)

		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/containers/%s", nCode))
	})

	e.GET("/containers/:nCode/episode/:episodeNumber", func(c echo.Context) error {
		nCode := c.Param("nCode")
		episodeNumberString := c.Param("episodeNumber")
		episodeNumber, err := strconv.Atoi(episodeNumberString)
		if err != nil {
			return c.String(http.StatusNotFound, fmt.Sprintf("not found POST /containers/%s/episode/%d/fetch", nCode, episodeNumber))
		}

		container, err := novel.GetContainer(nCode)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		}
		episode, err := container.GetEpisode(episodeNumber)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		}

		return c.JSONPretty(http.StatusOK, episode, "  ")
	})

	e.POST("/containers/:nCode/episode/:episodeNumber/fetch", func(c echo.Context) error {
		nCode := c.Param("nCode")
		episodeNumberString := c.Param("episodeNumber")
		episodeNumber, err := strconv.Atoi(episodeNumberString)
		if err != nil {
			return c.String(http.StatusNotFound, fmt.Sprintf("not found POST /containers/%s/episode/%d/fetch", nCode, episodeNumber))
		}

		container, err := novel.GetContainer(nCode)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		}
		_, err = container.GetEpisode(episodeNumber)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		}

		job.Enqueue(job.JobTypeFetchEpisode, nCode, episodeNumber)
		job.Enqueue(job.JobTypeBuildEpisode, nCode, episodeNumber)

		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/containers/%s", nCode))
	})

	e.POST("/containers/:nCode/episode/:episodeNumber/publish", func(c echo.Context) error {
		nCode := c.Param("nCode")
		episodeNumberString := c.Param("episodeNumber")
		episodeNumber, err := strconv.Atoi(episodeNumberString)
		if err != nil {
			return c.String(http.StatusNotFound, fmt.Sprintf("not found POST /containers/%s/episode/%d/publish", nCode, episodeNumber))
		}

		container, err := novel.GetContainer(nCode)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		}
		_, err = container.GetEpisode(episodeNumber)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		}

		job.Enqueue(job.JobTypeSendToKindleEpisode, nCode, episodeNumber)

		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/containers/%s", nCode))
	})

	e.Logger.Info(e.Start(fmt.Sprintf("%s:%d", host, port)))
}
