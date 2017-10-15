// +build go1.8

package core

import (
	"testing"
	"context"
)

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