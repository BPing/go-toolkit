package client

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

//
//  请求接口
//    Client接受处理的请求接口
type Request interface {
	//返回*http.Request
	HttpRequest() (*http.Request, error)
	//返回请求相关内容格式化字符串
	String() string
	// 获取超时时间
	// =0代表此请求不启用超时设置
	// <0代表默认使用全局
	// >0代表自定义超时时间
	// @deprecated  2017-09-26
	GetTimeOut() time.Duration
	//克隆
	Clone() interface{}
	//设置尝试次数
	SetReqCount(reqCount int)
}

// 请求基类
//    实现请求接口的基类，所有请求对象继承必须继承此基类
type BaseRequest struct {
	reqCount int
}

func (b *BaseRequest) HttpRequest() (*http.Request, error) {
	return nil, errors.New("implement Interface's Method::HttpRequest")
}

func (b *BaseRequest) String() string {
	return fmt.Sprintf("ReqCount:%d \n", b.reqCount)
}

func (b *BaseRequest) GetTimeOut() time.Duration {
	return -1
}

func (b *BaseRequest) Clone() interface{} {
	new_obj := (*b)
	return &new_obj
}

func (b *BaseRequest) SetReqCount(reqCount int) {
	b.reqCount = reqCount
}

func (b *BaseRequest) getReqCount() int {
	return b.reqCount
}
