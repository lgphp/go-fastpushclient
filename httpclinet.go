package fastpushclient

import (
	"errors"
	"github.com/lgphp/go-fastpushclient/utils"
	"github.com/parnurzeal/gorequest"
	"time"
)

type HTTPClient struct {
	info AppInfo
}

func NewFastLivePushHttpClient(info AppInfo) HTTPClient {
	return HTTPClient{
		info: info,
	}
}

func fmtErrors(errs []error) error {
	var errMsg string
	for _, e := range errs {
		errMsg += e.Error()
	}
	return errors.New(errMsg)
}

func (hc *HTTPClient) Post(url string, reqBody interface{}) (resp gorequest.Response, data []byte, err error) {
	signature := utils.MakeApiSignature(reqBody, utils.GetApiKey(hc.info.appKey))
	send := gorequest.New().Post(url).Timeout(time.Second * 5).Send(reqBody)
	response, body, errs := send.Set("API-SIGNATURE", signature).Set("APP-ID", hc.info.appID).End()
	return response, []byte(body), fmtErrors(errs)
}
