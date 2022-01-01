package fastpushclient

import (
	"github.com/go-netty/go-netty"
	"github.com/lgphp/go-fastpushclient/logger"
)

type ExceptionHandler struct {
	netty.ExceptionHandler
	name string
}

func newExceptionHandler(name string) *ExceptionHandler {
	return &ExceptionHandler{
		name: name,
	}
}

func (h *ExceptionHandler) HandleException(ctx netty.ExceptionContext, ex netty.Exception) {
	logger.Warnw("exception handler:", ex)
}
