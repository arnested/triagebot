package main

import (
	"net/http"
	"os"
	"strings"
	"time"
	_ "time/tzdata"

	_ "golang.org/x/crypto/x509roots/fallback"
)

func main() {
	missing := []string{}

	envs := []string{
		"TRIAGEBOT_ADDR",
		"TRIAGEBOT_JIRA_USER",
		"TRIAGEBOT_JIRA_PASS",
		"TRIAGEBOT_JIRA_FILTER",
		"TRIAGEBOT_JIRA_FILTER_UNRELEASED",
		"ZULIP_TOKEN",
		"ZULIP_BOT_MAIL",
		"ZULIP_BOT_APIKEY",
	}

	for _, env := range envs {
		_, ok := os.LookupEnv(env)

		if !ok {
			missing = append(missing, env)
		}
	}

	if len(missing) > 0 {
		panic("Missing environment variables: " + strings.Join(missing, ", "))
	}

	server := &http.Server{
		Addr:              os.Getenv("TRIAGEBOT_ADDR"),
		ReadHeaderTimeout: 3 * time.Second, //nolint:mnd
	}

	http.HandleFunc("/", Handle)

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// Handle is the entrypoint for the Google Cloud Function.
func Handle(resp http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	// Handle outgoing messages from Zulip.
	case "/outgoing":
		outgoingHandler := http.HandlerFunc(outgoing)
		chain := parseZulipMiddleware(
			authenticationOutgoingMiddleware(
				authorizationMiddleware(
					reactMiddleware(outgoingHandler))))
		chain.ServeHTTP(resp, req)

	// Handle schedules events from Google Cloud Scheduler.
	case "/schedule":
		scheduleHandler := http.HandlerFunc(schedule)
		chain := parseScheduleMiddleware(authenticationScheduleMiddleware(scheduleHandler))
		chain.ServeHTTP(resp, req)

	default:
		http.Error(resp, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}
