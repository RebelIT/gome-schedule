# gome-schedule
GoLang Home Automation - Job Scheduler

## An add-on to gome-core
* **[gome-core](https://github.com/RebelIT/gome-core)** is required for this to function

1. Set a desired state schedule and time to check, this add-on will enforce that state schedule at the start & end scheduled times
    * ex: ensure a roku TV remains off between 10PM and 9AM.  if the device is turned on it will be turned off.
2. Set a fire & forget (crontab) like action to a device (single schedule or recurring)
    * ex: turn on a roku tv at 8AM every weekday
3. In development: get and launch a roku app at a specific time every day after powering on (deep linking where applicable to launch content inside the app/channel as well)
    
## WIP
* only supports schedules for device power on and off on a schedule, see [gome-core doco](https://github.com/RebelIT/gome-core/blob/master/README.md) for supported devices.
* currently no method to update a schedule, only delete and post new with the updated body

## Application Configuration
### Defaults
* name: `"gome-schedule"`
* statsd: `""`
* dbPath: `"badgerDatabase"`
* authToken: `"changeMePlease"`
* port: `"6661"`
* coreToken: `"changeMePlease"`
* coreUrl: `"http://localhost`
* corePort: `"6660"`
* stateTimeSec `30`

### Flags
* Override any of the defaults with args
    * `-name "myApp"  -statsd "127.0.0.1:8125" -dbPath "./myDatabases" -authToken "ThisIsMyNewP@ssw0Rd" -port "8080"`

## Installation
### Docker
* **NOTE** if you change the default port you need to change it in the command below and pass in the `-port` flag
1. `docker build -t gome-schedule:latest .`
2. `docker run -it --rm -p 6661:6661 -v $PWD:/gome-schedule gome-v` - Run with defaults
1. `docker run -it --rm -p 6661:6661 -v $PWD:/gome-schedule gome-schedule -authToken "abc"` - Override defaults
1. `docker run -it --rm -p 8081:8081 -v $PWD:/gome-schedule gome-schedule -port "8081"` - Run on different port

### Manual
1. build it yourself with `go build` and any args you need 
2. execute it

## Supported schedulable actions
* [x] roku tv power on and off
* [] roku launch channel on schedule
* [] rpIoT actions scheduled GPIO controls
* [] rpIoT maintenance actions system updates, reboots
* [] tuya device actions lights on/off on schedule
* [] ecoVacs robot vacuum scheduled start/stop
* thats all I have in my house RN...  more developed as I purchase smart things. 

## Usage: 
* `/api/status`
    * **Method**: GET
        * **Purpose**: Checks the status of the web api
        * **Returns**: Status code

* `/api/schedule/state`
    * **Method**: GET
        * **Purpose**: Gets what state schedule (loaded in the database)
        * **Returns**: Array of schedule names

* `/api/schedule/state/{friendlyName}`
    * **Method**: GET
        * **Purpose**: Gets the state schedule details of schedule name
        * **Returns**: State schedule type json
     * **Method**: POST
         * **Purpose**: Adds a new state schedule
         * **Returns**: Status code
     * **Method**: PUT **NOT YET IMPLEMENTED**
         * **Purpose**: Updates an existing state schedule
         * **Returns**: Status code
     * **Method**: DELETE
         * **Purpose**: Removes a state schedule
         * **Returns**: Status code
         
* `/api/schedule/toggle`
    * **Method**: GET
        * **Purpose**: Gets what toggle schedule (loaded in the database)
        * **Returns**: Array of schedule names
        
* `/api/schedule/toggle/{friendlyName}`
    * **Method**: GET
        * **Purpose**: Gets the toggle schedule details of schedule name
        * **Returns**: Toggle schedule type json
     * **Method**: POST
        * **Purpose**: Adds a new toggle schedule
        * **Returns**: Status code
     * **Method**: PUT **NOT YET IMPLEMENTED**
        * **Purpose**: Updates an existing toggle schedule
        * **Returns**: Status code         
     * **Method**: DELETE
        * **Purpose**: Removes a toggle schedule
        * **Returns**: Status code
        
## QuickStart Playground
#### gome-core
1. `git clone git@github.com:RebelIT/gome-core.git`
2. `docker build -t gome-core:latest .`
1. `docker run -it --rm -p 6660:6660 -v $PWD:/gome-core gome-core`
1. `curl -i http://localhost:6660/api/status`
1. `curl -i -H "Content-Type: application/json" -H "Authorization: Bearer changeMePlease" -X POST -d '{"name": "basement","address": "192.168.1.10","port": "8060"}' http://localhost:6660/api/device/roku`
1. `curl -i -H "Content-Type: application/json" -H "Authorization: Bearer changeMePlease" http://localhost:6660/api/device/roku`
1. `curl -i -H "Content-Type: application/json" -H "Authorization: Bearer changeMePlease" -X POST http://localhost:6660/api/roku/basement/key/home`

#### gome-schedule
from another directory & terminal
1. `git clone git@github.com:RebelIT/gome-schedule.git`
2. `docker build -t gome-schedule:latest .`
1. `docker run -it --rm -p 6661:6661 -v $PWD:/gome-schedule gome-schedule`
1. `curl -i http://localhost:6661/api/status`
1. `curl -i -H "Content-Type: application/json" -H "Authorization: Bearer changeMePlease" -X POST -d '{"name":"family tv remains off overnight","device_name":"basement","device_type":"roku","enabled":true,"day":"sunday","device_action":"power","device_state":false,"start_time":"1900","end_time":"0800"}' http://localhost:6661/api/schedule/state/familyTvOff`
1. `curl -i -H "Content-Type: application/json" -H "Authorization: Bearer changeMePlease" http://localhost:6661/api/schedule/state`
1. `curl -i -H "Content-Type: application/json" -H "Authorization: Bearer changeMePlease" http://localhost:6661/api/schedule/state/familyTvOff`
1. `curl -i -H "Content-Type: application/json" -H "Authorization: Bearer changeMePlease" -X POST -d '{"name":"family turn on","device_name":"basement","device_type":"roku","enabled":true,"day":"monday","Reoccurring":false,"device_action":"power","device_state":true,"toggle_time":"0900"}' http://localhost:6661/api/schedule/toggle/familyTvOn`
1. `curl -i -H "Content-Type: application/json" -H "Authorization: Bearer changeMePlease" http://localhost:6661/api/schedule/toggle`
1. `curl -i -H "Content-Type: application/json" -H "Authorization: Bearer changeMePlease" http://localhost:6661/api/schedule/toggle/familyTvOn`