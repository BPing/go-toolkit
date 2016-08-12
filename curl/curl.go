package curl

import (
	"github.com/BPing/Golib/curl/client"
	"net/http"
	"net/url"
	"strings"
	"fmt"
)

const (
	GET = "GET"
	POST = "POST"
	PUT = "PUT"
	DELETE = "DELETE"
	HEAD = "HEAD"
)


//----------------------------------------------------------------------------------------------------------------------

//
//
type CurlRequest struct {
	client.BaseRequest
	// http
	Method  string
	Url     string

	Params  map[string]string
	Headers map[string]string

	// body不为nil,则params附带到Url上
	Body    []byte
}

//
func (curl *CurlRequest) HttpRequest() (req *http.Request, err error) {

	v := url.Values{}

	// set param
	for key, val := range curl.Params {
		v.Add(key, val)
	}

	curl.Method = strings.ToUpper(curl.Method)

	if curl.Method == GET {
		req, err = http.NewRequest(curl.Method, curl.Url + "?" + v.Encode(), nil)
	} else {
		if curl.Body != nil {
			// body不为nil,则params附带到Url上
			req, err = http.NewRequest(curl.Method, curl.Url + "?" + v.Encode(), strings.NewReader(string(curl.Body)))
			req.ContentLength = int64(len(curl.Body))
		}else {
			req, err = http.NewRequest(curl.Method, curl.Url, strings.NewReader(v.Encode()))
			req.ContentLength = int64(len([]byte(v.Encode())))
		}

	}

	// set header
	for key, val := range curl.Headers {
		req.Header.Add(key, val)
	}

	return req, err
}

func (curl *CurlRequest) String() string {
	return fmt.Sprintf("Url:%s,Method:%s,Header:%#v,Params:%#v,Body:%v", curl.Url, curl.Method, curl.Headers, curl.Params, nil != curl.Body)
}

//func (curl *CurlRequest) Clone() client.Request {
//	// var newBody []byte
//	//copy(newBody, curl.Body)
//	new_obj := (*curl)
//
//	return &new_obj
//}


//----------------------------------------------------------------------------------------------------------------------

//
//type Curl struct {
//	*client.Client
//}
//
//func (c *Curl)Query(url, method string, params, header map[string]string, body []byte) (resBody map[string]interface{}, resHeader map[string][]string, responseStatus string) {
//	if nil == c.Client {
//		c.Client = client.DefaultClient
//	}
//	curlReq := &CurlRequest{
//		Url:url,
//		Method:method,
//		Params:params,
//		Headers:header,
//		Body:body,
//	}
//
//	resp, err := c.Client.Query(curlReq)
//
//	if err != nil {
//		return nil, nil, "10003 remote req error::" + err.Error()
//	}
//
//	resp.ToJSON(&resBody)
//	resHeader = resp.Header
//	responseStatus = resp.Status
//	return
//}

// curl
func Curl(url, method string, params, header map[string]string, body []byte) (resBody map[string]interface{}, resHeader map[string][]string, responseStatus string) {

	curlReq := &CurlRequest{
		Url:url,
		Method:method,
		Params:params,
		Headers:header,
		Body:body,
	}

	resp, err := client.DoRequest(curlReq)

	if err != nil {
		return nil, nil, err.Error()
	}

	resp.ToJSON(&resBody)
	resHeader = resp.Header
	responseStatus = resp.Status
	return
}