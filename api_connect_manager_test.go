package fastpushclient

import (
	"github.com/lgphp/go-fastpushclient/logger"
	"testing"
)

// api 测试
var (
	// 测试环境的相关配置
	TEST_ENV_MERCHANT_ID = "a127297f117c4a3fb095a15443bc96fc"
	TEST_ENV_APP_ID      = "b4722bb12f30485582fb3e3a5c6157c6"
	TEST_ENV_APP_KEY     = "NUhONBRTxPxFtFkH78P9AJ2EDUJ1EeoaFzGVoJUz5BcYFtqiag0baRw61y1ycoZaYpkxp9BC08K2F8h2II4tyQ"
	TEST_ENV_USER_ID     = "8613810654610"

	//本地环境相关配置
	LOCAL_ENV_MERCHANT_ID = "3527348b4ff04e988f3fadd7f1e4f155"
	LOCAL_ENV_APP_ID      = "4fb367eabe8c45f2a1b6714c6e40fd19"
	LOCAL_ENV_APP_KEY     = "fCCHFkVtMk7sf5XQmfyTpPvuVH0PKmUd559HUtTlDxCBg5y4P3SAcxnAgxCG/AuRO0y//ZbgwRQg1wCJGOWw/w=="
	LOCAL_ENV_USER_ID     = "97158000000"
)

func TestHttpClient(t *testing.T) {
	info := NewAppInfo(TEST_ENV_MERCHANT_ID,
		TEST_ENV_APP_ID,
		TEST_ENV_APP_KEY)
	//NewFastLivePushHttpClient(info)
	_, _ = NewFastLivePushClient(info).AddInitializedListener(func(err error) {
		logger.Warnw("sss", err)
	}).AddNotificationStatusListener(func(messageId, toUserId, appId string, statusCode uint32, statusText string) {

	}).BuildConnect()
}
