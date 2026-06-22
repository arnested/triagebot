package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"arnested.dk/go/triagebot/internal/cal"
	"arnested.dk/go/triagebot/internal/jira"
	"arnested.dk/go/triagebot/internal/team"
)

// ZulipPayload is the outgoing message received from Zulip.
//
//nolint:tagliatelle
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

	response := response(r.Context())
	//nolint:errchkjson
	_ = json.NewEncoder(w).Encode(response)
}

func response(ctx context.Context) ZulipResponse {
	response := ZulipResponse{}

	issues, err := jira.GetIssues("TRIAGEBOT_JIRA_FILTER")
	if err != nil {
		response.Content = "error getting issues: " + err.Error()

		return response
	}

	unreleasedIssues, err := jira.GetIssues("TRIAGEBOT_JIRA_FILTER_UNRELEASED")
	if err != nil {
		response.Content = fmt.Sprintf("error getting unreleased issues: %s\n", err.Error())

		return response
	}

	needsAction := false

	if len(issues) == 0 {
		response.Content = fmt.Sprintln(NoIssuesNeedTriage)
	}

	if len(issues) > 0 {
		response.Content = fmt.Sprintf(LeadText+":\n\n%s", len(issues), jira.FormatIssues(issues))
		needsAction = true
	}

	if len(unreleasedIssues) > 0 {
		response.Content = fmt.Sprintf(
			"%s\n\n\n"+UnreleasedText+":\n\n%s",
			response.Content,
			len(unreleasedIssues),
			jira.FormatIssues(unreleasedIssues),
		)
		needsAction = true
	}

	// Only tag people if they need to do something - and if it's a work day.
	if needsAction && cal.IsWorkday(time.Now()) {
		response.Content += "\n\n" + team.Triage(ctx, true)
	}

	return response
}
