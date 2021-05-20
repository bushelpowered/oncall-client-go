package oncall

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func (c *Client) GetTeams() ([]string, error) {
	teamList := []string{}
	_, err := c.Get("/api/v0/teams", &teamList)
	return teamList, errors.Wrap(err, "Fetching list of teams")
}

func (c *Client) GetTeam(name string) (Team, error) {
	t := Team{}
	_, err := c.Get("/api/v0/teams/"+name, &t)
	for rosterName, roster := range t.Rosters {
		roster.Name = rosterName
		t.Rosters[rosterName] = roster
	}
	return t, errors.Wrapf(err, "Fetching team deatils for %s", name)
}

func (c *Client) CreateTeam(t TeamConfig) (Team, error) {
	if t.Name == "" || t.SchedulingTimezone == "" {
		return Team{}, errors.New("You must define both the team Name and SchedulingTimezone")
	}
	log.Tracef("Going to create team %+v", t)
	_, createErr := c.Post("/api/v0/teams", t, nil)
	if createErr != nil {
		if strings.Contains(createErr.Error(), "HTTP Request failed (422)") {
			log.Error("Team already created")
		} else {
			return Team{}, errors.Wrapf(createErr, "Creating team %s", t.Name)
		}
	} else {
		log.Tracef("Successfully created team %+v", t)
	}

	createdTeam, getErr := c.GetTeam(t.Name)
	if createErr != nil {
		if getErr != nil {
			log.Error(errors.Wrap(getErr, "Getting team after failed create"))
		}
		return createdTeam, errors.Wrapf(createErr, "Creating team %s", t.Name)
	}
	return createdTeam, errors.Wrap(getErr, "Getting team after create")
}

func (c *Client) UpdateTeam(name string, t TeamConfig) (Team, error) {
	_, err := c.Put("/api/v0/teams/"+name, t, nil)
	if err != nil {
		return Team{}, errors.Wrapf(err, "Updating team %s", name)
	}

	teamName := name
	if t.Name != "" {
		teamName = t.Name
	}
	ret, err := c.GetTeam(teamName)
	return ret, errors.Wrapf(err, "Updating team %s", name)
}

func (c *Client) DeleteTeam(name string) error {
	existingTeam, err := c.GetTeam(name)
	if err != nil {
		return errors.Wrapf(err, "Failed to fetch team %s when attempting to delete", name)
	}

	existingTeam.TeamConfig.Name = fmt.Sprintf("%s-deleted-%d", name, time.Now().Unix())
	updatedTeam, err := c.UpdateTeam(name, existingTeam.TeamConfig)
	if err != nil {
		return errors.Wrapf(err, "Failed to rename team from %s to %s before delete", name, existingTeam.TeamConfig.Name)
	}

	_, err = c.Delete("/api/v0/teams/"+updatedTeam.Name, nil, nil)
	return errors.Wrapf(err, "Deleting team %s", name)
}
