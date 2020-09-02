package schedule

import (
	"encoding/json"
	"fmt"
	"github.com/rebelit/gome-schedule/common/config"
	"github.com/rebelit/gome-schedule/database"
	"log"
	"net/http"
	"os"
)

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

func requiredPathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func createRequiredPath(path string) error {
	err := os.Mkdir(path, 0755)
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
