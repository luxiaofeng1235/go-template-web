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
	"go-web-template/internal/models"
	"go-web-template/utils"
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
//   - req: 获取聊天历史请求，包含用户ID、对方用户ID、分页等参数
//
// 返回值:
//   - list: 聊天历史记录列表
//   - total: 符合条件的总记录数
//   - err: 错误信息
func (s *MessageService) GetChatHistoryByParams(req *models.ChatHistoryReq) (list []models.ChatHistoryRes, total int64, err error) {
	// 参数验证
	if req.UserId <= 0 {
		err = fmt.Errorf("用户ID不能为空")
		return
	}
	if req.ToUserId <= 0 {
		err = fmt.Errorf("对方用户ID不能为空")
		return
	}

	// 设置默认分页参数
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// 构建查询条件 - 查询双向聊天记录
	db := global.DB.Model(&models.ChatMessage{}).
		Where("((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?))",
			req.UserId, req.ToUserId, req.ToUserId, req.UserId)

	// 添加时间范围筛选（如果提供）
	if req.StartTime > 0 {
		db = db.Where("created_at >= ?", req.StartTime)
	}
	if req.EndTime > 0 {
		db = db.Where("created_at <= ?", req.EndTime)
	}

	// 添加消息类型筛选（如果提供）
	if req.MessageType > 0 {
		db = db.Where("message_type = ?", req.MessageType)
	}

	// 按时间倒序排列
	db = db.Order("created_at desc")

	// 统计总条数
	err = db.Count(&total).Error
	if err != nil {
		global.Errlog.Error("查询聊天记录总数失败", "userId", req.UserId, "toUserId", req.ToUserId, "error", err)
		return
	}

	// 分页查询
	if req.PageNum > 0 && req.PageSize > 0 {
		err = db.Offset((req.PageNum - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error
	} else {
		err = db.Offset(req.PageNum).Find(&list).Error
	}

	if err != nil {
		global.Errlog.Error("查询聊天历史失败", "userId", req.UserId, "toUserId", req.ToUserId, "error", err)
		return
	}

	// 容错判断：空列表处理
	if len(list) <= 0 {
		return list, total, nil
	}

	// 将结果按时间正序排列（最新消息在最后）
	for i := 0; i < len(list)/2; i++ {
		j := len(list) - 1 - i
		list[i], list[j] = list[j], list[i]
	}

	return list, total, err
}

// SendMessage 发送消息
// 参数:
//   - req: 发送消息请求，包含发送者、接收者、消息内容等信息
//
// 返回值:
//   - result: 发送结果，包含消息ID等信息
//   - err: 错误信息
func (s *MessageService) SendMessage(req *models.SendMessageReq) (result *models.SendMessageRes, err error) {
	// 参数验证
	if req.SecretKey == "" {
		err = fmt.Errorf("密钥不能为空")
		return
	}
	if req.ReceiverID == "" {
		err = fmt.Errorf("接收者ID不能为空")
		return
	}
	if req.Content == "" && req.MessageType == 1 {
		err = fmt.Errorf("文本消息内容不能为空")
		return
	}

	// 验证SecretKey并获取发送者信息
	var secretKey models.SecretKey
	err = global.DB.Where("access_key = ? AND status = ?", req.SecretKey, models.SecretKeyStatusNormal).First(&secretKey).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("无效的密钥")
			return
		}
		global.Errlog.Error("查询密钥失败", "secretKey", req.SecretKey, "error", err)
		return
	}

	// 生成消息ID
	messageId := utils.GenerateMessageId()

	// 创建消息记录
	now := utils.GetUnix()
	message := models.ChatMessage{
		MessageID:   messageId,
		SecretKey:   req.SecretKey,
		SenderID:    secretKey.DeviceFingerprint,
		SenderName:  secretKey.Nickname,
		ReceiverID:  req.ReceiverID,
		MessageType: req.MessageType,
		Content:     req.Content,
		ImageURL:    req.ImageURL,
		VideoURL:    req.VideoURL,
		ChatType:    req.ChatType,
		ExtraData:   req.ExtraData,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = global.DB.Create(&message).Error
	if err != nil {
		global.Sqllog.Error("创建消息记录失败", "messageId", messageId, "error", err)
		return
	}

	result = &models.SendMessageRes{
		MessageID:   message.ID,
		UserId:      0, // 根据实际需求调整
		ToUserId:    0, // 根据实际需求调整
		Content:     message.Content,
		MessageType: message.MessageType,
		CreateTime:  message.CreatedAt,
	}

	return result, nil
}
