package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/leifarriens/gpu-leaf/internal/gpu"
	"github.com/leifarriens/gpu-leaf/internal/utils"
)

func main() {
	gpuPowerInfo, err := gpu.GetPowerInfo()

	// this flag decides if PowerLevel will go above DefaultPowerLevel up to MaxPowerLevel
	shouldOc := flag.Bool("oc", false, "Should gpu-leave raise power level above 100%")
	threshold := flag.Int("t", 95, "Utilization threshold")

	flag.Parse()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("GPU Info: %+v\n", gpuPowerInfo)

	var maxPowerLimit float64

	if *shouldOc {
		maxPowerLimit = gpuPowerInfo.MaxPowerLimit
	} else {
		maxPowerLimit = gpuPowerInfo.DefaultPowerLimit
	}

	gpuConfig := gpu.GPUConfig{
		MinPowerLimit: gpuPowerInfo.MinPowerLimit,
		MaxPowerLimit: maxPowerLimit,
		Threshold:     *threshold,
	}

	ticker := time.NewTicker(1000 * time.Millisecond)

	logFile, err := os.OpenFile("gpu_info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	defer logFile.Close()

	logWriter := io.MultiWriter(logFile, os.Stdout)
	logger := log.New(logWriter, "", log.LstdFlags)

	for range ticker.C {

		if err := logGPUStats(logger, &gpuConfig); err != nil {
			log.Fatalf("Error receiving GPU stats: %v", err)
		}
	}
}

func logGPUStats(logger *log.Logger, gpuConfig *gpu.GPUConfig) error {
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
			temperature := utils.ParseFloat(fields[0])
			powerDraw := utils.ParseFloat(fields[1])
			gpuUtil := int(utils.ParseFloat(fields[2]))
			powerLimit := utils.ParseFloat(fields[3])

			gpuStats := gpu.GPUStats{
				Temperature: temperature,
				PowerDraw:   powerDraw,
				Utilization: gpuUtil,
				PowerLimit:  powerLimit,
			}

			logger.Printf("GPU Info: %+v", gpuStats)
			gpu.Leaf(gpuConfig, &gpuStats)
		}
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
