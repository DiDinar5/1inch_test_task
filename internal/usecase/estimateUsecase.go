package usecase

import (
	"context"

	"github.com/DiDinar5/1inch_test_task/domain"
)

type EstimateUsecase struct {
}

func NewEstimateUsecase() *EstimateUsecase {
	return &EstimateUsecase{}
}

func (u *EstimateUsecase) Estimate(ctx context.Context, req domain.EstimateRequest) (domain.EstimateResponse, error) {
	return domain.EstimateResponse{}, nil
}
