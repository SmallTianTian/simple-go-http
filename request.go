package sgh

import (
	"context"
	"net/http"
	"strings"

	"github.com/SmallTianTian/simple-go-http/utils"
)

type Request struct {
	Method      HttpMethod
	URL         string
	Header      http.Header
	Body        interface{}
	RequestType BodyType
	Ctx         context.Context
}

func NewRequest(opts ...func(*Request)) *Request {
	req := &Request{}
	for _, f := range opts {
		f(req)
	}
	return req
}

func (req *Request) Post(body interface{}) *Request {
	req.Body = body
	return req.HttpMethod(POST)
}

func (req *Request) Get(url string) *Request {
	return req.HttpMethod(GET).Url(url)
}

func (req *Request) Url(url string) *Request {
	req.URL = url
	return req
}

func (req *Request) HttpMethod(meth HttpMethod) *Request {
	req.Method = meth
	return req
}

func (req *Request) HttpBody(body interface{}) *Request {
	req.Body = body
	return req
}

func (req *Request) SetHeader(key, value string) *Request {
	if req.Header == nil {
		req.Header = http.Header{}
	}
	req.Header.Set(key, value)
	return req
}

func (req *Request) Context(ctx context.Context) *Request {
	req.Ctx = ctx
	return req
}

func (req *Request) build() (method HttpMethod, url string, header http.Header, body []byte) {
	method = req.Method
	url = req.URL
	if req.Header == nil {
		header = http.Header{}
	} else {
		header = req.Header.Clone()
	}

	// safe check
	if req.Body == nil {
		return
	}

	// use default option.
	// method == get will use url query body
	// else use json body.
	rt := req.RequestType
	if rt == Default {
		if req.Method == GET {
			rt = UrlQuery
		} else {
			rt = Json
		}
	}

	switch rt {
	case Json:
		body = utils.Struct2Json(req.Body)
		header.Set("Content-Type", "application/json")
	case Xml:
		body = utils.Struct2Xml(req.Body)
		header.Set("Content-Type", "application/xml")
	case UrlQuery:
		query := utils.Struct2UrlQuery(req.Body)
		if strings.Contains(req.URL, "?") {
			url += "&" + string(query)
		} else {
			url += "?" + string(query)
		}
	}
	return
}

func (req *Request) Do(resp *Response, opts ...func(*Request, *Response)) error {
	return defaultClient.Do(req, resp, opts...)
}
