package fastpushclient

import (
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec"
	"github.com/go-netty/go-netty/utils"
	"github.com/lgphp/go-bytebuf"
	"github.com/lgphp/go-fastpushclient/logger"
	"github.com/pkg/errors"
	"io"
)

type CodecHandler struct {
	name           string
	c              *Client
	maxFrameLength int
	allbuf         []byte
}

func newCodecHandler(name string, maxFrameLength int, client *Client) codec.Codec {
	return &CodecHandler{
		name:           name,
		c:              client,
		maxFrameLength: maxFrameLength,
		allbuf:         make([]byte, 0),
	}
}

func (h *CodecHandler) CodecName() string {
	return h.name
}

// 解码
func (h *CodecHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {
	buffer := make([]byte, h.maxFrameLength)
	reader := utils.MustToReader(message)
	n, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		return
	}

	if n < MIN_TCP_PACkET_LENGTH {
		return
	}
	// 读取包体长度
	if n > MAX_TCP_PACkET_LENGTH-4 {
		return
	}
	//  handle  half packet and stick packet
	//  very important for socket communication
	//

	if len(h.allbuf) != 0 {
		h.allbuf = append(h.allbuf, buffer[:n]...)
	} else {
		h.allbuf = buffer[:n]
	}
	buf, _ := bytebuf.NewByteBuf(h.allbuf[:])
	//println("读到：n" , n , "buffer" , len(buffer) , "h.allbuf" , len(h.allbuf))
	defer func() {
		buf.Release()
	}()
	for {
		// if readablebytes  < length field length
		if buf.ReadableBytes() < 4 {
			return
		}
		pktLen, _ := buf.ReadUInt32BE()
		// if readablebytes <  packet Length
		if buf.ReadableBytes() < int(pktLen) {
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
		case PushMessageACKCode:
			amap := newAckMessageAckPayload()
			amap.Unpack(buf, h.c)
			ctx.HandleRead(amap)
			break
		default:
			logger.Warnw("decoder", errors.New("unknow PayloadCode"),
				"payloadCode", payloadCode)
			buf.SkipBytes(int(pktLen - 3))
			break
		}
		// just read complete
		if buf.ReaderIndex() == buf.WriterIndex() {
			// make allbuf = []
			h.allbuf = make([]byte, 0)
			return
		}
	}

}

// 编码
func (h *CodecHandler) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	buf, _ := bytebuf.NewByteBuf()
	defer func() {
		buf.Release()
	}()
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
		logger.Warnw("encoder", errors.New("unknow Payload"),
			"payload", message)
		break
	}
}
