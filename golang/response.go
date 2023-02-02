package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type PointResponse struct {
	response *PointDataResponse
	err      error
}

type PointDataResponse struct {
	Dimensions    map[string]Dimension `json:"dimensions"`
	Variables     map[string]Variable  `json:"variables"`
	NoDataReasons map[string]uint32    `json:"noDataReasons"`
}

type Dimension struct {
	Type  string `json:"type"`
	Units string `json:"units"`
	Data  any    `json:"data"`
}

type Variable struct {
	StandardName string    `json:"standardName"`
	Units        string    `json:"units"`
	SIUnits      string    `json:"siUnits"`
	Dimensions   []string  `json:"dimensions"`
	Data         []float64 `json:"data"`
	NoData       []uint8   `json:"noData"` // or uint32??
}

func (r *PointResponse) UnmarshalJSON(data []byte) error {
	var errors []string
	if err := json.Unmarshal(data, &errors); err == nil {
		r.err = fmt.Errorf(strings.Join(errors, ", "))
		return nil
	}
	dataRes := PointDataResponse{}
	err := json.Unmarshal(data, &dataRes)
	if err == nil {
		r.response = &dataRes
	}
	return err
}
