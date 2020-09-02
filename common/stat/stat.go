package stat

import (
	"github.com/rebelit/gome-schedule/common/config"
	"gopkg.in/alexcesaro/statsd.v2"
	"log"
)

const (
	STATEOK      = "ok"
	STATEFAILURE = "error"
	HTTPIN       = "inbound"
	HTTPOUT      = "outbound"
)

func disabled() bool {
	if config.App.StatAddr == "" {
		return true
	}
	return false
}

// Generic metrics types using influx line protocol
func counter(measurement string, tags statsd.Option) {
	if disabled() {
		return
	}

	addrOpt := statsd.Address(config.App.StatAddr)
	fmtOpt := statsd.TagsFormat(statsd.InfluxDB)
	s, err := statsd.New(addrOpt, fmtOpt, tags)
	if err != nil {
		log.Printf("ERROR: sending %s counter %s", measurement, err)
	}
	defer s.Close()

	s.Increment(config.App.Name + "." + measurement)
	return
}

func gauge(measurement string, tags statsd.Option, value int) {
	if disabled() {
		return
	}

	addrOpt := statsd.Address(config.App.StatAddr)
	fmtOpt := statsd.TagsFormat(statsd.InfluxDB)
	s, err := statsd.New(addrOpt, fmtOpt, tags)
	if err != nil {
		log.Printf("ERROR: sending %s gauge %s", measurement, err)
	}
	defer s.Close()

	s.Gauge(config.App.Name+"."+measurement, value)
	return
}

//Public package specific metrics
func Database(action string, state string) {
	measurement := "database"
	tags := statsd.Tags("action", action, "state", state)

	counter(measurement, tags)
	return
}

func Notify(service string, state string, statusCode int) {
	measurement := "notify"
	tags := statsd.Tags("service", service, "state", state)

	gauge(measurement, tags, statusCode)
	return
}

func Http(method string, direction string, url string, statusCode int) {
	measurement := "http"
	tags := statsd.Tags("direction", direction, "url", url)

	gauge(measurement, tags, statusCode)
	return
}
