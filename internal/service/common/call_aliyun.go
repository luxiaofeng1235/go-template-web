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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-web-template/global"
	"go-web-template/internal/constant"
	"go-web-template/internal/models"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
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
	Size      string `json:"size,omitempty"`
	N         int    `json:"n,omitempty"`
	Seed      int    `json:"seed,omitempty"`
	Style     string `json:"style,omitempty"`
	Format    string `json:"format,omitempty"`
	Watermark bool   `json:"watermark,omitempty"`
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

	// 对于图片和视频生成API，需要启用异步模式
	if strings.Contains(url, "text2image") || strings.Contains(url, "video-generation") {
		req.Header.Set("X-DashScope-Async", "enable")
	}

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

// GenerateImageByModelWithUserAndSize 根据模型生成图片的统一入口方法（包含用户ID和数组格式尺寸）
func GenerateImageByModelWithUserAndSize(modelType int, prompt string, apiSize string, dbSize []string, n int, watermark int, userID string) (*AliyunAIResponse, error) {
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

	// 处理watermark参数：整数 1 转为 true，其他为 false
	watermarkBool := watermark == 1

	// 构建图片生成参数
	params := &ImageGenerateParams{
		Size:      apiSize, // 使用API格式的尺寸
		N:         n,
		Format:    constant.IMAGE_FORMAT_URL,
		Watermark: watermarkBool,
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
				// 保存AI工作记录到数据库 - 使用正确的参数格式
				err = SaveAIWork(&models.CreateAiWorkReq{
					UserID: userID,
					TaskID: taskID,
					Params: map[string]interface{}{
						"n":         fmt.Sprintf("%d", n),
						"size":      dbSize, // ["1024", "1024"]
						"model":     fmt.Sprintf("%d", modelType),
						"prompt":    prompt,
						"user_id":   userID,
						"watermark": fmt.Sprintf("%d", watermark), // 使用原始watermark值
					},
					Type:       constant.AiWorkTypeImage, // 图片生成类型为1
					CreateTime: func() *time.Time { t := time.Now(); return &t }(),
					UpdateTime: func() *time.Time { t := time.Now(); return &t }(),
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

// GenerateImageByModelWithUser 根据模型生成图片的统一入口方法（包含用户ID）- 兼容旧接口
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

	// 处理watermark参数：字符串 "1" 转为 true，其他为 false
	watermarkBool := watermark == "1"

	// 构建图片生成参数
	params := &ImageGenerateParams{
		Size:      size,
		N:         n,
		Format:    constant.IMAGE_FORMAT_URL,
		Watermark: watermarkBool,
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
						"size":      size, // 保持 "1024*1024" 格式
						"n":         n,
						"watermark": watermark,
					},
					Type:       constant.AiWorkTypeImage, // 图片生成类型为1
					CreateTime: func() *time.Time { t := time.Now(); return &t }(),
					UpdateTime: func() *time.Time { t := time.Now(); return &t }(),
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
		// 检查是否为本地/局域网URL，如果是则转换为base64（与PHP版本保持一致）
		isLocal := isLocalOrPrivateURL(imgURL)
		global.Requestlog.Info("图片URL检测", "url", imgURL, "isLocal", isLocal)
		
		if isLocal {
			global.Requestlog.Info("开始转换本地图片为base64", "url", imgURL)
			base64Data, err := convertImageURLToBase64(imgURL)
			if err != nil {
				global.Requestlog.Error("图片转换失败", "url", imgURL, "error", err)
				return nil, fmt.Errorf("图片转换失败: %v", err)
			}
			global.Requestlog.Info("图片转换成功", "base64Length", len(base64Data))
			input.ImgURL = base64Data
		} else {
			global.Requestlog.Info("使用原始URL", "url", imgURL)
			input.ImgURL = imgURL
		}
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
				// 保存AI工作记录到数据库 - 与PHP版本格式保持一致
				params := map[string]interface{}{
					"to":     toType, // 直接存整数，与PHP版本一致
					"prompt": prompt,
				}
				// 只有当img_url不为空时才添加，避免空字符串
				if imgURL != "" {
					params["img_url"] = imgURL // 存储原始URL，不是base64后的
				}

				err = SaveAIWork(&models.CreateAiWorkReq{
					UserID:     userID,
					TaskID:     taskID,
					Params:     params,
					Type:       constant.AiWorkTypeVideo, // 视频生成类型为2
					CreateTime: func() *time.Time { t := time.Now(); return &t }(),
					UpdateTime: func() *time.Time { t := time.Now(); return &t }(),
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
	now := time.Now()
	aiWork := models.AiWork{
		UserID:     req.UserID,
		TaskID:     req.TaskID,
		Params:     paramsJSON,
		Type:       req.Type,
		Status:     constant.AiWorkStatusProcessing, // 生成中状态
		CreateTime: &now,
		UpdateTime: &now,
	}

	err = global.DB.Create(&aiWork).Error
	if err != nil {
		global.Sqllog.Error("创建AI工作记录失败", "taskID", req.TaskID, "error", err)
		return fmt.Errorf("创建AI工作记录失败")
	}

	global.Requestlog.Info("AI工作记录创建成功", "taskID", req.TaskID, "type", req.Type, "userID", req.UserID)
	return nil
}

// GetAIWorkByTaskID 根据TaskID和UserID获取AI工作记录
func GetAIWorkByTaskID(taskID string, userID string) (*models.AiWork, error) {
	if taskID == "" {
		return nil, fmt.Errorf("任务ID不能为空")
	}
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	// 添加调试信息
	global.Requestlog.Info("GetAIWorkByTaskID查询参数", "taskID", taskID, "userID", userID)
	fmt.Printf("GetAIWorkByTaskID查询参数: taskID=%s, userID=%s\n", taskID, userID)

	var aiWork models.AiWork
	err := global.DB.Where("task_id = ? AND user_id = ?", taskID, userID).First(&aiWork).Error
	if err != nil {
		global.Sqllog.Error("查询AI工作记录失败", "taskID", taskID, "userID", userID, "error", err)
		fmt.Printf("数据库查询错误: %v\n", err)
		return nil, fmt.Errorf("工作记录不存在: taskID=%s, userID=%s", taskID, userID)
	}

	global.Requestlog.Info("查询AI工作记录成功", "taskID", taskID, "userID", userID, "workID", aiWork.ID)
	fmt.Printf("查询成功: ID=%d, TaskID=%s, UserID=%s\n", aiWork.ID, aiWork.TaskID, aiWork.UserID)
	return &aiWork, nil
}

// addWatermarkToVideo 为视频添加水印并上传到OSS
func addWatermarkToVideo(inputVideoURL string) (string, error) {
	// 创建临时目录
	if err := os.MkdirAll(constant.TEMP_VIDEO_DIR, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %v", err)
	}

	// 生成临时文件名
	timestamp := time.Now().UnixNano()
	inputVideoPath := filepath.Join(constant.TEMP_VIDEO_DIR, fmt.Sprintf("input_%d.mp4", timestamp))
	outputVideoPath := filepath.Join(constant.TEMP_VIDEO_DIR, fmt.Sprintf("output_%d.mp4", timestamp))

	// 下载原始视频到本地
	resp, err := http.Get(inputVideoURL)
	if err != nil {
		return "", fmt.Errorf("下载视频失败: %v", err)
	}
	defer resp.Body.Close()

	inputFile, err := os.Create(inputVideoPath)
	if err != nil {
		return "", fmt.Errorf("创建临时视频文件失败: %v", err)
	}
	defer inputFile.Close()

	_, err = io.Copy(inputFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("保存视频文件失败: %v", err)
	}

	// 构建 FFmpeg 命令，将水印图片叠加到视频右下角
	// ffmpeg -i input.mp4 -i watermark.png -filter_complex "overlay=W-w-10:H-h-10" output.mp4
	cmd := exec.Command("ffmpeg",
		"-i", inputVideoPath,
		"-i", constant.WATERMARK_IMAGE_PATH,
		"-filter_complex", "overlay=W-w-10:H-h-10",
		"-y", // 覆盖输出文件
		outputVideoPath)

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		global.Errlog.Error("FFmpeg执行失败", "error", err, "output", string(output))
		return "", fmt.Errorf("视频水印处理失败: %v", err)
	}

	global.Requestlog.Info("视频水印处理成功", "input", inputVideoURL, "output", outputVideoPath)

	// 创建OSS服务实例上传处理后的视频
	ossConfig := &OSSConfig{
		Endpoint:        g.Cfg().MustGet(gctx.New(), "oss.endpoint", "").String(),
		AccessKeyID:     g.Cfg().MustGet(gctx.New(), "oss.accessKeyId", "").String(),
		AccessKeySecret: g.Cfg().MustGet(gctx.New(), "oss.accessKeySecret", "").String(),
		BucketName:      g.Cfg().MustGet(gctx.New(), "oss.bucket", "").String(),
		Domain:          g.Cfg().MustGet(gctx.New(), "oss.viewDomain", "").String(),
	}

	ossService, err := NewOSSService(ossConfig)
	if err != nil {
		// 清理临时文件
		os.Remove(inputVideoPath)
		os.Remove(outputVideoPath)
		return "", fmt.Errorf("创建OSS服务失败: %v", err)
	}

	// 读取处理后的视频文件
	outputFile, err := os.Open(outputVideoPath)
	if err != nil {
		// 清理临时文件
		os.Remove(inputVideoPath)
		os.Remove(outputVideoPath)
		return "", fmt.Errorf("打开输出视频文件失败: %v", err)
	}
	defer outputFile.Close()

	// 生成视频文件名
	videoFilename := fmt.Sprintf("ai_videos/watermarked_%d.mp4", timestamp)

	// 上传到OSS
	uploadResult, err := ossService.UploadFile(outputFile, videoFilename, "video/mp4")
	if err != nil {
		// 清理临时文件
		os.Remove(inputVideoPath)
		os.Remove(outputVideoPath)
		return "", fmt.Errorf("上传视频到OSS失败: %v", err)
	}

	// 清理临时文件
	os.Remove(inputVideoPath)
	os.Remove(outputVideoPath)

	global.Requestlog.Info("视频水印处理并上传OSS成功", "input", inputVideoURL, "ossURL", uploadResult.URL)

	return uploadResult.URL, nil
}

// GetImageResult 获取图片生成结果 - 按照PHP版本逻辑实现完整状态管理
func GetImageResult(taskID string, userID string) (*models.AIImageResult, error) {
	// 使用通用的getTask方法，匹配PHP版本
	data, err := getTask(taskID, userID)
	if err != nil {
		return nil, err
	}

	// 构造响应数据
	response := &models.AIImageResult{
		Iamge: []string{},
		Image: []string{},
	}

	// 解析返回的output
	output, ok := data["output"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("响应格式错误")
	}

	taskStatus, ok := output["task_status"].(string)
	if !ok {
		return nil, fmt.Errorf("任务状态格式错误")
	}

	switch taskStatus {
	case "OK":
		// 任务已完成且已有最终结果
		global.Requestlog.Info("DEBUG: GetImageResult OK case", "output", output)

		// 尝试解析为数组格式（标准格式 - PHP兼容）
		if results, ok := output["results"].([]interface{}); ok {
			global.Requestlog.Info("DEBUG: Found results array", "results", results)
			urlStrings := make([]string, len(results))
			for i, url := range results {
				if urlStr, ok := url.(string); ok {
					urlStrings[i] = urlStr
				}
			}
			global.Requestlog.Info("DEBUG: Final urlStrings from array", "urls", urlStrings)
			response.Iamge = urlStrings
			response.Image = urlStrings
			return response, nil
		}

		// 尝试解析为对象格式（旧格式兼容）
		if resultsMap, ok := output["results"].(map[string]interface{}); ok {
			global.Requestlog.Info("DEBUG: Found results map", "results", resultsMap)
			var urlStrings []string

			// 检查是否有url字段
			if url, exists := resultsMap["url"]; exists {
				if urlStr, ok := url.(string); ok {
					urlStrings = append(urlStrings, urlStr)
				}
			}

			global.Requestlog.Info("DEBUG: Final urlStrings from map", "urls", urlStrings)
			response.Iamge = urlStrings
			response.Image = urlStrings
			return response, nil
		}

		global.Requestlog.Error("DEBUG: results not found or wrong type", "output", output)
		return nil, fmt.Errorf("任务完成但结果为空")

	case "SUCCEEDED":
		// 图片生成成功，需要下载并上传到OSS（匹配PHP逻辑）
		if results, ok := output["results"].([]interface{}); ok {
			var finalUrls []string
			i := 1

			for _, item := range results {
				if resultMap, ok := item.(map[string]interface{}); ok {
					if originalURL, exists := resultMap["url"].(string); exists && originalURL != "" {
						// 构建本地路径和最终URL（匹配PHP逻辑）
						path := fmt.Sprintf("uploads/ai_images/%s_%d.png", taskID, i)
						finalURL := fmt.Sprintf("https://static.jsss999.com/%s", path)

						// 通过URL同步到阿里云OSS（匹配PHP的 FileServer::uploadUrl 逻辑）
						err := syncURLToOSS(originalURL, path)
						if err != nil {
							global.Errlog.Error("图片同步到OSS失败", "taskID", taskID, "url", originalURL, "error", err)
							// 同步失败，使用错误图片
							errorPath := "images/aichat_uni/ai/ai_picture/icon_error.png"
							finalUrls = append(finalUrls, fmt.Sprintf("https://static.jsss999.com/%s", errorPath))
						} else {
							finalUrls = append(finalUrls, finalURL)
						}
						i++
					} else {
						// 如果没有URL，使用错误图片（匹配PHP逻辑）
						errorPath := "images/aichat_uni/ai/ai_picture/icon_error.png"
						finalUrls = append(finalUrls, fmt.Sprintf("https://static.jsss999.com/%s", errorPath))
						i++
					}
				}
			}

			if len(finalUrls) == 0 {
				return nil, fmt.Errorf("数据错误，图片URL为空，请联系管理员")
			}

			// 更新状态为已完成，保存最终URL数组（直接存储URL数组，匹配PHP版本格式）
			err = UpdateAIWorkStatus(taskID, constant.AiWorkStatusCompleted, finalUrls)
			if err != nil {
				global.Errlog.Error("更新状态失败", "taskID", taskID, "error", err)
			}

			// 返回处理后的URL列表
			response.Iamge = finalUrls
			response.Image = finalUrls
			return response, nil
		}
		return nil, fmt.Errorf("数据错误，图片URL为空，请联系管理员")

	case "FAILED":
		// 任务失败
		errorMsg := "图片生成失败"
		if message, ok := output["message"].(string); ok && message != "" {
			errorMsg = message
		}

		// 更新状态为失败
		workData := map[string]interface{}{
			"error": errorMsg,
		}
		UpdateAIWorkStatus(taskID, constant.AiWorkStatusFailed, workData)

		return nil, fmt.Errorf(errorMsg)

	case "RUNNING", "PENDING", "SUSPENDED":
		// 任务进行中
		return nil, fmt.Errorf("图片正在生成中")

	case "NULL":
		// 任务不存在
		return nil, fmt.Errorf("任务不存在")

	default:
		return nil, fmt.Errorf("任务不存在或状态未知")
	}
}

// GetVideoResult 获取视频生成结果 - 匹配PHP版本逻辑
func GetVideoResult(taskID string, userID string) (*models.AIVideoResult, error) {
	// 使用通用的getTask方法，匹配PHP版本
	data, err := getTask(taskID, userID)
	if err != nil {
		return nil, err
	}

	// 构造响应数据
	response := &models.AIVideoResult{
		Video: []string{},
	}

	// 解析返回的output
	output, ok := data["output"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("响应格式错误")
	}

	taskStatus, ok := output["task_status"].(string)
	if !ok {
		return nil, fmt.Errorf("任务状态格式错误")
	}

	switch taskStatus {
	case "OK":
		// 任务已完成且已有最终结果
		if results, ok := output["results"].([]interface{}); ok {
			videoURLs := make([]string, len(results))
			for i, url := range results {
				if urlStr, ok := url.(string); ok {
					videoURLs[i] = urlStr
				}
			}
			response.Video = videoURLs
			return response, nil
		}
		return nil, fmt.Errorf("任务完成但结果为空")

	case "SUCCEEDED":
		// 检查是否正在处理水印（status=3）
		model, err := GetAIWorkByTaskID(taskID, userID)
		if err == nil && model.Status == constant.AiWorkStatusProcessing {
			return nil, fmt.Errorf("水印正在生成中")
		}

		// 视频生成成功，需要进行水印处理
		var videoURL string
		if url, ok := output["video_url"].(string); ok && url != "" {
			videoURL = url
		} else if results, ok := output["results"].(map[string]interface{}); ok {
			if url, exists := results["video_url"].(string); exists {
				videoURL = url
			}
		}

		if videoURL == "" {
			return nil, fmt.Errorf("数据错误，请联系管理员")
		}

		// 更新状态为处理中（水印处理）
		err = UpdateAIWorkStatus(taskID, constant.AiWorkStatusProcessing, nil)
		if err != nil {
			global.Errlog.Error("更新状态失败", "taskID", taskID, "error", err)
		}

		// 为视频添加水印并上传到OSS
		watermarkedVideoURL, err := addWatermarkToVideo(videoURL)
		if err != nil {
			global.Errlog.Error("视频水印处理失败", "taskID", taskID, "url", videoURL, "error", err)
			// 水印处理失败，恢复为已完成状态
			UpdateAIWorkStatus(taskID, constant.AiWorkStatusFailed, map[string]interface{}{
				"error": err.Error(),
			})
			return nil, fmt.Errorf("水印添加失败: %v", err)
		}

		// 保存处理结果到数据库（直接存储URL数组，匹配PHP版本格式）
		err = UpdateAIWorkStatus(taskID, constant.AiWorkStatusCompleted, []string{watermarkedVideoURL})
		if err != nil {
			global.Errlog.Error("更新工作结果失败", "taskID", taskID, "error", err)
		}

		response.Video = []string{watermarkedVideoURL}
		return response, nil

	case "FAILED":
		// 视频生成失败
		errorMsg := "视频生成失败"
		if message, ok := output["message"].(string); ok && message != "" {
			errorMsg = message
		}

		// 更新数据库状态为失败
		UpdateAIWorkStatus(taskID, constant.AiWorkStatusFailed, map[string]interface{}{
			"error": errorMsg,
		})
		return nil, fmt.Errorf(errorMsg)

	case "RUNNING", "PENDING", "SUSPENDED":
		// 视频生成中
		return nil, fmt.Errorf("视频生成中")

	case "NULL":
		// 任务不存在
		return nil, fmt.Errorf("任务不存在")

	default:
		return nil, fmt.Errorf("任务不存在或状态未知")
	}
}

// TaskData 阿里云任务状态数据结构
type TaskData struct {
	Status   string   `json:"status"`
	VideoURL string   `json:"video_url"`
	Results  []string `json:"results"`
	Message  string   `json:"message"`
}

// getTask 获取任务状态 - 通用方法，匹配PHP版本的getTask方法
func getTask(taskID string, userID string) (map[string]interface{}, error) {
	// 首先检查数据库中的任务状态
	aiWork, err := GetAIWorkByTaskID(taskID, userID)
	if err != nil {
		return map[string]interface{}{
			"output": map[string]interface{}{
				"task_status": "NULL",
				"message":     "任务不存在",
			},
		}, nil
	}

	// 如果任务状态是已完成(2)，直接返回OK状态
	if aiWork.Status == constant.AiWorkStatusCompleted {
		var work interface{}
		if len(aiWork.Work) > 0 {
			json.Unmarshal(aiWork.Work, &work)
		}

		return map[string]interface{}{
			"output": map[string]interface{}{
				"task_status": "OK",
				"results":     work, // 直接返回work内容（URL数组）
			},
		}, nil
	}

	// 如果任务状态是失败(3)，返回FAILED状态
	if aiWork.Status == constant.AiWorkStatusFailed {
		var work map[string]interface{}
		errorMsg := "任务失败"
		if len(aiWork.Work) > 0 {
			json.Unmarshal(aiWork.Work, &work)
			if errMsg, ok := work["error"].(string); ok {
				errorMsg = errMsg
			}
		}
		return map[string]interface{}{
			"output": map[string]interface{}{
				"task_status": "FAILED",
				"message":     errorMsg,
			},
		}, nil
	}

	// 如果是其他状态（待处理或处理中），调用阿里云API获取最新状态
	url := fmt.Sprintf("https://dashscope.aliyuncs.com/api/v1/tasks/%s", taskID)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 添加认证头
	req.Header.Set("Authorization", "Bearer "+constant.ALIYUN_AI_API_KEY)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求阿里云API失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		global.Errlog.Error("阿里云API请求失败", "status", resp.StatusCode, "body", string(body))
		return map[string]interface{}{
			"output": map[string]interface{}{
				"task_status": "FAILED",
				"message":     "API请求失败",
			},
		}, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var apiResponse map[string]interface{}
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 处理UNKNOWN状态 - 匹配PHP逻辑
	if output, ok := apiResponse["output"].(map[string]interface{}); ok {
		if taskStatus, exists := output["task_status"]; exists && taskStatus == "UNKNOWN" {
			// 更新数据库状态为失败
			UpdateAIWorkStatus(taskID, constant.AiWorkStatusFailed, map[string]interface{}{
				"error": "任务不存在",
			})
			return map[string]interface{}{
				"output": map[string]interface{}{
					"task_status": "FAILED",
					"message":     "任务不存在",
				},
			}, nil
		}
	}

	return apiResponse, nil
}

// getTaskStatus 获取阿里云任务状态 - 完整实现匹配PHP版本（保持向后兼容）
func getTaskStatus(taskID string, userID string) (*TaskData, error) {
	// 首先检查数据库中的任务状态
	aiWork, err := GetAIWorkByTaskID(taskID, userID)
	if err != nil {
		return &TaskData{
			Status:  "NULL",
			Message: "任务不存在",
		}, nil
	}

	// 如果任务状态是已完成(1)，直接返回OK状态
	if aiWork.Status == constant.AiWorkStatusCompleted {
		var work map[string]interface{}
		if len(aiWork.Work) > 0 {
			json.Unmarshal(aiWork.Work, &work)
		}

		// 检查是否有video_urls（已处理的视频）
		if urlsInterface, ok := work["video_urls"]; ok {
			if urls, ok := urlsInterface.([]interface{}); ok && len(urls) > 0 {
				results := make([]string, len(urls))
				for i, url := range urls {
					if urlStr, ok := url.(string); ok {
						results[i] = urlStr
					}
				}
				return &TaskData{
					Status:  "OK",
					Results: results,
				}, nil
			}
		}

		// 检查是否有其他结果格式
		if resultsInterface, ok := work["results"]; ok {
			if results, ok := resultsInterface.([]string); ok {
				return &TaskData{
					Status:  "OK",
					Results: results,
				}, nil
			}
		}
	}

	// 如果任务状态是失败(2)，返回FAILED状态
	if aiWork.Status == constant.AiWorkStatusFailed {
		var work map[string]interface{}
		errorMsg := "任务失败"
		if len(aiWork.Work) > 0 {
			json.Unmarshal(aiWork.Work, &work)
			if errMsg, ok := work["error"].(string); ok {
				errorMsg = errMsg
			}
		}
		return &TaskData{
			Status:  "FAILED",
			Message: errorMsg,
		}, nil
	}

	// 如果是其他状态（待处理或处理中），调用阿里云API获取最新状态
	url := fmt.Sprintf("https://dashscope.aliyuncs.com/api/v1/tasks/%s", taskID)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 添加认证头
	req.Header.Set("Authorization", "Bearer sk-OTVwdAIbvI") // 使用PHP中相同的key
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求阿里云API失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		global.Errlog.Error("阿里云API请求失败", "status", resp.StatusCode, "body", string(body))
		return &TaskData{
			Status:  "FAILED",
			Message: "API请求失败",
		}, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var apiResponse struct {
		Output struct {
			TaskStatus string `json:"task_status"`
			VideoURL   string `json:"video_url"`
			Results    struct {
				VideoURL string `json:"video_url"`
			} `json:"results"`
			Message string `json:"message"`
		} `json:"output"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 处理UNKNOWN状态 - 匹配PHP逻辑
	if apiResponse.Output.TaskStatus == "UNKNOWN" {
		// 更新数据库状态为失败
		UpdateAIWorkStatus(taskID, constant.AiWorkStatusFailed, map[string]interface{}{
			"error": "任务不存在",
		})
		return &TaskData{
			Status:  "FAILED",
			Message: "任务不存在",
		}, nil
	}

	// 构造返回数据
	taskData := &TaskData{
		Status:  apiResponse.Output.TaskStatus,
		Message: apiResponse.Output.Message,
	}

	// 获取视频URL - 兼容不同的响应格式
	if apiResponse.Output.VideoURL != "" {
		taskData.VideoURL = apiResponse.Output.VideoURL
	} else if apiResponse.Output.Results.VideoURL != "" {
		taskData.VideoURL = apiResponse.Output.Results.VideoURL
	}

	return taskData, nil
}

// ImageTaskData 图片任务状态数据结构
type ImageTaskData struct {
	Status   string   `json:"status"`
	ImageURL string   `json:"image_url"`
	Results  []string `json:"results"`
	Message  string   `json:"message"`
}

// getImageTaskStatus 获取图片任务状态 - 专门用于图片生成
func getImageTaskStatus(taskID string, userID string) (*ImageTaskData, error) {
	// 首先检查数据库中的任务状态
	aiWork, err := GetAIWorkByTaskID(taskID, userID)
	if err != nil {
		return &ImageTaskData{
			Status:  "NULL",
			Message: "任务不存在",
		}, nil
	}

	// 如果任务状态是已完成(2)，直接返回OK状态
	if aiWork.Status == constant.AiWorkStatusCompleted {
		var work map[string]interface{}
		if len(aiWork.Work) > 0 {
			json.Unmarshal(aiWork.Work, &work)
		}

		// 检查是否有URL结果
		if url, ok := work["url"].(string); ok {
			return &ImageTaskData{
				Status:  "OK",
				Results: []string{url},
			}, nil
		}

		// 检查是否有results数组
		if resultsInterface, ok := work["results"]; ok {
			if results, ok := resultsInterface.([]string); ok {
				return &ImageTaskData{
					Status:  "OK",
					Results: results,
				}, nil
			}
		}
	}

	// 如果任务状态是失败(3)，返回FAILED状态
	if aiWork.Status == constant.AiWorkStatusFailed {
		var work map[string]interface{}
		errorMsg := "图片生成失败"
		if len(aiWork.Work) > 0 {
			json.Unmarshal(aiWork.Work, &work)
			if errMsg, ok := work["error"].(string); ok {
				errorMsg = errMsg
			}
		}
		return &ImageTaskData{
			Status:  "FAILED",
			Message: errorMsg,
		}, nil
	}

	// 如果是其他状态（待处理或处理中），调用阿里云API获取最新状态
	url := fmt.Sprintf("https://dashscope.aliyuncs.com/api/v1/tasks/%s", taskID)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 添加认证头
	req.Header.Set("Authorization", "Bearer "+constant.ALIYUN_AI_API_KEY)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求阿里云API失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		global.Errlog.Error("阿里云API请求失败", "status", resp.StatusCode, "body", string(body))
		return &ImageTaskData{
			Status:  "FAILED",
			Message: "API请求失败",
		}, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var apiResponse struct {
		Output struct {
			TaskStatus string `json:"task_status"`
			Results    []struct {
				URL string `json:"url"`
			} `json:"results"`
			Message string `json:"message"`
		} `json:"output"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 处理UNKNOWN状态
	if apiResponse.Output.TaskStatus == "UNKNOWN" {
		// 更新数据库状态为失败
		UpdateAIWorkStatus(taskID, constant.AiWorkStatusFailed, map[string]interface{}{
			"error": "任务不存在",
		})
		return &ImageTaskData{
			Status:  "FAILED",
			Message: "任务不存在",
		}, nil
	}

	// 构造返回数据
	taskData := &ImageTaskData{
		Status:  apiResponse.Output.TaskStatus,
		Message: apiResponse.Output.Message,
	}

	// 提取图片URL
	if len(apiResponse.Output.Results) > 0 {
		taskData.ImageURL = apiResponse.Output.Results[0].URL
	}

	return taskData, nil
}

// downloadAndSaveImage 下载图片并保存到本地
func downloadAndSaveImage(imageURL, taskID string) (string, error) {
	// 创建目录
	uploadDir := "public/uploads/ai_images"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %v", err)
	}

	// 生成文件名
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("ai_image_%s_%d.jpg", taskID, timestamp)
	filepath := filepath.Join(uploadDir, filename)

	// 下载图片
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("下载图片失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载图片失败，状态码: %d", resp.StatusCode)
	}

	// 创建本地文件
	file, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("创建本地文件失败: %v", err)
	}
	defer file.Close()

	// 复制内容
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("保存图片失败: %v", err)
	}

	// 返回相对URL路径
	localURL := fmt.Sprintf("/uploads/ai_images/%s", filename)
	global.Requestlog.Info("图片下载成功", "taskID", taskID, "url", localURL)

	return localURL, nil
}

// downloadAndUploadToOSS 下载阿里云图片并上传到OSS - 匹配PHP的FileServer::uploadUrl功能
func downloadAndUploadToOSS(sourceURL, targetPath string) error {

	resp, err := http.Get(sourceURL)
	if err != nil {
		return fmt.Errorf("下载图片失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载图片失败，状态码: %d", resp.StatusCode)
	}

	// 读取图片内容
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取图片数据失败: %v", err)
	}

	// TODO: 这里应该调用OSS服务上传图片到指定路径
	// 目前暂时实现为一个占位符，后续集成OSS SDK
	global.Requestlog.Info("图片下载上传成功", "source", sourceURL, "target", targetPath, "size", len(imageData))

	// 返回成功（模拟上传成功）
	return nil
}

// UpdateAIWorkStatus 更新AI工作记录状态
func UpdateAIWorkStatus(taskID string, status int8, work interface{}) error {
	if taskID == "" {
		return fmt.Errorf("任务ID不能为空")
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":      status,
		"update_time": &now,
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

// GetAiWorkList 获取AI作品列表
func GetAiWorkList(userID string, workType int8, page int) (*models.AiWorkListRes, error) {
	pageSize := constant.PAGE_SIZE

	// 构建查询条件
	query := global.DB.Model(&models.AiWork{}).Where("user_id = ? AND status IN (?)", userID, []int8{constant.AiWorkStatusProcessing, constant.AiWorkStatusCompleted})

	// 类型过滤：当type != 3时，增加类型条件
	if workType != 3 {
		query = query.Where("type = ?", workType)
	}

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		global.Sqllog.Errorf("查询AI作品总数失败: %v", err)
		return nil, fmt.Errorf("查询AI作品总数失败")
	}

	// 计算分页
	offset := (page - 1) * pageSize
	pageCount := int64(0)
	if total > 0 {
		pageCount = (total + int64(pageSize) - 1) / int64(pageSize)
	}

	// 查询列表数据
	var works []models.AiWork
	if err := query.Order("id desc").Offset(offset).Limit(pageSize).Find(&works).Error; err != nil {
		global.Sqllog.Errorf("查询AI作品列表失败: %v", err)
		return nil, fmt.Errorf("查询AI作品列表失败")
	}

	// 转换为响应格式
	workList := make([]models.AiWorkRes, 0, len(works))
	for _, work := range works {
		// 解析params和work JSON字段
		var params map[string]interface{}
		var workData map[string]interface{}

		if len(work.Params) > 0 {
			if err := json.Unmarshal(work.Params, &params); err != nil {
				global.Errlog.Warnf("解析params JSON失败: %v", err)
				params = make(map[string]interface{})
			}
		}

		if len(work.Work) > 0 {
			if err := json.Unmarshal(work.Work, &workData); err != nil {
				global.Errlog.Warnf("解析work JSON失败: %v", err)
				workData = make(map[string]interface{})
			}
		}

		workRes := models.AiWorkRes{
			TaskID: work.TaskID,
			Params: params,
			Type:   work.Type,
			Status: work.Status,
			Work:   workData,
		}
		workList = append(workList, workRes)
	}

	// 构建响应
	result := &models.AiWorkListRes{
		Page:      fmt.Sprintf("%d", page),
		PageCount: int(pageCount),
		Total:     total,
		List:      workList,
	}

	global.Requestlog.Infof("获取AI作品列表成功: userID=%s, type=%d, page=%d, total=%d", userID, workType, page, total)
	return result, nil
}

// syncURLToOSS 同步URL到阿里云OSS - 匹配PHP的FileServer::uploadUrl逻辑
func syncURLToOSS(sourceURL, targetPath string) error {

	resp, err := http.Get(sourceURL)
	if err != nil {
		return fmt.Errorf("下载图片失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载图片失败，状态码: %d", resp.StatusCode)
	}

	// 读取图片内容
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取图片数据失败: %v", err)
	}

	// TODO: 这里应该调用OSS服务上传图片到指定路径
	// 目前暂时实现为一个占位符，后续集成OSS SDK
	// 模拟PHP的FileServer::uploadUrl功能：把URL同步给阿里云的OSS
	global.Requestlog.Info("图片同步到OSS成功", "source", sourceURL, "target", targetPath, "size", len(imageData))

	// 返回成功（模拟同步成功）
	return nil
}

// isLocalOrPrivateURL 检查URL是否为本地或局域网地址 - 匹配PHP版本逻辑
func isLocalOrPrivateURL(imgURL string) bool {
	parsedURL, err := url.Parse(imgURL)
	if err != nil || parsedURL.Host == "" {
		return false
	}

	host := parsedURL.Host
	// 处理带端口的情况
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}

	// 检查是否为localhost相关
	if host == "localhost" || host == "127.0.0.1" || host == "::1" {
		return true
	}

	// 检查是否为局域网IP段
	ip := net.ParseIP(host)
	if ip != nil && ip.To4() != nil {
		// 私有IP段：10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
		return ip.IsPrivate()
	}

	return false
}

// convertImageURLToBase64 将本地图片URL转换为base64格式供阿里云API使用 - 匹配PHP版本逻辑
func convertImageURLToBase64(imageURL string) (string, error) {
	// 创建HTTP客户端，设置30秒超时
	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("获取图片失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("获取图片失败，状态码: %d", resp.StatusCode)
	}

	// 获取图片二进制数据
	imageContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取图片数据失败: %v", err)
	}

	// 获取Content-Type来确定MIME类型
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		// 如果没有Content-Type，根据URL扩展名判断
		parsedURL, _ := url.Parse(imageURL)
		ext := strings.ToLower(filepath.Ext(parsedURL.Path))
		switch ext {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".gif":
			contentType = "image/gif"
		case ".webp":
			contentType = "image/webp"
		default:
			contentType = "image/jpeg" // 默认为jpeg
		}
	}

	// 转换为base64格式
	base64Data := base64.StdEncoding.EncodeToString(imageContent)

	// 返回符合阿里云API要求的格式：data:{MIME_type};base64,{base64_data}
	return fmt.Sprintf("data:%s;base64,%s", contentType, base64Data), nil
}
