/*
 * @file: device.go
 * @description: 设备管理，主要处理设备的 saveUserData deviceAuth getDeviceList三个接口的实现
 * @author: red
 * @created: 2025-09-10
 * @version: 1.0.0
 * @license: MIT License
 */

package api

import (
	"crypto/md5"
	"fmt"
	"go-web-template/global"
	"go-web-template/internal/constant"
	"go-web-template/internal/models"
	"go-web-template/utils"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"gorm.io/gorm"
)

// GetByDeviceFingerprint 根据设备指纹获取对应的密钥记录
// 参数:
//   - deviceFingerprint: 设备指纹标识
//
// 返回值:
//   - record: 找到的SecretKey记录
//   - err: 错误信息
func GetByDeviceFingerprint(deviceFingerprint string) (record *models.SecretKey, err error) {
	if deviceFingerprint == "" {
		return nil, fmt.Errorf("设备指纹不能为空")
	}

	var secretKey models.SecretKey
	err = global.DB.Where("device_fingerprint = ? AND status = ?", deviceFingerprint, constant.SecretKeyStatusNormal).First(&secretKey).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 未找到记录，返回nil而不是错误
		}
		global.Errlog.Error("查询设备指纹失败", "deviceFingerprint", deviceFingerprint, "error", err)
		return nil, err
	}

	return &secretKey, nil
}

// generateUniqueAccessKey 生成唯一的访问密钥（处理重复情况，注意性能）
// 参数:
//   - deviceFingerprint: 设备指纹
//   - systemSecret: 系统密钥
//
// 返回值:
//   - accessKey: 唯一的访问密钥
//   - err: 错误信息
func generateUniqueAccessKey(deviceFingerprint, systemSecret string) (string, error) {
	// 基础访问密钥生成 - 按照PHP版本逻辑
	hashInput := deviceFingerprint + systemSecret
	hash := fmt.Sprintf("%x", md5.Sum([]byte(hashInput)))
	// PHP版本使用 strtoupper(substr($hash, 0, 12))
	accessKey := "AK_" + strings.ToUpper(hash[0:12])

	// 检查是否重复（只检查一次，避免性能影响）
	var count int64
	err := global.DB.Model(&models.SecretKey{}).Where("access_key = ?", accessKey).Count(&count).Error
	if err != nil {
		global.Errlog.Error("检查访问密钥重复失败", "accessKey", accessKey, "error", err)
		return "", fmt.Errorf("检查访问密钥重复失败")
	}

	if count > 0 {
		// 极小概率重复，使用时间戳重新生成 - 按照PHP版本逻辑
		timestamp := fmt.Sprintf("%.3f", float64(utils.GetUnixNano())/1e9) // 使用microtime(true)等效
		hashInput = deviceFingerprint + systemSecret + timestamp
		hash = fmt.Sprintf("%x", md5.Sum([]byte(hashInput)))
		accessKey = "AK_" + strings.ToUpper(hash[0:12])
	}

	return accessKey, nil
}

// GenerateKeyMapping 通用的访问密钥生成逻辑（固定群组模式）
// 参数:
//   - deviceFingerprint: 设备指纹
//   - userInputKey: 用户输入的密钥（仅用于身份验证）
//
// 返回值:
//   - result: 包含访问密钥和群组ID的映射
//   - err: 错误信息
func GenerateKeyMapping(deviceFingerprint, userInputKey string) (map[string]string, error) {
	// 参数验证
	if deviceFingerprint == "" {
		return nil, fmt.Errorf("设备指纹不能为空")
	}

	// 获取系统密钥作为盐值
	systemSecret := g.Cfg().MustGet(gctx.New(), "security.secretKey", "").String()
	if systemSecret == "" {
		return nil, fmt.Errorf("系统密钥未配置")
	}

	// 生成访问密钥，处理重复情况
	accessKey, err := generateUniqueAccessKey(deviceFingerprint, systemSecret)
	if err != nil {
		return nil, err
	}

	// ⚠️ 重要：固定群组模式 - 所有用户永远在同一个群聊
	// 用户输入的密钥仅用于身份验证，不用于群组划分
	return map[string]string{
		"access_key": accessKey,             // 基于设备指纹生成唯一访问密钥
		"group_id":   constant.GROUP_ID_TAG, // 固定群组：GROUP_COMMON_CHAT
	}, nil
}

// GenerateNickname 基于设备指纹生成匿名昵称
// 参数:
//   - deviceFingerprint: 设备指纹
//
// 返回值:
//   - nickname: 生成的昵称
func GenerateNickname(deviceFingerprint string) string {
	// 基于设备指纹生成，确保同设备昵称一致且不重复
	hash := fmt.Sprintf("%x", md5.Sum([]byte(deviceFingerprint+"nickname_salt_2024")))
	randomNum := hash[0:8] // 取MD5前8位
	return fmt.Sprintf("匿名用户_%s", randomNum)
}

// SaveDeviceAccess 保存设备访问密钥信息到数据库
// 参数:
//   - deviceFingerprint: 设备指纹
//   - accessKey: 访问密钥
//   - groupId: 群组ID
//   - nickname: 昵称
//   - deviceInfo: 设备信息
//
// 返回值:
//   - result: 保存结果
//   - err: 错误信息
func SaveDeviceAccess(deviceFingerprint, accessKey, groupId, nickname, deviceInfo string) (*models.CreateSecretResp, error) {
	currentTime := utils.GetUnix()
	avatarURL := constant.DEFAULT_AVTAR

	data := models.SecretKey{
		DeviceFingerprint: deviceFingerprint,
		AccessKey:         accessKey,
		GroupID:           constant.GROUP_ID_TAG, // 默认先用这个默认的群组ID
		Nickname:          nickname,
		AvtarURL:          avatarURL, // 默认头像
		DeviceInfo:        deviceInfo,
		FirstVisitTime:    currentTime,
		LastVisitTime:     currentTime,
		VisitCount:        1,
		Status:            constant.STATUS_ENABLE,
		CreatedAt:         currentTime,
		UpdatedAt:         currentTime,
	}

	err := global.DB.Create(&data).Error
	if err != nil {
		global.Errlog.Error("保存设备访问信息失败", "deviceFingerprint", deviceFingerprint, "error", err)
		return nil, fmt.Errorf("保存设备访问信息失败")
	}

	// 确保获取到正确的ID
	global.Errlog.Info("设备创建成功", "deviceFingerprint", deviceFingerprint, "id", data.ID)

	return &models.CreateSecretResp{
		Id:        data.ID,
		AccessKey: accessKey,
		GROUPID:   groupId,
		AVTAR_URL: avatarURL,
		NickName:  nickname,
		IsNew:     true,
	}, nil
}

// GetOrCreateDeviceAccess 根据设备指纹获取或创建设备访问信息
// 参数:
//   - req: 创建密钥请求结构体
//
// 返回值:
//   - result: 设备访问信息
//   - err: 错误信息
func GetOrCreateDeviceAccess(req *models.CreateSecretKeyReq) (*models.CreateSecretResp, error) {
	// 参数验证
	if req == nil {
		return nil, fmt.Errorf("请求参数不能为空")
	}
	if req.DeviceFingerprint == "" {
		return nil, fmt.Errorf("设备指纹不能为空")
	}

	// 先检查设备是否已存在
	existing, err := GetByDeviceFingerprint(req.DeviceFingerprint)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		// 使用原子操作更新最后访问时间和访问次数
		currentTime := utils.GetUnix()
		err = global.DB.Model(&models.SecretKey{}).
			Where("device_fingerprint = ?", req.DeviceFingerprint).
			UpdateColumn("visit_count", gorm.Expr("visit_count + ?", 1)).Error

		if err == nil {
			// 然后更新其他字段
			err = global.DB.Model(&models.SecretKey{}).
				Where("device_fingerprint = ?", req.DeviceFingerprint).
				Updates(map[string]interface{}{
					"last_visit_time": currentTime,
					"updated_at":      currentTime,
				}).Error
		}

		if err != nil {
			global.Errlog.Error("更新访问信息失败", "deviceFingerprint", req.DeviceFingerprint, "error", err)
			return nil, fmt.Errorf("更新访问信息失败")
		}

		return &models.CreateSecretResp{
			Id:        existing.ID,
			AccessKey: existing.AccessKey,
			GROUPID:   existing.GroupID,
			AVTAR_URL: existing.AvtarURL,
			NickName:  existing.Nickname,
			IsNew:     false, // 默认存在就是false
		}, nil
	}

	// 生成新的访问密钥和群组信息
	keyMapping, err := GenerateKeyMapping(req.DeviceFingerprint, "")
	if err != nil {
		return nil, err
	}

	nickname := GenerateNickname(req.DeviceFingerprint)

	// 保存到数据库
	result, err := SaveDeviceAccess(
		req.DeviceFingerprint,
		keyMapping["access_key"],
		keyMapping["group_id"],
		nickname,
		req.DeviceInfo,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
