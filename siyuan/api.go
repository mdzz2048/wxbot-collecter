package siyuan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"wxbot/utils"
)

var config = getConfig(configPath)
var configPath = filepath.Join("./data", "config.json")

func getConfig(filePath string) utils.SiYuanConfig {
	data, _ := ioutil.ReadFile(filePath)

	config := utils.SiYuanConfig{}
	err := json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("解析配置文件失败", err)
		return utils.SiYuanConfig{}
	}

	return config
}

func request(option requestOption) (*SiYuanAPiResp, error) {
	// 发送请求体
	req, err := http.NewRequest(option.Method, option.URL, option.Body)
	if err != nil {
		return nil, err
	}

	if config.SiYuanToken != "" {
		req.Header.Add("Authorization", config.SiYuanToken)
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
