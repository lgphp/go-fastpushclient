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
