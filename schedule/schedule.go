package schedule

import (
	"github.com/rebelit/gome-schedule/common/config"
	"log"
	"sync"
)

func InitializeDatabases() {
	dbPath := config.App.DbPath
	if !requiredPathExists(dbPath) {
		err := createRequiredPath(dbPath)
		if err != nil {
			log.Fatalf("FATAL: Unable to create required database directory %s", dbPath)
		}
	}

	var waitgroup sync.WaitGroup
	waitgroup.Add(2)

	go stateDb(&waitgroup)
	go toggleDb(&waitgroup)

	waitgroup.Wait()
	return
}
