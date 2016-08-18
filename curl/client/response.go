package client

import (
	"net/http"
	"compress/gzip"
	"io/ioutil"
	"encoding/json"
	"encoding/xml"
	"os"
	"io"
)

type ResponseFormat string

const (
	JSONResponseFormat = ResponseFormat("JSON")
	XMLResponseFormat = ResponseFormat("XML")
)

// 封装标准库中的Response
// 方便处理响应内容信息
type Response struct {
	*http.Response
	body []byte //缓存响应的Response的body字节内容
}

// 响应的Response的body字节内容保存到文件中去(文件请求)
func (resp *Response) ToFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if resp.Response.Body == nil {
		return nil
	}
	defer resp.Response.Body.Close()
	_, err = io.Copy(f, resp.Response.Body)
	return err
}

// 返回响应的Response的body字节内容
func (resp *Response) Bytes() ([]byte, error) {
	if resp.body != nil {
		return resp.body, nil
	}

	if resp.Response.Body == nil {
		return nil, nil
	}

	var err error

	defer resp.Response.Body.Close()
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		resp.body, err = ioutil.ReadAll(reader)
	} else {
		resp.body, err = ioutil.ReadAll(resp.Body)
	}

	return resp.body, err
}

// 将响应的Response的body字节内容以JSON格式转化
func (resp *Response) ToJSON(v interface{}) error {
	data, err := resp.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// 将响应的Response的body字节内容以XML格式转化
func (resp *Response) ToXML(v interface{}) error {
	data, err := resp.Bytes()
	if err != nil {
		return err
	}
	return xml.Unmarshal(data, v)
}