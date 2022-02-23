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

	issues, _, err := jiraClient.Issue.Search(jql, nil)
	if err != nil {
		return nil, fmt.Errorf("getting jira issues: %w", err)
	}

	return issues, nil
}

// FormatIssues formats issues as Markdown.
func FormatIssues(issues []jira.Issue) string {
	output := make([]string, 0, len(issues))

	for _, issue := range issues {
		// We HTML and URL encode the dash in the issue key to
		// hide the issue from Hubot (otherwise Hubot will
		// spam with follow up comments).
		issueKeyHTML := strings.Replace(issue.Key, "-", "&#x2D;", 1)
		issueKeyURL := strings.Replace(issue.Key, "-", "%2d", 1)

		assignee := ""
		if issue.Fields.Assignee != nil {
			assignee = fmt.Sprintf(" - %s", issue.Fields.Assignee.DisplayName)
		}

		output = append(
			output,
			fmt.Sprintf("* [%s](%s/browse/%s) - %s%s", issueKeyHTML, baseURL, issueKeyURL, issue.Fields.Summary, assignee),
		)
	}

	return strings.Join(output, "\n")
}
