package tools

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

const (
	XMLHeader  = "xml"
	JSONHeader = "json"
)

const (
	POST = "POST"
	GET  = "GET"
)

type Curl interface {
	SetHeader(key, val string)
	Do() ([]byte, error)
}

type ReqParams struct {
	Url    string
	Method string
	Header string
	Params []byte
}

type reqObj struct {
	req *http.Request
}

func (p *ReqParams) InitRequest() (req Curl, err error) {
	var reqParams *bytes.Reader
	obj := new(reqObj)

	if p.Params != nil {
		reqParams = bytes.NewReader(p.Params)
		obj.req, err = http.NewRequest(p.Method, p.Url, reqParams)
	} else {
		obj.req, err = http.NewRequest(p.Method, p.Url, nil)
	}

	if err != nil {
		return nil, err
	}
	if p.Method == POST {
		switch p.Header {
		case JSONHeader:
			obj.req.Header.Set("Content-Type", "application/json;charset=UTF-8")
		case XMLHeader:
			obj.req.Header.Set("Accept", "application/xml")
			obj.req.Header.Set("Content-Status", "application/xml;charset=utf-8")
		default:
			obj.req.Header.Set("Content-Type", "application/json;charset=UTF-8")
		}
	}

	return obj, nil
}

func (obj *reqObj) SetHeader(key, val string) {
	obj.req.Header.Set(key, val)
}

func (obj *reqObj) Do() ([]byte, error) {
	defer func() {
		if er := recover(); er != nil {
			fmt.Print(fmt.Errorf("%v", er))
		}
	}()
	c := http.Client{}
	resp, err := c.Do(obj.req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
