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

	return formatByMonths(&commits)

}

func formatByMonths(query *Query) map[int][]bool {

	months := make(map[int][]bool)

	var currentMonth int
	var monthCheck = make([]bool, 31)

	for _, v := range query.Items {

		month, err := strconv.Atoi(v.Commit.Author.Date[5:7])
		if err != nil {
			log.Fatal(err)
		}

		day, err := strconv.Atoi(v.Commit.Author.Date[8:10])
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("mes: %v dia: %v\n", month, day)

		//está começando com o primeiro mes inteiro zerado devido a inicialização
		//não está verificando mês a mês, apenas os que tem pelo menos uma commit
		//fazer um for para ir zerando até chegar na primeira commit
		//mudar a query da api para só puxar do ano
		if month != currentMonth {
			totalDays := daysIn(time.Month(currentMonth), time.Now().Year())

			for i := range totalDays {

				if monthCheck[i] {
					continue
				}
				monthCheck[i] = false
			}

			months[currentMonth] = monthCheck
			monthCheck = make([]bool, 31)
			currentMonth = month

		}

		monthCheck[day] = true

	}

	return months

}

func daysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
