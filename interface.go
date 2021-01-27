package sgh

import "time"

// 客户端程序
type HttpClient interface {
	// 处理请求，并设置响应，可传入中间处理函数
	Do(*Request, *Response, ...Handler) error
	// 设置超时时间
	SetTimeout(time.Duration)
	// 设置处理函数，可以根据 HF 参数来确定如何合并`请求处理`中的处理函数
	Handlers([]Handler, HF)
	// 注册子客户端，如果请求能匹配上，则优先使用子客户端处理程序
	RegisterChild(MatchRequest, HttpClient)
}

// 中间处理程序
// !!!注意!!! `req` 无论在请求/响应都会传入，但 `resp` 在请求没有真正完成的时候，都是 `nil`
type Handler func(req *Request, resp *HandelResponse) error

type MatchRequest func(*Request) bool

type HF int

const (
	// 追加到前面
	AppendHead HF = iota
	// 追加到后面
	AppendTail
	// 覆盖当前
	Covered
	// 丢弃追加内容
	Discard
)

var (
	debug bool
)

func OpenDebug() {
	debug = true
}
