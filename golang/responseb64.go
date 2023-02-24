package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"
)

type VariableBase64 struct {
	Variable
	Data []byte `json:"data"` // golang's JSON package serialises []byte as base64 by default, and it will also deserialise base64 strings into []byte
}

type PointDataResponseBase64 struct {
	Dimensions  map[string]Dimension      `json:"dimensions"`
	Variables   map[string]VariableBase64 `json:"variables"`
	NoDataCodes map[string]uint32         `json:"noDataMask"`
}

type PointResponseBase64 struct {
	response *PointDataResponseBase64
	err      error
}

func (r *PointResponseBase64) UnmarshalJSON(data []byte) error {
	return UnmarshalJSONPointResponse(data, r)
}

func (r *PointResponseBase64) SetError(err error) {
	r.err = err
}

func (r *PointResponseBase64) Error() error {
	return r.err
}

func (r *PointResponseBase64) Response() any {
	return r.response
}

func (r *PointResponseBase64) ResponsePtr() any {
	return &r.response
}

func UnpackSlice(data []byte, nFloats uintptr, maskLut map[uint32]uint8) ([]float32, error) {
	floats := make([]float32, nFloats)
	return floats, binary.Read(bytes.NewReader(data), binary.LittleEndian, &floats)
}

func UnpackSliceUnsafeLE(data []byte, nFloats uintptr, maskLut map[uint32]uint8) ([]float32, error) {
	// use this if you know you are running on a little-endian platform
	return unsafe.Slice((*float32)(unsafe.Pointer(&data[0])), nFloats), nil
}

func (v VariableBase64) UnpackData(maskLut map[uint32]uint8) (VariableFloat, error) {
	regular := VariableFloat{Variable: v.Variable, Data: make([]float32, 100)}

	sizeOfF32 := unsafe.Sizeof(float32(0))
	nBytes := uintptr(len(v.Data))
	n := nBytes / sizeOfF32
	if n*sizeOfF32 != nBytes { // this should never happen - if it does, please report a bug
		return regular, fmt.Errorf("Cannot unpack data: %v is not a multiple of %v bytes", nBytes, sizeOfF32)
	}

	regular.NoData = make(Uint8SliceWrapper, n)
	for index := range regular.NoData {
		offset := uintptr(index) * sizeOfF32
		u32 := binary.LittleEndian.Uint32(v.Data[offset:])
		if mask := maskLut[u32]; mask != 0 {
			regular.NoData[index] = mask
		}
	}

	var err error
	regular.Data, err = UnpackSlice(v.Data, n, maskLut)
	return regular, err
}

func (r PointDataResponseBase64) UnpackData() (PointDataResponseFloat, error) {
	currCode := uint8(1)
	reasons := make(map[string]uint8, len(r.NoDataCodes)+1)
	maskLut := make(map[uint32]uint8, len(r.NoDataCodes))
	for reason, code := range r.NoDataCodes {
		reasons[reason] = currCode
		maskLut[code] = currCode
		currCode += 1
	}
	reasons["GOOD"] = 0

	regular := PointDataResponseFloat{
		Dimensions:    r.Dimensions,
		NoDataReasons: reasons,
		Variables:     make(map[string]VariableFloat, len(r.Variables)),
	}
	var err error
	for name, v := range r.Variables {
		if regular.Variables[name], err = v.UnpackData(maskLut); err != nil {
			return regular, fmt.Errorf("Failed to unpack %v: %v", name, err)
		}
	}
	return regular, nil
}
