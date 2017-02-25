package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"strings"
	"strconv"
)

type InputData struct {
	MaxCacheSize int
	VideoSizes   []int
	Endpoints    []Endpoint
	Requests     []Request
}

type Endpoint struct {
	DataCenterLatency int
	CacheToLatency    map[int]int
}

type Request struct {
	VideoID, Quantity, EndpointID int
}

type VideoScore struct {
	VideoID, Score int
}

func NewInputDataFromFile(inputPath string) InputData {
	return NewInputData(readFile(inputPath))
}

func NewInputData(input string) InputData {
	inputRows := strings.Split(input, "\n")
	inputHeaders := strings.Split(inputRows[0], " ")

	videoQuantity, _ := strconv.Atoi(inputHeaders[0])
	endpointQuantity, _ := strconv.Atoi(inputHeaders[1])
	requestQuantity, _ := strconv.Atoi(inputHeaders[2])
	maxCacheSize, _ := strconv.Atoi(inputHeaders[4])

	videoSizes := parseVideoSizes(inputRows[1], videoQuantity)
	endpoints := parseEndpoints(inputRows[2: len(inputRows)-requestQuantity-1], endpointQuantity)
	requests := parseRequests(inputRows[len(inputRows)-requestQuantity-1:], requestQuantity)

	return InputData{
		MaxCacheSize: maxCacheSize,
		VideoSizes:   videoSizes,
		Endpoints:    endpoints,
		Requests:     requests,
	}
}

func parseVideoSizes(input string, quantity int) []int {
	videoSizes := make([]int, quantity)

	videoSizeStrings := strings.Split(input, " ")
	for i := range videoSizes {
		videoSizes[i], _ = strconv.Atoi(videoSizeStrings[i])
	}

	return videoSizes
}

func parseEndpoints(input []string, quantity int) []Endpoint {
	endpoints := make([]Endpoint, quantity)

	endpointCounter, cacheCounter, cacheQuantity := 0, 0, 0
	for i := range input {
		data := strings.Split(input[i], " ")
		if cacheCounter == 0 {
			dataCenterLatency, _ := strconv.Atoi(data[0])
			cacheQuantity, _ = strconv.Atoi(data[1])
			endpoints[endpointCounter] = Endpoint{
				DataCenterLatency: dataCenterLatency,
				CacheToLatency:    make(map[int]int, cacheQuantity),
			}
		} else {
			cacheID, _ := strconv.Atoi(data[0])
			latency, _ := strconv.Atoi(data[1])
			endpoints[endpointCounter].CacheToLatency[cacheID] = latency
		}

		if cacheCounter < cacheQuantity {
			cacheCounter++
		} else {
			endpointCounter++
			cacheCounter = 0
		}
	}

	return endpoints
}

func parseRequests(input []string, quantity int) []Request {
	requests := make([]Request, quantity)

	for i := range requests {
		data := strings.Split(input[i], " ")
		videoID, _ := strconv.Atoi(data[0])
		endpointID, _ := strconv.Atoi(data[1])
		quantity, _ := strconv.Atoi(data[2])
		requests[i] = Request{
			VideoID:    videoID,
			EndpointID: endpointID,
			Quantity:   quantity,
		}
	}

	return requests
}

func main() {
	flag.Parse()
	input := flag.Args()

	inputPath := input[0]
	//outputPath := input[1]

	inputData := NewInputDataFromFile(inputPath)

	fmt.Printf("%+v", inputData)
}

func readFile(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(content)
}
