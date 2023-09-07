package simpread

import (
	"io/ioutil"
	"log"
	"net/http"
	"wxbot/utils"
)

var config = utils.Config{}
var simpreadConfig = config.GetConfig().SimpRead
var userAgent = utils.USER_AGENT

func request(option requestOption) (*APIResp, error) {
	client := http.DefaultClient

	// 发送请求体
	req, err := http.NewRequest(option.Method, option.URL, option.Body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	if simpreadConfig.SimpReadToken != "" {
		req.Header.Set("token", simpreadConfig.SimpReadToken)
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

func AddURL(url string, data *WebInfo) bool {
	var api = "https://simpread.ksria.cn/api/service/add"

	body := AddNewBody{
		URL:   url,
		Title: data.Title,
		Desc:  data.Description,
	}
	option := requestOption{
		URL:         api,
		Body:        utils.StructToReaderByJson(body),
		Method:      "POST",
		ContentType: "application/json",
	}
	result, err := request(option)
	if err != nil {
		log.Println("API 请求异常: ", err)
		return false
	}
	resp_code := result.Code
	if resp_code != 201 {
		log.Println("API 请求失败: ", result.Message)
		return false
	}
	log.Println("API 请求成功", api)
	return true
}
