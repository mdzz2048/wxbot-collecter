package cubox

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"wxbot/utils"
)

var baseHost = getCuboxHost()
var config = utils.Config{}
var cuboxConfig = config.GetConfig().Cubox
var userAgent = utils.USER_AGENT

const (
	login                  = "/c/api/login"
	userInfo               = "/c/api/userInfo"
	searchEngineExport     = "/c/api/v2/search_engine/inbox"
	searchEngineInbox      = "/c/api/v2/search_engine/my"
	searchEngineToday      = "/c/api/search_engines/export/text"
	searchEngineMy         = "/c/api/v2/search_engine/new"
	searchEngineNew        = "/c/api/search_engine/today"
	searchEngineUpdate     = "/c/api/v3/search_engine/update"
	searchEngineUpdateTags = "/c/api/v2/search_engines/updateTagsForName"
	searchEngineWebInfo    = "/c/api/v2/search_engine/webInfo"
	bookmarkContent        = "/c/api/bookmark/content"
	bookmarkDetail         = "/c/api/v2/bookmark/detail"
	bookmarkExist          = "/c/api/bookmark/exist"
	tagList                = "/c/api/v2/tag/list"
	tagRecent              = "/c/api/tag/use/recent"
	tagNew                 = "/c/api/v2/tag/new"
	tagsDelete             = "/c/api/tags/delete"
	tagUpdate              = "/c/api/tag/update"
	tagMerge               = "/c/api/tag/merge"
	tagMove                = "/c/api/tag/move"
	tagSort                = "/c/api/tag/sort"
	markLatest             = "/c/api/mark/search_engine/latest"
	markList               = "/c/api/v2/mark"
)

func getCuboxHost() string {
	host, exists := os.LookupEnv("CUBOX_HOST")
	if exists {
		fmt.Println("Get CUBOX_HOST: ", host)
	} else {
		host = "https://cubox.pro"
		fmt.Println("CUBOX_HOST is undefined, use default host: https://cubox.pro")
	}
	return host
}

func updateConfig(cuboxConfig utils.CuboxConfig) {
	mainConfig := config.GetConfig()
	mainConfig.Cubox = cuboxConfig
	config.UpdateConfig(mainConfig)
}

func request(option requestOption) (*APIResp, error) {
	client := http.DefaultClient

	// 发送请求体
	req, err := http.NewRequest(option.Method, option.URL, option.Body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	if cuboxConfig.CuboxToken != "" {
		req.Header.Set("Authorization", cuboxConfig.CuboxToken)
	}
	if option.ContentType != "" {
		req.Header.Set("Content-Type", option.ContentType)
	}

	resp, err := client.Do(req)
	log.Printf("发送 %s 请求到 %s", option.Method, option.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应内容
	var result APIResp
	respBody, _ := ioutil.ReadAll(resp.Body)
	utils.JsonUnmarshal(respBody, &result)

	return &result, nil
}

/* ------------------------ Login ------------------------ */

func Login() string {
	var api = baseHost + login
	var password = utils.GetMd5Str(cuboxConfig.CuboxPassword)

	body := LoginBody{
		UserName:  cuboxConfig.CuboxUser,
		Password:  password,
		AutoLogin: true,
	}
	option := requestOption{
		URL:         api,
		Body:        utils.StructToReaderByFormData(body),
		Method:      "POST",
		ContentType: "application/x-www-form-urlencoded",
	}
	result, err := request(option)
	if err != nil {
		log.Println("API 请求异常: ", err)
		return ""
	}
	if result.Message != "" {
		log.Println("API 请求失败: ", result.Message)
		return ""
	}

	cuboxConfig.CuboxToken = result.Token
	updateConfig(cuboxConfig)

	return result.Token
}

func RefreshToken() string {
	var token = cuboxConfig.CuboxToken
	_, err := UserInfo()
	if err == nil {
		return token
	}
	if err.Error() == "token error!" {
		token = Login()
		return token
	}
	log.Println("刷新 Token 失败: ", err.Error())
	return token
}

/* ------------------------ UserInfo ------------------------ */

func UserInfo() (*User, error) {
	var api = baseHost + userInfo

	option := requestOption{
		URL:    api,
		Method: "GET",
	}
	result, err := request(option)
	if err != nil {
		log.Println("API 请求异常: ", err)
		return nil, err
	}
	if result.Message != "" {
		log.Println("API 请求失败: ", err)
		utils.PrintResp(result)
		return nil, errors.New(result.Message)
	}

	user := User{}
	jsonData, _ := utils.JsonMarshal(result.Data)
	err = utils.JsonUnmarshal(jsonData, &user)
	if err != nil {
		log.Println("返回值解析失败: ", err)
		return nil, err
	}
	return &user, nil
}

/* ------------------------ SearchEngine ------------------------ */

func SearchEngineWebInfo(url string) (*WebInfo, error) {
	var api = baseHost + searchEngineWebInfo

	body := WebInfoBody{
		URL: url,
	}
	option := requestOption{
		URL:         api,
		Body:        utils.StructToReaderByFormData(body),
		Method:      "POST",
		ContentType: "application/x-www-form-urlencoded",
	}
	result, err := request(option)
	if err != nil {
		log.Println("API 请求异常: ", err)
		return nil, err
	}
	if result.Message != "" {
		log.Println("API 请求失败: ", result.Message)
		utils.PrintResp(result)
		return nil, errors.New(result.Message)
	}

	webInfo := WebInfo{}
	jsonData, _ := utils.JsonMarshal(result.Data)
	err = utils.JsonUnmarshal(jsonData, &webInfo)
	if err != nil {
		log.Println("返回值解析失败: ", err)
		return nil, err
	}

	return &webInfo, nil
}

func SearchEngineNew(url string, data *WebInfo) (*BookMark, error) {
	var api = baseHost + searchEngineNew

	body := AddNewBody{
		Type:        0,
		Title:       data.Title,
		Description: data.Description,
		TargetURL:   data.URL,
	}
	option := requestOption{
		URL:         api,
		Body:        utils.StructToReaderByFormData(body),
		Method:      "POST",
		ContentType: "application/x-www-form-urlencoded",
	}
	result, err := request(option)
	if err != nil {
		log.Println("API 请求异常: ", err)
		return nil, err
	}
	if result.Message != "" {
		log.Println("API 请求失败: ", result.Message)
		utils.PrintResp(result)
		return nil, errors.New(result.Message)
	}

	bookmark := BookMark{}
	jsonData, _ := utils.JsonMarshal(result.Data)
	err = utils.JsonUnmarshal(jsonData, &bookmark)
	if err != nil {
		log.Println("返回值解析失败: ", err)
		return nil, err
	}

	return &bookmark, nil
}
