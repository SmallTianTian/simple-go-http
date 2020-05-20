package sgh

import "time"

type SimpleHttp interface {
	Do(*Request, *Response) error
	SetTimeout(time.Duration)
}
