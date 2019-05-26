package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/andygrunwald/go-jira"
)

var baseURL = "https://reload.atlassian.net"

// getIssues gets issues.
func getIssues() []jira.Issue {
	tp := jira.BasicAuthTransport{
		Username: os.Getenv("TRIAGEBOT_JIRA_USER"),
		Password: os.Getenv("TRIAGEBOT_JIRA_PASS"),
	}

	jiraClient, err := jira.NewClient(tp.Client(), baseURL)
	if err != nil {
		panic(err)
	}

	jql := fmt.Sprintf("filter = %s", os.Getenv("TRIAGEBOT_JIRA_FILTER"))
	issues, _, err := jiraClient.Issue.Search(jql, nil)
	if err != nil {
		panic(err)
	}

	return issues
}

// formatIssues formats issues.
func formatIssues(issues []jira.Issue) string {
	var output []string
	for _, issue := range issues {
		// We lowercase the issue key to hide the issue from
		// Hubot (otherwise Hubot will spam with follow up
		// comments). Hubot is configured with
		// HUBOT_JIRA_IGNORECASE=false.
		output = append(output, fmt.Sprintf("* [%s: %s](%s/browse/%s)\n", strings.ToLower(issue.Key), issue.Fields.Summary, baseURL, strings.ToLower(issue.Key)))
	}

	return strings.Join(output, "")
}
