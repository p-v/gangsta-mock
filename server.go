package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

type conf struct {
	Delay          int64   `yaml:"delay"`
	Routes         []route `yaml:"routes"`
	CallbackPlugin string  `yaml:"callbackPlugin"`
}

type route struct {
	Delay          int64     `yaml:"delay"`
	Path           string    `yaml:"path"`
	Body           string    `yaml:"body"`
	Code           int       `yaml:"code"`
	Response       string    `yaml:"response"`
	Callback       *callback `yaml:"callback"`
	CallbackPlugin string    `yaml:"callbackPlugin"`
}

var c conf
var m map[string]route

func (c *conf) getConf() *conf {
	yamlFile, err := ioutil.ReadFile("gangsta.yml")
	if err != nil {
		log.Printf("yaml file not retrieved #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	routeData := m[string(ctx.Path())]
	delay := routeData.Delay
	if delay == 0 {
		delay = c.Delay
	}
	time.Sleep(time.Duration(delay) * time.Millisecond)
	ctx.SetStatusCode(routeData.Code)
	ctx.Write([]byte(routeData.Response))
	if routeData.Callback != nil {
		go makeHttpCall(routeData.Response, routeData.Callback)
	}

	if routeData.CallbackPlugin != "" {
		go makePluginCall(routeData.Body, routeData.CallbackPlugin, routeData.Callback)
	}
}

func main() {
	c.getConf()
	m = make(map[string]route)

	for _, r := range c.Routes {
		m[r.Path] = r
		initializePlugin(r.CallbackPlugin)
	}

	// Initialize generic plugin
	initializePlugin(c.CallbackPlugin)

	router := fasthttprouter.New()
	router.GET("/", fastHTTPHandler)

	log.Printf("Starting server")
	fasthttp.ListenAndServe(":8080", fastHTTPHandler)
}
