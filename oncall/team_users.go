package oncall

import (
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// GetRosters returns a list of rosters by team
func (c *Client) GetTeamUsers(team string) ([]string, error) {
	loggerTeamUser("get", team, "").Trace("Getting team users")
	rosterUserList := []string{}
	url := fmt.Sprintf("/api/v0/teams/%s/users", team)
	_, err := c.Get(url, &rosterUserList)
	return rosterUserList, errors.Wrapf(err, "Fetching list of rosters for %s", team)
}

// Set Team Users authoritatviely sets the list of users for a team
func (c *Client) SetTeamUsers(team string, usernames []string) error {
	log := loggerTeamUser("set", team, "")
	log.Tracef("Setting users: %v", team, usernames)
	currentUsers, err := c.GetTeamUsers(team)
	if err != nil {
		return errors.Wrap(err, "Getting current list of team users for "+team)
	}

	usersToRemove, usersToAdd, _, _ := getSetVennDiagram(currentUsers, usernames)

	for _, u := range usersToAdd {
		err := c.AddTeamUser(team, u)
		if err != nil {
			return errors.Wrapf(err, "Adding user %s to user %s", u, team)
		}
	}

	for _, u := range usersToRemove {
		err := c.RemoveTeamUser(team, u)
		if err != nil {
			return errors.Wrapf(err, "Removing user %s to user %s", u, team)
		}
	}

	return nil
}

func (c *Client) AddTeamUser(team, username string) error {
	userUser := User{
		Name: username,
	}
	loggerTeamUser("add", team, username).Tracef("Adding user")
	url := fmt.Sprintf("/api/v0/teams/%s/users", team)
	_, err := c.Post(url, userUser, nil)
	return errors.Wrapf(err, "Adding user %s as user on %s", username, team)
}

func (c *Client) RemoveTeamUser(team, username string) error {
	loggerTeamUser("remove", team, username).Tracef("Removing user")
	url := fmt.Sprintf("/api/v0/teams/%s/users/%s", team, username)
	_, err := c.Delete(url, nil, nil)
	return errors.Wrapf(err, "Removing user %s as user on %s", username, team)
}

func loggerTeamUser(action, team, username string) *log.Entry {
	logger := log.WithFields(log.Fields{
		"action":   action,
		"type":     "team_user",
		"team":     team,
		"username": username,
	})
	return logger
}
