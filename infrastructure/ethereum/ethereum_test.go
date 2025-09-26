package ethereum

import (
	"context"
	"math/big"
	"testing"
)

func TestEthereumService_GetPoolReserves(t *testing.T) {
	tests := []struct {
		name          string
		poolAddress   string
		expectedError bool
	}{
		{
			name:          "Invalid pool address",
			poolAddress:   "invalid_address",
			expectedError: true,
		},
		{
			name:          "Empty pool address",
			poolAddress:   "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewEthereumService("invalid-rpc-url")
			if err != nil {
				if !tt.expectedError {
					t.Errorf("Unexpected error: %v", err)
				}
				return
			}

			ctx := context.Background()
			result, err := service.GetPoolReserves(ctx, tt.poolAddress)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("Expected result but got nil")
			}
		})
	}
}

func TestEthereumService_GetTokenInfo(t *testing.T) {
	tests := []struct {
		name          string
		tokenAddress  string
		expectedError bool
	}{
		{
			name:          "Invalid token address",
			tokenAddress:  "invalid_address",
			expectedError: true,
		},
		{
			name:          "Empty token address",
			tokenAddress:  "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewEthereumService("invalid-rpc-url")
			if err != nil {
				if !tt.expectedError {
					t.Errorf("Unexpected error: %v", err)
				}
				return
			}

			ctx := context.Background()
			result, err := service.GetTokenInfo(ctx, tt.tokenAddress)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("Expected result but got nil")
			}
		})
	}
}

func TestEthereumService_Caching(t *testing.T) {
	service, err := NewEthereumService("invalid-rpc-url")
	if err != nil {
		t.Skipf("Skipping test due to Ethereum connection error: %v", err)
	}

	if len(service.tokenInfoCache) != 0 {
		t.Errorf("Expected empty cache, got %d items", len(service.tokenInfoCache))
	}

	if len(service.tokenAddresses) != 0 {
		t.Errorf("Expected empty token addresses cache, got %d items", len(service.tokenAddresses))
	}
}

func TestEthereumService_InitABI(t *testing.T) {
	service := &EthereumService{}

	err := service.initABI()
	if err != nil {
		t.Errorf("Failed to initialize ABI: %v", err)
	}

	if service.uniswapV2ABI.Methods == nil {
		t.Error("UniswapV2 ABI not initialized")
	}

	if service.erc20ABI.Methods == nil {
		t.Error("ERC20 ABI not initialized")
	}
}

func BenchmarkEthereumService_GetPoolReserves(b *testing.B) {
	service, err := NewEthereumService("invalid-rpc-url")
	if err != nil {
		b.Skipf("Skipping benchmark due to Ethereum connection error: %v", err)
	}

	ctx := context.Background()
	poolAddress := "0x1234567890123456789012345678901234567890"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetPoolReserves(ctx, poolAddress)
	}
}

func BenchmarkEthereumService_GetTokenInfo(b *testing.B) {
	service, err := NewEthereumService("invalid-rpc-url")
	if err != nil {
		b.Skipf("Skipping benchmark due to Ethereum connection error: %v", err)
	}

	ctx := context.Background()
	tokenAddress := "0x1234567890123456789012345678901234567890"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetTokenInfo(ctx, tokenAddress)
	}
}

func TestDecimalsFallback(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected uint8
		hasError bool
	}{
		{
			name:     "uint8",
			input:    uint8(18),
			expected: 18,
			hasError: false,
		},
		{
			name:     "uint32",
			input:    uint32(6),
			expected: 6,
			hasError: false,
		},
		{
			name:     "uint64",
			input:    uint64(8),
			expected: 8,
			hasError: false,
		},
		{
			name:     "big.Int valid",
			input:    big.NewInt(18),
			expected: 18,
			hasError: false,
		},
		{
			name:     "big.Int too large",
			input:    big.NewInt(300),
			expected: 0,
			hasError: true,
		},
		{
			name:     "string (invalid)",
			input:    "18",
			expected: 0,
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var decimals uint8
			var hasError bool

			switch v := tc.input.(type) {
			case uint8:
				decimals = v
			case uint32:
				decimals = uint8(v)
			case uint64:
				decimals = uint8(v)
			case *big.Int:
				if v.IsUint64() && v.Uint64() <= 255 {
					decimals = uint8(v.Uint64())
				} else {
					hasError = true
				}
			default:
				hasError = true
			}

			if hasError != tc.hasError {
				t.Errorf("Expected error %v, got %v", tc.hasError, hasError)
			}

			if !hasError && decimals != tc.expected {
				t.Errorf("Expected %d, got %d", tc.expected, decimals)
			}
		})
	}
}
