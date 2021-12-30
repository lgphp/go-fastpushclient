package fastpushclient

import (
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"push-sdk-go/client/utils"
	"push-sdk-go/logger"
	"strconv"
	"strings"
	"time"
)

const (
	Response_Success    = 0
	apiRoot             = "http://77.242.242.209:8080"
	pushListEndpoint    = "/biz/push/list"
	bindTokenEndpoint   = "/pushbind/bind"
	unBindTokenEndpoint = "/bizuser/unbind"
)

func PrintHttpResponseError(errs []error, req interface{}) {
	var errMsg string
	for _, e := range errs {
		errMsg += e.Error()
	}
	err := errors.New(errMsg)
	logger.Warnw("请求失败:", err, "req", req)
}

// 根据AppID获取PushGateList
func (c *Client) getPushGateIpList() ([]pushGateAddress, error) {
	pushListReq := newPushListRequest(c.appInfo.appID)
	signature := utils.MakeApiSignature(pushListReq, utils.GetApiKey(c.appInfo.appKey))
	request := gorequest.New()
	send := request.Post(fmt.Sprintf("%s%s", apiRoot, pushListEndpoint)).Timeout(time.Second * 5).Send(pushListReq)
	response, body, errs := send.Set("API-SIGNATURE", signature).Set("APP-ID", c.appInfo.appID).End()
	if errs == nil {
		if response.StatusCode == http.StatusOK {
			res, _ := simplejson.NewJson([]byte(body))
			code, _ := res.Get("code").Int()
			if code == Response_Success {
				data := res.Get("data").MustStringArray()
				pushList := make([]pushGateAddress, 0)
				for _, p := range data {
					ip := strings.Split(p, ":")
					port, _ := strconv.Atoi(ip[1])
					address := pushGateAddress{
						IP:   ip[0],
						Port: port,
					}
					pushList = append(pushList, address)
				}
				return pushList, nil
			} else {
				message, _ := res.Get("message").String()
				logger.Warnw("获取PushGateList失败", errors.New("Response.Code Not OK"), "req", pushListReq, "Code", code, "Message", message)
			}
		} else {
			logger.Warnw("获取PushGateList失败", errors.New("response.StatusCode Not OK"), "req", pushListReq, "StatusCode", response.StatusCode)
		}
	} else {
		PrintHttpResponseError(errs, pushListReq)
	}
	return nil, errors.New("获取PushGateList失败")
}

// 绑定Token ,  暂时不提供
func (c *Client) requestBindToken(userID, channel, token string) {
	bindReq := newBindTokenRequest(c.appInfo.appID, userID, channel, token)
	signature := utils.MakeApiSignature(bindReq, utils.GetApiKey(c.appInfo.appKey))
	send := gorequest.New().Set("API-SIGNATURE", signature).Set("APP-ID", c.appInfo.appID).Post(fmt.Sprintf("%s%s", apiRoot, bindTokenEndpoint)).Timeout(time.Second * 5).Send(bindReq)
	response, body, errs := send.End()
	if errs != nil {
		PrintHttpResponseError(errs, bindReq)
	} else {
		if response.StatusCode == http.StatusOK {
			res, _ := simplejson.NewJson([]byte(body))
			code, _ := res.Get("code").Int()
			if code == Response_Success {
				logger.Infow("BindToken 成功")
			} else {
				message, _ := res.Get("message").String()
				logger.Warnw("BindToken 失败", errors.New("Response.Code Not OK"), "req", bindReq, "Code", code, "Message", message)
			}
		} else {
			logger.Warnw("BindToken 失败", errors.New("response.StatusCode Not OK"), "req", bindReq, "StatusCode", response.StatusCode)
		}
	}
}

// 解绑token ， 暂时不提供
func (c *Client) requestUnBindToken(userID, channel, token string) {
	bindReq := newBindTokenRequest(c.appInfo.appID, userID, channel, token)
	signature := utils.MakeApiSignature(bindReq, utils.GetApiKey(c.appInfo.appKey))
	send := gorequest.New().Set("API-SIGNATURE", signature).Set("APP-ID", c.appInfo.appID).Post(fmt.Sprintf("%s%s", apiRoot, unBindTokenEndpoint)).Timeout(time.Second * 5).Send(bindReq)
	response, body, errs := send.End()
	if errs != nil {
		PrintHttpResponseError(errs, bindReq)
	} else {
		if response.StatusCode == http.StatusOK {
			res, _ := simplejson.NewJson([]byte(body))
			code, _ := res.Get("code").Int()
			if code == Response_Success {
				logger.Infow("UnBindToken  成功")
			} else {
				message, _ := res.Get("message").String()
				logger.Warnw("UnBindToken 失败", errors.New("Response.Code Not OK"), "req", bindReq, "Code", code, "Message", message)
			}
		} else {
			logger.Warnw("UnBindToken 失败", errors.New("response.StatusCode Not OK"), "req", bindReq, "StatusCode", response.StatusCode)
		}
	}
}
