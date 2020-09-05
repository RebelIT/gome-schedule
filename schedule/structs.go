package schedule

//DeviceState schedule ensures the device is in the state during
// the time periods defined in on | off.
// example: A TV should remain off from 10PM to 8AM.  if it is turned
// on between that time the DeviceState is enforced.
type DeviceState struct {
	Name         string `json:"name"`          //name of schedule
	DeviceName   string `json:"device_name"`   //name of device in gome-core database
	DeviceType   string `json:"device_type"`   //device name in gome-core database
	Enabled      bool   `json:"enabled"`       //schedule enabled or disabled
	Day          string `json:"day"`           //day of week to trigger
	DeviceAction string `json:"device_action"` //what action to do
	DeviceState  bool   `json:"device_state"`  //true|false device state on|off
	StartTime    string `json:"start_time"`    //time to turn on
	EndTime      string `json:"end_time"`      //time to turn off
}

//DeviceToggle schedule fires the action at the specified time
type DeviceToggle struct {
	Name         string `json:"name"`          //name of schedule
	DeviceName   string `json:"device_name"`   //name of device in gome-core database
	DeviceType   string `json:"device_type"`   //device name in gome-core database
	Enabled      bool   `json:"enabled"`       //schedule enabled or disabled
	Reoccurring  bool   `json:"reoccurring"`   //do this action every day and time defined
	Day          string `json:"day"`           //day of week to trigger
	DeviceAction string `json:"device_action"` //what action to do
	DeviceState  bool   `json:"device_state"`  //true|false device state on|off
	ToggleTime   string `json:"toggle_time"`   //time to toggle on
}
