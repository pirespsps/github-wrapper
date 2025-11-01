package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
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

func getCommits(user string) Query {

	var url = "https://api.github.com/search/commits?q=author:" + user + "&merge=true"

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

	return commits

}

func formatByMonths(*Query) {

	type Month struct {
		Month  int
		Commit struct {
			Message    string
			Date       string
			Repository string
		}
	}

	var months = []Month{}

	for _, v := range &Query.Items {

	}

}
