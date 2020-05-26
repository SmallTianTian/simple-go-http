package utils

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
)

func Struct2Json(obj interface{}) []byte {
	bs, _ := json.Marshal(obj)
	return bs
}

func Struct2Xml(obj interface{}) []byte {
	bs, _ := xml.Marshal(obj)
	return bs
}

func Struct2UrlQuery(obj interface{}) []byte {
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
