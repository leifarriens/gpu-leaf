package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/leifarriens/gpu-leaf/internal/gpu"
	"github.com/leifarriens/gpu-leaf/internal/utils"
	"github.com/leifarriens/gpu-leaf/internal/version"
)

func main() {
	interval := flag.Int("interval", 1000, "Polling interval in milliseconds")
	flag.IntVar(interval, "l", 1000, "(alias) Polling interval in milliseconds")
	threshold := flag.Int("threshold", 95, "Utilization percentage threshold to decide raising vs lowering power limit")
	flag.IntVar(threshold, "t", 95, "(alias) Utilization threshold")
	shouldOc := flag.Bool("overclock", false, "Allow using max power limit (potentially above default). Use cautiously.")
	flag.BoolVar(shouldOc, "oc", false, "(alias) Allow using max power limit")
	gpuIndex := flag.Int("gpu", 0, "Target GPU index (for multi-GPU systems)")
	dryRun := flag.Bool("dry-run", false, "Log intended power limit changes without applying them")
	logStdout := flag.Bool("log-stdout", true, "Enable logging to stdout")
	logFilePath := flag.String("log-file", "gpu_leaf.log", "Path to log file (set empty to disable file logging)")
	showVersion := flag.Bool("version", false, "Show version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Printf("gpu-leaf version: %s\n", version.Version)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nShutting down...")
		cancel()
	}()

	gpuPowerInfo, err := gpu.GetPowerInfo(*gpuIndex)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("GPU[%d] Power Info: %+v\n", *gpuIndex, gpuPowerInfo)

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
		GPUIndex:      *gpuIndex,
		DryRun:        *dryRun,
	}

	logger, logFile := utils.CreateLoggerWithPath(*logStdout, *logFilePath)
	if logFile != nil {
		defer logFile.Close()
	}

	ticker := time.NewTicker(time.Duration(*interval) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := gpu.WatchStats(ctx, &gpuConfig, logger, gpu.Leaf); err != nil {
				log.Fatalf("Error receiving GPU stats: %v", err)
			}
		}
	}
}
