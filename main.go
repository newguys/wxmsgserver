package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

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

func main() {
	router := gin.Default()
	router.POST("/wx", WXMsgReceive)
	log.Fatalln(router.Run(":80"))
}

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

func WXMsgReceive(c *gin.Context) {

	msgType := c.Params.ByName("MsgType")
	switch WxMsgType(msgType) {
	case WxMsgTypeText:
		WxTextMsgReply(c)
	case WxMsgTypeIMage:
		WxIMageMsgReply(c)
	case WxMsgTypeVoice:
		WxVoiceMsgReply(c)
	case WxMsgTypeVideo:
		WxVideoMsgReply(c)
	case WxMsgTypeShortVideo:
		WxShortVideoMsgReply(c)
	case WxMsgTypeLoaction:
		WxLocationMsgReply(c)
	case WxMsgTypeLink:
		WxLinkMsgReply(c)
	default:
		log.Printf("接收消息类型不合法 msgtype:%s", msgType)
	}

}

func WxTextMsgReply(c *gin.Context) {
	var textMsg WxTextMsgReceive
	err := c.ShouldBindXML(&textMsg)
	if err != nil {
		log.Printf("消息接收 xml数据包解析失败:%s", err.Error())
		return
	}
	log.Printf("消息接收 收到消息 消息类型：%s,消息内容:%s,发送者opened:%s,接收者微信名:%s", textMsg.MsgType, textMsg.Content, textMsg.FromUserName, textMsg.ToUserName)
	resp := WxRepTextMsg{
		ToUserName:   textMsg.FromUserName,
		FromUserName: textMsg.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      fmt.Sprintf("欢迎来到Evildoer的世界，现在时间：%s", time.Now().Format("2016-01-02 15:04:05")),
	}
	bs, err := xml.Marshal(&resp)
	if err != nil {
		log.Printf("消息回复 失败err:%s", err.Error())
		return
	}
	_, _ = c.Writer.Write(bs)

}

func WxIMageMsgReply(c *gin.Context) {
}

func WxVoiceMsgReply(c *gin.Context) {
}

func WxVideoMsgReply(c *gin.Context) {
}

func WxShortVideoMsgReply(c *gin.Context) {

}

func WxLocationMsgReply(c *gin.Context) {
}

func WxLinkMsgReply(c *gin.Context) {

}
