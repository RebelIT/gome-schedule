package device

import (
	"encoding/json"
	"fmt"
	"github.com/rebelit/gome-schedule/common/config"
	"github.com/rebelit/gome-schedule/common/httpRequest"
	"io/ioutil"
	"log"
	"net/http"
)

func setHeaders() (headers map[string]string) {
	headers = make(map[string]string)
	headers["Authorization"] = "Bearer " + config.App.CoreServiceToken
	headers["Content-Type"] = "application/json"

	return
}

func parseStateResponse(resp *http.Response) (state bool, error error) {
	s := DevPower{}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(data, &s); err != nil {
		return false, err
	}

	return s.State, nil
}

func GetDeviceActionState(devType string, devName string, devAction string) (state bool, error error) {
	baseUrl := fmt.Sprintf("%s:%s/api", config.App.CoreServiceUrl, config.App.CoreServicePort)
	uriPart := fmt.Sprintf("/%s/%s/%s", devType, devName, devAction)

	resp, err := httpRequest.Get(baseUrl+uriPart, setHeaders())
	if err != nil {
		return false, err
	}

	devState, err := parseStateResponse(&resp)
	if err != nil {
		return false, err
	}

	return devState, nil
}

func SetDeviceActionState(devType string, devName string, devAction string, devState string) error {
	baseUrl := fmt.Sprintf("%s:%s/api", config.App.CoreServiceUrl, config.App.CoreServicePort)
	uriPart := fmt.Sprintf("/%s/%s/%s/%s", devType, devName, devAction, devState)

	log.Printf("DEBUG:: %s", baseUrl+uriPart)
	resp, err := httpRequest.Put(baseUrl+uriPart, nil, setHeaders())
	if err != nil {
		log.Printf("DEBUG:: %s", err)
		return err
	}

	if resp.StatusCode != 200 {
		log.Printf("DEBUG:: non-200")
		return fmt.Errorf("non-200 response from %s", baseUrl+uriPart)
	}
	log.Printf("DEBUG:: I passed the http put...")
	return nil
}
