package sgh

import "net/http"

type HandelResponse struct {
	Header     http.Header
	Result     []byte
	ResultType BodyType
	Code       int
}

func (hr *HandelResponse) ToResponse(resp *Response) {
	if resp == nil {
		return
	}
	resp.Header = hr.Header
	resp.Code = hr.Code

}

type Response struct {
	Header     http.Header
	Result     interface{}
	ResultType BodyType
	Code       int
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
