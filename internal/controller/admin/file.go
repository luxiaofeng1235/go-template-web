/*
 * @file: file.go
 * @description: 文件上传控制器 - Admin后台管理模块
 * @author: red
 * @created: 2025-09-10
 * @version: 1.0.0
 * @license: MIT License
 */
package controller

import (
	"go-web-template/internal/service/common"
	"go-web-template/utils"

	"github.com/gogf/gf/v2/net/ghttp"
)

type FileController struct{}

// FormImage 表单方式上传图片
func (c *FileController) FormImage(r *ghttp.Request) {
	// 直接调用Service层处理图片上传，传入默认参数
	result, err := common.UploadImageSimple(r, 0, 0, 0)
	if err != nil {
		utils.FailEncrypt(r, err, "上传失败")
		return
	}

	utils.Success(r, result, "上传成功")
}