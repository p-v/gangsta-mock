package main

type CallbackHandler interface {
	Handle(request string) string
}

type handler string

func (h handler) Handle(request string) string {
	return "{\"message\": \"Wad up gangsta\"}"
}

var Handler handler
