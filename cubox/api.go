package cubox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"wxbot/utils"
)

// 本地测试为了方便，直接从父目录读取，打包运行的时候不行
// var configPath = filepath.Join("..", "config.json")
var urls = CuboxURL()
var config = getConfig()
var configPath = filepath.Join("./data", "config.json")

var userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36 Edg/116.0.1938.69"

type apiEndpoints struct {
	Login                  string
	SearchEngineExport     string
	SearchEngineInbox      string
	SearchEngineToday      string
	SearchEngineMy         string
	SearchEngineNew        string
	SearchEngineUpdate     string
	SearchEngineUpdateTags string
	SearchEngineWebInfo    string
	BookmarkContent        string
	BookmarkDetail         string
	BookmarkExist          string
	TagList                string
	TagRecent              string
	TagNew                 string
	TagsDelete             string
	TagUpdate              string
	TagMerge               string
	TagMove                string
	TagSort                string
	MarkLatest             string
	MarkList               string
}

func CuboxURL() apiEndpoints {
	host, exists := os.LookupEnv("CUBOX_HOST")
	if exists {
		fmt.Println("Get CUBOX_HOST: ", host)
	} else {
		host = "https://cubox.pro"
		fmt.Println("CUBOX_HOST is undefined, use default host: https://cubox.pro")
	}

	apiEndpoints := apiEndpoints{
		Login:                  host + "/c/api/login",
		SearchEngineInbox:      host + "/c/api/v2/search_engine/inbox",
		SearchEngineMy:         host + "/c/api/v2/search_engine/my",
		SearchEngineExport:     host + "/c/api/search_engines/export/text",
		SearchEngineNew:        host + "/c/api/v2/search_engine/new",
		SearchEngineToday:      host + "/c/api/search_engine/today",
		SearchEngineUpdate:     host + "/c/api/v3/search_engine/update",
		SearchEngineUpdateTags: host + "/c/api/v2/search_engines/updateTagsForName",
		SearchEngineWebInfo:    host + "/c/api/v2/search_engine/webInfo",
		BookmarkContent:        host + "/c/api/bookmark/content",
		BookmarkDetail:         host + "/c/api/v2/bookmark/detail",
		BookmarkExist:          host + "/c/api/bookmark/exist",
		TagList:                host + "/c/api/v2/tag/list",
		TagRecent:              host + "/c/api/tag/use/recent",
		TagNew:                 host + "/c/api/v2/tag/new",
		TagsDelete:             host + "/c/api/tags/delete",
		TagUpdate:              host + "/c/api/tag/update",
		TagMerge:               host + "/c/api/tag/merge",
		TagMove:                host + "/c/api/tag/move",
		TagSort:                host + "/c/api/tag/sort",
		MarkLatest:             host + "/c/api/mark/search_engine/latest",
		MarkList:               host + "/c/api/v2/mark",
	}

	return apiEndpoints
}

func getConfig() utils.CuboxConfig {
	data, _ := ioutil.ReadFile(configPath)

	config := utils.Config{}
	err := json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("解析配置文件失败", err)
		return utils.CuboxConfig{}
	}

	return config.Cubox
}

func updateConfig(config utils.CuboxConfig) {
	// 打开配置文件
	file, err := os.OpenFile(configPath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("打开配置文件失败: ", err)
		return
	}
	defer file.Close()

	// 读取原始的 JSON 数据
	mainConfig := utils.Config{}
	err = json.NewDecoder(file).Decode(&mainConfig)
	if err != nil {
		fmt.Println("解码配置时发生错误:", err)
		return
	}

	// 将文件指针移动到文件开头
	file.Seek(0, 0)

	// 清空文件内容
	err = file.Truncate(0)
	if err != nil {
		fmt.Println("清空配置文件时发生错误: ", err)
		return
	}

	// 更新配置文件
	mainConfig.Cubox = config
	err = json.NewEncoder(file).Encode(mainConfig)
	if err != nil {
		fmt.Println("写入配置文件错误: ", err)
		return
	}

	fmt.Println("更新配置文件成功")
}

func request(option requestOption) (*APIResp, error) {
	client := http.DefaultClient

	// 发送请求体
	req, err := http.NewRequest(option.Method, option.URL, option.Body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	if config.CuboxToken != "" {
		req.Header.Set("Authorization", config.CuboxToken)
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
	var api = urls.Login
	var password = utils.GetMd5Str(config.CuboxPassword)

	body := LoginBody{
		UserName:  config.CuboxUser,
		Password:  password,
		AutoLogin: true,
	}
	option := requestOption{
		URL:         api,
		Body:        utils.StructToReaderByFormData(body),
		Method:      "POST",
		ContentType: "application/x-www-form-urlencoded",
	}
	fmt.Println(option)
	result, err := request(option)
	if err != nil {
		log.Fatalln("API 请求失败: ", err)
		return ""
	}
	log.Println("API 请求成功")

	config.CuboxToken = result.Token
	updateConfig(config)

	return result.Token
}

/* ------------------------ SearchEngine ------------------------ */

func SearchEngineWebInfo(url string) (*WebInfo, error) {
	var api = urls.SearchEngineWebInfo

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
		log.Fatalln("API 请求失败: ", err)
		return nil, err
	}
	log.Println("API 请求成功")

	webInfo := WebInfo{}
	jsonData, _ := utils.JsonMarshal(result.Data)
	err = utils.JsonUnmarshal(jsonData, &webInfo)
	if err != nil {
		log.Fatalln("返回值解析失败: ", err)
		return nil, err
	}

	return &webInfo, nil
}

func SearchEngineNew(url string, data *WebInfo) (*BookMark, error) {
	var api = urls.SearchEngineNew

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
		log.Fatalln("API 请求失败: ", err)
		return nil, err
	}
	log.Println("API 请求成功")

	bookmark := BookMark{}
	jsonData, _ := utils.JsonMarshal(result.Data)
	err = utils.JsonUnmarshal(jsonData, &bookmark)
	if err != nil {
		log.Fatalln("返回值解析失败: ", err)
		return nil, err
	}
	// test
	// resp, _ := utils.JsonMarshalIndent(&result, "", "  ")
	// fmt.Println(string(resp))

	return &bookmark, nil
}
