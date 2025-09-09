package constant

// 文件类型常量 - 仅支持图片和视频
const (
	IMAGE_TYPE = 10 // 图片类型
	VIDEO_TYPE = 20 // 视频类型
)

// 上传配置常量
const (
	UPLOAD_BASE_DIR = "public/uploads/" // 上传基础目录
	MAX_IMAGE_SIZE  = 10                // 图片最大大小(MB)
	MAX_VIDEO_SIZE  = 100               // 视频最大大小(MB)
)

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
