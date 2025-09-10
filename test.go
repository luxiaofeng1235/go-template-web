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
)

func main() {

	a := utils.GetDatetime()
	fmt.Printf("a value is : %v\n", a)
}
