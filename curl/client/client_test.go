package client

import (
	"fmt"
	"net/http"
	"testing"
	"time"
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
	client.SetRetryCount(3)
	req := &TestRequest{RequestURL: "https://www.baidu.co/"}
	_, err := client.DoRequest(req)
	if nil == err {
		t.Fatal("weather Request", err)
	}

	if req.reqCount != 3 {
		t.Fatal("SetRetryCount", err)
	}
}

func TestClient_SetRecord(t *testing.T) {
	logMsg := ""
	client := NewClient("test", nil)
	client.SetRecord(func(tag, msg string) {
		logMsg = msg
	})
	client.SetRetryCount(1)
	req := &TestRequest{RequestURL: "http://www.weather.com.cn/aa"}
	client.DoRequest(req)
	if logMsg == "" {
		t.Fatal("set record fail")
	}
}

func TestClient_SetSlowReqLong(t *testing.T) {
	slowReqFlag := false
	client := NewClient("test", nil)
	client.SetRecord(func(tag, msg string) {
		if tag == SlowReqRecord {
			slowReqFlag = true
		}
		fmt.Printf("tag:%s msg:%s ", tag, msg)
	})
	client.SetSlowReqLong(time.Millisecond)
	req := &TestRequest{RequestURL: "http://www.weather.com.cn/data/cityinfo/101190408.html"}
	_, err := client.DoRequest(req)
	if nil != err {
		t.Fatal("weather Request", err)
	}
	if !slowReqFlag {
		t.Fatal("SetSlowReqLong fail")
	}
}
