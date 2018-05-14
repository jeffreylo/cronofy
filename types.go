package cronofy

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"
)

// Calendar are calendars supplied by providers linked to Cronofy.
type Calendar struct {
	ProviderName string `json:"provider_name"`
	ProfileID    string `json:"profile_id"`
	ProfileName  string `json:"profile_name"`
	CalendarID   string `json:"calendar_id"`
	CalendarName string `json:"calendar_name"`
	ReadOnly     bool   `json:"calendar_readonly"`
	Deleted      bool   `json:"calendar_deleted"`
}

type Pages struct {
	Current int    `json:"current"`
	Total   int    `json:"total"`
	Next    string `json:"next_page"`
}

type EventsResponse struct {
	Pages  *Pages   `json:"pages"`
	Events []*Event `json:"events"`
}

type EventsRequest struct {
	CalendarIDs    []string   `url:"calendar_ids[],omitempty"`
	From           *string    `url:"from,omitempty"`
	IncludeDeleted *bool      `url:"include_deleted,omitempty"`
	IncludeGeo     *bool      `url:"include_geo,omitempty"`
	IncludeMoved   *bool      `url:"include_moved,omitempty"`
	LastModified   *time.Time `url:"last_modified,omitempty"`
	LocalizedTimes *bool      `url:"localized_times,omitempty"`
	OnlyManaged    *bool      `url:"only_managed,omitempty"`
	TZID           string     `url:"tzid"`
	To             *string    `url:"to,omitempty"`
}

type Event struct {
	CalendarID  string    `json:"calendar_id"`
	EventUID    string    `json:"event_uid"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	Start       string    `json:"start"`
	End         string    `json:"end"`
	Deleted     bool      `json:"deleted"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Loc         struct {
		Description string `json:"description"`
	} `json:"location"`
	ParticipationStatus string `json:"participation_status"`
	Attendees           []struct {
		Email       string `json:"email"`
		DisplayName string `json:"display_name"`
		Status      string `json:"status"`
	} `json:"attendees"`
	Organizer struct {
		Email       string `json:"email"`
		DisplayName string `json:"display_name"`
	} `json:"organizer"`
	Transparency string   `json:"transparency"`
	Status       string   `json:"status"`
	Categories   []string `json:"categories"`
	Recurring    bool     `json:"recurring"`
	EventPrivate bool     `json:"event_private"`
	Options      struct {
		Delete                    bool `json:"delete"`
		Update                    bool `json:"update"`
		ChangeParticipationStatus bool `json:"change_participation_status"`
	} `json:"options"`

	StartTime *time.Time
	EndTime   *time.Time
	AllDay    bool
}

func (e *Event) UnmarshalJSON(b []byte) error {
	type event Event
	res := event{}
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}

	start, allDay, _ := parseDateTime(res.Start)
	end, _, _ := parseDateTime(res.End)
	res.StartTime = start
	res.EndTime = end
	res.AllDay = allDay
	*e = Event(res)
	return nil
}

func (e *Event) Accepted() bool {
	return e.ParticipationStatus == "accepted"
}

func (e *Event) Declined() bool {
	return e.ParticipationStatus == "declined"
}

func parseDateTime(v string) (*time.Time, bool, error) {
	var allDay bool
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		t, err = time.Parse("2006-01-02", v)
		if err != nil {
			panic(err)
		}
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		allDay = true
	}
	return &t, allDay, nil
}

var re = regexp.MustCompile(`\r?\n`)

func (e *Event) Location() string {
	description := strings.Replace(re.ReplaceAllString(e.Loc.Description, " "), ",", "", -1)
	return description
}
