package feishu

import (
	"fmt"
	"log"
	"sonar-webhook/sdks/util"
)

const RobotWebhookUrl = "https://open.feishu.cn/open-apis/bot/v2/hook/%s"

// RobotClient 机器人客户端
type RobotClient struct {
	Token string
}

func NewRobotClient(token string) RobotClient {
	return RobotClient{Token: token}
}

// SendMessage 发送消息
func (robot RobotClient) SendMessage(message Message) error {
	url := fmt.Sprintf(RobotWebhookUrl, robot.Token)
	content, err := util.HttpPostJson(url, message, nil)
	if err != nil {
		return err
	}

	log.Printf("send robot message response content: %s", content)
	return nil
}

// SendTextMessage 发送文本消息
func (robot RobotClient) SendTextMessage(text string) error {
	message := Message{
		MsgType: "text",
		Content: MessageContent{Text: text},
	}
	return robot.SendMessage(message)
}

// Message 机器人消息
type Message struct {
	MsgType string         `json:"msg_type"` // 消息类型
	Content MessageContent `json:"content"`  // 消息内容
}

type MessageContent struct {
	Text string `json:"text"` // 普通文本消息内容
}
