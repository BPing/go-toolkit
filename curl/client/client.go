// client 包 核心代码包
// @author cbping
package client

import (
	"errors"
	"net/http"
	"time"
	"net/url"
	"fmt"
)

const (
	SlowReqRecord = "SlowReqRecord"
	ReqRecord = "ReqRecord"
	ErrorReqRecord = "ErrorReqRecord"
)

func init() {
	SetDefaultClient("", http.DefaultClient)
}

//----------------------------------------------------------------------------------------------------------------------

//
//  客户端
//  处理http请求
type Client struct {
	// 采用默认&http.Client{}
	*http.Client

	//
	UserAgent   string

	// 超过SlowReqLong时间长度的请求，将记录为慢请求
	// 默认为2秒
	SlowReqLong time.Duration

	// 函数参数
	// 记录信息；如日志记录
	Record      func(tag, msg string)

	//版本号
	Version     string
	//
	debug       bool
}

func (c *Client) SetDebug(debug bool) {
	c.debug = debug
}

func (c *Client) SetRecord(record func(tag, msg string)) {
	c.Record = record
}

// 设置代理
// example:
//
//	func(req *http.Request) (*url.URL, error) {
// 		u, _ := url.ParseRequestURI("http://127.0.0.1:8118")
// 		return u, nil
// 	}
func (c *Client) SetProxy(proxy func(*http.Request) (*url.URL, error)) {
	//TODO::默认http.Client或者默认http.Transport时，是否值得改变代理（影响其他请求）？
	if nil != c.Client && nil != c.Client.Transport {
		c.Client.Transport.(http.Transport).Proxy = proxy
	}

	return
}

// 处理请求
func (c *Client) DoRequest(req Request) (resp *Response, err error) {
	if nil == c.Client {
		c.Client = http.DefaultClient
	}

	defer func() {
		if nil != err&&nil != c.Record {
			c.Record(ErrorReqRecord, fmt.Sprintf("query:: %s errorr:: %v) ", req.String(), err))
			err = clientError(err)
		}
	}()

	if nil == req {
		return nil, errors.New("Request is nil")
	}

	httpReq, err := req.HttpRequest()
	if nil != err {
		return nil, err
	}

	//必要头部信息设置
	httpReq.Header.Set("User-Agent", `Bping-Curl-` + c.UserAgent + "/" + c.Version)

	t0 := time.Now()
	httpResp, err := c.Client.Do(httpReq)
	t1 := time.Now()
	if nil != err {
		return nil, err
	}

	if nil != c.Record {
		if t1.Sub(t0) >= c.SlowReqLong {
			c.Record(SlowReqRecord, req.String())
		}
		c.Record(ReqRecord, fmt.Sprintf("http query:: %s %d (%v) ", req.String(), httpResp.StatusCode, t1.Sub(t0)))
	}

	return &Response{Response:httpResp}, nil

}

func NewClient(title string, client *http.Client) (c *Client) {
	return &Client{
		Client:client,
		Version:Version,
		UserAgent:title,
		debug:false,
		SlowReqLong:2 * time.Second,
	}
}

func clientError(err error) error {
	if nil == err {
		return nil
	}
	return errors.New("Bping-Curl-Client-Failure:" + err.Error())
}

//----------------------------------------------------------------------------------------------------------------------

var DefaultClient *Client

// 设置DefaultClient
func SetDefaultClient(title string, client *http.Client) {
	DefaultClient = NewClient(title, client)
}

// 设置代理
// example:
//
//	func(req *http.Request) (*url.URL, error) {
// 		u, _ := url.ParseRequestURI("http://127.0.0.1:8118")
// 		return u, nil
// 	}
func SetProxy(proxy func(*http.Request) (*url.URL, error)) {
	DefaultClient.SetProxy(proxy)
}

// 记录
func SetRecord(record func(tag, msg string)) {
	DefaultClient.SetRecord(record)
}

// 处理请求，内部调用DefaultClient
func DoRequest(req Request) (Response, error) {
	DefaultClient.Client.Transport
	return DefaultClient.DoRequest(req)
}