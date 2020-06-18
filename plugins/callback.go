package main

import (
	. "gangsta-mock/types"
)

type handler string

func (h handler) Handle(request HandlerRequest) HandlerResponse {
	return HandlerResponse{
		ResponseBody: "{\"message\": \"What up gangsta\"}",
		Path:         "http://localhost:8000/event",
	}
}

var Handler handler
