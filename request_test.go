package sgh

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestNewRequest(t *testing.T) {
	var testHeader http.Header = http.Header{}
	testHeader.Add("test-key", "test-value")
	var testContext context.Context = context.Background()
	testContext = context.WithValue(testContext, "test-ctx-key", "test-ctx-value")

	type args struct {
		opts []func(*Request)
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{
			name: "Change method.",
			args: args{
				opts: []func(*Request){
					func(req *Request) {
						req.Method = POST
					},
				},
			},
			want: &Request{Method: POST},
		},
		{
			name: "Change method many times.",
			args: args{
				opts: []func(*Request){
					func(req *Request) {
						req.Method = POST
					},
					func(req *Request) {
						req.Method = DELETE
					},
				},
			},
			want: &Request{Method: DELETE},
		},
		{
			name: "Change method AND Url.",
			args: args{
				opts: []func(*Request){
					func(req *Request) {
						req.Method = POST
						req.URL = "www.example.com"
					},
				},
			},
			want: &Request{Method: POST, URL: "www.example.com"},
		},
		{
			name: "Change method AND Url in two func.",
			args: args{
				opts: []func(*Request){
					func(req *Request) {
						req.Method = POST
					},
					func(req *Request) {
						req.URL = "www.example.com"
					},
				},
			},
			want: &Request{Method: POST, URL: "www.example.com"},
		},
		{
			name: "Change header",
			args: args{
				opts: []func(*Request){
					func(req *Request) {
						req.Header = testHeader
					},
				},
			},
			want: &Request{Header: testHeader},
		},
		{
			name: "Change body",
			args: args{
				opts: []func(*Request){
					func(req *Request) {
						req.Body = "this is body."
					},
				},
			},
			want: &Request{Body: "this is body."},
		},
		{
			name: "Change request type",
			args: args{
				opts: []func(*Request){
					func(req *Request) {
						req.RequestType = Xml
					},
				},
			},
			want: &Request{RequestType: Xml},
		},
		{
			name: "Change context",
			args: args{
				opts: []func(*Request){
					func(req *Request) {
						req.Ctx = testContext
					},
				},
			},
			want: &Request{Ctx: testContext},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRequest(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_Post(t *testing.T) {
	const (
		postString = "this is post body."
	)
	type fields struct {
		Method      HttpMethod
		URL         string
		Header      http.Header
		Body        interface{}
		RequestType BodyType
		Ctx         context.Context
	}
	type args struct {
		url  string
		body interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Request
	}{
		{
			name:   "Post string",
			fields: fields{},
			args:   args{body: postString},
			want: &Request{
				Method: POST,
				Body:   postString,
			},
		},
		{
			name:   "Post url string",
			fields: fields{},
			args:   args{url: "https://example.com", body: postString},
			want: &Request{
				URL:    "https://example.com",
				Method: POST,
				Body:   postString,
			},
		},
		{
			name:   "Post will reset http method",
			fields: fields{Method: DELETE},
			args:   args{body: postString},
			want: &Request{
				Method: POST,
				Body:   postString,
			},
		},
		{
			name:   "Post will reset http body",
			fields: fields{Body: "old body"},
			args:   args{body: postString},
			want: &Request{
				Method: POST,
				Body:   postString,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{
				Method:      tt.fields.Method,
				URL:         tt.fields.URL,
				Header:      tt.fields.Header,
				Body:        tt.fields.Body,
				RequestType: tt.fields.RequestType,
				Ctx:         tt.fields.Ctx,
			}
			if got := req.Post(tt.args.url, tt.args.body); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request.Post() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_Get(t *testing.T) {
	const (
		testUrl = "www.example.com"
	)

	type fields struct {
		Method      HttpMethod
		URL         string
		Header      http.Header
		Body        interface{}
		RequestType BodyType
		Ctx         context.Context
	}
	type args struct {
		url string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Request
	}{
		{
			name:   "Get url",
			fields: fields{},
			args:   args{url: testUrl},
			want: &Request{
				Method: GET,
				URL:    testUrl,
			},
		},
		{
			name:   "Get will reset http method",
			fields: fields{Method: POST},
			args:   args{url: testUrl},
			want: &Request{
				Method: GET,
				URL:    testUrl,
			},
		},
		{
			name:   "Get will reset url",
			fields: fields{URL: "www.bad-url.com"},
			args:   args{url: testUrl},
			want: &Request{
				Method: GET,
				URL:    testUrl,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{
				Method:      tt.fields.Method,
				URL:         tt.fields.URL,
				Header:      tt.fields.Header,
				Body:        tt.fields.Body,
				RequestType: tt.fields.RequestType,
				Ctx:         tt.fields.Ctx,
			}
			if got := req.Get(tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_Url(t *testing.T) {
	const (
		testUrl = "www.example.com"
	)

	type fields struct {
		Method      HttpMethod
		URL         string
		Header      http.Header
		Body        interface{}
		RequestType BodyType
		Ctx         context.Context
	}
	type args struct {
		url string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Request
	}{
		{
			name:   "Set request url",
			fields: fields{},
			args:   args{url: testUrl},
			want:   &Request{URL: testUrl},
		},
		{
			name:   "Set request url will reset url",
			fields: fields{URL: "www.bad-url.com"},
			args:   args{url: testUrl},
			want:   &Request{URL: testUrl},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{
				Method:      tt.fields.Method,
				URL:         tt.fields.URL,
				Header:      tt.fields.Header,
				Body:        tt.fields.Body,
				RequestType: tt.fields.RequestType,
				Ctx:         tt.fields.Ctx,
			}
			if got := req.Url(tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request.Url() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_HttpMethod(t *testing.T) {
	type fields struct {
		Method      HttpMethod
		URL         string
		Header      http.Header
		Body        interface{}
		RequestType BodyType
		Ctx         context.Context
	}
	type args struct {
		meth HttpMethod
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Request
	}{
		{
			name:   "Set http method",
			fields: fields{},
			args:   args{meth: PATCH},
			want:   &Request{Method: PATCH},
		},
		{
			name:   "Set http method will reset method",
			fields: fields{Method: POST},
			args:   args{meth: PATCH},
			want:   &Request{Method: PATCH},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{
				Method:      tt.fields.Method,
				URL:         tt.fields.URL,
				Header:      tt.fields.Header,
				Body:        tt.fields.Body,
				RequestType: tt.fields.RequestType,
				Ctx:         tt.fields.Ctx,
			}
			if got := req.HttpMethod(tt.args.meth); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request.HttpMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_HttpBody(t *testing.T) {
	const (
		bodyString = "this is string body."
	)
	type fields struct {
		Method      HttpMethod
		URL         string
		Header      http.Header
		Body        interface{}
		RequestType BodyType
		Ctx         context.Context
	}
	type args struct {
		body interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Request
	}{
		{
			name:   "Set http body",
			fields: fields{},
			args:   args{body: bodyString},
			want:   &Request{Body: bodyString},
		},
		{
			name:   "Set http body will reset body",
			fields: fields{Body: "old body."},
			args:   args{body: bodyString},
			want:   &Request{Body: bodyString},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{
				Method:      tt.fields.Method,
				URL:         tt.fields.URL,
				Header:      tt.fields.Header,
				Body:        tt.fields.Body,
				RequestType: tt.fields.RequestType,
				Ctx:         tt.fields.Ctx,
			}
			if got := req.HttpBody(tt.args.body); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request.HttpBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_SetHeader(t *testing.T) {
	var tmp, testHeader http.Header = http.Header{}, http.Header{}
	testHeader.Set("key", "value")
	tmp.Set("key", "old-value")

	type fields struct {
		Method      HttpMethod
		URL         string
		Header      http.Header
		Body        interface{}
		RequestType BodyType
		Ctx         context.Context
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Request
	}{
		{
			name:   "set header.",
			fields: fields{},
			args:   args{key: "key", value: "value"},
			want:   &Request{Header: testHeader},
		},
		{
			name:   "set header will reset header.",
			fields: fields{Header: tmp},
			args:   args{key: "key", value: "value"},
			want:   &Request{Header: testHeader},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{
				Method:      tt.fields.Method,
				URL:         tt.fields.URL,
				Header:      tt.fields.Header,
				Body:        tt.fields.Body,
				RequestType: tt.fields.RequestType,
				Ctx:         tt.fields.Ctx,
			}
			if got := req.SetHeader(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request.SetHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_Context(t *testing.T) {
	var tmp, testContext context.Context = context.Background(), context.Background()
	tmp = context.WithValue(tmp, "key", "old-value")
	testContext = context.WithValue(testContext, "key", "value")

	type fields struct {
		Method      HttpMethod
		URL         string
		Header      http.Header
		Body        interface{}
		RequestType BodyType
		Ctx         context.Context
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Request
	}{
		{
			name:   "set context.",
			fields: fields{},
			args:   args{ctx: testContext},
			want:   &Request{Ctx: testContext},
		},
		{
			name:   "set context will reset context.",
			fields: fields{Ctx: tmp},
			args:   args{ctx: testContext},
			want:   &Request{Ctx: testContext},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{
				Method:      tt.fields.Method,
				URL:         tt.fields.URL,
				Header:      tt.fields.Header,
				Body:        tt.fields.Body,
				RequestType: tt.fields.RequestType,
				Ctx:         tt.fields.Ctx,
			}
			if got := req.Context(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request.Context() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_build(t *testing.T) {
	oldHeader := make(http.Header)
	oldHeader.Set("Content-Type", "application/json")

	type XMlContent struct {
		Key string `xml:"key"`
	}

	type fields struct {
		Method      HttpMethod
		URL         string
		Header      http.Header
		Body        interface{}
		RequestType BodyType
		Ctx         context.Context
	}
	tests := []struct {
		name       string
		fields     fields
		wantMethod HttpMethod
		wantUrl    string
		wantHeader http.Header
		wantBody   []byte
	}{
		{
			name:       "default request in get will be url query.",
			fields:     fields{Method: GET, Body: map[string]string{"key": "value"}},
			wantMethod: GET,
			wantUrl:    "?key=value",
			wantHeader: http.Header{},
			wantBody:   nil,
		},
		{
			name:       "default request in post will be json body.",
			fields:     fields{Method: POST, Body: map[string]string{"key": "value"}},
			wantMethod: POST,
			wantUrl:    "",
			wantHeader: http.Header{"Content-Type": []string{"application/json"}},
			wantBody:   []byte(`{"key":"value"}`),
		},
		{
			name:       "default request in get will be url query AND append to url.",
			fields:     fields{Method: GET, Body: map[string]string{"key": "value"}, URL: "https://example.com?test=hello"},
			wantMethod: GET,
			wantUrl:    "https://example.com?test=hello&key=value",
			wantHeader: http.Header{},
			wantBody:   nil,
		},
		{
			name:       "get use json body.",
			fields:     fields{RequestType: Json, Method: GET, Body: map[string]string{"key": "value"}, URL: "https://example.com?test=hello"},
			wantMethod: GET,
			wantUrl:    "https://example.com?test=hello",
			wantHeader: http.Header{"Content-Type": []string{"application/json"}},
			wantBody:   []byte(`{"key":"value"}`),
		},
		{
			name:       "old.",
			fields:     fields{RequestType: Json, Method: GET, Body: map[string]string{"key": "value"}, URL: "https://example.com?test=hello"},
			wantMethod: GET,
			wantUrl:    "https://example.com?test=hello",
			wantHeader: http.Header{"Content-Type": []string{"application/json"}},
			wantBody:   []byte(`{"key":"value"}`),
		},
		{
			name:       "post use url query body.",
			fields:     fields{RequestType: UrlQuery, Method: POST, Body: map[string]string{"key": "value"}},
			wantMethod: POST,
			wantHeader: http.Header{},
			wantBody:   []byte(`key=value`),
		},
		{
			name:       "post use url query body NOT append to url.",
			fields:     fields{RequestType: UrlQuery, Method: POST, Body: map[string]string{"key": "value"}, URL: "http://example.com?test=hello"},
			wantMethod: POST,
			wantUrl:    "http://example.com?test=hello",
			wantHeader: http.Header{},
			wantBody:   []byte(`key=value`),
		},
		{
			name:       "use xml body",
			fields:     fields{RequestType: Xml, Method: POST, Body: XMlContent{Key: "value"}},
			wantMethod: POST,
			wantUrl:    "",
			wantHeader: http.Header{"Content-Type": []string{"application/xml"}},
			wantBody:   []byte(`<XMlContent><key>value</key></XMlContent>`),
		},
		{
			name:       "url query accept number",
			fields:     fields{RequestType: UrlQuery, Method: GET, Body: map[string]interface{}{"number": 1.234}, URL: "http://example.com?test=hello"},
			wantMethod: GET,
			wantUrl:    `http://example.com?test=hello&number=1.234`,
			wantHeader: http.Header{},
			wantBody:   nil,
		},
		{
			name:       "url query accept struct",
			fields:     fields{RequestType: UrlQuery, Method: GET, Body: map[string]interface{}{"struct": map[string]string{"sk": "sv"}}, URL: "http://example.com?test=hello"},
			wantMethod: GET,
			wantUrl:    `http://example.com?test=hello&struct={"sk":"sv"}`,
			wantHeader: http.Header{},
			wantBody:   nil,
		},
		{
			name: "Set xml request type will reset post default type.",
			fields: fields{
				Method:      POST,
				Body:        XMlContent{Key: "value"},
				RequestType: Xml,
			},
			wantMethod: POST,
			wantHeader: http.Header{"Content-Type": []string{"application/xml"}},
			wantBody:   []byte(`<XMlContent><key>value</key></XMlContent>`),
		},
		{
			name: "Set urlquery request type will reset post default type..",
			fields: fields{
				Method:      POST,
				Body:        map[string]string{"key": "value"},
				RequestType: UrlQuery,
			},
			wantMethod: POST,
			wantHeader: http.Header{},
			wantBody:   []byte(`key=value`),
		},
		{
			name:       "get method's string body",
			fields:     fields{Method: GET, Body: "single"},
			wantMethod: GET,
			wantUrl:    "?single=",
			wantHeader: http.Header{},
			wantBody:   nil,
		},
		{
			name:       "use string to post",
			fields:     fields{Method: POST, Body: "this is string."},
			wantMethod: POST,
			wantUrl:    "",
			wantHeader: http.Header{"Content-Type": []string{"application/json"}},
			wantBody:   []byte(`this is string.`),
		},
		{
			name:       "use string to xml post",
			fields:     fields{Method: POST, Body: "this is string.", RequestType: Xml},
			wantMethod: POST,
			wantUrl:    "",
			wantHeader: http.Header{"Content-Type": []string{"application/xml"}},
			wantBody:   []byte(`this is string.`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{
				Method:      tt.fields.Method,
				URL:         tt.fields.URL,
				Header:      tt.fields.Header,
				Body:        tt.fields.Body,
				RequestType: tt.fields.RequestType,
				Ctx:         tt.fields.Ctx,
			}
			gotMethod, gotUrl, gotHeader, gotBody := req.build()
			if gotMethod != tt.wantMethod {
				t.Errorf("Request.build() gotMethod = %v, want %v", gotMethod, tt.wantMethod)
			}
			if gotUrl != tt.wantUrl {
				t.Errorf("Request.build() gotUrl = %v, want %v", gotUrl, tt.wantUrl)
			}
			if !reflect.DeepEqual(gotHeader, tt.wantHeader) {
				t.Errorf("Request.build() gotHeader = %v, want %v", gotHeader, tt.wantHeader)
			}
			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("Request.build() gotBody = %v, want %v", string(gotBody), string(tt.wantBody))
			}
		})
	}
}

func TestRequest_SetRequestType(t *testing.T) {
	type fields struct {
		Method      HttpMethod
		URL         string
		Header      http.Header
		Body        interface{}
		RequestType BodyType
		Ctx         context.Context
	}
	type args struct {
		bt BodyType
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Request
	}{
		{
			name:   "can set body type.",
			fields: fields{RequestType: Json},
			args:   args{bt: Xml},
			want:   &Request{RequestType: Xml},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{
				Method:      tt.fields.Method,
				URL:         tt.fields.URL,
				Header:      tt.fields.Header,
				Body:        tt.fields.Body,
				RequestType: tt.fields.RequestType,
				Ctx:         tt.fields.Ctx,
			}
			if got := req.SetRequestType(tt.args.bt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request.SetRequestType() = %v, want %v", got, tt.want)
			}
		})
	}
}
