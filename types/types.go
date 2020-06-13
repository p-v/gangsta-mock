package types

import "github.com/valyala/fasthttp"

type HandlerRequest struct {
	Path        string
	RequestBody string
	QueryParams *fasthttp.Args
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
