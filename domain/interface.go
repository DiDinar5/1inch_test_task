package domain

import "context"

type UsecaseInterface interface {
	Estimate(ctx context.Context, req EstimateRequest) (EstimateResponse, error)
}
