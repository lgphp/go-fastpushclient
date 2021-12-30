package utils

import "github.com/lgphp/go-fastpushclient/bytebuf"

// 获取authKey
func GetAuthKey(appKey []byte) []byte {
	buf, _ := bytebuf.NewByteBuf(appKey)
	authKey := make([]byte, 16)
	buf.SkipBytes(16)
	_, _ = buf.ReadBytes(authKey)
	buf.Release()
	return authKey
}

// 获取apiKey
func GetApiKey(appKey []byte) []byte {
	buf, _ := bytebuf.NewByteBuf(appKey)
	apiKey := make([]byte, 16)
	_, _ = buf.ReadBytes(apiKey)
	buf.Release()
	return apiKey
}

// 获取消息加密Key
func GetMsgEncKey(appKey []byte) []byte {
	buf, _ := bytebuf.NewByteBuf(appKey)
	encKey := make([]byte, 16)
	buf.SkipBytes(32)
	_, _ = buf.ReadBytes(encKey)
	buf.Release()
	return encKey
}

// 获取AES IV
func GetMsgEncAesIV(appKey []byte) []byte {
	buf, _ := bytebuf.NewByteBuf(appKey)
	iv := make([]byte, 16)
	buf.SkipBytes(48)
	_, _ = buf.ReadBytes(iv)
	buf.Release()
	return iv
}
