package utils

import (
	"fmt"
	"strings"
)

func MegaBytesToBytes(mb int) int64 {
	if mb <= 0 {
		return 512 * 1024 * 1024
	}
	return int64(mb) * 1024 * 1024
}

func CPURequestToNanoCPUs(cores float64) int64 {
	if cores <= 0 {
		return 500_000_000
	}
	return int64(cores * 1_000_000_000)
}

func NormalizeContainerName(projectID string) string {
	return fmt.Sprintf("vessel-%s", strings.ToLower(strings.TrimSpace(projectID)))
}
