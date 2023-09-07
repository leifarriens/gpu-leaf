package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/leifarriens/gpu-leaf/internal/gpu"
	"github.com/leifarriens/gpu-leaf/internal/utils"
)

func main() {
	interval := flag.Int("l", 1000, "smi update interval")
	threshold := flag.Int("t", 95, "Utilization threshold")
	shouldOc := flag.Bool("oc", false, "Should gpu-leave raise power level above 100%")

	flag.Parse()

	gpuPowerInfo, err := gpu.GetPowerInfo()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("GPU Info: %+v\n", gpuPowerInfo)

	if !gpuPowerInfo.IsPowerManageable {
		log.Fatal("GPU does not support power management")
	}

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

	ticker := time.NewTicker(time.Duration(*interval) * time.Millisecond)

	logger, logFile := utils.CreateLogger(true, true) // TODO: make flag configurable

	defer logFile.Close()

	for range ticker.C {
		if err := gpu.WatchStats(&gpuConfig, logger, gpu.Leaf); err != nil {
			log.Fatalf("Error receiving GPU stats: %v", err)
		}
	}
}
