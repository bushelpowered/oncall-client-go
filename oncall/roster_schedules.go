package oncall

import (
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// GetRosterSchedules returrns a list of roster schedules with the key being the role and the value being the schedule
func (c *Client) GetRosterSchedules(team, roster string) (map[string]Schedule, error) {
	loggerRosterSchedules("getall", team, roster, "").Trace("Geting all roster schedules")
	ret := map[string]Schedule{}
	rosterScheduleList := []Schedule{}
	url := fmt.Sprintf("/api/v0/teams/%s/rosters/%s/schedules", team, &rosterScheduleList)
	_, err := c.Get(url, &rosterScheduleList)

	for _, s := range rosterScheduleList {
		loggerRosterSchedules("getall", team, roster, s.Role).Trace("Found role")
		ret[s.Role] = s
	}
	return ret, errors.Wrapf(err, "Fetching list of rosters for %s", team)
}

func (c *Client) AddRosterSchedule(team, roster string, schedule Schedule) error {
	loggerRosterSchedules("add", team, roster, schedule.Role).Trace("Going to add")
	url := fmt.Sprintf("/api/v0/teams/%s/rosters/%s/schedules", team, roster)
	_, err := c.Post(url, schedule, nil)
	return errors.Wrapf(err, "Adding schedule %s to roster %s/%s", schedule.Role, team, roster)
}

func (c *Client) GetRosterSchedule(team, roster, scheduleRole string) (Schedule, error) {
	loggerRosterSchedules("get", team, roster, scheduleRole).Trace("Getting schedule")

	allSchedules, err := c.GetRosterSchedules(team, roster)
	if err != nil {
		return Schedule{}, errors.Wrap(err, "getting all schedules for roster")
	}

	schedule, ok := allSchedules[scheduleRole]
	if !ok {
		return Schedule{}, errors.New("Did not find schedule")
	}

	return schedule, nil
}

func (c *Client) RemoveRosterScheduleByID(team, roster string, scheduleID int) error {
	logger := loggerRosterSchedules("deleteByID", team, roster, fmt.Sprintf("%d", scheduleID))
	logger.Trace("Going to delete")

	url := fmt.Sprintf("/api/v0/teams/%s/rosters/%s/schedules/%s", team, roster, scheduleID)
	_, err := c.Delete(url, roster, nil)
	return errors.Wrapf(err, "Removing schedule id %s from roster %s/%s", scheduleID, team, roster)
}

func (c *Client) RemoveRosterSchedule(team, roster, scheduleRole string) error {
	logger := loggerRosterSchedules("delete", team, roster, scheduleRole)
	logger.Trace("Fetching schedule for delete")

	schedule, err := c.GetRosterSchedule(team, roster, scheduleRole)
	if err != nil {
		logger.Trace("Could not find schedule for delete")
		return errors.Wrap(err, "Getting schedule for delete")
	}

	return c.RemoveRosterScheduleByID(team, roster, schedule.ID)
}

func loggerRosterSchedules(action, team, roster, schedule string) *log.Entry {
	logger := log.WithFields(log.Fields{
		"action":   action,
		"type":     "roster_schedules",
		"team":     team,
		"roster":   roster,
		"schedule": schedule,
	})
	return logger
}
