package core

import (
	"fmt"
	"net/http"
	"testing"
	"errors"
	"context"
)

type TestRequest struct {
	BaseRequest
	RequestURL string
}

func (b *TestRequest) HttpRequest() (*http.Request, error) {
	httpReq, err := http.NewRequest("GET", b.RequestURL, nil)
	return httpReq, err
}

func (b *TestRequest) String() string {
	return ""
}

type youdao struct {
	errorCode int    `xml:"errorCode"`
	query     string `xml:"query"`
}

func TestClient(t *testing.T) {
	end := make(chan int)
	client := NewClient("test", nil)
	client.SetDebug(true)
	go func() {
		respmap := make(map[string]interface{})
		req := &TestRequest{RequestURL: "http://www.weather.com.cn/data/cityinfo/101190408.html"}
		resp, err := client.DoRequest(req)
		if nil != err {
			t.Fatal("weather Request", err)
		}

		resp.ToJSON(&respmap)

		if nil != err {
			t.Fatal("JSON", err)
		}

		fmt.Println(respmap)

		end <- 1
	}()
	respxml := youdao{}
	req := &TestRequest{RequestURL: "http://fanyi.youdao.com/openapi.do?keyfrom=cbping&key=1366735279&type=data&doctype=xml&version=1.1&q=%E8%A6%81%E7%BF%BB%E8%AF%91%E7%9A%84%E6%96%87%E6%9C%AC"}
	resp1, err := client.DoRequest(req)
	if nil != err {
		t.Fatal("youdao Request", err)
	}
	err = resp1.ToXML(&respxml)

	if nil != err {
		t.Fatal("XML", err)
	}

	fmt.Println(respxml)

	<-end
}

func TestClient_SetRetryCount(t *testing.T) {
	client := NewClient("test", nil)
	client.SetDebug(true)
	client.SetMaxBadRetryCount(3)
	req := &TestRequest{RequestURL: "https://www.baidu.co/"}
	_, err := client.DoRequest(req)
	if nil == err {
		t.Fatal("weather Request", err)
	}

	if req.reqCount != 3 {
		t.Fatal("SetRetryCount", err)
	}
}

type TestHook struct {
}

func (log *TestHook) BeforeRequest(req Request, client Client) error {
	return errors.New("some error happen")
}

func (log *TestHook) AfterRequest(cErr error, req Request, client Client) {

}

func TestClient_AppendHook(t *testing.T) {
	client := NewClient("test", nil)
	testHook := &TestHook{}
	client.AppendHook(testHook)
	req := &TestRequest{RequestURL: "https://www.baidu.co/"}
	_, err := client.DoRequest(req)
	if err == nil || err.Error() != "some error happen" {
		t.Fatal("AppendHook", err)
	}
}

func TestNewClientCtx(t *testing.T) {
	ctx,cancelFunc:=context.WithCancel(context.Background())
	client:=NewClientCtx(ctx,"test", nil)
	cancelFunc()
	req := &TestRequest{RequestURL: "https://www.baidu.com/"}
	_, err := client.DoRequest(req)
	if err == nil || err.Error() != "context canceled"{
		t.Fatal("NewClientCtx", err)
	}
}