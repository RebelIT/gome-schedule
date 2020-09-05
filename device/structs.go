package device

type DevPower struct {
	Name  string `json:"name"`
	State bool   `json:"power_state"`
}
