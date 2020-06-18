package server

import (
	"gangsta-mock/types"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"plugin"
	"time"
)

type handler struct {
	Plugin   string `yaml:"plugin"`
	Response string `yaml:"response"`
}

func nonPluginHandler(callbackFunc CallbackFunc, routeData route) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		delay := routeData.Delay
		if delay == 0 {
			delay = c.Delay
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)
		ctx.SetStatusCode(routeData.Code)

		if routeData.Handler != nil {
			ctx.Write([]byte(routeData.Handler.Response))
		}

		callCb(routeData.Callback, callbackFunc, ctx, path)
	}
}

func pluginHandler(handlerFunc types.HandlerFunc, callbackFunc CallbackFunc, r route) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		pathParams := make(map[string]interface{})

		ctx.VisitUserValues(func(key []byte, value interface{}) {
			pathParams[string(key)] = value
		})
		resp := handlerFunc(types.Request{
			Path:        r.Path,
			Method:      r.Method,
			Body:        string(ctx.PostBody()),
			QueryParams: ctx.QueryArgs(),
			PathParams:  pathParams,
		})

		ctx.SetStatusCode(resp.Code)
		ctx.Write([]byte(resp.ResponseBody))

		callCb(r.Callback, callbackFunc, ctx, r.Path)
	}
}

func callCb(callback *callback, callbackFunc CallbackFunc, ctx *fasthttp.RequestCtx, path string) {
	if callback != nil {
		if callback.Plugin != "" {
			go makePluginCall(string(ctx.PostBody()), callbackFunc, path, ctx.QueryArgs())
		} else {
			go makeHttpCall(callback)
		}
	}
}

func initializeHandlerPlugin(hn *handler) types.HandlerFunc {
	if hn == nil || hn.Plugin == "" {
		return nil
	}

	pluginLoc := hn.Plugin

	handlerPlugin, err := plugin.Open(pluginLoc)
	if err != nil {
		log.Printf("%v", err)
		os.Exit(1)
	}

	plugHandler, err := handlerPlugin.Lookup("Handler")
	if err != nil {
		log.Printf("%v", err)
		os.Exit(1)
	}

	rhn, ok := plugHandler.(types.RequestHandler)
	if !ok {
		log.Printf("Bad type")
		log.Printf("Plugin Loc %s", pluginLoc)
		os.Exit(1)
	}

	return rhn.Handle
}

func gangstaHandler(handlerFunc types.HandlerFunc, callbackFunc CallbackFunc, r route) fasthttp.RequestHandler {
	if handlerFunc == nil {
		return nonPluginHandler(callbackFunc, r)
	}
	return pluginHandler(handlerFunc, callbackFunc, r)
}

func initializeHandler(r route, router *fasthttprouter.Router) {
	callbackFunc := initializeCallbackPlugin(r.Callback)
	handlerFunc := initializeHandlerPlugin(r.Handler)

	method := r.Method
	if method == "" {
		method = "GET"
	}

	// Register handler
	router.Handle(method, r.Path, gangstaHandler(handlerFunc, callbackFunc, r))
}
