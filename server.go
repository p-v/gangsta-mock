package main

import (
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

type conf struct {
	Delay  int64   `yaml:"delay"`
	Routes []route `yaml:"routes"`
}

type route struct {
	Delay    int64  `yaml:"delay"`
	Path     string `yaml:"path"`
	Response string `yaml:"response"`
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
	ctx.SetStatusCode(fasthttp.StatusOK)
	routeData := m[string(ctx.Path())]
	delay := routeData.Delay
	if delay == 0 {
		delay = c.Delay
	}
	time.Sleep(time.Duration(delay) * time.Millisecond)
	ctx.Write([]byte(routeData.Response))
}

func main() {
	c.getConf()
	m = make(map[string]route)

	for _, r := range c.Routes {
		m[r.Path] = r
	}

	log.Printf("Starting server")
	fasthttp.ListenAndServe(":8081", fastHTTPHandler)
}
