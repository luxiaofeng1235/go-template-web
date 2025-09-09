package controller

import (
	"go-web-template/internal/service"
	"go-web-template/utils"

	"github.com/gogf/gf/v2/net/ghttp"
)

type FileController struct{}

// FormImage 表单方式上传图片
func (c *FileController) FormImage(r *ghttp.Request) {
	// 直接调用Service层处理图片上传，传入默认参数
	result, err := service.UploadImageSimple(r, 0, 0, 0)
	if err != nil {
		utils.Fail(r, err, "上传失败")
		return
	}

	utils.Success(r, result, "上传成功")
}

// FormVideo 表单方式上传视频
func (c *FileController) FormVideo(r *ghttp.Request) {
	// 直接调用Service层处理视频上传，传入默认参数
	result, err := service.UploadVideoSimple(r, 0, 0, 0)
	if err != nil {
		utils.Fail(r, err, "上传失败")
		return
	}

	utils.Success(r, result, "上传成功")
}
