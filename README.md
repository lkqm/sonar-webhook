# sonar-webhook
SonarQube扫描检测结果webhook回调适配服务，实现钉钉机器人、企业微信机器人、企微应用消息通知等。

- 钉钉机器人：`/dingtalk/robot?key=xxx`
- 企微机器人：`/wecom/robot?key=xxx`
- 企微应用消息: `/wecom/message?corpId=xxx&corpSecret=xxx&agentId=x`
- 飞书机器人: `/feishu/robot?key=xxx`

备注：可选参数`sonarToken`设置后通知消息会包含预览信息 ”检测结果: Bugs(0), 漏洞(0), 异味(2), 重复率(1.2%)“
