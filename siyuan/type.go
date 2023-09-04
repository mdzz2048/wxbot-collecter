package siyuan

import (
	"io"
)

type SiYuanAPiResp struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

type requestOption struct {
	URL           string
	Method        string
	Body          io.Reader
	Authorization string
	ContentType   string
}
