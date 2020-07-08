package sgh

import (
	"net/http"
	"strings"
	"time"

	"github.com/SmallTianTian/simple-go-http/utils"
	"github.com/valyala/fasthttp"
)

var defaultClient = NewSimpleClient()

type SimpleClient struct {
	client  *fasthttp.Client
	timeout time.Duration
}

func NewSimpleClient() *SimpleClient {
	return &SimpleClient{
		client:  &fasthttp.Client{},
		timeout: 30 * time.Second,
	}
}

func (sc *SimpleClient) SetTimeout(timeout time.Duration) {
	sc.timeout = timeout
}

func (sc *SimpleClient) Do(req *Request, resp *Response, opts ...func(*Request, *Response)) error {
	for _, f := range opts {
		f(req, resp)
	}

	rq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(rq)
	rp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(rp)
	request2fastRequest(req, rq)

	timeout := sc.timeout
	if req.Timeout != 0 {
		timeout = req.Timeout
	}

	if err := sc.client.DoTimeout(rq, rp, timeout); err != nil {
		return err
	}
	if err := fastResponse2Response(rp, resp); err != nil {
		return err
	}

	for _, f := range opts {
		f(req, resp)
	}
	return nil
}

func request2fastRequest(req *Request, rq *fasthttp.Request) {
	method, url, header, body := req.build()
	rq.AppendBody(body)
	rq.SetRequestURI(url)
	for k, v := range header {
		if len(v) > 0 {
			rq.Header.Set(k, v[0])
		} else {
			rq.Header.Set(k, "")
		}
	}
	rq.Header.SetMethod(method.String())
}

func fastResponse2Response(rs *fasthttp.Response, resp *Response) (err error) {
	// safe check
	if resp == nil || resp.Result == nil {
		return nil
	}

	head := http.Header{}
	rs.Header.VisitAll(func(k, v []byte) {
		head.Add(string(k), string(v))
	})
	resp.Header = head

	rt := resp.ResultType
	if rt == Default {
		switch strings.ToLower(head.Get("Content-Type")) {
		case "application/json":
			rt = Json
		case "application/xml":
			rt = Xml
		default:
			rt = Json
		}
	}

	switch rt {
	case Json:
		err = utils.Json2Struct(rs.Body(), resp.Result)
	case Xml:
		err = utils.Xml2Struct(rs.Body(), resp.Result)
	}
	return
}
