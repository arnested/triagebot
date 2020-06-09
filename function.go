package triagebot

import (
	"encoding/json"
	"fmt"
	"net/http"

	"arnested.dk/go/triagebot/internal/jira"
)

// Payload is the outgoing message received from Zulip.
type Payload struct {
	Data    string `json:"data"`
	Token   string `json:"token"`
	Message struct {
		ID          int    `json:"id"`
		SenderEmail string `json:"sender_email"`
	} `json:"message"`
}

// Response is the data we send as an answer to the payload.
type Response struct {
	Content string `json:"content"`
}

// Handle is the entrypoint for the Google Cloud Function.
func Handle(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	// Handle outgoing messages from Zulip.
	case "/outgoing":
		outgoingHandler := http.HandlerFunc(outgoing)
		chain := parseMiddleware(authenticationMiddleware(authorizationMiddleware(reactMiddleware(outgoingHandler))))
		chain.ServeHTTP(w, r)

	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func outgoing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := response()
	json.NewEncoder(w).Encode(response)
}

func response() Response {
	response := Response{}

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
