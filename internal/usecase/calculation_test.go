package usecase

import (
	"math/big"
	"testing"
)

func bigIntFromString(s string) *big.Int {
	result, _ := new(big.Int).SetString(s, 10)
	return result
}

func TestCalculateAMMOutput(t *testing.T) {
	usecase := &EstimateUsecase{}

	tests := []struct {
		name        string
		input       *big.Int
		reserveIn   *big.Int
		reserveOut  *big.Int
		expected    *big.Int
		expectError bool
	}{
		{
			name:        "Normal swap 1 ETH",
			input:       big.NewInt(1000000000000000000),
			reserveIn:   bigIntFromString("10000000000000000000"),
			reserveOut:  bigIntFromString("20000000000000000000"),
			expected:    bigIntFromString("1813221787760298263"),
			expectError: false,
		},
		{
			name:        "Small swap",
			input:       big.NewInt(100000000000000000),
			reserveIn:   bigIntFromString("10000000000000000000"),
			reserveOut:  bigIntFromString("20000000000000000000"),
			expected:    bigIntFromString("197431606879412259"),
			expectError: false,
		},
		{
			name:        "Large swap",
			input:       bigIntFromString("10000000000000000000"),
			reserveIn:   bigIntFromString("10000000000000000000"),
			reserveOut:  bigIntFromString("20000000000000000000"),
			expected:    bigIntFromString("9984977466199298948"),
			expectError: false,
		},
		{
			name:        "Zero input",
			input:       big.NewInt(0),
			reserveIn:   bigIntFromString("10000000000000000000"),
			reserveOut:  bigIntFromString("20000000000000000000"),
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Zero reserve in",
			input:       big.NewInt(1000000000000000000),
			reserveIn:   big.NewInt(0),
			reserveOut:  bigIntFromString("20000000000000000000"),
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Zero reserve out",
			input:       big.NewInt(1000000000000000000),
			reserveIn:   bigIntFromString("10000000000000000000"),
			reserveOut:  big.NewInt(0),
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Very small reserves",
			input:       big.NewInt(1000),
			reserveIn:   big.NewInt(10000),
			reserveOut:  big.NewInt(20000),
			expected:    big.NewInt(1813),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := usecase.calculateAMMOutput(tt.input, tt.reserveIn, tt.reserveOut)

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

			if result.Cmp(tt.expected) != 0 {
				t.Errorf("Expected %s, got %s", tt.expected.String(), result.String())
			}
		})
	}
}

func TestParseAmount(t *testing.T) {
	usecase := &EstimateUsecase{}

	tests := []struct {
		name        string
		amountStr   string
		expected    *big.Int
		expectError bool
	}{
		{
			name:        "Valid amount",
			amountStr:   "1000000000000000000",
			expected:    big.NewInt(1000000000000000000),
			expectError: false,
		},
		{
			name:        "Valid amount with spaces",
			amountStr:   " 1000000000000000000 ",
			expected:    big.NewInt(1000000000000000000),
			expectError: false,
		},
		{
			name:        "Empty string",
			amountStr:   "",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Invalid format",
			amountStr:   "invalid",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Zero amount",
			amountStr:   "0",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Negative amount",
			amountStr:   "-1000",
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := usecase.parseAmount(tt.amountStr)

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

			if result.Cmp(tt.expected) != 0 {
				t.Errorf("Expected %s, got %s", tt.expected.String(), result.String())
			}
		})
	}
}

func BenchmarkCalculateAMMOutput(b *testing.B) {
	usecase := &EstimateUsecase{}
	input := big.NewInt(1000000000000000000)
	reserveIn := bigIntFromString("10000000000000000000")
	reserveOut := bigIntFromString("20000000000000000000")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := usecase.calculateAMMOutput(input, reserveIn, reserveOut)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCalculateAMMOutputBatch(b *testing.B) {
	usecase := &EstimateUsecase{}

	input := big.NewInt(1_000_000_000_000_000_000)
	reserveIn := bigIntFromString("10000000000000000000")
	reserveOut := bigIntFromString("20000000000000000000")

	const batchSize = 1000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < batchSize; j++ {
			_, err := usecase.calculateAMMOutput(input, reserveIn, reserveOut)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkParseAmount(b *testing.B) {
	usecase := &EstimateUsecase{}
	amountStr := "1000000000000000000"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := usecase.parseAmount(amountStr)
		if err != nil {
			b.Fatal(err)
		}
	}
}
