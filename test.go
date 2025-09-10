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
	"go-web-template/utils"
	"reflect"
	"time"
)

func main() {

	a := utils.GetDatetime()
	fmt.Printf("a value is : %v\n", a)
	t := time.Now() // 假设 t = 2025-08-28 17:13:47 +0800 CST

	// 只保留 "YYYY-MM-DD HH:MM:SS"
	formatted := t.Format("2006-01-02 15:04:05")

	fmt.Println(reflect.TypeOf(formatted)) // 输出: 2025-08-28 17:13:47
}
