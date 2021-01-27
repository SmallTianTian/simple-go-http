package sgh

import (
	"bufio"
	"bytes"
	"net/http"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

var defaultClient = NewSimpleClient()

// const maxDebugBody = 3 * 1024 * 1024
const maxDebugBody = 16 * 1024

func init() {
	var _ HttpClient = NewSimpleClient()
}

// 采用 fasthttp 的简单客户端
type SimpleClient struct {
	// http 客户端
	client *fasthttp.Client
	// 超时时间
	timeout time.Duration
	// 匹配请求的列表
	ms []MatchRequest
	// 子客户端，数量及顺序与 `ms` 一一对应
	cs []HttpClient

	// 中间处理程序
	hs []Handler
	// 与后续中间处理程序的合并方式
	hf HF
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

func (sc *SimpleClient) Handlers(h []Handler, f HF) {
	sc.hs = h
	sc.hf = f
}

func (sc *SimpleClient) RegisterChild(m MatchRequest, t HttpClient) {
	sc.ms = append(sc.ms, m)
	sc.cs = append(sc.cs, t)
}

func (sc *SimpleClient) Do(req *Request, resp *Response, opts ...Handler) error {
	// 如果能匹配上子客户端，使用子客户端处理程序
	for i, m := range sc.ms {
		if m(req) {
			return sc.cs[i].Do(req, resp, opts...)
		}
	}
	handlers := sc.mergeHandler(opts...)
	for _, f := range handlers {
		f(req, nil)
	}

	tmpResp, err := sc.dealReq(req)
	if err != nil {
		return err
	}

	for _, f := range handlers {
		f(req, tmpResp)
	}

	return BodyConvert(resp, tmpResp)
}

// 真正发送请求并构建响应
func (sc *SimpleClient) dealReq(req *Request) (*HandelResponse, error) {
	rq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(rq)
	rp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(rp)
	request2fastRequest(req, rq)

	var timeout time.Duration
	if timeout = req.Timeout; timeout <= 0 {
		timeout = sc.timeout
	}

	if err := sc.client.DoTimeout(rq, rp, timeout); err != nil {
		return nil, err
	}

	return fastResponse2Response(rp)
}

// 合并中间处理程序
func (sc *SimpleClient) mergeHandler(opts ...Handler) []Handler {
	switch sc.hf {
	case AppendHead:
		return append(opts, sc.hs...)
	case Discard:
		if len(sc.hs) == 0 {
			return opts
		}
		return sc.hs
	case Covered:
		if len(opts) == 0 {
			return sc.hs
		}
		return opts
	default:
		fallthrough
	case AppendTail:
		return append(sc.hs, opts...)
	}
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

	// debug info
	debugPrint(func() (h, b []byte) { return []byte(rq.Header.String()), rq.Body() }, false)
}

func fastResponse2Response(rs *fasthttp.Response) (resp *HandelResponse, err error) {
	resp = &HandelResponse{}
	// set body
	resp.Result = rs.Body()

	// debug info
	debugPrint(func() (h, b []byte) { return []byte(rs.Header.String()), resp.Result }, true)

	// set header
	head := http.Header{}
	rs.Header.VisitAll(func(k, v []byte) {
		head.Add(string(k), string(v))
	})
	resp.Header = head

	// set http code
	resp.Code = rs.StatusCode()

	ct := strings.ToLower(head.Get("Content-Type"))
	if cts := strings.Split(ct, ";"); len(cts) > 0 {
		ct = cts[0]
	}

	switch strings.Split(ct, ";")[0] {
	case "application/json":
		resp.ResultType = Json
	case "application/xml":
		resp.ResultType = Xml
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "application/vnd.ms-excel":
		resp.ResultType = Excel
	case "text/csv":
		resp.ResultType = CSV
	default:
		resp.ResultType = Default
	}
	return
}

func debugPrint(f func() (head, body []byte), isResp bool) {
	if !debug {
		return
	}
	h, b := f()

	pre := ">"
	if isResp {
		pre = "<"
	}

	sc := bufio.NewScanner(bytes.NewReader(h))
	for sc.Scan() {
		println(pre, sc.Text())
	}

	if len(b) > maxDebugBody {
		b = b[:maxDebugBody]
		b = append(b, []byte("...")...)
	}
	println(string(b))
}
