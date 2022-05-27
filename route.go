package main

import (
	"github.com/gin-gonic/gin"
	"sonar-webhook/sonar"
)

// 路由配置
func initRoute(g *gin.Engine) {
	g.POST("/dingtalk/robot", DingTalkRobotHandle)
	g.POST("/wecom/robot", WeComRobotHandle)
	g.POST("/wecom/message", WeComMessageHandle)
	g.POST("/feishu/robot", FeiShuRobotHandle)
}

// DingTalkRobotHandle 钉钉机器人通知处理方法
func DingTalkRobotHandle(c *gin.Context) {
	robotHandle(c, SendDingTalkRobotMessage)
}

// WeComRobotHandle 企微机器人处理方法
func WeComRobotHandle(c *gin.Context) {
	robotHandle(c, SendWeComRobotMessage)
}

// FeiShuRobotHandle 飞书器人通知处理方法
func FeiShuRobotHandle(c *gin.Context) {
	robotHandle(c, FeiShuRobotMessage)
}

func robotHandle(c *gin.Context, sender RobotMessageSender) {
	key := c.Query("key")
	sonarToken := c.Query("sonarToken")
	data := new(sonar.WebhookData)
	err := c.BindJSON(data)
	if err != nil {
		c.JSON(400, NewResultFail(400, "parse request body error: "+err.Error()))
		return
	}
	err = sender(key, sonarToken, data)
	if err != nil {
		c.JSON(500, NewResultFail(1, "request third failed: "+err.Error()))
		return
	}
	c.JSON(200, NewResultOkEmpty())
}

// WeComMessageHandle 企微应用消息通知处理方法
func WeComMessageHandle(c *gin.Context) {
	config := new(WeComConfig)
	sonarToken := c.Query("sonarToken")

	err := c.BindQuery(config)
	if err != nil {
		c.JSON(400, NewResultFail(400, "bind parameter failed: "+err.Error()))
		return
	}

	data := new(sonar.WebhookData)
	err = c.BindJSON(data)
	if err != nil {
		c.JSON(400, NewResultFail(400, "parse request body json error: "+err.Error()))
		return
	}

	err = SendWeComMessage(config, sonarToken, data)
	if err != nil {
		c.JSON(500, NewResultFail(1, "request third failed: "+err.Error()))
		return
	}
	c.JSON(200, NewResultOkEmpty())
}

// Result 接口层响应数据
type Result struct {
	Code int         `json:"code"` // 错误码
	Msg  string      `json:"msg"`  // 提示消息
	Data interface{} `json:"data"` // 数据
}

func NewResultFail(code int, msg string) *Result {
	return &Result{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

func NewResultOkEmpty() *Result {
	return &Result{
		Code: 0,
		Msg:  "ok",
		Data: nil,
	}
}
