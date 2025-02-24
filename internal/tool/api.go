package tool

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Tag struct {
	Name string `json:"name"`
}

func GetLatestTag() (string, error) {
	// 设置GitHub API的URL
	url := "https://api.github.com/repos/qycnet/palworld-server-tool-main/tags"
	// 创建HTTP客户端，设置超时时间为10秒
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	// 创建HTTP GET请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	// 发送HTTP请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	// 确保在函数结束时关闭响应体
	defer resp.Body.Close()
	// 读取响应体的全部内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// 定义Tag结构体切片
	var tags []Tag
	// 将响应体内容解析为Tag结构体切片
	err = json.Unmarshal(body, &tags)
	if err != nil {
		return "", err
	}
	// 检查是否有标签
	if len(tags) == 0 {
		return "", fmt.Errorf("未找到标签")
	}
	// 返回最新的标签名称
	return tags[0].Name, nil
}

func GetLatestTagFromGitee() (string, error) {
	// 设置请求的URL
	url := "https://gitee.com/api/v5/repos/qycnet/palworld-server-tool-main/tags"
	// 发起HTTP GET请求
	resp, err := http.Get(url)
	if err != nil {
		// 如果请求失败，返回错误
		return "", err
	}
	// 延迟关闭响应体
	defer resp.Body.Close()

	// 读取响应体的内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// 如果读取失败，返回错误
		return "", err
	}

	// 定义Tag切片
	var tags []Tag
	// 将响应体内容反序列化为Tag切片
	err = json.Unmarshal(body, &tags)
	if err != nil {
		// 如果反序列化失败，返回错误
		return "", err
	}

	// 如果tags不为空
	if len(tags) > 0 {
		// 返回最新标签的名称
		return tags[len(tags)-1].Name, nil
	}

	// 如果没有找到标签，返回错误
	return "", fmt.Errorf("未找到标签")
}
