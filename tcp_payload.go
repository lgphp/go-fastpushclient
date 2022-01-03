package fastpushclient

import (
	"encoding/json"
	"github.com/lgphp/go-bytebuf"
	"github.com/lgphp/go-fastpushclient/logger"
	"github.com/lgphp/go-fastpushclient/utils"
	"github.com/rogpeppe/fastuuid"
	"github.com/wumansgy/goEncrypt"
)

type Payloadble interface {
	Pack(buf *bytebuf.ByteBuf, client *Client)
	Unpack(buf *bytebuf.ByteBuf, client *Client)
}

const (
	MAX_TCP_PACkET_LENGTH = 65535
	MIN_TCP_PACkET_LENGTH = 8
	NEED_COMPRESS_SIZE    = 1024
)

type PayloadCode = uint16
type SPChannel = byte
type NotificationClassify = byte
type MessagePrior = byte

var (
	HeartBeatCode      PayloadCode = 10000
	ConnAuthCode       PayloadCode = 20000
	ConnAuthRespCode   PayloadCode = 20001
	PushMessageCode    PayloadCode = 30000
	PushMessageACKCode PayloadCode = 30001
	TokenUploadCode    PayloadCode = 40000

	// SpChannel enum
	APNS          SPChannel = 10
	FCM           SPChannel = 11
	HCM           SPChannel = 12
	APPLE_PUSHKIT SPChannel = 60
	HW_PUSHKIT    SPChannel = 61

	// Notification classifier enum
	Push         NotificationClassify = 1
	SMS          NotificationClassify = 2
	Email        NotificationClassify = 3
	InnerMessage NotificationClassify = 4
	VOIP         NotificationClassify = 5

	// Message priority enum
	HIGH   MessagePrior = 0
	MIDDLE MessagePrior = 1
	LOW    MessagePrior = 2
)

// tcp上报Token实体
type TokenUploadPayload struct {
	payloadCode uint16
	appId       string
	userId      string
	spChannel   SPChannel
	pushToken   string
	classifier  NotificationClassify
}

// 暂时不开放
func newPushTokenInfo(userId, pushToken string, spChannel SPChannel, classifier NotificationClassify) TokenUploadPayload {
	return TokenUploadPayload{
		payloadCode: TokenUploadCode,
		userId:      userId,
		spChannel:   spChannel,
		pushToken:   pushToken,
		classifier:  classifier,
	}
}

func (c *TokenUploadPayload) Pack(buf *bytebuf.ByteBuf, clinet *Client) {
	c.appId = clinet.appInfo.appID
	// 写包长度占位 , frame.LengthFieldCodec 编码器自动添加，所以不要写
	_ = buf.WriteUInt32BE(0)
	// 写版本号
	_ = buf.WriteByte(1)
	// 写类型码
	_ = buf.WriteUInt16BE(c.payloadCode)

	//写appID
	appIdBytes := []byte(c.appId)
	appIdLen := byte(len(appIdBytes))
	_ = buf.WriteByte(appIdLen)
	_ = buf.WriteBytes(appIdBytes)

	// 写UserID
	userIdBytes := []byte(c.userId)
	userIdLen := byte(len(userIdBytes))
	_ = buf.WriteByte(userIdLen)
	_ = buf.WriteBytes(userIdBytes)

	_ = buf.WriteByte(c.spChannel)

	// 写token
	pushTokenBytes := []byte(c.pushToken)
	pushTokenLen := uint32(len(pushTokenBytes))
	_ = buf.WriteUInt32BE(pushTokenLen)
	_ = buf.WriteBytes(pushTokenBytes)

	// 写classfier
	_ = buf.WriteByte(c.classifier)

	pktLen := buf.WriterIndex() - 4
	_ = buf.PutUInt32BE(0, uint32(pktLen))
}

func (c *TokenUploadPayload) Unpack(buf *bytebuf.ByteBuf, _ *Client) {
	println("无需实现")
}

type heartBeatPayload struct {
	payloadCode uint16
	zero        byte
}

func newHeartBeatPayload() heartBeatPayload {
	return heartBeatPayload{
		payloadCode: HeartBeatCode,
		zero:        0,
	}
}

// 心跳
func (c *heartBeatPayload) Pack(buf *bytebuf.ByteBuf, _ *Client) {
	// 写包长度占位 , frame.LengthFieldCodec 编码器自动添加，所以不要写
	_ = buf.WriteUInt32BE(0)
	// 写版本号
	_ = buf.WriteByte(1)
	// 写类型码
	_ = buf.WriteUInt16BE(c.payloadCode)
	_ = buf.WriteByte(c.zero)
	pktLen := buf.WriterIndex() - 4
	_ = buf.PutUInt32BE(0, uint32(pktLen))
}

// 无需解码
func (c *heartBeatPayload) Unpack(buf *bytebuf.ByteBuf, _ *Client) {
	logger.Debugw("no need tp implement")
}

// 连接鉴权包
type connAuthPayload struct {
	payloadCode uint16
	clientId    string
	merchantID  string
	appID       string
	authKey     []byte
}

func (c *connAuthPayload) Pack(buf *bytebuf.ByteBuf, _ *Client) {
	// 写包长度占位 , frame.LengthFieldCodec 编码器自动添加，所以不要写
	_ = buf.WriteUInt32BE(0)
	// 写版本号
	_ = buf.WriteByte(1)
	// 写类型码
	_ = buf.WriteUInt16BE(c.payloadCode)

	// 写ClientID
	_ = buf.WriteStringWithByteLength(c.clientId)

	// 写merchantID
	_ = buf.WriteStringWithByteLength(c.merchantID)

	//写appID
	_ = buf.WriteStringWithByteLength(c.appID)

	//写AuthKey
	_ = buf.WriteBytes(c.authKey)

	pktLen := buf.WriterIndex() - 4
	_ = buf.PutUInt32BE(0, uint32(pktLen))

}

func (c *connAuthPayload) Unpack(buf *bytebuf.ByteBuf, _ *Client) {
	logger.Debugw("no need tp implement")
}

func newConnAuthPayload(clientId string, info AppInfo) connAuthPayload {
	authKey := utils.GetAuthKey(info.GetAppKey())
	return connAuthPayload{
		payloadCode: ConnAuthCode,
		clientId:    clientId,
		merchantID:  info.merchantID,
		appID:       info.appID,
		authKey:     authKey,
	}
}

// 连接认证回应包
type connAuthRespPayload struct {
	payloadCode uint16
	statusCode  uint32
	message     string
	serverTime  uint64
}

func newConnAuthRespPayload() connAuthRespPayload {
	return connAuthRespPayload{
		payloadCode: ConnAuthRespCode,
	}
}

func (c *connAuthRespPayload) Pack(buf *bytebuf.ByteBuf, _ *Client) {
	logger.Debugw("no need tp implement")
}

func (c *connAuthRespPayload) Unpack(buf *bytebuf.ByteBuf, _ *Client) {
	c.statusCode, _ = buf.ReadUInt32BE()
	c.message = buf.ReadStringWithU32Length()
	stime, _ := buf.ReadUInt64BE()
	c.serverTime = stime
}

// PUSH消息包
type PushMessagePayload struct {
	payloadCode uint16
	messageID   string
	classifier  NotificationClassify
	merchantID  string
	appID       string
	priority    MessagePrior
	toUid       string
	messageBody []byte
}

func (p *PushMessagePayload) Pack(buf *bytebuf.ByteBuf, client *Client) {
	// 写包长度占位
	_ = buf.WriteUInt32BE(0)
	// 写版本号
	_ = buf.WriteByte(1)
	// 写类型码
	_ = buf.WriteUInt16BE(p.payloadCode)
	// 写messageID
	_ = buf.WriteStringWithByteLength(p.messageID)

	// 写通知分类
	_ = buf.WriteByte(p.classifier)
	// 写merchantID
	_ = buf.WriteStringWithByteLength(p.merchantID)
	//写appID
	_ = buf.WriteStringWithByteLength(p.appID)

	var encFlag byte = 2
	// 加密消息,获取消息加密key CBC 模式
	encMessageKey := utils.GetMsgEncKey(client.appInfo.GetAppKey())
	encAesIV := utils.GetMsgEncAesIV(client.appInfo.GetAppKey())
	messageBody, _ := goEncrypt.AesCbcEncrypt(p.messageBody, encMessageKey, encAesIV)
	if len(messageBody) > NEED_COMPRESS_SIZE {
		// 压缩
		messageBody = utils.Gzip(messageBody)
		encFlag = 3
	}
	// 写加密标志
	_ = buf.WriteByte(encFlag)
	// 写优先级
	_ = buf.WriteByte(p.priority)
	// 写toUserID
	_ = buf.WriteStringWithByteLength(p.toUid)

	// 写MessageBody
	_ = buf.WriteBytes(messageBody)
	pktLen := buf.WriterIndex() - 4
	_ = buf.PutUInt32BE(0, uint32(pktLen))

}

func (p *PushMessagePayload) Unpack(buf *bytebuf.ByteBuf, _ *Client) {
	logger.Debugw("no need tp implement")
}

// 创建一个新push通知
func NewPushMessagePayloadFromPushNotification(n PushNotification, classifier NotificationClassify, app *AppInfo) PushMessagePayload {
	messageId := fastuuid.MustNewGenerator().Hex128()
	messageBody, _ := json.Marshal(n.messageBody)
	return PushMessagePayload{
		payloadCode: PushMessageCode,
		messageID:   messageId,
		classifier:  classifier,
		merchantID:  app.GetMerchantID(),
		appID:       app.GetAppID(),
		priority:    n.priority,
		toUid:       n.toUid,
		messageBody: messageBody,
	}
}

// 消息回执包
type messageAckPayload struct {
	payloadCode   uint16
	messageID     string
	appId         string
	userId        string
	statusCode    uint32
	statusMessage string
}

func newAckMessageAckPayload() messageAckPayload {
	return messageAckPayload{
		payloadCode: PushMessageACKCode,
	}
}
func (c *messageAckPayload) Pack(buf *bytebuf.ByteBuf, _ *Client) {
	logger.Debugw("no need tp implement")
}

func (c *messageAckPayload) Unpack(buf *bytebuf.ByteBuf, _ *Client) {

	c.messageID = buf.ReadStringWithByteLength()
	c.appId = buf.ReadStringWithByteLength()
	c.userId = buf.ReadStringWithByteLength()
	code, _ := buf.ReadUInt32BE()
	c.statusCode = code

	c.statusMessage = buf.ReadStringWithU32Length()

}
