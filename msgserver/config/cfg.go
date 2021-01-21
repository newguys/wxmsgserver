package config

import (
	"fmt"
	"os"
	"encoding/json"
)

type Config struct {
	RabbitMQ struct{
		Url string `json:"url"`
		OrderTopic string `json:"order_topic"`
		RedPacketTopic string `json:"redpacket_topic"`
		AwardTopic string `json:"award_topic"`
	}`json:"rabbitmq"`

	WxCfg struct{
		Host string `json:"host"`
		AppID string `json:"app_id"`
		AppSecret string `json:"app_secret"`
		GrantType string `json:"grant_type"`
	}`json:"wxcfg"`

	Template struct{

		Order struct{
			TemplateId string `json:"template_id"`
			TemplateName string `json:"template_name"`
		}`json:"order"`

		RedPacket struct{
			TemplateId string `json:"template_id"`
			TemplateName string `json:"template_name"`
		}`json:"redpacket"`

		Award struct{
			TemplateId string `json:"template_id"`
			TemplateName string `json:"template_name"`
		}`json:"award"`

	}`json:"template"`
	DefaultOpenID string `json:"default_openid"`
}

//ParseConfig 解析config
func ParseConfig(filepath string, out interface{}) {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(out)
	if err != nil {
		panic(err)
	}
	fmt.Printf("config:%+v",out)
}
