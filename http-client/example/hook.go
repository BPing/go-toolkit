package example

import (
	"time"
	"fmt"
	"net/http"
	"net/url"
	"github.com/BPing/go-toolkit/http-client/core"
	"github.com/BPing/go-toolkit/http-client/hook"
)

type TestRequest struct {
	core.BaseRequest
	RequestURL string
}

func (b *TestRequest) HttpRequest() (*http.Request, error) {
	httpReq, err := http.NewRequest("GET", b.RequestURL, nil)
	return httpReq, err
}

func (b *TestRequest) ServerName() string {
	uri, _ := url.Parse(b.RequestURL)
	return uri.Hostname()
}

func SetLog() {
	logMsg := ""
	c := core.NewClient("test", nil)
	record := func(tag, msg string) {
		fmt.Print(tag, msg)
		logMsg = msg
	}
	c.AppendHook(hook.NewLogHook(3*time.Second, record))

	req := &TestRequest{RequestURL: "http://www.weather.baidu"}
	c.DoRequest(req)

}

func circuitTest() {
	settings := hook.CircuitSettings{
		Name: "test",
		ReadyToTrip: func(counts hook.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
		OnStateChange: func(name string, from hook.State, to hook.State) {
			fmt.Print(name, from, to, "\n")
		},
		Interval:    time.Second * 10,
		Timeout:     time.Second * 5,
		MaxRequests: 6,
	}

	circuitHook := hook.NewCircuitHook(settings)
	record := func(tag, msg string) {
		fmt.Print(tag, msg, "\n")
	}
	logHook := hook.NewLogHook(3*time.Second, record)
	c := core.NewClient("test", nil)
	c.AppendHook(logHook, circuitHook)
	//
	//end := make(chan int8)
	//go func() {
	for ; ; {
		req := &TestRequest{RequestURL: "http://127.0.0.1:6101/wechat/protection/common/offer_lost_protection_materials/?service_guid=DE6C2AE0205CC6620A55CCEF44A1B162"}
		fmt.Println(c.DoRequest(req))
		time.Sleep(time.Millisecond * 500)
	}
	//	end <- 1
	//}()
	//<-end
}