package main

import (
	"time"

	"github.com/tett23/narou_epub/src/config"
	"github.com/tett23/narou_epub/src/job"
	"github.com/tett23/narou_epub/src/server"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	go func() {
		job.ProcessJobQueue()
	}()
	go func() {
		server.Start("localhost", 1323)
	}()

	for _, nCode := range conf.NCodes {
		job.Enqueue(job.JobTypeFetchLatestEpisode, nCode, -1)
	}

	ch := time.Tick(1 * time.Hour)
	// ch := time.Tick(5 * time.Second)
	for {
		select {
		case <-ch:
			for _, nCode := range conf.NCodes {
				job.Enqueue(job.JobTypeFetchLatestEpisode, nCode, -1)
			}
		}
	}
}
