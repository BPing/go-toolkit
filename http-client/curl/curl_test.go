package curl

import (
	"fmt"
	"testing"
)

func TestCurl(t *testing.T) {
	respmap := make(map[string]interface{})
	resp, err := HttpCurl(HttpConfig{
		Method: GET,
		Url:    "http://www.weather.com.cn/data/cityinfo/101190408.html",
	})
	resp.ToJSON(&respmap)
	if nil != err {
		t.Fatal("JSON", err)
	}
	fmt.Println(respmap)

	resp, err = Do("http://www.weather.com.cn/data/cityinfo/101190408.html", GET, nil, nil, nil)
	resp.ToJSON(&respmap)
	if nil != err {
		t.Fatal("JSON", err)
	}
	fmt.Println(respmap)
}
