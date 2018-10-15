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

/*

Loadtest results:

[Sun Oct 14 2018 21:41:57 GMT-0300 (Brasilia Standard Time)] INFO Requests: 0, requests per second: 0, mean latency: 0 ms
[Sun Oct 14 2018 21:42:02 GMT-0300 (Brasilia Standard Time)] INFO Requests: 4568, requests per second: 914, mean latency: 3.5 ms
[Sun Oct 14 2018 21:42:07 GMT-0300 (Brasilia Standard Time)] INFO Requests: 9567, requests per second: 1000, mean latency: 3 ms
[Sun Oct 14 2018 21:42:12 GMT-0300 (Brasilia Standard Time)] INFO Requests: 14567, requests per second: 1000, mean latency: 2.8 ms
[Sun Oct 14 2018 21:42:17 GMT-0300 (Brasilia Standard Time)] INFO
[Sun Oct 14 2018 21:42:17 GMT-0300 (Brasilia Standard Time)] INFO Target URL:          http://localhost:8787/work
[Sun Oct 14 2018 21:42:17 GMT-0300 (Brasilia Standard Time)] INFO Max time (s):        20
[Sun Oct 14 2018 21:42:17 GMT-0300 (Brasilia Standard Time)] INFO Concurrency level:   100
[Sun Oct 14 2018 21:42:17 GMT-0300 (Brasilia Standard Time)] INFO Agent:               none
[Sun Oct 14 2018 21:42:17 GMT-0300 (Brasilia Standard Time)] INFO Requests per second: 1000
[Sun Oct 14 2018 21:42:17 GMT-0300 (Brasilia Standard Time)] INFO
[Sun Oct 14 2018 21:42:17 GMT-0300 (Brasilia Standard Time)] INFO Completed requests:  19565
[Sun Oct 14 2018 21:42:17 GMT-0300 (Brasilia Standard Time)] INFO Total errors:        0
[Sun Oct 14 2018 21:42:17 GMT-0300 (Brasilia Standard Time)] INFO Total time:          20.002837615 s
[Sun Oct 14 2018 21:42:17 GMT-0300 (Brasilia Standard Time)] INFO Requests per second: 978
[Sun Oct 14 2018 21:42:17 GMT-0300 (Brasilia Standard Time)] INFO Mean latency:        3 ms
[Sun Oct 14 2018 21:42:17 GMT-0300 (Brasilia Standard Time)] INFO
[Sun Oct 14 2018 21:42:17 GMT-0300 (Brasilia Standard Time)] INFO Percentage of the requests served within a certain time
[Sun Oct 14 2018 21:42:17 GMT-0300 (Brasilia Standard Time)] INFO   50%      2 ms
[Sun Oct 14 2018 21:42:17 GMT-0300 (Brasilia Standard Time)] INFO   90%      5 ms
[Sun Oct 14 2018 21:42:17 GMT-0300 (Brasilia Standard Time)] INFO   95%      6 ms
[Sun Oct 14 2018 21:42:17 GMT-0300 (Brasilia Standard Time)] INFO   99%      8 ms
[Sun Oct 14 2018 21:42:17 GMT-0300 (Brasilia Standard Time)] INFO  100%      21 ms (longest request)

*/

