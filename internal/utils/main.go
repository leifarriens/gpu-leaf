package utils

import (
	"log"
	"strconv"
	"strings"
)

func ParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		log.Printf("Error parsing float: %v", err)
	}
	return f
}

func ParseInt(s string) int {
	f, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("Error parsing int: %v", err)
	}
	return f
}
