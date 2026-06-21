package team

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"arnested.dk/go/triagebot/internal/zulip"
	"github.com/joefitzgerald/forecast"
)

func Triage(ctx context.Context) string {
	api := forecast.New(
		"https://api.forecastapp.com",
		os.Getenv("FORECAST_ACCOUNT_ID"),
		os.Getenv("FORECAST_ACCESS_TOKEN"),
	)

	today := time.Now().Format(time.DateOnly)

	assignments, err := api.AssignedPeople(today, today)
	if err != nil {
		slog.Error("failed to retrieve assigned people", "error", err)

		return ""
	}

	targetProjectID := os.Getenv("FORECAST_TEAM")

	var triageTeam strings.Builder

	if personIDs, ok := assignments[targetProjectID]; ok && len(personIDs) > 0 {
		for _, personID := range personIDs {
			person, err := api.Person(personID)
			if err != nil {
				slog.Warn("could not fetch person details", "personID", personID)

				continue
			}

			tag, err := zulip.UserByEmail(ctx, person.Email)
			if err != nil {
				triageTeam.WriteString(person.Email + "\n")
			} else {
				triageTeam.WriteString(tag + "\n")
			}
		}
	} else {
		triageTeam.WriteString("@**all**, no one is assigned to triage today!\n")
	}

	return triageTeam.String()
}
