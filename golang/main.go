package main

import (
	"bytes"
	"flag"
	"fmt"
	"math"
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
	request  PointRequest
	response PointResponse
	process  func(PointResponse) (any, error)
}

type PointResponse interface {
	SetError(error)
	Error() error
	Response() any
	ResponsePtr() any
}

type Dimension struct {
	Type  string `json:"type"`
	Units string `json:"units"`
	Data  any    `json:"data"`
}

type Variable struct {
	StandardName string   `json:"standardName"`
	Units        string   `json:"units"`
	SIUnits      string   `json:"siUnits"`
	Dimensions   []string `json:"dimensions"`
}

func UnmarshalJSONPointResponse(data []byte, out PointResponse) error {
	var errors []string
	if err := json.Unmarshal(data, &errors); err == nil {
		out.SetError(fmt.Errorf(strings.Join(errors, ", ")))
		return nil
	}
	responsePtr := out.ResponsePtr()
	return json.Unmarshal(data, &responsePtr)
}

func unpackBase64(r PointResponse) (any, error) {
	b64, ok := r.(*PointResponseBase64)
	if !ok {
		return nil, fmt.Errorf("Couldn't cast argument to a base64 response")
	}
	return b64.response.UnpackData()
}

func directionFromComponents(u, v float32) float32 {
	dividend := math.Atan2(-float64(u), -float64(v)) * (180 / math.Pi)
	return float32(math.Mod(dividend+360, 360))
}

func magnitudeFromComponents(u, v float32) float32 {
	return float32(math.Sqrt(math.Pow(float64(u), 2) + math.Pow(float64(v), 2)))
}

func windFromComponents(r PointResponse) (any, error) {
	wrapper, ok := r.(*PointResponseFloat)
	if !ok {
		return nil, fmt.Errorf("Couldn't cast argument to a float response")
	}
	regular := wrapper.response
	north, found := regular.Variables["wind.speed.northward.at-10m"]
	if !found {
		return nil, fmt.Errorf("Northward wind component not found")
	}
	east, found := regular.Variables["wind.speed.eastward.at-10m"]
	if !found {
		return nil, fmt.Errorf("Eastward wind component not found")
	}
	n := len(east.Data)
	if len(north.Data) != n {
		return nil, fmt.Errorf("Northward and eastward wind components have different lengths")
	}

	direction := make([]float32, n)
	magnitude := make([]float32, n)
	noData := make(Uint8SliceWrapper, n)
	for index := range east.Data {
		// the data are already initialised to 0 in the slice, just like Golang's JSON deserialisation if it encounters a null in an array of numbers
		if east.NoData[index] != 0 {
			noData[index] = east.NoData[index]
			continue
		}
		if north.NoData[index] != 0 {
			noData[index] = north.NoData[index]
			continue
		}
		direction[index] = directionFromComponents(east.Data[index], north.Data[index])
		magnitude[index] = magnitudeFromComponents(east.Data[index], north.Data[index])
	}
	regular.Variables["wind.direction.at-10m"] = VariableFloat{Data: direction, NoData: noData, Variable: Variable{Dimensions: east.Dimensions, Units: "meterPreSecond", SIUnits: "m.s^{-1}", StandardName: "wind_from_direction_at_10m_above_ground_level"}}
	regular.Variables["wind.speed.at-10m"] = VariableFloat{Data: magnitude, NoData: noData, Variable: Variable{Dimensions: east.Dimensions, Units: "degree", SIUnits: "degree", StandardName: "wind_speed_at_10m_above_ground_level"}}
	return regular, nil
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
	examples = map[string]Example{
		"pointTime":            {request: pointTime, response: &PointResponseFloat{}},
		"pointTimeBase64":      {request: pointTimeBase64, response: &PointResponseBase64{}, process: unpackBase64},
		"pointTimeWind":        {request: pointTimeWind, response: &PointResponseFloat{}},
		"pointTimeWindVectors": {request: pointTimeWindVectors, response: &PointResponseFloat{}, process: windFromComponents},
	}

	pointTime = PointRequest{
		Points: []Point{
			{Longitude: 174.7842, Latitude: -37.7935},
			{Longitude: 172, Latitude: -43},
		},
		Time: TimeSequence{
			From:     &now,
			Interval: ptr(Duration{time.Hour * 3}),
			Repeat:   ptr(uint32(2)),
		},
		Variables: []string{"wave.height"},
	}

	pointTimeBase64 = PointRequest{
		Points: []Point{
			{Longitude: 174.7842, Latitude: -37.7935},
			{Longitude: 172, Latitude: -43},
		},
		Time: TimeSequence{
			From:     &now,
			Interval: ptr(Duration{time.Hour * 3}),
			Repeat:   ptr(uint32(2)),
		},
		Variables: []string{"wave.height"},
		Format:    ptr("base64"),
	}

	pointTimeWind = PointRequest{
		Points: []Point{{Longitude: 174.7842, Latitude: -37.7935}},
		Time: TimeSequence{
			From:     &now,
			Interval: ptr(Duration{time.Hour * 3}),
			Repeat:   ptr(uint32(2)),
		},
		Variables: []string{"wind.speed.at-10m", "wind.direction.at-10m"},
	}

	pointTimeWindVectors = PointRequest{
		Points: []Point{{Longitude: 174.7842, Latitude: -37.7935}},
		Time: TimeSequence{
			From:     &now,
			Interval: ptr(Duration{time.Hour * 3}),
			Repeat:   ptr(uint32(2)),
		},
		Variables: []string{"wind.speed.northward.at-10m", "wind.speed.eastward.at-10m"},
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

	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&example.response); err != nil {
		fmt.Printf("Request did not return valid JSON, please report this: %v\n", err)
		return 1
	}
	if res.StatusCode != http.StatusOK {
		fmt.Printf("API returned an error: %v\n%v\n", res.StatusCode, example.response.Error())
		return 1
	}

	result := example.response.Response()
	if data, err := json.MarshalIndent(result, "", " "); err == nil {
		fmt.Printf("Received: %v\n", string(data))
	} else {
		fmt.Printf("Received: %v but couldn't encode it as JSON: %v\n", result, err)
	}
	if example.process == nil {
		return 0
	}

	if processed, err := example.process(example.response); err != nil {
		fmt.Printf("Failed to process: %v\n", err)
	} else if processed != nil {
		if data, err := json.MarshalIndent(processed, "", " "); err == nil {
			fmt.Printf("Processed: %v\n", string(data))
		} else {
			fmt.Printf("Processed: %v but couldn't encode it as JSON: %v\n", processed, err)
		}
	}

	return 0
}
