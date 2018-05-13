package cronofy

import (
	"log"
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
	From           *time.Time `url:"from,omitempty"`
	IncludeDeleted *bool      `url:"include_deleted,omitempty"`
	IncludeGeo     *bool      `url:"include_geo,omitempty"`
	IncludeMoved   *bool      `url:"include_moved,omitempty"`
	LastModified   *time.Time `url:"last_modified,omitempty"`
	LocalizedTimes *bool      `url:"localized_times,omitempty"`
	OnlyManaged    *bool      `url:"only_managed,omitempty"`
	TZID           string     `url:"tzid"`
	To             *time.Time `url:"to,omitempty"`
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
}

func (e *Event) Accepted() bool {
	return e.ParticipationStatus == "accepted"
}

func (e *Event) Declined() bool {
	return e.ParticipationStatus == "declined"
}

func (e *Event) StartTime() *time.Time {
	t, err := time.Parse(time.RFC3339, e.Start)
	if err != nil {
		t, err = time.Parse("2006-01-02", e.Start)
		if err != nil {
			log.Println(err)
		}
	}
	return &t
}

var re = regexp.MustCompile(`\r?\n`)

func (e *Event) Location() string {
	description := strings.Replace(re.ReplaceAllString(e.Loc.Description, " "), ",", "", -1)
	return description
}
