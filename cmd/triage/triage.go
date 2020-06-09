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

	issues, err := jira.GetIssues()

	if err != nil {
		fmt.Printf("error getting issues: %s\n", err.Error())
		os.Exit(1)
	}

	if len(issues) == 0 {
		fmt.Println(triagebot.NoIssuesNeedTriage)
		os.Exit(0)
	}

	result := c.Compile(jira.FormatIssues(issues))

	fmt.Printf("%s:\n\n%s", triagebot.LeadText, result)
}
