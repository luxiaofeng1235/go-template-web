package service

import (
	"fmt"
	"go-web-template/utils"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go-web-template/global"
	"go-web-template/internal/constant"
	"go-web-template/internal/models"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/grand"
)

// UploadFile 通用文件上传方法
func UploadFile(r *ghttp.Request, req *models.UploadFileReq) (*models.UploadFileRes, error) {
	// 获取上传的文件
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请选择要上传的文件")
	}

	// 获取文件扩展名
	originalName := file.Filename
	if originalName == "" {
		return nil, fmt.Errorf("文件名不能为空")
	}

	extension := strings.ToLower(filepath.Ext(originalName))
	if extension == "" {
		return nil, fmt.Errorf("文件必须有扩展名")
	}

	// 去掉点号
	ext := strings.TrimPrefix(extension, ".")

	// 检查文件类型是否被允许
	if !constant.IsAllowedExtAny(ext) {
		return nil, fmt.Errorf("不支持的文件格式，仅支持图片格式：%s 和视频格式：%s",
			strings.Join(constant.AllowedImageExts, ", "),
			strings.Join(constant.AllowedVideoExts, ", "))
	}

	// 获取文件类型
	fileType := constant.GetFileTypeByExt(ext)
	if fileType == 0 {
		return nil, fmt.Errorf("无法确定文件类型")
	}

	// 检查文件大小
	maxSize := constant.GetMaxSizeByType(fileType)
	if file.Size > maxSize {
		var typeName string
		if fileType == constant.IMAGE_TYPE {
			typeName = "图片"
		} else {
			typeName = "视频"
		}
		maxSizeMB := maxSize / (1024 * 1024)
		return nil, fmt.Errorf("%s文件大小不能超过%dMB", typeName, maxSizeMB)
	}

	// 记录上传请求日志
	global.Requestlog.Info("开始文件上传",
		"filename", originalName,
		"size", file.Size,
		"type", fileType,
		"userID", req.UserID)

	// 创建按日期分组的上传目录
	dateDir := time.Now().Format("20060102")
	uploadDir := fmt.Sprintf("%s%s/", constant.UPLOAD_BASE_DIR, dateDir)

	// 确保上传目录存在
	if !gfile.Exists(uploadDir) {
		if err := gfile.Mkdir(uploadDir); err != nil {
			global.Errlog.Error("创建上传目录失败", "dir", uploadDir, "error", err)
			return nil, fmt.Errorf("创建上传目录失败")
		}
	}

	// 生成唯一文件名
	uniqueFileName := fmt.Sprintf("%s_%d.%s", grand.S(10), time.Now().Unix(), ext)
	fullPath := filepath.Join(uploadDir, uniqueFileName)

	// 保存文件到本地
	if err := saveUploadedFile(file, fullPath); err != nil {
		global.Errlog.Error("保存文件失败", "path", fullPath, "error", err)
		return nil, fmt.Errorf("文件上传失败")
	}

	// 处理原始文件名长度限制
	processedName := originalName
	if len(originalName) > 100 {
		// 截取前95个字符 + 扩展名
		nameWithoutExt := strings.TrimSuffix(originalName, extension)
		if len(nameWithoutExt) > 95 {
			nameWithoutExt = nameWithoutExt[:95]
		}
		processedName = nameWithoutExt + extension
	}

	// 相对路径用于数据库存储
	relativeURI := fmt.Sprintf("/%s%s/%s", constant.UPLOAD_BASE_DIR, dateDir, uniqueFileName)

	// 保存文件记录到数据库
	fileRecord := models.File{
		Name:       processedName,
		CID:        req.CID,
		Type:       fileType,
		URI:        relativeURI,
		ShopID:     req.ShopID,
		UserID:     req.UserID,
		CreateTime: utils.GetUnix(),
		Del:        0,
	}

	if err := global.DB.Create(&fileRecord).Error; err != nil {
		// 数据库保存失败，删除已上传的文件
		gfile.Remove(fullPath)
		global.Errlog.Error("保存文件记录到数据库失败", "error", err)
		return nil, fmt.Errorf("文件记录保存失败")
	}

	// 生成访问URL - 使用URL工具函数统一生成静态资源URL
	fullURL := utils.GetStaticResourceURL(r, relativeURI)

	// 记录上传成功日志
	global.Requestlog.Info("文件上传成功",
		"fileID", fileRecord.ID,
		"filename", processedName,
		"path", relativeURI,
		"userID", req.UserID)

	// 构造返回结果
	result := &models.UploadFileRes{
		ID:       fileRecord.ID,
		Name:     processedName,
		Type:     fileType,
		TypeName: constant.GetFileTypeName(fileType),
		Size:     file.Size,
		URI:      fullURL,
		BaseURI:  relativeURI,
		CID:      req.CID,
		ShopID:   req.ShopID,
		UserID:   req.UserID,
	}

	return result, nil
}

// saveUploadedFile 保存上传的文件到指定路径
func saveUploadedFile(fileHeader *ghttp.UploadFile, destPath string) error {
	// 打开上传的文件
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// 复制文件内容
	_, err = io.Copy(dst, src)
	return err
}

// UploadImage 图片上传方法
func UploadImage(r *ghttp.Request, req *models.UploadFileReq) (*models.UploadFileRes, error) {
	// 获取上传的文件
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请选择要上传的图片")
	}

	// 获取文件扩展名
	extension := strings.ToLower(filepath.Ext(file.Filename))
	ext := strings.TrimPrefix(extension, ".")

	// 检查是否为图片格式
	if !constant.IsAllowedExt(ext, constant.IMAGE_TYPE) {
		return nil, fmt.Errorf("不支持的图片格式，仅支持：%s",
			strings.Join(constant.AllowedImageExts, ", "))
	}

	return UploadFile(r, req)
}

// UploadVideo 视频上传方法
func UploadVideo(r *ghttp.Request, req *models.UploadFileReq) (*models.UploadFileRes, error) {
	// 获取上传的文件
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请选择要上传的视频")
	}

	// 获取文件扩展名
	extension := strings.ToLower(filepath.Ext(file.Filename))
	ext := strings.TrimPrefix(extension, ".")

	// 检查是否为视频格式
	if !constant.IsAllowedExt(ext, constant.VIDEO_TYPE) {
		return nil, fmt.Errorf("不支持的视频格式，仅支持：%s",
			strings.Join(constant.AllowedVideoExts, ", "))
	}

	return UploadFile(r, req)
}

// UploadImageSimple 简化的图片上传方法，按照PHP模式
func UploadImageSimple(r *ghttp.Request, cid int, shopID int, userID int) (*models.UploadFileRes, error) {
	req := &models.UploadFileReq{
		CID:    uint(cid),
		ShopID: shopID,
		UserID: uint(userID),
	}
	return UploadImage(r, req)
}

// UploadVideoSimple 简化的视频上传方法，按照PHP模式
func UploadVideoSimple(r *ghttp.Request, cid int, shopID int, userID int) (*models.UploadFileRes, error) {
	req := &models.UploadFileReq{
		CID:    uint(cid),
		ShopID: shopID,
		UserID: uint(userID),
	}
	return UploadVideo(r, req)
}
