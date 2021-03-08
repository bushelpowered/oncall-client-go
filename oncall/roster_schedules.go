package oncall

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// GetRosterSchedules returrns a list of roster schedules with the key being the role and the value being the schedule
// GET /api/v0/teams/{team}/rosters/{roster}/schedules
// Get schedules for a given roster. Information on schedule attributes is detailed in the schedules POST endpoint documentation. Schedules can be filtered with the following parameters passed in the query string:
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

// AddRosterSchedule creates a new roster schedule
// POST /api/v0/teams/{team}/rosters/{roster}/schedules
// Details here: https://oncall.tools/docs/api.html#post--api-v0-teams-team-rosters-roster-schedules
func (c *Client) AddRosterSchedule(team, roster string, schedule Schedule) error {
	loggerRosterSchedules("add", team, roster, schedule.Role).Trace("Going to add")
	url := fmt.Sprintf("/api/v0/teams/%s/rosters/%s/schedules", team, roster)
	_, err := c.Post(url, schedule, nil)
	return errors.Wrapf(err, "Adding schedule %s to roster %s/%s", schedule.Role, team, roster)
}

// GetRosterSchedule loops through the existing schedules and finds the one that matches the role you're looking for
func (c *Client) GetRosterSchedule(team, roster, scheduleRole string) (Schedule, error) {
	loggerRosterSchedules("get", team, roster, scheduleRole).Trace("Getting schedule")

	allSchedules, err := c.GetRosterSchedules(team, roster)
	if err != nil {
		return Schedule{}, errors.Wrap(err, "getting all schedules for roster")
	}

	for _, sched := range allSchedules {
		if strings.ToLower(sched.Role) == strings.ToLower(scheduleRole) {
			return sched, nil
		}
	}
	return Schedule{}, errors.New("Did not find schedule")
}

// UpdateRosterSchedule updates a roster scheudle for a team/roster pair.
// PUT /api/v0/schedules/{schedule_id}
// Update a schedule. Allows editing of role, team, roster, auto_populate_threshold, events, and advanced_mode. Only allowed for team admins. Note that simple mode schedules must conform to simple schedule restrictions (described in documentation for the /api/v0/team/{team_name}/rosters/{roster_name}/schedules GET endpoint). This is checked on both “events” and “advanced_mode” edits.
func (c *Client) UpdateRosterSchedule(team, roster, role string, schedule Schedule) error {
	loggerRosterSchedules("update", team, roster, role).Trace("Getting existing schedule")
	currSchedule, err := c.GetRosterSchedule(team, roster, role)
	if err != nil {
		return errors.Wrapf(err, "Getting schedule for update")
	}

	url := fmt.Sprintf("/api/v0/schedules/%d", currSchedule.ID)
	_, err = c.Put(url, schedule, nil)
	return errors.Wrapf(err, "Adding schedule %s to roster %s/%s", schedule.Role, team, roster)
}

// PopulateRosterSchedule will trigger a plan for schedule starting from the time that is specified
// POST /api/v0/schedules/{schedule_id}/populate
// Run the scheduler on demand from a given point in time. Deletes existing schedule events if applicable. Given the start param, this will find the first schedule start time after start, then populate out to the schedule’s auto_populate_threshold. It will also clear the calendar of any events associated with the chosen schedule from the start of the first event it created onward. For example, if start is Monday, May 1 and the chosen schedule starts on Wednesday, this will create events starting from Wednesday, May 3, and delete any events that start after May 3 that are associated with the schedule.
// `start` should be a unix timestamp after time.Now()
func (c *Client) PopulateRosterSchedule(team, roster, role string, startTime time.Time) error {
	loggerRosterSchedules("populate", team, roster, role).Trace("Getting existing schedule")
	if !startTime.After(time.Now().Add(-1 * time.Second)) {
		return fmt.Errorf("Populate time must be after time.Now()")
	}

	currSchedule, err := c.GetRosterSchedule(team, roster, role)
	if err != nil {
		return errors.Wrapf(err, "Getting schedule for update")
	}

	populateBody := map[string]int{
		"start": int(startTime.Unix()),
	}
	loggerRosterSchedules("populate", team, roster, role).Trace("Going to populate")
	url := fmt.Sprintf("/api/v0/schedules/%d/populate", currSchedule.ID)
	_, err = c.Post(url, populateBody, nil)
	return errors.Wrapf(err, "Populating schedule %s to roster %s/%s", role, team, roster)
}

// RemoveRosterScheduleByID uses the id to delete rather than the name
// DELETE /api/v0/schedules/{schedule_id}
// Delete a schedule by id. Only allowed for team admins.
func (c *Client) RemoveRosterScheduleByID(scheduleID int) error {
	loggerRosterSchedules("deleteByID", "", "", fmt.Sprintf("%d", scheduleID)).Trace("Going to delete")

	url := fmt.Sprintf("/api/v0/schedules/%d", scheduleID)
	_, err := c.Delete(url, nil, nil)
	return errors.Wrapf(err, "Removing schedule id %s", scheduleID)
}

// RemoveRosterSchedule is a helper function for removeing a roster by id
func (c *Client) RemoveRosterSchedule(team, roster, scheduleRole string) error {
	logger := loggerRosterSchedules("delete", team, roster, scheduleRole)
	logger.Trace("Fetching schedule for delete")

	schedule, err := c.GetRosterSchedule(team, roster, scheduleRole)
	if err != nil {
		logger.Trace("Could not find schedule for delete")
		return errors.Wrap(err, "Getting schedule for delete")
	}

	return c.RemoveRosterScheduleByID(schedule.ID)
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
