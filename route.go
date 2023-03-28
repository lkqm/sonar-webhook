package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"sonar-webhook/sonar"
)

// 路由配置
func initRoute(g *gin.Engine) {
	g.POST("/dingtalk/robot", DingTalkRobotHandle)
	g.POST("/wecom/robot", WeComRobotHandle)
	g.POST("/wecom/message", WeComMessageHandle)
	g.POST("/feishu/robot", FeiShuRobotHandle)
	g.POST("/feishu/message", FeiShuMessageHandle)
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
	// 参数绑定
	p := new(RobotWebhookParameters)
	err := c.BindQuery(p)
	if err != nil {
		c.JSON(400, NewResultFail(400, "bind parameter failed: "+err.Error()))
		return
	}
	data := new(sonar.WebhookData)
	err = c.BindJSON(data)
	if err != nil {
		c.JSON(400, NewResultFail(400, "parse request body error: "+err.Error()))
		return
	}

	// 成功跳过
	if data.IsQualityGateSuccess() && p.SkipSuccess {
		log.Printf("skip send message when quality gate success")
		return
	}

	err = sender(p, data)
	if err != nil {
		c.JSON(500, NewResultFail(1, "request third failed: "+err.Error()))
		return
	}
	c.JSON(200, NewResultOkEmpty())
}

// WeComMessageHandle 企微应用消息通知处理方法
func WeComMessageHandle(c *gin.Context) {
	// 参数绑定
	p := new(WeComMessageParameters)
	err := c.BindQuery(p)
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

	// 成功跳过
	if data.IsQualityGateSuccess() && p.SkipSuccess {
		log.Printf("skip send message when quality gate success")
		return
	}

	err = SendWeComMessage(p, data)
	if err != nil {
		c.JSON(500, NewResultFail(1, "request third failed: "+err.Error()))
		return
	}
	c.JSON(200, NewResultOkEmpty())
}

// FeiShuMessageHandle 企微应用消息通知处理方法
func FeiShuMessageHandle(c *gin.Context) {
	// 参数绑定
	p := new(FeiShuMessageParameters)
	err := c.BindQuery(p)
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

	// 成功跳过
	if data.IsQualityGateSuccess() && p.SkipSuccess {
		log.Printf("skip send message when quality gate success")
		return
	}

	err = SendFeiShuMessage(p, data)
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
