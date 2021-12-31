package fastpushclient

import (
	"encoding/base64"
	"github.com/go-netty/go-netty"
	"github.com/rogpeppe/fastuuid"
	"sync"
)

type pushGateAddress struct {
	IP   string
	Port int
}

type Client struct {
	wg                    *sync.WaitGroup
	clientId              string
	appInfo               AppInfo
	ch                    netty.Channel
	pushGateIpList        []pushGateAddress
	initialListener       InitializedListener
	sendListener          ClientSendListener
	messageStatusListener NotificationStatusListener
	// 是否能发送消息
	isSendNotification bool
	// 与服务器时间的差值
	timeDiff   int64
	httpClient HTTPClient
}

func buildClient() *Client {
	return &Client{
		wg:             &sync.WaitGroup{},
		clientId:       fastuuid.MustNewGenerator().Hex128(),
		pushGateIpList: make([]pushGateAddress, 0),
	}
}

func (c *Client) setappinfo(appInfo AppInfo) {

	c.appInfo = appInfo
}

func (c *Client) AddSendListener(l ClientSendListener) *Client {
	c.sendListener = l
	return c
}

func (c *Client) AddInitializedListener(l InitializedListener) *Client {
	c.initialListener = l
	return c
}

func (c *Client) AddNotificationStatusListener(l NotificationStatusListener) *Client {
	c.messageStatusListener = l
	return c
}

func (c *Client) setPushGateIpList(ipList []pushGateAddress) {
	c.pushGateIpList = ipList
}

func (c *Client) setChannel(ch netty.Channel) {
	c.ch = ch
}

type AppInfo struct {
	merchantID string
	appID      string
	appKey     []byte
}

func NewAppInfo(merchantID, appID, appKey string) AppInfo {
	bytes, _ := base64.StdEncoding.DecodeString(appKey)
	return AppInfo{
		merchantID: merchantID,
		appID:      appID,
		appKey:     bytes,
	}
}

func (a *AppInfo) GetMerchantID() string {
	return a.merchantID
}
func (a *AppInfo) GetAppID() string {
	return a.appID
}
func (a *AppInfo) GetAppKey() []byte {
	return a.appKey
}
