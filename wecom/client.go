package wecom

import (
	"errors"
	"fmt"
	"log"
	"sonar-webhook/util"
	"time"
)

const BaseUrl = "https://qyapi.weixin.qq.com"

var tokenCaches = make(map[string]AccessToken)

// Client 企业微信客户端
type Client struct {
	CorpId     string `json:"corpId"`     // 企业id
	CorpSecret string `json:"corpSecret"` // 应用密钥
	AgentId    int    `json:"agentId"`    // 应用id
}

func (client Client) GetAccessToken() (string, error) {
	key := client.CorpId + ":" + client.CorpSecret
	if token, ok := tokenCaches[key]; ok && !token.IsExpired() {
		return token.Token, nil
	}

	// Note: 忽略并发问题
	response, err := client.doGetAccessToken()
	if err != nil {
		return "", err
	}
	token := AccessToken{
		Token:    response.AccessToken,
		ExpireAt: time.Now().UnixMilli() + response.ExpiresIn*1000,
	}
	tokenCaches[key] = token

	return token.Token, nil
}

func (client Client) SendMessage(message Message) (string, error) {
	token, err := client.GetAccessToken()
	if err != nil {
		return "", err
	}

	response := new(SendMessageResponse)
	url := BaseUrl + fmt.Sprintf("/cgi-bin/message/send?access_token=%s", token)
	if content, err := util.HttpPostEntity(url, message, response, nil); err != nil {
		log.Printf("list user by department failed, the result: %s", content)
		return "", err
	}
	if !response.IsOk() {
		errorMsg := fmt.Sprintf("list user by department failed: errcode=%d, errmsg=%s", response.ErrCode, response.ErrMsg)
		return response.MsgId, errors.New(errorMsg)
	}
	return response.MsgId, nil
}

func (client Client) ListUserByDepartment(departId int64, fetchChild int) ([]UserInfo, error) {
	token, err := client.GetAccessToken()
	if err != nil {
		return nil, err
	}

	response := new(ListUserResponse)
	url := BaseUrl + fmt.Sprintf("/cgi-bin/user/list?access_token=%s&department_id=%d&fetch_child=%d", token, departId, fetchChild)
	if content, err := util.HttpGetJson(url, response); err != nil {
		log.Printf("list user by department failed, the result: %s", content)
		return nil, err
	}
	if !response.IsOk() {
		errorMsg := fmt.Sprintf("list user by department failed: errcode=%d, errmsg=%s", response.ErrCode, response.ErrMsg)
		return nil, errors.New(errorMsg)
	}
	return response.UserList, nil
}

func (client Client) doGetAccessToken() (*GetAccessTokenResponse, error) {
	url := BaseUrl + fmt.Sprintf("/cgi-bin/gettoken?corpid=%s&corpsecret=%s", client.CorpId, client.CorpSecret)
	response := new(GetAccessTokenResponse)
	if content, err := util.HttpGetJson(url, response); err != nil {
		log.Printf("get access token failed, the result: %s", content)
		return nil, err
	}
	if !response.IsOk() {
		errorMsg := fmt.Sprintf("get access token failed: errcode=%d, errmsg=%s", response.ErrCode, response.ErrMsg)
		return nil, errors.New(errorMsg)
	}
	return response, nil
}
