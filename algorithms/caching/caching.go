package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"strings"
	"strconv"
	"sort"
)

type InputData struct {
	MaxCacheSize  int
	CacheQuantity int
	VideoSizes    []int
	Endpoints     []Endpoint
	Requests      []Request
}

type Endpoint struct {
	DataCenterLatency int
	CacheToLatency    map[int]int
}

type Request struct {
	VideoID, Quantity, EndpointID int
}

type videoScore struct {
	videoID int
	score   float32
}

type byScore []videoScore

func (v byScore) Len() int {
	return len(v)
}

func (v byScore) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v byScore) Less(i, j int) bool {
	return v[i].score < v[j].score
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
	cacheQuantity, _ := strconv.Atoi(inputHeaders[3])
	maxCacheSize, _ := strconv.Atoi(inputHeaders[4])

	videoSizes := parseVideoSizes(inputRows[1], videoQuantity)
	endpoints := parseEndpoints(inputRows[2: len(inputRows)-requestQuantity-1], endpointQuantity)
	requests := parseRequests(inputRows[len(inputRows)-requestQuantity-1:], requestQuantity)

	return InputData{
		MaxCacheSize:  maxCacheSize,
		CacheQuantity: cacheQuantity,
		VideoSizes:    videoSizes,
		Endpoints:     endpoints,
		Requests:      requests,
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

func (inputData *InputData) PrintVideoDistributionPerCache() {
	cacheToVideoScores := inputData.getScoresPerVideoPerCache()
	for _, videoScores := range cacheToVideoScores {
		sort.Sort(sort.Reverse(byScore(videoScores)))
	}

	fmt.Printf("%d\n", inputData.CacheQuantity)

	for cacheID, videoScores := range cacheToVideoScores {
		fmt.Printf("%d ", cacheID)

		currentCacheSize := 0
		videoSet := make(map[int]bool)
		for _, videoScore := range videoScores {
			videoID := videoScore.videoID
			videoSize := inputData.VideoSizes[videoID]
			if _, ok := videoSet[videoID]; !ok && currentCacheSize+videoSize <= inputData.MaxCacheSize {
				fmt.Printf("%d ", videoID)
				videoSet[videoID] = true
				currentCacheSize += videoSize
			}
		}

		fmt.Print("\n")
	}
}

func (inputData *InputData) getScoresPerVideoPerCache() map[int][]videoScore {
	cacheToVideoScores := make(map[int][]videoScore, len(inputData.Endpoints))

	for _, request := range inputData.Requests {
		endpoint := inputData.Endpoints[request.EndpointID]
		videoSize := inputData.VideoSizes[request.VideoID]

		for cacheID, latency := range endpoint.CacheToLatency {
			latencyGain := endpoint.DataCenterLatency - latency
			effectiveLatencyGain := latencyGain * request.Quantity
			score := float32(effectiveLatencyGain) / float32(videoSize)
			if _, ok := cacheToVideoScores[cacheID]; ok {
				cacheToVideoScores[cacheID] = append(cacheToVideoScores[cacheID], videoScore{score: score, videoID: request.VideoID})
			} else {
				// FIXED the first for each endpoint was never used
				cacheToVideoScores[cacheID] = []videoScore{{score: score, videoID: request.VideoID}}
			}
		}
	}

	return cacheToVideoScores
}

func main() {
	flag.Parse()
	input := flag.Args()

	inputPath := input[0]

	inputData := NewInputDataFromFile(inputPath)
	inputData.PrintVideoDistributionPerCache()
}

func readFile(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(content)
}
