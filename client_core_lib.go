package fastpushclient

import (
	"fmt"
	"github.com/lgphp/go-fastpushclient/logger"
	"github.com/pkg/errors"
)

func NewFastLivePushClient(appInfo AppInfo) *Client {
	logger.Infow("开始初始化SDK....")
	client := buildClient()
	client.setappinfo(appInfo)
	client.httpClient = NewFastLivePushHttpClient(appInfo)
	return client
}

func (c *Client) BuildConnect() *Client {
	// 获取pushList
	logger.Infow("开始获取服务网关地址")
	pushList, err := c.getPushGateIpList()
	if err != nil {
		c.initialListener(err)
		return nil
	}
	c.pushGateIpList = pushList
	// 连接服务端
	logger.Infow("开始连接服务网关")
	err = c.connectServer()
	// 发送ConnAuth
	if err == nil {
		c.sendConnAuth()
		return c
	}
	return nil
}

// 发送通知
func (c *Client) SendPushNotification(notification PushNotification) {
	if c.isSendNotification {
		pushMessage := NewPushMessagePayloadFromPushNotification(notification, Push, &c.appInfo)
		if c.ch != nil || c.ch.IsActive() {
			c.ch.Write(pushMessage)
			c.sendListener(fmt.Sprintf("%s", pushMessage.messageID), nil)
		} else {
			c.sendListener(fmt.Sprintf("%s", pushMessage.messageID), errors.New("通道已经关闭"))
		}
	} else {
		c.sendListener(fmt.Sprintf("不能发送消息:"), errors.New("连接鉴权未完成"))
	}

}

// 发送voip通知{pushkit / callkit }
func (c *Client) SendVoipNotification(notification PushNotification) {
	if c.isSendNotification {
		pushMessage := NewPushMessagePayloadFromPushNotification(notification, VOIP, &c.appInfo)
		if c.ch != nil || c.ch.IsActive() {
			c.ch.Write(pushMessage)
			c.sendListener(fmt.Sprintf("%s", pushMessage.messageID), nil)
		} else {
			c.sendListener(fmt.Sprintf("%s", pushMessage.messageID), errors.New("通道已经关闭"))
		}
	} else {
		c.sendListener(fmt.Sprintf("不能发送消息:"), errors.New("连接鉴权未完成"))
	}

}

// 发送sms
func (c *Client) SendSMSMessage(notification PushNotification) {
	if c.isSendNotification {
		pushMessage := NewPushMessagePayloadFromPushNotification(notification, SMS, &c.appInfo)
		if c.ch != nil || c.ch.IsActive() {
			c.ch.Write(pushMessage)
			c.sendListener(fmt.Sprintf("%s", pushMessage.messageID), nil)
		} else {
			c.sendListener(fmt.Sprintf("%s", pushMessage.messageID), errors.New("通道已经关闭"))
		}
	} else {
		c.sendListener(fmt.Sprintf("不能发送消息:"), errors.New("连接鉴权未完成"))
	}

}
