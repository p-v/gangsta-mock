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

var httpClient = &http.Client{Timeout: time.Second * 10}

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

func makePluginCall(request string, pluginLoc string, cb *callback, path string) {
	if cb.Delay != 0 {
		time.Sleep(time.Duration(cb.Delay) * time.Millisecond)
	}
	callbackHander := callbackMap[pluginLoc]
	handlerResponse := callbackHander(types.HandlerRequest{RequestBody: request, Path: path})

	callbackPath := cb.Path
	if handlerResponse.Path != "" {
		callbackPath = handlerResponse.Path
	}

	contentType := "application/json"
	if handlerResponse.ContentType != "" {
		contentType = handlerResponse.ContentType
	}

	req, err := http.NewRequest("POST", callbackPath, bytes.NewBuffer([]byte(handlerResponse.ResponseBody)))
	if err != nil {
		log.Printf("Error read request: %v", err)
		return
	}

	// set headers
	req.Header.Set("Content-Type", contentType)
	for header, value := range handlerResponse.Headers {
		req.Header.Set(header, value)
	}

	_, err = httpClient.Do(req)
	if err != nil {
		log.Printf("Error occurred while calling callback: %v", err)
	}
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
