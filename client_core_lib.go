package fastpushclient

import (
	"fmt"
	"github.com/lgphp/go-fastpushclient/logger"
	"github.com/pkg/errors"
	"time"
)

func NewFastLivePushClient(appInfo AppInfo) *Client {
	logger.Infow("Initial PushClient")
	client := buildClient()
	client.setAppinfo(appInfo)
	client.httpClient = NewFastLivePushHttpClient(appInfo)
	return client
}

func (c *Client) BuildConnect() (*Client, error) {

	if c.initialListener == nil {
		return nil, errors.New("InitialListener must not be nil")
	}
	if c.sendListener == nil {
		return nil, errors.New("SendListener must not be nil")
	}
	if c.messageStatusListener == nil {
		return nil, errors.New("MessageStatusListener must not be nil")
	}
	// 获取pushList
	logger.Infow("Will get PushGate server address")
	pushList, err := c.getPushGateIpList()
	if err != nil {
		c.initialListener(err)
		return nil, err
	}
	c.pushGateIpList = pushList
	// 连接服务端
	logger.Infow("Connect PushGate server")
	err = c.connectServer(c.ctx)
	// 发送ConnAuth
	if err == nil {
		c.sendConnAuth()
		return c, nil
	}
	return nil, err

}

func (c *Client) sendTask() {
	for msg := range c.sendQueue {
		pushmsg := msg
		c.workerpool.Submit(func() {
			if c.ch != nil || c.ch.IsActive() {
				c.ch.Write(pushmsg)
				c.sendListener(fmt.Sprintf("%s", pushmsg.messageID), nil)
			} else {
				c.sendListener(fmt.Sprintf("%s", pushmsg.messageID), errors.New("Connection has been closed"))
			}
		})
		// speed control
		c.promiseConnected()
		time.Sleep(time.Millisecond * time.Duration(c.sendSpeed))
	}
}

// 发送push通知
func (c *Client) SendPushNotification(pushNotification PushNotification) {

	if c.isSendNotification {
		pushMessage := NewPushMessagePayloadFromPushNotification(pushNotification, Push, &c.appInfo)
		c.sendQueue <- pushMessage
	} else {
		c.sendListener(fmt.Sprintf("Didn't send push message:"), errors.New("Authentication of connection not finished"))
	}

}

// 发送voip通知{pushkit / callkit }
func (c *Client) SendVoipNotification(pushNotification PushNotification) {

	if c.isSendNotification {
		pushMessage := NewPushMessagePayloadFromPushNotification(pushNotification, VOIP, &c.appInfo)
		c.sendQueue <- pushMessage
		c.sendListener(fmt.Sprintf("Didn't send push message:"), errors.New("Authentication of connection not finished"))
	}

}

type SmsMessage = PushNotification

// 发送sms
func (c *Client) SendSMSMessage(smsMessage SmsMessage) {

	if c.isSendNotification {
		pushMessage := NewPushMessagePayloadFromPushNotification(smsMessage, SMS, &c.appInfo)
		c.sendQueue <- pushMessage
	} else {
		c.sendListener(fmt.Sprintf("Didn't send push message:"), errors.New("Authentication of connection not finished"))
	}

}
