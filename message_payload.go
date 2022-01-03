package fastpushclient

import "fmt"

type MessageBody struct {
	Title string            `json:"title,omitempty"`
	Body  string            `json:"body,omitempty"`
	Data  map[string]string `json:"data,omitempty"`
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
	toUid       string
	priority    MessagePrior
	messageBody MessageBody
}

// 创建一条push通知
func NewPushNotification(toUid string, priority MessagePrior, body MessageBody) PushNotification {
	return PushNotification{
		toUid:       toUid,
		priority:    priority,
		messageBody: body,
	}

}
