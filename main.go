package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func main() {
	commitsPerMonth := getCommits(User)

	for i, v := range commitsPerMonth {

		month := time.Month(i).String()
		fmt.Print(month + "\n")

		for i, value := range v {

			if i == 31 {
				continue
			}

			if value {
				fmt.Print(lipgloss.NewStyle().Foreground(lipgloss.Color("#66ff00")).Render(strconv.Itoa(i+1) + " "))
			} else {
				fmt.Print(lipgloss.NewStyle().Foreground(lipgloss.Color("#d3000000")).Render(strconv.Itoa(i+1) + " "))
			}

		}
		fmt.Print("\n\n")
	}
}
