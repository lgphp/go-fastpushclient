package main

import (
	"fmt"
	pushSDK "github.com/lgphp/go-fastpushclient"
	"github.com/lgphp/go-fastpushclient/logger"
	"sync"
	"time"
)

func init() {
	logger.InitDevelopment("DEBUG")
}

var (
	ch = make(chan bool)
	c  = 0
	// 测试环境的相关配置
	TEST_ENV_MERCHANT_ID = "a127297f117c4a3fb095a15443bc96fc"
	TEST_ENV_APP_ID      = "b4722bb12f30485582fb3e3a5c6157c6"
	TEST_ENV_APP_KEY     = "NUhONBRTxPxFtFkH78P9AJ2EDUJ1EeoaFzGVoJUz5BcYFtqiag0baRw61y1ycoZaYpkxp9BC08K2F8h2II4tyQ=="
	TEST_ENV_USER_ID     = "8613810654610"

	//本地环境相关配置
	LOCAL_ENV_MERCHANT_ID = "3527348b4ff04e988f3fadd7f1e4f155"
	LOCAL_ENV_APP_ID      = "4fb367eabe8c45f2a1b6714c6e40fd19"
	LOCAL_ENV_APP_KEY     = "fCCHFkVtMk7sf5XQmfyTpPvuVH0PKmUd559HUtTlDxCBg5y4P3SAcxnAgxCG/AuRO0y//ZbgwRQg1wCJGOWw/w=="
	LOCAL_ENV_USER_ID     = "97158000000"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	info := pushSDK.NewAppInfo(TEST_ENV_MERCHANT_ID,
		TEST_ENV_APP_ID,
		TEST_ENV_APP_KEY)
	client, _ := pushSDK.NewFastLivePushClient(info).AddInitializedListener(initialSDKCallback).AddSendListener(sendCallBack).AddNotificationStatusListener(notificationCallBack).BuildConnect()
	<-ch
	sendNotification(client)
	wg.Wait()
}

func sendNotification(client *pushSDK.Client) {
	// 发送100条消息
	for i := 1; i <= 20; i++ {
		body, _ := pushSDK.NewMessageBody(fmt.Sprintf("%s+:%v", "标题", i), "消息体", nil)
		notification := pushSDK.NewPushNotification(TEST_ENV_USER_ID, pushSDK.LOW, body)
		client.SendPushNotification(notification)
		time.Sleep(time.Microsecond * 1)
	}

}

// SDK 初始化回调
func initialSDKCallback(err error) {
	if nil != err {
		logger.Warnw("无法初始化SDK ", err)
	} else {
		logger.Infow("SDK初始化成功")
		ch <- true
	}
}

func notificationCallBack(messageId, toUserId, appId string, statusCode uint32, statusText string) {
	logger.Infow("投递结果", "messageId", messageId, "toUserId", toUserId, "statusText", statusText)
	if statusCode == 1 {
		c = c + 1
		println("第", c, "条回执")
	}
}

// 发送消息回调
func sendCallBack(messageId string, err error) {
	//if err != nil {
	//	logger.Warnw("发送失败:", err, "messageId", messageId)
	//} else {
	//	logger.Infow("发送成功", "messageId", messageId)
	//}
}
