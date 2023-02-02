package main

import (
	"encoding/json"
	"time"
)

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
