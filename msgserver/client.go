package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"wxmsgserver/log"

	"github.com/gin-gonic/gin"
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

func NewClient() *Client {

	return &Client{
		HTTPClient: http.DefaultClient,
		AppID:      "wx8aa615f8c0728a31",
		AppSecret:  "76f18821982615844afd251c08ed74f1",
		GrantType:  "client_credential",
		WxHost:     "https://api.weixin.qq.com/cgi-bin/",
	}
}

func (c *Client) SendTempleMsg(context *gin.Context) {
	var sendmsg SendTemplateMsgReq
	err := context.ShouldBindJSON(&sendmsg)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errcode": 1,
			"errmsg":  "参数不合法",
		})
	}
	log.Logger().Infof("sendmsg :%+v", sendmsg)
	bs, err := json.Marshal(&sendmsg)
	if err != nil {
		log.Logger().Infof("SendTempleMsg failed  marshal failed err :%s", err.Error())
		context.JSON(http.StatusOK, gin.H{
			"errcode": 1,
			"errmsg":  "参数序列化失败",
		})
	}
	var accesstoken string = c.GetAccessToken()
	if accesstoken == "" {
		log.Logger().Errorf("SendTempleMsg accesstoken invaild")
		context.JSON(http.StatusOK, gin.H{
			"errcode": 1,
			"errmsg":  "accesstoken失效",
		})
	}
	log.Logger().Infof("sendtemplate accesstoken:%s", accesstoken)
	url := fmt.Sprintf("%smessage/template/send?access_token=%s", c.WxHost, accesstoken)
	log.Logger().Infof("sendtemplemsg url:%s,body:%s", url, string(bs))
	httpreq, err := http.NewRequest("POST", url, bytes.NewBuffer(bs))
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"errcode": 1,
			"errmsg":  fmt.Sprintf("发送模版信息失败原因1：%s", err.Error()),
		})
	}
	httpreq.Header.Set("Content-Type", "application/json")

	httpresp, err := c.HTTPClient.Do(httpreq)
	if httpresp != nil {
		defer httpresp.Body.Close()
	}
	if err != nil {
		log.Logger().Errorf("GetAccessToken response  error %s", err.Error())
		context.JSON(http.StatusOK, gin.H{
			"errcode": 1,
			"errmsg":  fmt.Sprintf("发送模版信息失败原因2：%s", err.Error()),
		})
	}
	if httpresp.StatusCode != http.StatusOK {
		log.Logger().Errorf("GetAccessToken failed errcode:%d", httpresp.StatusCode)
		context.JSON(httpresp.StatusCode, gin.H{
			"errcode": httpresp.StatusCode,
			"errmsg":  fmt.Sprintf("发送模版信息失败原因3：%s", err.Error()),
		})
	}
	response := &SendTemplateMsgResp{}
	decoder := json.NewDecoder(httpresp.Body)
	if err := decoder.Decode(&response); err != nil && err != io.EOF {
		log.Logger().Errorf("GetAccessToken response decode error %s", err.Error())
		context.JSON(http.StatusOK, gin.H{
			"errcode": httpresp.StatusCode,
			"errmsg":  fmt.Sprintf("发送模版信息失败原因4：%s", err.Error()),
		})
	}
	log.Logger().Infof("sendtemplemsg response :%+v", response)
	if response.ErrCode == 0 {
		context.JSON(http.StatusOK, &response)
	} else {

		context.JSON(http.StatusOK, gin.H{
			"errcode": 1,
			"errmsg":  fmt.Sprintf("发送模版信息失败原因5：%s", response.ErrMsg),
		})
	}
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
	httpreq, err := http.NewRequest("POST", url, bytes.NewBuffer(bs))
	if err != nil {
		return ""
	}
	httpreq.Header.Set("Content-Type", "application/json")

	httpresp, err := c.HTTPClient.Do(httpreq)
	if httpresp != nil {
		defer httpresp.Body.Close()
	}
	if err != nil {
		log.Logger().Errorf("GetAccessToken response  error %s", err.Error())
		return ""
	}
	if httpresp.StatusCode != http.StatusOK {
		log.Logger().Errorf("GetAccessToken failed errcode:%d", httpresp.StatusCode)
		return ""
	}
	response := &AccessTokenResp{}
	decoder := json.NewDecoder(httpresp.Body)
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
