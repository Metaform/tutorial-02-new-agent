package activity

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"tutorial-01-new-agent/model"

	"github.com/eclipse-cfm/cfm/common/system"
	"github.com/eclipse-cfm/cfm/pmanager/api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type UserInfoProcessor struct {
	api.BaseActivityProcessor
	Monitor system.LogMonitor
	// those are not exported:
	httpClient http.Client
	url        string
}

func NewProcessor(monitor system.LogMonitor, client http.Client, url string) UserInfoProcessor {
	return UserInfoProcessor{
		BaseActivityProcessor: api.BaseActivityProcessor{},
		Monitor:               monitor,
		httpClient:            client,
		url:                   url,
	}
}

func (h UserInfoProcessor) ProcessDeploy(activityContext api.ActivityContext) api.ActivityResult {

	tracer := otel.GetTracerProvider().Tracer("cfm.agent.user-info")
	_, span := tracer.Start(activityContext.Context(), "fetch-user-info")
	defer span.End()

	resp, err := h.httpClient.Get(h.url)
	if err != nil {
		span.RecordError(err)
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	span.AddEvent("user data fetched")

	if err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: err}
	}
	// this REST call returns a list of users, we just pick one at random, and set them on the output
	var result []model.User
	if err := json.Unmarshal(body, &result); err != nil {
		span.RecordError(err)
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: err}
	}

	// select one user at random
	if len(result) > 0 {
		randomIndex := rand.Intn(len(result))
		selectedUser := result[randomIndex]
		span.SetAttributes(attribute.Int("userIndex", randomIndex), attribute.String("userName", selectedUser.Name))
		activityContext.SetOutputValue("selectedUser", selectedUser)
		h.Monitor.Infof("Fetched User Info for: %s (%s)", selectedUser.Name, selectedUser.Username)
	}

	return api.ActivityResult{Result: api.ActivityResultComplete}
}

func (h UserInfoProcessor) ProcessDispose(activityContext api.ActivityContext) api.ActivityResult {
	//TODO implement me
	h.Monitor.Infof("Disposing hello world: can't dispose of the entire world...")
	return api.ActivityResult{Result: api.ActivityResultComplete}
}
