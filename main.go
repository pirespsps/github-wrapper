package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func main() {
	commitsPerMonth := getCommits(User)
	terminal(&commitsPerMonth)
}

func terminal(commits *map[int][]bool) {
	styleTrue := lipgloss.NewStyle().
		Background(lipgloss.Color("#004d02")).
		Foreground(lipgloss.Color("#ffffff")).
		Align(lipgloss.Center).
		Padding(0, 1)

	styleFalse := lipgloss.NewStyle().
		Background(lipgloss.Color("#800000")).
		Foreground(lipgloss.Color("#ffffff")).
		Align(lipgloss.Center).
		Padding(0, 1)

	for i, v := range *commits {

		month := time.Month(i).String()
		fmt.Print(month + "\n")

		for i, value := range v {

			if i == 31 {
				continue
			}

			if value {
				fmt.Print(styleTrue.Render(strconv.Itoa(i + 1)))
			} else {
				fmt.Print(styleFalse.Render(strconv.Itoa(i + 1)))
			}

		}
		fmt.Print("\n\n")
	}
}

func makeWidget(commits *map[int][]bool) []byte {

	jsonObj, err := json.Marshal(commits)
	if err != nil {
		log.Fatal(err)
	}

	return jsonObj

}
