package triagebot

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"arnested.dk/go/triagebot/internal/cal"
	"arnested.dk/go/triagebot/internal/jira"
	"github.com/containrrr/shoutrrr"
)

const tag = "@**all**"

// SchedulePayload is the message posted from Google Cloud Scheduler.
type SchedulePayload struct {
	Stream string `json:"stream"`
	Token  string `json:"token"`
}

func schedule(w http.ResponseWriter, r *http.Request) {
	payload := r.Context().Value(payloadKey{}).(SchedulePayload)
	t := time.Now()

	query := &url.Values{}
	query.Set("stream", payload.Stream)
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
		http.Error(w, err.Error(), http.StatusOK)

		return
	}

	issues, err := jira.GetIssues()

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)

		return
	}

	if len(issues) == 0 && !cal.IsFirstWorkdaySinceDrupalSecurityAnnouncements(time.Now()) {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)

		return
	}

	message := fmt.Sprintf("%s, %s:\n\n%s", LeadText, tag, jira.FormatIssues(issues))
	if len(issues) == 0 {
		message = fmt.Sprintln(NoIssuesNeedTriage)
	}

	errs := sender.Send(message, nil)

	if errs != nil && len(errs) > 0 && errs[0] != nil {
		http.Error(w, fmt.Sprintf("%v", errs), http.StatusServiceUnavailable)

		return
	}

	w.WriteHeader(http.StatusCreated)
}