package utils

import (
	"encoding/json"
	"encoding/xml"
)

func Json2Struct(bs []byte, obj interface{}) error {
	return json.Unmarshal(bs, obj)
}

func Xml2Struct(bs []byte, obj interface{}) error {
	return xml.Unmarshal(bs, obj)
}
