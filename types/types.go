package types

type HandlerRequest struct {
	Path        string
	RequestBody string
}

type HandlerResponse struct {
	ContentType  string
	Path         string
	ResponseBody string
	Headers      map[string]string
}

type CallbackHandler interface {
	Handle(request HandlerRequest) HandlerResponse
}
