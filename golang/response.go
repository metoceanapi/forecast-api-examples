package main

import (
	"fmt"
	"strings"
)

type Uint8SliceWrapper []uint8 // golang's JSON package serialises []byte/[]uint8 as base64 by default, but we want human-readable output

type VariableFloat struct {
	Variable
	Data   []float32         `json:"data"`   //  N.B. during decoding, nulls are skipped
	NoData Uint8SliceWrapper `json:"noData"` //
}

type PointDataResponseFloat struct {
	Dimensions    map[string]Dimension     `json:"dimensions"`
	Variables     map[string]VariableFloat `json:"variables"`
	NoDataReasons map[string]uint8         `json:"noDataReasons"`
}

type PointResponseFloat struct {
	response *PointDataResponseFloat
	err      error
}

func (r *PointResponseFloat) UnmarshalJSON(data []byte) error {
	return UnmarshalJSONPointResponse(data, r)
}

func (r *PointResponseFloat) SetError(err error) {
	r.err = err
}

func (r *PointResponseFloat) Error() error {
	return r.err
}

func (r *PointResponseFloat) Response() any {
	return r.response
}

func (r *PointResponseFloat) ResponsePtr() any {
	return &r.response
}

func (s Uint8SliceWrapper) MarshalJSON() ([]byte, error) {
	return []byte(strings.Join(strings.Fields(fmt.Sprint(s)), ",")), nil
}
