package fastpushclient

import (
	"github.com/go-netty/go-netty"
	"github.com/lgphp/go-fastpushclient/logger"
	"github.com/pkg/errors"
)

type BizProcessorHandler struct {
	netty.ChannelInboundHandler
	name   string
	client *Client
}

func newBizChannelHandler(name string, c *Client) BizProcessorHandler {
	return BizProcessorHandler{
		name:   name,
		client: c,
	}
}

func (h BizProcessorHandler) HandleActive(ctx netty.ActiveContext) {
	logger.Infow("连接成功,开始准备通信", "远端服务器地址", ctx.Channel().RemoteAddr())
}

func (h BizProcessorHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {
	switch message.(type) {
	case connAuthRespPayload:
		if payload, ok := message.(connAuthRespPayload); ok {
			h.client.dealConnAuthResp(payload)
		}
		break
	default:
		logger.Warnw("业务处理器", errors.New("无需处理的Payload"), "payload", message)
		break
	}
}

func (h BizProcessorHandler) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	logger.Warnw("与服务器断开连接", ex)
	// 重新连接
	h.client.isSendNotification = false
	h.client.reConnectServer()
}
