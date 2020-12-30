package main

import (
	"encoding/xml"
	"fmt"

	"time"

	"wxmsgserver/log"
	lutil "wxmsgserver/util"

	"github.com/gin-gonic/gin"
)

func main() {

	gin.SetMode(gin.DebugMode)
	gin.ForceConsoleColor()
	client := NewClient()
	router := gin.Default()
	router.Use(log.LoggerToFile())
	router.POST("/wx", WxTextMsgReply)
	router.GET("/wx", WXCheckSignature)
	router.POST("/lane/sendtemplate", client.SendTempleMsg)
	log.Logger().Fatal(router.Run(":80"))
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
		log.Logger().Infof("接收消息类型不合法 msgtype:%s", msgType)
	}

}
func WXCheckSignature(c *gin.Context) {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")

	ok := lutil.CheckSignature(signature, timestamp, nonce, WXToken)
	if !ok {
		log.Logger().Infof("微信公众号接入校验失败!")
		return
	}

	log.Logger().Infof("微信公众号接入校验成功!")
	_, _ = c.Writer.WriteString(echostr)
}

func WxTextMsgReply(c *gin.Context) {
	var textMsg WxTextMsgReceive
	err := c.ShouldBindXML(&textMsg)
	if err != nil {
		log.Logger().Infof("消息接收 xml数据包解析失败:%s", err.Error())
		return
	}
	log.Logger().Infof("消息接收 收到消息 消息类型：%s,消息内容:%s,发送者opened:%s,接收者微信名:%s", textMsg.MsgType, textMsg.Content, textMsg.FromUserName, textMsg.ToUserName)
	resp := WxRepTextMsg{
		ToUserName:   textMsg.FromUserName,
		FromUserName: textMsg.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      fmt.Sprintf("欢迎来到Evildoer的世界，现在时间：%s", time.Now().Format("2006-01-02 15:04:05")),
	}
	bs, err := xml.Marshal(&resp)
	if err != nil {
		log.Logger().Infof("消息回复 失败err:%s", err.Error())
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
