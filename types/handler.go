package types

import "github.com/valyala/fasthttp"

type Request struct {
	Path        string
	Method      string
	Body        string
	QueryParams *fasthttp.Args
	PathParams  map[string]interface{}
}

type Response struct {
	ContentType  string
	ResponseBody string
	Headers      map[string]string
	Code         int
}

type RequestHandler interface {
	Handle(request Request) Response
}

type HandlerFunc func(request Request) Response
