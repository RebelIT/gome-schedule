package notify

import (
	"encoding/json"
	"fmt"
	"github.com/rebelit/gome-schedule/common/config"
	"github.com/rebelit/gome-schedule/common/httpRequest"
	"github.com/rebelit/gome-schedule/common/stat"
	"log"
)

type SlackMsg struct {
	Text     string `json:"text"`
	Username string `json:"username"`
	IconPath string `json:"icon_path"`
}

func Slack(message string) {
	if config.App.SlackWebhook == "" {
		log.Printf("DEBUG: would have sent slack: %s", message)
		return
	}
	content := SlackMsg{}
	content.Text = message
	content.Username = "gome"
	content.IconPath = ""
	body, _ := json.Marshal(content)
	resp, err := httpRequest.Post(config.App.SlackWebhook, body, nil)
	if err != nil {
		log.Printf("ERROR: Slack, %s", err)
		stat.Notify("slack", stat.STATEFAILURE, 0)
		return
	}
	if resp.StatusCode != 200 {
		log.Printf("ERROR: Slack, %s", fmt.Errorf("slack returned a non 200 response"))
		stat.Notify("slack", stat.STATEFAILURE, resp.StatusCode)
		return
	}

	log.Printf("INFO: Slack, sent %s", message)
	stat.Notify("slack", stat.STATEOK, resp.StatusCode)
	return
}
