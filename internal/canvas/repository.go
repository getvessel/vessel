package canvas

import "context"

// Repository defines data operations for the canvas read-model domain.
type Repository interface {
	ListCanvasSummaries(ctx context.Context) ([]CanvasSummary, error)
	GetCanvasSummary(ctx context.Context, id string) (*CanvasSummary, error)
	GetEnvironmentCanvas(ctx context.Context, id string) (*EnvironmentCanvas, error)
}
