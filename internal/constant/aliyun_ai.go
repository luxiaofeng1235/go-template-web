/*
 * @file: aliyun_ai.go
 * @description: 阿里云AI服务相关常量配置
 * @author: red
 * @created: 2025-09-10
 * @version: 1.0.0
 * @license: MIT License
 */
package constant

// 阿里云AI服务配置
const (
	// API密钥
	ALIYUN_AI_API_KEY = "sk-OTVwdAIbvI"

	// 服务端点URL
	ALIYUN_CHAT_URL  = "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"
	ALIYUN_IMAGE_URL = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text2image/image-synthesis"
	ALIYUN_VIDEO_URL = "https://dashscope.aliyuncs.com/api/v1/services/aigc/video-generation/video-synthesis"
)

// AI工作类型常量 - 与PHP版本保持一致
const (
	AiWorkTypeImage int8 = 1 // 图片生成
	AiWorkTypeVideo int8 = 2 // 视频生成
)

// AI工作状态常量
const (
	AiWorkStatusPending    int8 = 0 // 待处理
	AiWorkStatusProcessing int8 = 1 // 处理中
	AiWorkStatusCompleted  int8 = 2 // 已完成
	AiWorkStatusFailed     int8 = 3 // 失败
)

// AI模型类型常量
const (
	AI_MODEL_ASSISTANT = 0 // AI咨询助手
	AI_MODEL_SEARCH    = 1 // AI搜索
	AI_MODEL_DOCTOR    = 2 // AI家庭医生
	AI_MODEL_PET       = 3 // AI宠物助手
)

// AI模型配置
type AIModelConfig struct {
	URL    string
	Model  string
	System string
	Search bool
}

// GetAIModelConfig 获取AI模型配置
func GetAIModelConfig(modelType int, isDeepReflection bool) AIModelConfig {
	baseSystem := "，全程使用简体中文回答，如果回答中有数学相关公式请使用双$符加换行的markdown语法"

	switch modelType {
	case AI_MODEL_ASSISTANT:
		return AIModelConfig{
			URL:    ALIYUN_CHAT_URL,
			Model:  "deepseek-v3",
			System: "你现在是AiChat AI咨询助手，来自AiChat私域聊天平台，是AiChat AI中的一种，请以AiChat AI咨询助手的身份进行回答" + baseSystem,
			Search: false,
		}
	case AI_MODEL_SEARCH:
		model := "qwen-max"
		if isDeepReflection {
			model = "qwq-plus"
		}
		return AIModelConfig{
			URL:    ALIYUN_CHAT_URL,
			Model:  model,
			System: "你现在是AiChat AI搜索，来自AiChat私域聊天平台，是AiChat AI中的一种，请以AiChat AI搜索的身份在互联网上搜索相关答案并返回相关内容，如果用户输入了链接那么请拒绝用户的请求并告知用户我无法直接访问网页" + baseSystem,
			Search: true,
		}
	case AI_MODEL_DOCTOR:
		return AIModelConfig{
			URL:    ALIYUN_CHAT_URL,
			Model:  "deepseek-v3",
			System: "你现在是AiChat AI家庭医生助手，来自AiChat私域聊天平台，是AiChat AI中的一种，请以AiChat AI家庭医生的身份进行各种医学知识相关的回答，以中医为主西医为辅" + baseSystem,
			Search: false,
		}
	case AI_MODEL_PET:
		return AIModelConfig{
			URL:    ALIYUN_CHAT_URL,
			Model:  "deepseek-v3",
			System: "你现在是AiChat AI宠物助手，来自AiChat私域聊天平台，是AiChat AI中的一种，请以AiChat AI宠物助手的身份进行各种宠物知识相关回答" + baseSystem,
			Search: false,
		}
	default:
		// 默认返回助手配置
		return GetAIModelConfig(AI_MODEL_ASSISTANT, false)
	}
}

// 图片生成相关常量
const (
	// 图片生成模型
	IMAGE_MODEL_TURBO = "wanx2.1-t2i-turbo"
	IMAGE_MODEL_PLUS  = "wanx2.1-t2i-plus"

	// 图片尺寸
	IMAGE_SIZE_1024x1024 = "1024*1024"
	IMAGE_SIZE_720x1280  = "720*1280"
	IMAGE_SIZE_1280x720  = "1280*720"

	// 图片格式
	IMAGE_FORMAT_URL    = "url"
	IMAGE_FORMAT_BASE64 = "base64"
)

// 视频生成相关常量
const (
	// 视频生成模型
	VIDEO_MODEL_I2V_PLUS  = "wanx2.1-i2v-plus"  // 图生视频
	VIDEO_MODEL_T2V_TURBO = "wanx2.1-t2v-turbo" // 文生视频
)
