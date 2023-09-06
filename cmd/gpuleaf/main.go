package main

import (
	"bufio"
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
	gpuInfo, err := gpu.GetPowerInfo()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("GPU Info: %+v\n", gpuInfo)

	ticker := time.NewTicker(100 * time.Millisecond)

	logFile, err := os.OpenFile("gpu_info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	logWriter := io.MultiWriter(logFile, os.Stdout)
	logger := log.New(logWriter, "", log.LstdFlags)

	for range ticker.C {
		if err := logGPUStats(logger, gpuInfo); err != nil {
			logger.Printf("Error logging GPU stats: %v", err)
		}
	}
}

func logGPUStats(logger *log.Logger, gpuPowerInfo *gpu.GPUPowerInfo) error {
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
			gpuUtil := utils.ParseFloat(fields[2])
			powerLimit := utils.ParseFloat(fields[3])

			gpuStats := gpu.GPUStats{
				Temperature: temperature,
				PowerDraw:   powerDraw,
				Utilization: gpuUtil,
				PowerLimit:  powerLimit,
			}

			logger.Printf("GPU Info: %+v", gpuStats)
			gpu.Leaf(gpuPowerInfo, &gpuStats)
		}
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
