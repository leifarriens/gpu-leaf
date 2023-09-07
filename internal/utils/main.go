package utils

import (
	"io"
	"log"
	"os"
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

func CreateLogger(stdout bool, file bool) (*log.Logger, *os.File) {
	var logFile *os.File

	var writers []io.Writer

	if file {
		logFile, err := os.OpenFile("gpu_leaf.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			log.Fatalf("Error opening log file: %v", err)
		}

		writers = append(writers, logFile)
	}

	if stdout {
		writers = append(writers, os.Stdout)
	}

	logWriter := io.MultiWriter(writers...)

	logger := log.New(logWriter, "", log.LstdFlags)

	return logger, logFile
}
