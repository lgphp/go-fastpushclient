package fastpushclient

type InitializedListener = func(err error)
type ClientSendListener = func(messageId string, err error)
type NotificationStatusListener = func(messageId, toUserId, appId string, statusCode uint32, statusText string)

type Listenerable interface {
	AddSendListener(l ClientSendListener) *Client
	AddInitializedListener(l InitializedListener) *Client
	AddNotificationStatusListener(l NotificationStatusListener) *Client
}
