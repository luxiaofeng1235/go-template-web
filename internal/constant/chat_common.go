package constant

const (
	GROUP_ID_TAG    = "GROUP_COMMON_CHAT"                                               //群聊标记,默认识别一个群聊
	DEFAULT_AVTAR   = "https://static.jsss999.com/static/common/image/default/user.png" //默认头像
	STATUS_ENABLE   = 1                                                                 //默认开启
	STATUS_DISABELD = 0                                                                 //默认禁用

	MESSAGE_TYPE_TEXT  = 1 // 纯文本消息
	MESSAGE_TYPE_IMAGE = 2 // 图片消息
	MESSAGE_TYPE_VIDEO = 3 // 视频消息/图文混合

	// 聊天类型常量 (对应chat_type字段)
	CHAT_TYPE_GROUP   = 1 // 群聊
	CHAT_TYPE_PRIVATE = 2 // 私聊

	// 通话类型常量 (对应call_type字段)
	CALL_TYPE_VOICE = 1 // 语音通话（只有音频）
	CALL_TYPE_VIDEO = 2 // 视频通话（音频+视频）

	// 会议状态常量 (对应meeting_status字段)
	MEETING_STATUS_WAITING = 1 // 等待中
	MEETING_STATUS_ACTIVE  = 2 // 进行中
	MEETING_STATUS_ENDED   = 3 // 已结束
)
