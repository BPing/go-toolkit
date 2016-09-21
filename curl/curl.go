// Copyright 2016  cbping. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//
// curl包
//
//   func Curl(url, method string, params, header map[string]string, body []byte) (resBody map[string]interface{}, resHeader map[string][]string, responseStatus string)
//
package curl

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/BPing/Golib/curl/client"
	"time"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	HEAD   = "HEAD"
)

//----------------------------------------------------------------------------------------------------------------------

//
//请求
// 实现Request接口
type CurlRequest struct {
	client.BaseRequest
	// http
	Method string
	Url    string

	Params  map[string]string
	Headers map[string]string

	// body不为nil,则params附带到Url上
	Body []byte

	// 超时时间
	// =0代表此请求不启用超时设置
	// <0代表默认使用全局
	// >0代表自定义超时时间
	Timeout time.Duration
}

//
func (curl *CurlRequest) HttpRequest() (req *http.Request, err error) {

	v := url.Values{}

	// set param
	for key, val := range curl.Params {
		v.Add(key, val)
	}

	curl.Method = strings.ToUpper(curl.Method)

	queryParam := ""
	if len(v) > 0 {
		queryParam = "?" + v.Encode()
	}

	if curl.Method == GET {
		req, err = http.NewRequest(curl.Method, curl.Url+queryParam, nil)
	} else {
		if curl.Body != nil && string(curl.Body) != "" {
			// body不为nil,则params附带到Url上
			req, err = http.NewRequest(curl.Method, curl.Url+queryParam, strings.NewReader(string(curl.Body)))
		} else {
			req, err = http.NewRequest(curl.Method, curl.Url, strings.NewReader(v.Encode()))
			//头部参数设置为form data 类型
			req.Header.Add("Content-type", "application/x-www-form-urlencoded")
		}

	}

	// set header
	for key, val := range curl.Headers {
		req.Header.Add(key, val)
	}

	return req, err
}

func (curl *CurlRequest) String() string {
	return fmt.Sprintf("\n Url:%s, \n Method:%s,\n Header:%#v,\n Params:%#v,\n Body:%v \n", curl.Url, curl.Method, curl.Headers, curl.Params, string(curl.Body))
}

func (curl *CurlRequest) GetTimeOut() time.Duration {
	return curl.Timeout
}

//func (curl *CurlRequest) Clone() client.Request {
//	// var newBody []byte
//	//copy(newBody, curl.Body)
//	new_obj := (*curl)
//
//	return &new_obj
//}

//----------------------------------------------------------------------------------------------------------------------

// curl
// @url string 请求Uri
// @method string 方法。GET，POST，PUT等
// @params map[string]string 参数。?a=b
// @header map[string]string 头部信息
// @body   []byte
func Curl(url, method string, params, header map[string]string, body []byte) (resp *client.Response, err error) {
	resp,err=CurlTimeout(url, method , params, header, body ,-1)
	return
}

// 超时处理
func CurlTimeout(url, method string, params, header map[string]string, body []byte,timeout time.Duration) (resp *client.Response, err error) {

	curlReq := &CurlRequest{
		Url:     url,
		Method:  method,
		Params:  params,
		Headers: header,
		Body:    body,
		Timeout: timeout,
	}

	resp, err = client.DoRequest(curlReq)
	return
}
