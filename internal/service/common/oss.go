/*
 * @file: oss.go
 * @description: 阿里云OSS对象存储服务封装 - 通用文件存储服务
 * @author: red
 * @created: 2025-09-10
 * @version: 1.0.0
 * @license: MIT License
 */
package common

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gogf/gf/v2/util/grand"
)

// OSSConfig OSS配置结构
type OSSConfig struct {
	Endpoint        string // OSS访问域名
	AccessKeyID     string // AccessKey ID
	AccessKeySecret string // AccessKey Secret
	BucketName      string // 存储桶名称
	Domain          string // 自定义域名(可选)
}

// OSSService OSS服务结构
type OSSService struct {
	client *oss.Client
	bucket *oss.Bucket
	config *OSSConfig
}

// OSSUploadResult 上传结果
type OSSUploadResult struct {
	Key      string `json:"key"`      // OSS对象键
	URL      string `json:"url"`      // 访问URL
	Size     int64  `json:"size"`     // 文件大小
	Filename string `json:"filename"` // 原始文件名
}

// NewOSSService 创建OSS服务实例
func NewOSSService(config *OSSConfig) (*OSSService, error) {
	client, err := oss.New(config.Endpoint, config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("创建OSS客户端失败: %v", err)
	}

	bucket, err := client.Bucket(config.BucketName)
	if err != nil {
		return nil, fmt.Errorf("获取OSS存储桶失败: %v", err)
	}

	return &OSSService{
		client: client,
		bucket: bucket,
		config: config,
	}, nil
}

// UploadFile 上传文件到OSS
func (s *OSSService) UploadFile(reader io.Reader, filename string, contentType string) (*OSSUploadResult, error) {
	// 生成唯一的对象键
	key := s.generateObjectKey(filename)

	// 上传选项
	options := []oss.Option{
		oss.ContentType(contentType),
		oss.ContentDisposition(fmt.Sprintf("inline; filename=\"%s\"", filename)),
	}

	// 上传文件
	err := s.bucket.PutObject(key, reader, options...)
	if err != nil {
		return nil, fmt.Errorf("上传文件到OSS失败: %v", err)
	}

	// 获取文件信息
	props, err := s.bucket.GetObjectDetailedMeta(key)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 构建访问URL
	url := s.getObjectURL(key)

	// 获取文件大小
	size := int64(0)
	if contentLength := props.Get("Content-Length"); contentLength != "" {
		fmt.Sscanf(contentLength, "%d", &size)
	}

	return &OSSUploadResult{
		Key:      key,
		URL:      url,
		Size:     size,
		Filename: filename,
	}, nil
}

// UploadFileWithPath 上传文件到OSS的指定路径（用于视频等需要固定路径的场景）
func (s *OSSService) UploadFileWithPath(reader io.Reader, objectPath string, contentType string) (*OSSUploadResult, error) {
	// 直接使用指定的路径作为对象键
	key := objectPath

	// 上传选项
	options := []oss.Option{
		oss.ContentType(contentType),
		oss.ContentDisposition(fmt.Sprintf("inline; filename=\"%s\"", filepath.Base(objectPath))),
	}

	// 上传文件
	err := s.bucket.PutObject(key, reader, options...)
	if err != nil {
		return nil, fmt.Errorf("上传文件到OSS失败: %v", err)
	}

	// 获取文件信息
	props, err := s.bucket.GetObjectDetailedMeta(key)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 构建访问URL
	url := s.getObjectURL(key)

	// 获取文件大小
	size := int64(0)
	if contentLength := props.Get("Content-Length"); contentLength != "" {
		fmt.Sscanf(contentLength, "%d", &size)
	}

	return &OSSUploadResult{
		Key:      key,
		URL:      url,
		Size:     size,
		Filename: filepath.Base(objectPath),
	}, nil
}

// UploadFileByURL 通过URL上传文件到OSS
func (s *OSSService) UploadFileByURL(fileURL string, filename string) (*OSSUploadResult, error) {
	// 使用HTTP客户端下载文件
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("从URL下载文件失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("下载文件失败，状态码: %d", resp.StatusCode)
	}

	// 获取内容类型
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 上传到OSS
	return s.UploadFile(resp.Body, filename, contentType)
}

// DeleteFile 删除OSS中的文件
func (s *OSSService) DeleteFile(key string) error {
	err := s.bucket.DeleteObject(key)
	if err != nil {
		return fmt.Errorf("删除OSS文件失败: %v", err)
	}
	return nil
}

// GetFileURL 获取文件访问URL
func (s *OSSService) GetFileURL(key string) string {
	return s.getObjectURL(key)
}

// GetSignedURL 获取签名URL(临时访问)
func (s *OSSService) GetSignedURL(key string, expires time.Duration) (string, error) {
	signedURL, err := s.bucket.SignURL(key, oss.HTTPGet, int64(expires.Seconds()))
	if err != nil {
		return "", fmt.Errorf("生成签名URL失败: %v", err)
	}
	return signedURL, nil
}

// ListFiles 列出文件
func (s *OSSService) ListFiles(prefix string, maxKeys int) ([]string, error) {
	marker := ""
	var keys []string

	for {
		lsRes, err := s.bucket.ListObjects(oss.Marker(marker), oss.Prefix(prefix), oss.MaxKeys(maxKeys))
		if err != nil {
			return nil, fmt.Errorf("列出文件失败: %v", err)
		}

		for _, object := range lsRes.Objects {
			keys = append(keys, object.Key)
		}

		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}

	return keys, nil
}

// generateObjectKey 生成唯一的对象键
func (s *OSSService) generateObjectKey(filename string) string {
	// 获取文件扩展名
	ext := filepath.Ext(filename)

	// 生成日期路径
	datePath := time.Now().Format("2006/01/02")

	// 生成随机文件名
	randomName := grand.S(16) + "_" + fmt.Sprintf("%d", time.Now().Unix())

	// 组合完整路径
	return fmt.Sprintf("uploads/%s/%s%s", datePath, randomName, ext)
}

// getObjectURL 构建对象访问URL
func (s *OSSService) getObjectURL(key string) string {
	if s.config.Domain != "" {
		// 使用自定义域名
		return fmt.Sprintf("https://%s/%s", s.config.Domain, key)
	}

	// 使用默认OSS域名
	endpoint := strings.TrimPrefix(s.config.Endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "http://")
	return fmt.Sprintf("https://%s.%s/%s", s.config.BucketName, endpoint, key)
}

// IsExist 检查文件是否存在
func (s *OSSService) IsExist(key string) (bool, error) {
	exist, err := s.bucket.IsObjectExist(key)
	if err != nil {
		return false, fmt.Errorf("检查文件存在失败: %v", err)
	}
	return exist, nil
}
