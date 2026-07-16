package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type MetricsService struct {
	tsdbURL    string
	httpClient *http.Client
}

func NewMetricsService() *MetricsService {
	return &MetricsService{
		tsdbURL:    "http://127.0.0.1:8428",
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// GetServiceMetrics queries the TSDB for a given serviceID within the specified time range.
func (s *MetricsService) GetServiceMetrics(ctx context.Context, serviceID string, start, end time.Time, step string) (map[string]any, error) {
	// e.g., CPU query
	cpuQuery := fmt.Sprintf(`container_cpu_usage_percent{service_id="%s"}`, serviceID)
	// Memory query
	memQuery := fmt.Sprintf(`container_memory_usage_bytes{service_id="%s"}`, serviceID)

	cpuData, err := s.queryRange(ctx, cpuQuery, start, end, step)
	if err != nil {
		return nil, err
	}

	memData, err := s.queryRange(ctx, memQuery, start, end, step)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"cpu":    cpuData,
		"memory": memData,
	}, nil
}

func (s *MetricsService) queryRange(ctx context.Context, query string, start, end time.Time, step string) (any, error) {
	u, _ := url.Parse(s.tsdbURL + "/api/v1/query_range")
	q := u.Query()
	q.Set("query", query)
	q.Set("start", fmt.Sprintf("%d", start.Unix()))
	q.Set("end", fmt.Sprintf("%d", end.Unix()))
	q.Set("step", step)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("tsdb returned status %d", resp.StatusCode)
	}

	var res struct {
		Data any `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res.Data, nil
}
