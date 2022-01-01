package fastpushclient

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/transport/tcp"
	"github.com/lgphp/go-fastpushclient/logger"
	"time"
)

func (c *Client) pipelineInitializer() func(channel netty.Channel) {
	return func(ch netty.Channel) {
		ch.Pipeline().AddLast(newCodecHandler("codec-handler", c))
		ch.Pipeline().AddLast(newBizProcessorHandler("biz-handler", c))
		ch.Pipeline().AddLast(newEventHandler("evnet-handler"))
		ch.Pipeline().AddLast(newExceptionHandler("exception-handler"))
	}
}

// connect PushGate
func (c *Client) connectServer(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(c.ctx, time.Second*5)
	defer cancel()
	for _, addr := range c.pushGateIpList {
		pipeLine := c.pipelineInitializer()
		bootstrap := netty.NewBootstrap(netty.WithChannel(netty.NewBufferedChannel(128, 1024)),
			netty.WithClientInitializer(pipeLine), netty.WithTransport(tcp.New()))
		serverAddress := fmt.Sprintf("tcp://%s:%d", addr.IP, addr.Port)
		channel, err := bootstrap.Connect(serverAddress, "FastLivePushClient")
		if nil != err {
			logger.Warnw("connect failed", err, "connect info:", addr)
		} else {
			c.setChannel(channel)
			break
		}
	}
	for {
		select {
		case <-ctx.Done():
			c.initialListener(errors.New("connect to server timeout > 5 secs"))
			return errors.New("connect to server timeout > 5 secs")
		default:
			if c.ch == nil || !c.ch.IsActive() {
				c.initialListener(errors.New("can't connect to PushGate server , socket channel has nil or  inactive "))
				return errors.New("can't connect to PushGate server, socket channel has nil or  inactive")
			}
			return nil
		}
	}
}

func (c *Client) reConnectServer() {
	logger.Warnw("reconnecting to PushGate server", errors.New("reason: disconnected"))
	pushList, err := c.getPushGateIpList()
	if err != nil {
		logger.Warnw("can't re-connect to PushGate server", err)
	}
	c.pushGateIpList = pushList
	// 连接服务端
	err = c.connectServer(c.ctx)
	// 发送ConnAuth
	if err == nil {
		c.sendConnAuth()
	}
}

func (c *Client) promiseConnected() {
	if c.ch == nil || !c.ch.IsActive() {
		c.reConnectServer()
	}
}

//  发送连接认证
func (c *Client) sendConnAuth() {
	c.promiseConnected()
	payload := newConnAuthPayload(c.clientId, c.appInfo)
	c.ch.Write(payload)
}

// 处理连接认证回复
func (c *Client) handleConnAuthResp(payload connAuthRespPayload) {
	if payload.statusCode == HTTP_RESPONSE_CODE_OK {
		// 成功 ==> 发送心跳
		c.startHeartbeatTask()
		stime := payload.serverTime
		ctime := time.Now().UnixNano() / 1e6
		c.timeDiff = int64(stime) - ctime
		c.initialListener(nil)
		c.isSendNotification = true
	} else {
		// 鉴权不通过
		c.isSendNotification = false
		logger.Warnw("authentication of connection failed", errors.New(fmt.Sprintf("code:%v , message:%s", payload.statusCode, payload.message)))
		c.initialListener(errors.New(fmt.Sprintf("code:%v , message:%s", payload.statusCode, payload.message)))
	}
}

// 处理消息回执
func (c *Client) handleMessageACK(payload messageAckPayload) {
	c.messageStatusListener(payload.messageID, payload.userId, payload.appId, payload.statusCode, payload.statusMessage)
}

// 发送心跳
func (c *Client) startHeartbeatTask() {
	go func() {
		for {
			payload := newHeartBeatPayload()
			if c.ch != nil && c.ch.IsActive() && c.isSendNotification {
				c.ch.Write(payload)
				time.Sleep(time.Second * 15)
			}
			break
		}
	}()

}
