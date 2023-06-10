package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"sonar-webhook/sdks/dingtalk"
	"sonar-webhook/sdks/feishu"
	"sonar-webhook/sdks/sonar"
	"sonar-webhook/sdks/wecom"
	"strings"
	"text/template"
)

type RobotMessageSender func(p *RobotWebhookParameters, data *sonar.WebhookData) error

// SendDingTalkRobotMessage 发送钉钉机器人消息
func SendDingTalkRobotMessage(p *RobotWebhookParameters, data *sonar.WebhookData) error {
	buildArgs := MessageBuildArgs{
		WebhookData: data,
		SonarToken:  p.SonarToken,
		NewCode:     p.NewCode,
	}
	message, err := buildMessageContent(buildArgs)
	if err != nil {
		return err
	}
	client := dingtalk.NewRobotClient(p.Key)
	return client.SendTextMessage(message)
}

// SendWeComRobotMessage 发送企微机器人消息
func SendWeComRobotMessage(p *RobotWebhookParameters, data *sonar.WebhookData) error {
	buildArgs := MessageBuildArgs{
		WebhookData: data,
		SonarToken:  p.SonarToken,
		NewCode:     p.NewCode,
	}
	message, err := buildMessageContent(buildArgs)
	if err != nil {
		return err
	}
	client := wecom.NewRobotClient(p.Key)
	return client.SendTextMessage(message)
}

// SendWeComMessage 发送企微应用消息
func SendWeComMessage(p *WeComMessageParameters, data *sonar.WebhookData) error {
	buildArgs := MessageBuildArgs{
		WebhookData: data,
		SonarToken:  p.SonarToken,
		NewCode:     p.NewCode,
	}
	text, err := buildMessageContent(buildArgs)
	if err != nil {
		return err
	}
	email := strings.TrimSpace(data.Properties.SonarAnalysisUserEmail)
	if email == "" {
		return errors.New("properties [sonar.analysis.userEmail] can't be empty")
	}
	userId := strings.Split(email, "@")[0]

	config := p.WeComConfig
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
func FeiShuRobotMessage(p *RobotWebhookParameters, data *sonar.WebhookData) error {
	buildArgs := MessageBuildArgs{
		WebhookData: data,
		SonarToken:  p.SonarToken,
		NewCode:     p.NewCode,
	}
	message, err := buildMessageContent(buildArgs)
	if err != nil {
		return err
	}
	client := feishu.NewRobotClient(p.Key)
	return client.SendTextMessage(message)
}

// SendFeiShuMessage 发送飞书应用消息
func SendFeiShuMessage(p *FeiShuMessageParameters, data *sonar.WebhookData) error {
	buildArgs := MessageBuildArgs{
		WebhookData: data,
		SonarToken:  p.SonarToken,
		NewCode:     p.NewCode,
	}
	text, err := buildMessageContent(buildArgs)
	if err != nil {
		return err
	}
	email := strings.TrimSpace(data.Properties.SonarAnalysisUserEmail)
	if email == "" {
		return errors.New("properties [sonar.analysis.userEmail] can't be empty")
	}

	config := p.FeiShuConfig
	client := feishu.Client{AppId: config.AppId, AppSecret: config.AppSecret}
	contentObj := feishu.MessageContentText{Text: text}
	contentBytes, err := json.Marshal(contentObj)
	if err != nil {
		return err
	}
	request := feishu.SendMessageRequest{
		ReceiveId: email,
		MsgType:   "text",
		Content:   string(contentBytes),
	}
	_, err = client.SendMessage("email", request)
	if err != nil {
		return err
	}
	return nil
}

func buildMessageContent(p MessageBuildArgs) (string, error) {
	data := p.WebhookData
	sonarToken := p.SonarToken
	newCode := p.NewCode

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
		measuresData := buildMeasuresData(measuresComponent, newCode)
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

func buildMeasuresData(measuresComponent *sonar.MeasuresComponent, newCode bool) *MeasuresData {
	measuresData := new(MeasuresData)
	if newCode {
		for _, measure := range measuresComponent.Measures {
			switch measure.Metric {
			case "new_bugs":
				measuresData.Bugs = measure.Period.Value
			case "new_vulnerabilities":
				measuresData.Vulnerabilities = measure.Period.Value
			case "new_code_smells":
				measuresData.CodeSmells = measure.Value
			case "new_duplicated_lines_density":
				measuresData.DuplicatedLinesDensity = measure.Period.Value
			}
		}
	} else {
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
	}
	if measuresData.Bugs == "" {
		measuresData.Bugs = "0"
	}
	if measuresData.Vulnerabilities == "" {
		measuresData.Vulnerabilities = "0"
	}
	if measuresData.CodeSmells == "" {
		measuresData.CodeSmells = "0"
	}
	if measuresData.DuplicatedLinesDensity == "" {
		measuresData.DuplicatedLinesDensity = "0.0"
	}
	return measuresData
}

func getMeasuresComponent(data *sonar.WebhookData, sonarToken string) *sonar.MeasuresComponent {
	if sonarToken == "" {
		return nil
	}

	client := sonar.SonarClient{
		ServerUrl: data.ServerURL,
		Token:     sonarToken,
	}
	metricKeys := "bugs,reliability_rating,vulnerabilities,security_rating,code_smells,sqale_rating,duplicated_lines_density,coverage,new_bugs,new_reliability_rating,new_vulnerabilities,new_security_rating,new_code_smells,new_duplicated_lines_density,new_coverage"
	response, err := client.GetMeasuresComponent(data.Project.Key, data.Branch.Name, metricKeys)

	if err != nil {
		log.Printf("get measures response failed: %v", err)
		return nil
	}
	return &response.Component
}
