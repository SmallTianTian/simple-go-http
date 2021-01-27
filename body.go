package sgh

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"errors"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/SmallTianTian/simple-go-http/utils"
)

var (
	NotSupport error = errors.New("NOT SUPPORT")
)

func BodyConvert(except *Response, have *HandelResponse) (err error) {
	var bs []byte
	if except.ResultType != have.ResultType {
		switch have.ResultType {
		case Xml:
			bs, err = Xml2(have.Result, except.ResultType)
		case Json:
			bs, err = Json2(have.Result, except.ResultType)
		case Excel:
			bs, err = Excel2(have.Result, except.ResultType)
		case JsonP:
			bs, err = JsonP2(have.Result, except.ResultType)
		case CSV:
			bs, err = Csv2(have.Result, except.ResultType)
		}
	}
	if err != nil {
		return
	}

	switch except.ResultType {
	case Json:
		return json.Unmarshal(bs, except.Result)
	case Xml:
		return xml.Unmarshal(bs, except.Result)
	}
	return NotSupport
}

func Xml2(body []byte, eType BodyType) ([]byte, error) {
	if eType == Xml {
		return body, nil
	}

	var m map[string]interface{}
	if err := utils.Xml2Struct(body, &m); err != nil {
		return nil, err
	}

	if eType == Json {
		return json.Marshal(m)
	}
	return nil, NotSupport
}

func Json2(body []byte, eType BodyType) ([]byte, error) {
	if eType == Json {
		return body, nil
	}

	var m map[string]interface{}
	if err := utils.Json2Struct(body, &m); err != nil {
		return nil, err
	}

	if eType == Xml {
		return xml.Marshal(m)
	}
	return nil, NotSupport
}

func Excel2(body []byte, eType BodyType) ([]byte, error) {
	f, err := excelize.OpenReader(bytes.NewReader(body))

	if err != nil {
		return nil, err
	}
	var name string
	if sheetNames := f.GetSheetList(); len(sheetNames) > 0 {
		name = sheetNames[0]
	}

	rows, err := f.GetRows(name)
	if err != nil {
		return nil, err
	}
	return rows2(rows, eType)
}

func JsonP2(body []byte, eType BodyType) ([]byte, error) {
	start := bytes.IndexByte(body, '(')
	end := bytes.LastIndexByte(body, ')')
	body = body[start+1 : end]
	return Json2(body, eType)
}

func Csv2(body []byte, eType BodyType) ([]byte, error) {
	c := csv.NewReader(bytes.NewReader(body))
	lines, err := c.ReadAll()
	if err != nil {
		return nil, err
	}
	return rows2(lines, eType)
}

func rows2(rows [][]string, eType BodyType) ([]byte, error) {
	title := rows[0]
	array := make([]map[string]string, 0, len(rows)-1)
	for _, row := range rows[1:] {
		m := make(map[string]string)
		for i, k := range title {
			m[k] = row[i]
		}
		array = append(array, m)
	}

	switch eType {
	case Json:
		return json.Marshal(array)
	case Xml:
		return xml.Marshal(array)
	default:
		return nil, NotSupport
	}
}
