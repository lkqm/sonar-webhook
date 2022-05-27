package main

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"sonar-webhook/dingtalk"
	"sonar-webhook/feishu"
	"sonar-webhook/sonar"
	"sonar-webhook/wecom"
	"strings"
)

type RobotMessageSender func(key string, sonarToken string, data *sonar.WebhookData) error

// SendDingTalkRobotMessage 发送钉钉机器人消息
func SendDingTalkRobotMessage(key string, sonarToken string, data *sonar.WebhookData) error {
	message, err := buildMessageContent(data, sonarToken)
	if err != nil {
		return err
	}
	client := dingtalk.NewRobotClient(key)
	return client.SendTextMessage(message)
}

// SendWeComRobotMessage 发送企微机器人消息
func SendWeComRobotMessage(key string, sonarToken string, data *sonar.WebhookData) error {
	message, err := buildMessageContent(data, sonarToken)
	if err != nil {
		return err
	}
	client := wecom.NewRobotClient(key)
	return client.SendTextMessage(message)
}

// SendWeComMessage 发送企微应用消息
func SendWeComMessage(config *WeComConfig, sonarToken string, data *sonar.WebhookData) error {
	text, err := buildMessageContent(data, sonarToken)
	if err != nil {
		return err
	}
	email := strings.TrimSpace(data.Properties.SonarAnalysisUserEmail)
	if email == "" {
		return errors.New("properties [sonar.analysis.userEmail] can't be empty")
	}
	userId := strings.Split(email, "@")[0]

	client := wecom.Client{
		CorpId:     config.CorpId,
		CorpSecret: config.CorpSecret,
		AgentId:    config.AgentId,
	}
	message := wecom.Message{
		ToUser:                 userId,
		MsgType:                "text",
		AgentId:                config.AgentId,
		Text:                   &wecom.Text{Content: text},
		EnableDuplicateCheck:   0,
		DuplicateCheckInterval: 0,
	}
	_, err = client.SendMessage(message)
	if err != nil {
		return err
	}
	return nil
}

// FeiShuRobotMessage 发送飞书机器人信息
func FeiShuRobotMessage(key string, sonarToken string, data *sonar.WebhookData) error {
	message, err := buildMessageContent(data, sonarToken)
	if err != nil {
		return err
	}
	client := feishu.NewRobotClient(key)
	return client.SendTextMessage(message)
}

func buildMessageContent(data *sonar.WebhookData, sonarToken string) (string, error) {
	tplData := new(MessageTemplateData)
	tplMsg := `【代码检测结果】{{.AnalysisResult}}
项目分支: {{.WebhookData.Project.Name}}({{.WebhookData.Branch.Name}}){{if .Measure}}
检测结果: Bugs({{.Measure.Bugs}}), 漏洞({{.Measure.Vulnerabilities}}), 异味({{.Measure.CodeSmells}}), 重复率({{.Measure.DuplicatedLinesDensity}}%){{end}}
检测链接: {{.Url}}`
	tplData.WebhookData = data
	tplData.AnalysisResult = "通过"
	if data.QualityGate.Status == sonar.StatusError {
		tplData.AnalysisResult = "不通过"
	}
	tplData.Url = data.Branch.URL
	if tplData.Url == "" {
		tplData.Url = data.Project.URL
	}

	// 项目指标信息
	measuresComponent := getMeasuresComponent(data, sonarToken)
	if measuresComponent != nil {
		measuresData := new(MeasuresData)
		for _, measure := range measuresComponent.Measures {
			switch measure.Metric {
			case "bugs":
				measuresData.Bugs = measure.Value
			case "vulnerabilities":
				measuresData.Vulnerabilities = measure.Value
			case "code_smells":
				measuresData.CodeSmells = measure.Value
			case "duplicated_lines_density":
				measuresData.DuplicatedLinesDensity = measure.Value
			}
		}
		tplData.Measure = measuresData
	}
	tpl, err := template.New("message").Parse(tplMsg)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err = tpl.Execute(&buf, tplData); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func getMeasuresComponent(data *sonar.WebhookData, sonarToken string) *sonar.MeasuresComponent {
	if sonarToken == "" {
		return nil
	}

	client := sonar.SonarClient{
		ServerUrl: data.ServerURL,
		Token:     sonarToken,
	}
	response, err := client.GetMeasuresComponent(data.Project.Key, data.Branch.Name,
		"bugs,reliability_rating,vulnerabilities,security_rating,code_smells,sqale_rating,duplicated_lines_density,coverage")

	if err != nil {
		log.Printf("get measures response failed: %v", err)
		return nil
	}
	return &response.Component
}

type WeComConfig struct {
	CorpId     string `json:"corpId" form:"corpId" binding:"required"`         // 企业id
	CorpSecret string `json:"corpSecret" form:"corpSecret" binding:"required"` // 企业密钥
	AgentId    int    `json:"agentId" form:"agentId" binding:"required"`       // 应用id
}

// MessageTemplateData 消息模板上下问数据
type MessageTemplateData struct {
	AnalysisResult    string                   // 质量门结果
	Url               string                   // 结果访问url
	WebhookData       *sonar.WebhookData       // webhook回调数据
	MeasuresComponent *sonar.MeasuresComponent // 项目分支指标数据（原始数据）
	Measure           *MeasuresData            // 项目分支指标数据
}

type MeasuresData struct {
	Bugs                   string // Bugs数
	Vulnerabilities        string // 漏洞数
	CodeSmells             string // 异味数
	DuplicatedLinesDensity string // 代码行重复率
	ReliabilityRating      string // 可维护性
	SecurityRating         string // 安全评级
	SqaleRating            string // Sqale评级
	Coverage               string // 单元测试覆盖率
}
