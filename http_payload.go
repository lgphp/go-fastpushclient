package fastpushclient

type HttpResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Request interface{}

type HttpRequest struct {
	Request
	Signature string `json:"signature,omitempty"`
}

func newHttpRequest(req Request, signature string) *HttpRequest {
	return &HttpRequest{
		req,
		signature,
	}
}

type pushList struct {
	AppID string `json:"appID,omitempty"`
}

func newPushListRequest(appID string) pushList {
	return pushList{
		AppID: appID,
	}
}

type bindToken struct {
	UserID  string `json:"userID,omitempty"`
	AppID   string `json:"appID,omitempty"`
	Channel string `json:"channel,omitempty"`
	Token   string `json:"token,omitempty"`
}

func newBindTokenRequest(appID, userID, channel, token string) bindToken {
	return bindToken{
		UserID:  userID,
		AppID:   appID,
		Channel: channel,
		Token:   token,
	}
}
