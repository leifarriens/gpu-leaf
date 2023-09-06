package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type GPUStats struct {
	Temperature float64
	PowerDraw   float64
	Utilization float64
	PowerLimit  float64
}

func main() {
	ticker := time.NewTicker(500 * time.Millisecond)

	logFile, err := os.OpenFile("gpu_info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	logWriter := io.MultiWriter(logFile, os.Stdout)
	logger := log.New(logWriter, "", log.LstdFlags)

	for range ticker.C {
		if err := logGPUInfo(logger); err != nil {
			logger.Printf("Error logging GPU info: %v", err)
		}
	}
}

func logGPUInfo(logger *log.Logger) error {
	cmd := exec.Command("nvidia-smi", "--query-gpu=temperature.gpu,power.draw,utilization.gpu,power.limit", "--format=csv,noheader,nounits")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")
		if len(fields) == 4 {
			temperature := parseFloat(fields[0])
			powerDraw := parseFloat(fields[1])
			gpuUtil := parseFloat(fields[2])
			powerLimit := parseFloat(fields[3])

			gpuStats := GPUStats{
				Temperature: temperature,
				PowerDraw:   powerDraw,
				Utilization: gpuUtil,
				PowerLimit:  powerLimit,
			}

			logger.Printf("GPU Info: %+v", gpuStats)
		}
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func parseFloat(s string) float64 {
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		log.Printf("Error parsing float: %v", err)
	}
	return f
}
