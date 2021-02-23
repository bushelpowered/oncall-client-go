package oncall

import (
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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

	// Create two "sets", one of current users and one of target users
	setCurrentUsers := map[string]bool{}
	setTargetUsers := map[string]bool{}
	for _, u := range usernames {
		setTargetUsers[u] = true
	}
	for _, u := range currentUsers {
		setCurrentUsers[u] = true
	}

	// If a user is in the list of current users but not in the list of target users
	// then we ned to remove that user
	usersToRemove := []string{}
	for u, _ := range setCurrentUsers {
		_, targetUser := setTargetUsers[u]
		if !targetUser {
			usersToRemove = append(usersToRemove, u)
		}
	}

	// If a user is in the list of target users but not in the list of current users
	// Then we need to add that user
	usersToAdd := []string{}
	for u, _ := range setTargetUsers {
		_, targetUser := setCurrentUsers[u]
		if !targetUser {
			usersToAdd = append(usersToAdd, u)
		}
	}

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
		Name: username,
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
