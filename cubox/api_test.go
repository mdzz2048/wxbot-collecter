package cubox

import (
	"fmt"
	"testing"
	"wxbot/utils"
)

func TestGetConfig(t *testing.T) {
	config := getConfig()

	if (config == utils.CuboxConfig{}) {
		t.Errorf("获取配置文件失败")
	}
	fmt.Println(config)
}

func TestUpdateConfig(t *testing.T) {
	config := getConfig()

	config.CuboxToken = "1"
	updateConfig(config)
}

func TestRequest(t *testing.T) {
	option := requestOption{
		URL:         "http://localhost:6806/api/notebook/lsNotebooks",
		ContentType: "application/json",
	}
	resp, err := request(option)
	if err != nil {
		t.Errorf("请求失败")
	}
	fmt.Println(resp)
}

func TestLogin(t *testing.T) {
	token := Login()

	if token == "" {
		t.Errorf("登录失败")
	}
	fmt.Println(token)
}

func TestSearchEngineWebInfo(t *testing.T) {
	url := "https://cloud.tencent.com/developer/article/1849807"
	webInfo, err := SearchEngineWebInfo(url)

	if err != nil {
		t.Errorf(err.Error())
	}
	data, _ := utils.JsonMarshalIndent(webInfo, "", "  ")
	fmt.Println(string(data))
}

func TestSearchEngineNew(t *testing.T) {
	url := "https://juejin.cn/post/7197053309106552888"
	webInfo, err := SearchEngineWebInfo(url)
	if err != nil {
		t.Errorf(err.Error())
	}

	new, err := SearchEngineNew(url, webInfo)
	if err != nil {
		t.Errorf(err.Error())
	}
	data, _ := utils.JsonMarshalIndent(new, "", "  ")
	fmt.Println(string(data))
}
