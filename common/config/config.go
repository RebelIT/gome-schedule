package config

import (
	"flag"
	"log"
	"os"
	"strconv"
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
	StateTimeSec     int64 //time to check device schedule state in seconds
	FullMemory       bool  //use this for full in memory badgerDB cache
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
	configEnvironment(c)
	configFlags(c)

	App = c
	return
}

func configDefaults(c *Conf) {
	c.Name = "gome-schedule"
	c.StatAddr = ""
	c.DbPath = "badgerDatabase"
	c.AuthToken = "changeMePlease"
	c.ListenPort = "6661"
	c.CoreServiceToken = "changeMePlease"
	c.CoreServiceUrl = "http://localhost"
	c.CoreServicePort = "6660"
	c.GenerateSpec = false
	c.StateTimeSec = 30
	c.FullMemory = false

	return
}

func configEnvironment(c *Conf) {
	name := os.Getenv("SCHEDULE_NAME")
	statsd := os.Getenv("SCHEDULE_STATSD")
	dbPath := os.Getenv("SCHEDULE_DBPATH")
	authToken := os.Getenv("SCHEDULE_TOKEN")
	port := os.Getenv("SCHEDULE_PORT")
	fullMemory := os.Getenv("SCHEDULE_MEMORY")
	coreServiceToken := os.Getenv("SCHEDULE_CORE_TOKEN")
	coreServiceUrl := os.Getenv("SCHEDULE_CORE_URL")
	coreServicePort := os.Getenv("SCHEDULE_CORE_PORT")
	stateTimeSec := os.Getenv("SCHEDULE_STATE_TIME_SEC")

	if name != "" {
		c.Name = name
	}
	if statsd != "" {
		c.StatAddr = statsd
	}
	if dbPath != "" {
		c.DbPath = dbPath
	}
	if authToken != "" {
		c.AuthToken = authToken
	}
	if port != "" {
		c.ListenPort = port
	}
	if fullMemory != "" {
		c.FullMemory = true
	}
	if coreServiceToken != "" {
		c.CoreServiceToken = coreServiceToken
	}
	if coreServiceUrl != "" {
		c.CoreServiceUrl = coreServiceUrl
	}
	if coreServicePort != "" {
		c.CoreServicePort = coreServicePort
	}
	if stateTimeSec != "" {
		i, _ := strconv.Atoi(stateTimeSec)
		c.StateTimeSec = int64(i)
	}

	return
}

func configFlags(c *Conf) {
	flag.StringVar(&c.StatAddr, "statsd", c.StatAddr, "statsd address")
	flag.StringVar(&c.Name, "name", c.Name, "application name")
	flag.StringVar(&c.DbPath, "dbPath", c.DbPath, "path to local database")
	flag.StringVar(&c.AuthToken, "authToken", c.AuthToken, "app authentication token")
	flag.StringVar(&c.ListenPort, "port", c.ListenPort, "http listener http port")
	flag.StringVar(&c.CoreServiceToken, "coreToken", c.CoreServiceToken, "auth token for gome-core")
	flag.StringVar(&c.CoreServiceUrl, "coreUrl", c.CoreServiceUrl, "address of gome-core http://dnsname")
	flag.StringVar(&c.CoreServicePort, "corePort", c.CoreServicePort, "tcp port for gome-core")
	flag.BoolVar(&c.GenerateSpec, "generateSpec", c.GenerateSpec, "print the http spec to console")
	flag.Int64Var(&c.StateTimeSec, "stateTimeSec", c.StateTimeSec, "tine to check device states")
	flag.BoolVar(&c.FullMemory, "fullMemory", c.FullMemory, "run database with full memory cache mem > 1GB required")
	flag.Parse()
	return
}
