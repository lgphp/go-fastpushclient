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

func newBizProcessorHandler(name string, c *Client) *BizProcessorHandler {
	return &BizProcessorHandler{
		name:   name,
		client: c,
	}
}

func (h *BizProcessorHandler) HandleActive(ctx netty.ActiveContext) {
	logger.Infow("Connected", "remoteAddr", ctx.Channel().RemoteAddr())
}

func (h *BizProcessorHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {
	switch message.(type) {
	case connAuthRespPayload:
		if payload, ok := message.(connAuthRespPayload); ok {
			h.client.handleConnAuthResp(payload)
		}
		break
	case messageAckPayload:
		payload := message.(messageAckPayload)
		h.client.handleMessageACK(payload)
		break
	default:
		logger.Warnw("bussiness handler:", errors.New("Unknow Payload"), "payload", message)
		break
	}
}

func (h *BizProcessorHandler) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	logger.Warnw("disconnect to remote server", ex)
	// 重新连接
	h.client.isSendNotification = false
	h.client.reConnectServer()
}
