package schedule

import (
	"github.com/rebelit/gome-schedule/common/config"
	"log"
	"time"
)

var retry = 5 //seconds

func Runner() {
	log.Println("INFO: starting schedule runners")

	for {
		go processStateSchedules()
		go processToggleSchedules()

		time.Sleep(time.Second * time.Duration(config.App.StateTimeSec))
	}
}
