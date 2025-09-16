package jira

import (
	"fmt"
	"os"
	"strings"

	"github.com/andygrunwald/go-jira"
)

const baseURL = "https://reload.atlassian.net"

// GetIssues gets issues.
func GetIssues(filterID string) ([]jira.Issue, error) {
	//nolint:exhaustivestruct
	tp := jira.BasicAuthTransport{
		Username: os.Getenv("TRIAGEBOT_JIRA_USER"),
		Password: os.Getenv("TRIAGEBOT_JIRA_PASS"),
	}

	jiraClient, err := jira.NewClient(tp.Client(), baseURL)
	if err != nil {
		return nil, fmt.Errorf("getting jira client: %w", err)
	}

	jql := fmt.Sprintf("filter = %s", os.Getenv(filterID))

	issues, _, err := jiraClient.Issue.SearchV2JQL(jql, nil)
	if err != nil {
		return nil, fmt.Errorf("getting jira issues: %w", err)
	}

	return issues, nil
}

// FormatIssues formats issues as Markdown.
func FormatIssues(issues []jira.Issue) string {
	output := make([]string, 0, len(issues))

	for _, issue := range issues {
		assignee := ""
		if issue.Fields.Assignee != nil {
			assignee = fmt.Sprintf(" - %s", issue.Fields.Assignee.DisplayName)
		}

		output = append(
			output,
			fmt.Sprintf("* [%s](%s/browse/%s) - %s%s", issue.Key, baseURL, issue.Key, issue.Fields.Summary, assignee),
		)
	}

	return strings.Join(output, "\n")
}
