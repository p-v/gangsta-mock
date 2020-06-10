package types

type HandlerRequest struct {
	Path        string
	RequestBody string
}

type HandlerResponse struct {
	ContentType  string
	Path         string
	ResponseBody string
}

type CallbackHandler interface {
	Handle(request HandlerRequest) HandlerResponse
}
