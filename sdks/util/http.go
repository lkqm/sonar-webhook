package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func HttpPostJson(url string, data interface{}, headers map[string]string) (string, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataBytes))
	if err != nil {
		return "", err
	}
	if headers != nil {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func HttpPostEntity(url string, data interface{}, returnData interface{}, headers map[string]string) (string, error) {
	content, err := HttpPostJson(url, data, headers)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal([]byte(content), returnData)
	if err != nil {
		return content, err
	}

	return content, nil
}

func HttpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func HttpGetJson(url string, returnData interface{}) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(bytes, returnData)
	if err != nil {
		return string(bytes), err
	}

	return string(bytes), nil
}

func HttpRequest(request *http.Request, returnData interface{}) (string, error) {
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if returnData != nil {
		err = json.Unmarshal(bytes, returnData)
		if err != nil {
			return string(bytes), err
		}
	}

	return string(bytes), nil
}
