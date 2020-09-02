package config

import (
	"flag"
	"log"
)

type Conf struct {
	StatAddr         string //127.0.0.1:8125
	SlackWebhook     string //https://hooks.slack.com/services/<ID>
	DbPath           string
	Name             string
	AuthToken        string
	ListenPort       string
	CoreServiceToken string
	CoreServiceUrl   string
	CoreServicePort  string
	GenerateSpec     bool
}

var App *Conf

const (
	STATEDBNAME  = "state"
	TOGGLEDBNAME = "toggle"
)

func Runtime() {
	log.Printf("INFO: loading runtime configuration")
	c := &Conf{}
	configDefaults(c)
	configFlags(c)

	App = c
	return
}

func configDefaults(c *Conf) {
	c.Name = "gome-core"
	c.StatAddr = ""
	c.DbPath = "badgerDatabase"
	c.AuthToken = "changeMePlease"
	c.ListenPort = "6661"
	c.CoreServiceToken = "changeMePlease"
	c.CoreServiceUrl = ""
	c.CoreServicePort = ""
	c.GenerateSpec = false
	return
}

func configFlags(c *Conf) {
	flag.StringVar(&c.StatAddr, "statsd", c.StatAddr, "statsd address")
	flag.StringVar(&c.Name, "name", c.Name, "application name")
	flag.StringVar(&c.DbPath, "dbPath", c.DbPath, "path to local database")
	flag.StringVar(&c.SlackWebhook, "slackWebhook", c.SlackWebhook, "slack webhook url")
	flag.StringVar(&c.AuthToken, "authToken", c.AuthToken, "app authentication token")
	flag.StringVar(&c.ListenPort, "port", c.ListenPort, "http listener http port")
	flag.BoolVar(&c.GenerateSpec, "generateSpec", c.GenerateSpec, "print the http spec to console")
	flag.Parse()
	return
}
