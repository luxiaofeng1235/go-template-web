# GoFrame Template For SingleRepo

Quick Start: 
- https://goframe.org/quick

基于goframe的实现的私聊和AI搜索的实现，方便进行二进制打包。



#websocket工作流程：

1. 用户打开网站 → 前端JS自动收集设备信息
2. 生成设备指纹 → 前端计算并发送给后端
3. 后端处理 → 基于指纹生成固定的 access_key
4. 建立连接 → 用户获得固定身份，可以聊天

相关功能：
#AI搜索相关（流式输出）
#AI生图、生视频、查看生成的任务状况
#文件上传（视频、图片）
#拉取群聊+私聊历史消息
#设备认证接口
#根据AK保存用户数据（昵称、头像等）
#websocket聊天（认证AK、私聊、群聊推送长链接消息）
#获取已加入群组的群聊设摆列表
#get_ai_work_list AI生图历史记录