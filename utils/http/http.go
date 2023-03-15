package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func HttpDo(method string, url string, msg []byte) (*http.Response, error) {
	client := &http.Client{}
	body := bytes.NewBuffer(msg)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//从 HTTP 响应中反序列化数据并返回反序列化后的对象
func DeserializeFromHttpResponse[T any](response *http.Response) (*T, error) {
	// 读取响应体
	response_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	// 反序列化为目标类型对象
	var result T
	err = json.Unmarshal(response_body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
