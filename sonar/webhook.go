package sonar

const (
	StatusSuccess = "SUCCESS"
	StatusError   = "ERROR"
	StatusOk      = "OK"
)

const (
	MetricReliabilityRating      = "reliability_rating"       // 可靠性等级
	MetricSecurityRating         = "security_rating"          // 安全性等级
	MetricSqaleRating            = "sqale_rating"             // sqale评级
	MetricDuplicatedLinesDensity = "duplicated_lines_density" // 代码行重复密度
	MetricBlockerViolations      = "blocker_violations"       // 阻塞的违规数
	MetricCriticalViolations     = "critical_violations"      // 严重的违规数
	MetricMajorViolations        = "major_violations"         // 主要的违规数
)

// WebhookData webhook请求参数
type WebhookData struct {
	ServerURL   string      `json:"serverUrl"`   // 服务地址
	TaskID      string      `json:"taskId"`      // 任务Id
	Status      string      `json:"status"`      // 任务状态
	AnalysedAt  string      `json:"analysedAt"`  // 分析时间
	Revision    string      `json:"revision"`    // 修订号
	Project     Project     `json:"project"`     // 项目信息
	Branch      Branch      `json:"branch"`      // 分支信息
	QualityGate QualityGate `json:"qualityGate"` // 质量阈值信息
	Properties  Properties  `json:"properties"`  // 附加属性
}

type Project struct {
	Key  string `json:"key"`  // 项目标识
	Name string `json:"name"` // 项目名称
	URL  string `json:"url"`  // 项目访问url
}

type Branch struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	IsMain bool   `json:"isMain"`
	URL    string `json:"url"`
}

type QualityGate struct {
	Name       string       `json:"name"`       // 质量配置名称
	Status     string       `json:"status"`     // 质量阈值是否通过
	Conditions []Conditions `json:"conditions"` // 质量配置条件
}

type Conditions struct {
	Metric         string `json:"metric"`         // 指标名称
	Operator       string `json:"operator"`       // 操作符
	ErrorThreshold string `json:"errorThreshold"` // 阈值
	Value          string `json:"value"`          // 指标值
	Status         string `json:"status"`         // 状态
}

type Properties struct {
	SonarAnalysisUserEmail string `json:"sonar.analysis.userEmail"` // 用户邮件
}
