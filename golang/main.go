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

type Point struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

// You must specify at least one of from and to, unless you specify an explicit list of times
type TimeSequence struct {
	From     *time.Time `json:"from,omitempty"`
	To       *time.Time `json:"to,omitempty"`
	Interval *Duration  `json:"interval,omitempty"`
	Repeat   *uint32    `json:"repeat,omitempty"`
}

type PointRequest struct {
	Time      TimeSequence `json:"time"`
	Points    []Point      `json:"points"`
	Variables []string     `json:"variables"`
	Format    *string      `json:"outputFormat,omitempty"`
}

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

type Example PointRequest // TODO we will need more than this

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

// TODO need to parse API response

var (
	now      = time.Now()
	examples = map[string]Example{"pointTime": Example(pointTime)}

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

	reqData, err := json.Marshal(example)
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

	var data json.RawMessage
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&data); err != nil {
		fmt.Printf("Request did not return valid JSON, please report this\n")
		return 1
	}
	if res.StatusCode != http.StatusOK {
		fmt.Printf("API returned an error: %v\n%v\n", res.StatusCode, string(data))
		return 1
	}

	fmt.Printf("Received: %v\n", string(data))
	return 0
}
