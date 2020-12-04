package httpmock

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	ErrInternalFailure = errors.New("we did something wrong")
	ErrUpstreamFailure = errors.New("someone else did something wrong")
)

type FooClient struct {
	url      string
	username string
	password string
	client   *http.Client
}

func NewFooClient(url, username, password string) *FooClient {
	return &FooClient{
		url:      url,
		username: username,
		password: password,
		client:   &http.Client{},
	}
}

func (fc *FooClient) DoStuff(param string) (string, error) {

	reqBody := struct {
		Param string `json:"param"`
	}{
		Param: param,
	}
	jReqBody, err := json.Marshal(reqBody)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/path", fc.url), bytes.NewBuffer(jReqBody))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInternalFailure, err)
	}
	req.SetBasicAuth(fc.username, fc.password)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := fc.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrUpstreamFailure, err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w, server responded with status %d", ErrUpstreamFailure, resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrUpstreamFailure, err)
	}

	var expResp struct {
		Result string `json:"result"`
	}
	err = json.Unmarshal(body, &expResp)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrUpstreamFailure, err)
	}

	return expResp.Result, nil
}
