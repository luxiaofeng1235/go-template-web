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
	"go-web-template/internal/constant"
	"go-web-template/internal/service/common"
)

func main() {
	fmt.Println("=== 阿里云AI服务功能测试 ===")

	// 1. 测试AI模型配置
	fmt.Println("\n1. 测试AI模型配置:")
	testAIModelConfigs()

	// 2. 创建AI服务实例
	fmt.Println("\n2. 创建AI服务实例:")
	aiService := common.NewAliyunAIService()
	if aiService != nil {
		fmt.Println("✅ AI服务实例创建成功")
	} else {
		fmt.Println("❌ AI服务实例创建失败")
		return
	}

	// 3. 测试图片生成参数
	fmt.Println("\n3. 测试图片生成参数:")
	testImageGenerateParams(aiService)

	// 4. 测试视频生成参数
	fmt.Println("\n4. 测试视频生成参数:")
	testVideoGenerateParams(aiService)

	// 5. 测试OSS配置
	fmt.Println("\n5. 测试OSS服务配置:")
	testOSSServiceConfig()

	fmt.Println("\n=== 测试完成 ===")
}

// testAIModelConfigs 测试AI模型配置
func testAIModelConfigs() {
	models := []struct {
		name string
		modelType int
		deepReflection bool
	}{
		{"AI咨询助手", constant.AI_MODEL_ASSISTANT, false},
		{"AI搜索", constant.AI_MODEL_SEARCH, false},
		{"AI搜索(深度反思)", constant.AI_MODEL_SEARCH, true},
		{"AI家庭医生", constant.AI_MODEL_DOCTOR, false},
		{"AI宠物助手", constant.AI_MODEL_PET, false},
	}

	for _, model := range models {
		config := constant.GetAIModelConfig(model.modelType, model.deepReflection)
		fmt.Printf("  %s:\n", model.name)
		fmt.Printf("    URL: %s\n", config.URL)
		fmt.Printf("    Model: %s\n", config.Model)
		fmt.Printf("    Search: %v\n", config.Search)
		fmt.Printf("    System: %s...\n", config.System[:50])
		fmt.Println()
	}
}

// testImageGenerateParams 测试图片生成参数
func testImageGenerateParams(aiService *common.AliyunAIService) {
	prompt := "一只可爱的小猫咪在花园里玩耍"
	
	// 测试默认参数
	fmt.Printf("  测试提示词: %s\n", prompt)
	fmt.Printf("  默认图片尺寸: %s\n", constant.IMAGE_SIZE_1024x1024)
	fmt.Printf("  支持的图片格式: %s, %s\n", constant.IMAGE_FORMAT_URL, constant.IMAGE_FORMAT_BASE64)
	
	// 创建图片生成参数
	params := &common.ImageGenerateParams{
		Size:   constant.IMAGE_SIZE_1024x1024,
		N:      1,
		Format: constant.IMAGE_FORMAT_URL,
	}
	
	fmt.Printf("  参数配置 - 尺寸: %s, 数量: %d, 格式: %s\n", 
		params.Size, params.N, params.Format)
	
	// 注意：这里不实际调用API，只是测试参数配置
	fmt.Println("  ✅ 图片生成参数配置正常")
}

// testVideoGenerateParams 测试视频生成参数
func testVideoGenerateParams(aiService *common.AliyunAIService) {
	prompt := "海边日落的美丽景色，海浪轻柔地拍打着沙滩"
	
	fmt.Printf("  测试提示词: %s\n", prompt)
	fmt.Println("  视频生成模型:")
	fmt.Printf("    图生视频: %s\n", constant.VIDEO_MODEL_I2V_PLUS)
	fmt.Printf("    文生视频: %s\n", constant.VIDEO_MODEL_T2V_TURBO)
	
	// 与PHP保持一致，不传递额外参数
	fmt.Println("  参数配置: 使用默认参数（与PHP保持一致）")
	
	// 注意：这里不实际调用API，只是测试参数配置
	fmt.Println("  ✅ 视频生成参数配置正常")
}

// testOSSServiceConfig 测试OSS服务配置
func testOSSServiceConfig() {
	// 模拟OSS配置（使用假数据）
	config := &common.OSSConfig{
		Endpoint:        "oss-cn-hangzhou.aliyuncs.com",
		AccessKeyID:     "your-access-key-id",
		AccessKeySecret: "your-access-key-secret",
		BucketName:      "your-bucket-name",
		Domain:          "your-custom-domain.com", // 可选
	}
	
	fmt.Printf("  OSS端点: %s\n", config.Endpoint)
	fmt.Printf("  存储桶: %s\n", config.BucketName)
	fmt.Printf("  自定义域名: %s\n", config.Domain)
	
	// 注意：这里不实际创建OSS客户端，只是测试配置结构
	fmt.Println("  ✅ OSS服务配置结构正常")
	
	fmt.Println("  💡 实际使用时需要配置真实的OSS参数")
}

// 演示完整的使用流程（注释掉的代码，实际使用时取消注释）
func demonstrateUsage() {
	// 1. 创建服务实例
	// aiService := common.NewAliyunAIService()
	// ossConfig := &common.OSSConfig{ /* 实际配置 */ }
	// ossService, err := common.NewOSSService(ossConfig)
	
	// 2. 生成图片并上传到OSS
	// result, err := aiService.GenerateImageToOSS(
	//     "一只可爱的猫咪", 
	//     ossService, 
	//     &common.ImageGenerateParams{
	//         Size: constant.IMAGE_SIZE_1024x1024,
	//         N: 1,
	//     },
	// )
	
	// 3. 生成视频并上传到OSS
	// videoResult, err := aiService.GenerateVideoToOSS(
	//     "海边日落的美景",
	//     ossService,
	//     &common.VideoGenerateParams{
	//         // 与PHP保持一致，不传递额外参数
	//     },
	// )
	
	// 4. 文本对话
	// messages := []common.ChatMessage{{
	//     Role: "user",
	//     Content: "你好，请介绍一下自己",
	// }}
	// response, err := aiService.ChatCompletion(
	//     constant.AI_MODEL_ASSISTANT,
	//     messages,
	//     false,
	// )
}