package oncall

type Team struct {
	TeamConfig
	Admins   []User            `json:"admins"`
	ID       int               `json:"id"`
	Rosters  map[string]Roster `json:"rosters"`
	Services []string          `json:"services"`
	Users    map[string]User   `json:"users"`
}

type TeamConfig struct {
	Name               string `json:"name"`
	Email              string `json:"email"`
	SlackChannel       string `json:"slack_channel"`
	IrisPlan           string `json:"iris_plan"`
	SchedulingTimezone string `json:"scheduling_timezone"`
}

type ScheduleEvent struct {
	Duration int `json:"duration"`
	Start    int `json:"start"`
}

type ScheduleScheduler struct {
	Name string `json:"name"`
}
type Schedule struct {
	AdvancedMode          int               `json:"advanced_mode"`
	AutoPopulateThreshold int               `json:"auto_populate_threshold"`
	Scheduler             ScheduleScheduler `json:"scheduler"`
	Events                []ScheduleEvent   `json:"events"`
	ID                    int               `json:"id"`
	Role                  string            `json:"role"`
	Roster                string            `json:"roster"`
	Team                  string            `json:"team"`
	Timezone              string            `json:"timezone"`
}

type RosterUser struct {
	InRotationInt int    `json:"in_rotation"`
	InRotation    bool   `json:"in_rotation"`
	Name          string `json:"name"`
}

type Roster struct {
	Name      string       `json:"name"`
	ID        int          `json:"id"`
	Schedules []Schedule   `json:"schedules"`
	Users     []RosterUser `json:"users"`
}

type Contacts struct {
	Call  string `json:"call"`
	Email string `json:"email"`
	Im    string `json:"im"`
	Sms   string `json:"sms"`
}

type User struct {
	Active   int      `json:"active"`
	Contacts Contacts `json:"contacts"`
	FullName string   `json:"full_name"`
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	PhotoURL string   `json:"photo_url"`
	TimeZone string   `json:"time_zone"`
}
