package siyuan

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"wxbot/utils"
)

var config = utils.Config{}
var siyuanConfig = config.GetConfig().SiYuan

func request(option requestOption) (*SiYuanAPiResp, error) {
	// 发送请求体
	req, err := http.NewRequest(option.Method, option.URL, option.Body)
	if err != nil {
		return nil, err
	}

	if siyuanConfig.SiYuanToken != "" {
		req.Header.Add("Authorization", siyuanConfig.SiYuanToken)
	}
	if option.ContentType != "" {
		req.Header.Add("Content-Type", option.ContentType)
	}

	resp, err := http.DefaultClient.Do(req)
	log.Printf("发送 %s 请求到 %s", option.Method, option.URL)
	if err != nil {
		fmt.Println("请求失败: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应内容
	respStruct := SiYuanAPiResp{}
	respBody, _ := ioutil.ReadAll(resp.Body)
	utils.JsonUnmarshal(respBody, &respStruct)

	return &respStruct, nil
}
