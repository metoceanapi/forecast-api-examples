package main

import (
	"unsafe"
)

type VariableBase64 struct {
	Variable
	Data []byte `json:"data"` // golang treats []byte as base64-encoded by default
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

func (v VariableBase64) UnpackData(maskLut map[uint32]uint8) VariableFloat {
	regular := VariableFloat{Variable: v.Variable}
	sizeOfF32 := unsafe.Sizeof(float32(0))
	n := uintptr(len(v.Data)) / sizeOfF32
	regular.Data = unsafe.Slice((*float32)(unsafe.Pointer(&v.Data[0])), n)
	regular.NoData = make([]uint8, n)
	for index := range regular.NoData {
		offset := uintptr(index) * sizeOfF32
		u32 := *(*uint32)(unsafe.Pointer(&v.Data[offset]))
		if mask := maskLut[u32]; mask != 0 {
			regular.NoData[index] = mask
		}
	}

	return regular
}

func (r PointDataResponseBase64) UnpackData() PointDataResponseFloat {
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
	for name, v := range r.Variables {
		regular.Variables[name] = v.UnpackData(maskLut)
	}
	return regular
}
