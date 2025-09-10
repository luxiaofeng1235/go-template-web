/*
 * @Description: 消息管理服务层，处理消息群聊私聊、获取消息列表设备等功能
 * @Version: 1.0.0
 * @Author: red
 * @Date: 2025-01-10
 * @LastEditTime: 2025-01-10
 */

package api

import (
	"errors"
	"fmt"
	"go-web-template/global"
	"go-web-template/internal/constant"
	"go-web-template/internal/models"

	"gorm.io/gorm"
)

// MessageService 消息服务结构体
type MessageService struct{}

// GetChatUserList 获取聊天用户列表（对应PHP的getChatUserList）
// 参数:
//   - req: 获取设备列表请求，包含访问密钥和分页参数
//
// 返回值:
//   - list: SecretKey列表，直接返回数据库字段
//   - total: 符合条件的总记录数
//   - err: 错误信息
func (s *MessageService) GetChatUserList(req *models.ChatUserListReq) (list []models.SecretKey, total int64, err error) {
	// 参数验证
	if req.AccessKey == "" {
		err = fmt.Errorf("访问密钥不能为空")
		return
	}

	// 设置默认分页参数
	if req.PageNo <= 0 {
		req.PageNo = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20 // 默认20条
	}

	// 构建查询条件 - 查询SecretKey表获取设备列表，选择指定字段
	db := global.DB.Model(&models.SecretKey{}).
		Select("id, device_fingerprint, access_key, nickname, user_note, avtar_url").
		Where("status = ?", models.SecretKeyStatusNormal)

	// 统计总条数
	err = db.Count(&total).Error
	if err != nil {
		global.Errlog.Error("查询聊天用户总数失败", "accessKey", req.AccessKey, "error", err)
		return
	}

	// 分页查询
	err = db.Offset((req.PageNo - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error
	if err != nil {
		global.Errlog.Error("查询聊天用户列表失败", "accessKey", req.AccessKey, "error", err)
		return
	}

	// 容错判断：空列表处理
	if len(list) <= 0 {
		return list, total, nil
	}

	// 如果传了access_key，关联查询用户备注信息
	if req.AccessKey != "" {
		for i, secretKey := range list {
			// 检查是否为自己，如果不是自己则查询用户备注
			if secretKey.AccessKey != req.AccessKey {
				var userNote models.UserNote
				err = global.DB.Where("access_key = ? AND to_access_key = ?", req.AccessKey, secretKey.AccessKey).First(&userNote).Error
				if err == nil && userNote.UserNote != "" {
					// 如果找到用户备注，则使用备注替换user_note字段
					list[i].UserNote = userNote.UserNote
				} else if !errors.Is(err, gorm.ErrRecordNotFound) {
					global.Errlog.Error("查询用户备注失败", "accessKey", req.AccessKey, "toAccessKey", secretKey.AccessKey, "error", err)
				}
			}
		}
	}

	return list, total, nil
}

// GetChatHistoryByParams 根据参数获取聊天历史记录（对应PHP的getChatHistoryByParams）
// 参数:
//   - req: 获取聊天历史请求，支持群聊和私聊两种模式
//
// 返回值:
//   - list: 聊天历史记录列表
//   - total: 符合条件的总记录数
//   - err: 错误信息
func (s *MessageService) GetChatHistoryByParams(req *models.ChatHistoryReq) (list []models.ChatMessage, total int64, err error) {
	// 设置默认聊天类型为群聊
	if req.ChatType <= 0 {
		req.ChatType = constant.CHAT_TYPE_GROUP // 默认群聊
	}
	// 设置默认分页参数
	if req.PageNo <= 0 {
		req.PageNo = constant.PAGE_NO
	}
	if req.PageSize <= 0 {
		req.PageSize = constant.PAGE_SIZE
	}

	// 构建查询条件
	db := global.DB.Model(&models.ChatMessage{}).Where("chat_type = ?", req.ChatType)

	// 根据聊天类型构建不同的查询条件
	if req.ChatType == 2 { // 私聊模式
		// 私聊时需要验证参数
		if req.AccessKey == "" || req.ReceiverID == "" {
			err = fmt.Errorf("私聊模式下访问密钥和接收者ID不能为空")
			return
		}

		// 验证AccessKey，获取发送者信息
		var secretKey models.SecretKey
		err = global.DB.Where("access_key = ? AND status = ?", req.AccessKey, models.SecretKeyStatusNormal).First(&secretKey).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = fmt.Errorf("无效的访问密钥")
				return
			}
			global.Errlog.Error("查询访问密钥失败", "accessKey", req.AccessKey, "error", err)
			return
		}

		// 查询双向私聊记录（发送者与接收者之间的对话）
		db = db.Where("((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?))",
			secretKey.DeviceFingerprint, req.ReceiverID, req.ReceiverID, secretKey.DeviceFingerprint)

	} else { // 群聊模式（chat_type = 1）
		// 群聊可以不需要特定的发送者和接收者限制，显示所有群聊消息
		// 如果提供了access_key，可以用于权限验证
		if req.AccessKey != "" {
			var secretKey models.SecretKey
			err = global.DB.Where("access_key = ? AND status = ?", req.AccessKey, models.SecretKeyStatusNormal).First(&secretKey).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					err = fmt.Errorf("无效的访问密钥")
					return
				}
				global.Errlog.Error("查询访问密钥失败", "accessKey", req.AccessKey, "error", err)
				return
			}
		}
	}
	// 按时间倒序排列
	db = db.Order("created_at desc")
	// 统计总条数
	err = db.Count(&total).Error
	if err != nil {
		global.Errlog.Error("查询聊天记录总数失败", "chatType", req.ChatType, "error", err)
		return
	}
	// 分页查询
	err = db.Offset((req.PageNo - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error
	if err != nil {
		global.Errlog.Error("查询聊天历史失败", "chatType", req.ChatType, "error", err)
		return
	}

	// 容错判断：空列表处理
	if len(list) <= 0 {
		return list, total, nil
	}
	return list, total, err
}
