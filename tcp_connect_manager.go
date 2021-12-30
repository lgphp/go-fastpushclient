package fastpushclient

import (
	"errors"
	"fmt"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/transport/tcp"
	"github.com/lgphp/go-fastpushclient/logger"
	"time"
)

func (c *Client) pipeLineInitializer() func(channel netty.Channel) {
	return func(ch netty.Channel) {
		// 基于长度字段解码器
		//ch.Pipeline().AddLast(frame.LengthFieldCodec(binary.BigEndian, 65535, 0, 4, 0, 0))
		ch.Pipeline().AddLast(newCodecHandler("编解码器", c))
		ch.Pipeline().AddLast(newBizChannelHandler("业务处理器", c))
		ch.Pipeline().AddLast(newEventHandler("事件处理器"))
		ch.Pipeline().AddLast(newExceptionHandler("异常处理器"))
	}
}

// 连接pushGate
func (c *Client) connectServer() error {
	//addr := pushGateAddress{
	//	IP:   "10.110.240.49",
	//	Port: 4442,
	//}
	//c.pushGateIpList = append(c.pushGateIpList, addr)
	for _, addr := range c.pushGateIpList {
		pipeLine := c.pipeLineInitializer()
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
	if c.ch == nil || !c.ch.IsActive() {
		c.initialListener(errors.New("无法连接远程服务器"))
		return errors.New("无法连接远程服务器")
	}
	return nil
}

func (c *Client) reConnectServer() {
	logger.Warnw("重新连接服务器", errors.New("服务器断开连接"))
	pushList, err := c.getPushGateIpList()
	if err != nil {
		logger.Warnw("无法重新连接服务器", err)
	}
	c.pushGateIpList = pushList
	// 连接服务端
	err = c.connectServer()
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
func (c *Client) dealConnAuthResp(payload connAuthRespPayload) {
	if payload.StatusCode == Response_Success {
		// 成功 ==> 发送心跳
		c.startHeartbeatTask()
		stime := payload.ServerTime
		ctime := time.Now().UnixNano() / 1e6
		c.timeDiff = int64(stime) - ctime
		c.initialListener(nil)
		c.isSendNotification = true
	} else {
		// 鉴权不通过
		c.isSendNotification = false
		logger.Warnw("长连接鉴权不通过", errors.New(fmt.Sprintf("code:%v , message:%s", payload.StatusCode, payload.Message)))
		c.initialListener(errors.New(fmt.Sprintf("code:%v , message:%s", payload.StatusCode, payload.Message)))
	}
}

// 发送心跳
func (c *Client) startHeartbeatTask() {
	go func() {
		for {
			payload := newHeartBeatPayload()
			if c.ch == nil || !c.ch.IsActive() {
				return
			}
			c.ch.Write(payload)
			time.Sleep(time.Second * 15)
		}
	}()

}
