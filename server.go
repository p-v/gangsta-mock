package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type conf struct {
	Delay  int64   `yaml:"delay"`
	Routes []route `yaml:"routes"`
}

type route struct {
	Delay    int64     `yaml:"delay"`
	Path     string    `yaml:"path"`
	Method   string    `yaml:"method"`
	Code     int       `yaml:"code"`
	Handler  *handler  `yaml:"handler"`
	Callback *callback `yaml:"callback"`
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

func main() {
	c.getConf()
	m = make(map[string]route)

	router := fasthttprouter.New()

	for _, r := range c.Routes {
		m[r.Path] = r
		initializeHandler(r, router)
	}

	log.Printf("Starting server")
	fasthttp.ListenAndServe(":8080", router.Handler)
}
