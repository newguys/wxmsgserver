package main

import (
	"encoding/xml"
)

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
type SubData struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}

type Data struct {
	First SubData `json:"first"`
	Keyword1 SubData `json:"keyword1,omitempty"`
	Keyword2 SubData `json:"keyword2,omitempty"`
	Keyword3 SubData `json:"keyword3,omitempty"`
	Keyword4 SubData `json:"keyword4,omitempty"`
	Keyword5 SubData `json:"keyword5,omitempty"`
	Keyword6 SubData `json:"keyword6,omitempty"`
	Keyword7 SubData `json:"keyword7,omitempty"`
	Keyword8 SubData `json:"keyword8,omitempty"`
	Keyword9 SubData `json:"keyword9,omitempty"`
	Remark SubData `json:"remark"`
}

//SendTemplateMsgReq
type SendTemplateMsgReq struct {
	Touser      string      `json:"touser"`
	TemplateID  string      `json:"template_id"`
	URL         string      `json:"url,omitempty"`
	MiniProgram MiniProgram `json:"miniprogram,omitempty"`
	Data        Data        `json:"data"`
}

type OrderNTSt struct {
	OpenId string `json:"openId"`
	OrderCode string `json:"orderCode"`
	Price string `json:"price"`
	StoreName string `json:"storeName"`
	StoreId int64 `json:"storeId"`
	ActivityId int64 `json:"activityId"`
	ActivityName string `json:"activityName"`
	CommodityName string `json:"commodityName"`
}

type RedPackNTSt struct {
	OpenId string `json:"openId"`
	StoreName string `json:"storeName"`
	StoreId int64 `json:"storeId"`
	Amount string `json:"amount"`
}

type AwardNTSt struct {
	OpenId string `json:"openId"`
	StoreName string `json:"storeName"`
	StoreId int64 `json:"storeId"`
	ActivityId int64 `json:"activityId"`
	ActivityName string `json:"activityName"`
	PrizeName string `json:"prizeName"`
}