package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
)

/* ------------------------ 判断 ------------------------ */

func IsURL(str string) bool {
	pattern := `^(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(?:/[^/]*)*$`
	match, _ := regexp.MatchString(pattern, str)
	return match
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
	resp, _ := JsonMarshalIndent(v, "", "  ")
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
