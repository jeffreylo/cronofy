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
	to := now.Add(7 * 24 * time.Hour)
	res, err := client.GetEvents(&cronofy.EventsRequest{
		TZID:        "UTC",
		From:        &now,
		To:          &to,
		CalendarIDs: calendarIDs,
	})
	if err != nil {
		log.Fatal(err)
	}
	tz, _ := time.LoadLocation(*timezone)
	for _, v := range res.Events {
		if v.Accepted() && v.StartTime().After(now) {
			if v.Location() != "" {
				fmt.Printf("%s: %s [%s]\n", v.StartTime().In(tz).Format(time.RFC822), v.Summary, v.Location())
			} else {
				fmt.Printf("%s: %s\n", v.StartTime().In(tz).Format(time.RFC822), v.Summary)
			}
		}
	}
}
