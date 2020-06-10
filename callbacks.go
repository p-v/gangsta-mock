package main

import (
	"bytes"
	"gangsta-mock/types"
	"log"
	"net/http"
	"os"
	"plugin"
	"time"
)

type CustomCallback func(request types.HandlerRequest) types.HandlerResponse

type callback struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
	Body   string `yaml:"body"`
	Delay  int64  `yaml:"delay"`
}

var callbackMap map[string]CustomCallback

func makeHttpCall(response string, cb *callback) {
	if cb.Delay != 0 {
		time.Sleep(time.Duration(cb.Delay) * time.Millisecond)
	}

	http.Post(cb.Path, "application/json", bytes.NewBuffer([]byte(cb.Body)))
}

func makePluginCall(request string, pluginLoc string, cb *callback) {
	callbackHander := callbackMap[pluginLoc]
	handlerResponse := callbackHander(types.HandlerRequest{RequestBody: request, Path: cb.Path})

	callbackPath := cb.Path
	if handlerResponse.Path != "" {
		callbackPath = handlerResponse.Path
	}

	contentType := "application/json"
	if handlerResponse.ContentType != "" {
		contentType = handlerResponse.ContentType
	}
	http.Post(callbackPath, contentType, bytes.NewBuffer([]byte(handlerResponse.ResponseBody)))
}

func initializePlugin(pluginLoc string) {
	if callbackMap == nil {
		callbackMap = make(map[string]CustomCallback)
	}

	if pluginLoc == "" {
		return
	}

	callbackPlugin, err := plugin.Open(pluginLoc)
	if err != nil {
		log.Printf("%v", err)
		os.Exit(1)
	}

	plugHandler, err := callbackPlugin.Lookup("Handler")
	if err != nil {
		log.Printf("%v", err)
		os.Exit(1)
	}

	handler, ok := plugHandler.(types.CallbackHandler)
	if !ok {
		log.Printf("Bad type")
		log.Printf("Plugin Loc %s", pluginLoc)
		os.Exit(1)
	}

	callbackMap[pluginLoc] = handler.Handle
}
