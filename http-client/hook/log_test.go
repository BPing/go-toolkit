package hook

import (
	"fmt"
	"time"
	"net/http"
	"testing"
	"github.com/BPing/go-toolkit/http-client/core"
)

type TestRequest struct {
	core.BaseRequest
	RequestURL string
}

func (b *TestRequest) HttpRequest() (*http.Request, error) {
	httpReq, err := http.NewRequest("GET", b.RequestURL, nil)
	return httpReq, err
}

func TestClient_SetRecord(t *testing.T) {
	logMsg := ""
	c := core.NewClient("test", nil)
	record := func(tag, msg string) {
		logMsg = msg
	}
	c.AppendHook(NewLogHook(time.Duration(0), record))

	req := &TestRequest{RequestURL: "http://www.weather.com.cn/aa"}
	c.DoRequest(req)
	if logMsg == "" {
		t.Fatal("set record fail")
	}
	logMsg = ""
	req = &TestRequest{RequestURL: "http://www.bai.kk"}
	c.DoRequest(req)
	if logMsg == "" {
		t.Fatal("set record fail")
	}
}

func TestClient_SetSlowReqLong(t *testing.T) {
	slowReqFlag := false
	c := core.NewClient("test", nil)
	record := func(tag, msg string) {
		if tag == SlowReqRecord {
			slowReqFlag = true
		}
		fmt.Printf("tag:%s msg:%s ", tag, msg)
	}
	c.AppendHook(NewLogHook(time.Millisecond, record))
	req := &TestRequest{RequestURL: "http://www.weather.com.cn/data/cityinfo/101190408.html"}
	_, err := c.DoRequest(req)
	if nil != err {
		t.Fatal("weather Request", err)
	}
	if !slowReqFlag {
		t.Fatal("SetSlowReqLong fail")
	}
}
