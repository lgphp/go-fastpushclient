package fastpushclient

import (
	"github.com/go-netty/go-netty/codec"
	"github.com/go-netty/go-netty/utils"
	"github.com/pkg/errors"
	"push-sdk-go/bytebuf"
	"push-sdk-go/logger"
)

type CodecHandler struct {
	name           string
	c              *Client
	maxFrameLength int
	buffer         []byte
}

func newCodecHandler(name string, client *Client) codec.Codec {
	return &CodecHandler{
		name:           name,
		c:              client,
		maxFrameLength: 1024,
		buffer:         make([]byte, 1024),
	}
}

func (h CodecHandler) CodecName() string {
	return h.name
}

// 解码
func (h CodecHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {
	reader := utils.MustToReader(message)
	n, _ := reader.Read(h.buffer)
	buf, _ := bytebuf.NewByteBuf(h.buffer[:n])
	if len(buf.AvailableBytes()) < MIN_TCP_PACkET_LENGTH {
		return
	}
	// 读取包体长度
	pktLen, _ := buf.ReadUInt32BE()
	if pktLen > MAX_TCP_PACkET_LENGTH-4 {
		return
	}
	// 读取版本号
	_, _ = buf.ReadByte()
	// 读取payloadCode
	payloadCode, _ := buf.ReadUInt16BE()
	switch payloadCode {
	case ConnAuthRespCode:
		carp := newConnAuthRespPayload()
		carp.Unpack(buf, h.c)
		ctx.HandleRead(carp)
		break
	default:
		logger.Warnw("解码器", errors.New("不需要解码的PayloadCode"),
			"payloadCode", payloadCode)
		break
	}
	buf.Release()
}

// 编码
func (h CodecHandler) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	buf, _ := bytebuf.NewByteBuf()
	switch message.(type) {
	case heartBeatPayload:
		beatPayload := message.(heartBeatPayload)
		beatPayload.Pack(buf, h.c)
		ctx.Write(buf.AvailableBytes())
		break
	case connAuthPayload:
		authPayload := message.(connAuthPayload)
		authPayload.Pack(buf, h.c)
		ctx.Write(buf.AvailableBytes())
		break
	case PushMessagePayload:
		pushPayload := message.(PushMessagePayload)
		pushPayload.Pack(buf, h.c)
		ctx.Write(buf.AvailableBytes())
		break
	case TokenUploadPayload:
		tokenUploadPayload := message.(TokenUploadPayload)
		tokenUploadPayload.Pack(buf, h.c)
		ctx.Write(buf.AvailableBytes())
		break
	default:
		logger.Warnw("编码器", errors.New("不能识别的Payload,不能编码"),
			"payload", message)
		break
	}
	buf.Release()
}
