package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"wxmsgserver/log"
)


type Client struct {
	HTTPClient    *http.Client
	AccessToken   string
	AppID         string
	AppSecret     string
	AccessTokenTs int64 //token 失效时间
	GrantType     string
	WxHost        string
}

func NewClient(appId,appSecret,grantType,host string) *Client {
	fmt.Print("new http client\n")
	return &Client{
		HTTPClient: http.DefaultClient,
		AppID:      appId,
		AppSecret:  appSecret,
		GrantType:  grantType,
		WxHost:     host,
	}
}
func (c *Client)SendTempleMsgV2(bs []byte) (err error,response SendTemplateMsgResp) {
	var accessToken  = c.GetAccessToken()

	if accessToken == "" {
		log.Logger().Errorf("SendTempleMsg accessToken invaild")
		response.ErrMsg = "accessToken失效"
		response.ErrCode = 1
		return fmt.Errorf(response.ErrMsg),response
	}
	log.Logger().Infof("sendtemplate accessToken:%s", accessToken)
	url := fmt.Sprintf("%smessage/template/send?access_token=%s", c.WxHost, accessToken)
	log.Logger().Infof("sendtemplemsg url:%s,body:%s", url, string(bs))
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(bs))
	if err != nil {
		response.ErrMsg = fmt.Sprintf("发送模版信息失败原因1：%s", err.Error())
		response.ErrCode = 1
		return fmt.Errorf(response.ErrMsg),response
	}
	if httpReq == nil {
		response.ErrMsg = "send template req is nil"
		response.ErrCode = 1
		return fmt.Errorf(response.ErrMsg),response
	}
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.HTTPClient.Do(httpReq)
	if httpResp != nil {
		defer httpResp.Body.Close()
	}else {
		response.ErrMsg ="send template resp is nil"
		response.ErrCode = 1
		return fmt.Errorf(response.ErrMsg),response
	}
	if err != nil {
		log.Logger().Errorf("GetAccessToken response  error %s", err.Error())
		response.ErrMsg = fmt.Sprintf("发送模版信息失败原因2：%s", err.Error())
		response.ErrCode = 1
		return fmt.Errorf(response.ErrMsg),response
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Logger().Errorf("GetAccessToken failed errcode:%d", httpResp.StatusCode)
		response.ErrMsg = fmt.Sprintf("发送模版信息失败原因3：%d", httpResp.StatusCode)
		response.ErrCode = 1
		return fmt.Errorf(response.ErrMsg),response
	}
	decoder := json.NewDecoder(httpResp.Body)
	if err := decoder.Decode(&response); err != nil && err != io.EOF {
		log.Logger().Errorf("GetAccessToken response decode error %s", err.Error())
		response.ErrMsg = fmt.Sprintf("发送模版信息失败原因4：%s", err.Error())
		response.ErrCode = 1
		return fmt.Errorf(response.ErrMsg),response
	}
	log.Logger().Infof("sendtemplemsg response :%+v", response)
	if response.ErrCode != 0 {
		return fmt.Errorf(response.ErrMsg),response
	}
	return nil,response
}

func (c *Client) GetAccessToken() string {
	log.Logger().Infof("AccessTokenTs:%d,accesstoken:%s", c.AccessTokenTs, c.AccessToken)
	if c.AccessToken != "" && (c.AccessTokenTs-time.Now().Unix() > 7200) {
		return c.AccessToken
	}
	params := &AccessTokenReq{
		APPID:     c.AppID,
		APPSecret: c.AppSecret,
		GrantType: c.GrantType,
	}
	bs, err := json.Marshal(params)
	if err != nil {
		log.Logger().Errorf("getaccesstokken failed  marshal failed err :%s", err.Error())
		return ""
	}
	url := fmt.Sprintf("%stoken?grant_type=%s&appid=%s&secret=%s", c.WxHost, c.GrantType, c.AppID, c.AppSecret)
	log.Logger().Infof("url:%s", url)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(bs))
	if err != nil {
		return ""
	}
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.HTTPClient.Do(httpReq)
	if httpResp != nil {
		defer httpResp.Body.Close()
	}
	if err != nil {
		log.Logger().Errorf("GetAccessToken response  error %s", err.Error())
		return ""
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Logger().Errorf("GetAccessToken failed errcode:%d", httpResp.StatusCode)
		return ""
	}
	response := &AccessTokenResp{}
	decoder := json.NewDecoder(httpResp.Body)
	if err := decoder.Decode(&response); err != nil && err != io.EOF {
		log.Logger().Errorf("GetAccessToken response decode error %s", err.Error())
		return ""
	}
	if response.ErrCode == 0 {
		c.AccessToken = response.AccessToken
		c.AccessTokenTs = time.Now().Unix() + response.ExpriesIn
		return c.AccessToken
	}
	log.Logger().Errorf("GetAccessToken failed errcode:%d errmsg:%s", response.ErrCode, response.ErrMsg)
	return ""
}

func (c *Client) GetTempleList() {

}
