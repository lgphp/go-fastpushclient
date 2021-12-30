package fastpushclient

type MessageBody struct {
	Title string            `json:"title"`
	Body  string            `json:"body"`
	Data  map[string]string `json:"data"`
}

func NewMessageBody(title, body string, attachmentData map[string]string) MessageBody {
	return MessageBody{
		Title: title,
		Body:  body,
		Data:  attachmentData,
	}
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
