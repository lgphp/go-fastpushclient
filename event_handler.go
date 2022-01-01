package fastpushclient

import "github.com/go-netty/go-netty"

type EventHandler struct {
	netty.EventHandler
	name string
}

func newEventHandler(name string) *EventHandler {
	return &EventHandler{
		name: name,
	}
}

func (h *EventHandler) HandleEvent(ctx netty.EventContext, event netty.Event) {

}
