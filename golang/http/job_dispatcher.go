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

type Dispatcher struct {
	jobQueue chan Job
	maxWorkers int
}

func (d Dispatcher) Start() {
	for i := 1; i <= d.maxWorkers; i++ {
		go func(i int) {
			for j := range d.jobQueue {
				go func(job Job) {
					logger.Infof("Started job '%s'...", job.Name)
					time.Sleep(job.Duration)
					logger.Infof("Finished job '%s'!", job.Name)
				}(j)
			}
		}(i)
	}
}

func (d Dispatcher) ScheduleJob(j Job) {
	d.jobQueue <- j
	logger.Infof("Job '%s' scheduled...", j.Name)
}

func NewDispatcher(maxWorkers, maxQueueSize int) *Dispatcher{
	return &Dispatcher{
		maxWorkers: maxWorkers,
		jobQueue: make(chan Job, maxQueueSize),
	}
}

func main() {
	dispatcher := NewDispatcher(*maxWorkers, *maxQueueSize)
	dispatcher.Start()
	http.HandleFunc("/startJob", func (writer http.ResponseWriter, request * http.Request) {
		dispatcher.ScheduleJob(Job {
			Name: request.URL.Path,
			Duration: time.Second,
		})
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d",*serverPort), nil))
}
