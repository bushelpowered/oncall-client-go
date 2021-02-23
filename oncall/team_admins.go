package oncall

import (
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// GetRosters returns a list of rosters by team
func (c *Client) GetTeamAdmins(team string) ([]string, error) {
	loggerTeamAdmin("get", team, "").Trace("Getting team admins")
	rosterUserList := []string{}
	url := fmt.Sprintf("/api/v0/teams/%s/admins", team)
	_, err := c.Get(url, &rosterUserList)
	return rosterUserList, errors.Wrapf(err, "Fetching list of rosters for %s", team)
}

// Set Team Admins authoritatviely sets the list of admins for a team
func (c *Client) SetTeamAdmins(team string, usernames []string) error {
	log := loggerTeamAdmin("set", team, "")
	log.Tracef("Setting admins: %v", team, usernames)
	currentUsers, err := c.GetTeamAdmins(team)
	if err != nil {
		return errors.Wrap(err, "Getting current list of team admins for "+team)
	}

	usersToRemove, usersToAdd, _, _ := getSetVennDiagram(currentUsers, usernames)

	for _, u := range usersToAdd {
		err := c.AddTeamAdmin(team, u)
		if err != nil {
			return errors.Wrapf(err, "Adding user %s to admin %s", u, team)
		}
	}

	for _, u := range usersToRemove {
		err := c.RemoveTeamAdmin(team, u)
		if err != nil {
			return errors.Wrapf(err, "Removing user %s to admin %s", u, team)
		}
	}

	return nil
}

func (c *Client) AddTeamAdmin(team, username string) error {
	adminUser := User{
		Name: username,
	}
	loggerTeamAdmin("add", team, username).Tracef("Adding admin")
	url := fmt.Sprintf("/api/v0/teams/%s/admins", team)
	_, err := c.Post(url, adminUser, nil)
	return errors.Wrapf(err, "Adding user %s as admin on %s", username, team)
}

func (c *Client) RemoveTeamAdmin(team, username string) error {
	loggerTeamAdmin("remove", team, username).Tracef("Removing admin")
	url := fmt.Sprintf("/api/v0/teams/%s/admins/%s", team, username)
	_, err := c.Delete(url, nil, nil)
	return errors.Wrapf(err, "Removing user %s as admin on %s", username, team)
}

func loggerTeamAdmin(action, team, username string) *log.Entry {
	logger := log.WithFields(log.Fields{
		"action":   action,
		"type":     "teamAdmin",
		"team":     team,
		"username": username,
	})
	return logger
}
