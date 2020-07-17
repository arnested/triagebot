package triagebot

import (
	"encoding/json"
	"fmt"
	"net/http"

	"arnested.dk/go/triagebot/internal/jira"
)

// ZulipPayload is the outgoing message received from Zulip.
type ZulipPayload struct {
	Data    string `json:"data"`
	Token   string `json:"token"`
	Message struct {
		ID          int    `json:"id"`
		SenderEmail string `json:"sender_email"`
	} `json:"message"`
}

// ZulipResponse is the data we send as an answer to the payload.
type ZulipResponse struct {
	Content string `json:"content"`
}

func outgoing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	response := response()
	_ = json.NewEncoder(w).Encode(response)
}

func response() ZulipResponse {
	response := ZulipResponse{}

	issues, err := jira.GetIssues()
	if err != nil {
		response.Content = fmt.Sprintf("error getting issues: %s", err.Error())

		return response
	}

	if len(issues) == 0 {
		response.Content = fmt.Sprintln(NoIssuesNeedTriage)

		return response
	}

	response.Content = fmt.Sprintf("%s:\n\n%s", LeadText, jira.FormatIssues(issues))

	return response
}
