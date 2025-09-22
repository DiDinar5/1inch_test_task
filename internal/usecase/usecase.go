package usecase

import (
	"github.com/DiDinar5/1inch_test_task/domain"
)

func NewUsecase(ethereumService domain.EthereumServiceInterface) domain.UsecaseInterface {
	return NewEstimateUsecase(ethereumService)
}
