package main

import (
	"os"
	"strings"

	"github.com/bushelpowered/oncall-client-go/oncall"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.TraceLevel)
	oc, err := oncall.New(nil, oncall.Config{
		Username:   os.Getenv("ONCALL_USERNAME"),
		Password:   os.Getenv("ONCALL_PASSWORD"),
		Endpoint:   os.Getenv("ONCALL_ENDPOINT"),
		AuthMethod: oncall.AuthMethodUser,
	})
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to create oncall client"))
	}

	err = createDeleteTeam(oc)
	if err != nil {
		log.Fatal(err)
	}
}

func createDeleteTeam(oc *oncall.Client) error {
	t, err := oc.CreateTeam(oncall.TeamConfig{
		Name:               "go-client-test-team",
		SchedulingTimezone: "US/Central",
	})
	if err != nil {
		if !strings.Contains(err.Error(), "HTTP Request failed (422)") {
			return errors.Wrap(err, "Creating team")
		}
		log.Info("Team was already created")
	}
	log.Infof("Created team: %+v", t)
	err = oc.DeleteTeam(t.Name)
	return errors.Wrap(err, "Deleting team "+t.Name)
}

func listAllTeams(oc *oncall.Client) error {
	teams, err := oc.GetTeams()
	if err != nil {
		return (errors.Wrap(err, "Failed to get teams"))
	}

	for _, t := range teams {
		log.Print("Found team " + t)
		team, err := oc.GetTeam(t)
		if err != nil {
			return errors.Wrapf(err, "Failed to get team %s", t)
		}
		log.Debugf("Team: %+v", team)
	}

	return nil
}