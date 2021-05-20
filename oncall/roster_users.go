package oncall

import (
	"fmt"

	"github.com/pkg/errors"
)

// GetRosters returns a list of rosters by team
func (c *Client) GetRosterUsers(team, roster string) ([]string, error) {
	rosterUserList := []string{}
	url := fmt.Sprintf("/api/v0/teams/%s/rosters/%s/users", team, roster)
	_, err := c.Get(url, &rosterUserList)
	return rosterUserList, errors.Wrapf(err, "Fetching list of rosters for %s", team)
}

func (c *Client) SetRosterUsers(team, roster string, usernames []string) error {
	log.Tracef("Goign to set roster %s/%s users to: %v", team, roster, usernames)
	currentUsers, err := c.GetRosterUsers(team, roster)
	if err != nil {
		return errors.Wrap(err, "Getting current list of roster users")
	}

	usersToRemove, usersToAdd, _, _ := getSetVennDiagram(currentUsers, usernames)

	for _, u := range usersToAdd {
		err := c.AddRosterUser(team, roster, u)
		if err != nil {
			return errors.Wrapf(err, "Adding user %s", u)
		}
	}

	for _, u := range usersToRemove {
		err := c.RemoveRosterUser(team, roster, u)
		if err != nil {
			return errors.Wrapf(err, "Removing user %s", u)
		}
	}

	return nil
}

func (c *Client) AddRosterUser(team, roster, username string) error {
	rosterUser := RosterUser{
		Name:       username,
		InRotation: true,
	}

	log.Tracef("Going to add %s to roster %s/%s", username, team, roster)
	url := fmt.Sprintf("/api/v0/teams/%s/rosters/%s/users", team, roster)
	_, err := c.Post(url, rosterUser, nil)
	return errors.Wrapf(err, "Adding user %s to roster %s/%s", username, team, roster)
}

func (c *Client) RemoveRosterUser(team, roster, username string) error {
	url := fmt.Sprintf("/api/v0/teams/%s/rosters/%s/users/%s", team, roster, username)
	_, err := c.Delete(url, roster, nil)
	return errors.Wrapf(err, "Removing user %s from roster %s/%s", username, team, roster)
}
