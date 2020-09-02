package schedule

import "time"

//DeviceState schedule ensures the device is in the state during
// the time periods defined in on | off.
// example: A TV should remain off from 10PM to 8AM.  if it is turned
// on between that time the DeviceState is enforced.
type DeviceState struct {
	Name       string    `json:"name"`        //name of schedule
	DeviceName string    `json:"device_name"` //name of device in gome-core database
	DeviceType string    `json:"device_type"` //device name in gome-core database
	Status     string    `json:"status"`      //schedule enabled or disabled
	Day        string    `json:"day"`         //day of week to trigger
	On         time.Time `json:"on"`          //time to turn on
	OnAction   string    `json:"on_action"`   //what action to perform (see gome-core uri's | /api/roku/{name}/power/on)
	Off        time.Time `json:"off"`         //time to turn off
	OffAction  string    `json:"off_action"`  //what action to perform (see gome-core uri's | /api/roku/{name}/power/off)
}

//DeviceToggle schedule fires the action at the specified time
type DeviceToggle struct {
	Name         string      `json:"name"`          //name of schedule
	DeviceName   string      `json:"device_name"`   //name of device in gome-core database
	DeviceType   string      `json:"device_type"`   //device name in gome-core database
	Status       string      `json:"status"`        //schedule enabled or disabled
	Reoccurring  bool        `json:"reoccurring"`   //do this action every day and time defined
	Day          []string    `json:"day"`           //day of week to trigger
	ToggleTime   []time.Time `json:"toggle_time"`   //time to toggle on
	ToggleAction string      `json:"toggle_action"` //what action to perform (see gome-core uri's | /api/roku/{name}/power/on)
}
