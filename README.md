# sonar-webhook

SonarQube扫描检测结果webhook回调适配服务，实现钉钉机器人、企业微信机器人、企微应用消息通知等。

- 钉钉机器人：`/dingtalk/robot?key=xxx`
- 企微机器人：`/wecom/robot?key=xxx`
- 企微应用消息: `/wecom/message?corpId=xxx&corpSecret=xxx&agentId=x`
- 飞书机器人: `/feishu/robot?key=xxx`

## 安装

- Docker: `docker run -p 8080:8080 lkqm/sonar-webhook`

## 其他参数

- sonarToken: SonarQube认证token，string类型，非必须，指定后会包含详细
- skipSuccess: 是否跳过质量门成功情况，boolean类型，非必须，默认false
