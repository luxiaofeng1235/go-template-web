/*
 * @file: call_aliyun.go
 * @description: 阿里云AI服务调用封装 - 包含文本生成、图片生成、视频生成
 * @author: red
 * @created: 2025-09-10
 * @version: 1.0.0
 * @license: MIT License
 */
package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-web-template/global"
	"go-web-template/internal/constant"
	"go-web-template/internal/models"
	"go-web-template/utils"
	"io"
	"net/http"
	"time"
)

// 通用响应结构
type AliyunAIResponse struct {
	RequestID string      `json:"request_id"`
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Output    interface{} `json:"output"`
	Usage     interface{} `json:"usage"`
}

// 文本对话请求结构
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 图片生成请求结构
type ImageGenerateRequest struct {
	Model      string              `json:"model"`
	Input      ImageGenerateInput  `json:"input"`
	Parameters ImageGenerateParams `json:"parameters"`
}

type ImageGenerateInput struct {
	Prompt string `json:"prompt"`
}

type ImageGenerateParams struct {
	Size   string `json:"size,omitempty"`
	N      int    `json:"n,omitempty"`
	Seed   int    `json:"seed,omitempty"`
	Style  string `json:"style,omitempty"`
	Format string `json:"format,omitempty"`
}

// 图片生成响应结构
type ImageGenerateOutput struct {
	TaskID     string        `json:"task_id"`
	TaskStatus string        `json:"task_status"`
	Results    []ImageResult `json:"results,omitempty"`
	Message    string        `json:"message,omitempty"`
}

type ImageResult struct {
	URL     string `json:"url,omitempty"`
	B64JSON string `json:"b64_json,omitempty"`
}

// 视频生成请求结构
type VideoGenerateRequest struct {
	Model      string              `json:"model"`
	Input      VideoGenerateInput  `json:"input"`
	Parameters VideoGenerateParams `json:"parameters"`
}

type VideoGenerateInput struct {
	Prompt string `json:"prompt"`
	ImgURL string `json:"img_url,omitempty"`
}

type VideoGenerateParams struct {
	Duration   int    `json:"duration,omitempty"`
	Resolution string `json:"resolution,omitempty"`
	FrameRate  int    `json:"frame_rate,omitempty"`
}

// 视频生成响应结构
type VideoGenerateOutput struct {
	TaskID     string        `json:"task_id"`
	TaskStatus string        `json:"task_status"`
	Results    []VideoResult `json:"results,omitempty"`
	Message    string        `json:"message,omitempty"`
}

type VideoResult struct {
	URL string `json:"url,omitempty"`
}

// AliyunAIService 阿里云AI服务
type AliyunAIService struct {
	apiKey string
	client *http.Client
}

// NewAliyunAIService 创建阿里云AI服务实例
func NewAliyunAIService() *AliyunAIService {
	return &AliyunAIService{
		apiKey: constant.ALIYUN_AI_API_KEY,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ChatCompletion 文本对话生成
func (s *AliyunAIService) ChatCompletion(modelType int, messages []ChatMessage, isDeepReflection bool) (*AliyunAIResponse, error) {
	config := constant.GetAIModelConfig(modelType, isDeepReflection)

	// 添加系统消息
	allMessages := append([]ChatMessage{{
		Role:    "system",
		Content: config.System,
	}}, messages...)

	request := ChatRequest{
		Model:       config.Model,
		Messages:    allMessages,
		Temperature: 0.7,
		MaxTokens:   2000,
		Stream:      false,
	}

	return s.makeRequest("POST", config.URL, request)
}

// GenerateImage 生成图片 - toImage方法的Go实现
func (s *AliyunAIService) GenerateImage(prompt string, options *ImageGenerateParams) (*AliyunAIResponse, error) {
	return s.GenerateImageWithModel(prompt, constant.IMAGE_MODEL_TURBO, options)
}

// GenerateImageWithModel 使用指定模型生成图片
func (s *AliyunAIService) GenerateImageWithModel(prompt string, model string, options *ImageGenerateParams) (*AliyunAIResponse, error) {
	// 设置默认参数
	if options == nil {
		options = &ImageGenerateParams{}
	}
	if options.Size == "" {
		options.Size = constant.IMAGE_SIZE_1024x1024
	}
	if options.N == 0 {
		options.N = 1
	}
	if options.Format == "" {
		options.Format = constant.IMAGE_FORMAT_URL
	}

	request := ImageGenerateRequest{
		Model: model,
		Input: ImageGenerateInput{
			Prompt: prompt,
		},
		Parameters: *options,
	}

	return s.makeRequest("POST", constant.ALIYUN_IMAGE_URL, request)
}

// GenerateVideo 生成视频 - toVideo方法的Go实现
func (s *AliyunAIService) GenerateVideo(prompt string, options *VideoGenerateParams) (*AliyunAIResponse, error) {
	return s.GenerateVideoWithModel(prompt, constant.VIDEO_MODEL_T2V_TURBO, options)
}

// GenerateVideoWithModel 使用指定模型生成视频
func (s *AliyunAIService) GenerateVideoWithModel(prompt string, model string, options *VideoGenerateParams) (*AliyunAIResponse, error) {
	// 与PHP保持一致，简化请求结构
	request := VideoGenerateRequest{
		Model: model,
		Input: VideoGenerateInput{
			Prompt: prompt,
		},
	}

	// 只有在options不为空且有实际参数时才添加Parameters
	if options != nil && (options.Duration != 0 || options.Resolution != "" || options.FrameRate != 0) {
		request.Parameters = *options
	}

	return s.makeRequest("POST", constant.ALIYUN_VIDEO_URL, request)
}

// GetTaskStatus 获取异步任务状态
func (s *AliyunAIService) GetTaskStatus(taskURL string) (*AliyunAIResponse, error) {
	return s.makeRequest("GET", taskURL, nil)
}

// makeRequest 发起HTTP请求的通用方法
func (s *AliyunAIService) makeRequest(method, url string, payload interface{}) (*AliyunAIResponse, error) {
	var body io.Reader

	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("序列化请求数据失败: %v", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 发起请求
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	// 解析响应
	var result AliyunAIResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &result, nil
}

// 便捷方法：直接生成图片并上传到OSS
func (s *AliyunAIService) GenerateImageToOSS(prompt string, ossService *OSSService, options *ImageGenerateParams) (*OSSUploadResult, error) {
	// 生成图片
	resp, err := s.GenerateImage(prompt, options)
	if err != nil {
		return nil, fmt.Errorf("生成图片失败: %v", err)
	}

	// 解析图片生成结果
	output, ok := resp.Output.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("图片生成响应格式错误")
	}

	results, ok := output["results"].([]interface{})
	if !ok || len(results) == 0 {
		return nil, fmt.Errorf("未获取到生成的图片")
	}

	result := results[0].(map[string]interface{})
	imageURL, ok := result["url"].(string)
	if !ok {
		return nil, fmt.Errorf("未获取到图片URL")
	}

	// 上传图片到OSS
	filename := fmt.Sprintf("ai_image_%d.jpg", time.Now().Unix())
	ossResult, err := ossService.UploadFileByURL(imageURL, filename)
	if err != nil {
		return nil, fmt.Errorf("上传图片到OSS失败: %v", err)
	}

	return ossResult, nil
}

// GenerateImageByModel 根据模型生成图片的统一入口方法
func GenerateImageByModel(modelType int, prompt string, size string, n int, watermark string) (*AliyunAIResponse, error) {
	return GenerateImageByModelWithUser(modelType, prompt, size, n, watermark, "0")
}

// GenerateImageByModelWithUser 根据模型生成图片的统一入口方法（包含用户ID）
func GenerateImageByModelWithUser(modelType int, prompt string, size string, n int, watermark string, userID string) (*AliyunAIResponse, error) {
	// 创建AI服务实例
	aiService := NewAliyunAIService()
	if aiService == nil {
		return nil, fmt.Errorf("AI服务初始化失败")
	}

	// 根据模型类型选择正确的模型
	var model string
	switch modelType {
	case 1:
		model = constant.IMAGE_MODEL_TURBO
	case 2:
		model = constant.IMAGE_MODEL_PLUS
	default:
		model = constant.IMAGE_MODEL_TURBO
	}

	// 构建图片生成参数
	params := &ImageGenerateParams{
		Size:   size,
		N:      n,
		Format: constant.IMAGE_FORMAT_URL,
	}

	// 调用阿里云API生成图片
	result, err := aiService.GenerateImageWithModel(prompt, model, params)
	if err != nil {
		return nil, err
	}

	// 检查响应中是否包含task_id
	if result != nil && result.Output != nil {
		output, ok := result.Output.(map[string]interface{})
		if ok {
			if taskID, exists := output["task_id"].(string); exists && taskID != "" {
				// 保存AI工作记录到数据库
				err = SaveAIWork(&models.CreateAiWorkReq{
					UserID: userID,
					TaskID: taskID,
					Params: map[string]interface{}{
						"model":     model,
						"prompt":    prompt,
						"size":      size,
						"n":         n,
						"watermark": watermark,
					},
					Type:       models.AiWorkTypeImage, // 图片生成类型为2
					CreateTime: utils.GetUnix(),
					UpdateTime: utils.GetUnix(),
				})
				if err != nil {
					global.Errlog.Error("保存AI图片生成工作记录失败", "taskID", taskID, "error", err)
					// 不影响主流程，继续返回结果
				}
			}
		}
	}

	return result, nil
}

// GenerateVideoByType 根据类型生成视频的统一入口方法
func GenerateVideoByType(toType int, prompt string, imgURL string) (*AliyunAIResponse, error) {
	return GenerateVideoByTypeWithUser(toType, prompt, imgURL, "0")
}

// GenerateVideoByTypeWithUser 根据类型生成视频的统一入口方法（包含用户ID）
func GenerateVideoByTypeWithUser(toType int, prompt string, imgURL string, userID string) (*AliyunAIResponse, error) {
	// 创建AI服务实例
	aiService := NewAliyunAIService()
	if aiService == nil {
		return nil, fmt.Errorf("AI服务初始化失败")
	}

	// 根据to参数选择正确的模型
	var model string
	switch toType {
	case 1:
		model = constant.VIDEO_MODEL_I2V_PLUS // 图生视频
	case 2:
		model = constant.VIDEO_MODEL_T2V_TURBO // 文生视频
	default:
		model = constant.VIDEO_MODEL_T2V_TURBO
	}

	// 构建请求输入
	input := VideoGenerateInput{
		Prompt: prompt,
	}

	// 如果是图生视频且提供了图片URL，添加到输入中
	if toType == 1 && imgURL != "" {
		// 注意：这里需要根据阿里云API的实际要求来处理图片URL
		// 可能需要转换为base64或其他格式
		input.ImgURL = imgURL
	}

	// 与PHP保持一致，不传递额外参数
	request := VideoGenerateRequest{
		Model: model,
		Input: input,
	}

	// 调用阿里云API生成视频
	result, err := aiService.makeRequest("POST", constant.ALIYUN_VIDEO_URL, request)
	if err != nil {
		return nil, err
	}

	// 检查响应中是否包含task_id
	if result != nil && result.Output != nil {
		output, ok := result.Output.(map[string]interface{})
		if ok {
			if taskID, exists := output["task_id"].(string); exists && taskID != "" {
				// 保存AI工作记录到数据库
				err = SaveAIWork(&models.CreateAiWorkReq{
					UserID: userID,
					TaskID: taskID,
					Params: map[string]interface{}{
						"to":      toType,
						"prompt":  prompt,
						"img_url": imgURL,
						"model":   model,
					},
					Type:       models.AiWorkTypeVideo, // 视频生成类型为4
					CreateTime: utils.GetUnix(),
					UpdateTime: utils.GetUnix(),
				})
				if err != nil {
					global.Errlog.Error("保存AI视频生成工作记录失败", "taskID", taskID, "error", err)
					// 不影响主流程，继续返回结果
				}
			}
		}
	}

	return result, nil
}

// 便捷方法：直接生成视频并上传到OSS
func (s *AliyunAIService) GenerateVideoToOSS(prompt string, ossService *OSSService, options *VideoGenerateParams) (*OSSUploadResult, error) {
	// 生成视频
	resp, err := s.GenerateVideo(prompt, options)
	if err != nil {
		return nil, fmt.Errorf("生成视频失败: %v", err)
	}

	// 解析视频生成结果
	output, ok := resp.Output.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("视频生成响应格式错误")
	}

	results, ok := output["results"].([]interface{})
	if !ok || len(results) == 0 {
		return nil, fmt.Errorf("未获取到生成的视频")
	}

	result := results[0].(map[string]interface{})
	videoURL, ok := result["url"].(string)
	if !ok {
		return nil, fmt.Errorf("未获取到视频URL")
	}

	// 上传视频到OSS
	filename := fmt.Sprintf("ai_video_%d.mp4", time.Now().Unix())
	ossResult, err := ossService.UploadFileByURL(videoURL, filename)
	if err != nil {
		return nil, fmt.Errorf("上传视频到OSS失败: %v", err)
	}

	return ossResult, nil
}

// SaveAIWork 保存AI工作记录到数据库
func SaveAIWork(req *models.CreateAiWorkReq) error {
	if req.TaskID == "" {
		return fmt.Errorf("任务ID不能为空")
	}

	// 检查任务是否已存在
	var count int64
	err := global.DB.Model(&models.AiWork{}).Where("task_id = ?", req.TaskID).Count(&count).Error
	if err != nil {
		global.Sqllog.Error("检查AI工作任务失败", "taskID", req.TaskID, "error", err)
		return fmt.Errorf("检查AI工作任务失败")
	}
	if count > 0 {
		global.Requestlog.Info("AI工作任务已存在", "taskID", req.TaskID)
		return nil // 任务已存在，直接返回成功
	}

	// 转换参数为JSON
	paramsJSON, err := json.Marshal(req.Params)
	if err != nil {
		global.Errlog.Error("序列化参数失败", "taskID", req.TaskID, "error", err)
		return fmt.Errorf("序列化参数失败")
	}

	// 创建AI工作记录
	now := utils.GetUnix()
	aiWork := models.AiWork{
		UserID:     req.UserID,
		TaskID:     req.TaskID,
		Params:     paramsJSON,
		Type:       req.Type,
		Status:     models.AiWorkStatusPending, // 待处理状态
		CreateTime: now,
		UpdateTime: now,
	}

	err = global.DB.Create(&aiWork).Error
	if err != nil {
		global.Sqllog.Error("创建AI工作记录失败", "taskID", req.TaskID, "error", err)
		return fmt.Errorf("创建AI工作记录失败")
	}

	global.Requestlog.Info("AI工作记录创建成功", "taskID", req.TaskID, "type", req.Type, "userID", req.UserID)
	return nil
}

// GetAIWorkByTaskID 根据TaskID获取AI工作记录
func GetAIWorkByTaskID(taskID string) (*models.AiWork, error) {
	if taskID == "" {
		return nil, fmt.Errorf("任务ID不能为空")
	}

	var aiWork models.AiWork
	err := global.DB.Where("task_id = ?", taskID).First(&aiWork).Error
	if err != nil {
		global.Sqllog.Error("查询AI工作记录失败", "taskID", taskID, "error", err)
		return nil, fmt.Errorf("工作记录不存在")
	}

	return &aiWork, nil
}

// UpdateAIWorkStatus 更新AI工作记录状态
func UpdateAIWorkStatus(taskID string, status int8, work map[string]interface{}) error {
	if taskID == "" {
		return fmt.Errorf("任务ID不能为空")
	}

	updates := map[string]interface{}{
		"status":      status,
		"update_time": utils.GetUnix(),
	}

	// 如果有工作结果，则更新
	if work != nil {
		workJSON, err := json.Marshal(work)
		if err != nil {
			global.Errlog.Error("序列化工作结果失败", "taskID", taskID, "error", err)
			return fmt.Errorf("序列化工作结果失败")
		}
		updates["work"] = workJSON
	}

	err := global.DB.Model(&models.AiWork{}).Where("task_id = ?", taskID).Updates(updates).Error
	if err != nil {
		global.Sqllog.Error("更新AI工作记录失败", "taskID", taskID, "error", err)
		return fmt.Errorf("更新AI工作记录失败")
	}

	global.Requestlog.Info("AI工作记录更新成功", "taskID", taskID, "status", status)
	return nil
}
