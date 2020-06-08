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

func (req *Request) Post(url string, body interface{}) *Request {
	return req.HttpMethod(POST).Url(url).HttpBody(body)
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

func (req *Request) SetRequestType(bt BodyType) *Request {
	req.RequestType = bt
	return req
}

func (req *Request) build() (method HttpMethod, url string, header http.Header, body []byte) {
	defer func() { reqFormatPrint(method, url, header, body) }()
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
		if req.Method == GET {
			if strings.Contains(req.URL, "?") {
				url += "&" + string(query)
			} else {
				url += "?" + string(query)
			}
		} else {
			body = query
		}
	}
	return
}

func (req *Request) Do(resp *Response, opts ...func(*Request, *Response)) error {
	return defaultClient.Do(req, resp, opts...)
}

func reqFormatPrint(method HttpMethod, url string, header http.Header, body []byte) {
	if !debug {
		return
	}

	sb := strings.Builder{}

	sb.WriteString("==== http request info ====\n")
	sb.WriteString(method.String())
	sb.WriteString(" ")
	sb.WriteString(url)
	sb.WriteString("\n")
	for k := range header {
		sb.WriteString(k)
		sb.WriteString(": ")
		sb.WriteString(header.Get(k))
		sb.WriteString("\n")
	}
	if len(body) > 0 {
		sb.WriteString("\n")
		sb.Write(body)
		sb.WriteString("\n")
	}
	print(sb.String())
}
