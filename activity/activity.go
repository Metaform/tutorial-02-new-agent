package activity

import (
	"encoding/json"
	"hello-world-agent/model"
	"io"
	"math/rand"
	"net/http"

	"github.com/eclipse-cfm/cfm/common/system"
	"github.com/eclipse-cfm/cfm/pmanager/api"
)

type HelloWorldProcessor struct {
	api.BaseActivityProcessor
	Monitor system.LogMonitor
	// those are not exported:
	httpClient http.Client
	url        string
}

func NewProcessor(monitor system.LogMonitor, client http.Client, url string) HelloWorldProcessor {
	return HelloWorldProcessor{
		BaseActivityProcessor: api.BaseActivityProcessor{},
		Monitor:               monitor,
		httpClient:            client,
		url:                   url,
	}
}

func (h HelloWorldProcessor) ProcessDeploy(activityContext api.ActivityContext) api.ActivityResult {

	resp, err := h.httpClient.Get(h.url)
	if err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: err}
	}
	// this REST call returns a list of users, we just pick one at random, and set them on the output
	var result []model.User
	if err := json.Unmarshal(body, &result); err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: err}
	}

	// select one user at random
	if len(result) > 0 {
		randomIndex := rand.Intn(len(result))
		selectedUser := result[randomIndex]
		activityContext.SetOutputValue("selectedUser", selectedUser)
		h.Monitor.Infof("Fetched User Info for: %s (%s)", selectedUser.Name, selectedUser.Username)
	}

	return api.ActivityResult{Result: api.ActivityResultComplete}
}

func (h HelloWorldProcessor) ProcessDispose(activityContext api.ActivityContext) api.ActivityResult {
	//TODO implement me
	h.Monitor.Infof("Disposing hello world: can't dispose of the entire world...")
	return api.ActivityResult{Result: api.ActivityResultComplete}
}
