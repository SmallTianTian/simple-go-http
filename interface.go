package sgh

import "time"

type SimpleHttp interface {
	Do(*Request, *Response) error
	SetTimeout(time.Duration)
}

var (
	debug bool
)

func OpenDebug() {
	debug = true
}
