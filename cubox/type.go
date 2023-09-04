package cubox

import "io"

/* ------------------------ 请求体 ------------------------ */

type requestOption struct {
	URL           string
	Method        string
	Body          io.Reader
	Authorization string
	ContentType   string
}

// path: /c/api/login
type LoginBody struct {
	UserName  string `json:"userName"`
	Password  string `json:"password"`
	AutoLogin bool   `json:"autoLogin"`
}

// path: /c/api/v2/search_engine/webInfo
type WebInfoBody struct {
	URL string `json:"url"`
}

// path: /c/api/v2/search_engine/new
type AddNewBody struct {
	Type        int    `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TargetURL   string `json:"targetURL"`
	WebContent  string `json:"webContent"`
}

/* ------------------------ 返回值 ------------------------ */

// The basic return structure
type APIResp struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Token   string      `json:"token"`
	UserId  string      `json:"userId"`
}

// path: /c/api/v2/search_engine/webInfo
type WebInfoCover struct {
	Key string `json:"key"`
	Src string `json:"src"`
}

type WebInfo struct {
	URL         string         `json:"url"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Covers      []WebInfoCover `json:"covers"`
	Tags        []string       `json:"tags"`
}

// path: /c/api/v2/search_engine/new
type BookMark struct {
	HasMark            bool     `json:"hasMark"`
	InBlackOrWhiteList bool     `json:"inBlackOrWhiteList"`
	GroupID            string   `json:"groupId"`
	IsRead             bool     `json:"isRead"`
	Description        string   `json:"description"`
	Title              string   `json:"title"`
	Archiving          bool     `json:"archiving"`
	Type               int      `json:"type"`
	Content            string   `json:"content"`
	Cover              string   `json:"cover"`
	ResourceURL        string   `json:"resourceURL"`
	ArticleWordCount   int      `json:"articleWordCount"`
	Byline             string   `json:"byline"`
	UserSearchEngineID string   `json:"userSearchEngineID"`
	ArticleName        string   `json:"articleName"`
	ArchiveName        string   `json:"archiveName"`
	UpdateTime         string   `json:"updateTime"` // 使用 time.Time 类型来解析时间字符串
	Finished           bool     `json:"finished"`
	Marks              []string `json:"marks"` // 使用字符串切片来解析 marks 字段
	Tags               []string `json:"tags"`  // 使用字符串切片来解析 tags 字段
	HomeURL            string   `json:"homeURL"`
	GroupName          string   ` json:"groupName"`
	MarkCount          int      ` json:"markCount"`
	CreateTime         string   ` json:"createTime"`
	StarTarget         bool     ` json:"starTarget"`
	TargetURL          string   ` json:"targetURL"`
	Status             string   ` json:"status"`
}

// path: /c/api/group/my/other
type MyOtherInfo struct {
	AllSize   int `json:"allSize"`
	InboxSize int `json:"inboxSize"`
	TodaySize int `json:"todaySize"`
	StarSize  int `json:"starSize"`
}
