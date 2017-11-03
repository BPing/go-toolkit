# httpClient

# core

* Client：封装http.Client.修饰http.Client。处理请求过程加入一些必要的自定义处理。如：失败尝试。
          亦可以通过钩子，添加额外功能。
* Request：接口类型。
* Response：封装http.Response.修饰http.Response。集成一些常用的处理响应内容方法。如：`ToJson()` 返回json格式内容


### 配置

- `失败尝试次数`:默认2次
```go
 core.SetMaxBadRetryCount(2)
```

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
// Name  名字，请务必保障名字的唯一性
//
// MaxRequests Half-Open状态下允许通过的最大请求数
//
// Interval 重置时间间隔（Closed状态下有效）。如果为零，永远不重置。
//
// Timeout  超时时间（Open状态下有效）。
//          超时之后，状态将转变为Half-Open状态。
//          如果为零，默认为60秒
//
// ReadyToTrip   测试是否应该从Closed状态转变为Open状态。
//               true 表示可以转变，否则不可以。
//               如果不配置，则采用默认的。默认失败次数达到5次则进入Open状
//
// OnStateChange 状态变化将调用此方法。
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
package main

import  (
	"github.com/BPing/go-toolkit/http-client/curl"
	)

curl.Do( '', curl.POST , make(map[string]string), make(map[string]string) ,nil)
```

`或者`

```go
package main

import  (
	"github.com/BPing/go-toolkit/http-client/curl"
	)

// 如果Body不为nil，则会覆盖Data数据，也就是说Body优先级高于Data
curl.HttpCurl(HttpConfig{
			Url:     '',
			Method:  curl.POST,
			Params:  make(map[string]string),
			Data:    make(map[string]string),
			Headers: make(map[string]string),
			Body:    nil})
```