package fastpushclient

import (
	"context"
	"encoding/base64"
	"github.com/gammazero/workerpool"
	"github.com/go-netty/go-netty"
	"github.com/rogpeppe/fastuuid"
	"go.uber.org/atomic"
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
	// is auth of socket connection  passed ?
	isSendNotification bool
	// time difference between client and server
	timeDiff             int64
	httpClient           HTTPClient
	ctx                  context.Context
	sendQueue            chan PushMessagePayload
	workerpool           *workerpool.WorkerPool
	sendSpeed            uint16
	bootstrap            netty.Bootstrap
	isRetryConnecting    *atomic.Bool
	retryCnt             *atomic.Int32
	max_Send_Buffer_Size uint32
}

func buildClient() *Client {
	client := &Client{
		wg:                   &sync.WaitGroup{},
		clientId:             fastuuid.MustNewGenerator().Hex128(),
		pushGateIpList:       make([]pushGateAddress, 0),
		ctx:                  context.Background(),
		workerpool:           workerpool.New(10),
		sendSpeed:            uint16(30),
		isRetryConnecting:    atomic.NewBool(false),
		retryCnt:             atomic.NewInt32(0),
		max_Send_Buffer_Size: 10000,
	}
	client.sendQueue = make(chan PushMessagePayload, client.max_Send_Buffer_Size)
	return client
}

func (c *Client) SetSendBuffSize(size uint32) *Client {
	c.sendQueue = make(chan PushMessagePayload, size)
	return c
}

func (c *Client) setAppinfo(appInfo AppInfo) {
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
