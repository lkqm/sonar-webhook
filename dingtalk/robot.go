package dingtalk

import (
	"fmt"
	"log"
	"sonar-webhook/util"
)

const RobotWebhookUrl = "https://oapi.dingtalk.com/robot/send?access_token=%s"

// RobotClient 机器人客户端
type RobotClient struct {
	AccessToken string
}

func NewRobotClient(accessToken string) RobotClient {
	return RobotClient{AccessToken: accessToken}
}

// SendMessage 发送消息
func (robot RobotClient) SendMessage(message RobotMessage) error {
	url := fmt.Sprintf(RobotWebhookUrl, robot.AccessToken)
	content, err := util.HttpPostJson(url, message)
	if err != nil {
		return err
	}

	log.Printf("send robot message response content: %s", content)
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
func (robot RobotClient) SendMarkdownMessage(title, text string) error {
	message := RobotMessage{
		MsgType: "markdown",
		Markdown: &Markdown{
			Title: title,
			Text:  text,
		},
	}
	return robot.SendMessage(message)
}

// RobotMessage 机器人消息
type RobotMessage struct {
	MsgType  string    `json:"msgtype"`  // 消息类型
	Text     *Text     `json:"text"`     // 文本消息
	Markdown *Markdown `json:"markdown"` // markdown消息
}

type Text struct {
	Content string `json:"content"` // 文本内容
}

type Markdown struct {
	Title string `json:"title"` // 标题
	Text  string `json:"text"`  // 内容
}
