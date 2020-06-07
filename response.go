package sgh

import "net/http"

type Response struct {
	Header     http.Header
	Result     interface{}
	ResultType BodyType
}

func NewResponse(resultStruct interface{}, resultType BodyType) *Response {
	return &Response{
		Result:     resultStruct,
		ResultType: resultType,
	}
}

func NewDefaultResponse(resultStruct interface{}) *Response {
	return NewResponse(resultStruct, Default)
}

func NewJsonResponse(resultStruct interface{}) *Response {
	return NewResponse(resultStruct, Json)
}

func NewXmlResponse(resultStruct interface{}) *Response {
	return NewResponse(resultStruct, Xml)
}
