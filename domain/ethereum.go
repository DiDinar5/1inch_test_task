package domain

import (
	"math/big"
)

type PoolReserves struct {
	Reserve0    *big.Int `json:"reserve0"`
	Reserve1    *big.Int `json:"reserve1"`
	BlockNumber uint64   `json:"block_number"`
}

type TokenInfo struct {
	Address  string `json:"address"`
	Symbol   string `json:"symbol"`
	Decimals uint8  `json:"decimals"`
}
