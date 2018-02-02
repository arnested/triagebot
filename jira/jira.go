package jira

import (
	"fmt"
	"os"
	"strings"

	"github.com/andygrunwald/go-jira"
)

var baseURL = "https://reload.atlassian.net"

// GetIssues gets issues.
func GetIssues() []jira.Issue {
	jiraClient, err := jira.NewClient(nil, baseURL)
	if err != nil {
		panic(err)
	}

	res, err := jiraClient.Authentication.AcquireSessionCookie(os.Getenv("TRIAGEBOT_JIRA_USER"), os.Getenv("TRIAGEBOT_JIRA_PASS"))
	if err != nil || res == false {
		fmt.Printf("Result: %v\n", res)
		panic(err)
	}

	jql := fmt.Sprintf("filter = %s", os.Getenv("TRIAGEBOT_JIRA_FILTER"))
	issues, _, err := jiraClient.Issue.Search(jql, nil)
	if err != nil {
		panic(err)
	}

	return issues
}

// FormatIssues formats issues.
func FormatIssues(issues []jira.Issue) string {
	var output []string
	for _, issue := range issues {
		output = append(output, fmt.Sprintf("* [%s: %s](%s/browse/%s)\n", issue.Key, issue.Fields.Summary, baseURL, issue.Key))
	}

	return strings.Join(output, "")
}
