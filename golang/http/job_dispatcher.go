/*
* A simple test for scaling easily to ~1000 http reqs/s using goroutines and channels.
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"flag"
	"time"
	"github.com/google/logger"
)

var (
	maxQueueSize = flag.Int("max_queue_size", 100, "The maximum size of the job queue.")
	maxWorkers = flag.Int("max_workers", 20, "The maximum number of workers.")
	serverPort = flag.Int("port", 8787, "The server port.")
)

func init() {
	flag.Parse()
	logger.Init("Logger", true, false, ioutil.Discard)
	logger.Infof("Maximum queue size: %d", *maxQueueSize)
	logger.Infof("Maximum number of workers: %d", *maxWorkers)
	logger.Infof("Starting go-vote webserver on port %d ...", *serverPort)
}

type Job struct {
	Name string
	Duration time.Duration	
}

func doWork(job Job) {
	logger.Infof("Started job '%s'...", job.Name)
	time.Sleep(job.Duration)
	logger.Infof("Finished job '%s'!", job.Name)
}

func main() {
	jobs := make(chan Job, *maxQueueSize)
	for i := 1; i <= *maxWorkers; i++ {
		go func(i int) {
			for j := range jobs {
				go doWork(j)
			}
		}(i)
	}
	http.HandleFunc("/work", func (writer http.ResponseWriter, request * http.Request) {
		jobs <- Job {
			Name: request.URL.Path,
			Duration: time.Second,
		}
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d",*serverPort), nil))
}
