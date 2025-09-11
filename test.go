/*
 * @file: test.go
 * @description: 阿里云AI服务功能测试
 * @author: red
 * @created: 2025-09-10
 * @version: 1.0.0
 * @license: MIT License
 */
package main

import (
	"fmt"
	"go-web-template/internal/service/common"
	"go-web-template/utils"
	"reflect"
	"time"
)

func main() {
	// 测试isLocalOrPrivateURL函数
	testURLs := []string{
		"http://192.168.0.53:8006/uploads/file/20250908/68be7a5870c62_1757313624.jpg",
	}

	fmt.Println("测试 isLocalOrPrivateURL 函数:")
	fmt.Println("========================================")

	for _, testURL := range testURLs {
		result := common.IsLocalOrPrivateURL(testURL)
		fmt.Printf("URL: %s\n", testURL)
		fmt.Printf("是否为本地/私有地址: %v\n", result)
		fmt.Println("----------------------------------------")
	}

	return

	base64Str, err := utils.ImageURLToBase64("http://192.168.0.53:8006/uploads/file/20250908/68be7a5870c62_1757313624.jpg")
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
		return
	}

	// 打印完整的 Base64 数据 URL（可用于 HTML 的 src）
	fmt.Println("Base64 Data URL:")
	fmt.Println(base64Str[:100] + "...") // 只打印前100字符避免输出过长
	return

	a := utils.GetDatetime()
	fmt.Printf("a value is : %v\n", a)
	t := time.Now() // 假设 t = 2025-08-28 17:13:47 +0800 CST

	// 只保留 "YYYY-MM-DD HH:MM:SS"
	formatted := t.Format("2006-01-02 15:04:05")

	fmt.Println(reflect.TypeOf(formatted)) // 输出: 2025-08-28 17:13:47
}
