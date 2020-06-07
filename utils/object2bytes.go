package utils

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
)

func Struct2Json(obj interface{}) []byte {
	// single
	switch obj.(type) {
	case string, int8, int16, int32, int, int64, float32, float64,
		uint8, uint16, uint32, uint, uint64, bool, uintptr:
		return []byte(fmt.Sprintf("%v", obj))
	}
	bs, _ := json.Marshal(obj)
	return bs
}

func Struct2Xml(obj interface{}) []byte {
	// single
	switch obj.(type) {
	case string, int8, int16, int32, int, int64, float32, float64,
		uint8, uint16, uint32, uint, uint64, bool, uintptr:
		return []byte(fmt.Sprintf("%v", obj))
	}
	bs, _ := xml.Marshal(obj)
	return bs
}

func Struct2UrlQuery(obj interface{}) []byte {
	// single
	switch obj.(type) {
	case string, int8, int16, int32, int, int64, float32, float64,
		uint8, uint16, uint32, uint, uint64, bool, uintptr:
		return []byte(fmt.Sprintf("%v=", obj))
	}

	var m map[string]interface{}
	json.Unmarshal(Struct2Json(obj), &m)
	bb := bytes.Buffer{}
	for k, v := range m {
		bb.WriteString(k)
		bb.WriteString("=")
		var valueStr string
		switch v.(type) {
		case string:
			valueStr = v.(string)
		default:
			vbs, _ := json.Marshal(v)
			valueStr = string(vbs)
		}
		bb.WriteString(valueStr)
		bb.WriteString("&")
	}
	return bb.Bytes()[:bb.Len()-1]
}
