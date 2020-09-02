package schedule

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rebelit/gome-schedule/common/config"
	"github.com/rebelit/gome-schedule/database"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

func unmarshalToggleSchedule(data []byte) (schedule DeviceToggle, error error) {
	if err := json.Unmarshal(data, &schedule); err != nil {
		return schedule, err
	}

	return
}

func toggleDb(wg *sync.WaitGroup) {
	defer wg.Done()

	if err := newDatabase(config.TOGGLEDBNAME); err != nil {
		log.Fatalf("FATAL: %s failed with %s", config.TOGGLEDBNAME, err)
		return
	}

	return
}

func getToggleSchedule(name string) (schedule DeviceToggle, error error) {
	db, err := database.Open(fmt.Sprintf("%s/%s", config.App.DbPath, config.TOGGLEDBNAME))
	if err != nil {
		return schedule, err
	}

	value, err := db.Get(name)
	if err != nil {
		return schedule, err
	}

	s, err := unmarshalToggleSchedule(value)
	if err != nil {
		log.Printf("WARN: getToggleSchedule failed to process %s", name)
		//ToDo: metric
		return schedule, err
	}

	return s, nil
}

func getAllToggleSchedules() (schedules []DeviceToggle, error error) {
	db, err := database.Open(fmt.Sprintf("%s/%s", config.App.DbPath, config.TOGGLEDBNAME))
	if err != nil {
		return schedules, err
	}

	keys, err := db.GetAllKeys()
	if err != nil {
		return schedules, err
	}

	for _, k := range keys {
		data, err := db.Get(k)
		if err != nil {
			log.Printf("WARN: getAllToggleSchedules failed to get %s", k)
			//ToDo: metric
			continue
		}

		s, err := unmarshalToggleSchedule(data)
		if err != nil {
			log.Printf("WARN: getAllToggleSchedules failed to process %s", k)
			//ToDo: metric
			continue
		}
		schedules = append(schedules, s)
	}

	return schedules, nil
}

//http handlers
func ToggleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := strings.ToLower(vars["friendlyName"])

	response, err := getToggleSchedule(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("ERROR: ToggleGet %s", err)
		//ToDo: Metric
		return
	}

	respondHttpBody(w, r, response)
	//ToDo: Metric
	return
}

func ToggleGetAll(w http.ResponseWriter, r *http.Request) {
	response, err := getAllToggleSchedules()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("ERROR: ToggleGetAll %s", err)
		//ToDo: Metric
		return
	}

	respondHttpBody(w, r, response)
	//ToDo: Metric
	return
}

//http handlers
func ToggleNew(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := strings.ToLower(vars["friendlyName"])

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//ToDo: Metric
		log.Printf("ERROR: ToggleNew %s", err)
		return
	}

	if err := newSchedule(name, data, config.TOGGLEDBNAME); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//ToDo: Metric
		log.Printf("ERROR: ToggleNew %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	//ToDo: Metric
	return
}

func ToggleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := strings.ToLower(vars["friendlyName"])

	if err := deleteSchedule(name, config.TOGGLEDBNAME); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//ToDo: Metric
		log.Printf("ERROR: ToggleDelete %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	//ToDo: Metric
	return
}
