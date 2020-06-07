# simple-go-http

一个直接使用 `Golang struct` 发送 `http` 请求的客户端，更适合于 `RESTFul` 请求。

### 特性

- 快，默认客户端底层使用 `fasthttp`
- 简单
- 容易操作 `JSON`、`XML` 和 `URL parameters`
- 支持设置超时

### 安装

```shell
go get github.com/SmallTianTian/simple-go-http
```

### 例子

[基础用法](#基础用法)
[请求/响应默认行为](#请求/响应默认行为)
[添加钩子](#添加钩子)

### 基础用法

```golang
import "client" github.com/SmallTianTian/simple-go-http
```

#### 1. 极简 GET 请求

```golang
err := client.NewRequest()
             .Get("https://example.com")
```

#### 2. 携带请求体的 POST 请求

```golang
err := client.NewRequest()
             .POST("https://example.com", "this is string body.")
```

#### 3. 设置请求头

```golang
client.NewRequest()
      .SetHeader("header key", "header value")
      .SetHeader("other key", "other value")
```

#### 4. 更多的 RESTFUL 请求，详见[枚举值](./enum.go#30)

```golang
client.NewRequest()
      .HttpMethod(client.DELETE)
      .Url("https://example.com")
      .SetHeader("other key", "other value")
```

#### 5. 设置上下文

```golang
client.NewRequest()
      .Context(ctx)
```

#### 6. 设置请求体格式

```golang
type XMlContent struct {
    Key string `xml:"key"`
}

// url: https://example.com
// body: `<XMlContent><key>value</key></XMlContent`
// header: Content-Type: application/xml
client.NewRequest()
      .Post("https://example.com", XMlContent{Key: "value"})
      .SetRequestType(Xml)
```

#### 7. 响应体自动写入 struct

```golang
var res struct {
    RespKey string `json:"resp_key"`
}

// 1. json response
client.NewRequest()
      .Post("https://example.com", nil) // will return `{"resp_key": "this is fake"}`
      .Do(client.NewJsonResponse(&res)) // res{RespKey: "this is fake"}
```

### 请求/响应默认行为

#### 1. 利用 HttpMethod 的默认功能

```golang
// 1. Get 默认将 Body 转为 UrlQuery 并追加到 URL 后。
client.NewRequest()
      .Get("https://example.com")
      .Body("single") // url: https://example.com?single=

client.NewRequest()
      .Get("https://example.com")
      .Body(map[string]string{"key": "value"}) // url: https://example.com?key=value

// 2. POST 默认将 Body 转为 Json 并设置 json 请求头
// url: https://example.com
// body: `single`
// header: Content-Type: application/json
client.NewRequest()
      .Post("https://example.com", "single")

// url: https://example.com
// body: `{"key": "value"}`
// header: Content-Type: application/json
client.NewRequest()
      .Post("https://example.com", map[string]string{"key": "value"})
```

#### 2. 利用响应头来自动写入 strcut

```golang
var res struct {
    RespKey string `json:"resp_key"`
}

// 1. json response
client.NewRequest()
      // will return:
      // header: Content-Type: application/json
      // body: `{"resp_key": "this is fake"}`
      .Post("https://example.com", nil)
      .Do(client.NewDefaultResponse(&res)) // res{RespKey: "this is fake"}
```

### 添加钩子

```golang
// 可以自由的添加钩子，在请求开始前和请求结束后，按顺序执行。
setReqHook := func(req *Request, resp *Response) {
    req.Body = "hook"
}

setRespHook := func(req *Request, resp *Response) {
    resp.Header.Set("hook", "value")
}

client.NewRequest()
      .Get("https://example.com")
      .Do(client.NewDefaultResponse(&res), setReqHook, setRespHook)
// req:
// url: https://example.com?hook=

// resp:
// header: hook:value
```
