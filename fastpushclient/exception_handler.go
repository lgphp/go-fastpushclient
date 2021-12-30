package fastpushclient

type ExceptionHandler struct {
	netty.ExceptionHandler
	name string
}

func newExceptionHandler(name string) ExceptionHandler {
	return ExceptionHandler{
		name: name,
	}
}

func (h ExceptionHandler) HandleException(ctx netty.ExceptionContext, ex netty.Exception) {
	//logger.Warnw("异常处理器获得异常" , ex)
}
