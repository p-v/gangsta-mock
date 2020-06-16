package server

import (
	"bytes"
	"gangsta-mock/types"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
	"os"
	"plugin"
	"time"
)

var httpClient = &http.Client{Timeout: time.Second * 10}

type CallbackFunc func(request types.HandlerRequest) types.HandlerResponse

type callback struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
	Body   string `yaml:"body"`
	Delay  int64  `yaml:"delay"`
	Plugin string `yaml:"plugin"`
}

func makeHttpCall(cb *callback) {
	if cb.Delay != 0 {
		time.Sleep(time.Duration(cb.Delay) * time.Millisecond)
	}

	http.Post(cb.Path, "application/json", bytes.NewBuffer([]byte(cb.Body)))
}

func makePluginCall(request string, callbackFunc CallbackFunc, path string, queryParams *fasthttp.Args) {
	handlerResponse := callbackFunc(types.HandlerRequest{
		RequestBody: request,
		Path:        path,
		QueryParams: queryParams,
	})

	contentType := "application/json"
	if handlerResponse.ContentType != "" {
		contentType = handlerResponse.ContentType
	}

	req, err := http.NewRequest("POST", handlerResponse.Path, bytes.NewBuffer([]byte(handlerResponse.ResponseBody)))
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

func initializeCallbackPlugin(cb *callback) CallbackFunc {
	if cb == nil || cb.Plugin == "" {
		return nil
	}

	pluginLoc := cb.Plugin

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

	return handler.Handle
}
