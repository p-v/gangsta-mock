package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"plugin"
	"time"
)

type CustomCallback func(request string) string

type callback struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
	Body   string `yaml:"body"`
	Delay  int64  `yaml:"delay"`
}

type CallbackHandler interface {
	Handle(request string) string
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
	resp := callbackHander(request)
	http.Post(cb.Path, "application/json", bytes.NewBuffer([]byte(resp)))
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

	handler, ok := plugHandler.(CallbackHandler)
	if !ok {
		log.Printf("Bad type")

	}

	callbackMap[pluginLoc] = handler.Handle
}
