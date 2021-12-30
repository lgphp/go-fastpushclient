package fastpushclient

type Listener = func(err error)
type SendListener = func(messageId string, err error)

type Listenerable interface {
	AddSendListener(l SendListener) *Client
	AddInitializedListener(l Listener) *Client
}
