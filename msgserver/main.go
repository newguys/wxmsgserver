package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"net/http"
	"wxmsgserver/msgserver/config"

	"time"

	"wxmsgserver/log"
	lutil "wxmsgserver/util"

	"github.com/gin-gonic/gin"
)

var cfgfile = flag.String("config_path","msgserver.json","config file")

func main() {
	log.LoggerToFile()
	flag.Parse()
	fmt.Printf("cfgfile:%s\n",*cfgfile)
/*	gin.SetMode(gin.DebugMode)
	gin.ForceConsoleColor()
	router := gin.Default()
	router.Use(log.LoggerToFile())*/
	var cfg config.Config

	config.ParseConfig(*cfgfile,&cfg)
	basehandler := newBaseHandler(cfg)
	//router.POST("/wx", WxTextMsgReply)
	//router.GET("/wx", WXCheckSignature)
	//router.POST("/lane/sendtemplate", basehandler.SendTemplateMsg)
	//log.Logger().Fatal(router.Run(":80"))

	basehandler.ConsumeOrderMqMsg()
	//basehandler.ConsumeRedpackMqMsg()
	//basehandler.ConsumeAwardMqMsg()

}

type BaseHandler struct {
	NetClient *Client
	Config config.Config
}

func newBaseHandler(cfg config.Config) *BaseHandler {
	return &BaseHandler{
		NetClient: NewClient(cfg.WxCfg.AppID,cfg.WxCfg.AppSecret,cfg.WxCfg.GrantType,cfg.WxCfg.Host),
		Config: cfg,
	}
}

func failOnError(err error,msg string)  {
	if err != nil{
		fmt.Printf("%s:%s\n",msg,err.Error())
		log.Logger().Errorf("%s:%s",msg,err.Error())
	}
}

func (baseHandler *BaseHandler) ConsumeOrderMqMsg(){
	conn,err := amqp.Dial("amqp://admin:admin@localhost:5672/test")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		baseHandler.Config.RabbitMQ.OrderTopic, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")
	err = ch.QueueBind(
		q.Name,       // queue name
		baseHandler.Config.RabbitMQ.OrderTopic,            // routing key
		baseHandler.Config.RabbitMQ.OrderTopic, // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range msgs{
			baseHandler.ParseAndSendOrderTemplateMsg(d)
		}
	}()
	<-forever
	fmt.Println("end")
}

func (baseHandler *BaseHandler)ParseAndSendOrderTemplateMsg(d amqp.Delivery)  {
	fmt.Printf("receive d:%s",d.Body)
	order := &OrderNTSt{}
	err := json.Unmarshal(d.Body,order)
	if err != nil {
		failOnError(err,fmt.Sprintf("rabbitmq order msg :%s unmarshal failed err:%s",string(d.Body)))
		err = d.Ack(true)
		if err !=nil {
			failOnError(err,"rabbitmq order msg ack failed ")
		}
		return
	}

	err,bs := baseHandler.translateOrderMsgData(order)
	if err != nil{
		failOnError(err,"rabbit order msg translate template msg failed err")
	}else{
		err,resp :=baseHandler.NetClient.SendTempleMsgV2(bs)
		if err !=nil {
			failOnError(err,fmt.Sprintf("send to openid:%s order template msg failed err:",order.OpenId))
			//log.Logger().Errorf("send to openid:%s order template msg failed err:%s",order.OpenId,err.Error())
		}
		if resp.ErrCode != 0{
			failOnError(err,fmt.Sprintf("send to openid:%s order template msg response failed err:",order.OpenId))
			//log.Logger().Errorf("send to openid:%s order template msg response failed err:%s",order.OpenId,resp.ErrMsg)
		}
	}
	err = d.Ack(true)
	if err !=nil {
		failOnError(err,"rabbitmq order msg ack failed")
		//log.Logger().Errorf("rabbitmq order msg ack failed err:%s",err.Error())
	}
}

func (baseHandler *BaseHandler) ConsumeRedpackMqMsg(){
	conn,err := amqp.Dial("amqp://admin:admin@localhost:5672/test")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		baseHandler.Config.RabbitMQ.OrderTopic, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")
	err = ch.QueueBind(
		q.Name,       // queue name
		baseHandler.Config.RabbitMQ.OrderTopic,            // routing key
		baseHandler.Config.RabbitMQ.OrderTopic, // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	forever := make(chan bool)
	go func() {
		for d := range msgs{
			red := &RedPackNTSt{}
			err := json.Unmarshal(d.Body,red)
			if err != nil {
				log.Logger().Error("rabbitmq red packet msg :%s unmarshal failed err:%s",string(d.Body),err.Error())
				err = d.Ack(true)
				if err !=nil {
					log.Logger().Errorf("rabbitmq red packet msg ack failed err:%s",err.Error())
				}
				continue
			}
			err,bs := baseHandler.translateRedMsgData(red)
			if err != nil{
				log.Logger().Errorf("rabbit red packet msg translate template msg failed err")
			}else{
				err,resp := baseHandler.NetClient.SendTempleMsgV2(bs)
				if err !=nil {
					log.Logger().Errorf("send to openid:%s red template msg failed err:%s",red.OpenId,err.Error())
				}
				if resp.ErrCode != 0{
					log.Logger().Errorf("send to openid:%s red template msg response failed err:%s",red.OpenId,resp.ErrMsg)
				}
			}
			err = d.Ack(true)
			if err !=nil {
				log.Logger().Errorf("rabbitmq red msg ack failed err:%s",err.Error())
			}
		}
	}()
	<-forever
}


func (baseHandler *BaseHandler) ConsumeAwardMqMsg(){
	conn,err := amqp.Dial("amqp://admin:admin@localhost:5672/test")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		baseHandler.Config.RabbitMQ.OrderTopic, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")
	err = ch.QueueBind(
		q.Name,       // queue name
		baseHandler.Config.RabbitMQ.OrderTopic,            // routing key
		baseHandler.Config.RabbitMQ.OrderTopic, // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	forever := make(chan bool)
	go func() {
		for d := range msgs{
			award := &AwardNTSt{}
			err := json.Unmarshal(d.Body,award)
			if err != nil {
				log.Logger().Error("rabbitmq award msg :%s unmarshal failed err:%s",string(d.Body),err.Error())
				err = d.Ack(true)
				if err !=nil {
					log.Logger().Errorf("rabbitmq award  msg ack failed err:%s",err.Error())
				}
				continue
			}
			err,bs := baseHandler.translateAwardMsgData(award)
			if err != nil{
				log.Logger().Errorf("rabbit award msg translate template msg failed err")
			}else {
				err, resp := baseHandler.NetClient.SendTempleMsgV2(bs)
				if err != nil {
					log.Logger().Errorf("send to openid:%s red template msg failed err:%s", award.OpenId, err.Error())
				}
				if resp.ErrCode != 0 {
					log.Logger().Errorf("send to openid:%s red template msg response failed err:%s", award.OpenId, resp.ErrMsg)
				}
			}
			err = d.Ack(true)
			if err !=nil {
				log.Logger().Errorf("rabbitmq award msg ack failed err:%s",err.Error())
			}
		}
	}()
	<-forever
}


//func WXMsgReceive(c *gin.Context) {
//
//	msgType := c.Params.ByName("MsgType")
//	switch WxMsgType(msgType) {
//	case WxMsgTypeText:
//		WxTextMsgReply(c)
//	case WxMsgTypeIMage:
//		WxIMageMsgReply(c)
//	case WxMsgTypeVoice:
//		WxVoiceMsgReply(c)
//	case WxMsgTypeVideo:
//		WxVideoMsgReply(c)
//	case WxMsgTypeShortVideo:
//		WxShortVideoMsgReply(c)
//	case WxMsgTypeLoaction:
//		WxLocationMsgReply(c)
//	case WxMsgTypeLink:
//		WxLinkMsgReply(c)
//	default:
//		log.Logger().Infof("接收消息类型不合法 msgtype:%s", msgType)
//	}
//}
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
//
//func WxIMageMsgReply(c *gin.Context) {
//}
//
//func WxVoiceMsgReply(c *gin.Context) {
//}
//
//func WxVideoMsgReply(c *gin.Context) {
//}
//
//func WxShortVideoMsgReply(c *gin.Context) {
//
//}
//
//func WxLocationMsgReply(c *gin.Context) {
//}

//func WxLinkMsgReply(c *gin.Context) {
//
//}

func (baseHandler *BaseHandler)SendTemplateMsg(context *gin.Context)  {
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
	_ ,response := baseHandler.NetClient.SendTempleMsgV2(bs)
	context.JSON(http.StatusOK,response)
}

func (baseHandler *BaseHandler)translateOrderMsgData(msg *OrderNTSt) (error,[]byte) {
	keyword1 := msg.OrderCode
	keyword2 := msg.Price
	keyword3 := msg.StoreName
	keyword4 := msg.StoreId
	keyword5 := msg.ActivityId
	keyword6 := msg.ActivityName
	keyword7 := msg.CommodityName
	remark := "如有任何疑问欢迎再与客服联络"

	templateMsg := SendTemplateMsgReq{
		Touser: "OPENID",
		TemplateID: baseHandler.Config.Template.Order.TemplateId,
	}
	templateMsg.Data.First.Value = baseHandler.Config.Template.Order.TemplateName
	templateMsg.Data.Keyword1.Value = keyword1
	templateMsg.Data.Keyword2.Value = keyword2
	templateMsg.Data.Keyword3.Value = keyword3
	templateMsg.Data.Keyword4.Value = fmt.Sprintf("%d",keyword4)
	templateMsg.Data.Keyword5.Value = fmt.Sprintf("%d",keyword5)
	templateMsg.Data.Keyword6.Value = keyword6
	templateMsg.Data.Keyword7.Value = keyword7
	templateMsg.Data.Remark.Value = remark

	bs,err := json.Marshal(templateMsg)
	if err != nil {
		log.Logger().Errorf("translate template msg failed err :%s",err.Error())
		return err,[]byte{}
	}
	return nil,bs
}


func (baseHandler *BaseHandler)translateAwardMsgData(msg *AwardNTSt) (error,[]byte) {

	keyword3 := msg.StoreName
	keyword4 := msg.StoreId
	keyword5 := msg.ActivityId
	keyword6 := msg.ActivityName

	remark := "如有任何疑问欢迎再与客服联络"

	templateMsg := SendTemplateMsgReq{
		Touser: "OPENID",
		TemplateID: baseHandler.Config.Template.Award.TemplateId,
	}
	templateMsg.Data.First.Value = baseHandler.Config.Template.Award.TemplateName

	templateMsg.Data.Keyword3.Value = keyword3
	templateMsg.Data.Keyword4.Value = fmt.Sprintf("%d",keyword4)
	templateMsg.Data.Keyword5.Value = fmt.Sprintf("%d",keyword5)
	templateMsg.Data.Keyword6.Value = keyword6
	templateMsg.Data.Remark.Value = remark

	bs,err := json.Marshal(templateMsg)
	if err != nil {
		log.Logger().Errorf("translate template msg failed err :%s",err.Error())
		return err,[]byte{}
	}
	return nil,bs
}


func (baseHandler *BaseHandler)translateRedMsgData(msg *RedPackNTSt) (error,[]byte) {

	keyword3 := msg.StoreName
	keyword4 := msg.StoreId
	keyword9 := msg.Amount

	remark := "如有任何疑问欢迎再与客服联络"

	templateMsg := SendTemplateMsgReq{
		Touser: "OPENID",
		TemplateID: baseHandler.Config.Template.RedPacket.TemplateId,
	}
	templateMsg.Data.First.Value = baseHandler.Config.Template.RedPacket.TemplateName
	templateMsg.Data.Keyword3.Value = keyword3
	templateMsg.Data.Keyword4.Value = fmt.Sprintf("%d",keyword4)
	templateMsg.Data.Keyword9.Value = keyword9
	templateMsg.Data.Remark.Value = remark

	bs,err := json.Marshal(templateMsg)
	if err != nil {
		log.Logger().Errorf("translate template msg failed err :%s",err.Error())
		return err,[]byte{}
	}
	return nil,bs
}