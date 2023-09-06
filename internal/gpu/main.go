package gpu

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/leifarriens/gpu-leaf/internal/utils"
)

type GPUPowerInfo struct {
	MinPowerLimit     float64
	MaxPowerLimit     float64
	DefaultPowerLimit float64
}

type GPUStats struct {
	Temperature float64
	PowerDraw   float64
	Utilization float64
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

func Leaf(gpuPowerInfo *GPUPowerInfo, gpuStats *GPUStats) {
	if gpuStats.PowerLimit > gpuPowerInfo.MinPowerLimit && gpuStats.PowerLimit <= gpuPowerInfo.MaxPowerLimit {
		if gpuStats.Utilization < 95 {
			s := fmt.Sprintf("%f", gpuStats.PowerLimit-1)
			cmd := exec.Command("nvidia-smi", "-i", "0", "-pl", s)

			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}
		}

		if gpuStats.Utilization > 95 {
			s := fmt.Sprintf("%f", gpuStats.PowerLimit+1)
			cmd := exec.Command("nvidia-smi", "-i", "0", "-pl", s)

			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}
		}

	}
}
