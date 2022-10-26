package main

import (
	"fmt"
	"os"

	"arnested.dk/go/triagebot"
	"arnested.dk/go/triagebot/internal/jira"
)

func main() {
	issues, err := jira.GetIssues("TRIAGEBOT_JIRA_FILTER")
	if err != nil {
		//nolint:forbidigo
		fmt.Printf("error getting issues: %s\n", err.Error())
		os.Exit(1)
	}

	unreleasedIssues, err := jira.GetIssues("TRIAGEBOT_JIRA_FILTER_UNRELEASED")
	if err != nil {
		//nolint:forbidigo
		fmt.Printf("error getting unreleased issues: %s\n", err.Error())
		os.Exit(1)
	}

	if len(issues) == 0 {
		//nolint:forbidigo
		fmt.Println(triagebot.NoIssuesNeedTriage)
	}

	if len(issues) > 0 {
		result := jira.FormatIssues(issues)
		//nolint:forbidigo
		fmt.Printf("%s:\n\n%s", triagebot.LeadText, result)
	}

	if len(unreleasedIssues) > 0 {
		unreleasedResult := jira.FormatIssues(unreleasedIssues)
		//nolint:forbidigo
		fmt.Printf("\n%s:\n\n%s", triagebot.UnreleasedText, unreleasedResult)
	}
}
