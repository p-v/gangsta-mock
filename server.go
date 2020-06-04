package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type conf struct {
	Routes []route `yaml:"routes"`
}

type route struct {
	Path     string `yaml:"path"`
	Response string `yaml:"response"`
}

var m map[string]route

func (c *conf) getConf() *conf {
	fmt.Println("Reading...")
	yamlFile, err := ioutil.ReadFile("data.yml")
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
	ctx.SetBody([]byte(routeData.Response))
}

func main() {
	var c conf
	c.getConf()
	m = make(map[string]route)

	for _, r := range c.Routes {
		m[r.Path] = r
	}

	fasthttp.ListenAndServe(":8081", fastHTTPHandler)
}
