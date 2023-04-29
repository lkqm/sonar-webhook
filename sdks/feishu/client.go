package feishu

import (
	"errors"
	"fmt"
	"log"
	"sonar-webhook/sdks/util"
	"time"
)

const BaseUrl = "https://open.feishu.cn/open-apis"

var tokenCaches = make(map[string]AccessToken)

// Client 飞书客户端
type Client struct {
	AppId     string `json:"appId"`     // 应用ID
	AppSecret string `json:"appSecret"` // 应用密钥
}

func (client Client) GetAccessToken() (string, error) {
	key := client.AppId + ":" + client.AppSecret
	if token, ok := tokenCaches[key]; ok && !token.IsExpired() {
		return token.Token, nil
	}

	// Note: 忽略并发问题
	response, err := client.doGetAccessToken()
	if err != nil {
		return "", err
	}
	token := AccessToken{
		Token:    "Bearer " + response.TenantAccessToken,
		ExpireAt: time.Now().UnixMilli() + int64(response.Expire)*1000,
	}
	tokenCaches[key] = token

	return token.Token, nil
}

func (client Client) SendMessage(receiveIdType string, request SendMessageRequest) (*SendMessageResponse, error) {
	accessToken, err := client.GetAccessToken()
	if err != nil {
		return nil, err
	}
	url := BaseUrl + "/im/v1/messages?receive_id_type=" + receiveIdType
	headers := map[string]string{"Authorization": accessToken}
	response := new(SendMessageResponse)
	if content, err := util.HttpPostEntity(url, request, response, headers); err != nil {
		log.Printf("send message failed, the result: %s", content)
		return nil, err
	}
	if !response.IsOk() {
		errorMsg := fmt.Sprintf("send message failed: code=%d, msg=%s", response.Code, response.Msg)
		return nil, errors.New(errorMsg)
	}
	return response, nil
}

func (client Client) doGetAccessToken() (*GetAccessTokenResponse, error) {
	url := BaseUrl + "/auth/v3/tenant_access_token/internal"
	request := &GetAccessTokenRequest{AppId: client.AppId, AppSecret: client.AppSecret}
	response := new(GetAccessTokenResponse)
	if content, err := util.HttpPostEntity(url, request, response, nil); err != nil {
		log.Printf("get access token failed, the result: %s", content)
		return nil, err
	}
	if !response.IsOk() {
		errorMsg := fmt.Sprintf("get access token failed: code=%d, msg=%s", response.Code, response.Msg)
		return nil, errors.New(errorMsg)
	}
	return response, nil
}
