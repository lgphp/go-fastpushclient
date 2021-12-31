package fastpushclient

import "fmt"

type MessageBody struct {
	Title string            `json:"title"`
	Body  string            `json:"body"`
	Data  map[string]string `json:"data"`
}

func NewMessageBody(title, body string, attachmentData map[string]string) (MessageBody, error) {
	if len(title) == 0 {
		return MessageBody{}, fmt.Errorf("title must be specified")
	}
	if len(body) == 0 {
		return MessageBody{}, fmt.Errorf("body must be specified")
	}
	return MessageBody{
		Title: title,
		Body:  body,
		Data:  attachmentData,
	}, nil
}

type PushNotification struct {
	ToUid    string
	Priority MessagePrior
	Body     MessageBody
}

// 创建一条push通知
func NewPushNotification(toUid string, priority MessagePrior, body MessageBody) PushNotification {
	return PushNotification{
		ToUid:    toUid,
		Priority: priority,
		Body:     body,
	}

}
