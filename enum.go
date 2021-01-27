package sgh

type HttpMethod uint8

func (hm HttpMethod) String() string {
	switch hm {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case DELETE:
		return "DELETE"
	case PATCH:
		return "PATCH"
	case HEAD:
		return "HEAD"
	case CONNECT:
		return "CONNECT"
	case OPTIONS:
		return "OPTIONS"
	case TRACE:
		return "TRACE"
	}
	return ""
}

const (
	GET HttpMethod = iota
	POST
	PUT
	DELETE
	PATCH
	// 不常使用
	HEAD
	CONNECT
	OPTIONS
	TRACE
)

type BodyType uint8

const (
	Default BodyType = iota
	Json
	Xml
	UrlQuery

	// NOT USE IN REQUEST
	Excel
	JsonP
	CSV
)
