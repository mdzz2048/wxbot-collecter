package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
)

/* ------------------------ 常量 ------------------------ */
var CONFIG_PATH = filepath.Join("./data", "config.json")
var DB_PATH = filepath.Join("./data", "wxbot.db")
var USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36 Edg/116.0.1938.69"

/* ------------------------ 判断 ------------------------ */

func IsURL(str string) bool {
	pattern := `^(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(?:/[^/]*)*$`
	match, _ := regexp.MatchString(pattern, str)
	return match
}

/* ------------------------ 环境 ------------------------ */

func GetCollectType() string {
	collectType, exist := os.LookupEnv("COLLECT_TYPE")
	if exist {
		return collectType
	} else {
		return "cubox"
	}
}

/* ------------------------ 配置 ------------------------ */
func (config Config) IsConfigExist() bool {
	if _, err := os.Stat(CONFIG_PATH); os.IsNotExist(err) {
		return false
	}
	return true
}

func (config Config) InitConfig() {
	configData := Config{
		Cubox: CuboxConfig{
			CuboxUser:     "",
			CuboxPassword: "",
			CuboxToken:    "",
		},
		SiYuan: SiYuanConfig{
			SiYuanHost:  "",
			SiYuanToken: "",
		},
		SimpRead: SimpReadConfig{
			SimpReadToken: "",
		},
		WeChat: WeChatConfig{
			WeChatName: "",
		},
	}
	jsonData, err := JsonMarshalIndent(configData, "", "\t")
	if err != nil {
		log.Fatalln("初始化配置文件失败: ", err)
	}
	err = ioutil.WriteFile(CONFIG_PATH, jsonData, 0644)
	if err != nil {
		log.Fatalln("写入配置文件失败: ", err)
	}
}

func (config Config) GetConfig() Config {
	data, _ := ioutil.ReadFile(CONFIG_PATH)

	configData := Config{}
	err := json.Unmarshal(data, &configData)
	if err != nil {
		fmt.Println("解析配置文件失败", err)
		return configData
	}

	return configData
}

func (config Config) UpdateConfig(configData Config) {
	// 打开配置文件
	file, err := os.OpenFile(CONFIG_PATH, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("打开配置文件失败: ", err)
		return
	}
	defer file.Close()

	// 将文件指针移动到文件开头
	file.Seek(0, 0)

	// 清空文件内容
	err = file.Truncate(0)
	if err != nil {
		fmt.Println("清空配置文件时发生错误: ", err)
		return
	}

	// 更新配置文件
	err = json.NewEncoder(file).Encode(configData)
	if err != nil {
		fmt.Println("写入配置文件错误: ", err)
		return
	}

	fmt.Println("更新配置文件成功")
}

/* ------------------------ 加密 ------------------------ */

func GetMd5Str(str string) string {
	hash := md5.Sum([]byte(str))
	md5Str := fmt.Sprintf("%x", hash)
	return md5Str
}

/* ------------------------ JSON 解析 ------------------------ */

func JsonUnmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		fmt.Println("json unmarshal failed: ", err)
	}
	return err
}

func JsonMarshal(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		fmt.Println("json marshal failed: ", err)
	}
	return data, err
}

func JsonMarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	data, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		fmt.Println("json marshalIndent failed: ", err)
	}
	return data, err
}

func PrintResp(v interface{}) {
	resp, _ := JsonMarshalIndent(v, "", "\t")
	fmt.Println(string(resp))
}

/* ------------------------ 结构体解析 ------------------------ */

func StructToReaderByJson(structData interface{}) io.Reader {
	data, err := json.Marshal(structData)
	if nil != err {
		fmt.Println("序列化失败", err)
		return &bytes.Reader{}
	}
	return bytes.NewReader(data)
}

func StructToReaderByFormData(structData interface{}) io.Reader {
	v := reflect.ValueOf(structData)
	if v.Kind() != reflect.Struct {
		fmt.Println("data must be a struct")
		return nil
	}

	values := url.Values{}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}

		switch value.Kind() {
		case reflect.String:
			values.Add(tag, value.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			values.Add(tag, strconv.FormatInt(value.Int(), 10))
		case reflect.Bool:
			values.Add(tag, strconv.FormatBool(value.Bool()))
		default:
			fmt.Printf("unsupported field type: %s", value.Kind())
			return nil
		}
	}

	body := bytes.NewBufferString(values.Encode())
	return body
}
