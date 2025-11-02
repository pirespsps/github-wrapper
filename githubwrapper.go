package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Query struct {
	Items []struct {
		Commit struct {
			Message string `json:"message"`
			Author  struct {
				Date string `json:"date"`
			}
		}
		Repository struct {
			FullName string `json:"full_name"`
		}
	}
}

func getCommits(user string) map[int][]bool {

	year := fmt.Sprint(time.Now().Year())

	var url = "https://api.github.com/search/commits?q=author:" + user + "&merge=true&since:" + year + "01-01T00:00:00Z"

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var commits Query

	json.Unmarshal(data, &commits)

	return formatByMonths(&commits)

}

func formatByMonths(query *Query) map[int][]bool {

	months := make(map[int][]bool)

	currentMonth := 1
	var monthCheck = make([]bool, 32)

	for _, v := range query.Items {

		month, err := strconv.Atoi(v.Commit.Author.Date[5:7])
		if err != nil {
			log.Fatal(err)
		}

		day, err := strconv.Atoi(v.Commit.Author.Date[8:10])
		if err != nil {
			log.Fatal(err)
		}

		if month > currentMonth {

			for j := range month - currentMonth {

				totalDays := daysIn(time.Month(currentMonth+j), time.Now().Year())

				for i := range totalDays {

					if monthCheck[i] {
						continue
					}
					monthCheck[i] = false
				}
				months[currentMonth+j] = monthCheck
			}

			monthCheck = make([]bool, daysIn(time.Month(month), time.Now().Year()))
			currentMonth = month

		}

		monthCheck[day] = true

	}

	months[currentMonth] = monthCheck

	return months

}

func daysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
