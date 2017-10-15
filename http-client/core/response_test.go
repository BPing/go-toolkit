package core

import (
	"testing"
	"net/http"
	"io/ioutil"
	"bytes"
	"os"
)

type tmp struct {
	Name string `json:"Name" xml:"Name"`
}

func TestResponse(t *testing.T) {
	tmpVal := &tmp{}
	resp := &Response{Response: nil}
	err := resp.ToJSON(tmpVal)
	if err != RawRespNilErr {
		t.Fatal("RawRespNilErr", err)
	}
	rawResp := &http.Response{}
	rawResp.Body = nil
	resp.Response = rawResp
	err = resp.ToXML(tmpVal)
	if err != RawRespBodyNilErr {
		t.Fatal("RawRespBodyNilErr", err)
	}

	t.Log(resp.ToString())

	rawResp.Body = ioutil.NopCloser(bytes.NewReader([]byte("{\"Name\":\"cbping\"}")))
	resp.ToJSON(tmpVal)
	if tmpVal.Name != "cbping" {
		t.Fatal("ToJSON", tmpVal.Name)
	}

	resp.ToString()
	if resp.ToString() != "{\"Name\":\"cbping\"}" {
		t.Fatal("ToString", resp.ToString())
	}

	rawResp.Body = ioutil.NopCloser(bytes.NewReader([]byte("<Name>it's my \"node\" & i like it<Name>")))
	resp.ToXML(tmpVal)
	if tmpVal.Name != "cbping" {
		t.Fatal("ToXML", tmpVal.Name)
	}

	resp.ToFile("cbping.xml")
	xmlFile, err := os.Open("cbping.xml")
	if err != nil {
		t.Fatal("ToFile", tmpVal.Name)
	}

	str := make([]byte, 1024)
	xmlFile.Read(str)
	if string(str) == "" {
		t.Fatal("ToFile", string(str))
	}
	t.Log(string(str))
	resp.Close()
}
