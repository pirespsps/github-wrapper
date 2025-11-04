package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Query struct {
	Items []struct {
		Commit struct {
			Author struct {
				Date string `json:"date"`
			}
		}
	}
}

func getCommits(user string) map[int][]bool {

	var wg sync.WaitGroup
	queries := make(chan Query, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			query := paginatedFetch(user, page)
			queries <- query
		}(i + 1)
	}

	go func() {
		wg.Wait()
		close(queries)
	}()

	var RealQuery Query

	for query := range queries {
		RealQuery.Items = append(RealQuery.Items, query.Items...)
	}

	return formatByMonths(&RealQuery)

}

func paginatedFetch(user string, page int) Query {

	year := fmt.Sprint(time.Now().Year())
	pageStr := strconv.Itoa(page)

	var url = "https://api.github.com/search/commits?q=author:" + user + "+committer-date:>=" + year + "-01-01T00:00:00Z" + "&sort=author-date&order=asc&per_page=100&page=" + pageStr

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accept", "application/vnd.github.cloak-preview")

	client := &http.Client{Timeout: 10 * time.Second}

	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var commits Query

	json.Unmarshal(data, &commits)

	return commits
}

func formatByMonths(query *Query) map[int][]bool {
	months := make(map[int][]bool)
	year := time.Now().Year()

	for _, v := range query.Items {

		month, err := strconv.Atoi(v.Commit.Author.Date[5:7])
		if err != nil {
			log.Fatal(err)
		}

		day, err := strconv.Atoi(v.Commit.Author.Date[8:10])
		if err != nil {
			log.Fatal(err)
		}

		if _, value := months[month]; !value {
			months[month] = make([]bool, daysIn(time.Month(month), year))
		}

		if (day - 1) < len(months[month]) {
			months[month][day-1] = true
		}

	}

	return months
}

func daysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
