package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/DiDinar5/1inch_test_task/domain"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumService struct {
	client           *ethclient.Client
	uniswapV2ABI     abi.ABI
	erc20ABI         abi.ABI
	tokenAddresses   map[string]string
	tokenAddressesMu sync.RWMutex
	abiInitOnce      sync.Once
}

const uniswapV2PairABI = `[
	{
		"inputs": [],
		"name": "getReserves",
		"outputs": [
			{"internalType": "uint112", "name": "_reserve0", "type": "uint112"},
			{"internalType": "uint112", "name": "_reserve1", "type": "uint112"},
			{"internalType": "uint32", "name": "_blockTimestampLast", "type": "uint32"}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "token0",
		"outputs": [{"internalType": "address", "name": "", "type": "address"}],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "token1",
		"outputs": [{"internalType": "address", "name": "", "type": "address"}],
		"stateMutability": "view",
		"type": "function"
	}
]`

const erc20ABI = `[
	{
		"inputs": [],
		"name": "symbol",
		"outputs": [{"internalType": "string", "name": "", "type": "string"}],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "decimals",
		"outputs": [{"internalType": "uint8", "name": "", "type": "uint8"}],
		"stateMutability": "view",
		"type": "function"
	}
]`

func NewEthereumService(rpcURL string) (*EthereumService, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum: %w", err)
	}

	service := &EthereumService{
		client:         client,
		tokenAddresses: make(map[string]string),
	}

	service.abiInitOnce.Do(func() {
		service.initABI()
	})

	return service, nil
}

func (e *EthereumService) initABI() {
	var err error

	e.uniswapV2ABI, err = abi.JSON(strings.NewReader(uniswapV2PairABI))
	if err != nil {
		panic(fmt.Sprintf("failed to parse Uniswap V2 ABI: %v", err))
	}

	e.erc20ABI, err = abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		panic(fmt.Sprintf("failed to parse ERC20 ABI: %v", err))
	}
}

func (e *EthereumService) GetPoolReserves(ctx context.Context, poolAddress string) (*domain.PoolReserves, error) {
	if !common.IsHexAddress(poolAddress) {
		return nil, fmt.Errorf("invalid pool address: %s", poolAddress)
	}

	poolContract := common.HexToAddress(poolAddress)

	e.tokenAddressesMu.RLock()
	cachedAddresses, exists := e.tokenAddresses[poolAddress]
	e.tokenAddressesMu.RUnlock()

	var token0Address, token1Address common.Address
	var err error

	if exists {
		addresses := strings.Split(cachedAddresses, ",")
		if len(addresses) == 2 {
			token0Address = common.HexToAddress(addresses[0])
			token1Address = common.HexToAddress(addresses[1])
		}
	}

	if token0Address == (common.Address{}) || token1Address == (common.Address{}) {
		token0Data, err := e.callContract(ctx, poolContract, e.uniswapV2ABI, "token0")
		if err != nil {
			return nil, fmt.Errorf("failed to get token0 address: %w", err)
		}

		token1Data, err := e.callContract(ctx, poolContract, e.uniswapV2ABI, "token1")
		if err != nil {
			return nil, fmt.Errorf("failed to get token1 address: %w", err)
		}

		if err := e.uniswapV2ABI.UnpackIntoInterface(&token0Address, "token0", token0Data); err != nil {
			return nil, fmt.Errorf("failed to unpack token0 address: %w", err)
		}

		if err := e.uniswapV2ABI.UnpackIntoInterface(&token1Address, "token1", token1Data); err != nil {
			return nil, fmt.Errorf("failed to unpack token1 address: %w", err)
		}

		e.tokenAddressesMu.Lock()
		e.tokenAddresses[poolAddress] = token0Address.Hex() + "," + token1Address.Hex()
		e.tokenAddressesMu.Unlock()
	}

	reservesData, err := e.callContract(ctx, poolContract, e.uniswapV2ABI, "getReserves")
	if err != nil {
		return nil, fmt.Errorf("failed to get pool reserves: %w", err)
	}

	var reserves struct {
		Reserve0           *big.Int
		Reserve1           *big.Int
		BlockTimestampLast uint32
	}

	if err := e.uniswapV2ABI.UnpackIntoInterface(&reserves, "getReserves", reservesData); err != nil {
		return nil, fmt.Errorf("failed to unpack reserves data: %w", err)
	}

	blockNumber, err := e.client.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current block number: %w", err)
	}

	return &domain.PoolReserves{
		Reserve0:    reserves.Reserve0,
		Reserve1:    reserves.Reserve1,
		Token0:      token0Address.Hex(),
		Token1:      token1Address.Hex(),
		BlockNumber: blockNumber,
	}, nil
}

func (e *EthereumService) GetTokenInfo(ctx context.Context, tokenAddress string) (*domain.TokenInfo, error) {
	if !common.IsHexAddress(tokenAddress) {
		return nil, fmt.Errorf("invalid token address: %s", tokenAddress)
	}

	tokenContract := common.HexToAddress(tokenAddress)

	symbolData, err := e.callContract(ctx, tokenContract, e.erc20ABI, "symbol")
	if err != nil {
		return nil, fmt.Errorf("failed to get token symbol: %w", err)
	}

	decimalsData, err := e.callContract(ctx, tokenContract, e.erc20ABI, "decimals")
	if err != nil {
		return nil, fmt.Errorf("failed to get token decimals: %w", err)
	}

	var symbol string
	var decimals uint8

	if err := e.erc20ABI.UnpackIntoInterface(&symbol, "symbol", symbolData); err != nil {
		return nil, fmt.Errorf("failed to unpack symbol: %w", err)
	}

	if err := e.erc20ABI.UnpackIntoInterface(&decimals, "decimals", decimalsData); err != nil {
		return nil, fmt.Errorf("failed to unpack decimals: %w", err)
	}

	return &domain.TokenInfo{
		Address:  tokenAddress,
		Symbol:   symbol,
		Decimals: decimals,
	}, nil
}

func (e *EthereumService) callContract(ctx context.Context, contract common.Address, parsedABI abi.ABI, method string) ([]byte, error) {
	data, err := parsedABI.Pack(method)
	if err != nil {
		return nil, fmt.Errorf("failed to pack method %s: %w", method, err)
	}

	result, err := e.client.CallContract(ctx, ethereum.CallMsg{
		To:   &contract,
		Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract method %s: %w", method, err)
	}

	return result, nil
}
