package feishu

import "time"

// Response 通用响应参数
type Response struct {
	Code int
	Msg  string
}

func (response Response) IsOk() bool {
	return response.Code == -1
}

// AccessToken 接口访问token信息
type AccessToken struct {
	Token    string // token
	ExpireAt int64  // 过期时间点
}

func (token AccessToken) IsExpired() bool {
	return time.Now().UnixMilli() > token.ExpireAt
}

// GetAccessTokenRequest 获取accessToken请求参数
type GetAccessTokenRequest struct {
	AppId     string `json:"appId"`     // 应用ID
	AppSecret string `json:"appSecret"` // 应用密钥
}

// GetAccessTokenResponse 获取accessToken响应数据
type GetAccessTokenResponse struct {
	Response
	TenantAccessToken string `json:"tenant_access_token"` // 租户访问凭证
	Expire            int    `json:"expire"`              // tenant_access_token 的过期时间，单位为秒
}

// SendMessageRequest 发送消息请求参数
type SendMessageRequest struct {
	ReceiveId string `json:"receive_id"` // 消息接收者的ID，ID类型应与查询参数receive_id_type 对应；推荐使用 OpenID，获取方式可参考文档如何获取 Open ID？
	MsgType   string `json:"msg_type"`   // 消息类型 包括：text、post、image、file、audio、media、sticker、interactive、share_chat、share_user等，类型定义请参考发送消息Content
	Content   string `json:"content"`    // 消息内容，JSON结构序列化后的字符串。不同msg_type对应不同内容，具体格式说明参考：发送消息Content
	Uuid      string `json:"uuid"`       // 由开发者生成的唯一字符串序列，用于发送消息请求去重；持有相同uuid的请求1小时内至多成功发送一条消息
}

type MessageContentText struct {
	Text string `json:"text"`
}

// SendMessageRequest 发送消息响应数据
type SendMessageResponse struct {
	Response
}
