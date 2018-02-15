package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/arnested/triagebot/jira"
	"github.com/eko/flowbot"
	"github.com/jasonlvhit/gocron"
	"github.com/kardianos/service"
	"github.com/rickar/cal"
)

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	flowbot.FlowdockStreamUrl = fmt.Sprintf("https://stream.flowdock.com/flows/%s", os.Getenv("TRIAGEBOT_FLOW"))
	flowbot.FlowdockFlowToken = os.Getenv("TRIAGEBOT_FLOW_TOKEN")

	flowbot.FlowdockAuthUsername = os.Getenv("TRIAGEBOT_FLOW_USER")
	flowbot.FlowdockAuthPassword = os.Getenv("TRIAGEBOT_FLOW_PASS")
	flowbot.FlowdockRobotName = "TriageBot"

	go triagebot()

	location, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		log.Println("Unfortunately can't load a location")
		log.Println(err)
	} else {
		gocron.ChangeLoc(location)
	}

	gocron.Every(1).Day().At("08:40").Do(dailyStatus)
	gocron.Every(1).Day().At("12:20").Do(dailyStatus)

	<-gocron.Start()

	select {}
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "triagebot",
		DisplayName: "TriageBot",
		Description: "This is a robot monitoring triage needs.",
		UserName:    "nobody",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			err := s.Install()
			if err != nil {
				log.Fatal(err)
			}
		case "uninstall":
			err := s.Uninstall()
			if err != nil {
				log.Fatal(err)
			}
		case "start":
			err := s.Start()
			if err != nil {
				log.Fatal(err)
			}
		case "restart":
			err := s.Restart()
			if err != nil {
				log.Fatal(err)
			}
		case "stop":
			err := s.Stop()
			if err != nil {
				log.Fatal(err)
			}
		}

		os.Exit(0)
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

func dailyStatus() {
	flowbot.SendChat(getStatus(false))
}

func triagebot() {

	flowbot.AddCommand("^[Tt]riage\\?$", func(command flowbot.Command, entry flowbot.Entry) {
		flowbot.SendThreadChat(entry.ThreadId, getStatus(true))
	})

	flowbot.Stream()
}

func getStatus(force bool) string {
	status := ""

	issues := jira.GetIssues()
	c := cal.NewCalendar()

	cal.AddDanishHolidays(c)
	c.AddHoliday(
		cal.DKJuleaften,
		cal.DKNytaarsaften,
	)

	now := time.Now()

	if len(issues) > 0 && (force || c.IsWorkday(now)) {
		status = fmt.Sprintf("Følgende issues mangler triage:\n\n%s", jira.FormatIssues(issues))
		if !force {
			status = fmt.Sprintf("@team, %s", status)
		}
	} else {
		// Security issues are announced every
		// Wednesday. Calculate our latest Wednesdays.
		since := ((int(now.Weekday()) - int(time.Wednesday)) + 7) % 7
		lastWednesday := now.AddDate(0, 0, -since)

		// Calculate how many workdays have passed since last
		// Wednesday.
		workdays := c.CountWorkdays(lastWednesday, now) - 1

		// If this is forced or the first workday since
		// last Wednesday output that there are no issues.
		if force || workdays == 1 {
			status = fmt.Sprintf("Ingen issues mangler triage.")
		}
	}
	return status
}
