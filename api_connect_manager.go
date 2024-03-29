package fastpushclient

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/lgphp/go-fastpushclient/logger"
	"net/http"
	"strconv"
	"strings"
)

const (
	HTTP_RESPONSE_CODE_OK = 0
	pushListEndpoint      = "/biz/push/list"
)

func PrintHttpResponseError(err error, req interface{}) {
	logger.Warnw("Request failed:", err, "req", req)
}

// 根据AppID获取PushGateList
func (c *Client) getPushGateIpList() ([]pushGateAddress, error) {
	pushListReq := newPushListRequest(c.appInfo.appID)
	response, body, err := c.httpClient.Post(fmt.Sprintf("%s%s", c.apiRoot, pushListEndpoint), pushListReq)
	if err == nil {
		if response.StatusCode == http.StatusOK {
			res, _ := simplejson.NewJson(body)
			code, _ := res.Get("code").Int()
			if code == HTTP_RESPONSE_CODE_OK {
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
				logger.Warnw("Request PushGate server failed", errors.New("Response.Code Not OK"), "req", pushListReq, "Code", code, "Message", message)
			}
		} else {
			logger.Warnw("Request PushGate server failed", errors.New("Response.StatusCode Not OK"), "req", pushListReq, "StatusCode", response.StatusCode)
		}
	} else {
		PrintHttpResponseError(err, pushListReq)
	}
	return nil, errors.New(fmt.Sprintf("Request PushGate server failed: %s", err))
}
