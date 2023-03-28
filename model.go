package main

import "sonar-webhook/sonar"

//-----------------------------------
// 请求参数
//-----------------------------------

// RobotWebhookParameters 机器人回掉请求参数
type RobotWebhookParameters struct {
	Key         string `json:"key" form:"key" binding:"required"` // 访问密钥
	SonarToken  string `json:"sonarToken" form:"sonarToken"`      // SonarQue认证密钥
	SkipSuccess bool   `json:"skipSuccess" form:"skipSuccess"`    // 是否跳过质量门成功
}

// WeComMessageParameters 机器人回掉请求参数
type WeComMessageParameters struct {
	WeComConfig
	SonarToken  string `json:"sonarToken" form:"sonarToken"`   // SonarQue认证密钥
	SkipSuccess bool   `json:"skipSuccess" form:"skipSuccess"` // 是否跳过质量门成功
}

type WeComConfig struct {
	CorpId     string `json:"corpId" form:"corpId" binding:"required"`         // 企业id
	CorpSecret string `json:"corpSecret" form:"corpSecret" binding:"required"` // 企业密钥
	AgentId    int    `json:"agentId" form:"agentId" binding:"required"`       // 应用id
}

// FeiShuMessageParameters 回调请求飞书消息参数
type FeiShuMessageParameters struct {
	FeiShuConfig
	SonarToken  string `json:"sonarToken" form:"sonarToken"`   // SonarQue认证密钥
	SkipSuccess bool   `json:"skipSuccess" form:"skipSuccess"` // 是否跳过质量门成功
}

type FeiShuConfig struct {
	AppId     string `json:"appId" form:"appId" binding:"required"`         // 应用ID
	AppSecret string `json:"appSecret" form:"appSecret" binding:"required"` // 应用密钥
}

//-----------------------------------
// 模版消息上下文
//-----------------------------------

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
