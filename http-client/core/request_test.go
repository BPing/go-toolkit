package core

import (
	"testing"
	"time"
)

func TestBaseRequest(t *testing.T) {
	base := &BaseRequest{}

	if _, ok := base.HookData("testKey"); ok {
		t.Fatal("HookData false")
	}

	base.SetHookData("testKey", "testVal")
	if val, ok := base.HookData("testKey"); !ok && val.(string) != "testVal" {
		t.Fatal("HookData", val)
	}

	base.setReqLongTime(time.Second)
	if base.ReqLongTime() != time.Second {
		t.Fatal("ReqLongTime", base.ReqLongTime())
	}

	base.setReqCount(3)
	if base.ReqCount() != 3 {
		t.Fatal("ReqCount", base.ReqCount())
	}
	cloneReq := base.Clone().(*BaseRequest)
	if cloneReq.ReqCount() != 3 {
		t.Fatal("cloneReq ReqCount", cloneReq.ReqCount())
	}

	cloneReq.setResponse(nil)
	if cloneReq.Response() != nil {
		t.Fatal("cloneReq Response", cloneReq.Response())
	}

	t.Log(base.String())
	t.Log(base.ServerName())
	t.Log(base.HttpRequest())
	t.Log(base.TimeOut())
}
