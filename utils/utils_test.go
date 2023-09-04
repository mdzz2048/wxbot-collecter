package utils

import (
	"fmt"
	"testing"
)

type APIResp struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Token   string      `json:"token"`
	UserId  string      `json:"userId"`
}

type LoginBody struct {
	UserName  string `json:"userName"`
	Password  string `json:"password"`
	AutoLogin bool   `json:"autoLogin"`
}

func TestStructToReaderByFormData(t *testing.T) {
	body := LoginBody{
		UserName:  "email@gmail.com",
		Password:  "password",
		AutoLogin: true,
	}
	str := StructToReaderByFormData(body)
	fmt.Println("str: ", str)
}

func TestJsonMarshal(t *testing.T) {
	type Person struct {
		Name  string
		Age   int
		Email string
	}
	person := Person{
		Name:  "Alice",
		Age:   30,
		Email: "alice@example.com",
	}
	jsonData, _ := JsonMarshal(person)
	fmt.Println(jsonData)
}

func TestJsonUnmarshal(t *testing.T) {
	jsonData := Cubox_APIResp
	jsonStruct := APIResp{}
	err := JsonUnmarshal([]byte(jsonData), &jsonStruct)
	fmt.Println(jsonStruct)

	if err != nil {
		t.Errorf("解析失败")
	}
}

func TestJsonMarshalIndent(t *testing.T) {
	jsonData := Cubox_APIResp
	jsonStruct := APIResp{}
	_ = JsonUnmarshal([]byte(jsonData), &jsonStruct)
	fmt.Println(jsonStruct)

	data, err := JsonMarshalIndent(jsonStruct, "", "  ")
	fmt.Println(string(data))

	if err != nil {
		t.Errorf("解析失败")
	}
}

func TestGetMd5Str(t *testing.T) {
	str := "password"
	hash := GetMd5Str(str)
	fmt.Print(hash)
}

func TestIsURL(t *testing.T) {
	url := "https://www.zhihu.com/answer/3195549312"
	is_url := IsURL(url)

	fmt.Println(is_url)
}
