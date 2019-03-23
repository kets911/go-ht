package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

type Arguments map[string]interface{}

type Info struct {
	maxRequestTime int
	minRequestTime int
	averageTime    float64
	timeouts       int
}

var wg sync.WaitGroup
var mutex sync.Mutex
var timeouts int

func buildInfo(infoChannel chan Info, durationsChannel chan int) {
	var once sync.Once
	var allRequestsTime, maxRequestTime, minRequestTime, count int
	for requestDuration := range durationsChannel {
		once.Do(func() {
			minRequestTime = requestDuration
			maxRequestTime = requestDuration
		})
		allRequestsTime += requestDuration
		count++
		if maxRequestTime < requestDuration {
			maxRequestTime = requestDuration
		}
		if minRequestTime > requestDuration {
			minRequestTime = requestDuration
		}
	}
	infoChannel <- Info{
		minRequestTime: minRequestTime,
		maxRequestTime: maxRequestTime,
		averageTime:    float64(allRequestsTime) / float64(count),
		timeouts:       timeouts,
	}
}

func checkDuration(durationsChannel chan int, url string, client http.Client) {
	defer wg.Done()
	mutex.Lock()
	defer mutex.Unlock()
	start := time.Now()
	_, err := client.Get(url)
	requestDuration := int(time.Since(start).Nanoseconds())
	if e, ok := err.(net.Error); ok && e.Timeout() {
		timeouts++
		return
	}
	if err != nil {
		panic(err)
	}
	durationsChannel <- requestDuration
}

func parseArgs() Arguments {
	var url = flag.String("url", "", "url")
	var requestCount = flag.Int("requestCount", 0, "requestCount")
	var timeout = flag.Int("timeout", 0, "timeout")
	flag.Parse()
	return Arguments{
		"url":          *url,
		"requestCount": *requestCount,
		"timeout":      *timeout,
	}
}

func main() {
	args := parseArgs()
	url := args["url"].(string)
	requestCount := args["requestCount"].(int)
	timeout := time.Duration(args["timeout"].(int))

	durationsChannel := make(chan int)
	infoChannel := make(chan Info)
	client := http.Client{
		Timeout: timeout,
	}
	wg.Add(requestCount)
	go buildInfo(infoChannel, durationsChannel)
	for i := 0; i < requestCount; i++ {
		go checkDuration(durationsChannel, url, client)
	}
	wg.Wait()
	close(durationsChannel)
	info := <-infoChannel
	fmt.Printf("maxRequestTime: %10d\nminRequestTime: %10d\naverageTime: %20f\ntimeouts: %8d\n", info.maxRequestTime, info.minRequestTime, info.averageTime, info.timeouts)
}
