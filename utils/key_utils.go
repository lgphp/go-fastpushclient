package utils

import "github.com/lgphp/go-fastpushclient/bytebuf"

// authKey
func GetAuthKey(appKey []byte) []byte {
	buf, _ := bytebuf.NewByteBuf(appKey)
	authKey := make([]byte, 16)
	buf.SkipBytes(16)
	_, _ = buf.ReadBytes(authKey)
	defer func() {
		buf.Release()
	}()
	return authKey
}

// apiKey
func GetApiKey(appKey []byte) []byte {
	buf, _ := bytebuf.NewByteBuf(appKey)
	apiKey := make([]byte, 16)
	_, _ = buf.ReadBytes(apiKey)
	defer func() {
		buf.Release()
	}()
	return apiKey
}

// message encryption Key
func GetMsgEncKey(appKey []byte) []byte {
	buf, _ := bytebuf.NewByteBuf(appKey)
	encKey := make([]byte, 16)
	buf.SkipBytes(32)
	_, _ = buf.ReadBytes(encKey)
	defer func() {
		buf.Release()
	}()
	return encKey
}

// AES IV
func GetMsgEncAesIV(appKey []byte) []byte {

	buf, _ := bytebuf.NewByteBuf(appKey)
	iv := make([]byte, 16)
	buf.SkipBytes(48)
	_, _ = buf.ReadBytes(iv)
	defer func() {
		buf.Release()
	}()
	return iv
}
