package siyuan

import (
	"fmt"
	"testing"
	"wxbot/utils"
)

type getNotebookBody struct {
	Notebook string `json:"notebook"`
}

func TestRequest(t *testing.T) {
	body := getNotebookBody{
		Notebook: "20221203002144-vffam8j",
	}
	option := requestOption{
		// URL:         "http://127.0.0.1:6806/api/notebook/lsNotebooks",
		URL:         "http://127.0.0.1:6806/api/notebook/getNotebookConf",
		Method:      "POST",
		ContentType: "application/json",
		Body:        utils.StructToReaderByJson(body),
	}
	resp, err := request(option)
	if err != nil {
		fmt.Println(err)
		t.Errorf("请求失败")
	}
	data, _ := utils.JsonMarshalIndent(resp, "", "    ")
	fmt.Println(string(data))
}
