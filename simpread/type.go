package simpread

import "io"

type requestOption struct {
	URL           string
	Method        string
	Body          io.Reader
	Authorization string
	ContentType   string
}

type APIResp struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

type AddNewBody struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	Desc  string `json:"desc,omitempty"`
	Tags  string `json:"tags,omitempty"`
	Note  string `json:"note,omitempty"`
}

type WebInfo struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
