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

func unmarshalStateSchedule(data []byte) (schedule DeviceState, error error) {
	if err := json.Unmarshal(data, &schedule); err != nil {
		return schedule, err
	}

	return
}

func stateDb(wg *sync.WaitGroup) {
	defer wg.Done()

	if err := newDatabase(config.STATEDBNAME); err != nil {
		log.Fatalf("FATAL: %s failed with %s", config.STATEDBNAME, err)
		return
	}

	return
}

func getStateSchedule(name string) (schedule DeviceState, error error) {
	db, err := database.Open(fmt.Sprintf("%s/%s", config.App.DbPath, config.STATEDBNAME))
	if err != nil {
		return schedule, err
	}

	value, err := db.Get(name)
	if err != nil {
		return schedule, err
	}

	s, err := unmarshalStateSchedule(value)
	if err != nil {
		log.Printf("WARN: getAllStateSchedules failed to process %s", name)
		//ToDo: metric
		return schedule, err
	}

	return s, nil
}

func getAllStateSchedules() (schedules map[string]DeviceState, error error) {
	db, err := database.Open(fmt.Sprintf("%s/%s", config.App.DbPath, config.STATEDBNAME))
	if err != nil {
		return schedules, err
	}

	keys, err := db.GetAllKeys()
	if err != nil {
		return schedules, err
	}

	schedules = make(map[string]DeviceState)
	//ToDo: need to pool the connections in the database. this is kinda hacky...
	for _, k := range keys {
		db, err := database.Open(fmt.Sprintf("%s/%s", config.App.DbPath, config.STATEDBNAME))
		if err != nil {
			return schedules, err
		}

		data, err := db.Get(k)
		if err != nil {
			log.Printf("WARN: getAllStateSchedules failed to get %s", k)
			//ToDo: metric
			continue
		}

		s, err := unmarshalStateSchedule(data)
		if err != nil {
			log.Printf("WARN: getAllStateSchedules failed to process %s", k)
			//ToDo: metric
			continue
		}
		schedules[k] = s
	}

	return schedules, nil
}

//http handlers
func StateGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := strings.ToLower(vars["friendlyName"])

	response, err := getStateSchedule(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("ERROR: StateGet %s", err)
		//ToDo: Metric
		return
	}

	respondHttpBody(w, r, response)
	//ToDo: Metric
	return
}

func StateGetAll(w http.ResponseWriter, r *http.Request) {
	response, err := getAllStateSchedules()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("ERROR: StateGetAll %s", err)
		//ToDo: Metric
		return
	}

	respondHttpBody(w, r, response)
	//ToDo: Metric
	return
}

func StateNew(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := strings.ToLower(vars["friendlyName"])

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//ToDo: Metric
		log.Printf("ERROR: StateNew %s", err)
		return
	}

	if err := newSchedule(name, data, config.STATEDBNAME); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//ToDo: Metric
		log.Printf("ERROR: StateNew %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	//ToDo: Metric
	return
}

func StateUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := strings.ToLower(vars["friendlyName"])

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//ToDo: Metric
		log.Printf("ERROR: StateUpdate %s", err)
		return
	}

	if err := updateSchedule(name, data, config.STATEDBNAME); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//ToDo: Metric
		log.Printf("ERROR: StateUpdate %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	//ToDo: Metric
	return
}

func StateDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := strings.ToLower(vars["friendlyName"])

	if err := deleteSchedule(name, config.STATEDBNAME); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//ToDo: Metric
		log.Printf("ERROR: StateDelete %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	//ToDo: Metric
	return
}
