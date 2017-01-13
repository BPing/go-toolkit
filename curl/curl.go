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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BPing/Golib/curl/client"
	"net/http"
	"net/url"
	"strings"
	"io"
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
	// http config
	HttpConfig
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
		// 请求中的body 数据 来自 curl.Data或者curl.Body
		var bodyData io.Reader
		if curl.Body != nil && string(curl.Body) != "" {
			// body不为nil,则params附带到Url上
			bodyData = strings.NewReader(string(curl.Body))
		} else if curl.Data != nil && len(curl.Data) > 0 {
			contentType, _ := curl.Headers["Content-type"]
			switch strings.ToLower(contentType) {
			// json格式
			case "application/json":
				curl.Headers["Content-type"] = "application/json"
				var dataJson []byte
				dataJson, err = json.Marshal(curl.Data)
				if err != nil {
					return
				}
				bodyData = strings.NewReader(string(dataJson))

			//头部参数设置为form data 类型
			case "application/x-www-form-urlencoded":
				fallthrough
			default:
				curl.Headers["Content-type"] = "application/x-www-form-urlencoded"
				dataForm := url.Values{}
				for key, val := range curl.Data {
					dataForm.Add(key, val)
				}
				bodyData = strings.NewReader(dataForm.Encode())
			}
		}
		//fmt.Println(curl.Method, curl.Url + queryParam, bodyData)
		req, err = http.NewRequest(curl.Method, curl.Url+queryParam, bodyData)
	}

	// set header
	for key, val := range curl.Headers {
		req.Header.Add(key, val)
	}

	return req, err
}

func (curl *CurlRequest) String() string {
	return fmt.Sprintf("\n Url:%s, \n Method:%s,\n Header:%#v,\n Params:%#v,\n Data:%#v,\n Body:%v \n", curl.Url, curl.Method, curl.Headers, curl.Params, curl.Data, string(curl.Body))
}

//----------------------------------------------------------------------------------------------------------------------

// curl
// @url string 请求Uri
// @method string 方法。GET，POST，PUT等
// @params map[string]string 参数。?a=b
// @header map[string]string 头部信息
// @body   []byte
// @Deprecated 建议使用 HttpCurl
func Curl(url, method string, params, header map[string]string, body []byte) (resp *client.Response, err error) {
	resp, err = CurlWithClient(url, method, params, header, body, client.DefaultClient)
	return
}

// 自定义的client执行curl请求
// @Deprecated 建议使用 HttpCurl
func CurlWithClient(url, method string, params, header map[string]string, body []byte, c *client.Client) (resp *client.Response, err error) {
	if c == nil {
		err = errors.New("*client.Client is nil")
	}
	curlReq := &CurlRequest{
		HttpConfig: HttpConfig{
			Url:     url,
			Method:  method,
			Params:  params,
			Data:    params,
			Headers: header,
			Body:    body},
	}

	resp, err = c.DoRequest(curlReq)
	return
}

// http config
type HttpConfig struct {
	client *client.Client

	// 方法。GET，POST，PUT等
	Method string
	// 请求Uri
	Url string
	// query 参数
	Params map[string]string
	// body内容 from或者json
	Data    map[string]string
	Headers map[string]string
	// 如果Body不为nil，则会覆盖Data数据，也就是说Body优先级高于Data
	Body []byte
}

// http 请求
func HttpCurl(config HttpConfig) (resp *client.Response, err error) {
	clientTmp := config.client
	if clientTmp == nil {
		clientTmp = client.DefaultClient
	}
	curlReq := &CurlRequest{
		HttpConfig: config,
	}
	resp, err = clientTmp.DoRequest(curlReq)
	return
}
