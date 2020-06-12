package types

type Request struct {
	Path   string
	Method string
	Body   string
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
