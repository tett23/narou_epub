package main

import (
	"time"

	"github.com/tett23/narou_epub/src/config"
	"github.com/tett23/narou_epub/src/job"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	go func() {
		job.ProcessJobQueue()
	}()

	for _, nCode := range conf.NCodes {
		job.Enqueue(job.JobTypeFetchLatestEpisode, nCode, -1)
	}

	ch := time.Tick(1 * time.Hour)
	for {
		select {
		case <-ch:
			for _, nCode := range conf.NCodes {
				job.Enqueue(job.JobTypeFetchLatestEpisode, nCode, -1)
			}
		}
	}
}
