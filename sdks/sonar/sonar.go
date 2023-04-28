package sonar

import (
	"fmt"
	"net/http"
	"sonar-webhook/sdks/util"
)

type SonarClient struct {
	ServerUrl string // 服务地址
	Token     string // 访问token
}

// GetMeasuresComponent 获取项目指标
func (client *SonarClient) GetMeasuresComponent(component string, branch string, metricKeys string) (*MeasuresComponentResponse, error) {
	path := "/api/measures/component?component=%s&branch=%s&metricKeys=%s"
	url := client.ServerUrl + fmt.Sprintf(path, component, branch, metricKeys)
	response := new(MeasuresComponentResponse)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(client.Token, "")
	req.Header.Set("Host", client.ServerUrl)

	_, err = util.HttpRequest(req, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
