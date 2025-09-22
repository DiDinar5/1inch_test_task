package ethereum

import (
	"context"
	"fmt"

	"github.com/DiDinar5/1inch_test_task/domain"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumService struct {
	client *ethclient.Client
}

func NewEthereumService(rpcURL string) (*EthereumService, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum: %w", err)
	}

	return &EthereumService{
		client: client,
	}, nil
}

func (e *EthereumService) GetPoolReserves(ctx context.Context, poolAddress string) (*domain.PoolReserves, error) {
	return &domain.PoolReserves{}, nil
}

func (e *EthereumService) GetTokenInfo(ctx context.Context, tokenAddress string) (*domain.TokenInfo, error) {
	return &domain.TokenInfo{}, nil
}
