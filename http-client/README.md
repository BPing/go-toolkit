# httpClient

# core

* Client：封装http.Client.修饰http.Client。处理请求过程加入一些必要的自定义处理。如：失败尝试。
          亦可以通过钩子，添加额外功能。
* Request：接口类型。
* Response：封装http.Response.修饰http.Response。集成一些常用的处理响应内容方法。如：`ToJson()` 返回json格式内容

# hook

## 系统钩子

* LogHook 日志记录，包括慢请求

```go
	record := func(tag, msg string) {
		logMsg = msg
	}
	core.AppendHook(NewLogHook(time.Duration(0), record))
```

* CircuitHook 断路器（熔断处理）

```go
	settings := CircuitSettings{
		Name: "test",
		ReadyToTrip: func(counts Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
		OnStateChange: func(name string, from State, to State) {
			//fmt.Println(name, from, to)
		},
		Interval:    time.Second * 20,
		Timeout:     time.Second * 2,
		MaxRequests: 6,
	}
	circuitHook := NewCircuitHook(settings)
	core.AppendHook(circuitHook)
```

## 自定义钩子

```go
  type Hook interface {
  	// 请求处理前执行
  	// 如果返回错误
  	// 将提前终止请求
  	// 并将此错误返回
  	BeforeRequest(req Request, client Client) error

  	// 请求处理后执行
  	// @params err 请求处理错误信息，如果不为nil，代表请求失败
  	AfterRequest(cErr error, req Request, client Client)
  }
```

# curl

* 发起请求
```go
curl.Do(url, method string, params, header map[string]string, body []byte)
```
