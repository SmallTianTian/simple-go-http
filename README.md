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
[辅助调试](#辅助调试)

### 基础用法

```golang
import (
    "context"
    "time"
    client "github.com/SmallTianTian/simple-go-http"
)

var resp string
err := client.NewRequest().
		// 1. 设置上下文
		Context(context.Background()).

		// 2. 极简 GET 请求
		Get("https://example.com").
		// 2.1 携带请求体的 POST 请求
		// 	   Post("https://post.example.com", "this is request body.").
		// 2.2 更多的 RESTFUL 请求，详见[枚举值](./enum.go#30)
		//     HttpMethod(client.DELETE).
		//     Url("https://example.com").

		// 3. 设置请求头
		SetHeader("header key", "header value").
		// SetHeader("other key", "other value").

		// 4. 设置请求体格式
		SetRequestType(client.Xml).

		// 5. 为 `Request` 设置单独的 timeout
		SetTimeout(3 * time.Second).
		// 6. 响应体自动写入 struct
		Do(client.NewJsonResponse(&resp))
```



### 请求/响应默认行为

#### 1. 利用 HttpMethod 的默认功能

```golang
// 1. Get 默认将 Body 转为 UrlQuery 并追加到 URL 后。
client.NewRequest().
       Get("https://example.com").
       HttpBody("single") // url: https://example.com?single=

client.NewRequest().
       Get("https://example.com").
       HttpBody(map[string]string{"key": "value"}) // url: https://example.com?key=value

// 2. POST 默认将 Body 转为 Json 并设置 json 请求头
// url: https://example.com
// body: `single`
// header: Content-Type: application/json
client.NewRequest().
       Post("https://example.com", "single")

// url: https://example.com
// body: `{"key": "value"}`
// header: Content-Type: application/json
client.NewRequest().
       Post("https://example.com", map[string]string{"key": "value"})
```

#### 2. 利用响应头来自动写入 strcut

```golang
var res struct {
    RespKey string `json:"resp_key"`
}

// 1. json response
client.NewRequest().
       // will return:
       // header: Content-Type: application/json
       // body: `{"resp_key": "this is fake"}`
       Post("https://example.com", nil).
       Do(client.NewDefaultResponse(&res)) // res{RespKey: "this is fake"}
```

### 添加钩子

```golang
// 可以自由的添加钩子，在请求开始前和请求结束后，按顺序执行。
setReqHook := func(req *client.Request, resp *client.Response) {
    req.Body = "hook"
}

setRespHook := func(req *client.Request, resp *client.Response) {
    resp.Header.Set("hook", "value")
}

var resp string
client.NewRequest().
       Get("https://example.com").
       Do(client.NewDefaultResponse(&resp), setReqHook, setRespHook)
// req:
// url: https://example.com?hook=

// resp:
// header: hook:value
```

### 辅助调试

```golang
client.OpenDebug() // 开启 Debug 调试信息
```
