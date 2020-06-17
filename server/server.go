package server

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
var server *fasthttp.Server

func (c *conf) getConf(configFile string) *conf {
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Unable to fetch configuration #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}

func StartServer(configFile string) {
	c.getConf(configFile)

	router := fasthttprouter.New()

	for _, r := range c.Routes {
		initializeHandler(r, router)
	}

	log.Printf("Starting server")
	server = &fasthttp.Server{Handler: router.Handler}
	server.ListenAndServe(":8080")
}

func StopServer() {
	server.Shutdown()
}
