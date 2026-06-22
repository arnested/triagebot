package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"arnested.dk/go/triagebot/internal/jira"
	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/router"
)

// SchedulePayload is the message posted from Google Cloud Scheduler.
type SchedulePayload struct {
	Stream string `json:"stream"`
	Token  string `json:"token"`
}

func shoutrrrSender(stream string) (*router.ServiceRouter, error) {
	t := time.Now()

	query := &url.Values{}
	query.Set("stream", stream)
	query.Set("topic", fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day()))

	zulipShoutrrrServiceURL := url.URL{
		User:     url.UserPassword(os.Getenv("ZULIP_BOT_MAIL"), os.Getenv("ZULIP_BOT_APIKEY")),
		Host:     "reload.zulipchat.com",
		Path:     "/",
		Scheme:   "zulip",
		RawQuery: query.Encode(),
	}

	sender, err := shoutrrr.CreateSender(zulipShoutrrrServiceURL.String())
	if err != nil {
		return nil, fmt.Errorf("creating shoutrrr sender: %w", err)
	}

	return sender, nil
}

func message() (string, bool, error) {
	issues, err := jira.GetIssues("TRIAGEBOT_JIRA_FILTER")
	if err != nil {
		return "", false, fmt.Errorf("jira filter: %w", err)
	}

	unreleasedIssues, err := jira.GetIssues("TRIAGEBOT_JIRA_FILTER_UNRELEASED")
	if err != nil {
		return "", false, fmt.Errorf("jira unreleased filter: %w", err)
	}

	message := fmt.Sprintln(NoIssuesNeedTriage)

	if len(issues) > 0 {
		message = fmt.Sprintf(LeadText+":\n\n%s", len(issues), jira.FormatIssues(issues))
	}

	if len(unreleasedIssues) > 0 {
		message = fmt.Sprintf(
			"%s\n\n\n"+UnreleasedText+":\n\n%s",
			message,
			len(unreleasedIssues),
			jira.FormatIssues(unreleasedIssues),
		)
	}

	return message, true, nil
}

func schedule(resp http.ResponseWriter, req *http.Request) {
	payloadData := req.Context().Value(payloadKey{})

	payload, ok := payloadData.(SchedulePayload)
	if !ok {
		http.Error(resp, http.StatusText(http.StatusNoContent), http.StatusNoContent)

		return
	}

	sender, err := shoutrrrSender(payload.Stream)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusOK)

		return
	}

	message, ok, err := message()
	if err != nil {
		http.Error(resp, err.Error(), http.StatusServiceUnavailable)

		return
	}

	if !ok {
		http.Error(resp, http.StatusText(http.StatusNoContent), http.StatusNoContent)

		return
	}

	errs := sender.Send(message, nil)

	if len(errs) > 0 && errs[0] != nil {
		http.Error(resp, fmt.Sprintf("%v", errs), http.StatusServiceUnavailable)

		return
	}

	resp.WriteHeader(http.StatusCreated)
}
