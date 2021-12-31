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
	apiRoot               = "http://77.242.242.209:8080"
	pushListEndpoint      = "/biz/push/list"
)

func PrintHttpResponseError(err error, req interface{}) {
	logger.Warnw("请求失败:", err, "req", req)
}

// 根据AppID获取PushGateList
func (c *Client) getPushGateIpList() ([]pushGateAddress, error) {
	pushListReq := newPushListRequest(c.appInfo.appID)
	response, body, err := c.httpClient.Post(fmt.Sprintf("%s%s", apiRoot, pushListEndpoint), pushListReq)
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
				logger.Warnw("获取PushGateList失败", errors.New("Response.Code Not OK"), "req", pushListReq, "Code", code, "Message", message)
			}
		} else {
			logger.Warnw("获取PushGateList失败", errors.New("response.StatusCode Not OK"), "req", pushListReq, "StatusCode", response.StatusCode)
		}
	} else {
		PrintHttpResponseError(err, pushListReq)
	}
	return nil, errors.New("获取PushGateList失败")
}
