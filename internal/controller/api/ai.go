/*
 * @file: ai.go
 * @description: AI服务API控制器 - 实现图片生成、视频生成等AI功能
 * @author: red
 * @created: 2025-09-10
 * @version: 1.0.0
 * @license: MIT License
 */
package controller

import (
	"fmt"
	"go-web-template/global"
	"go-web-template/internal/models"
	"go-web-template/internal/service/api"
	"go-web-template/internal/service/common"
	"go-web-template/utils"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type AiController struct{}

// ToImage 生成图片接口
func (c *AiController) ToImage(r *ghttp.Request) {
	var req *models.ToImageReq
	if err := r.Parse(&req); err != nil {
		utils.FailEncrypt(r, err, "参数解析失败")
		return
	}

	// 参数验证
	if err := g.Validator().Data(req).Run(r.Context()); err != nil {
		utils.FailEncrypt(r, err, "参数验证失败")
		return
	}

	// 验证size格式：应该是 "width~height" 的形式
	sizeParts := strings.Split(req.Size, "~")
	if len(sizeParts) != 2 {
		utils.FailEncrypt(r, fmt.Errorf("size格式错误"), "size格式错误")
		return
	}

	// 检查尺寸是否为数字且在合理范围内
	var width, height int
	var err error
	if width, err = strconv.Atoi(sizeParts[0]); err != nil {
		utils.FailEncrypt(r, err, "size格式错误")
		return
	}
	if height, err = strconv.Atoi(sizeParts[1]); err != nil {
		utils.FailEncrypt(r, err, "size格式错误")
		return
	}

	// 验证尺寸范围：720-1440，总像素不超过200万
	if width < 720 || width > 1440 || height < 720 || height > 1440 {
		utils.FailEncrypt(r, fmt.Errorf("尺寸范围错误"), "尺寸范围必须在720-1440之间")
		return
	}
	if width*height > 2000000 {
		utils.FailEncrypt(r, fmt.Errorf("尺寸过大"), "尺寸过大，像素不能超过200万")
		return
	}

	// 构建图片尺寸字符串（用于阿里云API）
	imageSize := fmt.Sprintf("%d*%d", width, height)

	// 构建尺寸数组（用于数据库存储）
	sizeArray := []string{fmt.Sprintf("%d", width), fmt.Sprintf("%d", height)}

	// 调用图片生成服务 - 传递数组格式的size用于数据库存储，使用原始watermark值
	result, err := common.GenerateImageByModelWithUserAndSize(req.Model, req.Prompt, imageSize, sizeArray, req.N, req.Watermark, fmt.Sprintf("%d", req.UserID))
	if err != nil {
		global.Errlog.Error(r.Context(), "图片生成失败: %v", err)
		utils.FailEncrypt(r, err, "图片生成失败")
		return
	}

	// 提取task_id从阿里云响应中
	taskID := ""
	if result != nil && result.Output != nil {
		output, ok := result.Output.(map[string]interface{})
		if ok {
			if id, exists := output["task_id"].(string); exists {
				taskID = id
			}
		}
	}

	// 如果没有获取到task_id，返回错误
	if taskID == "" {
		utils.FailEncrypt(r, fmt.Errorf("获取任务ID失败"), "图片生成请求失败")
		return
	}

	// 返回期望格式：只返回task_id字符串
	utils.Success(r, taskID, "图片生成中")
}

// GetImage 获取图片生成结果
func (c *AiController) GetImage(r *ghttp.Request) {
	var req *models.GetTaskResultReq
	if err := r.Parse(&req); err != nil {
		utils.FailEncrypt(r, err, "参数解析失败")
		return
	}

	// 参数验证
	if err := g.Validator().Data(req).Run(r.Context()); err != nil {
		utils.FailEncrypt(r, err, "参数验证失败")
		return
	}

	// 调用Service层获取图片结果
	result, err := common.GetImageResult(req.TaskId, req.UserID)
	if err != nil {
		utils.FailEncrypt(r, err, "获取图片结果失败")
		return
	}

	utils.Success(r, result, "获取图片结果成功")
}

// ToVideo 生成视频接口
func (c *AiController) ToVideo(r *ghttp.Request) {
	var req *models.ToVideoReq
	if err := r.Parse(&req); err != nil {
		utils.FailEncrypt(r, err, "参数解析失败")
		return
	}

	// 参数验证
	if err := g.Validator().Data(req).Run(r.Context()); err != nil {
		utils.FailEncrypt(r, err, "参数验证失败")
		return
	}

	// 根据to参数进行额外验证
	if req.To == 1 {
		// 图生视频，需要img_url
		if req.ImgURL == "" {
			utils.FailEncrypt(r, fmt.Errorf("img_url必须输入"), "图生视频时img_url必须输入")
			return
		}
	} else {
		// 文生视频，需要prompt
		if req.Prompt == "" {
			utils.FailEncrypt(r, fmt.Errorf("prompt必须输入"), "文生视频时prompt必须输入")
			return
		}
	}

	// 调用视频生成服务 - 统一在service层处理复杂逻辑
	result, err := common.GenerateVideoByTypeWithUser(req.To, req.Prompt, req.ImgURL, req.UserID)
	if err != nil {
		global.Errlog.Error(r.Context(), "视频生成失败: %v", err)
		utils.FailEncrypt(r, err, "视频生成失败")
		return
	}

	// 提取task_id从阿里云响应中
	taskID := ""
	if result != nil && result.Output != nil {
		output, ok := result.Output.(map[string]interface{})
		if ok {
			if id, exists := output["task_id"].(string); exists {
				taskID = id
			}
		}
	}

	// 如果没有获取到task_id，返回错误
	if taskID == "" {
		utils.FailEncrypt(r, fmt.Errorf("获取任务ID失败"), "视频生成请求失败")
		return
	}

	// 返回期望格式：只返回task_id字符串
	utils.Success(r, taskID, "视频生成中")
}

// GetVideo 获取视频生成结果
func (c *AiController) GetVideo(r *ghttp.Request) {
	var req *models.GetTaskResultReq
	if err := r.Parse(&req); err != nil {
		utils.FailEncrypt(r, err, "参数解析失败")
		return
	}

	// 参数验证
	if err := g.Validator().Data(req).Run(r.Context()); err != nil {
		utils.FailEncrypt(r, err, "参数验证失败")
		return
	}

	// 调用Service层获取视频结果
	result, err := common.GetVideoResult(req.TaskId, req.UserID)
	if err != nil {
		// 按照PHP逻辑，特定错误消息需要特殊处理
		errorMsg := err.Error()
		if errorMsg == "视频生成中" || errorMsg == "水印正在生成中" {
			// 按照PHP逻辑返回特殊格式：code=0, show=1, msg=错误消息
			utils.SuccessWithShow(r, req.TaskId, errorMsg)
			return
		}
		utils.FailEncrypt(r, err, errorMsg)
		return
	}

	utils.Success(r, result, "获取视频结果成功")
}

// GetAiWorkList 获取AI作品历史记录
func (c *AiController) GetAiWorkList(r *ghttp.Request) {
	var req *models.AIWorkListReq
	if err := r.Parse(&req); err != nil {
		utils.FailEncrypt(r, err, "参数解析失败")
		return
	}

	// 参数验证
	if err := g.Validator().Data(req).Run(r.Context()); err != nil {
		utils.FailEncrypt(r, err, "参数验证失败")
		return
	}

	// 调用Service层获取AI作品列表
	result, err := common.GetAiWorkList(fmt.Sprintf("%d", req.UserID), int8(req.Type), req.Page)
	if err != nil {
		utils.FailEncrypt(r, err, "获取AI作品列表失败")
		return
	}

	utils.Success(r, result, "获取AI作品列表成功")
}

func (c *AiController) SendStream2(r *ghttp.Request) {
	var req *models.AIStreamReq
	if err := r.Parse(&req); err != nil {
		utils.FailEncrypt(r, err, "参数解析失败")
		return
	}

	if err := g.Validator().Data(req).Run(r.Context()); err != nil {
		utils.FailEncrypt(r, err, "参数验证失败")
		return
	}
	// 调用Stream服务处理流式响应
	streamService := &api.StreamService{}
	if err := streamService.SendStream2(r, req); err != nil {
		global.Errlog.Error(r.Context(), "流式响应失败: %v", err)
		utils.FailEncrypt(r, err, "流式响应失败")
		return
	}
}
