package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"encoding/json"
	"net/http"
)

const URL string = "https://forecast-v2.metoceanapi.com/point/time"

type Example struct {
	request PointRequest
	process func(PointResponse) any
}

func ptr[V any](v V) *V {
	return &v
}

func mapKeys(m any) []string {
	value := reflect.ValueOf(m)
	if value.Kind() != reflect.Map {
		return nil
	}
	keys := value.MapKeys()
	strings := make([]string, len(keys))
	for index, key := range keys {
		strings[index] = fmt.Sprintf("%v", key)
	}
	sort.Strings(strings)
	return strings
}

var (
	now      = time.Now()
	examples = map[string]Example{"pointTime": {request: pointTime}}

	pointTime = PointRequest{
		Points: []Point{{Longitude: 0, Latitude: 0}},
		Time: TimeSequence{
			From:     &now,
			Interval: ptr(Duration{time.Hour * 3}),
			Repeat:   ptr(uint32(2)),
		},
		Variables: []string{"wave.height"},
	}

	apiKey  = flag.String("apikey", "", "API key to use for requests.")
	example = flag.String("example", "", "Example to run. One of "+strings.Join(mapKeys(examples), ", "))
	help    = flag.Bool("h", false, "Show this message.")
)

func main() {
	os.Exit(mainWithCode())
}

func mainWithCode() int {
	flag.Parse()
	if *help {
		flag.Usage()
		return 0
	}

	example, exampleFound := examples[*example]
	if !exampleFound {
		flag.Usage()
		return 1
	}

	reqData, err := json.Marshal(example.request)
	if err != nil {
		fmt.Printf("Failed to serialise example: %v\n", err)
		return 1
	}

	client := &http.Client{Timeout: time.Second * 12}

	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(reqData))
	req.Header.Set("x-api-key", *apiKey)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Printf("Couldn't create request struct: %v\n", err)
		return 1
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return 1
	}
	defer res.Body.Close()

	var resData PointResponse
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&resData); err != nil {
		fmt.Printf("Request did not return valid JSON, please report this: %v\n", err)
		return 1
	}
	if res.StatusCode != http.StatusOK {
		fmt.Printf("API returned an error: %v\n%v\n", res.StatusCode, resData.err)
		return 1
	}

	fmt.Printf("Received: %v\n", resData.response)
	if example.process == nil {
		return 0
	}

	if processed := example.process(resData); processed != nil {
		fmt.Printf("Processed: %v\n", processed)
	}

	return 0
}
