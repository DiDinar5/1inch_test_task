package usecase

import (
	"context"

	"github.com/DiDinar5/1inch_test_task/domain"
)

type EstimateUsecase struct {
	ethereumService domain.EthereumServiceInterface
}

func NewEstimateUsecase(ethereumService domain.EthereumServiceInterface) *EstimateUsecase {
	return &EstimateUsecase{
		ethereumService: ethereumService,
	}
}

func (u *EstimateUsecase) Estimate(ctx context.Context, req domain.EstimateRequest) (domain.EstimateResponse, error) {
	return domain.EstimateResponse{}, nil
}
