package main

import "encoding/xml"

const WXToken string = "Evildoer"

type Configs struct {
	APPID     string
	APPSecret string
	GrantType string
}

//AccessTokenReq 获取accesstoken请求
type AccessTokenReq struct {
	APPID     string `json:"appid"`
	APPSecret string `json:"secret"`
	GrantType string `json:"grant_type"`
}

//AccessTokenResp 获取accesstoken返回
type AccessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpriesIn   int64  `json:"expires_in"`
	ErrCode     int32  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

type MiniProgram struct {
	AppID    string `json:"appid"`
	PagePath string `json:"pagepath,omitempty"`
}

type Data struct {
}

//SendTemplateMsgReq
type SendTemplateMsgReq struct {
	Touser      string      `json:"touser"`
	TemplateID  string      `json:"template_id"`
	URL         string      `json:"url,omitempty"`
	MiniProgram MiniProgram `json:"miniprogram,omitempty"`
	Data        Data        `json:"data"`
}

type SendTemplateMsgResp struct {
	ErrCode int32  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	MsgId   int64  `json:"msgid"`
}

type WxMsgType string

const (
	WxMsgTypeText       WxMsgType = "text"
	WxMsgTypeIMage      WxMsgType = "image"
	WxMsgTypeVoice      WxMsgType = "voice"
	WxMsgTypeVideo      WxMsgType = "video"
	WxMsgTypeShortVideo WxMsgType = "shortvideo"
	WxMsgTypeLoaction   WxMsgType = "location"
	WxMsgTypeLink       WxMsgType = "link"
)

//文本消息
type WxTextMsgReceive struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`   //接收方微信号
	FromUserName string   `xml:"FromUserName"` //发送方微信号，若为普通用户，则是一个OpenID
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
	MsgId        int64    `xml:"MsgId"`
}

type WxRepTextMsg struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`   //接收方微信号
	FromUserName string   `xml:"FromUserName"` //发送方微信号，若为普通用户，则是一个OpenID
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
}

//图片消息
type WxImageMsgReceive struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`   //接收方微信号
	FromUserName string   `xml:"FromUserName"` //发送方微信号，若为普通用户，则是一个OpenID
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
	MsgId        int64    `xml:"MsgId"`
}

type WxRepImageMsg struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`   //接收方微信号
	FromUserName string   `xml:"FromUserName"` //发送方微信号，若为普通用户，则是一个OpenID
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
}
