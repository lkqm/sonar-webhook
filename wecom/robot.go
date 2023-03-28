package wecom

import (
	"fmt"
	"log"
	"sonar-webhook/util"
)

const RobotWebhookUrl = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s"

// RobotClient 机器人客户端
type RobotClient struct {
	Key string
}

func NewRobotClient(key string) RobotClient {
	return RobotClient{Key: key}
}

// SendMessage 发送消息
func (robot RobotClient) SendMessage(message RobotMessage) error {
	url := fmt.Sprintf(RobotWebhookUrl, robot.Key)
	content, err := util.HttpPostJson(url, message, nil)
	if err != nil {
		return err
	}

	log.Printf("send wecom message response content: %s", content)
	return nil
}

// SendTextMessage 发送markdown消息
func (robot RobotClient) SendTextMessage(content string) error {
	message := RobotMessage{
		MsgType: "text",
		Text: &Text{
			Content: content,
		},
	}
	return robot.SendMessage(message)
}

// SendMarkdownMessage 发送markdown消息
func (robot RobotClient) SendMarkdownMessage(content string) error {
	message := RobotMessage{
		MsgType: "markdown",
		Markdown: &Markdown{
			Content: content,
		},
	}
	return robot.SendMessage(message)
}

// RobotMessage 机器人webhook请求参数
type RobotMessage struct {
	MsgType  string    `json:"msgtype"`  // 消息类型
	Text     *Text     `json:"text"`     // 文本消息
	Markdown *Markdown `json:"markdown"` // markdown消息
}

type Markdown struct {
	Content string `json:"content"` // markdown内容
}

type Text struct {
	Content string `json:"content"` // 文本内容
}
