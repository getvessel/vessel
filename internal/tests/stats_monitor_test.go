package tests

import (
	"math"
	"testing"

	"github.com/docker/docker/api/types"
	"vessel.dev/vessel/internal/orchestrator"
)

func TestCalculateCPUPercentage(t *testing.T) {
	tests := []struct {
		name  string
		stats *types.StatsJSON
		want  float64
	}{
		{
			name: "zero delta returns zero",
			stats: &types.StatsJSON{
				Stats: types.Stats{
					CPUStats: types.CPUStats{
						CPUUsage: types.CPUUsage{
							TotalUsage: 100,
						},
						SystemUsage: 200,
					},
					PreCPUStats: types.CPUStats{
						CPUUsage: types.CPUUsage{
							TotalUsage: 100,
						},
						SystemUsage: 200,
					},
				},
			},
			want: 0.0,
		},
		{
			name: "negative delta returns zero",
			stats: &types.StatsJSON{
				Stats: types.Stats{
					CPUStats: types.CPUStats{
						CPUUsage: types.CPUUsage{
							TotalUsage: 50,
						},
						SystemUsage: 100,
					},
					PreCPUStats: types.CPUStats{
						CPUUsage: types.CPUUsage{
							TotalUsage: 100,
						},
						SystemUsage: 200,
					},
				},
			},
			want: 0.0,
		},
		{
			name: "online CPUs specified",
			stats: &types.StatsJSON{
				Stats: types.Stats{
					CPUStats: types.CPUStats{
						CPUUsage: types.CPUUsage{
							TotalUsage: 200,
						},
						SystemUsage: 1000,
						OnlineCPUs:  4,
					},
					PreCPUStats: types.CPUStats{
						CPUUsage: types.CPUUsage{
							TotalUsage: 100,
						},
						SystemUsage: 500,
					},
				},
			},
			want: 80.0,
		},
		{
			name: "fallback to PercpuUsage length when OnlineCPUs is zero",
			stats: &types.StatsJSON{
				Stats: types.Stats{
					CPUStats: types.CPUStats{
						CPUUsage: types.CPUUsage{
							TotalUsage:  200,
							PercpuUsage: []uint64{100, 100},
						},
						SystemUsage: 1000,
						OnlineCPUs:  0,
					},
					PreCPUStats: types.CPUStats{
						CPUUsage: types.CPUUsage{
							TotalUsage: 100,
						},
						SystemUsage: 500,
					},
				},
			},
			want: 40.0,
		},
		{
			name: "fallback to 1 core when OnlineCPUs is zero and PercpuUsage is empty",
			stats: &types.StatsJSON{
				Stats: types.Stats{
					CPUStats: types.CPUStats{
						CPUUsage: types.CPUUsage{
							TotalUsage: 200,
						},
						SystemUsage: 1000,
						OnlineCPUs:  0,
					},
					PreCPUStats: types.CPUStats{
						CPUUsage: types.CPUUsage{
							TotalUsage: 100,
						},
						SystemUsage: 500,
					},
				},
			},
			want: 20.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := orchestrator.CalculateCPUPercentage(tt.stats)
			if math.Abs(got-tt.want) > 1e-6 {
				t.Errorf("CalculateCPUPercentage() = %v, want %v", got, tt.want)
			}
		})
	}
}
