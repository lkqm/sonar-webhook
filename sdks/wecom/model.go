package wecom

import "time"

// AccessToken 接口访问token信息
type AccessToken struct {
	Token    string // token
	ExpireAt int64  // 过期时间点
}

func (token AccessToken) IsExpired() bool {
	return time.Now().UnixMilli() > token.ExpireAt
}

// Result 企业微信响应结果
type Result struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (result Result) IsOk() bool {
	return result.ErrCode == 0
}

type GetAccessTokenResponse struct {
	Result
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type ListUserResponse struct {
	Result
	UserList []UserInfo `json:"userlist"`
}

type UserInfo struct {
	Userid           string   `json:"userid"`
	Name             string   `json:"name"`
	Department       []int    `json:"department"`
	Order            []int    `json:"order"`
	Position         string   `json:"position"`
	Mobile           string   `json:"mobile"`
	Gender           string   `json:"gender"`
	Email            string   `json:"email"`
	BizMail          string   `json:"biz_mail"`
	IsLeaderInDept   []int    `json:"is_leader_in_dept"`
	DirectLeader     []string `json:"direct_leader"`
	Avatar           string   `json:"avatar"`
	ThumbAvatar      string   `json:"thumb_avatar"`
	Telephone        string   `json:"telephone"`
	Alias            string   `json:"alias"`
	Status           int      `json:"status"`
	Address          string   `json:"address"`
	EnglishName      string   `json:"english_name"`
	OpenUserid       string   `json:"open_userid"`
	MainDepartment   int      `json:"main_department"`
	QrCode           string   `json:"qr_code"`
	ExternalPosition string   `json:"external_position"`
}

type Message struct {
	ToUser                 string    `json:"touser"`
	ToParty                string    `json:"toparty"`
	ToTag                  string    `json:"totag"`
	MsgType                string    `json:"msgtype"`
	AgentId                int       `json:"agentid"`
	Text                   *Text     `json:"text"`
	Markdown               *Markdown `json:"markdown"`
	EnableDuplicateCheck   int       `json:"enable_duplicate_check"`
	DuplicateCheckInterval int       `json:"duplicate_check_interval"`
}

type SendMessageResponse struct {
	Result
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
	MsgId        string `json:"msgid"`
	ResponseCode string `json:"response_code"`
}
