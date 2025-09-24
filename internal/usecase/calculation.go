package usecase

import (
	"fmt"
	"math/big"
	"strings"
)

func (u *EstimateUsecase) parseAmount(amountStr string) (*big.Int, error) {
	amountStr = strings.TrimSpace(amountStr)

	if amountStr == "" {
		return nil, fmt.Errorf("amount cannot be empty")
	}

	amount, ok := new(big.Int).SetString(amountStr, 10)
	if !ok {
		return nil, fmt.Errorf("invalid amount format: %s", amountStr)
	}

	if amount.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	return amount, nil
}

func (u *EstimateUsecase) calculateAMMOutput(input, reserveIn, reserveOut *big.Int) (*big.Int, error) {
	if reserveIn.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("invalid reserve in: must be positive")
	}
	if reserveOut.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("invalid reserve out: must be positive")
	}

	if input.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("input amount must be positive")
	}

	feeNumerator := big.NewInt(UniswapV2FeeNumerator)
	feeDenominator := big.NewInt(UniswapV2FeeDenominator)

	inputWithFee := new(big.Int).Mul(input, feeNumerator)

	numerator := new(big.Int).Mul(inputWithFee, reserveOut)

	reserveInWithFee := new(big.Int).Mul(reserveIn, feeDenominator)

	denominator := new(big.Int).Add(reserveInWithFee, inputWithFee)

	if denominator.Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf("division by zero in AMM calculation")
	}

	output := new(big.Int).Div(numerator, denominator)

	if output.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("calculated output amount is not positive")
	}

	return output, nil
}
