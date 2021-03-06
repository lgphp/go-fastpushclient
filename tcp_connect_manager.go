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
		ch.Pipeline().AddLast(newCodecHandler("codec-handler", 65535, c))
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
		bootstrap := netty.NewBootstrap(netty.WithChannel(netty.NewBufferedChannel(1, 4096)),
			netty.WithClientInitializer(pipeLine), netty.WithTransport(tcp.New()))
		serverAddress := fmt.Sprintf("tcp://%s:%d", addr.IP, addr.Port)
		channel, err := bootstrap.Connect(serverAddress, "FastLivePushClient")
		if nil != err {
			logger.Warnw("connect failed", err, "connect info:", addr)
		} else {
			c.setChannel(channel)
			c.bootstrap = bootstrap
			break
		}
	}
	for {
		select {
		case <-ctx.Done():
			c.initialListener(errors.New("connect to server timeout > 5 secs"))
			c.reConnectServer()
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
	if c.isRetryConnecting.Load() {
		return
	}
	c.isRetryConnecting.Store(true)
	go func() {
		for {
			c.retryCnt.Add(1)
			logger.Warnw("reconnecting to PushGate server", errors.New("reason: disconnected"), "times", c.retryCnt)
			// 重试连接
			time.Sleep(time.Second * time.Duration((c.retryCnt.Load())*10))

			if c.retryCnt.Load() > 61 {
				c.initialListener(errors.New("out of re-connect times, client will Shutdown!!!!!"))
				c.bootstrap.Shutdown()
				return
			}
			pushList, err := c.getPushGateIpList()
			if err != nil {
				logger.Warnw("can't re-connect to PushGate server", err)
				continue
			}
			c.pushGateIpList = pushList
			// 连接服务端
			err = c.connectServer(c.ctx)
			// 发送ConnAuth
			if err == nil {
				c.sendConnAuth()
				c.isRetryConnecting.Store(false)
				c.retryCnt.Store(0)
				return
			}
		}

	}()

}

func (c *Client) promiseConnected() {
	if c.ch == nil || !c.ch.IsActive() {
		c.reConnectServer()
	}
}

//  发送连接认证
func (c *Client) sendConnAuth() {
	payload := newConnAuthPayload(c.clientId, c.appInfo)
	c.ch.Write(payload)
}

// 处理连接认证回复
func (c *Client) handleConnAuthResp(payload connAuthRespPayload) {
	if payload.statusCode == HTTP_RESPONSE_CODE_OK {
		stime := payload.serverTime
		ctime := time.Now().UnixNano() / 1e6
		c.timeDiff = int64(stime) - ctime
		c.isSendNotification = true
		//设置发送速度
		c.sendSpeed = payload.speedLimit
		logger.Infow(fmt.Sprintf("Connect Authentication Success , Send Speed Limited : %d /sec", 1000/c.sendSpeed))
		// 发送心跳
		go c.startHeartbeatTask()
		// 启动发送任务
		go c.sendTask()
		// 设置成功回调
		c.initialListener(nil)
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
	for {
		payload := newHeartBeatPayload()
		if c.ch != nil && c.ch.IsActive() && c.isSendNotification {
			c.ch.Write(payload)
			time.Sleep(time.Second * 15)
		} else {
			c.reConnectServer()
			return
		}

	}
}
