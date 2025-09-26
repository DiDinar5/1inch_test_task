package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/DiDinar5/1inch_test_task/domain"
)

type mockEthereumService struct {
	poolReserves *domain.PoolReserves
	error        error
}

func (m *mockEthereumService) GetPoolReserves(ctx context.Context, poolAddress string) (*domain.PoolReserves, error) {
	return m.poolReserves, m.error
}

func TestEstimate(t *testing.T) {
	tests := []struct {
		name           string
		request        domain.EstimateRequest
		mockReserves   *domain.PoolReserves
		mockError      error
		expectedAmount string
		expectError    bool
	}{
		{
			name: "Valid swap Token0 to Token1",
			request: domain.EstimateRequest{
				Pool:      "0x1234567890123456789012345678901234567890",
				Src:       "0x1111111111111111111111111111111111111111",
				Dst:       "0x2222222222222222222222222222222222222222",
				SrcAmount: "1000000000000000000",
			},
			mockReserves: &domain.PoolReserves{
				Reserve0:    bigIntFromString("10000000000000000000"),
				Reserve1:    bigIntFromString("20000000000000000000"),
				Token0:      "0x1111111111111111111111111111111111111111",
				Token1:      "0x2222222222222222222222222222222222222222",
				BlockNumber: 12345,
			},
			expectedAmount: "1813221787760298263",
			expectError:    false,
		},
		{
			name: "Valid swap Token1 to Token0",
			request: domain.EstimateRequest{
				Pool:      "0x1234567890123456789012345678901234567890",
				Src:       "0x2222222222222222222222222222222222222222",
				Dst:       "0x1111111111111111111111111111111111111111",
				SrcAmount: "2000000000000000000",
			},
			mockReserves: &domain.PoolReserves{
				Reserve0:    bigIntFromString("10000000000000000000"),
				Reserve1:    bigIntFromString("20000000000000000000"),
				Token0:      "0x1111111111111111111111111111111111111111",
				Token1:      "0x2222222222222222222222222222222222222222",
				BlockNumber: 12345,
			},
			expectedAmount: "906610893880149131",
			expectError:    false,
		},
		{
			name: "Ethereum service error",
			request: domain.EstimateRequest{
				Pool:      "0x1234567890123456789012345678901234567890",
				Src:       "0x1111111111111111111111111111111111111111",
				Dst:       "0x2222222222222222222222222222222222222222",
				SrcAmount: "1000000000000000000",
			},
			mockError:   errors.New("failed to get reserves"),
			expectError: true,
		},
		{
			name: "Invalid amount format",
			request: domain.EstimateRequest{
				Pool:      "0x1234567890123456789012345678901234567890",
				Src:       "0x1111111111111111111111111111111111111111",
				Dst:       "0x2222222222222222222222222222222222222222",
				SrcAmount: "invalid",
			},
			mockReserves: &domain.PoolReserves{
				Reserve0:    bigIntFromString("10000000000000000000"),
				Reserve1:    bigIntFromString("20000000000000000000"),
				Token0:      "0x1111111111111111111111111111111111111111",
				Token1:      "0x2222222222222222222222222222222222222222",
				BlockNumber: 12345,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockEthereumService{
				poolReserves: tt.mockReserves,
				error:        tt.mockError,
			}

			usecase := NewEstimateUsecase(mockService)
			result, err := usecase.Estimate(context.Background(), tt.request)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result.DstAmount != tt.expectedAmount {
				t.Errorf("Expected %s, got %s", tt.expectedAmount, result.DstAmount)
			}
		})
	}
}
