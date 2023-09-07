package cubox

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"wxbot/utils"
)

// 本地测试为了方便，直接从父目录读取，打包运行的时候不行
// var configPath = filepath.Join("..", "config.json")
var urls = CuboxURL()
var config = utils.Config{}
var cuboxConfig = config.GetConfig().Cubox
var userAgent = utils.USER_AGENT

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
	var api = urls.Login
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
	fmt.Println(option)
	result, err := request(option)
	if err != nil {
		log.Fatalln("API 请求失败: ", err)
		return ""
	}
	log.Println("API 请求成功")

	cuboxConfig.CuboxToken = result.Token
	updateConfig(cuboxConfig)

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
