package cli

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func ToDateString(s uint64) string {
	hours := s / 3600
	minutes := (s - (3600 * hours)) / 60
	seconds := s - (3600 * hours) - (minutes * 60)
	if hours != 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
}

func ToInt(s string) (int, error) {
	i, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return -1, err
	}
	return int(i), nil
}

func AskConfirm() bool {
	var response string

	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		fmt.Println("Please type (y)es or (n)o and then press enter:")
		return AskConfirm()
	}
}
