package fastpushclient

import "fmt"

type MessageBody struct {
	title string
	body  string
	data  map[string]string
}

func NewMessageBody(title, body string, attachmentData map[string]string) (MessageBody, error) {
	if len(title) == 0 {
		return MessageBody{}, fmt.Errorf("title must be specified")
	}
	if len(body) == 0 {
		return MessageBody{}, fmt.Errorf("body must be specified")
	}
	return MessageBody{
		title: title,
		body:  body,
		data:  attachmentData,
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
