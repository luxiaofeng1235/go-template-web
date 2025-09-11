package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// ImageURLToBase64 将网络图片 URL 转换为 Base64 编码字符串
// 返回值：data URL 格式的字符串（如：data:image/jpeg;base64,...），以及可能的错误
func ImageURLToBase64(imageURL string) (string, error) {
	// 1. 发起 HTTP GET 请求获取图片
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("请求图片失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP 请求失败，状态码: %d", resp.StatusCode)
	}

	// 2. 读取响应体（图片二进制数据）
	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取图片数据失败: %w", err)
	}

	// 3. 推断图片的 MIME 类型（Content-Type）
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		// 如果服务器未返回 Content-Type，可以基于 URL 或文件扩展名简单推断
		switch {
		case endsWith(imageURL, ".jpg"), endsWith(imageURL, ".jpeg"):
			contentType = "image/jpeg"
		case endsWith(imageURL, ".png"):
			contentType = "image/png"
		case endsWith(imageURL, ".gif"):
			contentType = "image/gif"
		case endsWith(imageURL, ".webp"):
			contentType = "image/webp"
		default:
			contentType = "image/jpeg" // 默认
		}
	}

	// 4. 将二进制数据编码为 Base64
	base64Data := base64.StdEncoding.EncodeToString(buf.Bytes())

	// 5. 返回 data URL 格式字符串
	return fmt.Sprintf("data:%s;base64,%s", contentType, base64Data), nil
}

// 简单判断字符串是否以某后缀结尾（不区分大小写）
func endsWith(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

/*
* @note 数值切割成字符串
* @param nums int64 输入字符
* @return object
 */
func JoinInt64ToString(nums []int64) string {
	var stringSlice []string
	for _, num := range nums {
		stringSlice = append(stringSlice, strconv.FormatInt(num, 10))
	}
	return strings.Join(stringSlice, ", ")
}

/*
* @note 判断最后两位字符是否为==或者最后一位是=
* @param s string string 输入字符
* @return bool
 */
func RemoveEqualSigns(s string) string {
	if s == "" {
		return ""
	}
	n := len(s)
	if n == 0 {
		return s
	}
	// 检查最后两个字符是否为 ==
	if n > 1 && s[n-2] == '=' && s[n-1] == '=' {
		newStr := s[:n-2]
		log.Printf("查找到【%s】包含末尾有两个【=】 ,新的字符为：%s", s, newStr)
		return newStr // 去掉最后两个字符
	}
	// 检查最后一个字符是否为 =
	if s[n-1] == '=' {
		secondStr := s[:n-1]
		log.Printf("查找到【%s】包含末尾有一个【=】 ,新的字符为：%s", s, secondStr)
		return secondStr // 去掉最后一个字符
	}
	log.Printf("未发现特殊字符【=】，不需要替换，原样返回字符串 【%s】", s)
	return s // 如果都不是，返回原字符串
}
