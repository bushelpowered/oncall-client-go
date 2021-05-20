package oncall

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// GetRosters returns a list of rosters by team
func (c *Client) GetRosters(team string) ([]string, error) {
	rosterList := make(map[string]interface{})
	url := fmt.Sprintf("/api/v0/teams/%s/rosters", team)
	_, err := c.Get(url, &rosterList)
	ret := []string{}
	for r := range rosterList {
		ret = append(ret, r)
	}
	return ret, errors.Wrapf(err, "Fetching list of rosters for %s", team)
}

func (c *Client) GetRoster(team, name string) (Roster, error) {
	roster := Roster{}
	url := fmt.Sprintf("/api/v0/teams/%s/rosters/%s", team, name)
	_, err := c.Get(url, &roster)
	roster.Name = name
	return roster, errors.Wrapf(err, "Fetching roster deatils for %s/%s", team, name)
}

func (c *Client) CreateRoster(team, name string) (Roster, error) {
	roster := Roster{
		Name: name,
	}

	log.Tracef("Going to create roster %s/%s", team, name)
	url := fmt.Sprintf("/api/v0/teams/%s/rosters", team)
	_, createErr := c.Post(url, roster, nil)
	if createErr != nil {
		if strings.Contains(createErr.Error(), "HTTP Request failed (422)") {
			log.Error("Roster already created")
		} else {
			return roster, errors.Wrapf(createErr, "Creating roster %s", roster.Name)
		}
	} else {
		log.Tracef("Successfully created roster %s/%s", team, name)
	}

	createdRoster, getErr := c.GetRoster(team, roster.Name)
	if createErr != nil {
		if getErr != nil {
			log.Error(errors.Wrap(getErr, "Getting roster after failed create"))
		}
		return createdRoster, errors.Wrapf(createErr, "Creating roster %s", roster.Name)
	}
	return createdRoster, errors.Wrap(getErr, "Getting roster after create")
}

func (c *Client) UpdateRoster(team, name string, roster Roster) (Roster, error) {
	url := fmt.Sprintf("/api/v0/teams/%s/rosters/%s", team, name)
	_, err := c.Put(url, roster, nil)
	if err != nil {
		return roster, errors.Wrapf(err, "Updating roster %s/%s", team, name)
	}

	rosterName := name
	if roster.Name != "" {
		rosterName = roster.Name
	}
	ret, err := c.GetRoster(team, rosterName)
	return ret, errors.Wrapf(err, "Updating roster %s/%s", team, name)
}

func (c *Client) DeleteRoster(team, name string) error {
	url := fmt.Sprintf("/api/v0/teams/%s/rosters/%s", team, name)
	_, err := c.Delete(url, nil, nil)
	return errors.Wrapf(err, "Deleting roster %s/%s", team, name)
}
