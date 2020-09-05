package schedule

import (
	"github.com/rebelit/gome-schedule/common/config"
	"github.com/rebelit/gome-schedule/device"
	"log"
	"strconv"
	"sync"
)

//Public functions
func InitializeDatabases() {
	log.Printf("INFO: Initialize the databases")

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

//State Schedule processing
func processStateSchedules() {
	log.Println("INFO: state schedules processing start")

	schedules, err := getAllStateSchedules()
	if err != nil {
		log.Printf("ERROR: process getAllStateSchedules %s", err)
		//ToDo: metric
		return
	}

	//ToDo: metric on schedule length
	if len(schedules) <= 0 {
		log.Printf("WARN: no state schedules to process, should there be?")
		return
	}

	for _, s := range schedules {
		//ToDo: maybe combine these in to an okToProcess function.
		if !s.Enabled {
			//skip disabled
			continue
		}

		if !isToday(s.Day) {
			//skip if not today
			continue
		}

		enforceStart, _ := strconv.Atoi(s.StartTime)
		enforceEnd, _ := strconv.Atoi(s.EndTime)
		if !inScheduleBlock(enforceStart, enforceEnd) {
			continue
		}

		go s.doStateSchedule()
	}

	return
}

func (s *DeviceState) doStateSchedule() {
	desiredState := s.DeviceState

	//get the current action state
	currentState, err := device.GetDeviceActionState(s.DeviceType, s.DeviceName, s.DeviceAction)
	if err != nil {
		log.Printf("ERROR: schedule, powerState %s/%s, %s", s.DeviceType, s.DeviceName, err)
		//ToDo: metric
	}

	if currentState != s.DeviceState {
		if err := device.SetDeviceActionState(s.DeviceType, s.DeviceName, s.DeviceAction, boolToState(desiredState)); err != nil {
			log.Printf("ERROR: setting %s %s %s %s desired state", s.DeviceType, s.DeviceName,
				s.DeviceAction, boolToState(desiredState))
			//toDo: metric
		}
	}

	return
}

func inScheduleBlock(enforceStart int, enforceEnd int) (inScheduleBlock bool) {
	_, currentTime, _, _ := currentTimeSplit()
	reverseCheck := false
	inScheduleBlock = false

	if enforceEnd <= enforceStart {
		//spans a day PM to AM on schedule
		reverseCheck = true
	}

	if !reverseCheck {
		//does not span PM to AM
		inScheduleBlock = inBetween(currentTime, enforceStart, enforceEnd)
	} else {
		//spans a day PM to AM reverse check the schedule
		inScheduleBlock = inBetweenReverse(currentTime, enforceEnd, enforceStart)
	}

	return inScheduleBlock
}

//Toggle Schedule processing
func processToggleSchedules() {
	log.Println("INFO: toggle schedules processing start")

	schedules, err := getAllToggleSchedules()
	if err != nil {
		log.Printf("ERROR: process getAllToggleSchedules %s", err)
		//ToDo: metric
		return
	}

	//ToDo: metric on schedule length
	if len(schedules) <= 0 {
		log.Printf("WARN: no toggle schedules to process, should there be?")
		return
	}

	_, now, _, _ := currentTimeSplit()

	for _, s := range schedules {
		//ToDo: maybe combine these in to an okToProcess function.
		if !s.Enabled {
			//skip disabled
			continue
		}

		if !isToday(s.Day) {
			//skip if not today
			continue
		}

		toggleTime, _ := strconv.Atoi(s.ToggleTime)
		if now != toggleTime {
			//skip time does not match
			continue
		}

		go s.doToggleSchedule()
	}

	return
}

func (s *DeviceToggle) doToggleSchedule() {
	if err := device.SetDeviceActionState(s.DeviceType, s.DeviceName, s.DeviceAction, boolToState(s.DeviceState)); err != nil {
		log.Printf("ERROR: setting %s %s %s %s desired state", s.DeviceType, s.DeviceName,
			s.DeviceAction, boolToState(s.DeviceState))
		//toDo: metric
	}

	return
}
