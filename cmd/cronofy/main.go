package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dcowgill/envflag"
	"github.com/jeffreylo/cronofy"
)

const dateFmt = "02 Jan --:--"
const dateTimeFmt = "02 Jan 15:04"

func main() {
	var (
		accessToken = flag.String("access-token", "", "cronofy.com")
		ids         = flag.String("calendar-ids", "", "calendar ids (comma separated)")
		timezone    = flag.String("timezone", "UTC", "iana recognized timezone")
	)
	flag.Parse()
	envflag.SetPrefix("cronofy")
	envflag.Parse()
	if *accessToken == "" {
		log.Fatal("access-token is required")
	}

	calendarIDs := strings.Split(*ids, ",")
	client := cronofy.NewClient(&cronofy.Config{
		AccessToken: *accessToken,
	})
	now := time.Now().UTC()
	from := now.Format("2006-01-02")
	end := now.Add(7 * 24 * time.Hour)
	to := end.Format("2006-01-02")
	res, err := client.GetEvents(&cronofy.EventsRequest{
		TZID:        "UTC",
		From:        &from,
		To:          &to,
		CalendarIDs: calendarIDs,
	})
	if err != nil {
		log.Fatal(err)
	}
	tz, _ := time.LoadLocation(*timezone)
	for _, v := range res.Events {
		if v.Accepted() {
			if v.AllDay {
				fmt.Printf("%s: %s\n", v.StartTime.Format(dateFmt), v.Summary)
			} else {
				if v.Location() != "" {
					fmt.Printf("%s: %s [%s]\n", v.StartTime.In(tz).Format(dateTimeFmt), v.Summary, v.Location())
				} else {
					fmt.Printf("%s: %s\n", v.StartTime.In(tz).Format(dateTimeFmt), v.Summary)
				}
			}
		}
	}
}
