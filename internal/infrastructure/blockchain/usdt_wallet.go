package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

// USDT ABI (ERC20 Token)
const usdtABI = `[
	{
		"constant": true,
		"inputs": [{"name": "_owner", "type": "address"}],
		"name": "balanceOf",
		"outputs": [{"name": "balance", "type": "uint256"}],
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{"name": "_to", "type": "address"},
			{"name": "_value", "type": "uint256"}
		],
		"name": "transfer",
		"outputs": [{"name": "", "type": "bool"}],
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{"name": "_spender", "type": "address"},
			{"name": "_value", "type": "uint256"}
		],
		"name": "approve",
		"outputs": [{"name": "", "type": "bool"}],
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{"name": "_owner", "type": "address"},
			{"name": "_spender", "type": "address"}
		],
		"name": "allowance",
		"outputs": [{"name": "", "type": "uint256"}],
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{"name": "_from", "type": "address"},
			{"name": "_to", "type": "address"},
			{"name": "_value", "type": "uint256"}
		],
		"name": "transferFrom",
		"outputs": [{"name": "", "type": "bool"}],
		"type": "function"
	}
]`

// USDTWalletClient handles USDT token operations
type USDTWalletClient struct {
	client       *ethclient.Client
	privateKey   *ecdsa.PrivateKey
	fromAddr     common.Address
	usdtContract common.Address
	chainID      *big.Int
	auth         *bind.TransactOpts
}

// NewUSDTWalletClient creates a new USDT wallet client
func NewUSDTWalletClient(rawUrl string, pkHex string, usdtContractAddr string) (*USDTWalletClient, error) {
	client, err := ethclient.Dial(rawUrl)
	if err != nil {
		return nil, fmt.Errorf("dial eth: %w", err)
	}

	pk, err := crypto.HexToECDSA(strings.TrimPrefix(pkHex, "0x"))
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get chain id: %w", err)
	}

	publicKey := pk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid public key")
	}
	fromAddr := crypto.PubkeyToAddress(*publicKeyECDSA)

	auth, err := bind.NewKeyedTransactorWithChainID(pk, chainID)
	if err != nil {
		return nil, fmt.Errorf("create transactor: %w", err)
	}

	// Set reasonable gas price and limit
	auth.GasLimit = 100000                       // Standard ERC20 transfer
	auth.GasPrice = big.NewInt(params.GWei * 20) // 20 Gwei

	return &USDTWalletClient{
		client:       client,
		privateKey:   pk,
		fromAddr:     fromAddr,
		usdtContract: common.HexToAddress(usdtContractAddr),
		chainID:      chainID,
		auth:         auth,
	}, nil
}

// GetBalance returns the USDT balance of the wallet
func (c *USDTWalletClient) GetBalance(ctx context.Context) (*big.Int, error) {
	parsedABI, err := abi.JSON(strings.NewReader(usdtABI))
	if err != nil {
		return nil, fmt.Errorf("parse USDT ABI: %w", err)
	}

	boundContract := bind.NewBoundContract(c.usdtContract, parsedABI, c.client, c.client, c.client)

	var results []interface{}
	err = boundContract.Call(nil, &results, "balanceOf", c.fromAddr)
	if err != nil {
		return nil, fmt.Errorf("call balanceOf: %w", err)
	}

	if len(results) == 0 {
		return big.NewInt(0), nil
	}

	balance, ok := results[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid balance type")
	}

	return balance, nil
}

// Transfer sends USDT to another address
func (c *USDTWalletClient) Transfer(ctx context.Context, toAddress string, amount *big.Int) (string, error) {
	parsedABI, err := abi.JSON(strings.NewReader(usdtABI))
	if err != nil {
		return "", fmt.Errorf("parse USDT ABI: %w", err)
	}

	boundContract := bind.NewBoundContract(c.usdtContract, parsedABI, c.client, c.client, c.client)

	toAddr := common.HexToAddress(toAddress)

	auth := *c.auth
	auth.Context = ctx
	auth.Value = big.NewInt(0) // No ETH for ERC20 transfer

	tx, err := boundContract.Transact(&auth, "transfer", toAddr, amount)
	if err != nil {
		return "", fmt.Errorf("transfer USDT: %w", err)
	}

	if tx == nil {
		return "", fmt.Errorf("transaction was nil")
	}

	return tx.Hash().Hex(), nil
}

// Approve allows a spender to use USDT from this wallet
func (c *USDTWalletClient) Approve(ctx context.Context, spender string, amount *big.Int) (string, error) {
	parsedABI, err := abi.JSON(strings.NewReader(usdtABI))
	if err != nil {
		return "", fmt.Errorf("parse USDT ABI: %w", err)
	}

	boundContract := bind.NewBoundContract(c.usdtContract, parsedABI, c.client, c.client, c.client)

	spenderAddr := common.HexToAddress(spender)

	auth := *c.auth
	auth.Context = ctx
	auth.Value = big.NewInt(0)

	tx, err := boundContract.Transact(&auth, "approve", spenderAddr, amount)
	if err != nil {
		return "", fmt.Errorf("approve USDT: %w", err)
	}

	if tx == nil {
		return "", fmt.Errorf("transaction was nil")
	}

	return tx.Hash().Hex(), nil
}

// GetAllowance returns the allowance for a spender
func (c *USDTWalletClient) GetAllowance(ctx context.Context, spender string) (*big.Int, error) {
	parsedABI, err := abi.JSON(strings.NewReader(usdtABI))
	if err != nil {
		return nil, fmt.Errorf("parse USDT ABI: %w", err)
	}

	boundContract := bind.NewBoundContract(c.usdtContract, parsedABI, c.client, c.client, c.client)

	spenderAddr := common.HexToAddress(spender)

	var results []interface{}
	err = boundContract.Call(nil, &results, "allowance", c.fromAddr, spenderAddr)
	if err != nil {
		return nil, fmt.Errorf("call allowance: %w", err)
	}

	if len(results) == 0 {
		return big.NewInt(0), nil
	}

	allowance, ok := results[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid allowance type")
	}

	return allowance, nil
}

// TransferFrom transfers USDT from one address to another (with approval)
func (c *USDTWalletClient) TransferFrom(ctx context.Context, from string, to string, amount *big.Int) (string, error) {
	parsedABI, err := abi.JSON(strings.NewReader(usdtABI))
	if err != nil {
		return "", fmt.Errorf("parse USDT ABI: %w", err)
	}

	boundContract := bind.NewBoundContract(c.usdtContract, parsedABI, c.client, c.client, c.client)

	fromAddr := common.HexToAddress(from)
	toAddr := common.HexToAddress(to)

	auth := *c.auth
	auth.Context = ctx
	auth.Value = big.NewInt(0)

	tx, err := boundContract.Transact(&auth, "transferFrom", fromAddr, toAddr, amount)
	if err != nil {
		return "", fmt.Errorf("transferFrom USDT: %w", err)
	}

	if tx == nil {
		return "", fmt.Errorf("transaction was nil")
	}

	return tx.Hash().Hex(), nil
}

// GetAddress returns the wallet address
func (c *USDTWalletClient) GetAddress() string {
	return c.fromAddr.Hex()
}

// EstimateTransferGas estimates gas cost for a transfer
func (c *USDTWalletClient) EstimateTransferGas(ctx context.Context, toAddress string, amount *big.Int) (*big.Int, error) {
	// Standard ERC20 transfer gas limit
	return big.NewInt(100000), nil
}

// WaitForTransaction waits for a transaction to be mined
func (c *USDTWalletClient) WaitForTransaction(ctx context.Context, txHash string) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, c.client, &types.Transaction{})
	if err != nil {
		return nil, fmt.Errorf("wait for transaction: %w", err)
	}

	return receipt, nil
}

// GetTransactionStatus checks if a transaction was successful
func (c *USDTWalletClient) GetTransactionStatus(ctx context.Context, txHash string) (bool, error) {
	receipt, err := c.client.TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		return false, fmt.Errorf("get transaction receipt: %w", err)
	}

	if receipt == nil {
		return false, fmt.Errorf("transaction not found")
	}

	return receipt.Status == 1, nil
}

// USDT decimals (USDT has 6 decimals)
const USDTDecimals = 6

// ToUSDTAmount converts float amount to USDT smallest unit
func ToUSDTAmount(amount float64) *big.Int {
	// Convert to smallest unit (6 decimals)
	amountStr := fmt.Sprintf("%.6f", amount)
	amountBig, _ := new(big.Float).SetString(amountStr)

	multiplier := big.NewFloat(1e6)
	result := new(big.Float).Mul(amountBig, multiplier)

	// Convert to big.Int
	amountInt := new(big.Int)
	result.Int(amountInt)

	return amountInt
}

// FromUSDTAmount converts USDT smallest unit to float amount
func FromUSDTAmount(amount *big.Int) float64 {
	if amount == nil {
		return 0
	}

	amountFloat := new(big.Float).SetInt(amount)
	divisor := big.NewFloat(1e6)
	result := new(big.Float).Quo(amountFloat, divisor)

	resultFloat, _ := result.Float64()
	return resultFloat
}
