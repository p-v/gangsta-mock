package main

import (
	"gangsta-mock/types"
)

type handler string

func (h handler) Handle(request types.Request) types.Response {
	return types.Response{
		ResponseBody: "{\"message\": \"What up gangsta\"}",
	}
}

var Handler handler
