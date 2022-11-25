package oauth

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	config *Config
}

type Config struct {
	ServerUrl    string `yaml:"serverUrl"`
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
}

func NewClient(config *Config) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) httpPost(url, data string, header http.Header) ([]byte, error) {
	method := "POST"

	payload := strings.NewReader(data)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	if header == nil {
		req.Header = http.Header{}
	} else {
		req.Header = header.Clone()
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Client) RequestAccessToken(code string) (*AccessTokenRespond, error) {
	resp, err := c.httpPost(c.config.ServerUrl+"/plugin.php?id=codfrm_oauth2:server&op=access_token", "client_id="+c.config.ClientID+"&client_secret="+
		c.config.ClientSecret+"&code="+code, nil)
	if err != nil {
		return nil, err
	}
	ret := &AccessTokenRespond{}
	if err := json.Unmarshal(resp, ret); err != nil {
		return nil, err
	}
	if ret.Code != 0 {
		return nil, ret
	}
	return ret, nil
}

func (c *Client) RequestUser(access_token string) (*UserRespond, error) {
	resp, err := c.httpPost(c.config.ServerUrl+"/plugin.php?id=codfrm_oauth2:server&op=user", "access_token="+access_token, nil)
	if err != nil {
		return nil, err
	}
	ret := &UserRespond{}
	if err := json.Unmarshal(resp, ret); err != nil {
		return nil, err
	}
	if ret.Code != 0 {
		return nil, ret
	}
	return ret, nil
}
