package schedule

import (
	"encoding/json"
	"fmt"
	"github.com/rebelit/gome-schedule/common/config"
	"github.com/rebelit/gome-schedule/database"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var weekend = []string{"saturday", "sunday"}
var weekday = []string{"monday", "tuesday", "wednesday", "thursday", "friday"}

func newDatabase(dbName string) error {
	db, err := database.Open(fmt.Sprintf("%s/%s", config.App.DbPath, dbName))
	if err != nil {
		return err
	}

	_, err = db.GetAllKeys()
	if err != nil {
		return err
	}

	return nil
}

func deleteSchedule(name string, dbName string) error {
	db, err := database.Open(fmt.Sprintf("%s/%s", config.App.DbPath, dbName))
	if err != nil {
		return err
	}

	if err := db.Delete(name); err != nil {
		return err
	}

	return nil
}

func newSchedule(name string, schedule []byte, dbName string) error {
	db, err := database.Open(fmt.Sprintf("%s/%s", config.App.DbPath, dbName))
	if err != nil {
		return err
	}

	if err := db.Set(name, schedule); err != nil {
		return err
	}

	return nil
}

func updateSchedule(name string, schedule []byte, dbName string) error {
	if err := deleteSchedule(name, dbName); err != nil {
		return err
	}
	if err := newSchedule(name, schedule, dbName); err != nil {
		return err
	}
	return nil
}

func requiredPathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func createRequiredPath(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}

	return nil
}

func respondHttpBody(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("ERROR: devices handler response, %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			//ToDo: Metric
			return
		}
	}

	//ToDo: Metric
	return
}

func currentTimeSplit() (strTime string, intTime int, weekday string, now time.Time) {
	Now := time.Now()
	NowMinute := Now.Minute()
	NowHour := Now.Hour()
	NowDay := Now.Weekday()

	sTime := ""
	singleMinute := inBetween(NowMinute, 0, 9)
	if singleMinute {
		sTime = strconv.Itoa(NowHour) + "0" + strconv.Itoa(NowMinute)
	} else {
		sTime = strconv.Itoa(NowHour) + strconv.Itoa(NowMinute)
	}

	iTime, _ := strconv.Atoi(sTime)
	day := strings.ToLower(NowDay.String())

	return sTime, iTime, day, Now
}

func inBetween(i, min, max int) bool {
	if (i >= min) && (i <= max) {
		return true
	} else {
		return false
	}
}

func inBetweenReverse(i, min, max int) bool {
	if (i >= min) && (i <= max) {
		return false
	} else {
		return true
	}
}

func boolToState(boolState bool) string {
	if boolState {
		return "on"
	}
	return "off"
}

func isToday(day string) bool {
	switch day {
	case "daily": //process every day
		log.Printf("DEBUG: schedule for %s days, processing", day)
		return true

	case "weekend": //process weekends only
		return isWeekend()

	case "weekday": //procecss weekdays only
		return isWeekday()

	default: //process day by day
		_, _, today, _ := currentTimeSplit()
		if today == day {
			log.Printf("DEBUG: schedule for %s and it's %s, processing", day, today)
			return true
		}
		log.Printf("DEBUG: schedule for %s and it's %s, ignore it", day, today)
		return false
	}
}

func isWeekend() bool {
	_, _, today, _ := currentTimeSplit()
	for _, d := range weekend {
		if today == d {
			return true
		}
	}
	return false
}

func isWeekday() bool {
	_, _, today, _ := currentTimeSplit()
	for _, d := range weekday {
		if today == d {
			return true
		}
	}
	return false
}
