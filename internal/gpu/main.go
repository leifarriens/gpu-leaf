package gpu

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/leifarriens/gpu-leaf/internal/utils"
)

type GPUPowerInfo struct {
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
	cmd := exec.Command("nvidia-smi", "--query-gpu=power.min_limit,power.max_limit,power.default_limit", "--format=csv,noheader,nounits")

	stdout, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	output := strings.TrimSpace(string(stdout))

	fields := strings.Split(output, ",")

	if len(fields) != 3 {
		return nil, fmt.Errorf("unexpected output format")
	}

	minPowerLimit := utils.ParseFloat(fields[0])

	maxPowerLimit := utils.ParseFloat(fields[1])

	defaultPowerLimit := utils.ParseFloat(fields[2])

	gpuInfo := &GPUPowerInfo{
		MinPowerLimit:     minPowerLimit,
		MaxPowerLimit:     maxPowerLimit,
		DefaultPowerLimit: defaultPowerLimit,
	}

	return gpuInfo, nil
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
