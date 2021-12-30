package fastpushclient

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/rogpeppe/fastuuid"
	"github.com/wumansgy/goEncrypt"
	"push-sdk-go/bytebuf"
	"push-sdk-go/client/utils"
	"push-sdk-go/logger"
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
	HeartBeatCode    PayloadCode = 10000
	ConnAuthCode     PayloadCode = 20000
	ConnAuthRespCode PayloadCode = 20001
	PushMessageCode  PayloadCode = 30000
	TokenUploadCode  PayloadCode = 40000
	APNS             SPChannel   = 10
	FCM              SPChannel   = 11
	HCM              SPChannel   = 12

	Push         NotificationClassify = 1
	SMS          NotificationClassify = 2
	EMAIL        NotificationClassify = 3
	InnerMessage NotificationClassify = 4
	VOIP         NotificationClassify = 5

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
	// 写包长度占位
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
	PayloadCode uint16
	Zero        byte
}

func newHeartBeatPayload() heartBeatPayload {
	return heartBeatPayload{
		PayloadCode: HeartBeatCode,
		Zero:        0,
	}
}

// 心跳
func (c *heartBeatPayload) Pack(buf *bytebuf.ByteBuf, _ *Client) {
	// 写包长度占位
	_ = buf.WriteUInt32BE(0)
	// 写版本号
	_ = buf.WriteByte(1)
	// 写类型码
	_ = buf.WriteUInt16BE(c.PayloadCode)
	_ = buf.WriteByte(c.Zero)
	pktLen := buf.WriterIndex() - 4
	_ = buf.PutUInt32BE(0, uint32(pktLen))
}

// 无需解码
func (c *heartBeatPayload) Unpack(buf *bytebuf.ByteBuf, _ *Client) {
	println("无需实现")
}

// 连接鉴权包
type connAuthPayload struct {
	PayloadCode uint16
	ClientId    string
	MerchantID  string
	AppID       string
	AuthKey     []byte
}

func (c *connAuthPayload) Pack(buf *bytebuf.ByteBuf, _ *Client) {
	// 写包长度占位
	_ = buf.WriteUInt32BE(0)
	// 写版本号
	_ = buf.WriteByte(1)
	// 写类型码
	_ = buf.WriteUInt16BE(c.PayloadCode)

	// 写ClientID
	clientIdBytes := []byte(c.ClientId)
	clientIdLen := byte(len(clientIdBytes))
	_ = buf.WriteByte(clientIdLen)
	_ = buf.WriteBytes(clientIdBytes)

	// 写merchantID
	merchantIdBytes := []byte(c.MerchantID)
	merchantIdLen := byte(len(merchantIdBytes))
	_ = buf.WriteByte(merchantIdLen)
	_ = buf.WriteBytes(merchantIdBytes)

	//写appID
	appIdBytes := []byte(c.AppID)
	appIdLen := byte(len(appIdBytes))
	_ = buf.WriteByte(appIdLen)
	_ = buf.WriteBytes(appIdBytes)

	//写AuthKey
	_ = buf.WriteBytes(c.AuthKey)
	pktLen := buf.WriterIndex() - 4
	_ = buf.PutUInt32BE(0, uint32(pktLen))

}

func (c *connAuthPayload) Unpack(buf *bytebuf.ByteBuf, _ *Client) {
	println("无需实现")
}

func newConnAuthPayload(clientId string, info AppInfo) connAuthPayload {
	authKey := utils.GetAuthKey(info.GetAppKey())
	return connAuthPayload{
		PayloadCode: ConnAuthCode,
		ClientId:    clientId,
		MerchantID:  info.merchantID,
		AppID:       info.appID,
		AuthKey:     authKey,
	}
}

// 鉴权回复包
type connAuthRespPayload struct {
	PayloadCode uint16
	StatusCode  uint32
	Message     string
	ServerTime  uint64
}

func newConnAuthRespPayload() connAuthRespPayload {
	return connAuthRespPayload{
		PayloadCode: ConnAuthRespCode,
	}
}

func (c *connAuthRespPayload) Pack(buf *bytebuf.ByteBuf, _ *Client) {
	println("无需实现")
}

func (c *connAuthRespPayload) Unpack(buf *bytebuf.ByteBuf, _ *Client) {
	// 读取相关
	code, _ := buf.ReadUInt32BE()
	c.StatusCode = code
	messageLen, _ := buf.ReadUInt32BE()
	msg := make([]byte, messageLen)
	_, _ = buf.ReadBytes(msg)
	c.Message = string(msg)
	stime, _ := buf.ReadUInt64BE()
	c.ServerTime = stime
}

type PushMessagePayload struct {
	PayloadCode uint16
	MessageID   string
	Classifier  NotificationClassify
	MerchantID  string
	AppID       string
	Priority    MessagePrior
	ToUid       string
	MessageBody []byte
}

func (p *PushMessagePayload) Pack(buf *bytebuf.ByteBuf, client *Client) {
	// 写包长度占位
	_ = buf.WriteUInt32BE(0)
	// 写版本号
	_ = buf.WriteByte(1)
	// 写类型码
	_ = buf.WriteUInt16BE(p.PayloadCode)
	// 写messageID
	messageIdBytes := []byte(fastuuid.MustNewGenerator().Hex128())
	messageIdLen := byte(len(messageIdBytes))
	_ = buf.WriteByte(messageIdLen)
	_ = buf.WriteBytes(messageIdBytes)
	// 写通知分类
	_ = buf.WriteByte(p.Classifier)

	// 写merchantID
	merchantIdBytes := []byte(p.MerchantID)
	merchantIdLen := byte(len(merchantIdBytes))
	_ = buf.WriteByte(merchantIdLen)
	_ = buf.WriteBytes(merchantIdBytes)

	//写appID
	appIdBytes := []byte(p.AppID)
	appIdLen := byte(len(appIdBytes))
	_ = buf.WriteByte(appIdLen)
	_ = buf.WriteBytes(appIdBytes)
	var encFlag byte = 2
	// 加密消息,获取消息加密key CBC 模式
	encMessageKey := utils.GetMsgEncKey(client.appInfo.GetAppKey())
	encAesIV := utils.GetMsgEncAesIV(client.appInfo.GetAppKey())
	messageBody, _ := goEncrypt.AesCbcEncrypt(p.MessageBody, encMessageKey, encAesIV)
	if len(messageBody) > NEED_COMPRESS_SIZE {
		// 压缩
		messageBody = utils.Gzip(messageBody)
		encFlag = 3
	}
	// 写加密标志
	_ = buf.WriteByte(encFlag)
	// 写优先级
	_ = buf.WriteByte(p.Priority)

	// 写toUserID
	toUserIdBytes := []byte(p.ToUid)
	toUserIdLen := byte(len(toUserIdBytes))
	_ = buf.WriteByte(toUserIdLen)
	_ = buf.WriteBytes(toUserIdBytes)
	// 写MessageBody
	_ = buf.WriteBytes(messageBody)
	pktLen := buf.WriterIndex() - 4
	_ = buf.PutUInt32BE(0, uint32(pktLen))

}

func (p *PushMessagePayload) Unpack(buf *bytebuf.ByteBuf, _ *Client) {
	logger.Warnw("PushMessagePayload#Unpack 客户端无需实现", errors.New("客户端无需实现消息解码"))
}

// 创建一个新push通知
func NewPushMessagePayloadFromPushNotification(n PushNotification, classifier NotificationClassify, app *AppInfo) PushMessagePayload {
	messageId := fastuuid.MustNewGenerator().Hex128()
	messageBody, _ := json.Marshal(n.Body)
	return PushMessagePayload{
		PayloadCode: PushMessageCode,
		MessageID:   messageId,
		Classifier:  classifier,
		MerchantID:  app.GetMerchantID(),
		AppID:       app.GetAppID(),
		Priority:    n.Priority,
		ToUid:       n.ToUid,
		MessageBody: messageBody,
	}
}
