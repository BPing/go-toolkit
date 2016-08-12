package client

import (
	"net/http"
	"fmt"
	"errors"
)

//
//  请求接口
//    所有请求对象继承的接口，也是Client接受处理的请求接口
type Request interface {
	//返回*http.Request
	HttpRequest() (*http.Request, error)
	//返回请求相关内容格式化字符串
	String() string
	//克隆
	Clone() interface{}
}

type BaseRequest struct {

}

func (b *BaseRequest) HttpRequest() (req *http.Request, err error) {
	return nil, errors.New("Implement Interface's Method::HttpRequest")
}

func (b *BaseRequest) String() string {
	return fmt.Sprintf("Request:%v", b)
}

func (b *BaseRequest) Clone() interface{} {
	new_obj := (*b)
	return &new_obj
}