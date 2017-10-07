package hook

import (
	"net/http"
	"github.com/BPing/go-toolkit/http-client/core"
	"testing"
	"fmt"
	"time"
	"errors"
)

type TestCircuitRequest struct {
	core.BaseRequest
	RequestURL string
}

func (b *TestCircuitRequest) HttpRequest() (*http.Request, error) {
	httpReq, err := http.NewRequest("GET", b.RequestURL, nil)
	return httpReq, err
}

func TestNewCircuitBreaker(t *testing.T) {
	settings := CircuitSettings{
		Name: "test",
		ReadyToTrip: func(counts Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
		OnStateChange: func(name string, from State, to State) {
			fmt.Print(name, from, to)
		},
		Interval:    time.Second * 2,
		Timeout:     time.Second * 2,
		MaxRequests: 6,
	}

	excuteFunc := func(cb *CircuitBreaker, req func() (interface{}, error)) (interface{}, error) {
		generation, err := cb.beforeRequest()
		if err != nil {
			return nil, err
		}

		defer func() {
			e := recover()
			if e != nil {
				cb.afterRequest(generation, false)
				panic(e)
			}
		}()

		result, err := req()
		cb.afterRequest(generation, err == nil)
		return result, err
	}

	reqFailFunc := func() (interface{}, error) {
		return nil, errors.New("some error happen")
	}

	reqSuccessFunc := func() (interface{}, error) {
		return nil, nil
	}

	reqManyFunc := func() (interface{}, error) {
		time.Sleep(time.Second)
		return nil, nil
	}

	breaker := NewCircuitBreaker(settings)

	excuteFunc(breaker, reqFailFunc)
	time.Sleep(time.Second * 2)
	for i := 0; i < 2; i++ {
		excuteFunc(breaker, reqFailFunc)
		time.Sleep(time.Millisecond * 100)
	}
	if breaker.State() != StateClosed {
		t.Fatal("StateOpen", "Interval to reset，so state of breaker already is closed")
	}
	excuteFunc(breaker, reqFailFunc)
	if breaker.State() != StateOpen {
		t.Fatal("StateOpen", "from closed to open")
	}

	if _, err := excuteFunc(breaker, reqManyFunc); err != ErrOpenState {
		t.Fatal("StateOpen", "circuit breaker is open")
	}
	time.Sleep(time.Second * 3)
	if breaker.State() != StateHalfOpen {
		t.Fatal("StateHalfOpen", "Timeout: from open to halfOpen")
	}

	// 在halfOpen状态下请求失败，再次进入open状态
	excuteFunc(breaker, reqFailFunc)
	if breaker.State() != StateOpen {
		t.Fatal("StateOpen", "req fail:from halfOpen to open")
	}

	time.Sleep(time.Second * 2)
	if breaker.State() != StateHalfOpen {
		t.Fatal("StateHalfOpen", "from open to halfOpen again")
	}

	// 过多的请求
	for i := 0; i < 20; i++ {
		go func() {
			// ErrTooManyRequests
			fmt.Println(excuteFunc(breaker, reqManyFunc))
		}()
		time.Sleep(time.Millisecond * 100)
	}
	// 成功请求达到设置次数
	for i := 0; i < 6; i++ {
		excuteFunc(breaker, reqSuccessFunc)
		time.Sleep(time.Millisecond * 100)
	}
	if breaker.State() != StateClosed {
		t.Fatal("StateClosed", "the all req success,from open to closed")
	}

}

func TestNewCircuitHook(t *testing.T) {
	settings := CircuitSettings{
		Name: "test",
		ReadyToTrip: func(counts Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
		OnStateChange: func(name string, from State, to State) {
			//fmt.Println(name, from, to)
		},
		Interval:    time.Second * 20,
		Timeout:     time.Second * 2,
		MaxRequests: 6,
	}

	handleFailFunc := func(cErr error, req core.Request) error {
		return errors.New("some error happen")
	}

	handleSuccessFunc := func(cErr error, req core.Request) error {
		return nil
	}

	circuitHook := NewCircuitHook(settings)
	c := core.NewClient("test", nil)
	c.AppendHook(circuitHook)
	circuitHook.SetHandleCErr(handleFailFunc)
	req := &TestRequest{RequestURL: "http://127.0.0.1:6101/"}
	fmt.Println(c.DoRequest(req))

	var err error
	for i := 0; i < 3; i++ {
		_, err = c.DoRequest(req)
	}
	_, err = c.DoRequest(req)
	if err != ErrOpenState {
		t.Fatal("StateOpen", "from closed to open")
	}

	time.Sleep(time.Second * 3)
	circuitHook.SetHandleCErr(handleSuccessFunc)
	_, err = c.DoRequest(req)
	fmt.Println(err)
	if err == ErrOpenState {
		t.Fatal("StateHalfOpen", "Timeout: from open to halfOpen")
	}

	// 过多的请求
	for i := 0; i < 10; i++ {
		// ErrTooManyRequests
		go func() {
			_, err = c.DoRequest(req)
			//fmt.Println(err)
		}()
		time.Sleep(time.Millisecond * 100)
	}
	time.Sleep(time.Second * 5)
	// 成功请求达到设置次数
	for i := 0; i < 6; i++ {
		_, err = c.DoRequest(req)
		fmt.Println(err)
	}
	if err == ErrOpenState || err == ErrTooManyRequests {
		t.Fatal("StateClosed", "the all req success,from open to closed")
	}
}
