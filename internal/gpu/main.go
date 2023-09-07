package gpu

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/leifarriens/gpu-leaf/internal/utils"
)

type GPUPowerInfo struct {
	IsPowerManageable bool
	MinPowerLimit     float64
	MaxPowerLimit     float64
	DefaultPowerLimit float64
}

type GPUConfig struct {
	MinPowerLimit float64
	MaxPowerLimit float64
	Threshold     int
}

type GPUStats struct {
	Temperature float64
	PowerDraw   float64
	Utilization int
	PowerLimit  float64
}

func GetPowerInfo() (*GPUPowerInfo, error) {
	cmd := exec.Command("nvidia-smi", "--query-gpu=power.management,power.min_limit,power.max_limit,power.default_limit", "--format=csv,noheader,nounits")

	stdout, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	output := strings.TrimSpace(string(stdout))

	fields := strings.Split(output, ",")

	if len(fields) != 4 {
		return nil, fmt.Errorf("query gpu unexpected output format")
	}

	isPowerManageable := fields[0] == "Supported" || fields[0] == "Enabled"

	minPowerLimit := utils.ParseFloat(fields[1])

	maxPowerLimit := utils.ParseFloat(fields[2])

	defaultPowerLimit := utils.ParseFloat(fields[3])

	gpuInfo := &GPUPowerInfo{
		IsPowerManageable: isPowerManageable,
		MinPowerLimit:     minPowerLimit,
		MaxPowerLimit:     maxPowerLimit,
		DefaultPowerLimit: defaultPowerLimit,
	}

	return gpuInfo, nil
}

func WatchStats(gpuConfig *GPUConfig, logger *log.Logger, callback func(*GPUConfig, *GPUStats)) error {
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

			gpuStats := GPUStats{
				Temperature: temperature,
				PowerDraw:   powerDraw,
				Utilization: gpuUtil,
				PowerLimit:  powerLimit,
			}

			logger.Printf("GPU Info: %+v", gpuStats)
			callback(gpuConfig, &gpuStats)
		}
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func Leaf(gpuConfig *GPUConfig, gpuStats *GPUStats) {
	baseFactor := gpuConfig.MaxPowerLimit / 10

	var newValue float64

	if gpuStats.Utilization < gpuConfig.Threshold {
		utilDiff := (100 - float64(gpuStats.Utilization)) / 100
		newValue = gpuStats.PowerLimit - (baseFactor * utilDiff)
	}

	if gpuStats.Utilization >= gpuConfig.Threshold {
		newValue = gpuStats.PowerLimit + (baseFactor * (float64(gpuStats.Utilization) / 100))
	}

	if newValue > gpuConfig.MinPowerLimit && newValue < gpuConfig.MaxPowerLimit {
		setPowerLimit(newValue)
	} else {
		fmt.Println(newValue)
	}
}

func setPowerLimit(limit float64) {
	// s := fmt.Sprintf("%", limit)
	s := strconv.Itoa(int(limit))

	cmd := exec.Command("nvidia-smi", "-i", "0", "-pl", s)

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
