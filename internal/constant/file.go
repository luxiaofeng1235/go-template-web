/*
 * @file: file.go
 * @description: 文件相关常量定义
 * @author: go-web-template
 * @created: 2025-09-09
 * @version: 1.0.0
 * @license: MIT License
 */

package constant

// 文件类型常量
const (
	FILE_TYPE_IMAGE    uint8 = 10 // 图片
	FILE_TYPE_DOCUMENT uint8 = 20 // 文档
	FILE_TYPE_VIDEO    uint8 = 30 // 视频
	FILE_TYPE_AUDIO    uint8 = 40 // 音频
	FILE_TYPE_ARCHIVE  uint8 = 50 // 压缩包
	FILE_TYPE_OTHER    uint8 = 60 // 其他
)

// 文件分类类型常量
const (
	FileCateTypeImage    uint8 = 10 // 图片分类
	FileCateTypeDocument uint8 = 20 // 文档分类
	FileCateTypeVideo    uint8 = 30 // 视频分类
	FileCateTypeAudio    uint8 = 40 // 音频分类
	FileCateTypeArchive  uint8 = 50 // 压缩包分类
	FileCateTypeOther    uint8 = 60 // 其他分类
)

// 删除状态常量
const (
	FILE_DEL_NORMAL  uint8 = 0 // 正常
	FILE_DEL_DELETED uint8 = 1 // 删除
)

// 文件分类删除状态常量
const (
	FileCateDelNormal  uint8 = 0 // 正常
	FileCateDelDeleted uint8 = 1 // 删除
)

// 向后兼容的别名
const (
	IMAGE_TYPE = FILE_TYPE_IMAGE // 图片类型（向后兼容）
	VIDEO_TYPE = FILE_TYPE_VIDEO // 视频类型（向后兼容）
)

// 上传配置常量
const (
	UPLOAD_BASE_DIR = "uploads/" // 上传基础目录（相对于static serverRoot）
	MAX_IMAGE_SIZE  = 10         // 图片最大大小(MB)
	MAX_VIDEO_SIZE  = 100        // 视频最大大小(MB)
)

// 文件类型结构体
type FileType struct {
	Value uint8  `json:"value"` // 类型值
	Label string `json:"label"` // 类型名称
}

// 文件类型枚举列表（前端下拉框直接遍历）
var FileTypeList = []FileType{
	{Value: FILE_TYPE_IMAGE, Label: "图片"},
	{Value: FILE_TYPE_DOCUMENT, Label: "文档"},
	{Value: FILE_TYPE_VIDEO, Label: "视频"},
	{Value: FILE_TYPE_AUDIO, Label: "音频"},
	{Value: FILE_TYPE_ARCHIVE, Label: "压缩包"},
	{Value: FILE_TYPE_OTHER, Label: "其他"},
}

// 删除状态结构体
type FileDelStatus struct {
	Value uint8  `json:"value"` // 状态值
	Label string `json:"label"` // 状态名称
}

// 删除状态枚举列表（前端下拉框直接遍历）
var FileDelStatusList = []FileDelStatus{
	{Value: FILE_DEL_NORMAL, Label: "正常"},
	{Value: FILE_DEL_DELETED, Label: "已删除"},
}

// 允许的文件扩展名 - 仅图片和视频
var AllowedImageExts = []string{"jpg", "jpeg", "png", "gif", "webp", "bmp"}
var AllowedVideoExts = []string{"mp4", "avi", "mov", "wmv", "flv", "webm", "mkv"}

// 文件MIME类型映射 - 仅图片和视频
var MimeTypeMap = map[string]string{
	// 图片类型
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"gif":  "image/gif",
	"webp": "image/webp",
	"bmp":  "image/bmp",

	// 视频类型
	"mp4":  "video/mp4",
	"avi":  "video/avi",
	"mov":  "video/quicktime",
	"wmv":  "video/x-ms-wmv",
	"flv":  "video/x-flv",
	"webm": "video/webm",
	"mkv":  "video/x-matroska",
}

// GetFileTypeName 获取文件类型名称
func GetFileTypeName(fileType uint8) string {
	switch fileType {
	case FILE_TYPE_IMAGE:
		return "图片"
	case FILE_TYPE_DOCUMENT:
		return "文档"
	case FILE_TYPE_VIDEO:
		return "视频"
	case FILE_TYPE_AUDIO:
		return "音频"
	case FILE_TYPE_ARCHIVE:
		return "压缩包"
	case FILE_TYPE_OTHER:
		return "其他"
	default:
		return "未知类型"
	}
}

// GetFileDelName 获取删除状态名称
func GetFileDelName(del uint8) string {
	switch del {
	case FILE_DEL_NORMAL:
		return "正常"
	case FILE_DEL_DELETED:
		return "已删除"
	default:
		return "未知状态"
	}
}

// GetFileTypeByValue 根据值获取文件类型信息
func GetFileTypeByValue(value uint8) *FileType {
	for _, fileType := range FileTypeList {
		if fileType.Value == value {
			return &fileType
		}
	}
	return nil
}

// GetFileDelStatusByValue 根据值获取删除状态信息
func GetFileDelStatusByValue(value uint8) *FileDelStatus {
	for _, status := range FileDelStatusList {
		if status.Value == value {
			return &status
		}
	}
	return nil
}

// GetFileTypeByExt 根据文件扩展名获取文件类型
func GetFileTypeByExt(ext string) uint8 {
	for _, imageExt := range AllowedImageExts {
		if ext == imageExt {
			return IMAGE_TYPE
		}
	}

	for _, videoExt := range AllowedVideoExts {
		if ext == videoExt {
			return VIDEO_TYPE
		}
	}

	return 0 // 不支持的文件类型
}

// GetMaxSizeByType 根据文件类型获取最大允许大小(字节)
func GetMaxSizeByType(fileType uint8) int64 {
	switch fileType {
	case IMAGE_TYPE:
		return MAX_IMAGE_SIZE * 1024 * 1024
	case VIDEO_TYPE:
		return MAX_VIDEO_SIZE * 1024 * 1024
	default:
		return 0 // 不支持的文件类型
	}
}

// IsAllowedExt 检查文件扩展名是否被允许
func IsAllowedExt(ext string, fileType uint8) bool {
	switch fileType {
	case IMAGE_TYPE:
		for _, allowedExt := range AllowedImageExts {
			if ext == allowedExt {
				return true
			}
		}
	case VIDEO_TYPE:
		for _, allowedExt := range AllowedVideoExts {
			if ext == allowedExt {
				return true
			}
		}
	}
	return false
}

// IsAllowedExtAny 检查文件扩展名是否为允许的任何类型
func IsAllowedExtAny(ext string) bool {
	// 检查是否为图片格式
	for _, allowedExt := range AllowedImageExts {
		if ext == allowedExt {
			return true
		}
	}

	// 检查是否为视频格式
	for _, allowedExt := range AllowedVideoExts {
		if ext == allowedExt {
			return true
		}
	}

	return false
}
