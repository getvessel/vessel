package canvas

import "context"

type Repository interface {
	ListCanvasSummaries(ctx context.Context) ([]CanvasSummary, error)
	GetCanvasSummary(ctx context.Context, id string) (*CanvasSummary, error)
	GetEnvironmentCanvas(ctx context.Context, id string) (*EnvironmentCanvas, error)
}
