package main

import (
	"fmt"
	"os"

	"arnested.dk/go/triagebot"
	"arnested.dk/go/triagebot/internal/jira"
	"github.com/tj/go-termd"
)

func main() {
	var c termd.Compiler

	issues, err := jira.GetIssues("TRIAGEBOT_JIRA_FILTER")
	if err != nil {
		fmt.Printf("error getting issues: %s\n", err.Error())
		os.Exit(1)
	}

	unreleasedIssues, err := jira.GetIssues("TRIAGEBOT_JIRA_FILTER_UNRELEASED")
	if err != nil {
		fmt.Printf("error getting unreleased issues: %s\n", err.Error())
		os.Exit(1)
	}

	if len(issues) == 0 {
		fmt.Println(triagebot.NoIssuesNeedTriage)
	}

	if len(issues) > 0 {
		result := c.Compile(jira.FormatIssues(issues))
		fmt.Printf("%s:\n\n%s", triagebot.LeadText, result)
	}

	if len(unreleasedIssues) > 0 {
		unreleasedResult := c.Compile(jira.FormatIssues(unreleasedIssues))
		fmt.Printf("\n%s:\n\n%s", triagebot.UnreleasedText, unreleasedResult)
	}
}
